# Output Security events (logon, failed logon, group add, service install). JSON array. Arg: minutes.
$min = if ($args[0]) { [int]$args[0] } else { 30 }
if ($min -lt 1) { $min = 30 }
$start = (Get-Date).AddMinutes(-$min)
$ids = 4624, 4625, 4728, 4732, 4756, 7045, 1102
$events = @(Get-WinEvent -FilterHashtable @{ LogName = 'Security'; Id = $ids; StartTime = $start } -MaxEvents 200 -ErrorAction SilentlyContinue)
if ($events.Count -eq 0) { Write-Output '[]'; exit }
$events | Select-Object TimeCreated, Id, Message, LogName | ConvertTo-Json -Compress
