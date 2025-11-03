# OpenTelemetry Monitoring Architecture

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                        Go Application                                │
│                     (ZOGTest-Golang API)                             │
│                                                                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │
│  │   HTTP       │  │   Business   │  │   Database   │              │
│  │   Handlers   │  │   Logic      │  │   Layer      │              │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘              │
│         │                  │                  │                       │
│         └──────────────────┼──────────────────┘                      │
│                            │                                          │
│  ┌─────────────────────────▼──────────────────────────┐             │
│  │      OpenTelemetry Instrumentation                  │             │
│  │  - Traces (spans, context propagation)             │             │
│  │  - Metrics (counters, histograms, gauges)          │             │
│  │  - Resource attributes                             │             │
│  └─────────────────────────┬──────────────────────────┘             │
│                            │                                          │
└────────────────────────────┼──────────────────────────────────────────┘
                             │ OTLP (gRPC/HTTP)
                             │ Port 4317/4318
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│              OpenTelemetry Collector                                 │
│                                                                       │
│  ┌──────────┐      ┌──────────────┐      ┌──────────────┐          │
│  │ Receiver │ ───> │  Processors  │ ───> │  Exporters   │          │
│  │  OTLP    │      │  - Batch     │      │  - Jaeger    │          │
│  │ gRPC/HTTP│      │  - Memory    │      │  - Prometheus│          │
│  │          │      │    Limiter   │      │  - Logging   │          │
│  └──────────┘      │  - Resource  │      └──────────────┘          │
│                    └──────────────┘                                  │
│                                                                       │
│  Health Check: :13133                                                │
│  Metrics: :8888, :8889                                               │
└───────────────────────┬──────────────────┬──────────────────────────┘
                        │                  │
                        │ Traces           │ Metrics
                        ▼                  ▼
        ┌───────────────────────┐  ┌──────────────────────┐
        │      Jaeger           │  │     Prometheus        │
        │   (All-in-One)        │  │   (Time Series DB)    │
        │                       │  │                       │
        │  - Collector          │  │  - Scraper            │
        │  - Query Service      │  │  - Storage            │
        │  - UI                 │  │  - Query Engine       │
        │                       │  │  - Alerting           │
        │  Port: 16686 (UI)     │  │  Port: 9090 (UI)      │
        │        14250 (gRPC)   │  │                       │
        └───────────┬───────────┘  └──────────┬────────────┘
                    │                         │
                    │                         │
                    └────────────┬────────────┘
                                 │
                                 ▼
                    ┌─────────────────────────┐
                    │       Grafana           │
                    │   (Visualization)       │
                    │                         │
                    │  - Dashboards           │
                    │  - Alerts               │
                    │  - Data Source Queries  │
                    │                         │
                    │  Port: 3000 (UI)        │
                    │  Login: admin/admin     │
                    └─────────────────────────┘
```

## Data Flow

### Traces Flow
```
HTTP Request
    ↓
Application Handler (creates root span)
    ↓
Business Logic (creates child spans)
    ↓
Database Query (creates db span)
    ↓
OpenTelemetry SDK (batches spans)
    ↓
OTLP Exporter (sends to collector)
    ↓
OTel Collector (processes & exports)
    ↓
Jaeger (stores & indexes traces)
    ↓
Jaeger UI (visualizes traces)
```

### Metrics Flow
```
Application Events
    ↓
OTel Metrics Instruments (counter, histogram, gauge)
    ↓
Meter Provider (aggregates metrics)
    ↓
OTLP Exporter (sends to collector)
    ↓
OTel Collector (processes & exports)
    ↓
Prometheus Exporter (:8889)
    ↓
Prometheus (scrapes & stores)
    ↓
