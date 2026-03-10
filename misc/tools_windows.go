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
	zipData, err := downloadData("https://download.sysinternals.com/files/SysinternalsSuite.zip")
	if err != nil {
		fmt.Println("Error downloading Sysinternals Suite:", err)
		return
	}
	extractZip(zipData, targetDir)
	fmt.Println("Sysinternals Suite installed successfully")
}

func installFirefox() {
	tmpPath, err := downloadToTemp(
		"https://download.mozilla.org/?product=firefox-latest&os=win64&lang=en-US",
		"firefox-*.exe",
	)
	if err != nil {
		fmt.Println("Error downloading Firefox:", err)
		return
	}
	defer os.Remove(tmpPath)
	if err := exec.Command(tmpPath, "/silent").Run(); err != nil {
		fmt.Println("Error running installer:", err)
		return
	}
	fmt.Println("Firefox installed successfully")
}

func installNxlog() {
	tmpPath, err := downloadToTemp(nxlogUrl, "nxlog-*.msi")
	if err != nil {
		fmt.Println("Error downloading Nxlog:", err)
		return
	}
	if err := exec.Command("msiexec", "/i", tmpPath, "/quiet").Run(); err != nil {
		fmt.Println("Error running installer:", err, tmpPath)
		return
	}
	if err := exec.Command("sc", "stop", "nxlog").Run(); err != nil {
		fmt.Println("Error stopping Nxlog service:", err)
		return
	}
	tmpConfPath := path.Join(os.TempDir(), "nxlog-landschaft-temp.conf")
	embed.ExtractFile("misc/nxlog.conf", tmpConfPath)
	defer os.Remove(tmpConfPath)
	loadNxlogConfig(tmpConfPath)
	fmt.Println("Nxlog installed successfully")
}

func loadNxlogCert(src string) {
	if err := copyFileTo(`C:\Program Files\nxlog\cert\graylog-ca.pem`, src); err != nil {
		fmt.Println("Error loading certificate:", err)
		return
	}
	fmt.Println("Certificate loaded successfully")
}

func loadNxlogConfig(src string) {
	if err := copyFileTo(`C:\Program Files\nxlog\conf\nxlog.conf`, src); err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}
	fmt.Println("Configuration loaded successfully")
}

// copyFileTo copies the file at src to dst, overwriting dst if it exists.
func copyFileTo(dst, src string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("reading %s: %w", src, err)
	}
	return os.WriteFile(dst, data, 0644)
}

// downloadToTemp downloads url and writes the result to a temp file with the given name pattern.
// The caller is responsible for removing the returned path.
func downloadToTemp(url, pattern string) (string, error) {
	data, err := downloadData(url)
	if err != nil {
		return "", err
	}
	f, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", fmt.Errorf("creating temp file: %w", err)
	}
	defer f.Close()
	if _, err := f.Write(data); err != nil {
		os.Remove(f.Name())
		return "", fmt.Errorf("writing temp file: %w", err)
	}
	return f.Name(), nil
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
		fmt.Println("Error reading Sysinternals Suite:", err)
		return
	}
	if err = os.MkdirAll(targetDir, 0755); err != nil {
		fmt.Println("Error creating target directory")
		return
	}
	for _, file := range zipReader.File {
		fileReader, err := file.Open()
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer fileReader.Close()
		extractedPath := filepath.Join(targetDir, file.Name)
		if file.FileInfo().IsDir() {
			if err = os.MkdirAll(extractedPath, os.ModePerm); err != nil {
				fmt.Println("Error creating directory:", err)
				return
			}
		} else {
			extractedFile, err := os.Create(extractedPath)
			if err != nil {
				fmt.Println("Error creating file:", err)
				return
			}
			defer extractedFile.Close()
			if _, err = io.Copy(extractedFile, fileReader); err != nil {
				fmt.Println("Error copying file:", err)
				return
			}
		}
	}
}
