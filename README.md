# 🏞️ Landschaft

> Changing the CCDC Landscape, one script at a time.

# Landschaft

Landschaft is a cross-platform cybersecurity tool designed for rapid system hardening, triage, and monitoring in CCDC (Collegiate Cyber Defense Competition) environments. It provides teams with essential capabilities to secure systems, detect vulnerabilities, and maintain operational awareness during competitions.

## Features

### System Triage
- [x] **Network Information**: Display DNS names, IP addresses, and open ports with service identification
- [ ] **Software Analysis**: View OS versions and required updates (in development)
- [ ] **Firewall Status**: Check active firewalls and unblocked ports
- [ ] **User Management**: Display user/group information, sudo capabilities, and domain join status (in development)
- [x] **System Information**: View hostname and other basic system details

### System Hardening
- [ ] **Admin Management**: Add local admin users with randomized credentials and pre-seeded SSH keys
- [ ] **Password Management**: Rotate local user passwords with confirmation safeguards
- [ ] **Logging Infrastructure**:
  - [ ] Hidden bash logging for Linux systems
  - [ ] Configurable logging agents with TLS certificate support
  - [ ] Integration with centralized logging (Graylog)
- [ ] **Security Tools**:
  - [ ] Fail2ban configuration
  - [ ] Auditing agents (sysmon for Windows, auditd for Linux)

### Audit Capabilities
- [ ] **SSH Configuration Analysis**: Detect weaknesses in SSH configurations
- [ ] **GPO Analysis**: Identify common Group Policy weaknesses
- [ ] **Vulnerability Scanning**: Check system software for known vulnerabilities
- [ ] **Baseline Comparison**: Create and compare system baselines to detect changes

### LDAP Tools
- [ ] **Password Management**: Generate CSVs with new random passwords for domain users
- [ ] **User Onboarding**: Generate LDIF files from CSV templates
- [ ] **Password Updates**: Create password change LDIFs with proper hashing

### Advanced Features (Planned)
- [ ] **Uptime Monitoring**: Agent for host uptime and tampering detection
- [ ] **Code Signing**: Protection against binary tampering
- [ ] **Secure File Serving**: HTTPS server with self-signed CA and basic authentication
- [ ] **Firewall Management**: TUI for managing firewall settings
- [ ] **Service Analysis**: Detection of non-default services

## Usage

```bash
# Perform initial system triage
landschaft triage

# Harden a system with recommended settings
landschaft harden

# Audit system for vulnerabilities
landschaft audit

# Manage system baselines
landschaft baseline create
landschaft baseline compare

# Generate TLS certificates for secure logging
landschaft gen-cert

# LDAP management tools
landschaft ldap generate-passwords
landschaft ldap generate-ldif
```

## Platform Support

Landschaft is designed to work for only Linux and Windows systems. MacOS support is not planned.

## Installation

### Building from Source

Requirements:
- Go 1.18 or higher

```bash
# Clone the repository
git clone https://github.com/UT-CTF/landschaft.git
cd landschaft

# Run the build script (Linux and Windows binaries)
./build.sh

# Or build manually for your platform
go build -o landschaft
```

## Development

Landschaft is built with Golang using the Cobra CLI framework. Contributions are welcome through pull requests.

Current Project Status:
- [x] Initialize codebase and select tech stack
- [x] Define project structure and create repository
- [ ] Set up CI

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
