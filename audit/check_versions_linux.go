package audit

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func checkSoftwareVersions(opts Options) {
	eco, err := detectOSVEcosystem()
	if err != nil {
		fmt.Println("OSV: unable to detect distro ecosystem:", err)
		return
	}

	pkgs, err := listDpkgPackages()
	if err != nil {
		fmt.Println("OSV: unable to list installed packages:", err)
		return
	}

	// Curated set of packages that matter most in CCDC.
	targets := []string{
		"openssh-server",
		"openssh-client",
		"openssl",
		"sudo",
		"curl",
		"wget",
		"nginx",
		"apache2",
		"postgresql",
		"mysql-server",
		"mariadb-server",
		"bind9",
		"samba",
		"rsyslog",
	}

	var queries []osvQuery
	for _, name := range targets {
		ver, ok := pkgs[name]
		if !ok || strings.TrimSpace(ver) == "" {
			continue
		}
		queries = append(queries, osvQuery{
			Package: osvPackage{Ecosystem: eco, Name: name},
			Version: ver,
		})
	}

	if len(queries) == 0 {
		fmt.Printf("OSV: no matching packages found for ecosystem %q\n", eco)
		return
	}

	max := opts.MaxPackages
	if max <= 0 {
		max = DefaultOptions().MaxPackages
	}
	if len(queries) > max {
		queries = queries[:max]
	}

	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = DefaultOptions().Timeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	results, err := osvQueryBatch(ctx, queries)
	if err != nil {
		fmt.Println("OSV: query failed:", err)
		return
	}

	printedAny := false
	for i, r := range results {
		if len(r.Vulns) == 0 {
			continue
		}

		q := queries[i]
		if !printedAny {
			fmt.Printf("OSV: potential matches for %q packages (ecosystem: %s)\n", len(results), eco)
			printedAny = true
		}

		fmt.Printf("\n- %s %s\n", q.Package.Name, q.Version)
		for _, v := range r.Vulns {
			summary := strings.TrimSpace(firstLine(v.Summary))
			if summary == "" {
				summary = "No summary"
			}
			fmt.Printf("  - %s: %s\n", v.ID, summary)
		}
	}

	if !printedAny {
		fmt.Println("OSV: no vulnerabilities returned for selected packages (this may be a false negative depending on distro mapping).")
	}
}

func detectOSVEcosystem() (string, error) {
	content, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "", err
	}

	m := make(map[string]string)
	sc := bufio.NewScanner(strings.NewReader(string(content)))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		val := strings.TrimSpace(parts[1])
		val = strings.Trim(val, "\"")
		m[key] = val
	}

	id := strings.ToLower(m["id"])
	idLike := strings.ToLower(m["id_like"])

	switch {
	case id == "ubuntu" || strings.Contains(idLike, "ubuntu"):
		return "Ubuntu", nil
	case id == "debian" || strings.Contains(idLike, "debian"):
		return "Debian", nil
	case id == "alpine" || strings.Contains(idLike, "alpine"):
		return "Alpine", nil
	case id == "rocky" || strings.Contains(idLike, "rocky"):
		return "Rocky Linux", nil
	case id == "almalinux" || strings.Contains(idLike, "almalinux"):
		return "AlmaLinux", nil
	default:
		return "", fmt.Errorf("unsupported distro (ID=%q ID_LIKE=%q)", m["id"], m["id_like"])
	}
}

func listDpkgPackages() (map[string]string, error) {
	// Debian/Ubuntu package listing.
	cmd := exec.Command("dpkg-query", "-W", "-f=${Package}\t${Version}\n")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	pkgs := make(map[string]string)
	sc := bufio.NewScanner(strings.NewReader(string(out)))
	for sc.Scan() {
		line := sc.Text()
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			continue
		}
		name := strings.TrimSpace(parts[0])
		ver := strings.TrimSpace(parts[1])
		if name == "" || ver == "" {
			continue
		}
		pkgs[name] = ver
	}
	return pkgs, nil
}

func firstLine(s string) string {
	if idx := strings.IndexByte(s, '\n'); idx >= 0 {
		return s[:idx]
	}
	return s
}

