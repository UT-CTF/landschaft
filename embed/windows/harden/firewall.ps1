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


# $params = @{
#     Test='hello'
#     Second='bye'
# }
# $rules[0].PSObject.Properties | % { Write-Host $_.Value}
# Write-Host $rules[0].GetType()
# Write-Host $params.GetType()

# try {
#     $rules | % {
#         Write-Host "Creating rule for $($_.Name) ..."
#         $_ | % {
#             Write-Host $_
#         }
#         $params = @{
#             DisplayName = $_.Name
#             Direction   = $_.Direction
#             Action      = $_.Action
#             Protocol    = $_.Protocol
#             LocalPort   = $_.LocalPort
#             Profile     = "Any"
#             Enabled     = "True"
#             ErrorAction = "Stop"
#         }

#         if ($_.Program) {
#             $params.Program = $_.Program
#         }

#         # Write-Host "New-NetFirewallRule $params"
#         # foreach ($key in $params.Keys) {
#         #     Write-Host "$key : $($params[$key])"
#         # }
#         New-NetFirewallRule @params
#     }
# }
# catch {
#     Write-Host "Error creating rule: $($_.Exception.Message)"
#     exit 1
# }

Write-Host "Firewall rules applied successfully."

Write-Host "Writing old rules IDs to $OldRulesPath"
$OldRules | Select-Object InstanceID | % {$_.InstanceID} | Out-File -FilePath $OldRulesPath

# Write-Host "Removing existing inbound rules ..."
# $OldRules | Remove-NetFirewallRule
