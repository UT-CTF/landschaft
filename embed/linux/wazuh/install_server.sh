#!/bin/bash
set -e

# Usage: install_server.sh <num_agents> <ip1,ip2,...>
# Installs the Wazuh manager and registers agents, saving keys to ./agent_keys/

if [ $# -lt 2 ]; then
    echo "Usage: $0 <num_agents> <ip1,ip2,...>"
    exit 1
fi

NUM_AGENTS="$1"
IP_LIST="$2"
KEY_DIR="./agent_keys"

IFS=',' read -ra IPS <<< "$IP_LIST"

if [ "${#IPS[@]}" -ne "$NUM_AGENTS" ]; then
    echo "[!] Number of IPs (${#IPS[@]}) must match num_agents ($NUM_AGENTS)."
    exit 1
fi

mkdir -p "$KEY_DIR"

# -------------------------------------------------------
# INSTALL WAZUH MANAGER
# -------------------------------------------------------
echo "[+] Updating system packages..."
apt-get update -y
apt-get install -y curl gnupg lsb-release apt-transport-https

echo "[+] Adding Wazuh repository..."
curl -s https://packages.wazuh.com/key/GPG-KEY-WAZUH | gpg --dearmor \
    | tee /usr/share/keyrings/wazuh.gpg > /dev/null
echo "deb [signed-by=/usr/share/keyrings/wazuh.gpg] https://packages.wazuh.com/4.x/apt stable main" \
    | tee /etc/apt/sources.list.d/wazuh.list
apt-get update -y

echo "[+] Installing Wazuh manager..."
apt-get install -y wazuh-manager

echo "[+] Enabling and starting Wazuh manager..."
systemctl enable wazuh-manager
systemctl start wazuh-manager

# Allow agent communication port
if command -v ufw &>/dev/null; then
    ufw allow 1514/tcp
    ufw allow 1514/udp
    ufw allow 1515/tcp
fi

# -------------------------------------------------------
# REGISTER AGENTS
# -------------------------------------------------------
for ((i = 0; i < NUM_AGENTS; i++)); do
    AGENT_NAME="agent$((i + 1))"
    AGENT_IP="${IPS[$i]}"

    echo "[+] Adding agent '$AGENT_NAME' with IP $AGENT_IP..."
    /var/ossec/bin/manage_agents <<EOF
A
$AGENT_NAME
$AGENT_IP
y
Q
EOF

    sleep 1

    echo "[+] Retrieving ID for $AGENT_NAME..."
    AGENT_ID=$(printf "L\nQ\n" | /var/ossec/bin/manage_agents \
        | grep "$AGENT_NAME" | awk -F'[:,]' '{gsub(/ /,"",$2); print $2}')

    if [ -z "$AGENT_ID" ]; then
        echo "[!] Could not find ID for $AGENT_NAME"
        exit 1
    fi

    echo "[+] Exporting key for $AGENT_NAME (ID: $AGENT_ID)..."
    RAW_OUTPUT=$(printf "E\n%s\n\nQ\n" "$AGENT_ID" | /var/ossec/bin/manage_agents)
    AGENT_KEY=$(echo "$RAW_OUTPUT" | grep -E '^[A-Za-z0-9+/=]{20,}$')

    if [ -z "$AGENT_KEY" ]; then
        echo "[!] Failed to extract key for $AGENT_NAME"
        exit 1
    fi

    echo "$AGENT_KEY" > "$KEY_DIR/${AGENT_NAME}_key.txt"
    chmod 600 "$KEY_DIR/${AGENT_NAME}_key.txt"
    echo "[+] Key saved to $KEY_DIR/${AGENT_NAME}_key.txt"
done

echo ""
echo "[+] Wazuh manager installed. $NUM_AGENTS agent(s) registered."
echo "[+] Keys stored in: $KEY_DIR"
echo "[+] Deploy agents using: landschaft wazuh install-agent --manager-ip <this_ip> --agent-name <name> --server-user <user> --key-dir <path>"
