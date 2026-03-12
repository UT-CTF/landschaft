$ErrorActionPreference = "Stop"

Write-Host "[*] Checking for Npcap installation..."
try {
    Get-Service -Name "npcap" -ErrorAction Stop | Out-Null
    Write-Host "[+] Npcap is already installed."
    exit 0
} catch {
    Write-Host "[*] Npcap is not installed."
}

Write-Host "[*] Downloading Npcap..."

$npcapUrl = "https://npcap.com/dist/npcap-1.79.exe"
$installerPath = "$env:TEMP\npcap-installer.exe"

Invoke-WebRequest -Uri $npcapUrl -OutFile $installerPath

Write-Host "[*] Installing Npcap (requires GUI)..."

Start-Process -FilePath $installerPath `
    -ArgumentList "/winpcap_mode=yes", "/loopback_support=yes" `
    -Wait -Verb RunAs

Write-Host "[+] Npcap installation complete."
