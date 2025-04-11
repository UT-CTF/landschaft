package harden

import (
	"fmt"

	"github.com/UT-CTF/landschaft/embed"
)

func backup_etc(backupDirectory string) {
	embed.ExecuteScript("harden/backup_etc.sh", false, fmt.Sprintf("-Path '%s'", backupDirectory))
}

func restore_etc(restoreDirectory string) {
	embed.ExecuteScript("harden/restore_etc.sh", false, fmt.Sprintf("-Path '%s'", restoreDirectory))
}
