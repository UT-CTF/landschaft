param(
    [switch]$ClearScheduledTask,
    [string]$RulesFile,
    [string]$BackupFile,
    [string]$Direction
)

$scheduledTaskName = "LSCWFW - Restore Firewall Rules"

if ($ClearScheduledTask) {
    try {
        Unregister-ScheduledTask -TaskName $scheduledTaskName -Confirm:$false -ErrorAction Stop
        Write-Host "Cleared scheduled task: $scheduledTaskName"
    }
    catch {
        Write-Host "Error clearing scheduled task: $($_.Exception.Message)"
        exit 1
    }
    exit 0
}

try {
    if(Test-Path $BackupFile) {
        Write-Host "Backup file already exists at $BackupFile. Please remove it before running this script to avoid overwriting the backup."
        exit 1
    }
    netsh advfirewall export $BackupFile | Out-Null
    Write-Host "Successfully backed up current firewall rules to $BackupFile"
    
    if (Get-ScheduledTask -TaskName $scheduledTaskName -ErrorAction SilentlyContinue) {
        Unregister-ScheduledTask -TaskName $scheduledTaskName -Confirm:$false
    }

    $action = New-ScheduledTaskAction -Execute "netsh" -Argument "advfirewall import $BackupFile" -ErrorAction Stop
    $timeTrigger = New-ScheduledTaskTrigger -Once -At (Get-Date).AddMinutes(3) -ErrorAction Stop
    $rebootTrigger = New-ScheduledTaskTrigger -AtStartup -ErrorAction Stop
    $principal = New-ScheduledTaskPrincipal -UserId "SYSTEM" -LogonType ServiceAccount -ErrorAction Stop
    Register-ScheduledTask -TaskName $scheduledTaskName -Action $action -Trigger $timeTrigger, $rebootTrigger -Principal $principal -ErrorAction Stop | Out-Null
    Write-Host "Scheduled task created to restore firewall rules in 3 minutes or on next reboot: $scheduledTaskName"
}
catch {
    Write-Host "Error creating scheduled backup task: $($_.Exception.Message)"
    exit 1
}

$oldRules = Get-NetFirewallRule | Where-Object { $_.Direction -eq $Direction }

try {
    (Get-Content $RulesFile | ConvertFrom-Json -ErrorAction Stop) | ForEach-Object {
        $params = @{}
        $_.PSObject.Properties | ForEach-Object {
            $params[$_.Name] = $_.Value
        }
        $params.Enabled = "True"
        $params.Profile = "Any"
        $params.ErrorAction = "Stop"
        Write-Host "Creating rule: $($params.DisplayName)"
        New-NetFirewallRule @params | Out-Null
    }
}
catch {
    Write-Host "Error applying firewall rules: $($_.Exception.Message)"
    exit 1
}

Write-Host "Successfully added new firewall rules"
Write-Host "Removing old firewall rules ..."
try {
    $oldRules | Remove-NetFirewallRule -ErrorAction Stop
}
catch {
    Write-Host "Error removing old rules: $($_.Exception.Message)"
    exit 1
}
Write-Host "Successfully removed old firewall rules"
Write-Host "Firewall rules updated successfully. If you need to restore the previous rules, a scheduled task has been created to do so in 3 minutes. You can also manually run 'netsh advfirewall import $BackupFile' to restore the backup immediately."
