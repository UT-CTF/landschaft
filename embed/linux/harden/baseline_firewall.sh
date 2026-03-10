#!/bin/bash
# Conservative baseline: allow established, allow SSH (22), and common ports. Requires root.
set -e
if command -v ufw &>/dev/null; then
  ufw default deny incoming
  ufw default allow outgoing
  ufw allow 22/tcp
  ufw allow from 127.0.0.1
  ufw --force enable || true
  echo "UFW baseline applied."
elif command -v firewall-cmd &>/dev/null && systemctl is-active --quiet firewalld 2>/dev/null; then
  firewall-cmd -q --permanent --add-service=ssh 2>/dev/null || true
  firewall-cmd -q --reload 2>/dev/null || true
  echo "Firewalld: ssh service allowed. Add more as needed."
else
  echo "No ufw or firewalld found. Configure iptables or install ufw."
fi
