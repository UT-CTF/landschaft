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
        Name                              = $_.Name
        SamAccountName                    = $_.SamAccountName
        Enabled                           = $_.Enabled
        Groups                            = $groupNames -join ';'
        AccountNotDelegated               = $_.AccountNotDelegated
        AllowReversiblePasswordEncryption = $_.AllowReversiblePasswordEncryption
        DoesNotRequirePreAuth             = $_.DoesNotRequirePreAuth
        PasswordNeverExpires              = $_.PasswordNeverExpires
        PasswordNotRequired               = $_.PasswordNotRequired
        TrustedForDelegation              = $_.TrustedForDelegation
        TrustedToAuthForDelegation        = $_.TrustedToAuthForDelegation
    }
} | Sort-Object Name | Export-Csv "$BaselinePath\ad-users.csv" -NoTypeInformation
