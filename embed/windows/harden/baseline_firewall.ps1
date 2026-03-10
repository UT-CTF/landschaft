# Conservative baseline: allow established, allow RDP (3389) and common management. Requires admin.
$ErrorActionPreference = "SilentlyContinue"
$rules = Get-NetFirewallRule -PolicyStore ActiveStore | Where-Object { $_.Enabled -eq 'True' }
Write-Host "Current firewall rules count: $($rules.Count). Use landschaft harden firewall --inbound --apply -f <json> for custom baseline."
