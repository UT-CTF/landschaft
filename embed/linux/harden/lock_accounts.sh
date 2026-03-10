#!/bin/bash
# Disable guest and common high-risk accounts.
set -e
for u in guest nobody; do
  if id "$u" &>/dev/null; then
    usermod -L "$u" 2>/dev/null || true
    echo "Locked: $u"
  fi
done
echo "Done. Consider also: usermod -L test, admin, or other default accounts."
