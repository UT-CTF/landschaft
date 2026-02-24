param(
    [Parameter(Mandatory = $true)]
    [string]$BaselinePath
)

if (-not (Test-Path $BaselinePath)) {
    New-Item -ItemType Directory -Path $BaselinePath | Out-Null
}

Import-Module ActiveDirectory
Get-ADObject -Filter * -Properties * | ForEach-Object {
    [PSCustomObject]@{
        Name           = $_.Name
        DistinguishedName = $_.DistinguishedName
        ObjectClass    = $_.ObjectClass
    }
} | Sort-Object Name | Export-Csv "$BaselinePath\ad-objects.csv" -NoTypeInformation
