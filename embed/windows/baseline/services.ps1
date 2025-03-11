param (
    [Parameter(Mandatory = $true)]
    [String]$ExportPath
)
function Get-HashBackup($filePath, $hashType = "SHA256") {
    $hasher = [System.Security.Cryptography.HashAlgorithm]::Create($hashType)
    $stream = [System.IO.File]::OpenRead($filePath)
    $hashBytes = $hasher.ComputeHash($stream)
    $stream.Close()
    return -join ($hashBytes | ForEach-Object { "{0:X2}" -f $_ })
}
function Get-ServiceInfoList {
    $serviceList = Get-WmiObject -Class Win32_Service | Select-Object Name, StartMode, State, PathName, DisplayName

    $outList = @()

    foreach ($service in $serviceList) {
        $pathName = $service.PathName
        if ($null -eq $pathName) {
            $pathName = "N/A"
        }
        $exePath = $pathName
        if ($exePath -match '"([^"]+)"') {
            $exePath = $Matches[1]
        }
        else {
            $spaceInd = $exePath.IndexOf(' ')
            if ($spaceInd -gt 0) {
                $exePath = $exePath.Substring(0, $spaceInd)
            }
        }

        if ($exePath -eq "N/A") {
            $tmpHash = "N/A"
        }
        else {
            try {
                $tmpHash = (Get-FileHash -Path $exePath -Algorithm SHA256 -ErrorAction SilentlyContinue).Hash
            }
            catch {
                try {
                    $tmpHash = Get-HashBackup -filePath $exePath
                }
                catch {
                    $tmpHash = "N/A"
                }
            }
            if ($null -eq $tmpHash) {
                Write-Host "Error hashing $exePath"
                $tmpHash = "N/A"
            }
        }
        $outList += [PSCustomObject]@{
            Name           = $service.Name
            DisplayName    = $service.DisplayName
            StartMode      = $service.StartMode
            State          = $service.State
            Path           = $pathName
            Executable     = $exePath
            ExecutableHash = $tmpHash
        }
    }

    return $outList
}

Get-ServiceInfoList | ConvertTo-Csv -NoTypeInformation | Set-Content -Path $ExportPath
Write-Host "Wrote service baseline to $ExportPath"
