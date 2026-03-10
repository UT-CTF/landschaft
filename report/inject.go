package report

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/UT-CTF/landschaft/util"
)

// InjectReport writes a Markdown report for CCDC injects from action log and optional triage TSV.
func InjectReport(actionLogPath, triagePath, outPath string) error {
	var buf strings.Builder

	buf.WriteString("# Landschaft inject report\n\n")
	buf.WriteString("Generated: " + time.Now().UTC().Format(time.RFC3339) + "\n\n")

	// Summary
	buf.WriteString("## Summary\n\n")
	hostname, _ := os.Hostname()
	buf.WriteString(fmt.Sprintf("- **Host:** %s\n", hostname))
	buf.WriteString(fmt.Sprintf("- **Action log:** %s\n", actionLogPath))
	buf.WriteString(fmt.Sprintf("- **Triage file:** %s\n\n", triagePath))

	// Actions taken
	if actionLogPath != "" {
		buf.WriteString("## Actions taken\n\n")
		entries, err := readActionLog(actionLogPath)
		if err != nil {
			buf.WriteString(fmt.Sprintf("*(Error reading action log: %v)*\n\n", err))
		} else {
			for _, e := range entries {
				buf.WriteString(fmt.Sprintf("- **%s** `%s %s` (exit %d, %d ms)\n", e.Timestamp, e.Command, strings.Join(e.Args, " "), e.ExitCode, e.DurationMs))
			}
			buf.WriteString("\n")
		}
	}

	// Current state (from triage)
	if triagePath != "" {
		buf.WriteString("## Current state snapshot\n\n")
		content, err := os.ReadFile(triagePath)
		if err != nil {
			buf.WriteString(fmt.Sprintf("*(Error reading triage file: %v)*\n\n", err))
		} else {
			lines := strings.Split(string(content), "\n")
			if len(lines) > 0 {
				buf.WriteString("```\n")
				for _, l := range lines {
					if len(l) > 200 {
						l = l[:200] + "..."
					}
					buf.WriteString(l + "\n")
				}
				buf.WriteString("```\n\n")
			}
		}
	}

	// Findings
	buf.WriteString("## Findings\n\n")
	buf.WriteString("- Review non-default services in triage output.\n")
	buf.WriteString("- Review action log for hardening steps applied.\n")
	buf.WriteString("- If present, review landschaft-detections.jsonl for suspicious activity timeline.\n\n")

	outDir := filepath.Dir(outPath)
	_ = os.MkdirAll(outDir, 0755)
	return os.WriteFile(outPath, []byte(buf.String()), 0644)
}

func readActionLog(path string) ([]util.ActionLogEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var entries []util.ActionLogEntry
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}
		var e util.ActionLogEntry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			continue
		}
		entries = append(entries, e)
	}
	return entries, sc.Err()
}
