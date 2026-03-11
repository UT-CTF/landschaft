package baseline

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/UT-CTF/landschaft/util"
)

func runBaseline(dir string) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		log.Error("Failed to resolve path", "err", err)
		return
	}
	util.RunAndRedirectScript("baseline/baseline.sh", absDir)
}

func compareSnapshots(baseDir string, oldNum, newNum int) {
	oldPath := filepath.Join(baseDir, "baseline", fmt.Sprintf("%d", oldNum))
	newPath := filepath.Join(baseDir, "baseline", fmt.Sprintf("%d", newNum))
	out, err := exec.Command("diff", "-r", "-u", oldPath, newPath).CombinedOutput()
	s := strings.TrimSpace(string(out))
	if err != nil && len(s) == 0 {
		log.Error("diff failed", "err", err)
		return
	}
	if s == "" {
		fmt.Println("No differences found.")
	} else {
		fmt.Println(s)
	}
}
