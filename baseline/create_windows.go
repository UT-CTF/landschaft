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
		Short: "Create all baselines into a timestamped directory",
		Long: `Create all baselines (services, processes, autoruns, users, ports, wmi) and save CSV
files into the output directory. Use -o/--output to specify the directory; if omitted
a new directory named baseline-MMDD-HHMM is created in the current working directory.`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			out, _ := cmd.Flags().GetString("output")
			if out == "" {
				now := time.Now()
				out = fmt.Sprintf("baseline-%02d%02d-%02d%02d", int(now.Month()), now.Day(), now.Hour(), now.Minute())
			}
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
	createAllCmd.Flags().StringP("output", "o", "", "Output folder (defaults to timestamped folder)")
	createCmd.AddCommand(createAllCmd)

	for name := range baselineComponents {
		name := name
		cmdC := &cobra.Command{
			Use:   name,
			Short: fmt.Sprintf("Create %s baseline", name),
			Args:  cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				out, _ := cmd.Flags().GetString("output")
				if out == "" {
					now := time.Now()
					out = fmt.Sprintf("baseline-%02d%02d-%02d%02d", int(now.Month()), now.Day(), now.Hour(), now.Minute())
				}
				out, err := filepath.Abs(out)
				if err != nil {
					fmt.Println("Error resolving output path:", err)
					return
				}
				_ = os.MkdirAll(out, 0755)
				createSingleBaseline(name, out)
			},
		}
		cmdC.Flags().StringP("output", "o", "", "Output folder (defaults to timestamped folder)")
		createCmd.AddCommand(cmdC)
	}

	cmd.AddCommand(createCmd)
}

func checkIfDomainController() bool {
	out, err := exec.Command("wmic", "computersystem", "get", "domainrole").Output()
	if err != nil {
		return false
	}
	s := string(out)
	return strings.Contains(s, "4") || strings.Contains(s, "5")
}

func createSingleBaseline(name, baselineDir string) {
	scriptName, ok := baselineComponents[name]
	if !ok {
		fmt.Printf("Unknown baseline component: %s\n", name)
		return
	}
	if name == "autoruns" {
		_ = misc.EnsureSysinternals(sysinternalsDirectory)
		util.RunAndRedirectScript(
			fmt.Sprintf("baseline/%s", scriptName),
			"-BaselinePath", fmt.Sprintf("'%s'", baselineDir),
			"-SysinternalsPath", fmt.Sprintf("'%s'", sysinternalsDirectory),
		)
	} else {
		util.RunAndRedirectScript(
			fmt.Sprintf("baseline/%s", scriptName),
			"-BaselinePath", fmt.Sprintf("'%s'", baselineDir),
		)
	}
}
