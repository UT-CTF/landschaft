package baseline

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// compareCSVDirs compares two baseline directories. Each directory is expected to
// contain CSV files named services.csv, processes.csv, autoruns.csv, users.csv,
// adobjects.csv, ports.csv. The function will print added/removed/changed items
// per-file using component-specific primary keys and prints full objects for
// additions/removals.
func compareCSVDirs(dirA, dirB string) {
	dirA, _ = filepath.Abs(dirA)
	dirB, _ = filepath.Abs(dirB)

	files := []string{"services.csv", "processes.csv", "autoruns.csv", "users.csv", "adobjects.csv", "ports.csv"}

	for _, f := range files {
		pathA := filepath.Join(dirA, f)
		pathB := filepath.Join(dirB, f)

		keyCols := keyColumnsForFile(f)
		mA, errA := loadCSVWithKey(pathA, keyCols)
		mB, errB := loadCSVWithKey(pathB, keyCols)

		fmt.Println(strings.Repeat("=", 60))
		fmt.Printf("Comparing %s\n", f)

		if errA != nil {
			fmt.Printf("Could not load %s: %v\n", pathA, errA)
		}
		if errB != nil {
			fmt.Printf("Could not load %s: %v\n", pathB, errB)
		}
		if errA != nil || errB != nil {
			continue
		}

		addedKeys, removedKeys, changed := diffMaps(mA, mB)

		if len(addedKeys) > 0 {
			sort.Strings(addedKeys)
			fmt.Printf("Added in %s:\n", dirB)
			for _, k := range addedKeys {
				fmt.Printf("\t%s\n", formatObject(k, mB[k]))
			}
		}
		if len(removedKeys) > 0 {
			sort.Strings(removedKeys)
			fmt.Printf("Removed from %s:\n", dirB)
			for _, k := range removedKeys {
				fmt.Printf("\t%s\n", formatObject(k, mA[k]))
			}
		}
		if len(changed) > 0 {
			fmt.Println("Changed entries:")
			for _, c := range changed {
				fmt.Printf("\t%s\n", c)
			}
		}
	}
}

// keyColumnsForFile returns the list of CSV column names to be used as the
// primary key for a given file name.
func keyColumnsForFile(file string) []string {
	switch file {
	case "adobjects.csv":
		return []string{"DistinguishedName"}
	case "autoruns.csv":
		return []string{"Location", "Name", "LaunchString"}
	case "ports.csv":
		return []string{"LocalAddress", "LocalPort"}
	case "processes.csv":
		return []string{"Name", "Path"}
	case "services.csv":
		return []string{"Name"}
	case "users.csv":
		return []string{"SamAccountName"}
	default:
		return nil
	}
}

// loadCSVWithKey loads a CSV into a map keyed by the composite key defined by
// keyCols. If keyCols is nil or empty, the first column is used as the key.
func loadCSVWithKey(path string, keyCols []string) (map[string]map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	recs, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(recs) < 1 {
		return nil, fmt.Errorf("empty csv: %s", path)
	}
	headers := recs[0]
	out := make(map[string]map[string]string)

	buildKey := func(row []string) string {
		if len(keyCols) == 0 {
			if len(row) > 0 {
				return row[0]
			}
			return ""
		}
		vals := make([]string, 0, len(keyCols))
		for _, kc := range keyCols {
			idx := -1
			for i, h := range headers {
				if h == kc {
					idx = i
					break
				}
			}
			if idx >= 0 && idx < len(row) {
				vals = append(vals, row[idx])
			} else {
				vals = append(vals, "")
			}
		}
		return strings.Join(vals, "|")
	}

	for _, row := range recs[1:] {
		if len(row) == 0 {
			continue
		}
		m := make(map[string]string)
		for i, cell := range row {
			if i < len(headers) {
				m[headers[i]] = cell
			} else {
				m[fmt.Sprintf("col_%d", i)] = cell
			}
		}
		key := buildKey(row)
		out[key] = m
	}
	return out, nil
}

// formatObject returns a single-line representation of the object's fields in
// key:value pairs separated by ", ".
func formatObject(key string, obj map[string]string) string {
	if obj == nil {
		return key
	}
	// pretty multi-line output: key on first line, then each field on its own indented line
	lines := []string{key}
	keys := make([]string, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		lines = append(lines, fmt.Sprintf("\t\t%s: %s", k, obj[k]))
	}
	return strings.Join(lines, "\n")
}

// loadGenericCSV remains as a convenience wrapper using first-column key.
func loadGenericCSV(path string) (map[string]map[string]string, error) {
	return loadCSVWithKey(path, nil)
}

// diffMaps returns added keys (in b but not a), removed keys (in a but not b), and
// changed descriptions for keys present in both where values differ.
func diffMaps(a, b map[string]map[string]string) (added, removed []string, changed []string) {
	for k := range a {
		if _, ok := b[k]; !ok {
			removed = append(removed, k)
		}
	}
	for k := range b {
		if _, ok := a[k]; !ok {
			added = append(added, k)
		}
	}
	for k := range a {
		if vb, ok := b[k]; ok {
			va := a[k]
			for hk, hv := range va {
				if vb[hk] != hv {
					changed = append(changed, fmt.Sprintf("%s: %s changed from '%s' to '%s'", k, hk, hv, vb[hk]))
				}
			}
		}
	}
	return
}
