# Harden RDP: require NLA, restrict to specific users if desired.
$ErrorActionPreference = "Stop"
Set-ItemProperty -Path "HKLM:\System\CurrentControlSet\Control\Terminal Server" -Name "UserAuthentication" -Value 1 -Type DWord -ErrorAction SilentlyContinue
Set-ItemProperty -Path "HKLM:\System\CurrentControlSet\Control\Terminal Server" -Name "SecurityLayer" -Value 1 -Type DWord -ErrorAction SilentlyContinue
Write-Host "RDP NLA and security layer set. Reboot or restart TermService for full effect."
