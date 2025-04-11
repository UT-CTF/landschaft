param (
    [switch]$Apply,
    [switch]$Prune,
    [string]$RulePath,
    [string]$BackupPath,
    [string]$OldRulesPath,
    [string]$Direction
)

if ($Prune) {
    Write-Host "Removing all old inbound rules ..."
    Get-Content $OldRulesPath | % {
        if ($_.Trim().Length -eq 0) {
            return
        }
        Write-Host "Removing rule with ID: $_"
        try {
            Remove-NetFirewallRule $_
        }
        catch {
            Write-Host "Error removing rule: $($_.Exception.Message)"
        }
    }
    exit 0
}

if (-not $Apply) {
    Get-Content $RulePath | Write-Host
    exit 0
}

try {
    Write-Host "Backing up current rules to $BackupPath ..."
    netsh advfirewall export $BackupPath
}
catch {
    Write-Host "Error backing up rules: $($_.Exception.Message)"
    exit 1
}

try {
    $OldRules = Get-NetFirewallRule | ? { $_.Direction -eq $Direction }
}
catch {
    Write-Host "Error getting current rules: $($_.Exception.Message)"
    exit 1
}

try {
(Get-Content $RulePath | ConvertFrom-Json -ErrorAction Stop) | % {
        $params = @{}
        $_.PSObject.Properties | % {
            $params[$_.Name] = $_.Value
        }
        $params.Enabled = "True"
        $params.Profile = "Any"
        $params.ErrorAction = "Stop"
        Write-Host "Creating rule: $($params.DisplayName)"
        New-NetFirewallRule @params
    }
}
catch {
    Write-Host "Error creating rules: $($_.Exception.Message)"
    exit 1
}

Write-Host "Firewall rules applied successfully."

Write-Host "Writing old rules IDs to $OldRulesPath"
$OldRules | Select-Object InstanceID | % {$_.InstanceID} | Out-File -FilePath $OldRulesPath
