package baseline

import (
	"fmt"

	"github.com/UT-CTF/landschaft/util"
	"github.com/spf13/cobra"
)

func SetupCommand(cmd *cobra.Command) {
	var dir string
	cmd.PersistentFlags().StringVarP(&dir, "dir", "d", ".", "Directory to store baseline snapshots")

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a baseline snapshot",
		Args:  cobra.NoArgs,
		Run: func(c *cobra.Command, args []string) {
			runBaseline(dir)
		},
	}

	var oldNum, newNum int
	compareCmd := &cobra.Command{
		Use:   "compare",
		Short: "Diff two baseline snapshots",
		Args:  cobra.NoArgs,
		Run: func(c *cobra.Command, args []string) {
			compareSnapshots(dir, oldNum, newNum)
		},
	}
	compareCmd.Flags().IntVar(&oldNum, "old", 1, "Older snapshot number")
	compareCmd.Flags().IntVar(&newNum, "new", 2, "Newer snapshot number")

	cmd.AddCommand(createCmd, compareCmd)
}

func Run(cmd *cobra.Command) {
	fmt.Println(util.ErrorStyle.Render("Error: No subcommand specified\n"))
	_ = cmd.Usage()
}
