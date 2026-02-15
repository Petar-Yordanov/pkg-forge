$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $PSScriptRoot
$examples = Join-Path $root "examples"

$files = Get-ChildItem $examples -File -Recurse | Where-Object {
  $_.Extension -in @(".yml", ".yaml")
}

if ($files.Count -eq 0) {
  Write-Host "No example files found in $examples"
  exit 0
}

$failed = 0

foreach ($f in $files) {
  Write-Host "==> $($f.FullName)"
  & go run . plan -c $f.FullName | Out-Null
  if ($LASTEXITCODE -ne 0) {
    $failed++
  }
}

if ($failed -ne 0) {
  Write-Host "FAILED: $failed file(s)"
  exit 1
}

Write-Host "OK: $($files.Count) file(s)"
