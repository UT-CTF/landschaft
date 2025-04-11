#!/bin/bash
# iptables-setup.sh
# Hardened firewall rules for Linux servers
# Customize as needed by commenting/uncommenting lines

set -e

# Initial Setup
echo "[+] Installing iptables and iptables-persistent (if not already installed)..."
apt-get update
apt-get install -y iptables
apt-get install -y iptables-persistent

echo "[+] Flushing existing rules..."
iptables -F
iptables -X
iptables -Z

ip6tables -F
ip6tables -X
ip6tables -Z

# ANTI LOCKOUT RULE: WHITELIST SSH FROM UR IP
# COMMENT THIS OUT AT YOUR OWN RISK!!!!!!!!! BUT ALSO I HAVE BACKUP FLUSH ALL RULES SO ITS NOT THAT RISKY TBH
read -p "Enter the IP address to whitelist for SSH. Make sure u type this in right because I can't be bothered to do regex to check. " SSH_IP
# if nothing is entered, skip whitelist
if [[ -z "$SSH_IP" ]]; then
  echo "[!] No IP entered. Skipping SSH allow rule."
else
  echo "[+] Allowing SSH from $SSH_IP..."
  iptables -A INPUT -p tcp -s "$SSH_IP" --dport 22 -j ACCEPT
fi

# setting default politcies
echo "[+] Setting default policies to DROP... (dont worry there is anti lockout)"
iptables -P INPUT DROP
iptables -P FORWARD DROP
iptables -P OUTPUT DROP

ip6tables -P INPUT DROP
ip6tables -P FORWARD DROP
ip6tables -P OUTPUT DROP

echo "[+] allowing loopback traffic..."
iptables -A INPUT -i lo -j ACCEPT
iptables -A OUTPUT -o lo -j ACCEPT

echo "[+] handling connection tracking..."
iptables -A INPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
iptables -A INPUT -m conntrack --ctstate INVALID -j DROP
iptables -A OUTPUT -m conntrack --ctstate ESTABLISHED -j ACCEPT

echo "[+] allowing DHCP (client)..."
iptables -A INPUT -p udp --dport 68 --sport 67 -j ACCEPT
iptables -A OUTPUT -p udp --dport 67 --sport 68 -j ACCEPT

echo "[+] allowing DNS resolution..."
iptables -A OUTPUT -p udp --dport 53 -j ACCEPT
iptables -A OUTPUT -p tcp --dport 53 -j ACCEPT

echo "[+] allowing NTP (time sync)..."
iptables -A OUTPUT -p udp --dport 123 -j ACCEPT
iptables -A OUTPUT -p tcp --dport 123 -j ACCEPT

# ---------- incoming traffic UNCOMMENT AS NEEDED----------
#
# echo "[+] allowing SSH from subnet or ip"
# iptables -A INPUT -p tcp -s X.X.X.X/Y --dport 22 -j ACCEPT

# echo "[+] allowing HTTP and HTTPS (Web Server)"
# iptables -A INPUT -p tcp --dport 80 -j ACCEPT # HTTP
# iptables -A INPUT -p tcp --dport 443 -j ACCEPT # HTTPS

# echo "[+] allowing MySQL"
# iptables -A INPUT -p tcp -s X.X.X.X/Y --dport 3306 -j ACCEPT # from subnet
# iptables -A INPUT -p tcp --dport 3306 -j ACCEPT # from anywhere

# echo "[+] allowing LDAP (unencrypted)"
# iptables -A INPUT -p tcp --dport 389 -j ACCEPT

# echo "[+] allowing LDAPS (encrypted)"
# iptables -A INPUT -p tcp --dport 636 -j ACCEPT

# ---------- outgoing traffic ----------
# Add what outbound services your system needs

echo "[+] allowing outbound web traffic (apt, curl, etc.)..."
iptables -A OUTPUT -p tcp --dport 80 -j ACCEPT
iptables -A OUTPUT -p tcp --dport 443 -j ACCEPT
iptables -A OUTPUT -p udp --dport 443 -j ACCEPT  # For HTTP/3 (QUIC)

# echo "[+] allowing outbound SMTP for sending email (choose one)"
# iptables -A OUTPUT -p tcp --dport 25 -j ACCEPT    # SMTP
# iptables -A OUTPUT -p tcp --dport 465 -j ACCEPT   # SMTPS
# iptables -A OUTPUT -p tcp --dport 587 -j ACCEPT   # Submission

# echo "[+] allowing outbound LDAP / AD"
# iptables -A OUTPUT -p tcp --dport 389 -j ACCEPT
# iptables -A OUTPUT -p tcp --dport 636 -j ACCEPT

# echo "[+] allowing outbound git"
# iptables -A OUTPUT -p tcp --dport 9418 -j ACCEPT

# echo "[+] allowing outbound FTP (command/control only)"
# iptables -A OUTPUT -p tcp --dport 21 -j ACCEPT

# ---------- saving rules ----------
echo "[+] Saving rules..."
netfilter-persistent save

echo "[!] Setting iptables flush in 5 minutes as backup..."
(sleep 300 && iptables -F && echo "[+] Auto-flush: iptables rules have been cleared after 5 minutes.") &
echo "[!] To cancel the flush, run kill <pid>. You can find the PID by running ps aux| grep 'sleep 300"

echo "[✓] Firewall setup complete. Review rules and reboot if needed."

