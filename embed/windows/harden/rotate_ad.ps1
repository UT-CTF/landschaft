param(
    [Parameter(Mandatory = $true)]
    [string]$Path
)

# file from path line by line
$lines = Get-Content -Path $Path
foreach ($line in $lines) {
    $parts = $line.Split(",")
    $user = $parts[0]
    $pass = $parts[1]
    $secstring = ConvertTo-SecureString $pass -AsPlainText -Force
    Set-ADAccountPassword -Identity $user -NewPassword $secstring
    Write-Host "Rotated password for $user"
}