Grafana (queries & visualizes)
```

## Network Ports

| Service | Port | Protocol | Purpose |
|---------|------|----------|---------|
| Application | 8000 | HTTP | API Server |
| Application | 8000/metrics | HTTP | Prometheus Metrics |
| OTel Collector | 4317 | gRPC | OTLP Receiver |
| OTel Collector | 4318 | HTTP | OTLP Receiver |
| OTel Collector | 8888 | HTTP | Collector Metrics |
| OTel Collector | 8889 | HTTP | Prometheus Exporter |
| OTel Collector | 13133 | HTTP | Health Check |
| Jaeger | 16686 | HTTP | Jaeger UI |
| Jaeger | 14250 | gRPC | Jaeger Collector |
| Jaeger | 14268 | HTTP | Jaeger Collector |
| Prometheus | 9090 | HTTP | Prometheus UI |
| Grafana | 3000 | HTTP | Grafana UI |
| PostgreSQL | 5432 | TCP | Database |

## Instrumentation Points

### Automatic Instrumentation
```
┌─────────────────────────────────────┐
│  HTTP Middleware                     │
│  - Request/Response timing           │
│  - Status codes                      │
│  - Route patterns                    │
│  - In-flight requests                │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│  Database Layer (otelpgx)            │
│  - Query execution time              │
│  - Connection pool metrics           │
│  - Query statements                  │
└─────────────────────────────────────┘
```

### Custom Instrumentation
```
┌─────────────────────────────────────┐
│  Business Metrics                    │
│  - User operations counter           │
│  - Topic operations counter          │
│  - News operations counter           │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│  Custom Spans                        │
│  - Service layer operations          │
│  - Repository operations             │
│  - Business logic steps              │
└─────────────────────────────────────┘
```

## Metric Types

### Counters (monotonically increasing)
- `http.server.requests.total` - Total HTTP requests
- `db.queries.total` - Total database queries
- `business.user.operations.total` - User operations
- `business.topic.operations.total` - Topic operations
- `business.news.operations.total` - News operations

### Histograms (distributions)
- `http.server.request.duration` - HTTP request latency
- `db.query.duration` - Database query latency

### Gauges (current values)
- `http.server.requests.in_flight` - Active HTTP requests
- `db.connections.open` - Open database connections

## Trace Structure

```
Trace: HTTP GET /api/v1/users
├─ Span: HTTP GET /api/v1/users [ROOT]
│  ├─ Attributes: http.method, http.route, http.status_code
│  ├─ Duration: 45ms
│  │
│  ├─ Span: UserService.GetUsers
│  │  ├─ Attributes: service.operation
│  │  ├─ Duration: 42ms
│  │  │
│  │  └─ Span: UserRepository.FindAll
│  │     ├─ Attributes: db.system, db.operation
│  │     ├─ Duration: 38ms
│  │     │
│  │     └─ Span: pgx.Query
│  │        ├─ Attributes: db.statement, db.table
│  │        └─ Duration: 35ms
│  │
│  └─ Status: OK
```

## Dashboard Structure

```
┌──────────────────────────────────────────┐
│  API Overview Dashboard                  │
├──────────────────────────────────────────┤
│                                           │
│  ┌─────────────┐  ┌─────────────┐       │
│  │ Request Rate│  │ Latency p95 │       │
│  └─────────────┘  └─────────────┘       │
│                                           │
│  ┌─────────────┐  ┌─────────────┐       │
│  │ Error Rate  │  │ In-Flight   │       │
│  └─────────────┘  └─────────────┘       │
│                                           │
│  ┌─────────────────────────────┐        │
│  │   Request Duration Heatmap   │        │
│  └─────────────────────────────┘        │
│                                           │
│  ┌─────────────────────────────┐        │
│  │   Top Endpoints by Traffic   │        │
│  └─────────────────────────────┘        │
│                                           │
└──────────────────────────────────────────┘
```

## Sampling Strategy

```
┌─────────────────────────────────────────┐
│  Environment: Development                │
│  Sampling: AlwaysSample (100%)          │
│  - Trace all requests                    │
│  - Good for debugging                    │
└─────────────────────────────────────────┘

┌─────────────────────────────────────────┐
│  Environment: Production                 │
│  Sampling: ParentBased(TraceIDRatio)    │
│  - Sample 10-70% of traces              │
│  - Respect parent sampling decision      │
│  - Reduce storage & bandwidth            │
└─────────────────────────────────────────┘
```

## Security Considerations

```
┌─────────────────────────────────────────┐
│  Development                             │
│  - No authentication                     │
│  - HTTP (insecure)                       │
│  - All ports exposed                     │
└─────────────────────────────────────────┘

┌─────────────────────────────────────────┐
│  Production (Recommended)                │
│  - Grafana authentication                │
│  - HTTPS/TLS for collectors              │
│  - Network isolation                     │
│  - API tokens                            │
│  - Firewall rules                        │
└─────────────────────────────────────────┘
```
