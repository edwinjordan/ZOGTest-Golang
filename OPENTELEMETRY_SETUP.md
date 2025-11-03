# OpenTelemetry Monitoring Setup

This project includes comprehensive observability using OpenTelemetry for traces, metrics, and logs.

## Architecture

```
┌─────────────────┐
│  Go Application │
│  (Port 8000)    │
└────────┬────────┘
         │ OTLP (gRPC/HTTP)
         ▼
┌─────────────────────┐
│ OTel Collector      │
│ Ports: 4317, 4318   │
└──┬──────────────┬───┘
   │              │
   │ Traces       │ Metrics
   ▼              ▼
┌─────────┐  ┌─────────────┐
│ Jaeger  │  │ Prometheus  │
│ :16686  │  │ :9090       │
└────┬────┘  └──────┬──────┘
     │              │
     └──────┬───────┘
            ▼
     ┌─────────────┐
     │  Grafana    │
     │  :3000      │
     └─────────────┘
```

## Components

### 1. OpenTelemetry Collector
- **Ports**: 
  - 4317 (OTLP gRPC)
  - 4318 (OTLP HTTP)
  - 8888 (Metrics about collector)
  - 8889 (Prometheus exporter)
- **Purpose**: Receives telemetry data from the application and exports to backends

### 2. Jaeger
- **Port**: 16686 (UI)
- **Purpose**: Distributed tracing visualization
- **Access**: http://localhost:16686

### 3. Prometheus
- **Port**: 9090
- **Purpose**: Metrics collection and storage
- **Access**: http://localhost:9090

### 4. Grafana
- **Port**: 3000
- **Purpose**: Unified observability dashboard
- **Access**: http://localhost:3000
- **Credentials**: admin/admin

## Quick Start

### 1. Start the Observability Stack

```powershell
# Start all services
docker-compose up -d

# Check services are running
docker-compose ps

# View logs
docker-compose logs -f otel-collector
```

### 2. Configure Your Application

Copy the example environment file:
```powershell
Copy-Item .env.example .env
```

Edit `.env` and set:
```bash
ENABLE_INSTRUMENTATION=true
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
SERVICE_NAME=zogtest-golang-api
```

### 3. Run Your Application

```powershell
go run main.go
```

## Available Endpoints

### Application
- **API**: http://localhost:8000
- **Swagger**: http://localhost:8000/swagger/index.html
- **Metrics**: http://localhost:8000/metrics

### Monitoring
- **Jaeger UI**: http://localhost:16686
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000
- **OTel Collector Health**: http://localhost:13133

## Features

### Distributed Tracing
- Automatic HTTP request tracing
- Custom span creation in business logic
- Trace context propagation across services
- Database query tracing

### Metrics Collection
- HTTP request metrics (count, duration, in-flight)
- Database query metrics
- Business operation metrics (user, topic, news operations)
- Custom application metrics

### What's Instrumented

1. **HTTP Requests**
   - Method, path, status code
   - Request/response duration
   - In-flight requests

2. **Database Operations**
   - Query execution time
   - Connection pool metrics
   - Operation success/failure

3. **Business Metrics**
   - User operations (create, read, update, delete)
   - Topic operations
   - News operations

## Usage Examples

### View Traces in Jaeger
1. Open http://localhost:16686
2. Select service: `zogtest-golang-api`
3. Click "Find Traces"
4. Explore trace details and spans

### Query Metrics in Prometheus
1. Open http://localhost:9090
2. Try these queries:
   ```promql
   # HTTP request rate
   rate(http_server_requests_total[5m])
   
   # Request duration 95th percentile
   histogram_quantile(0.95, http_server_request_duration_bucket)
   
   # Database query rate
   rate(db_queries_total[5m])
   ```

### Create Dashboards in Grafana
1. Open http://localhost:3000 (admin/admin)
2. Add Prometheus as data source (already configured)
3. Create dashboard with panels:
   - HTTP request rate
   - Response time percentiles
   - Error rate
   - Database performance

## Custom Metrics

Add custom metrics in your code:

```go
import (
    "github.com/edwinjordan/ZOGTest-Golang.git/internal/metrics"
)

// In your handler/service
func (s *Service) DoSomething(ctx context.Context) error {
    // Record custom metric
    s.otelMetrics.RecordUserOperation(ctx, "custom_operation", true)
    
    // Your business logic
    return nil
}
```

## Custom Tracing

Add custom spans:

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
)

func (s *Service) ProcessData(ctx context.Context, data string) error {
    tracer := otel.Tracer("my-service")
    ctx, span := tracer.Start(ctx, "ProcessData")
    defer span.End()
    
    span.SetAttributes(
        attribute.String("data.size", fmt.Sprintf("%d", len(data))),
    )
    
    // Your processing logic
    
    return nil
}
```

## Sampling Configuration

### Development (100% sampling)
```bash
TRACING_SAMPLE_RATE=1.0
```

### Production (70% sampling)
```bash
APP_ENVIRONMENT=production
TRACING_SAMPLE_RATE=0.7
```

## Troubleshooting

### Check OTel Collector Health
```powershell
curl http://localhost:13133
```

### View Collector Logs
```powershell
docker-compose logs -f otel-collector
```

### Verify Application is Sending Data
```powershell
# Check Prometheus targets
# Open http://localhost:9090/targets
# Check if otel-metrics is UP

# Check Jaeger for traces
# Open http://localhost:16686
```

### Common Issues

1. **No traces appearing in Jaeger**
   - Check `ENABLE_INSTRUMENTATION=true` in .env
   - Verify OTel Collector is running
   - Check application logs for connection errors

2. **No metrics in Prometheus**
   - Verify Prometheus is scraping OTel Collector
   - Check http://localhost:9090/targets
   - Ensure application is sending metrics

3. **Connection refused errors**
   - Ensure all Docker containers are running
   - Check if ports are available (not used by other services)

## Performance Impact

- **Traces**: Minimal overhead (<1% CPU)
- **Metrics**: Negligible overhead
- **Sampling**: Adjustable based on environment
- **Batching**: Enabled to reduce network calls

## Production Recommendations

1. **Use sampling** to reduce data volume (0.1-0.7 recommended)
2. **Set resource limits** for OTel Collector
3. **Configure retention policies** in Prometheus
4. **Enable authentication** for Grafana
5. **Use persistent volumes** for data storage
6. **Monitor collector performance** via its own metrics

## Additional Resources

- [OpenTelemetry Documentation](https://opentelemetry.io/docs/)
- [Jaeger Documentation](https://www.jaegertracing.io/docs/)
- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
