#!/bin/bash
# List running service unit names (one per line) for triage non-default comparison.
systemctl list-units --type=service --state=running --no-legend --no-pager -o 'unit' 2>/dev/null | sed 's/\.service$//' || true
