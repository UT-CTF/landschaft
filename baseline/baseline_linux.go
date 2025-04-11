package baseline

import (
	"log"
	"os"

	"github.com/UT-CTF/landschaft/util"
)

func runBaseline() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err, "Failed to get current working directory")
		return
	}
	util.RunAndPrintScript("baseline/baseline.sh", cwd)
}
