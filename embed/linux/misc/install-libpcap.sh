#!/usr/bin/env bash
set -e

echo "[*] Detecting distribution..."

if [ -f /etc/debian_version ]; then
    echo "[*] Debian/Ubuntu detected"
    apt update
    apt install -y libpcap0.8 libpcap-dev

elif [ -f /etc/redhat-release ]; then
    echo "[*] RHEL/CentOS/Fedora detected"
    yum install -y libpcap libpcap-devel || dnf install -y libpcap libpcap-devel

elif [ -f /etc/arch-release ]; then
    echo "[*] Arch Linux detected"
    pacman -Sy --noconfirm libpcap

else
    echo "[!] Unknown distro. Please install libpcap manually."
    exit 1
fi

echo "[*] Verifying installation..."
ldconfig -p | grep libpcap || echo "[!] libpcap not found in ldconfig cache"

echo "[+] libpcap installation complete."
echo "[!] Reminder: Packet capture requires root or CAP_NET_RAW capability."
