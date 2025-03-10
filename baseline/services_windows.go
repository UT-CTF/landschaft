package baseline

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/UT-CTF/landschaft/util"
	"github.com/spf13/cobra"
)

type baselineConfig struct {
	create  bool
	compare bool
	output  string
	files   []string
}

var baselineCfg baselineConfig

var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Create or compare services baseline",
	Run: func(cmd *cobra.Command, args []string) {
		baselineServices(cmd, baselineCfg)
	},
}

func baselineServices(cmd *cobra.Command, cfg baselineConfig) {
	if cfg.create {
		createBaseline(cfg.output)
	} else if cfg.compare {
		if len(cfg.files) != 2 {
			fmt.Printf("Expected 2 files, got %d ([%s])\n", len(cfg.files), strings.Join(cfg.files, ", "))
			_ = cmd.Usage()
			return
		}
		compareServices(cfg.files[0], cfg.files[1])
	} else {
		fmt.Println("Invalid options")
	}
}

func setupServicesCmd(cmd *cobra.Command) {
	servicesCmd.Flags().BoolVarP(&baselineCfg.create, "create", "c", false, "Create baseline")
	servicesCmd.Flags().BoolVarP(&baselineCfg.compare, "compare", "m", false, "Compare baseline")
	servicesCmd.Flags().StringVarP(&baselineCfg.output, "output", "o", "", "Output file (expected .csv file)")
	servicesCmd.Flags().StringSliceVarP(&baselineCfg.files, "files", "f", []string{}, "Files to compare (expected 2 .csv files)")
	servicesCmd.MarkFlagsMutuallyExclusive("create", "compare")
	servicesCmd.MarkFlagsOneRequired("create", "compare")
	servicesCmd.MarkFlagsRequiredTogether("create", "output")
	servicesCmd.MarkFlagsRequiredTogether("compare", "files")

	cmd.AddCommand(servicesCmd)
}

func createBaseline(csvPath string) {
	csvPath, err := filepath.Abs(csvPath)
	if err != nil {
		fmt.Println("Could not get absolute path: ", err)
		return
	}
	util.RunAndPrintScript("baseline/services.ps1", "-ExportPath", fmt.Sprintf("'%s'", csvPath))
}

func compareServices(csvPath1 string, csvPath2 string) {
	services1, err := loadServices(csvPath1)
	if err != nil {
		fmt.Println(err)
		return
	}

	services2, err := loadServices(csvPath2)
	if err != nil {
		fmt.Println(err)
		return
	}

	newServices := make([]string, 0)
	removedServices := make([]string, 0)
	commonServices := make([]string, 0)

	for name := range services1 {
		if _, ok := services2[name]; !ok {
			removedServices = append(removedServices, name)
		} else {
			commonServices = append(commonServices, name)
		}
	}

	for name := range services2 {
		if _, ok := services1[name]; !ok {
			newServices = append(newServices, name)
		}
	}

	fmt.Printf("New services: \n\t%s\n", strings.Join(newServices, "\n\t"))
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("Removed services: \n\t%s\n", strings.Join(removedServices, "\n\t"))
	fmt.Println(strings.Repeat("-", 50))
	fmt.Println("Changed services:")
	for _, name := range commonServices {
		service1 := services1[name]
		service2 := services2[name]

		for key, value := range service1 {
			if service2[key] != value {
				fmt.Printf("\tService %s: %s changed from %s to %s\n", name, key, value, service2[key])
			}
		}
	}
}

func loadServices(csvFile string) (map[string]map[string]string, error) {
	file, err := os.Open(csvFile)
	if err != nil {
		return nil, fmt.Errorf("could not open %s: %w", csvFile, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %w", csvFile, err)
	}

	services := make(map[string]map[string]string)

	rowNames := records[0]
	for _, row := range records[1:] {
		service := make(map[string]string)
		for i, cell := range row {
			service[rowNames[i]] = cell
		}
		// fmt.Println(service)
		services[service["Name"]] = service
	}
	return services, nil
}
