package baseline

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

func setupCompareCmd(cmd *cobra.Command) {
	compareCmd := &cobra.Command{
		Use:   "compare",
		Short: "Compare baselines",
	}

	compareAllCmd := &cobra.Command{
		Use:   "all <baselineA> <baselineB>",
		Short: "Compare two baseline directories",
		Long:  "Compare two directories produced by 'baseline create all' and report added/removed/changed entries for each component.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			compareCSVDirs(args[0], args[1])
		},
	}
	compareCmd.AddCommand(compareAllCmd)

	for name := range baselineComponents {
		cmdCmp := &cobra.Command{
			Use:   name + " <baselineA> <baselineB>",
			Short: fmt.Sprintf("Compare %s baselines", name),
			Long:  fmt.Sprintf("Compare the %s.csv files in two baseline directories and report added/removed/changed entries.", name),
			Args:  cobra.ExactArgs(2),
			Run: func(cmd *cobra.Command, args []string) {
				err := compareCSVFiles(fmt.Sprintf("%s.csv", name), args[0], args[1])
				if err != nil {
					fmt.Printf("Error comparing %s: %v\n", name, err)
				}
			},
		}

		compareCmd.AddCommand(cmdCmp)
	}

	cmd.AddCommand(compareCmd)
}

func getAllCSVFiles(dir string) ([]string, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("could not get absolute path of %s: %v", dir, err)
	}
	files := []string{}
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("could not read directory %s: %v", dir, err)
	}
	for _, entry := range dirEntries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".csv") {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

func compareCSVFiles(fileName, dirA, dirB string) error {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	keyCols := keyColumnsForFile(fileName)
	pathA := filepath.Join(dirA, fileName)
	pathB := filepath.Join(dirB, fileName)
	mapA, err := loadCSVWithKey(pathA, keyCols)
	if err != nil {
		return fmt.Errorf("could not load %s: %v", pathA, err)
	}
	mapB, err := loadCSVWithKey(pathB, keyCols)
	if err != nil {
		return fmt.Errorf("could not load %s: %v", pathB, err)
	}

	fmt.Printf("Comparing %s\n", fileName)

	addedKeys, removedKeys, changed := diffMaps(mapA, mapB)

	if len(addedKeys) > 0 {
		sort.Strings(addedKeys)
		fmt.Printf("Added in %s:\n", dirB)
		for _, k := range addedKeys {
			fmt.Printf("\t%s\n", formatObject(k, mapB[k]))
		}
	}
	if len(removedKeys) > 0 {
		sort.Strings(removedKeys)
		fmt.Printf("Removed from %s:\n", dirB)
		for _, k := range removedKeys {
			fmt.Printf("\t%s\n", formatObject(k, mapA[k]))
		}
	}
	if len(changed) > 0 {
		fmt.Println("Changed entries:")
		for _, c := range changed {
			fmt.Printf("\t%s\n", c)
		}
	}
	return nil
}

func compareCSVDirs(dirA, dirB string) {
	filesA, err := getAllCSVFiles(dirA)
	if err != nil {
		fmt.Printf("Error getting CSV files from directory %s: %v\n", dirA, err)
		return
	}
	filesB, err := getAllCSVFiles(dirB)
	if err != nil {
		fmt.Printf("Error getting CSV files from directory %s: %v\n", dirB, err)
		return
	}

	sharedFiles := []string{}
	for _, fA := range filesA {
		if slices.Contains(filesB, fA) {
			sharedFiles = append(sharedFiles, fA)
		}
	}

	for _, f := range sharedFiles {
		err := compareCSVFiles(f, dirA, dirB)
		if err != nil {
			fmt.Printf("Error comparing file %s: %v\n", f, err)
		}
	}
}

func keyColumnsForFile(file string) []string {
	switch file {
	case "ad-objects.csv":
		return []string{"DistinguishedName"}
	case "autoruns.csv":
		return []string{"Location", "Name", "LaunchString"}
	case "ports.csv":
		return []string{"LocalAddress", "LocalPort"}
	case "processes.csv":
		return []string{"Name", "Path"}
	case "services.csv":
		return []string{"Name"}
	case "ad-users.csv":
		return []string{"SamAccountName"}
	case "local-users.csv":
		return []string{"Name"}
	default:
		return nil
	}
}

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
		return strings.Join(vals, " | ")
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

func formatObject(key string, obj map[string]string) string {
	if obj == nil {
		return key
	}
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

func diffMaps(a, b map[string]map[string]string) (added []string, removed []string, changed []string) {
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
