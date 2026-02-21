package triage

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/UT-CTF/landschaft/util"
)

func Run() {
	file, err := os.Create("triage.tsv")
	if err != nil {
		fmt.Println("Error creating csv file:", err)
	}

	fmt.Println(util.TitleColor.Render("Network"))

	hostname, csv := runNetworkTriage()

	fmt.Println(util.TitleColor.Render("OS Version"))

	if _, err := file.Write([]byte(hostname + "\t" + runOSVersionTriage() + "\t" + csv)); err != nil {
		// ignore error
	}

	fmt.Println(util.TitleColor.Render("Users & Groups"))

	if _, err := file.Write([]byte(runUsersTriage())); err != nil {
		// ignore error
	}

	fmt.Println(util.TitleColor.Render("Firewall"))

	if _, err := file.Write([]byte(runFirewallTriage())); err != nil {
		// ignore error
	}

	if _, err := file.Write([]byte(getDomainInfo())); err != nil {
		// ignore error
	}

	errF := file.Close()
	if errF != nil {

	}

	printCopyInstructions()

}

func printCopyInstructions() {
	sshConn := os.Getenv("SSH_CONNECTION")
	fields := strings.Fields(sshConn)
	serverUser := os.Getenv("USER")
	serverIP := "<ip>"
	if len(fields) >= 3 {
		serverIP = fields[2]
	}

	exePath, _ := os.Executable()
	tsvPath := filepath.Join(filepath.Dir(exePath), "triage.tsv")

	catCmd := "cat"
	if runtime.GOOS == "windows" {
		catCmd = "type"
		serverUser = os.Getenv("USERNAME")
	}

	fmt.Println("To copy to clipboard:")
	fmt.Printf("\tIf your host is linux: ssh %s@%s \"%s %s\" | xclip -selection clipboard\n\n", serverUser, serverIP, catCmd, tsvPath)
	fmt.Printf("\tIf your host is windows: ssh %s@%s \"%s %s\" | clip.exe\n\n\n", serverUser, serverIP, catCmd, tsvPath)
}
