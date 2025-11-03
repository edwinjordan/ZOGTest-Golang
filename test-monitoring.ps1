# Test OpenTelemetry Integration

Write-Host "=== Testing OpenTelemetry Integration ===" -ForegroundColor Cyan
Write-Host ""

# Function to check if a URL is accessible
function Test-Endpoint {
    param(
        [string]$Url,
        [string]$Name
    )
    
    try {
        $response = Invoke-WebRequest -Uri $Url -UseBasicParsing -TimeoutSec 5 -ErrorAction Stop
        Write-Host "✓ $Name is accessible (Status: $($response.StatusCode))" -ForegroundColor Green
        return $true
    }
    catch {
        Write-Host "✗ $Name is NOT accessible" -ForegroundColor Red
        return $false
    }
}

# Check observability services
Write-Host "Checking observability services..." -ForegroundColor Yellow
Write-Host ""

$jaegerOk = Test-Endpoint "http://localhost:16686" "Jaeger UI"
$prometheusOk = Test-Endpoint "http://localhost:9090" "Prometheus"
$grafanaOk = Test-Endpoint "http://localhost:3000" "Grafana"
$otelHealthOk = Test-Endpoint "http://localhost:13133" "OTel Collector Health"

Write-Host ""

# Check if application is running
Write-Host "Checking application..." -ForegroundColor Yellow
$appRunning = Test-Endpoint "http://localhost:8000" "Application API"

if ($appRunning) {
    Write-Host ""
    Write-Host "Sending test requests to generate telemetry data..." -ForegroundColor Yellow
    
    # Make multiple requests to generate data
    for ($i = 1; $i -le 10; $i++) {
        try {
            Invoke-WebRequest -Uri "http://localhost:8000" -UseBasicParsing -ErrorAction SilentlyContinue | Out-Null
            Write-Host "." -NoNewline -ForegroundColor Gray
        }
        catch {
            # Ignore errors
        }
        Start-Sleep -Milliseconds 100
    }
    
    Write-Host ""
    Write-Host "✓ Test requests sent" -ForegroundColor Green
    
    # Check metrics endpoint
    Write-Host ""
    Write-Host "Checking metrics endpoint..." -ForegroundColor Yellow
    $metricsOk = Test-Endpoint "http://localhost:8000/metrics" "Prometheus Metrics Endpoint"
    
    if ($metricsOk) {
        Write-Host ""
        Write-Host "Sample metrics:" -ForegroundColor Cyan
        $metrics = Invoke-WebRequest -Uri "http://localhost:8000/metrics" -UseBasicParsing
        $metricsLines = ($metrics.Content -split "`n" | Select-String "^http_" | Select-Object -First 5)
        foreach ($line in $metricsLines) {
            Write-Host "  $line" -ForegroundColor Gray
        }
    }
}
else {
    Write-Host ""
    Write-Host "⚠ Application is not running. Start it with: go run main.go" -ForegroundColor Yellow
}

# Summary
Write-Host ""
Write-Host "=== Summary ===" -ForegroundColor Cyan
Write-Host ""

if ($jaegerOk -and $prometheusOk -and $grafanaOk -and $otelHealthOk) {
    Write-Host "✓ All monitoring services are running!" -ForegroundColor Green
    
    if ($appRunning) {
        Write-Host "✓ Application is running and instrumented!" -ForegroundColor Green
        Write-Host ""
        Write-Host "Next steps:" -ForegroundColor White
        Write-Host "1. Open Jaeger to view traces: http://localhost:16686" -ForegroundColor Gray
        Write-Host "2. Open Prometheus to query metrics: http://localhost:9090" -ForegroundColor Gray
        Write-Host "3. Open Grafana to create dashboards: http://localhost:3000" -ForegroundColor Gray
    }
    else {
        Write-Host "⚠ Application needs to be started" -ForegroundColor Yellow
        Write-Host "Run: go run main.go" -ForegroundColor Gray
    }
}
else {
    Write-Host "⚠ Some services are not available" -ForegroundColor Yellow
    Write-Host "Try running: docker-compose up -d" -ForegroundColor Gray
}

Write-Host ""
