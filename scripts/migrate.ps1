param (
    [Parameter(Position = 0)]
    [string]$Action = "up",

    [Parameter(Position = 1)]
    [string]$Service = "all",

    [string]$HostName = "localhost",
    [int]$Port = 5432
)

$ErrorActionPreference = "Stop"

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host " SELARAS HEALTHCARE ECOSYSTEM - UNIFIED MIGRATION RUNNER  " -ForegroundColor Cyan
Write-Host " Action: $Action | Service: $Service | Host: ${HostName}:$Port" -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan

go run ./scripts/migrator "-action=$Action" "-service=$Service" "-db-host=$HostName" "-db-port=$Port"
