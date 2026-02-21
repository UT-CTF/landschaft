param(
    [Parameter(Mandatory = $true)]
    [string]$BaselinePath
)

if (-not (Test-Path $BaselinePath)) {
    New-Item -ItemType Directory -Path $BaselinePath | Out-Null
}

Import-Module ActiveDirectory
Get-ADUser -Filter * -Properties * | ForEach-Object {
    $groupNames = $_.MemberOf | ForEach-Object { (Get-ADGroup $_).Name } | Sort-Object
    [PSCustomObject]@{
        Name           = $_.Name
        SamAccountName = $_.SamAccountName
        Enabled        = $_.Enabled
        Groups         = $groupNames -join ';'
    }
} | Sort-Object Name | Export-Csv "$BaselinePath\users.csv" -NoTypeInformation
