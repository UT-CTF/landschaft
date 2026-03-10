package cmd

import (
	"fmt"

	"github.com/UT-CTF/landschaft/score"
	"github.com/UT-CTF/landschaft/util"
	"github.com/spf13/cobra"
)

var scoredServicesCmd = &cobra.Command{
	Use:   "scored-services",
	Short: "List and explain candidate scored services (listening ports)",
	Long:  `Auto-discovers listening ports and maps them to common protocols. Scoring may check different endpoints.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := score.RunList(); err != nil {
			fmt.Println(util.ErrorStyle.Render("Error: " + err.Error()))
		}
	},
}

var scoredListCmd = &cobra.Command{
	Use:   "list",
	Short: "List listening ports",
	Run: func(cmd *cobra.Command, args []string) {
		if err := score.RunList(); err != nil {
			fmt.Println(util.ErrorStyle.Render("Error: " + err.Error()))
		}
	},
}

var scoredExplainCmd = &cobra.Command{
	Use:   "explain",
	Short: "List listening ports with protocol explanations",
	Run: func(cmd *cobra.Command, args []string) {
		if err := score.RunExplain(); err != nil {
			fmt.Println(util.ErrorStyle.Render("Error: " + err.Error()))
		}
	},
}

func init() {
	scoredServicesCmd.AddCommand(scoredListCmd)
	scoredServicesCmd.AddCommand(scoredExplainCmd)
	rootCmd.AddCommand(scoredServicesCmd)
}
