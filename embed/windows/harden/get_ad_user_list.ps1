Get-ADUser -Filter {Enabled -eq $true} | % { Write-Host $_.Name }
