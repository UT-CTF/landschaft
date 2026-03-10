package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/UT-CTF/landschaft/util"
	"github.com/spf13/cobra"
)

// Version information - will be set during build
var (
	Version   = "dev"
	BuildTime = "unknown"
)

var actionLogPath string
var actionLogStart time.Time

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "landschaft",
	Short: "A cybersecurity system tool",
	Long: `Landschaft is a cross-platform cybersecurity tool designed for rapid system
hardening, triage, and monitoring.`,
	Version: fmt.Sprintf("%s (built: %s)", Version, BuildTime),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		actionLogStart = time.Now()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		path := util.ParseActionLogPath(actionLogPath)
		util.AppendActionLog(path, util.NewActionLogEntry(os.Args, 0, actionLogStart))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		path := util.ParseActionLogPath(actionLogPath)
		util.AppendActionLog(path, util.NewActionLogEntry(os.Args, 1, actionLogStart))
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&actionLogPath, "action-log", "", "Path to action log JSONL (default: ./landschaft-actions.jsonl or LANDSCHAFT_ACTION_LOG)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
