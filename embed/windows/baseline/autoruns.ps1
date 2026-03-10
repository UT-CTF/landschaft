param(
    [Parameter(Mandatory = $true)]
    [string]$BaselinePath,
    [Parameter(Mandatory = $true)]
    [string]$SysinternalsPath
)

if (-not (Test-Path $BaselinePath)) {
    New-Item -ItemType Directory -Path $BaselinePath | Out-Null
}

& "$SysinternalsPath\autorunsc64.exe" -accepteula -a * -x * -h -nobanner > "$BaselinePath\autoruns.xml"
$xml = [xml](Get-Content "$BaselinePath\autoruns.xml")
$xml.autoruns.item | ForEach-Object {
    [PSCustomObject]@{
        Location = $_.location
        Name = $_.itemname
        Enabled = $_.enabled
        Profile = $_.profile
        LaunchString = $_.launchstring
        Description = $_.description
        Company = $_.company
        ImagePath = $_.imagepath
        Hash = $_.sha256hash
    }
} | Sort-Object Location, Name | Export-Csv "$BaselinePath\autoruns.csv" -NoTypeInformation
Remove-Item "$BaselinePath\autoruns.xml" -Force
