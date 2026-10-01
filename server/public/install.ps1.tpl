# anon installer for Windows.
# Usage (PowerShell): irm __SERVER__/install.ps1 | iex
$ErrorActionPreference = "Stop"

$Server = "__SERVER__"
$Dir = Join-Path $env:LOCALAPPDATA "anon"
$Arch = if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "amd64" }

Write-Host "Downloading anon (windows/$Arch)..."
New-Item -ItemType Directory -Force -Path $Dir | Out-Null
Invoke-WebRequest -UseBasicParsing -Uri "$Server/dl/anon-windows-$Arch.exe" -OutFile (Join-Path $Dir "anon.exe")

# Add install dir to the user PATH if it is not there yet.
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if (($UserPath -split ";") -notcontains $Dir) {
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$Dir", "User")
}
$env:Path = "$env:Path;$Dir"

$Key = $env:ANON_KEY
if (-not $Key) { $Key = Read-Host "Team key" }
& (Join-Path $Dir "anon.exe") config $Server $Key | Out-Null

Write-Host ""
Write-Host "anon installed!" -ForegroundColor Green
Write-Host "  1. Telegram e bot ke /start pathan (message receive korar jonno)"
Write-Host "  2. Notun terminal khule try korun: anon members"
