# Docker & Kubernetes Deployment Guide

This guide provides instructions for building and deploying the Blog Go application using Docker and Kubernetes.

## Prerequisites

- Docker (for building images)
- Kubernetes cluster (for deployment)
- kubectl configured to access your cluster
- Container registry access (Docker Hub, GCR, ECR, etc.)

## Building the Docker Image

### 1. Build the image locally

```bash
docker build -t blog-go:latest .
```

### 2. Tag for your registry

```bash
# For Docker Hub
docker tag blog-go:latest your-username/blog-go:latest

# For Google Container Registry
docker tag blog-go:latest gcr.io/your-project/blog-go:latest

# For AWS ECR
docker tag blog-go:latest your-account.dkr.ecr.region.amazonaws.com/blog-go:latest
```

### 3. Push to registry

```bash
# For Docker Hub
docker push your-username/blog-go:latest

# For Google Container Registry
docker push gcr.io/your-project/blog-go:latest

# For AWS ECR
docker push your-account.dkr.ecr.region.amazonaws.com/blog-go:latest
```

## Running with Docker

### Basic run

```bash
docker run -p 8080:8080 blog-go:latest
```

### Run with custom configuration

```bash
docker run -p 8080:8080 \
  -e SERVER_ADDR=:8080 \
  -v $(pwd)/data:/app/data \
  blog-go:latest
```

### Run with Docker Compose

Create a `docker-compose.yml`:

```yaml
version: '3.8'

services:
  blog-go:
    image: blog-go:latest
    ports:
      - "8080:8080"
    environment:
      - SERVER_ADDR=:8080
      - DB_PATH=/app/data/blog.db
    volumes:
      - ./data:/app/data
    restart: unless-stopped
```

Then run:

```bash
docker-compose up -d
```

## Deploying to Kubernetes

### 1. Update the deployment configuration

Edit `k8s/deployment.yaml` and update the image reference:

```yaml
image: your-registry/blog-go:latest
```

### 2. Apply Kubernetes manifests

```bash
# Create the persistent volume claim
kubectl apply -f k8s/pvc.yaml

# Deploy the application
kubectl apply -f k8s/deployment.yaml

# Create the service
kubectl apply -f k8s/service.yaml

# (Optional) Create ingress for external access
kubectl apply -f k8s/ingress.yaml
```

### 3. Verify deployment

```bash
# Check pod status
kubectl get pods -l app=blog-go

# Check service
kubectl get svc blog-go

# View logs
kubectl logs -f deployment/blog-go
```

### 4. Access the application

#### Port forwarding (for testing)

```bash
kubectl port-forward svc/blog-go 8080:80
# Access at http://localhost:8080
```

#### Via Ingress

If you configured ingress, access via your configured domain.

#### Via LoadBalancer

If you changed the service type to LoadBalancer:

```bash
kubectl get svc blog-go
# Use the EXTERNAL-IP shown
```

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_PATH` | `/app/data/blog.db` | Path to SQLite database file |
| `SERVER_ADDR` | `:8080` | Server listen address |

### Storage

The application uses a SQLite database that requires persistent storage. The Kubernetes deployment uses a PersistentVolumeClaim to ensure data persists across pod restarts.

**Important**: SQLite doesn't support concurrent writes well, so keep the deployment at 1 replica.

## Kubernetes Resources

### Adjusting Resource Limits

Edit `k8s/deployment.yaml` to adjust CPU and memory limits:

```yaml
resources:
  requests:
    memory: "64Mi"
    cpu: "100m"
  limits:
    memory: "128Mi"
    cpu: "200m"
```

### Scaling Considerations

⚠️ **Important**: This application uses SQLite, which doesn't support concurrent writes. Keep replicas set to 1.

If you need horizontal scaling, consider:
- Migrating to PostgreSQL or MySQL
- Using a read-replica pattern
- Implementing connection pooling

### Storage Class

If your cluster requires a specific storage class, uncomment and set it in `k8s/pvc.yaml`:

```yaml
storageClassName: your-storage-class
```

## Health Checks

The Dockerfile includes a health check. Kubernetes deployment also configures:

- **Liveness probe**: Checks if the app is running (restarts if failing)
- **Readiness probe**: Checks if the app is ready to serve traffic

## Troubleshooting

### Pod not starting

```bash
# Check pod events
kubectl describe pod -l app=blog-go

# Check logs
kubectl logs -f deployment/blog-go
```

### Database permission issues

Ensure the PVC has correct permissions. The app runs as user 1000:1000.

### Connection issues

```bash
# Test from within the cluster
kubectl run -it --rm debug --image=busybox --restart=Never -- wget -O- http://blog-go/
```

## Security Considerations

1. **Non-root user**: The container runs as a non-root user (UID 1000)
2. **Read-only filesystem**: Consider mounting the root filesystem as read-only
3. **Network policies**: Implement network policies to restrict traffic
4. **TLS**: Use ingress with TLS for production deployments
5. **Secrets**: Store sensitive configuration in Kubernetes Secrets

## Backup and Recovery

### Backup the database

```bash
# Copy database from pod
kubectl cp blog-go-pod-name:/app/data/blog.db ./backup-blog.db
```

### Restore database

```bash
# Copy database to pod
kubectl cp ./backup-blog.db blog-go-pod-name:/app/data/blog.db

# Restart pod to use new database
kubectl rollout restart deployment/blog-go
```

## CI/CD Integration

### GitHub Actions Example

Add to `.github/workflows/deploy.yml`:

```yaml
name: Build and Deploy to Kubernetes

on:
  push:
    branches: [master]

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Build and push Docker image
        run: |
          docker build -t your-registry/blog-go:${{ github.sha }} .
          docker push your-registry/blog-go:${{ github.sha }}

      - name: Deploy to Kubernetes
        run: |
          kubectl set image deployment/blog-go \
            blog-go=your-registry/blog-go:${{ github.sha }}
```

## Monitoring

Consider adding:
- Prometheus metrics endpoint
- Logging aggregation (ELK, Loki)
- Application Performance Monitoring (APM)

## Support

For issues or questions about the application itself, refer to the main README.md.
