Get-NetTCPConnection -State Listen | Select-Object LocalPort, OwningProcess, LocalAddress | ConvertTo-Json -Compress
