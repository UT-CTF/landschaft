package misc

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/UT-CTF/landschaft/embed"
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

var nxlogUrl string
var install bool
var certFile string
var configFile string
var nxlogCmd = &cobra.Command{
	Use:  "nxlog",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if install {
			installNxlog()
		}
		if len(certFile) > 0 {
			loadNxlogCert(certFile)
		}
		if len(configFile) > 0 {
			loadNxlogConfig(configFile)
		}
	},
}

var toolsGroup = &cobra.Group{
	ID:    "tools",
	Title: "Commands for installing tools",
}

var toolCommands = []*cobra.Command{
	sysinternalsCmd,
	firefoxCmd,
	nxlogCmd,
}

func setupToolsCommand(cmd *cobra.Command) {

	nxlogCmd.Flags().StringVarP(&nxlogUrl, "url", "u", "", "URL to download Nxlog from")
	nxlogCmd.Flags().BoolVarP(&install, "install", "i", false, "Install Nxlog")
	nxlogCmd.Flags().StringVarP(&certFile, "cert", "c", "", "Path to the certificate file")
	nxlogCmd.Flags().StringVarP(&configFile, "config", "f", "", "Path to the configuration file")

	nxlogCmd.MarkFlagsRequiredTogether("install", "url")

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

func installNxlog() {
	url := nxlogUrl
	data, err := downloadData(url)
	if err != nil {
		fmt.Println("Error downloading Nxlog: ", err)
		return
	}

	installFile, err := os.CreateTemp("", "nxlog-*.msi")
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
	// defer os.Remove(installFile.Name())

	cmd := exec.Command("msiexec", "/i", installFile.Name(), "/quiet")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error running installer: ", err)
		fmt.Println(installFile.Name())
		return
	}

	cmd = exec.Command("sc", "stop", "nxlog")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error stopping Nxlog service: ", err)
		return
	}

	tmpConfPath := path.Join(os.TempDir(), "nxlog-landschaft-temp.conf")
	embed.ExtractFile("misc/nxlog.conf", tmpConfPath)
	defer os.Remove(tmpConfPath)
	loadNxlogConfig(tmpConfPath)

	fmt.Println("Nxlog installed successfully")
}

func loadNxlogCert(certFile string) {
	cfile, err := os.Create("C:/Program Files/nxlog/cert/graylog-ca.pem")
	if err != nil {
		fmt.Println("Error creating certificate file: ", err)
		return
	}
	defer cfile.Close()

	data, err := os.ReadFile(certFile)
	if err != nil {
		fmt.Println("Error reading certificate file: ", err)
		return
	}

	_, err = cfile.Write(data)
	if err != nil {
		fmt.Println("Error writing certificate file: ", err)
		return
	}

	fmt.Println("Certificate loaded successfully")
}

func loadNxlogConfig(configFile string) {
	cfile, err := os.Create("C:/Program Files/nxlog/conf/nxlog.conf")
	if err != nil {
		fmt.Println("Error creating configuration file: ", err)
		return
	}
	defer cfile.Close()

	data, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Println("Error reading configuration file: ", err)
		return
	}

	_, err = cfile.Write(data)
	if err != nil {
		fmt.Println("Error writing configuration file: ", err)
		return
	}

	fmt.Println("Configuration loaded successfully")
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
