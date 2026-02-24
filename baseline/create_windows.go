package baseline

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

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
		Use:     "all <output-dir>",
		Short:   "Create all baselines into a directory",
		Long:    "Create all baselines (services, processes, autoruns, users, adobjects, ports) and save CSV files into the provided output directory.",
		Args:    cobra.ExactArgs(1),
		Example: "  landschaft baseline create all C:\\baselines",
		Run: func(cmd *cobra.Command, args []string) {
			out := args[0]
			out, _ = filepath.Abs(out)
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

	for name := range baselineComponents {
		cmdC := &cobra.Command{
			Use:   name + " <baseline-dir>",
			Short: fmt.Sprintf("Create %s baseline", name),
			Args:  cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				out := args[0]
				out, err := filepath.Abs(out)
				if err != nil {
					fmt.Printf("Error getting absolute path for output directory: %v\n", err)
					return
				}
				createSingleBaseline(name, out)
			},
		}

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
	if name == "autoruns" {
		_ = misc.EnsureSysinternals(sysinternalsDirectory)
		util.RunAndRedirectScript(fmt.Sprintf("baseline/%s", baselineComponents[name]), "-BaselinePath", fmt.Sprintf("'%s'", baselineDir), "-SysinternalsPath", fmt.Sprintf("'%s'", sysinternalsDirectory))
	} else {
		util.RunAndRedirectScript(fmt.Sprintf("baseline/%s", baselineComponents[name]), "-BaselinePath", fmt.Sprintf("'%s'", baselineDir))
	}
}
