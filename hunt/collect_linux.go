package hunt

import (
	"bufio"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/UT-CTF/landschaft/util"
)

func getDetections(since time.Duration) ([]Detection, error) {
	return collectLinux(since)
}

// collectLinux reads auth.log / secure and journalctl, returns normalized detections.
func collectLinux(since time.Duration) ([]Detection, error) {
	var out []Detection
	host, _ := os.Hostname()

	// Try auth.log (Debian/Ubuntu) or secure (RHEL)
	for _, path := range []string{"/var/log/auth.log", "/var/log/secure"} {
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		sc := bufio.NewScanner(f)
		cutoff := time.Now().Add(-since)
		for sc.Scan() {
			line := sc.Text()
			t, msg := parseAuthLogLine(line)
			if t.Before(cutoff) {
				continue
			}
			d := Detection{
				Timestamp: t.UTC().Format(time.RFC3339),
				Host:      host,
				OS:        runtime.GOOS,
				Source:    path,
				Message:   msg,
			}
			TagSuspicious(&d)
			out = append(out, d)
		}
		f.Close()
		if len(out) > 0 {
			break
		}
	}

	// journalctl if available
	jout, err := util.RunScriptQuiet("hunt/events.sh")
	if err == nil && strings.TrimSpace(jout) != "" {
		// events.sh outputs one JSON object per line
		for _, line := range strings.Split(strings.TrimSpace(jout), "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			d := Detection{
				Host:   host,
				OS:     runtime.GOOS,
				Source: "journalctl",
				Message: line,
			}
			TagSuspicious(&d)
			out = append(out, d)
		}
	}
	return out, nil
}

var authLogRE = regexp.MustCompile(`^(\w{3}\s+\d+\s+[\d:]+)\s+.*?\s+(.*)$`)

func parseAuthLogLine(line string) (t time.Time, msg string) {
	msg = line
	loc := time.Local
	if m := authLogRE.FindStringSubmatch(line); len(m) >= 3 {
		parsed, err := time.ParseInLocation("Jan _2 15:04:05", m[1], loc)
		if err == nil {
			year := time.Now().Year()
			parsed = time.Date(year, parsed.Month(), parsed.Day(), parsed.Hour(), parsed.Minute(), parsed.Second(), 0, loc)
			return parsed, m[2]
		}
	}
	return time.Now(), line
}
