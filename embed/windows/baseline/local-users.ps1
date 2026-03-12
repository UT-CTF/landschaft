param(
    [Parameter(Mandatory = $true)]
    [string]$BaselinePath
)

if (-not (Test-Path $BaselinePath)) {
    New-Item -ItemType Directory -Path $BaselinePath | Out-Null
}

Get-LocalUser | Select-Object Name, Enabled, Description | Sort-Object Name | Export-Csv "$BaselinePath\local-users.csv" -NoTypeInformation
