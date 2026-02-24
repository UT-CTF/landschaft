param(
    [Parameter(Mandatory = $true)]
    [string]$BaselinePath
)

if (-not (Test-Path $BaselinePath)) {
    New-Item -ItemType Directory -Path $BaselinePath | Out-Null
}

try {
    Get-Process -IncludeUserName | 
    Select-Object Name, Path, UserName | 
    Sort-Object Name | 
    Export-Csv "$BaselinePath\processes.csv" -NoTypeInformation
}
catch {
    Get-Process | ForEach-Object {
        try {
            $owner = (Get-WmiObject -Class Win32_Process -Filter "ProcessId = $($_.Id)").GetOwner()
            $user = "$($owner.Domain)\$($owner.User)"
        }
        catch { $user = "N/A" }
        $path = try { $_.Path } catch { "N/A" }
        [PSCustomObject]@{
            Name = $_.Name
            Path = $path
            User = $user
        }
    } | Sort-Object Name | Export-Csv "$BaselinePath\processes.csv" -NoTypeInformation 
}
