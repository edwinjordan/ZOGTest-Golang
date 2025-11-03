# OpenTelemetry Monitoring Quick Start Script for Windows PowerShell

Write-Host "=== ZOGTest-Golang OpenTelemetry Setup ===" -ForegroundColor Cyan
Write-Host ""

# Check if Docker is running
Write-Host "Checking Docker..." -ForegroundColor Yellow
try {
    $dockerInfo = docker info 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Host "ERROR: Docker is not running." -ForegroundColor Red
        Write-Host ""
        Write-Host "Please start Docker Desktop and wait for it to initialize, then try again." -ForegroundColor Yellow
        Write-Host ""
        Write-Host "To start Docker Desktop:" -ForegroundColor White
        Write-Host "1. Open Docker Desktop application" -ForegroundColor Gray
        Write-Host "2. Wait for the whale icon in system tray to turn green" -ForegroundColor Gray
        Write-Host "3. Run this script again" -ForegroundColor Gray
        Write-Host ""
        pause
        exit 1
    }
    Write-Host "✓ Docker is running" -ForegroundColor Green
}
catch {
    Write-Host "ERROR: Docker is not installed or not accessible." -ForegroundColor Red
    Write-Host "Please install Docker Desktop for Windows." -ForegroundColor Yellow
    exit 1
}

# Check if .env exists
if (-Not (Test-Path ".env")) {
    Write-Host "Creating .env file from .env.example..." -ForegroundColor Yellow
    Copy-Item ".env.example" ".env"
    Write-Host "✓ .env file created. Please review and update as needed." -ForegroundColor Green
} else {
    Write-Host "✓ .env file exists" -ForegroundColor Green
}

# Start observability stack
Write-Host ""
Write-Host "Starting observability stack (Jaeger, Prometheus, Grafana, OTel Collector)..." -ForegroundColor Yellow
docker-compose up -d

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Observability stack started successfully!" -ForegroundColor Green
    
    Write-Host ""
    Write-Host "=== Services Available ===" -ForegroundColor Cyan
    Write-Host "Jaeger UI:           http://localhost:16686" -ForegroundColor White
    Write-Host "Prometheus:          http://localhost:9090" -ForegroundColor White
    Write-Host "Grafana:             http://localhost:3000 (admin/admin)" -ForegroundColor White
    Write-Host "OTel Collector:      localhost:4317 (gRPC)" -ForegroundColor White
    Write-Host "Collector Health:    http://localhost:13133" -ForegroundColor White
    
    Write-Host ""
    Write-Host "Waiting for services to be ready..." -ForegroundColor Yellow
    Start-Sleep -Seconds 5
    
    # Check if services are healthy
    Write-Host "Checking service health..." -ForegroundColor Yellow
    $otelHealth = Invoke-WebRequest -Uri "http://localhost:13133" -UseBasicParsing -ErrorAction SilentlyContinue
    if ($otelHealth.StatusCode -eq 200) {
        Write-Host "✓ OpenTelemetry Collector is healthy" -ForegroundColor Green
    } else {
        Write-Host "⚠ OpenTelemetry Collector may not be ready yet" -ForegroundColor Yellow
    }
    
    Write-Host ""
    Write-Host "=== Next Steps ===" -ForegroundColor Cyan
    Write-Host "1. Start your application:" -ForegroundColor White
    Write-Host "   go run main.go" -ForegroundColor Gray
    Write-Host ""
    Write-Host "2. Make some requests to your API:" -ForegroundColor White
    Write-Host "   Invoke-WebRequest http://localhost:8000" -ForegroundColor Gray
    Write-Host ""
    Write-Host "3. View traces in Jaeger:" -ForegroundColor White
    Write-Host "   http://localhost:16686" -ForegroundColor Gray
    Write-Host ""
    Write-Host "4. View metrics in Prometheus:" -ForegroundColor White
    Write-Host "   http://localhost:9090" -ForegroundColor Gray
    Write-Host ""
    Write-Host "5. Create dashboards in Grafana:" -ForegroundColor White
    Write-Host "   http://localhost:3000" -ForegroundColor Gray
    Write-Host ""
    
    Write-Host "To stop all services, run: docker-compose down" -ForegroundColor Yellow
    Write-Host ""
    
} else {
    Write-Host "ERROR: Failed to start observability stack" -ForegroundColor Red
    Write-Host "Check docker-compose logs: docker-compose logs" -ForegroundColor Yellow
    exit 1
}
