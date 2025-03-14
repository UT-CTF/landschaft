package main

import (
	"os"

	"github.com/UT-CTF/landschaft/cmd"
	"github.com/charmbracelet/log"
)

var logger *log.Logger = log.NewWithOptions(os.Stdout, log.Options{
	ReportTimestamp: false,
	ReportCaller:    true,
})

func main() {
	log.SetDefault(logger)

	cmd.Execute()
}
