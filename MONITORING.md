# 🔭 OpenTelemetry Monitoring

## Overview

This application is fully instrumented with OpenTelemetry for comprehensive observability:
- **Distributed Tracing** via Jaeger
- **Metrics Collection** via Prometheus  
- **Visualization** via Grafana
- **Telemetry Pipeline** via OpenTelemetry Collector

## Quick Start

### 1️⃣ Start Monitoring Stack

```powershell
# Windows PowerShell
.\start-monitoring.ps1

# Or using Make
make monitoring-start

# Or manually
docker-compose up -d
```

### 2️⃣ Configure Environment

```powershell
# Copy example environment file
Copy-Item .env.example .env

# Edit .env and set:
# ENABLE_INSTRUMENTATION=true
# OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
```

### 3️⃣ Run Application

```powershell
go run main.go
```

### 4️⃣ Test Setup

```powershell
.\test-monitoring.ps1
```

## 📊 Access Monitoring Tools

| Tool | URL | Purpose |
|------|-----|---------|
| **Application** | http://localhost:8000 | API Server |
| **Swagger UI** | http://localhost:8000/swagger/index.html | API Documentation |
| **Metrics Endpoint** | http://localhost:8000/metrics | Prometheus Metrics |
| **Jaeger UI** | http://localhost:16686 | Distributed Tracing |
| **Prometheus** | http://localhost:9090 | Metrics & Queries |
| **Grafana** | http://localhost:3000 | Dashboards (admin/admin) |

## 🎯 What's Monitored

### HTTP Requests
- Request count by method, route, status
- Response time (histogram with percentiles)
- In-flight requests
- Error rates

### Database Operations
- Query count by operation and table
- Query duration
- Connection pool metrics
- Success/failure rates

### Business Metrics
- User operations (create, read, update, delete)
- Topic operations
- News operations

### System Metrics
- Application uptime
- Memory usage
- Goroutines
- GC statistics

## 📈 Sample Queries

### Prometheus

```promql
# Request rate
rate(http_server_requests_total[5m])

# 95th percentile latency
histogram_quantile(0.95, rate(http_server_request_duration_bucket[5m]))

# Error rate
rate(http_server_requests_total{http_status_code=~"5.."}[5m])
```

### Jaeger

1. Open http://localhost:16686
2. Select service: `zogtest-golang-api`
3. Click "Find Traces"
4. Explore trace details

## 🛠️ Configuration

### Environment Variables

```bash
# Enable/disable instrumentation
ENABLE_INSTRUMENTATION=true

# OpenTelemetry Collector endpoint
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317

# Service name for telemetry
SERVICE_NAME=zogtest-golang-api

# Sampling rate (0.0 to 1.0)
TRACING_SAMPLE_RATE=1.0

# Application environment
APP_ENVIRONMENT=development
```

### Sampling Strategies

**Development (100% sampling)**
```bash
APP_ENVIRONMENT=development
TRACING_SAMPLE_RATE=1.0
```

**Production (Reduced sampling)**
```bash
APP_ENVIRONMENT=production
TRACING_SAMPLE_RATE=0.1  # 10% sampling
```

## 📚 Documentation

- **[OPENTELEMETRY_SETUP.md](OPENTELEMETRY_SETUP.md)** - Detailed setup guide
- **[IMPLEMENTATION_SUMMARY.md](IMPLEMENTATION_SUMMARY.md)** - Implementation details

## 🔧 Commands

```powershell
# Start monitoring
.\start-monitoring.ps1

# Test setup
.\test-monitoring.ps1

# Stop monitoring
docker-compose down

# View logs
docker-compose logs -f otel-collector

# Check service health
docker-compose ps
```

## 🏗️ Architecture

```
Application → OpenTelemetry Collector → Jaeger (Traces)
                                      → Prometheus (Metrics)
                                      
Grafana ← Prometheus + Jaeger
```

## 📦 Components

- **OpenTelemetry SDK** - Instrumentation library
- **OTel Collector** - Telemetry data pipeline
- **Jaeger** - Distributed tracing backend
- **Prometheus** - Time-series metrics database
- **Grafana** - Visualization and dashboards

## 🚀 Features

✅ Automatic HTTP request tracing  
✅ Database query instrumentation  
✅ Custom metrics support  
✅ Custom span creation  
✅ Context propagation  
✅ Prometheus metrics export  
✅ Pre-configured dashboards  
✅ Health checks  

## 🐛 Troubleshooting

**No traces in Jaeger?**
- Check `ENABLE_INSTRUMENTATION=true`
- Verify OTel Collector is running
- Check logs: `docker-compose logs otel-collector`

**No metrics in Prometheus?**
- Check targets: http://localhost:9090/targets
- Verify metrics endpoint: http://localhost:8000/metrics

**Build errors?**
```powershell
go mod tidy
go build -v .
```

## 📖 Resources

- [OpenTelemetry Documentation](https://opentelemetry.io/docs/)
- [Jaeger Documentation](https://www.jaegertracing.io/docs/)
- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
