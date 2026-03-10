param(
    [switch]$Outbound
)

$inboundRulesFile = "firewall_rules_inbound.json"
$outboundRulesFile = "firewall_rules_outbound.json"

if ($Outbound) {
    $openPorts = @()
    $rulesFile = $outboundRulesFile
}
else {
    $openPorts = (Get-NetTCPConnection -State Listen | 
        Where-Object { $_.LocalAddress -notin @('127.0.0.1', '::1') } |
        Select-Object -Unique LocalPort).LocalPort
    $rulesFile = $inboundRulesFile
}

(Get-Content $rulesFile | ConvertFrom-Json -ErrorAction Stop) | 
Select-Object *, @{Name = "Enabled"; Expression = { $_.LocalPort -in $openPorts } } |
Sort-Object @{Expression = "Enabled"; Desc = $true }, @{Expression = "LocalPort"; Descending = $false } |
ConvertTo-Json -Depth 10 
