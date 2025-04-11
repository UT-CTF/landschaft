package harden

import (
	"fmt"

	"github.com/UT-CTF/landschaft/embed"
)

func backup_etc(backupDirectory string) {
	fmt.Print("attempting to backup to ", backupDirectory, "\n")
	embed.ExecuteScript("harden/backup_etc.sh", true, backupDirectory)
}
func restore_etc(restoreDirectory string) {
	embed.ExecuteScript("harden/restore_etc.sh", true, restoreDirectory)
}
