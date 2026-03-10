param(
    [Parameter(Mandatory=$true)]
    [string]$ManagerIP,

    [Parameter(Mandatory=$true)]
    [string]$AgentName,

    [string]$WazuhVersion = "4.9.2"
)

$ErrorActionPreference = "Stop"

$InstallerPath = "$env:TEMP\wazuh-agent-$WazuhVersion.msi"
$DownloadUrl = "https://packages.wazuh.com/4.x/windows/wazuh-agent-$WazuhVersion-1.msi"

# -------------------------------------------------------
# DOWNLOAD INSTALLER
# -------------------------------------------------------
Write-Host "[+] Downloading Wazuh agent $WazuhVersion..."
try {
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $InstallerPath -UseBasicParsing
} catch {
    Write-Error "[!] Failed to download Wazuh installer: $_"
    exit 1
}

if (-not (Test-Path $InstallerPath)) {
    Write-Error "[!] Installer not found after download."
    exit 1
}

# -------------------------------------------------------
# INSTALL
# -------------------------------------------------------
Write-Host "[+] Installing Wazuh agent '$AgentName' (manager: $ManagerIP)..."
$msiArgs = @(
    "/i", "`"$InstallerPath`"",
    "/q",
    "WAZUH_MANAGER=`"$ManagerIP`"",
    "WAZUH_AGENT_NAME=`"$AgentName`"",
    "WAZUH_REGISTRATION_SERVER=`"$ManagerIP`""
)

$process = Start-Process -FilePath "msiexec.exe" -ArgumentList $msiArgs -Wait -PassThru
if ($process.ExitCode -ne 0) {
    Write-Error "[!] Installation failed with exit code: $($process.ExitCode)"
    exit 1
}

# -------------------------------------------------------
# START SERVICE
# -------------------------------------------------------
Write-Host "[+] Starting Wazuh agent service..."
try {
    Start-Service -Name "WazuhSvc"
    Set-Service -Name "WazuhSvc" -StartupType Automatic
} catch {
    Write-Error "[!] Failed to start Wazuh service: $_"
    exit 1
}

$svc = Get-Service -Name "WazuhSvc" -ErrorAction SilentlyContinue
if ($svc -and $svc.Status -eq "Running") {
    Write-Host "[+] Wazuh agent '$AgentName' is running and connected to $ManagerIP"
} else {
    Write-Warning "[!] Wazuh service may not be running. Check: Get-Service WazuhSvc"
}

# Clean up installer
Remove-Item -Path $InstallerPath -Force -ErrorAction SilentlyContinue
