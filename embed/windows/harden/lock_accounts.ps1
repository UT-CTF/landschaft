# Disable Guest and other high-risk default accounts.
$ErrorActionPreference = "SilentlyContinue"
Disable-LocalUser -Name "Guest" -ErrorAction SilentlyContinue
Write-Host "Disabled Guest account. Consider disabling other default accounts as per policy."
