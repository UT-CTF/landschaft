# Landschaft inject report

Generated: 2026-03-04T01:02:03Z

## Summary

- **Host:** AkashLaptop
- **Action log:** landschaft-actions.jsonl
- **Triage file:** triage.tsv

## Actions taken

- **2026-03-04T01:00:41Z** `C:\Users\aakas\Documents\landschaft\landschaft.exe audit --versions` (exit 0, 1215 ms)
- **2026-03-04T01:00:42Z** `C:\Users\aakas\Documents\landschaft\landschaft.exe audit --sshd` (exit 0, 878 ms)
- **2026-03-04T01:01:06Z** `C:\Users\aakas\Documents\landschaft\landschaft.exe triage` (exit 0, 8824 ms)
- **2026-03-04T01:01:28Z** `C:\Users\aakas\Documents\landschaft\landschaft.exe harden --plan` (exit 0, 0 ms)
- **2026-03-04T01:01:29Z** `C:\Users\aakas\Documents\landschaft\landschaft.exe harden rdp --plan` (exit 0, 0 ms)
- **2026-03-04T01:01:33Z** `C:\Users\aakas\Documents\landschaft\landschaft.exe harden lock-accounts --plan` (exit 0, 0 ms)

## Current state snapshot

```
AkashLaptop	"Microsoft Windows 11 Pro
10.0.26200 N/A Build 26200

Standalone Workstation"	N/A	"Tailscale
	169.254.83.107

OpenVPN Data Channel Offload
	169.254.243.243

Local Area Connection* 1
	169.254.75.72

Local Area Connection* 2
	169.254.250.239

Wi-Fi
	10.148.197.160

Bluetooth Network Connection
	169.254.88.188

"	"## TCP ##
135	2012/svchost.exe
445	4/System
5040	2372/svchost.exe
5357	4/System
5426	4/System
7680	4152/svchost.exe
8000	6088/splunkd.exe
8089	6088/splunkd.exe
8191	26408/mongod.exe
49664	1760/lsass.exe
49665	1656/wininit.exe
49666	2324/svchost.exe
49667	2620/svchost.exe
49668	5268/spoolsv.exe
49672	5836/GameManagerService3.exe
49677	1728/services.exe
54235	4/System
139	4/System

## UDP ##
NONE"	"Enabled Local Users(1): 
	aakas

Disabled Local Users (4): 
	Administrator
	DefaultAccount
	Guest
	WDAGUtilityAccount"	"docker-users (1): 
	AKASHLAPTOP\aakas

Administrators (2): 
	AKASHLAPTOP\aakas
	AkashLaptop\Administrator

Guests (1): 
	AkashLaptop\Guest

System Managed Accounts Group (1): 
	AkashLaptop\DefaultAccount

Users (1): 
	AKASHLAPTOP\aakas"	"utexas 9 - Wi-Fi 
	Public - Enabled 
	State ON 
	Firewall Policy: BlockInbound,AllowOutbound

 Tailscale - Tailscale 
	Private - Enabled 
	State ON 
	Firewall Policy: BlockInbound,AllowOutbound"	"Appinfo, AppXSvc, BDESVC, BluetoothUserService_183d02, BrokerInfrastructure, BTAGService, BthAvctpSvc, bthserv, CAMService, camsvc, cbdhsvc_183d02, CDPSv...
```

## Findings

- Review non-default services in triage output.
- Review action log for hardening steps applied.
- If present, review landschaft-detections.jsonl for suspicious activity timeline.

