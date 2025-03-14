package embed

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"
)

func ReadScriptDir(directory string) ([]fs.DirEntry, error) {
	return scriptsFS.ReadDir(path.Join(scriptRootDir, directory))
}

func ExtractFile(scriptPath string, targetPath string) {
	file, err := scriptsFS.Open(path.Join(scriptRootDir, scriptPath))
	if err != nil {
		fmt.Println("Error opening embedded file: ", err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	targetFile, err := os.Create(targetPath)
	if err != nil {
		fmt.Println("Error creating target file: ", err)
		return
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, reader)
	if err != nil {
		fmt.Println("Error copying file: ", err)
		return
	}
}

func extractEmbeddedDir(scriptDirectory string) (string, error) {
	// Create a temporary folder to hold script dependencies
	tmpDir, err := os.MkdirTemp("", "ls-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %w", err)
	}

	// Extract the directory containing the embedded script
	embeddedDirPath := path.Join(scriptRootDir, scriptDirectory)

	embeddedDir, err := scriptsFS.ReadDir(embeddedDirPath)
	if err != nil {
		return "", fmt.Errorf("failed to read embedded scripts in %s: %w", embeddedDirPath, err)
	}
	for _, embeddedDirEntry := range embeddedDir {
		if embeddedDirEntry.IsDir() {
			fmt.Printf("WARNING: not extracting embedded directory %s in %s\n", embeddedDirEntry.Name(), embeddedDirPath)
			continue
		}

		// Read the script file's content
		embeddedFilePath := path.Join(embeddedDirPath, embeddedDirEntry.Name())
		file, err := scriptsFS.Open(embeddedFilePath)
		if err != nil {
			return "", fmt.Errorf("failed to open embedded file %s: %w", embeddedFilePath, err)
		}
		defer file.Close()

		// Create the file in the temporary directory
		reader := bufio.NewReader(file)
		tmpFilePath := path.Join(tmpDir, embeddedDirEntry.Name())

		tmpFile, err := os.Create(tmpFilePath)
		if err != nil {
			return "", fmt.Errorf("failed to create embedded file %s in temporary directory: %w", embeddedDirEntry.Name(), err)
		}
		defer tmpFile.Close()

		// Make script files executable if needed
		mode := os.FileMode(0600)
		fileName := embeddedDirEntry.Name()
		lowerFileName := strings.ToLower(fileName)
		if strings.HasSuffix(lowerFileName, ".ps1") || strings.HasSuffix(lowerFileName, ".sh") {
			mode = os.FileMode(0700) // Make script files executable
		}

		err = tmpFile.Chmod(mode)
		if err != nil {
			return "", fmt.Errorf("failed to set permissions on %s: %w", embeddedDirEntry.Name(), err)
		}

		// Write embedded file to disk
		_, err = io.Copy(tmpFile, reader)
		if err != nil {
			return "", fmt.Errorf("failed to write embedded file %s to temporary file: %w", embeddedDirEntry.Name(), err)
		}
	}

	return tmpDir, nil
}
