package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/UT-CTF/landschaft/hunt"
	"github.com/UT-CTF/landschaft/util"
	"github.com/spf13/cobra"
)

var huntArgs struct {
	since        string
	detectionsLog string
}

var huntCmd = &cobra.Command{
	Use:   "hunt",
	Short: "Collect and highlight suspicious log activity (no Graylog)",
	Long:  `Reads Windows Event Log or Linux auth/journal, applies heuristics, appends to detections JSONL.`,
	Run: func(cmd *cobra.Command, args []string) {
		since := 30 * time.Minute
		if huntArgs.since != "" {
			d, err := time.ParseDuration(huntArgs.since)
			if err != nil {
				fmt.Println(util.ErrorStyle.Render("Invalid --since: " + err.Error()))
				return
			}
			since = d
		}
		logPath := huntArgs.detectionsLog
		if logPath == "" {
			logPath = os.Getenv("LANDSCHAFT_DETECTIONS_LOG")
			if logPath == "" {
				logPath = "landschaft-detections.jsonl"
			}
		}
		if err := hunt.Run(since, logPath); err != nil {
			fmt.Println(util.ErrorStyle.Render("Error: " + err.Error()))
		}
	},
}

func init() {
	huntCmd.Flags().StringVar(&huntArgs.since, "since", "30m", "Time window (e.g. 30m, 1h)")
	huntCmd.Flags().StringVar(&huntArgs.detectionsLog, "detections-log", "", "Path to detections JSONL (default: landschaft-detections.jsonl)")
	rootCmd.AddCommand(huntCmd)
}
