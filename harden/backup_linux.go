package harden

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func takeBackup(path string, backupPath string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("'%s' does not exist", path)
	}

	basename := filepath.Base(path)
	bkupNum := 0

	err := filepath.Walk(backupPath, func(root string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file := info.Name()
			fileHyphenIdx := strings.LastIndex(file, "-")
			if file[:fileHyphenIdx] == basename {
				num, err := strconv.Atoi(file[fileHyphenIdx+1:])
				if err == nil {
					bkupNum = max(bkupNum, num)
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	bkupNum++
	helperPath := filepath.Join(backupPath, fmt.Sprintf("%s-%d", basename, bkupNum))
	err = copyFile(path, helperPath)
	if err != nil {
		return fmt.Errorf("error copying file: %v", err)
	}

	if bkupNum > 1 {
		prevPath := filepath.Join(backupPath, fmt.Sprintf("%s-%d", basename, bkupNum-1))
		cmd := exec.Command("diff", "-s", helperPath, prevPath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			err = fmt.Errorf("diff: %s", string(output))
		}
		fmt.Println(string(output))
	}

	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
