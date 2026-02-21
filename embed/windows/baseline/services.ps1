param(
    [Parameter(Mandatory = $true)]
    [string]$BaselinePath
)

if (-not (Test-Path $BaselinePath)) {
    New-Item -ItemType Directory -Path $BaselinePath | Out-Null
}

Get-Service | Select-Object Name, DisplayName, Status, StartType | Sort-Object Name | Export-Csv "$BaselinePath\services.csv" -NoTypeInformation
