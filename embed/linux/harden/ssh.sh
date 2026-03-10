#!/bin/bash
# Harden sshd: ensure PermitRootLogin no and safe defaults. Requires root.
set -e
CFG="${1:-/etc/ssh/sshd_config}"
if [ ! -f "$CFG" ]; then
  echo "No sshd_config at $CFG"
  exit 1
fi
# Backup
cp -a "$CFG" "${CFG}.landschaft.bak"
# Ensure critical directives (append if missing, or use sed to replace)
for pair in "PermitRootLogin no" "PasswordAuthentication no" "PermitEmptyPasswords no"; do
  key="${pair%% *}"
  if grep -qE "^[[:space:]]*#?[[:space:]]*${key}" "$CFG"; then
    sed -i "s/^[[:space:]]*#?[[:space:]]*${key}.*/$pair/" "$CFG"
  else
    echo "$pair" >> "$CFG"
  fi
done
echo "Restart sshd for changes to take effect (e.g. systemctl restart sshd)."
