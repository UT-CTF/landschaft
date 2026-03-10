#!/bin/bash
# Output recent journal lines (high priority or sshd) for hunt. One line per event.
journalctl -p warning..alert --since "30 min ago" -o short 2>/dev/null | head -100
journalctl _SYSTEMD_UNIT=sshd.service --since "30 min ago" -o short 2>/dev/null | head -50
