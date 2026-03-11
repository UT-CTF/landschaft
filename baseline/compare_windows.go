package baseline

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
		Use:  "all",
		Short: "Compare two baseline directories",
		Long: `Compare two directories produced by 'baseline create all' and report
added/removed/changed entries for each component. Specify directories with
-a/--baseline-a and -b/--baseline-b, or omit to auto-discover the two most recent.`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			dirA, _ := cmd.Flags().GetString("baseline-a")
			dirB, _ := cmd.Flags().GetString("baseline-b")
			if dirA == "" || dirB == "" {
				foundA, foundB, err := findLatestBaselinesByName(".")
				if err != nil {
					fmt.Println("Error finding latest baselines:", err)
					return
				}
				dirA, dirB = foundA, foundB
				fmt.Printf("Comparing baselines: %q with %q\n", dirA, dirB)
			}
			compareCSVDirs(dirA, dirB)
		},
	}
	compareAllCmd.Flags().StringP("baseline-a", "a", "", "Baseline A directory")
	compareAllCmd.Flags().StringP("baseline-b", "b", "", "Baseline B directory")
	compareAllCmd.MarkFlagsRequiredTogether("baseline-a", "baseline-b")
	compareCmd.AddCommand(compareAllCmd)

	for name := range baselineComponents {
		name := name
		cmdCmp := &cobra.Command{
			Use:   name,
			Short: fmt.Sprintf("Compare %s baselines", name),
			Args:  cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				dirA, _ := cmd.Flags().GetString("baseline-a")
				dirB, _ := cmd.Flags().GetString("baseline-b")
				if dirA == "" || dirB == "" {
					foundA, foundB, err := findLatestBaselinesByName(".")
					if err != nil {
						fmt.Println("Error finding latest baselines:", err)
						return
					}
					dirA, dirB = foundA, foundB
					fmt.Printf("Comparing baselines: %q with %q\n", dirA, dirB)
				}
				if err := compareCSVFiles(fmt.Sprintf("%s.csv", name), dirA, dirB); err != nil {
					fmt.Printf("Error comparing %s: %v\n", name, err)
				}
			},
		}
		cmdCmp.Flags().StringP("baseline-a", "a", "", "Baseline A directory")
		cmdCmp.Flags().StringP("baseline-b", "b", "", "Baseline B directory")
		cmdCmp.MarkFlagsRequiredTogether("baseline-a", "baseline-b")
		compareCmd.AddCommand(cmdCmp)
	}

	cmd.AddCommand(compareCmd)
}

func compareCSVDirs(dirA, dirB string) {
	filesA, err := getCSVFileNames(dirA)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", dirA, err)
		return
	}
	filesB, err := getCSVFileNames(dirB)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", dirB, err)
		return
	}
	for _, f := range filesA {
		if slices.Contains(filesB, f) {
			if err := compareCSVFiles(f, dirA, dirB); err != nil {
				fmt.Printf("Error comparing %s: %v\n", f, err)
			}
		}
	}
}

func compareCSVFiles(fileName, dirA, dirB string) error {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	keyCols := keyColumnsForFile(fileName)
	mapA, err := loadCSVWithKey(filepath.Join(dirA, fileName), keyCols)
	if err != nil {
		return fmt.Errorf("could not load %s/%s: %v", dirA, fileName, err)
	}
	mapB, err := loadCSVWithKey(filepath.Join(dirB, fileName), keyCols)
	if err != nil {
		return fmt.Errorf("could not load %s/%s: %v", dirB, fileName, err)
	}
	fmt.Printf("Comparing %s\n", fileName)
	added, removed, changed := diffMaps(mapA, mapB)
	if len(added) > 0 {
		sort.Strings(added)
		fmt.Printf("Added in %s:\n", dirB)
		for _, k := range added {
			fmt.Printf("\t%s\n", formatObject(k, mapB[k]))
		}
	}
	if len(removed) > 0 {
		sort.Strings(removed)
		fmt.Printf("Removed from %s:\n", dirB)
		for _, k := range removed {
			fmt.Printf("\t%s\n", formatObject(k, mapA[k]))
		}
	}
	for _, c := range changed {
		fmt.Printf("\t%s\n", c)
	}
	return nil
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
	case "wmi-bindings.csv":
		return []string{"Filter", "Consumer"}
	default:
		return nil
	}
}

// findLatestBaselinesByName finds the two most recent baseline-MMDD-HHMM dirs in baseDir.
func findLatestBaselinesByName(baseDir string) (string, string, error) {
	re := regexp.MustCompile(`^baseline-\d{4}-\d{4}$`)
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return "", "", fmt.Errorf("could not read %s: %v", baseDir, err)
	}
	var candidates []string
	for _, e := range entries {
		if e.IsDir() && re.MatchString(e.Name()) {
			candidates = append(candidates, e.Name())
		}
	}
	if len(candidates) < 2 {
		return "", "", fmt.Errorf("need at least 2 baseline folders in %s", baseDir)
	}
	sort.Strings(candidates)
	n := len(candidates)
	return filepath.Join(baseDir, candidates[n-2]), filepath.Join(baseDir, candidates[n-1]), nil
}

func getCSVFileNames(dir string) ([]string, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".csv") {
			files = append(files, e.Name())
		}
	}
	return files, nil
}

func loadCSVWithKey(path string, keyCols []string) (map[string]map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	recs, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}
	if len(recs) < 1 {
		return nil, fmt.Errorf("empty csv: %s", path)
	}
	headers := recs[0]
	out := make(map[string]map[string]string)
	for _, row := range recs[1:] {
		if len(row) == 0 {
			continue
		}
		m := make(map[string]string)
		for i, cell := range row {
			if i < len(headers) {
				m[headers[i]] = cell
			}
		}
		key := buildCSVKey(row, headers, keyCols)
		out[key] = m
	}
	return out, nil
}

func buildCSVKey(row, headers, keyCols []string) string {
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

func formatObject(key string, obj map[string]string) string {
	if obj == nil {
		return key
	}
	keys := make([]string, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	lines := []string{key}
	for _, k := range keys {
		lines = append(lines, fmt.Sprintf("\t\t%s: %s", k, obj[k]))
	}
	return strings.Join(lines, "\n")
}

func diffMaps(a, b map[string]map[string]string) (added, removed, changed []string) {
	for k := range b {
		if _, ok := a[k]; !ok {
			added = append(added, k)
		}
	}
	for k := range a {
		if _, ok := b[k]; !ok {
			removed = append(removed, k)
		}
	}
	for k, va := range a {
		if vb, ok := b[k]; ok {
			for hk, hv := range va {
				if vb[hk] != hv {
					changed = append(changed,
						fmt.Sprintf("%s: %s changed from %q to %q", k, hk, hv, vb[hk]))
				}
			}
		}
	}
	return
}
