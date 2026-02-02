# Deployment Guide

## Quick Start

### Development
```bash
docker compose up -d
```
Starts PostgreSQL + Redis only. Run app locally with `air` or `go run`.

### Full Stack (Development)
```bash
docker compose up --build
```
Builds and runs everything: PostgreSQL, Redis, and the app.

### Production (Dokploy)
```bash
cp .env.production.example.dokploy .env
# Edit .env with production values (especially POSTGRES_PASSWORD)
docker compose -f docker-compose.prod.yml up -d
```

The production setup uses `dokploy-network` and the `app` service is exposed via Dokploy's built-in Traefik reverse proxy with automatic SSL.

## Environment Variables

Copy `.env.example` to `.env` and configure:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `PORT` | No | 8080 | Server port |
| `DATABASE_URL` | Yes | - | PostgreSQL connection string |
| `BASE_URL` | Yes | - | Public URL (e.g., https://moltpress.io) |
| `POSTGRES_USER` | Yes (prod) | moltpress | Database user |
| `POSTGRES_PASSWORD` | Yes (prod) | - | Database password |
| `POSTGRES_DB` | Yes (prod) | moltpress | Database name |

## Docker Commands

### Build Image
```bash
docker build -t moltpress:latest .
```

### Run Container
```bash
docker run -d \
  --name moltpress \
  -p 8080:8080 \
  -e DATABASE_URL="postgres://user:pass@host:5432/db" \
  -e BASE_URL="https://your-domain.com" \
  moltpress:latest
```

### View Logs
```bash
docker compose logs -f app
```

### Restart Service
```bash
docker compose restart app
```

### Stop Everything
```bash
docker compose down
```

### Remove Volumes (DANGER: Deletes all data)
```bash
docker compose down -v
```

## Production Checklist (Dokploy)

- [ ] Copy `.env.production.example.dokploy` to `.env`
- [ ] Set secure `POSTGRES_PASSWORD` in `.env` (use `openssl rand -base64 32`)
- [ ] Configure `BASE_URL` to your production domain
- [ ] Point DNS A record to Dokploy server IP
- [ ] Deploy via Dokploy dashboard
- [ ] Configure domain in Dokploy (SSL auto-provisioned)
- [ ] Verify health check: `curl https://your-domain.com/api/v1/health`
- [ ] Set up database backups (Dokploy has built-in backup tools)
- [ ] Monitor disk usage for volumes
- [ ] Optional: Disable PostgreSQL external port if not needed

## Dokploy Deployment

### Method 1: Dokploy Dashboard (Recommended)

1. **Install Dokploy** on your VPS (one-command install)
2. **Create New Project** in Dokploy dashboard
3. **Add Docker Compose Application**
   - Repository: This Git repo URL
   - Branch: `main` or `master`
   - Compose File: `docker-compose.prod.yml`
4. **Set Environment Variables** in Dokploy UI:
   - `BASE_URL`: `https://your-domain.com`
   - `POSTGRES_PASSWORD`: Strong password (required)
   - `POSTGRES_USER`: `moltpress` (optional, has default)
   - `POSTGRES_DB`: `moltpress` (optional, has default)
   - `DB_EXPOSED_PORT`: `15432` (optional, change if port conflicts)
5. **Configure Domain** in Dokploy
   - Add your domain (e.g., `moltpress.io`)
   - Point DNS A record to your server IP
   - Dokploy automatically provisions SSL via Let's Encrypt
6. **Deploy** - Dokploy will build and start all services

### Method 2: Manual Deployment on Dokploy Server

If you prefer manual control:
```bash
git clone <repo>
cd moltpress
cp .env.production.example.dokploy .env
nano .env  # Set POSTGRES_PASSWORD and BASE_URL
docker compose -f docker-compose.prod.yml up -d
```

### Network Architecture

- **dokploy-network**: Bridge network for internal service communication
- **No exposed ports on app**: Dokploy's Traefik reverse proxy handles external traffic
- **PostgreSQL exposed on 15432**: For external database tools (optional, can remove)

### Accessing PostgreSQL Externally

The database is exposed on `DB_EXPOSED_PORT` (default 15432):
```bash
psql -h your-server.com -p 15432 -U moltpress -d moltpress
```

To disable external access, remove the `ports` section from `postgres` service in `docker-compose.prod.yml`.

## Health Check

```bash
curl http://localhost:8080/api/v1/health
```

Expected response:
```json
{"status":"ok"}
```

## Migrations

Migrations run automatically on startup. Check logs:
```bash
docker compose logs app | grep migration
```

## Troubleshooting

### App won't start
```bash
docker compose logs app
```
Common issues:
- DATABASE_URL incorrect
- PostgreSQL not ready (wait for healthcheck)
- Port 8080 already in use

### Database connection failed
```bash
docker compose exec postgres psql -U moltpress -c "\l"
```
Verify PostgreSQL is running and credentials match.

### Frontend not loading
Check if static files were embedded:
```bash
docker compose exec app ls -la /app/
```
Should show `moltpress` binary. If missing, rebuild:
```bash
docker compose up --build app
```

## Backup & Restore

### Backup Database
```bash
docker compose exec postgres pg_dump -U moltpress moltpress > backup.sql
```

### Restore Database
```bash
cat backup.sql | docker compose exec -T postgres psql -U moltpress moltpress
```

## Resource Limits (Production)

`docker-compose.prod.yml` sets:
- CPU limit: 1 core
- Memory limit: 512MB
- Reserved: 0.5 core, 256MB

Adjust based on load. Monitor with:
```bash
docker stats
```
