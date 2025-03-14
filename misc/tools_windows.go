package misc

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var toolsDict = map[string]func(){
	"sysinternals": installSysinternals,
	"python":       installPython,
	"firefox":      installFirefox,
}

var targetDir = ""

func getValidKeys(dict map[string]func()) []string {
	keys := make([]string, 0, len(dict))
	for k := range dict {
		keys = append(keys, k)
	}
	return keys
}

var toolsCmd = &cobra.Command{
	Use:       "tools",
	Short:     "Install various tools",
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: getValidKeys(toolsDict),
	Run: func(cmd *cobra.Command, args []string) {
		toolsDict[args[0]]()
	},
}

func setupToolsCommand(cmd *cobra.Command) {
	toolsCmd.Flags().StringVarP(&targetDir, "dir", "d", "", "Target directory to install tools")

	cmd.AddCommand(toolsCmd)
}

func installSysinternals() {
	if len(targetDir) == 0 {
		fmt.Println("Error: No target directory specified")
		return
	}

	url := "https://download.sysinternals.com/files/SysinternalsSuite.zip"

	zipData, err := downloadData(url)
	if err != nil {
		fmt.Println("Error downloading Sysinternals Suite: ", err)
		return
	}

	extractZip(zipData, targetDir)
}

func installPython() {
	fmt.Println("Not implemented")
}

func installFirefox() {
}

func downloadData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error downloading data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error downloading data: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading data: %w", err)
	}

	return data, nil
}

func extractZip(zipData []byte, targetDir string) {
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		fmt.Println("Error reading Sysinternals Suite: ", err)
		return
	}

	err = os.MkdirAll(targetDir, 0755)
	if err != nil {
		fmt.Println("Error creating target directory")
		return
	}

	for _, file := range zipReader.File {
		fileReader, err := file.Open()
		if err != nil {
			fmt.Println("Error opening file: ", err)
			return
		}
		defer fileReader.Close()

		extractedPath := filepath.Join(targetDir, file.Name)

		if file.FileInfo().IsDir() {
			err = os.MkdirAll(extractedPath, os.ModePerm)
			if err != nil {
				fmt.Println("Error creating directory: ", err)
				return
			}
		} else {
			// If the file is not a directory, create and write the file
			extractedFile, err := os.Create(extractedPath)
			if err != nil {
				fmt.Println("Error creating file: ", err)
				return
			}
			defer extractedFile.Close()

			_, err = io.Copy(extractedFile, fileReader)
			if err != nil {
				fmt.Println("Error copying file: ", err)
				return
			}
		}
	}
}
