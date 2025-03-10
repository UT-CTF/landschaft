Get-LocalUser | ? { $_.Enabled } | % { Write-Host $_.Name }
