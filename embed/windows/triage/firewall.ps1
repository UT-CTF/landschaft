Write-Host "Interfaces:"
Get-NetConnectionProfile | Format-Table -AutoSize -Property Name, InterfaceAlias, NetworkCategory

Get-NetFirewallProfile | % { 
    if ($_.Enabled) {
        Write-Host "Profile: $($_.Name) - Enabled" 
    } else {
        Write-Host "Profile: $($_.Name) - Disabled" 
    }
}

Write-Host

netsh advfirewall show allprofiles | % {
    if ($_.Contains("Profile Settings")) {
        Write-Host $_
    } elseif ($_.Contains("Firewall Policy") -or $_.Contains("State")) {
        Write-Host "`t$_"
    }
}
