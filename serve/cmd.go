package serve

import (
	"os"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var certDirectory = "~/.landschaft"

func SetupCommand(cmd *cobra.Command) {
	cmd.Args = cobra.MaximumNArgs(1)
	cmd.Flags().IntP("port", "p", 8443, "Port to serve on")
}

func Run(cmd *cobra.Command, args []string) {
	// Check if the directory argument is provided
	port, err := strconv.Atoi(cmd.Flag("port").Value.String())
	if err != nil {
		log.Error("Invalid port value", "err", err)
		return
	}

	var dir string
	if len(args) < 1 {
		dir = "."
	} else {
		dir = args[0]
	}

	// Check if the directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Error("Directory does not exist", "dir", dir)
		return
	}

	// Start the HTTPS server
	err = ServeDirectoryWithHTTPS(dir, port)
	if err != nil {
		log.Error("Failed to start HTTPS server", "err", err)
	}
}
