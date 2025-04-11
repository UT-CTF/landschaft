#!/bin/bash
# ufw-setup.sh
# Hardened firewall rules using UFW (Uncomplicated Firewall)

set -e

echo "[+] Installign UFW if not already installed.."
apt-get update
apt-get install -y ufw

echo "[+] Resetting UFW to default..."
ufw --force reset

# ANTI LOCKOUT RULE: WHITELIST SSH FROM YOUR IP
read -p "Enter the IP address to whitelist for SSH. Make sure u type this in right becuase idk regex so im not gonna check. " SSH_IP
if [[ -z "$SSH_IP" ]]; then
  echo "[!] No IP entered. Skipping SSH allow rule."
else
  echo "[+] Allowing SSH from $SSH_IP..."
  ufw allow from "$SSH_IP" to any port 22 proto tcp
fi

# Default deny everything
echo "[+] Setting default policies to deny..."
ufw default deny incoming
ufw default deny outgoing
ufw default deny routed

# Allow loopback
echo "[+] Allowing loopback..."
ufw allow in on lo
ufw allow out on lo

# Allow established & related connections (UFW does this automatically)

# ---------- incoming traffic ----------
#
# echo "[+] allowing HTTP and HTTPS"
# ufw allow 80/tcp # HTTP
# ufw allow 443/tcp # HTTPS
#
# echo "[+] allowing SSH from subnet"
# ufw allow from X.X.X.X/Y to any port 22 proto tcp
#
# echo "[+] allowing MySQL"
# ufw allow from X.X.X.X/Y to any port 3306 proto tcp # from subnet
# ufw allow 3306/tcp # from subnet
#
#
# echo "[+] allowing LDAP / secure LDAPS"
# ufw allow 389/tcp
# ufw allow 636/tcp

# ---------- outgoing traffic ----------

echo "[+] allowing DHCP client..."
ufw allow out 67/udp
ufw allow out 68/udp

echo "[+] allowing DNS resolution..."
ufw allow out 53/tcp
ufw allow out 53/udp

echo "[+] allowing NTP time sync..."
ufw allow out 123/udp
ufw allow out 123/tcp

echo "[+] allowing outbound web traffic..."
ufw allow out 80/tcp
ufw allow out 443/tcp
ufw allow out 443/udp

# echo "[+] allowing outbound SMTP (only do one i think?????)..."
# ufw allow out 25/tcp
# ufw allow out 465/tcp
# ufw allow out 587/tcp

# echo "[+] allowing outbound LDAP / AD"
# ufw allow out 389/tcp
# ufw allow out 636/tcp

# echo "[+] allowing outbound Git"
# ufw allow out 9418/tcp

# echo "[+] allwoing outbound FTP (control only because data transfer is random port)"
# ufw allow out 21/tcp
ufw allow out 50000:51000/tcp # Passive data ports for specific range DOUBLE CHECK THIS SAHANA

# ---------- enable and trhe firewall ----------
echo "[!] Enabling UFW..."
ufw --force enable

echo "[!] Scheduling UFW reset in 5 minutes as anti-lockout..."
echo "[!] Setting iptables flush in 5 minutes as backup..."
(sleep 300 && ufw --force reset && echo "[+] Auto-flush: iptables rules have been cleared after 5 minutes.") &
echo "[!] To cancel the flush, run kill <pid>. You can find the PID by running ps aux| grep 'sleep 300"


echo "[✓] UFW firewall setup complete. Use 'ufw status numbered' to review rules."

