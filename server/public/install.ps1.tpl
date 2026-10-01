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

# Tab completion via the PowerShell profile. Profiles do not load under the
# default "Restricted" policy, so allow local scripts for this user.
if ((Get-ExecutionPolicy) -eq "Restricted") {
    try { Set-ExecutionPolicy -Scope CurrentUser RemoteSigned -Force } catch {}
}
$Line = 'if (Get-Command anon -ErrorAction SilentlyContinue) { anon completion powershell | Out-String | Invoke-Expression }'
if (-not (Test-Path $PROFILE)) { New-Item -ItemType File -Force -Path $PROFILE | Out-Null }
if (-not (Select-String -Path $PROFILE -SimpleMatch "anon completion" -Quiet)) {
    Add-Content -Path $PROFILE -Value "`n$Line"
}

Write-Host ""
Write-Host "anon installed!" -ForegroundColor Green
Write-Host "  1. Telegram e bot ke /start pathan (message receive korar jonno)"
Write-Host "  2. Notun terminal khule try korun: anon members"
Write-Host "  3. anon send likhe naam er prothom akkhor diye Tab chapun (PowerShell e @ chara likhun)"
