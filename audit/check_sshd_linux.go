package audit

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func checkSSHD() {
	cmd := exec.Command("sshd", "-T")
	sshd_config, err := cmd.Output()
	if err != nil {
		fmt.Println("Error reading sshd_config release:", err)
		fmt.Println("This could be caused by running w/o sudo")
		return
	}

	scanner := bufio.NewScanner(bytes.NewReader(sshd_config))

	fmt.Println("listing bad sshd_config rules:")
	sshd_map := makeOutputMap()
	for scanner.Scan() {
		line := strings.ToLower(scanner.Text())
		line = strings.TrimRight(line, " \t\r\n")
		fields := strings.Fields(line)
		val, ok := sshd_map[fields[0]]
		if ok && val == fields[1] {
			fmt.Println(line)
		} else {
			arr := strings.Split(line, " ")
			if arr[0] == "ciphers" {
				checkCiphers(arr[1])
			} else if arr[0] == "protocol" {
				checkProtocols(arr[1])
			}
		}
	}
}

func makeOutputMap() map[string]string {
	sshdMap := map[string]string{
		"permitrootlogin":       "yes",
		"permitemptypasswords":  "yes",
		"x11forwarding":         "yes",
		"ignorerhosts":          "no",
		"permituserenvironment": "yes",
	}
	return sshdMap
}

func checkProtocols(line string) bool {
	return strings.Contains(line, "1")
}

func checkCiphers(line string) {
	arr := strings.Split(line, ",")
	badCiphers := map[string]bool{
		"des":          true,
		"3des-cbc":     true,
		"arcfour":      true,
		"arcfour128":   true,
		"arcfour256":   true,
		"blowfish-cbc": true,
		"hmac-md5":     true,
		"hmac-md5-96":  true,
		"hmac-sha1":    true,
		"hmac-sha1-96": true,
	}
	ciphersFound := make([]string, 0)
	for i := 0; i < len(arr); i++ {
		cipher := arr[i]
		if badCiphers[cipher] {
			ciphersFound = append(ciphersFound, cipher)
		}
	}
	if len(ciphersFound) > 0 {
		fmt.Println("Ciphers to consider removing:")
		fmt.Println(ciphersFound)
	}
}
