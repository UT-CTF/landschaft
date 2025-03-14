#!/bin/bash

set -e

# Check if arguments are provided
if [ $# -ne 2 ]; then
    echo "Usage: $0 <remote_log_server_ip> <path_to_ca_cert>"
    echo "Example: $0 192.168.1.10 /path/to/ca.crt"
    exit 1
fi

REMOTE_IP="$1"
CA_CERT_PATH="$2"
GRAYLOG_HOSTNAME="graylog.internal"

# Check if CA certificate exists
if [ ! -f "$CA_CERT_PATH" ]; then
    echo "Error: CA certificate not found at $CA_CERT_PATH"
    exit 1
fi

echo "Setting up rsyslog to forward logs to $GRAYLOG_HOSTNAME ($REMOTE_IP) using CA at $CA_CERT_PATH"

# Add entry to /etc/hosts
echo "Adding $GRAYLOG_HOSTNAME to /etc/hosts..."
# Check if entry already exists and remove it
if grep -q "$GRAYLOG_HOSTNAME" /etc/hosts; then
    sed -i "/$GRAYLOG_HOSTNAME/d" /etc/hosts
fi
echo "$REMOTE_IP $GRAYLOG_HOSTNAME" >> /etc/hosts

# Detect distribution
if [ -f /etc/debian_version ]; then
    echo "Debian-based distribution detected"
    apt-get update
    apt-get install -y rsyslog rsyslog-gnutls
elif [ -f /etc/redhat-release ] || [ -f /etc/fedora-release ]; then
    echo "Red Hat-based distribution detected"
    yum install -y rsyslog
elif [ -f /etc/arch-release ]; then
    echo "Arch Linux detected"
    pacman -Sy --noconfirm rsyslog
elif [ -f /etc/alpine-release ]; then
    echo "Alpine Linux detected"
    apk add --no-cache rsyslog
elif [ -f /etc/SuSE-release ] || [ -f /etc/openSUSE-release ]; then
    echo "SuSE-based distribution detected"
    zypper install -y rsyslog
else
    echo "Distribution not recognized. Attempting to install rsyslog anyway..."
    # Try common package managers
    if command -v apt-get >/dev/null; then
        apt-get update && apt-get install -y rsyslog
    elif command -v yum >/dev/null; then
        yum install -y rsyslog
    elif command -v dnf >/dev/null; then
        dnf install -y rsyslog
    elif command -v zypper >/dev/null; then
        zypper install -y rsyslog
    elif command -v pacman >/dev/null; then
        pacman -Sy --noconfirm rsyslog
    elif command -v apk >/dev/null; then
        apk add --no-cache rsyslog
    else
        echo "Could not find a package manager to install rsyslog"
        exit 1
    fi
fi

# Create directory for certificates
mkdir -p /etc/rsyslog.d/certs/
cp "$CA_CERT_PATH" /etc/rsyslog.d/certs/ca.pem
chmod +r /etc/rsyslog.d/certs/ca.pem

# Create configuration file
echo "Creating rsyslog forwarding configuration..."
cat > /etc/rsyslog.d/10-remote-logging.conf << EOF
# Configure remote forwarding with TLS
module(load="omfwd")

global(DefaultNetstreamDriverCAFile="")

action(type="omfwd" protocol="tcp" target="$GRAYLOG_HOSTNAME" port="5555"
       StreamDriver="gtls" StreamDriverMode="1" StreamDriverAuthMode="anon")
EOF

# Restart rsyslog service
echo "Restarting rsyslog service..."
if command -v systemctl >/dev/null; then
    systemctl restart rsyslog
    systemctl enable rsyslog
elif command -v service >/dev/null; then
    service rsyslog restart
    # Enable on boot for systems without systemd
    if [ -f /etc/debian_version ]; then
        update-rc.d rsyslog enable
    elif [ -f /etc/redhat-release ]; then
        chkconfig rsyslog on
    fi
else
    echo "Warning: Could not restart rsyslog automatically. Please restart it manually."
fi

echo "Rsyslog has been installed and configured successfully to forward logs to $GRAYLOG_HOSTNAME."
