param (
    [switch]$Apply,
    [string]$RulePath,
    [string]$BackupPath
)

if (-not $Apply) {
    Get-Content firewall_rules.json | Write-Host
    exit 0
}

Write-Host "Backing up current rules to $BackupPath ..."
netsh advfirewall export $BackupPath

$OldRules = Get-NetFirewallRule | ? { $_.Direction -eq "Inbound" }

$rules = Get-Content $RulePath | ConvertFrom-Json

$rules | % {
    Write-Host "Creating rule for $($_.Name) ..."
    $params = @{
        DisplayName = $_.Name
        Name = $_.Name
        Direction = $_.Direction
        Action = $_.Action
        Protocol = $_.Protocol
        LocalPort = $_.LocalPort
        Profile = "Any"
        Enabled = "True"
    }

    if ($_.Program) {
        $params.Program = $_.Program
    }

    # Write-Host "New-NetFirewallRule $params"
    # foreach ($key in $params.Keys) {
    #     Write-Host "$key : $($params[$key])"
    # }
    New-NetFirewallRule @params
}

Write-Host "Removing existing inbound rules ..."
$OldRules | Remove-NetFirewallRule
