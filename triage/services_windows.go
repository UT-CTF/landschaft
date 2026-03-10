package triage

import (
	"fmt"
	"strings"

	"github.com/UT-CTF/landschaft/util"
)

// defaultServicesWindows is a minimal allowlist of services commonly present on a default Windows install.
var defaultServicesWindows = map[string]bool{
	"EventLog": true, "FontCache": true, "PlugPlay": true, "RpcSs": true,
	"Schedule": true, "Spooler": true, "Wcmsvc": true, "Winmgmt": true,
	"wuauserv": true, "BITS": true, "CryptSvc": true, "DcomLaunch": true,
	"Dhcp": true, "Dnscache": true, "LanmanServer": true, "LanmanWorkstation": true,
	"lmhosts": true, "NlaSvc": true, "nsi": true, "Power": true,
	"ProfSvc": true, "SamSs": true, "SessionEnv": true, "StateRepository": true,
	"SystemEventsBroker": true, "Themes": true, "TimeBrokerSvc": true,
	"TokenBroker": true, "UmRdpService": true, "UsoSvc": true, "W32Time": true,
	"WdiServiceHost": true, "WdiSystemHost": true,
	"WinDefend": true, "WlanSvc": true, "WpnService": true,
	"Audiosrv": true, "AudioEndpointBuilder": true, "BFE": true, "CoreMessagingRegistrar": true,
	"DPS": true, "gpsvc": true, "hidserv": true, "IKEEXT": true,
	"iphlpsvc": true, "KeyIso": true, "LSM": true, "NcbService": true,
	"Netlogon": true, "netprofm": true, "NgcCtnrSvc": true, "NgcSvc": true,
	"PolicyAgent": true, "SENS": true, "ShellHWDetection": true, "StorSvc": true,
	"SysMain": true, "TieringEngineService": true, "TrkWks": true, "VaultSvc": true,
	"WbioSrvc": true, "Wecsvc": true, "WpnUserService": true,
}

func runServicesTriage() string {
	out, err := util.RunScriptQuiet("triage/services.ps1")
	if err != nil {
		return "err\t"
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	var nonDefault []string
	for _, line := range lines {
		name := strings.TrimSpace(line)
		if name == "" {
			continue
		}
		if !defaultServicesWindows[name] {
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
