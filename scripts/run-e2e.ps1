param (
    [string]$BaseUrl = "http://localhost:8080",
    [switch]$Verbose = $false
)

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host "  Selaras SIMRS - Executing E2E Patient Journey (PowerShell)" -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host "Target Gateway: $BaseUrl" -ForegroundColor Yellow

$VerboseFlag = ""
if ($Verbose) {
    $VerboseFlag = "-verbose"
}

# Run the standalone E2E CLI test runner
go run ./tests/e2e/cmd/main.go -base-url "$BaseUrl" $VerboseFlag
if ($LASTEXITCODE -ne 0) {
    Write-Host "E2E Patient Journey test failed with exit code $LASTEXITCODE" -ForegroundColor Red
    exit $LASTEXITCODE
}
