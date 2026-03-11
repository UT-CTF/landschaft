# Landschaft

> Changing the security landscape, one script at a time.

Landschaft is a cross-platform cybersecurity tool built for rapid system hardening, triage, and monitoring — designed for use in CCDC and similar time-constrained environments. It compiles to a single self-contained binary that embeds all scripts, so nothing else needs to be installed on the target host.

**Platform support:** Linux and Windows (x86-64). Most features have native implementations for both; where a feature is OS-specific, it prints a clear message and exits gracefully on the unsupported platform.

---

## Table of Contents

- [Building](#building)
- [Global flags](#global-flags)
- [Commands](#commands)
  - [triage](#triage)
  - [audit](#audit)
  - [harden](#harden)
  - [hunt](#hunt)
  - [scored-services](#scored-services)
  - [wazuh](#wazuh)
  - [graylog](#graylog)
  - [ldap](#ldap)
  - [misc](#misc)
  - [serve](#serve)
  - [baseline](#baseline)
  - [report](#report)
- [Action log](#action-log)
- [Workflow tutorial](#workflow-tutorial)

---

## Building

Requires Go 1.23+.

```bash
go build -o landschaft .
```

Cross-compile for a different platform:

```bash
GOOS=linux  GOARCH=amd64 go build -o landschaft-linux  .
GOOS=windows GOARCH=amd64 go build -o landschaft.exe   .
```

Embed build metadata:

```bash
go build -ldflags "-X github.com/UT-CTF/landschaft/cmd.Version=1.0.0 -X 'github.com/UT-CTF/landschaft/cmd.BuildTime=$(date -u)'" -o landschaft .
```

---

## Global flags

| Flag | Description |
|---|---|
| `--action-log <path>` | Path to the JSONL action log (default: `./landschaft-actions.jsonl`, or env `LANDSCHAFT_ACTION_LOG`) |
| `--version` | Print version and build time |

---

## Commands

### triage

Collects a snapshot of the current host state and writes it to `triage.tsv`. Run this first on every machine.

```bash
landschaft triage
```

Collects and prints:
- Hostname and FQDN
- Network interfaces and IPv4 addresses
- Open TCP/UDP listening ports (netstat)
- OS name and version
- Local users and groups (filters out nologin/system accounts)
- Firewall status and rules
- Running non-default services
- Domain membership (LDAP/SSSD on Linux; Active Directory/Workgroup on Windows)

Output is written to `triage.tsv` in the current directory. At the end it prints the `scp` command to copy the file back to your local machine (already on clipboard on Windows).

---

### audit

Audits the system for security issues and known vulnerabilities.

```bash
landschaft audit [flags]
```

| Flag | Default | Description |
|---|---|---|
| `--sshd` | true | Check SSHD configuration for insecure settings |
| `--versions` | true | Check installed packages against the OSV vulnerability database |
| `--max <n>` | 50 | Maximum number of packages to query |
| `--timeout <duration>` | 5s | OSV API request timeout |

**SSHD check (Linux):** runs `sshd -T` and flags dangerous settings — PermitRootLogin, PermitEmptyPasswords, Protocol 1, weak ciphers, X11 forwarding, etc.

**Version check (Linux):** detects the distro ecosystem (Ubuntu, Debian, Alpine, Rocky, AlmaLinux), queries installed packages via `dpkg-query`, and checks a curated set of security-relevant packages against [osv.dev](https://osv.dev).

**Version check (Windows):** enumerates all installed software from the registry and displays it as a table. OSV lookup is not performed (no reliable package-to-ecosystem mapping exists for arbitrary Windows software).

```bash
# Run only SSHD check, skip OSV
landschaft audit --versions=false

# Run OSV check with longer timeout and larger batch
landschaft audit --sshd=false --max 100 --timeout 30s
```

---

### harden

Applies hardening changes to the system. All subcommands support `--plan` to preview what would change without applying anything.

```bash
landschaft harden --plan          # preview all suggested subcommands
landschaft harden <subcommand> --plan
```

#### rotate-pwd

Rotates passwords for local (or domain) users in two steps.

```bash
# Step 1: generate a CSV of usernames and new passwords
landschaft harden rotate-pwd --generate -f passwords.csv

# Step 2: apply passwords from the CSV
landschaft harden rotate-pwd --apply -f passwords.csv
```

| Flag | Default | Description |
|---|---|---|
| `--generate` / `-g` | | Generate new passwords and write to file |
| `--apply` | | Apply passwords from file |
| `-f <path>` | (required) | CSV file path |
| `--length <n>` | 16 | Password length |
| `--blacklist <user,...>` | | Additional users to skip (always skips `root`, `blackteam`, `sync`) |
| `--allowed-chars <str>` | alphanumeric + symbols | Character set for generated passwords |
| `--domain` | | Rotate domain users instead of local (AD on Windows, LDAP on Linux) |

On Linux passwords are applied via `chpasswd`. On Windows via a PowerShell script targeting local or AD users.

#### firewall

Manages inbound and outbound firewall rules from a JSON rule file.

```bash
landschaft harden firewall --inbound  --apply   -f rules.json
landschaft harden firewall --outbound --apply   -f rules.json
landschaft harden firewall --inbound  --remove  -f rules.json
```

| Flag | Description |
|---|---|
| `--inbound` / `--outbound` | Direction of rules |
| `--apply` | Apply rules from file |
| `--remove` | Remove rules defined in file |
| `--backup` | Backup current rules before applying |
| `-f <path>` | JSON rule file |

Rule JSON format (Windows):

```json
[
  { "DisplayName": "Allow SSH", "LocalPort": 22, "Protocol": "TCP", "Action": "Allow", "Direction": "Inbound" }
]
```

#### backup-etc / restore-etc

Backs up and restores `/etc` (Linux only).

```bash
landschaft harden backup-etc
landschaft harden restore-etc
```

Incremental backups are numbered (`filename-1`, `filename-2`, …) and diffs against the previous version are printed on each run.

#### configure-shell

Configures shell command logging (Linux only).

```bash
landschaft harden configure-shell --shell logger --directory /dev/shm/ --backup backup/
```

| Flag | Default | Description |
|---|---|---|
| `--shell` / `-s` | `logger` | Shell type: `logger` or `ssh` |
| `--directory` / `-d` | `/dev/shm/` | Directory to write log files to |
| `--backup` / `-b` | `backup` | Backup directory for modified config files |

- **`logger`**: appends `PROMPT_COMMAND` to `/etc/environment`, `/etc/profile`, and `/etc/bash.bashrc` so every shell command is sent to syslog via `logger`.
- **`ssh`**: appends a `ForceCommand` to `/etc/ssh/sshd_config` that logs every SSH session command.

#### add-local-admin

Adds a new local administrator account and sets its password interactively.

```bash
landschaft harden add-local-admin <username>
```

On Linux: creates the user with `useradd`, adds to `sudo` (or `wheel` as fallback), then runs `passwd`.

#### ssh

Hardens the SSH daemon configuration (Linux only).

```bash
landschaft harden ssh
```

Applies a secure `sshd_config` — disables root login, enforces strong ciphers, disables X11 forwarding and empty passwords.

#### rdp

Hardens Remote Desktop settings (Windows only).

```bash
landschaft harden rdp
```

Enables Network Level Authentication (NLA) and applies secure RDP policy settings.

#### lock-accounts

Disables guest and other high-risk default accounts.

```bash
landschaft harden lock-accounts
```

#### baseline-firewall

Applies a conservative baseline firewall that allows established connections and common scored service ports, and blocks everything else.

```bash
landschaft harden baseline-firewall
```

Uses `ufw`/`iptables` on Linux and Windows Firewall on Windows.

#### enable-auditing

Enables host-level audit logging.

```bash
landschaft harden enable-auditing
```

On Windows: enables the full Windows audit policy (logon events, object access, process tracking, etc.) via `auditpol`.
On Linux: ensures `auditd` is installed and running with a baseline ruleset.

---

### hunt

Searches host logs for suspicious events and writes detections to a JSONL file.

```bash
landschaft hunt [flags]
```

| Flag | Default | Description |
|---|---|---|
| `--since <duration>` | `30m` | How far back to look (e.g. `1h`, `24h`) |
| `--detections-log <path>` | `./landschaft-detections.jsonl` | Output file (or env `LANDSCHAFT_DETECTIONS_LOG`) |

On Linux: reads from `journalctl` and `/var/log/auth.log` (or `/var/log/secure`).
On Windows: reads from the Windows Event Log via `Get-WinEvent`.

Detections are tagged with severity and human-readable explanations for:

| Event | Tag |
|---|---|
| Failed SSH/RDP login | `ssh-fail` / `rdp-fail` |
| Successful login | `login-success` |
| User added to privileged group | `group-add` |
| New service installed | `new-service` |
| Audit log cleared | `audit-cleared` |

Output JSONL format:

```json
{"timestamp":"2025-03-10T12:34:56Z","host":"web01","os":"linux","source":"auth","event_id":"","account":"root","remote_ip":"10.0.0.5","severity":"high","tags":["ssh-fail"],"explain":"Failed SSH login attempt"}
```

---

### scored-services

Lists ports currently listening on the host and explains what protocol/service they serve.

```bash
landschaft scored-services list     # show listening ports
landschaft scored-services explain  # show ports with protocol explanations
```

Useful for verifying that scored services are still up after hardening.

---

### wazuh

Installs Wazuh for centralized log collection and intrusion detection.

#### install-server (Linux only)

Installs the Wazuh manager, registers agents, and saves their authentication keys.

```bash
landschaft wazuh install-server -n 3 -i 10.0.0.11,10.0.0.12,10.0.0.13
```

| Flag | Description |
|---|---|
| `-n` / `--num-agents` | Number of agents to register |
| `-i` / `--ips` | Comma-separated list of agent IP addresses |

Keys are saved to `./agent_keys/agent1_key.txt`, `agent2_key.txt`, etc.

#### install-agent

Installs the Wazuh agent on this host and registers it with the manager.

**Linux:**

```bash
landschaft wazuh install-agent \
  --manager-ip 10.0.0.10 \
  --agent-name web01 \
  --server-user admin \
  --key-dir /home/admin
```

Fetches the pre-generated key from the manager via SCP, installs `wazuh-agent` from the official repo, imports the key, and starts the service.

**Windows:**

```bash
landschaft wazuh install-agent --manager-ip 10.0.0.10 --agent-name win-client1
```

Downloads the Wazuh MSI installer, installs silently with `WAZUH_MANAGER` set, and starts `WazuhSvc`.

| Flag | Description |
|---|---|
| `--manager-ip` | Wazuh manager IP (required) |
| `--agent-name` | Name for this agent (required) |
| `--server-user` | SSH username on manager (Linux only) |
| `--key-dir` | Remote directory containing `agent_keys/` (Linux only) |
| `--wazuh-version` | Agent version to download, default `4.9.2` (Windows only) |

---

### graylog

Installs a Graylog server for log aggregation (Linux only, Docker-based).

#### gen-certs

Generates a self-signed CA and server certificate chain for Graylog TLS.

```bash
landschaft graylog gen-certs
```

Creates in the current directory:
- `ca.crt` / `ca.key` — root CA
- `ca-bundle.key` — CA cert + key (for upload to Graylog setup wizard)
- `graylog.internal.crt` / `graylog.internal.key` — server certificate
- `graylog.internal.bundle.crt` — server cert + CA cert chain

#### install-server

Installs Graylog via Docker Compose, configures TLS, and opens firewall ports.

```bash
landschaft graylog install-server \
  --tls-public-chain graylog.internal.bundle.crt \
  --tls-private-key graylog.internal.key
```

After running, follow the printed instructions to start Graylog with `docker compose up -d`.

---

### ldap

Utilities for managing LDAP / Active Directory.

```bash
landschaft ldap --help
```

Subcommands include LDIF generation for bulk user creation and password management. Run `landschaft ldap --help` for the current subcommand list.

---

### misc

Miscellaneous helper tools.

#### tools (Linux)

Installs `rsyslog` for syslog forwarding.

```bash
landschaft misc tools
```

#### tools (Windows) — sysinternals, firefox, nxlog

```bash
# Download and extract Sysinternals Suite
landschaft misc sysinternals C:\Tools\Sysinternals

# Install Firefox silently
landschaft misc firefox

# Install and configure Nxlog (for log forwarding to Graylog)
landschaft misc nxlog --install --url "https://nxlog.co/system/files/products/files/348/nxlog-ce-3.2.2329.msi"
landschaft misc nxlog --cert path/to/graylog-ca.pem   # load CA cert into Nxlog
landschaft misc nxlog --config path/to/nxlog.conf      # load custom config
```

#### extract

Dumps all embedded scripts from the binary to disk.

```bash
landschaft misc extract ./scripts/
```

---

### serve

Starts a TLS-protected HTTPS file server in a local directory.

```bash
landschaft serve <directory> [--port 8443]
```

- Generates a self-signed TLS certificate automatically (ECDSA P-256, valid 1 year)
- Requires HTTP Basic Auth — a random password is generated and printed at startup, or you can type your own
- Useful for quickly transferring files to/from a compromised host

```bash
landschaft serve /opt/tools --port 9443
```

---

### baseline

Runs a baseline configuration check.

**Linux:** executes `baseline.sh` against the current working directory and reports differences from a known-good state.

**Windows:** enumerates running services and compares against a baseline allowlist, printing any non-default services.

```bash
landschaft baseline
landschaft baseline services   # Windows: service baseline only
```

---

### report

Generates reports from collected data.

#### inject

Produces a Markdown incident report from the action log and triage snapshot.

```bash
landschaft report inject \
  --action-log landschaft-actions.jsonl \
  --triage triage.tsv \
  --out report.md
```

The report includes:
- A summary section
- Table of all actions taken (command, timestamp, duration, exit code)
- Current state snapshot from triage
- Findings and notes section

---

## Action log

Every command automatically appends an entry to `landschaft-actions.jsonl` (or `LANDSCHAFT_ACTION_LOG`). This provides a tamper-evident audit trail of everything the tool has done on the host.

```json
{"timestamp":"2025-03-10T12:00:00Z","hostname":"web01","user":"root","os":"linux","command":"harden","args":["rotate-pwd","--apply","-f","passwords.csv"],"exit_code":0,"duration_ms":1243}
```

Use `landschaft report inject` to turn this log into a readable Markdown report.

---

## Workflow tutorial

This is a recommended sequence for bringing up a box during a competition.

### 1. Drop the binary

```bash
scp landschaft-linux user@10.0.0.11:/tmp/landschaft
ssh user@10.0.0.11 "chmod +x /tmp/landschaft"
```

### 2. Triage the host

```bash
/tmp/landschaft triage
```

Review the output. Copy `triage.tsv` back to your local machine.

### 3. Audit for vulnerabilities

```bash
/tmp/landschaft audit
```

Note any flagged SSHD settings or vulnerable packages.

### 4. Harden

Run in order — each step is reversible via the backup:

```bash
# Backup /etc first
/tmp/landschaft harden backup-etc

# Rotate all local user passwords
/tmp/landschaft harden rotate-pwd --generate -f /tmp/passwords.csv
/tmp/landschaft harden rotate-pwd --apply    -f /tmp/passwords.csv

# Apply SSH hardening
/tmp/landschaft harden ssh

# Apply baseline firewall
/tmp/landschaft harden baseline-firewall

# Lock high-risk default accounts
/tmp/landschaft harden lock-accounts

# Enable auditd
/tmp/landschaft harden enable-auditing

# Set up shell command logging via syslog
/tmp/landschaft harden configure-shell --shell logger
```

Preview any step first with `--plan`:

```bash
/tmp/landschaft harden rotate-pwd --apply -f /tmp/passwords.csv --plan
```

### 5. Install Wazuh (on manager)

```bash
# On the Wazuh manager machine (Linux):
/tmp/landschaft wazuh install-server -n 5 -i 10.0.0.11,10.0.0.12,10.0.0.13,10.0.0.14,10.0.0.15

# On each Linux agent:
/tmp/landschaft wazuh install-agent --manager-ip 10.0.0.10 --agent-name web01 --server-user admin --key-dir /home/admin

# On each Windows agent (run as Administrator):
landschaft.exe wazuh install-agent --manager-ip 10.0.0.10 --agent-name win-client1
```

### 6. Check scored services

```bash
/tmp/landschaft scored-services explain
```

Make sure all expected services are still listening after hardening.

### 7. Hunt for intrusions

```bash
/tmp/landschaft hunt --since 2h
```

Review `landschaft-detections.jsonl` for suspicious activity.

### 8. Generate a report

```bash
/tmp/landschaft report inject \
  --action-log landschaft-actions.jsonl \
  --triage triage.tsv \
  --out incident-report.md
```
