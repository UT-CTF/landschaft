package baseline

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/UT-CTF/landschaft/misc"
	"github.com/UT-CTF/landschaft/util"
	"github.com/spf13/cobra"
)

func setupCreateCmd(cmd *cobra.Command) {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create baselines",
	}

	createAllCmd := &cobra.Command{
		Use:   "all",
		Short: "Create all baselines into a directory",
		Long:  "Create all baselines (services, processes, autoruns, users, adobjects, ports) and save CSV files into the provided output directory. Use -o/--output to specify the output directory; if omitted a new directory named baseline-MMDD-HHMM will be created in the current working directory.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			out := ""
			// output flag selects the output folder
			if of, _ := cmd.Flags().GetString("output"); of != "" {
				out = of
			}
			if out == "" {
				// create default folder baseline-MMDD-HHMM in cwd
				now := time.Now()
				out = fmt.Sprintf("baseline-%02d%02d-%02d%02d", int(now.Month()), now.Day(), now.Hour(), now.Minute())
			}
			// ensure absolute path and create directory
			out, _ = filepath.Abs(out)
			_ = os.MkdirAll(out, 0755)
			_ = misc.EnsureSysinternals(sysinternalsDirectory)
			componentList := localScripts
			if checkIfDomainController() {
				componentList = dcScripts
			}

			for _, component := range componentList {
				fmt.Printf("Creating baseline for %s...\n", component)
				createSingleBaseline(component, out)
			}
		},
	}
	createCmd.AddCommand(createAllCmd)
	createAllCmd.Flags().StringP("output", "o", "", "Custom output folder name")

	for name := range baselineComponents {
		cmdC := &cobra.Command{
			Use:   name,
			Short: fmt.Sprintf("Create %s baseline", name),
			Args:  cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				out, _ := cmd.Flags().GetString("output")
				if out == "" {
					// create default folder baseline-MMDD-HHMM in cwd
					now := time.Now()
					out = fmt.Sprintf("baseline-%02d%02d-%02d%02d", int(now.Month()), now.Day(), now.Hour(), now.Minute())
				}
				out, err := filepath.Abs(out)
				if err != nil {
					fmt.Printf("Error getting absolute path for output directory: %v\n", err)
					return
				}
				_ = os.MkdirAll(out, 0755)
				createSingleBaseline(name, out)
			},
		}

		cmdC.Flags().StringP("output", "o", "", "Output folder for this baseline (defaults to timestamped folder)")
		createCmd.AddCommand(cmdC)
	}

	cmd.AddCommand(createCmd)
}

func checkIfDomainController() bool {
	cmd := exec.Command("wmic", "computersystem", "get", "domainrole")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error checking domain role: %v\n", err)
		return false
	}
	outStr := string(output)
	return strings.Contains(outStr, "4") || strings.Contains(outStr, "5")
}

func createSingleBaseline(name, baselineDir string) {
	scriptName, ok := baselineComponents[name]
	if !ok {
		fmt.Printf("Unknown baseline component: %s\n", name)
		return
	}
	if name == "autoruns" {
		_ = misc.EnsureSysinternals(sysinternalsDirectory)
		util.RunAndRedirectScript(fmt.Sprintf("baseline/%s", scriptName), "-BaselinePath", fmt.Sprintf("'%s'", baselineDir), "-SysinternalsPath", fmt.Sprintf("'%s'", sysinternalsDirectory))
	} else {
		util.RunAndRedirectScript(fmt.Sprintf("baseline/%s", scriptName), "-BaselinePath", fmt.Sprintf("'%s'", baselineDir))
	}
}
