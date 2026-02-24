param (
    [string]$RulesFile
)

$portInfo = Get-NetTCPConnection -State Listen | 
Where-Object { $_.LocalAddress -notin @('127.0.0.1', '::1') } |
Sort-Object LocalPort |
Select-Object -Unique LocalPort, @{Name = 'ProcessName'; Expression = { (Get-Process -Id $_.OwningProcess).Name } } 
$openPorts = $portInfo.LocalPort

$prebuiltRulesFile = "firewall_rules_inbound.json"

try {
    $loadedRules = Get-Content $prebuiltRulesFile | ConvertFrom-Json
}
catch {
    Write-Host "Error loading prebuilt rules: $($_.Exception.Message)"
    exit 1
}

# Find all prebuilt rules that match the open ports
$matchingRules = $loadedRules | Where-Object { $_.LocalPort -in $openPorts }
$unmatchedPorts = $openPorts | Where-Object { $_ -notin $matchingRules.LocalPort }

$matchingRules | ConvertTo-Json -Depth 5 | Out-File $RulesFile -Encoding UTF8
if ($unmatchedPorts.Count -gt 0) {
    Write-Host "Warning: The following open ports do not have matching prebuilt rules and will not be included in the auto-generated firewall rules file:"
    $unmatchedPorts | ForEach-Object { Write-Host " - Port $_" }
}

# map of common windows ports to service name
# $portMap = @{
#     22 = @('SSH', 'TCP')
#     25 = @('SMTP', 'TCP')
#     53 = @('DNS', 'BOTH')
#     80 = @('HTTP', 'TCP')
#     88 = @('Kerberos', 'BOTH')
#     110 = @('POP3', 'TCP')
#     123 = @('NTP', 'TCP')
#     135 = @('RPC', 'TCP')
#     139 = @('NetBIOS', 'TCP')
#     143 = @('IMAP', 'TCP')
#     389 = @('LDAP', 'BOTH')
#     443 = @('HTTPS', 'TCP')
#     445 = @('SMB', 'TCP')
#     465 = @('SMTPS', 'TCP')
#     587 = @('SMTP', 'TCP')
#     636 = @('LDAPS', 'TCP')
#     993 = @('IMAPS', 'TCP')
#     995 = @('POP3S', 'TCP')
#     3268 = @('LDAP GC', 'TCP')
#     3269 = @('LDAPS GC', 'TCP')
#     3389 = @('RDP', 'TCP')
#     5985 = @('WinRM', 'TCP')
#     5986 = @('WinRM SSL', 'TCP')
#     6001 = @('MAPI', 'TCP')
# }
