package misc

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var sysinternalsCmd = &cobra.Command{
	Use:                   "sysinternals [target-directory]",
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		installSysinternals(args[0])
	},
}

var firefoxCmd = &cobra.Command{
	Use:                   "firefox",
	Args:                  cobra.ExactArgs(0),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		installFirefox()
	},
}

var toolsGroup = &cobra.Group{
	ID:    "tools",
	Title: "Commands for installing tools",
}

var toolCommands = []*cobra.Command{
	sysinternalsCmd,
	firefoxCmd,
}

func setupToolsCommand(cmd *cobra.Command) {
	cmd.AddGroup(toolsGroup)

	for _, toolCmd := range toolCommands {
		toolCmd.GroupID = toolsGroup.ID
		cmd.AddCommand(toolCmd)
	}
}

func installSysinternals(targetDir string) {
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

	fmt.Println("Sysinternals Suite installed successfully")
}

// func installPython() {
// 	fmt.Println("Not implemented")
// }

func installFirefox() {
	url := "https://download.mozilla.org/?product=firefox-latest&os=win64&lang=en-US"

	data, err := downloadData(url)
	if err != nil {
		fmt.Println("Error downloading Firefox: ", err)
		return
	}

	installFile, err := os.CreateTemp("", "firefox-*.exe")
	if err != nil {
		fmt.Println("Error creating temporary file: ", err)
		return
	}

	_, err = installFile.Write(data)
	if err != nil {
		fmt.Println("Error writing data to file: ", err)
		return
	}
	installFile.Close()
	defer os.Remove(installFile.Name())

	cmd := exec.Command(installFile.Name(), "/silent")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error running installer: ", err)
		return
	}

	fmt.Println("Firefox installed successfully")
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
