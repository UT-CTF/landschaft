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
		return
	}

	scanner := bufio.NewScanner(bytes.NewReader(sshd_config))
	
	fmt.Println("listing bad sshd_config rules:")
	sshd_map := makeOutputMap()
	for scanner.Scan() {
		line := strings.ToLower(scanner.Text())
		line = strings.TrimRight(line, " \t\r\n")
		ok := sshd_map[line]
		if ok {
			fmt.Println(line)
		}
		arr := strings.Split(line, " ")
		if arr[0] == "ciphers" {
			checkCiphers(arr[1])
		}
	}
}

func makeOutputMap() map[string]bool {
	sshdMap := map[string]bool{
		"permitrootlogin yes":       true,
		"permitemptypasswords yes":  true,
		"x11forwarding yes":         true,
		"ignorerhosts no":           true,
		"permituserenvironment yes": true,
	}
	return sshdMap
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
