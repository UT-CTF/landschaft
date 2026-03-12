param(
    [Parameter(Mandatory = $true)]
    [string]$BaselinePath
)

$startupSubkeys = @(
    "Software\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run",
    "Software\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run32",
    "Software\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\StartupFolder"
)

$results = @()

function Get-StartupApproved {
    param(
        [string]$Root,
        [string]$UserSid,
        [string]$Username = $null
    )

    # Write-Host "Parsing StartupApproved entries under $Root for user $UserSid..." -ForegroundColor Cyan
    if (-not $Username) {
        try {
            $Username = (New-Object System.Security.Principal.SecurityIdentifier($UserSid)).
                        Translate([System.Security.Principal.NTAccount]).Value
        }
        catch {
            $Username = "Unknown"
        }
    }
    # Write-Host "Username: $Username" -ForegroundColor Green

    foreach ($sub in $startupSubkeys) {

        $path = "$Root\$sub"

        # Write-Host "Checking $path..." -ForegroundColor Yellow

        if (!(Test-Path $path)) { 
            # Write-Host "Path not found: $path" -ForegroundColor Red
            continue 
        }

        $props = Get-ItemProperty -Path $path

        foreach ($prop in $props.PSObject.Properties | Where-Object {$_.MemberType -eq "NoteProperty"}) {

            $name = $prop.Name
            $data = $prop.Value

            if ($data -is [byte[]] -and $data.Length -gt 0) {

                $statusByte = $data[0]

                if ($statusByte % 2 -eq 0) {
                    $status = "Enabled"
                }
                else {
                    $status = "Disabled"
                }

                $timestamp = $null
                if ($data.Length -ge 12) {
                    try {
                        $filetime = [BitConverter]::ToInt64($data,4)
                        $timestamp = [DateTime]::FromFileTimeUtc($filetime)
                    }
                    catch {}
                }

                $script:results += [PSCustomObject]@{
                    UserSID     = $UserSid
                    Username    = $Username
                    Name        = $name
                    Location    = $path
                    StatusByte  = ('0x{0:X2}' -f $statusByte)
                    Status      = $status
                    LastChange  = $timestamp
                }
            }
        }
    }
}

# --- Enumerate all user profiles ---
$profiles = Get-ItemProperty "HKLM:\Software\Microsoft\Windows NT\CurrentVersion\ProfileList\*"

foreach ($prof in $profiles) {

    $sid = $prof.PSChildName
    $profilePath = $prof.ProfileImagePath
    $ntuser = Join-Path $profilePath "NTUSER.DAT"

    if (!(Test-Path $ntuser)) { continue }

    $hiveRoot = "Registry::HKEY_USERS\$sid"
    $hiveLoaded = Test-Path $hiveRoot
    $tempLoaded = $false

    if (!$hiveLoaded) {
        reg load "HKU\$sid" $ntuser | Out-Null
        $tempLoaded = $true
    }

    Get-StartupApproved -Root $hiveRoot -UserSid $sid

    if ($tempLoaded) {
        reg unload "HKU\$sid" | Out-Null
    }
}

# --- Machine-wide StartupApproved ---
$machineRoot = "Registry::HKEY_LOCAL_MACHINE"

Get-StartupApproved -Root $machineRoot -UserSid "SYSTEM" -Username "SYSTEM"

$results | Sort-Object Username, Name | Select-Object Username, Name, StatusByte, Status | Export-Csv "$BaselinePath\startup-status.csv" -NoTypeInformation
