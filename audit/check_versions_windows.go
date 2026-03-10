package audit

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/UT-CTF/landschaft/embed"
	"github.com/UT-CTF/landschaft/util"
)

func checkSoftwareVersions(opts Options) {
	_ = opts

	out, err := embed.ExecuteScript("audit/software.ps1", false)
	if err != nil {
		fmt.Println("OSV: failed to collect software inventory:", err)
		return
	}

	type softwareRow struct {
		Name      string `json:"Name"`
		Version   string `json:"Version"`
		Publisher string `json:"Publisher"`
	}

	trimmed := strings.TrimSpace(out)
	if trimmed == "" {
		fmt.Println("OSV: no software inventory returned")
		return
	}

	var rows []softwareRow
	if err := json.Unmarshal([]byte(trimmed), &rows); err != nil {
		// Some PowerShell configurations emit a single object instead of an array.
		var single softwareRow
		if err2 := json.Unmarshal([]byte(trimmed), &single); err2 != nil {
			fmt.Println("OSV: failed to parse software inventory JSON:", err)
			return
		}
		rows = []softwareRow{single}
	}

	// Deduplicate and sort for stable output.
	seen := make(map[string]bool)
	deduped := make([]softwareRow, 0, len(rows))
	for _, r := range rows {
		key := strings.ToLower(strings.TrimSpace(r.Name)) + "\x00" + strings.TrimSpace(r.Version)
		if strings.TrimSpace(r.Name) == "" || strings.TrimSpace(r.Version) == "" {
			continue
		}
		if seen[key] {
			continue
		}
		seen[key] = true
		deduped = append(deduped, r)
	}

	sort.Slice(deduped, func(i, j int) bool {
		a, b := strings.ToLower(deduped[i].Name), strings.ToLower(deduped[j].Name)
		if a == b {
			return deduped[i].Version < deduped[j].Version
		}
		return a < b
	})

	fmt.Println(util.TitleColor.Render("Windows software inventory"))
	tableRows := make([][]string, 0, len(deduped))
	for _, r := range deduped {
		tableRows = append(tableRows, []string{r.Name, r.Version})
	}

	t := util.StyledTable().Headers("name", "version").Rows(tableRows...)
	fmt.Println(t.Render())

	fmt.Println("Note: OSV.dev does not reliably map arbitrary Windows installed software to an ecosystem/package name.")
	fmt.Println("      This section is inventory-only for now; Linux package checks use OSV directly.")
}

