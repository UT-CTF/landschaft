package triage

import (
	"fmt"
	"strings"

	"github.com/UT-CTF/landschaft/util"
)

// defaultServicesLinux is a minimal allowlist of services commonly present on a default Linux install.
// Services not in this set are printed as "non-default".
var defaultServicesLinux = map[string]bool{
	"ssh": true, "sshd": true,
	"cron": true, "crond": true,
	"rsyslog": true, "systemd-journald": true,
	"dbus": true, "dbus-org.freedesktop.NetworkManager": true,
	"network-manager": true, "NetworkManager": true,
	"systemd-networkd": true, "systemd-resolved": true,
	"systemd-timesyncd": true, "systemd-udevd": true,
	"systemd-logind": true, "user@": true,
	"getty@": true, "serial-getty@": true,
	"systemd-user-sessions": true, "-.slice": true,
	"system.slice": true, "user.slice": true,
	"polkit": true, "accounts-daemon": true,
	"atd": true, "cups-browsed": true, "cups": true,
	"irqbalance": true, "keyboard-setup": true,
	"networkd-dispatcher": true, "packagekit": true,
	"plymouth": true, "plymouth-quit-wait": true,
	"snapd": true, "snapd.apparmor": true,
	"snapd.autoimport": true, "snapd.core-fixup": true,
	"snapd.recovery": true, "snapd.seeded": true,
	"systemd-fsck": true, "systemd-random-seed": true,
	"ufw": true, "unattended-upgrades": true,
	"apparmor": true, "containerd": true, "docker": true,
	"gdm": true, "display-manager": true,
	"multi-user": true, "graphical": true,
	"remote-fs": true, "remote-fs.target": true,
}

func runServicesTriage() string {
	out, err := util.RunScriptQuiet("triage/services.sh")
	if err != nil {
		return "err\t"
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	var nonDefault []string
	for _, line := range lines {
		name := strings.TrimSpace(strings.TrimSuffix(line, ".service"))
		if name == "" {
			continue
		}
		matched := defaultServicesLinux[name]
		if !matched {
			for k := range defaultServicesLinux {
				if strings.HasPrefix(name, k) {
					matched = true
					break
				}
			}
		}
		if !matched {
			nonDefault = append(nonDefault, name)
		}
	}
	result := strings.Join(nonDefault, ", ")
	if result == "" {
		result = "(none)"
	}
	fmt.Println("Non-default running services:", result)
	return "\"" + result + "\"\t"
}
