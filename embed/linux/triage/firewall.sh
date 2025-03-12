#!/bin/bash

# Script to print allowed firewall ports cross-platform on Linux
# Supports iptables, firewalld and ufw

check_iptables() {
  if command -v iptables &>/dev/null; then
    echo "=== IPTABLES RULES ==="
    echo "INPUT Chain:"
    iptables -L INPUT -n | grep -E 'ACCEPT.*dpt:[0-9]+'
    echo
    echo "FORWARD Chain:"
    iptables -L FORWARD -n | grep -E 'ACCEPT.*dpt:[0-9]+'
    return 0
  fi
  return 1
}

check_firewalld() {
  if command -v firewall-cmd &>/dev/null && systemctl is-active --quiet firewalld; then
    echo "=== FIREWALLD RULES ==="
    echo "Active zones: $(firewall-cmd --get-active-zones | grep -v '^[[:space:]]')"
    echo "Open ports:"
    firewall-cmd --list-all | grep ports
    return 0
  fi
  return 1
}

check_ufw() {
  if command -v ufw &>/dev/null && ufw status &>/dev/null; then
    echo "=== UFW RULES ==="
    ufw status | grep -E '(ALLOW|DENY)'
    return 0
  fi
  return 1
}

found_firewall=0
higher_level_firewall=0

# Check for higher-level firewalls first
if check_firewalld; then
  found_firewall=1
  higher_level_firewall=1
fi

if check_ufw; then
  found_firewall=1
  higher_level_firewall=1
fi

# Only check iptables if no higher-level firewall was found
if [ $higher_level_firewall -eq 0 ]; then
  if check_iptables; then
    found_firewall=1
  fi
fi

if [ $found_firewall -eq 0 ]; then
  echo "No supported firewall (iptables, firewalld, ufw) detected or permissions insufficient."
  echo "This script requires root privileges to view firewall rules."
fi
