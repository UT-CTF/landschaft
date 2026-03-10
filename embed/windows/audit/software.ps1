$ErrorActionPreference = "SilentlyContinue"

$items = @()

$paths = @(
  "HKLM:\Software\Microsoft\Windows\CurrentVersion\Uninstall\*",
  "HKLM:\Software\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall\*",
  "HKCU:\Software\Microsoft\Windows\CurrentVersion\Uninstall\*"
)

foreach ($p in $paths) {
  Get-ItemProperty $p | ForEach-Object {
    $name = $_.DisplayName
    $ver  = $_.DisplayVersion
    if ([string]::IsNullOrWhiteSpace($name) -or [string]::IsNullOrWhiteSpace($ver)) {
      return
    }
    $items += [pscustomobject]@{
      Name      = $name
      Version   = $ver
      Publisher = $_.Publisher
    }
  }
}

# Include OpenSSH capability state if present (helps in many Windows CCDC images)
try {
  $caps = Get-WindowsCapability -Online | Where-Object { $_.Name -like "OpenSSH*" } | Select-Object Name, State
  foreach ($c in $caps) {
    $items += [pscustomobject]@{
      Name      = $c.Name
      Version   = $c.State
      Publisher = "WindowsCapability"
    }
  }
} catch {}

$items |
  Sort-Object Name, Version -Unique |
  ConvertTo-Json -Depth 3 -Compress

