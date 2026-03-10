# Enable key audit policies for CCDC visibility (logon failures, account management). Requires admin.
$ErrorActionPreference = "Stop"
auditpol /set /subcategory:"Logon" /failure:enable /success:enable
auditpol /set /subcategory:"Account Management" /failure:enable /success:enable
auditpol /set /subcategory:"Privilege Use" /failure:enable /success:enable
Write-Host "Audit policy updated. Security log will record logon, account management, and privilege use."
