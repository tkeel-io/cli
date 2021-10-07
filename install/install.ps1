# ------------------------------------------------------------
# Copyright 2021 The TKeel Contributors.
# Licensed under the Apache License.
# ------------------------------------------------------------
param (
    [string]$Version,
    [string]$TKeelRoot = "c:\tkeel"
)

Write-Output ""
$ErrorActionPreference = 'stop'

#Escape space of TKeelRoot path
$TKeelRoot = $TKeelRoot -replace ' ', '` '

# Constants
$TKeelCliFileName = "tkeel.exe"
$TKeelCliFilePath = "${TKeelRoot}\${TKeelCliFileName}"

# GitHub Org and repo hosting TKeel CLI
$GitHubOrg = "tkeel-io"
$GitHubRepo = "cli"

# Set Github request authentication for basic authentication.
if ($Env:GITHUB_USER) {
    $basicAuth = [System.Convert]::ToBase64String([System.Text.Encoding]::ASCII.GetBytes($Env:GITHUB_USER + ":" + $Env:GITHUB_TOKEN));
    $githubHeader = @{"Authorization" = "Basic $basicAuth" }
}
else {
    $githubHeader = @{}
}

if ((Get-ExecutionPolicy) -gt 'RemoteSigned' -or (Get-ExecutionPolicy) -eq 'ByPass') {
    Write-Output "PowerShell requires an execution policy of 'RemoteSigned'."
    Write-Output "To make this change please run:"
    Write-Output "'Set-ExecutionPolicy RemoteSigned -scope CurrentUser'"
    break
}

# Change security protocol to support TLS 1.2 / 1.1 / 1.0 - old powershell uses TLS 1.0 as a default protocol
[Net.ServicePointManager]::SecurityProtocol = "tls12, tls11, tls"

# Check if TKeel CLI is installed.
if (Test-Path $TKeelCliFilePath -PathType Leaf) {
    Write-Warning "TKeel is detected - $TKeelCliFilePath"
    Invoke-Expression "$TKeelCliFilePath --version"
    Write-Output "Reinstalling TKeel..."
}
else {
    Write-Output "Installing TKeel..."
}

# Create TKeel Directory
Write-Output "Creating $TKeelRoot directory"
New-Item -ErrorAction Ignore -Path $TKeelRoot -ItemType "directory"
if (!(Test-Path $TKeelRoot -PathType Container)) {
    throw "Cannot create $TKeelRoot"
}

# Get the list of release from GitHub
$releases = Invoke-RestMethod -Headers $githubHeader -Uri "https://api.github.com/repos/${GitHubOrg}/${GitHubRepo}/releases" -Method Get
if ($releases.Count -eq 0) {
    throw "No releases from github.com/tkeel-io/cli repo"
}

# Filter windows binary and download archive
if (!$Version) {
    $windowsAsset = $releases | Where-Object { $_.tag_name -notlike "*rc*" } | Select-Object -First 1 | Select-Object -ExpandProperty assets | Where-Object { $_.name -Like "*windows_amd64.zip" }
    if (!$windowsAsset) {
        throw "Cannot find the windows TKeel CLI binary"
    }
    $zipFileUrl = $windowsAsset.url
    $assetName = $windowsAsset.name
} else {
    $assetName = "tkeel_windows_amd64.zip"
    $zipFileUrl = "https://github.com/${GitHubOrg}/${GitHubRepo}/releases/download/v${Version}/${assetName}"
}

$zipFilePath = $TKeelRoot + "\" + $assetName
Write-Output "Downloading $zipFileUrl ..."

$githubHeader.Accept = "application/octet-stream"
Invoke-WebRequest -Headers $githubHeader -Uri $zipFileUrl -OutFile $zipFilePath
if (!(Test-Path $zipFilePath -PathType Leaf)) {
    throw "Failed to download TKeel Cli binary - $zipFilePath"
}

# Extract TKeel CLI to $TKeelRoot
Write-Output "Extracting $zipFilePath..."
Microsoft.Powershell.Archive\Expand-Archive -Force -Path $zipFilePath -DestinationPath $TKeelRoot
if (!(Test-Path $TKeelCliFilePath -PathType Leaf)) {
    throw "Failed to download TKeel Cli archieve - $zipFilePath"
}

# Check the TKeel CLI version
Invoke-Expression "$TKeelCliFilePath --version"

# Clean up zipfile
Write-Output "Clean up $zipFilePath..."
Remove-Item $zipFilePath -Force

# Add TKeelRoot directory to User Path environment variable
Write-Output "Try to add $TKeelRoot to User Path Environment variable..."
$UserPathEnvironmentVar = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($UserPathEnvironmentVar -like '*tkeel*') {
    Write-Output "Skipping to add $TKeelRoot to User Path - $UserPathEnvironmentVar"
}
else {
    [System.Environment]::SetEnvironmentVariable("PATH", $UserPathEnvironmentVar + ";$TKeelRoot", "User")
    $UserPathEnvironmentVar = [Environment]::GetEnvironmentVariable("PATH", "User")
    Write-Output "Added $TKeelRoot to User Path - $UserPathEnvironmentVar"
}

Write-Output "`r`nTKeel CLI is installed successfully."
Write-Output "To get started with TKeel, please visit https://docs.tkeel.io/getting-started/ ."
Write-Output "Ensure that Docker Desktop is set to Linux containers mode when you run TKeel in self hosted mode."
