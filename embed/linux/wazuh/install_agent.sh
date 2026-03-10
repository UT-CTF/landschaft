#!/bin/bash
set -e

# Usage: install_agent.sh <agent_name> <manager_ip> <server_user> <remote_key_dir>
# Installs the Wazuh agent and registers it with the manager using a pre-generated key.

if [ $# -lt 4 ]; then
    echo "Usage: $0 <agent_name> <manager_ip> <server_user> <remote_key_dir>"
    echo "Example: $0 agent1 192.168.1.10 admin /home/admin"
    exit 1
fi

AGENT_NAME="$1"
MANAGER_IP="$2"
SERVER_USER="$3"
REMOTE_KEY_DIR="$4"
KEY_DIR="$HOME/agent_keys"
REMOTE_KEY_PATH="${REMOTE_KEY_DIR}/agent_keys/${AGENT_NAME}_key.txt"

mkdir -p "$KEY_DIR"
KEY_FILE="$KEY_DIR/${AGENT_NAME}_key.txt"

# -------------------------------------------------------
# FETCH KEY FROM MANAGER
# -------------------------------------------------------
echo "[+] Fetching key for $AGENT_NAME from $MANAGER_IP..."
scp "${SERVER_USER}@${MANAGER_IP}:${REMOTE_KEY_PATH}" "$KEY_FILE"

if [ ! -f "$KEY_FILE" ]; then
    echo "[!] Failed to retrieve key file."
    exit 1
fi

AGENT_KEY=$(<"$KEY_FILE")

# -------------------------------------------------------
# INSTALL WAZUH AGENT
# -------------------------------------------------------
echo "[+] Adding Wazuh repository..."
apt-get update -y
apt-get install -y curl gnupg apt-transport-https

curl -s https://packages.wazuh.com/key/GPG-KEY-WAZUH | gpg --dearmor \
    | tee /usr/share/keyrings/wazuh.gpg > /dev/null
echo "deb [signed-by=/usr/share/keyrings/wazuh.gpg] https://packages.wazuh.com/4.x/apt stable main" \
    | tee /etc/apt/sources.list.d/wazuh.list
apt-get update -y

echo "[+] Installing Wazuh agent..."
WAZUH_MANAGER="$MANAGER_IP" apt-get install -y wazuh-agent

# -------------------------------------------------------
# REGISTER WITH MANAGER
# -------------------------------------------------------
echo "[+] Registering agent with manager..."
systemctl stop wazuh-agent || true

/var/ossec/bin/manage_agents <<EOF
I
$AGENT_KEY
y
Q
EOF

# Set manager IP in config
sed -i "s|MANAGER_IP|$MANAGER_IP|g" /var/ossec/etc/ossec.conf

# -------------------------------------------------------
# ENABLE AND START
# -------------------------------------------------------
systemctl enable wazuh-agent
systemctl start wazuh-agent

echo "[+] Wazuh agent '$AGENT_NAME' registered and running (manager: $MANAGER_IP)"
