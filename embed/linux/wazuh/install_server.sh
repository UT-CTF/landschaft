#!/bin/bash
set -e

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

if command -v ufw &>/dev/null; then
    ufw allow 1514/tcp
    ufw allow 1514/udp
    ufw allow 1515/tcp
fi

echo "[+] Wazuh manager installed and running."
