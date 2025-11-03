# Docker-Compose Troubleshooting Guide

## Common Error: Docker Desktop Not Running

### Error Message
```
error during connect: Get "http://%2F%2F.%2Fpipe%2FdockerDesktopLinuxEngine/v1.51/...": 
open //./pipe/dockerDesktopLinuxEngine: The system cannot find the file specified.
```

### Solution

#### 1. Start Docker Desktop
- Open **Docker Desktop** application on Windows
- Wait for Docker to fully initialize (icon in system tray turns green)
- The whale icon should be steady, not animated

#### 2. Verify Docker is Running
```powershell
docker ps
```

If successful, you should see a list of running containers (or empty list if none running).

#### 3. Try Again
```powershell
docker-compose up -d
```

---

## Other Common Issues

### Issue: Port Already in Use

**Error:**
```
Bind for 0.0.0.0:5432 failed: port is already allocated
```

**Solution:**
```powershell
# Check what's using the port
netstat -ano | findstr :5432

# Option 1: Stop the process using the port
# Option 2: Change the port in docker-compose.yml
```

### Issue: Network Conflicts

**Error:**
```
network ... has active endpoints
```

**Solution:**
```powershell
# Remove all stopped containers
docker-compose down

# Remove networks
docker network prune

# Try again
docker-compose up -d
```

### Issue: Volume Permissions

**Error:**
```
permission denied while trying to connect to the Docker daemon socket
```

**Solution:**
- Make sure Docker Desktop is running
- Run PowerShell/Command Prompt as Administrator
- Restart Docker Desktop

### Issue: Images Not Pulling

**Error:**
```
unable to get image '...': error during connect
```

**Solution:**
```powershell
# Check internet connection
# Check Docker Hub is accessible
ping hub.docker.com

# Try manual pull
docker pull prom/prometheus:latest
docker pull grafana/grafana:latest
docker pull jaegertracing/all-in-one:latest
docker pull otel/opentelemetry-collector-contrib:latest
docker pull postgres:15
```

---

## Quick Fix Commands

### Completely Reset Everything
```powershell
# Stop all containers
docker-compose down

# Remove volumes (WARNING: This deletes data!)
docker-compose down -v

# Remove all unused containers, networks, images
docker system prune -a

# Start fresh
docker-compose up -d
```

### View Logs
```powershell
# All services
docker-compose logs

# Specific service
docker-compose logs otel-collector
docker-compose logs jaeger
docker-compose logs prometheus
docker-compose logs grafana
docker-compose logs postgres

# Follow logs (real-time)
docker-compose logs -f otel-collector
```

### Restart a Single Service
```powershell
docker-compose restart otel-collector
docker-compose restart prometheus
docker-compose restart grafana
```

### Check Service Status
```powershell
# List all containers
docker-compose ps

# Check specific service health
docker inspect --format='{{.State.Health.Status}}' postgres
```

---

## Verification Steps

After starting services, verify everything is working:

### 1. Check All Containers Are Running
```powershell
docker-compose ps
```

All services should show "Up" status.

### 2. Check Health Endpoints
```powershell
# OTel Collector Health
curl http://localhost:13133

# Prometheus
curl http://localhost:9090

# Grafana
curl http://localhost:3000

# Jaeger
curl http://localhost:16686
```

### 3. Check Logs for Errors
```powershell
docker-compose logs | Select-String "error"
docker-compose logs | Select-String "ERROR"
```

---

## Prevention Tips

1. **Always start Docker Desktop before running docker-compose**
2. **Wait for Docker to fully initialize** (green icon in system tray)
3. **Check available disk space** - Docker needs space for images and volumes
4. **Keep Docker Desktop updated** to the latest version
5. **Configure Docker Desktop resources** (CPU, Memory) appropriately

---

## Getting Help

If issues persist:

1. Check Docker Desktop logs:
   - Click Docker icon in system tray
   - Click "Troubleshoot"
   - View logs

2. Verify Docker version:
   ```powershell
   docker --version
   docker-compose --version
   ```

3. Check system requirements:
   - Windows 10/11 Pro, Enterprise, or Education
   - WSL 2 enabled (for Docker Desktop)
   - Virtualization enabled in BIOS

4. Restart Docker Desktop:
   - Right-click Docker icon
   - Select "Restart Docker Desktop"

---

## Emergency: Complete Docker Reset

If nothing works:

```powershell
# Stop Docker Desktop

# Remove Docker data (Windows)
# This will delete all images, containers, and volumes!
# Settings > Resources > Advanced > "Clean / Purge data"

# Restart Docker Desktop

# Pull images again
docker-compose pull

# Start services
docker-compose up -d
```
