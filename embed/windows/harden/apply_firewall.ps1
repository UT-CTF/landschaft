param(
    [switch]$ClearScheduledTask,
    [string]$RulesFile,
    [string]$BackupFile,
    [string]$Direction
)

$scheduledTaskName = "LSCWFW - Restore Firewall Rules"

if ($ClearScheduledTask) {
    Unregister-ScheduledTask -TaskName $scheduledTaskName -Confirm:$false -ErrorAction SilentlyContinue
    Write-Host "Cleared scheduled task: $scheduledTaskName"
    exit 0
}

try {
netsh advfirewall export $BackupFile
# create a scheduled task to restore the firewall rules in 10 minutes
$action = New-ScheduledTaskAction -Execute "netsh" -Argument "advfirewall import $BackupFile"
$trigger = New-ScheduledTaskTrigger -Once -At (Get-Date).AddMinutes(5)
$principal = New-ScheduledTaskPrincipal -UserId "SYSTEM" -LogonType ServiceAccount
Register-ScheduledTask -TaskName $scheduledTaskName -Action $action -Trigger $trigger -Principal $principal
} catch {
    Write-Host "Error creating scheduled backup task: $($_.Exception.Message)"
    exit 1
}

$oldRules = Get-NetFirewallRule -Direction $Direction

try {
    (Get-Content $RulesFile | ConvertFrom-Json -ErrorAction Stop) | % {
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
} catch {
    Write-Host "Error applying firewall rules: $($_.Exception.Message)"
    exit 1
}

Write-Host "Successfully added new firewall rules"
Write-Host "Removing old firewall rules ..."
try {
    $oldRules | Remove-NetFirewallRule
} catch {
    Write-Host "Error removing old rules: $($_.Exception.Message)"
    exit 1
}
Write-Host "Successfully removed old firewall rules"
Write-Host "Firewall rules updated successfully. If you need to restore the previous rules, a scheduled task has been created to do so in 5 minutes. You can also manually run 'netsh advfirewall import $BackupFile' to restore the backup immediately."
