package baseline

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/UT-CTF/landschaft/misc"
	"github.com/UT-CTF/landschaft/util"
	"github.com/spf13/cobra"
)

func SetupCommand(cmd *cobra.Command) {
	// create and compare parent commands
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create baselines",
	}
	compareCmd := &cobra.Command{
		Use:   "compare",
		Short: "Compare baselines",
	}

	// list of components and their script names
	components := map[string]string{
		"services":  "services.ps1",
		"processes": "processes.ps1",
		"autoruns":  "autoruns.ps1",
		"users":     "users.ps1",
		"adobjects": "adobjects.ps1",
		"ports":     "ports.ps1",
	}

	// create all (positional: output-dir)
	createAllCmd := &cobra.Command{
		Use:     "all <output-dir>",
		Short:   "Create all baselines into a directory",
		Long:    "Create all baselines (services, processes, autoruns, users, adobjects, ports) and save CSV files into the provided output directory.",
		Args:    cobra.ExactArgs(1),
		Example: "  landschaft baseline create all C:\\baselines",
		Run: func(cmd *cobra.Command, args []string) {
			out := args[0]
			out, _ = filepath.Abs(out)
			// ensure sysinternals
			siPath := `C:\\ProgramData\\landschaft\\sysinternals`
			_ = misc.EnsureSysinternals(siPath)
			for name, script := range components {
				fmt.Printf("Running %s baseline...\n", name)
				scriptPath := fmt.Sprintf("baseline/%s", script)
				// autoruns requires sysinternals path
				if name == "autoruns" {
					util.RunAndRedirectScript(scriptPath, "-BaselinePath", fmt.Sprintf("'%s'", out), "-SysinternalsPath", fmt.Sprintf("'%s'", siPath))
				} else {
					util.RunAndRedirectScript(scriptPath, "-BaselinePath", fmt.Sprintf("'%s'", out))
				}
			}
		},
	}

	createCmd.AddCommand(createAllCmd)

	// per-component create (positional: output-dir)
	for name, script := range components {
		n := name
		s := script
		cmdC := &cobra.Command{
			Use:   n + " <output-dir>",
			Short: fmt.Sprintf("Create %s baseline", n),
			Args:  cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				out := args[0]
				out, _ = filepath.Abs(out)
				siPath := `C:\\ProgramData\\landschaft\\sysinternals`
				_ = misc.EnsureSysinternals(siPath)
				scriptPath := fmt.Sprintf("baseline/%s", s)
				if n == "autoruns" {
					util.RunAndRedirectScript(scriptPath, "-BaselinePath", fmt.Sprintf("'%s'", out), "-SysinternalsPath", fmt.Sprintf("'%s'", siPath))
				} else {
					util.RunAndRedirectScript(scriptPath, "-BaselinePath", fmt.Sprintf("'%s'", out))
				}
			},
		}

		createCmd.AddCommand(cmdC)
	}

	// compare all
	compareAllCmd := &cobra.Command{
		Use:   "all <dirA> <dirB>",
		Short: "Compare two baseline directories",
		Long:  "Compare two directories produced by 'baseline create all' and report added/removed/changed entries for each component.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			compareCSVDirs(args[0], args[1])
		},
	}
	compareCmd.AddCommand(compareAllCmd)

	// per-component compare (positional: fileA fileB)
	for name := range components {
		n := name
		cmdCmp := &cobra.Command{
			Use:   n + " <fileA> <fileB>",
			Short: fmt.Sprintf("Compare %s baselines", n),
			Args:  cobra.ExactArgs(2),
			Run: func(cmd *cobra.Command, args []string) {
				fileA := args[0]
				fileB := args[1]
				// for services we have a specialized comparator
				if n == "services" {
					compareServices(fileA, fileB)
					return
				}
				mA, errA := loadGenericCSV(fileA)
				mB, errB := loadGenericCSV(fileB)
				if errA != nil || errB != nil {
					fmt.Printf("Error loading files: %v %v\n", errA, errB)
					return
				}
				added, removed, changed := diffMaps(mA, mB)
				if len(added) > 0 {
					fmt.Printf("Added:\n\t%s\n", strings.Join(added, "\n\t"))
				}
				if len(removed) > 0 {
					fmt.Printf("Removed:\n\t%s\n", strings.Join(removed, "\n\t"))
				}
				if len(changed) > 0 {
					fmt.Println("Changed entries:")
					for _, c := range changed {
						fmt.Printf("\t%s\n", c)
					}
				}
			},
		}

		compareCmd.AddCommand(cmdCmp)
	}

	// register with parent
	cmd.AddCommand(createCmd)
	cmd.AddCommand(compareCmd)
}

func Run(cmd *cobra.Command) {
	fmt.Println(util.ErrorStyle.Render("Error: No subcommand specified\n"))
	_ = cmd.Usage()
}
