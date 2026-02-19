package triage

import (
	"fmt"
	"os"
	"os/exec"
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

	if _, err := file.Write([]byte(runNetworkTriage())); err != nil {
		// ignore error
	}

	fmt.Println(util.TitleColor.Render("Users & Groups"))

	if _, err := file.Write([]byte(runUsersTriage())); err != nil {
		// ignore error
	}

	fmt.Println(util.TitleColor.Render("OS Version"))

	if _, err := file.Write([]byte(runOSVersionTriage())); err != nil {
		// ignore error
	}

	fmt.Println(util.TitleColor.Render("Firewall"))

	if _, err := file.Write([]byte(runFirewallTriage())); err != nil {
		// ignore error
	}

	errF := file.Close()
	if errF != nil {

	}

	openOrCopyFile(file.Name())

}

func openOrCopyFile(filename string) {
	if runtime.GOOS == "windows" {
		println("\nOpen " + filename + " in notepad to copy to sheets")
		err := exec.Command("notepad.exe", filename).Start()
		if err != nil {
			err := exec.Command("explorer.exe", filename).Start()
			if err != nil {
				println(err.Error())
				return
			}
		}
		return
	}

	// Linux - install and use xclip
	err := exec.Command("sudo", "apt-get", "install", "-y", "xclip").Run()
	if err != nil {
		println("\nCannot install xclip! must use xclip or xsel to copy to sheets")
		println(err.Error())
		return
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	cmd := exec.Command("xclip", "-selection", "clipboard")
	cmd.Stdin = strings.NewReader(string(data))
	cmd.Run()

	println("\n\n ****** Triage copied to clipboard ******")
	println("\n\n Triage saved to " + filename + ". To copy:")
	println("\n\tcat " + filename + " | xclip -selection clipboard\n\n")
}
