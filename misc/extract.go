package misc

import (
	"fmt"
	"os"
	"path"

	"github.com/UT-CTF/landschaft/embed"
	"github.com/spf13/cobra"
)

var extractCmd = &cobra.Command{
	Use:                   "extract [target-directory]",
	Short:                 "Extracts embedded files from the binary",
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		extractFiles(args[0], "")
	},
}

func setupExtractCommand(cmd *cobra.Command) {
	cmd.AddCommand(extractCmd)
}

func extractFiles(targetDir string, curPath string) {
	dirEntries, err := embed.ReadScriptDir(curPath)
	if err != nil {
		fmt.Println("Error reading embedded directory: ", err)
		return
	}
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			err := os.MkdirAll(path.Join(targetDir, dirEntry.Name()), 0755)
			if err != nil {
				fmt.Println("Error creating directory: ", err)
				continue
			}
			extractFiles(path.Join(targetDir, dirEntry.Name()), path.Join(curPath, dirEntry.Name()))
		} else {
			embed.ExtractFile(path.Join(curPath, dirEntry.Name()), path.Join(targetDir, dirEntry.Name()))
		}
	}
}
