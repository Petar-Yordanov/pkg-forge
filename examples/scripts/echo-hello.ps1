[CmdletBinding()]
param()

Write-Host "hello from echo-hello.ps1"
Write-Host ("PSVersion=" + $PSVersionTable.PSVersion.ToString())
Write-Host ("PWD=" + (Get-Location).Path)
