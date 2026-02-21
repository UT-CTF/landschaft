package triage

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/UT-CTF/landschaft/util"
)

func runUsersTriage() string {
	csv := printUsers()
	csv += "\t"
	csv += printGroups()
	csv += "\t"
	return csv
}

type user struct {
	name  string
	uid   string
	gid   string
	shell string
}

type group struct {
	name  string
	gid   string
	users []string
}

func printUsers() string {
	file, err := os.Open("/etc/passwd")
	if err != nil {
		fmt.Println("Error reading /etc/passwd:", err)
		return "err"
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	userList := make([]user, 0)
	for scanner.Scan() {
		line := scanner.Text()
		arr := strings.Split(line, ":")
		// avoid all users that can't login
		if arr[len(arr)-1] != "/bin/false" && arr[len(arr)-1] != "/usr/sbin/nologin" {
			userList = append(userList, user{
				name:  arr[0],
				uid:   arr[2],
				gid:   arr[3],
				shell: arr[len(arr)-1],
			})
		}
	}

	// convert to [][]string
	var result string
	userListStr := make([][]string, len(userList))
	for i, u := range userList {
		userListStr[i] = []string{u.name, u.uid, u.gid, u.shell}
		result += fmt.Sprintf("%s (%s, %s, %s) \n", u.name, u.uid, u.gid, u.shell)
	}

	t := util.StyledTable().Rows(userListStr...)
	fmt.Println(t.Render())
	return "\"" + result + "\""
}

func printGroups() string {
	file, err := os.Open("/etc/group")
	if err != nil {
		fmt.Println("Error reading /etc/group:", err)
		return "err"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	groupList := make([]group, 0)

	for scanner.Scan() {
		line := scanner.Text()
		arr := strings.Split(line, ":")
		if len(arr) < 4 {
			continue
		}
		// avoid empty groups
		if arr[3] == "" {
			continue
		}
		groupList = append(groupList, group{
			name:  arr[0],
			gid:   arr[2],
			users: strings.Split(arr[3], ","),
		})
	}

	// convert to [][]string
	var result string
	groupListStr := make([][]string, len(groupList))
	for i, g := range groupList {
		groupListStr[i] = []string{g.name, g.gid, strings.Join(g.users, ",")}
		result += fmt.Sprintf("%s (GID: %s) - %s\n\n", g.name, g.gid, strings.Join(g.users, ", "))
	}

	t := util.StyledTable().Headers("group", "gid", "users").Rows(groupListStr...)
	fmt.Println(t.Render())
	return "\"" + result + "\""
}
