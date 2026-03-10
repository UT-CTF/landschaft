# List running Windows service names (one per line) for triage non-default comparison.
Get-Service | Where-Object { $_.Status -eq 'Running' } | ForEach-Object { $_.Name }
