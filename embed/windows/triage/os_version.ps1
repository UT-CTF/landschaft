$sysinfo = systeminfo
$sysinfo | ? { $_.StartsWith("OS Name") -or $_.StartsWith("OS Version") -or $_.StartsWith("OS Config") } | % { Write-Host $_ }
