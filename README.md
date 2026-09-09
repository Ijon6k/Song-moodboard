# Song-moodboard

<p align="center">
  <img src="asset/image.png" alt="Song-moodboard" width="100%" />
</p>

Personal project. Built to learn CI/CD pipelines, GitHub Actions, and deploying multi-container apps with Docker on a VPS.

Self-hosted music library with YouTube/stream extraction, object storage, and mood-based organization.

---

## Quick Start (Self-Hosted)

Run the entire app with a single `docker-compose.yml`. No source code, no build steps required.

### 1. Download Compose & Example Config

```bash
mkdir song-moodboard && cd song-moodboard

# Download the standalone compose file
curl -O https://raw.githubusercontent.com/Ijon6k/Song-moodboard/main/docker-compose.yml

# (Optional) Download .env.example if you want to customize credentials
curl -O https://raw.githubusercontent.com/Ijon6k/Song-moodboard/main/.env.example
cp .env.example .env
```

> **Note:** Creating `.env` is **optional**. If omitted, safe default local credentials and port `1122` will be used automatically.

### 2. Start

```bash
docker compose pull
docker compose up -d
```

Open **`http://localhost:1122`** (or `http://your-server-ip:1122`).

---

## Maintenance & Commands

### Updating to Latest Release

```bash
docker compose pull
docker compose up -d --remove-orphans
```

### Useful Commands

```bash
# View logs
docker compose logs -f

# Check container status
docker compose ps

# Stop containers
docker compose down

# Stop and wipe database/storage volumes (destructive!)
docker compose down -v
```

### Rollback to a Specific Version

```bash
# By semver tag
IMAGE_BACKEND=ghcr.io/ijon6k/song-moodboard-backend:v0.1.0 \
IMAGE_FRONTEND=ghcr.io/ijon6k/song-moodboard-frontend:v0.1.0 \
  docker compose up -d

# By commit SHA
IMAGE_BACKEND=ghcr.io/ijon6k/song-moodboard-backend:sha-a3f9c21 \
IMAGE_FRONTEND=ghcr.io/ijon6k/song-moodboard-frontend:sha-a3f9c21 \
  docker compose up -d
```

---

## Development (Build from Source)

If you are cloning this repository to modify the code and build locally:

```bash
git clone https://github.com/Ijon6k/Song-moodboard.git
cd Song-moodboard

# Build and run with dev compose
docker compose -f docker-compose.dev.yml up -d --build
```

---

## Architecture & Stack

- **Frontend**: SvelteKit (Bun)
- **Backend**: Go 1.23, Chi router
- **Audio**: yt-dlp + FFmpeg (extracts to 320k MP3)
- **Database**: PostgreSQL 16
- **Cache**: Redis 7
- **Storage**: SeaweedFS (S3-compatible)
- **Proxy**: Nginx on port 1122
- **CI/CD**: GitHub Actions → GHCR

```
Browser → Nginx (:1122)
              /api/* → Go backend (:8080) → Postgres, Redis, SeaweedFS
              /*     → SvelteKit SSR (:3000)
```

Only port `1122` is exposed to the host. Internal services (Postgres, Redis, SeaweedFS, Go API) communicate over an isolated Docker bridge network.

---

## Image Tags

Images are built and pushed to GHCR automatically via GitHub Actions:

| Trigger | Tags |
|---|---|
| Push to `dev` | `:dev`, `:sha-<commit>` |
| Push to `main` | `:sha-<commit>` |
| Git tag `v*.*.*` | `:latest`, `:<version>`, `:<major>.<minor>`, `:sha-<commit>` |
