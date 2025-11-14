# Docker Deployment

This guide provides instructions for building and running the Blog Go application using Docker.

## Building the Docker Image

```bash
docker build -t blog-go:latest .
```

## Running the Container

### Basic usage

```bash
docker run -p 8080:8080 blog-go:latest
```

The application will be available at `http://localhost:8080`

### With persistent storage

```bash
docker run -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  blog-go:latest
```

This mounts a local `data` directory to persist the SQLite database.

### With custom configuration

```bash
docker run -p 8080:8080 \
  -e SERVER_ADDR=:8080 \
  -e DB_PATH=/app/data/blog.db \
  -v $(pwd)/data:/app/data \
  blog-go:latest
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_PATH` | `/app/data/blog.db` | Path to SQLite database file |
| `SERVER_ADDR` | `:8080` | Server listen address |

## Docker Compose

Create a `docker-compose.yml`:

```yaml
version: '3.8'

services:
  blog-go:
    build: .
    ports:
      - "8080:8080"
    environment:
      - SERVER_ADDR=:8080
      - DB_PATH=/app/data/blog.db
    volumes:
      - ./data:/app/data
    restart: unless-stopped
```

Run with:

```bash
docker-compose up -d
```

## Pushing to a Registry

```bash
# Tag the image
docker tag blog-go:latest your-registry/blog-go:latest

# Push to registry
docker push your-registry/blog-go:latest
```

## Health Check

The container includes a built-in health check that monitors the application on port 8080.

Check container health:

```bash
docker ps
# Look for "healthy" status in the STATUS column
```
