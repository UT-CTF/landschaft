package cmd

import (
	"fmt"

	"github.com/UT-CTF/landschaft/report"
	"github.com/UT-CTF/landschaft/util"
	"github.com/spf13/cobra"
)

var reportArgs struct {
	actionLog string
	triage    string
	out       string
}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate reports for injects and documentation",
}

var reportInjectCmd = &cobra.Command{
	Use:   "inject",
	Short: "Generate Markdown inject report from action log and triage",
	Run: func(cmd *cobra.Command, args []string) {
		actionLog := util.ParseActionLogPath(reportArgs.actionLog)
		out := reportArgs.out
		if out == "" {
			out = "inject-report.md"
		}
		if err := report.InjectReport(actionLog, reportArgs.triage, out); err != nil {
			fmt.Println(util.ErrorStyle.Render("Error: " + err.Error()))
			return
		}
		fmt.Println("Wrote inject report to", out)
	},
}

func init() {
	reportInjectCmd.Flags().StringVar(&reportArgs.actionLog, "action-log", "", "Path to landschaft-actions.jsonl")
	reportInjectCmd.Flags().StringVar(&reportArgs.triage, "triage", "triage.tsv", "Path to triage TSV")
	reportInjectCmd.Flags().StringVarP(&reportArgs.out, "out", "o", "inject-report.md", "Output Markdown file")
	reportCmd.AddCommand(reportInjectCmd)
	rootCmd.AddCommand(reportCmd)
}
