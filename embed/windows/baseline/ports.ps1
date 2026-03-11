param(
    [Parameter(Mandatory = $true)]
    [string]$BaselinePath
)

if (-not (Test-Path $BaselinePath)) {
    New-Item -ItemType Directory -Path $BaselinePath | Out-Null
}

Get-NetTCPConnection -State Listen |
    Select-Object LocalAddress, LocalPort, @{Name = "Process"; Expression = {(Get-Process -Id $_.OwningProcess -ea 0).ProcessName}} |
    Where-Object {$_.LocalAddress -notin @("127.0.0.1","::1")} |
    Sort-Object LocalPort |
    Export-Csv "$BaselinePath\ports.csv" -NoTypeInformation
