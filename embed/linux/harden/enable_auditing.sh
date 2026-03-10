#!/bin/bash
# Ensure auth logging is present (rsyslog/journald). Optionally enable auditd.
set -e
if command -v systemctl &>/dev/null; then
  systemctl is-active --quiet systemd-journald && echo "journald is active"
  systemctl is-active --quiet rsyslog 2>/dev/null && echo "rsyslog is active" || true
fi
if command -v auditctl &>/dev/null; then
  auditctl -s 2>/dev/null || echo "auditd not running; consider: systemctl enable --now auditd"
else
  echo "auditd not installed; auth.log/journal will still capture SSH and sudo."
fi
