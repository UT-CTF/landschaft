package baseline

import (
	"fmt"

	"github.com/UT-CTF/landschaft/util"
	"github.com/spf13/cobra"
)

var baselineComponents = map[string]string{
	"services":    "services.ps1",
	"processes":   "processes.ps1",
	"autoruns":    "autoruns.ps1",
	"ad-users":    "ad-users.ps1",
	"local-users": "local-users.ps1",
	"ad-objects":  "ad-objects.ps1",
	"ports":       "ports.ps1",
}

var dcScripts = []string{"services", "processes", "autoruns", "ad-users", "adobjects", "ports"}
var localScripts = []string{"services", "processes", "autoruns", "local-users", "ports"}

var sysinternalsDirectory = `C:\\ProgramData\\landschaft\\sysinternals`

func SetupCommand(cmd *cobra.Command) {
	setupCompareCmd(cmd)
	setupCreateCmd(cmd)
}

func Run(cmd *cobra.Command) {
	fmt.Println(util.ErrorStyle.Render("Error: No subcommand specified\n"))
	_ = cmd.Usage()
}
