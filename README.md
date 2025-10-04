# Go Monitoring Service

A Go-based monitoring service with Prometheus metrics support, OpenTelemetry tracing, and AWS S3 integration.

### Prerequisites
- Go 1.22 or higher
- Docker (optional)

### Installation and Setup

1. **Clone the repository**
```bash
git clone <repository-url>
cd go-monitoring
```

2. **Initialize Go module and install dependencies**
```bash
# Navigate to the application directory
cd go-app

# Initialize module (if not already initialized)
go mod init go-monitoring

# Install all dependencies
go mod tidy

# Check dependencies
go list -m all
```

3. **Code verification and formatting**
```bash
# Format code
go fmt ./...

# Check syntax
go vet ./...

# Verify code compiles
go build -o /dev/null .
```

4. **Build application**
```bash
# Regular build for development
go build -o app .

# Production build (static, without CGO)
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Check built file
ls -la app
file app
```

5. **Run application**
```bash
# Local run
./app

# Or directly through go
go run .

# Run with environment variables
APP_PORT=3000 ./app
```

6. **API Testing**
```bash
# Health check
curl http://localhost:8000/health

# Get Prometheus metrics (separate port!)
curl http://localhost:8081/metrics

# List devices (15 IoT devices with UUID, MAC, firmware)
curl http://localhost:8000/api/devices

# List images
curl http://localhost:8000/api/images

# Pretty JSON output
curl http://localhost:8000/api/devices | jq .
curl http://localhost:8000/health | jq .

# Check metrics (Prometheus format - not JSON)
curl http://localhost:8081/metrics
```

7. **Working with Docker**
```bash
# Build Docker image
docker build -t go-monitoring .

# Run container
docker run -p 8000:8000 go-monitoring

# Run with environment variables
docker run -p 8000:8000 -e APP_PORT=8000 -e LOG_LEVEL=debug go-monitoring

# Run in background
docker run -d -p 8000:8000 --name monitoring go-monitoring

# View logs
docker logs monitoring

# Stop container
docker stop monitoring
docker rm monitoring
```

## 📦 Dependencies

The project uses the following main libraries:
- `github.com/gin-gonic/gin` - HTTP web framework
- `github.com/prometheus/client_golang` - Prometheus metrics
- `go.opentelemetry.io/otel` - OpenTelemetry tracing
- `github.com/aws/aws-sdk-go` - AWS SDK
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `gopkg.in/yaml.v2` - YAML configuration

## ⚙️ Configuration

Application settings are configured through `config.yaml` file:

```yaml
appPort: 8000
otlpEndpoint: "http://localhost:4318"
s3:
  region: "us-east-1"
  bucketName: "my-bucket"
postgres:
  host: "localhost"
  port: 5432
  database: "monitoring"
  username: "user"
  password: "password"
```

### Environment Variables

You can override settings through environment variables:
```bash
export APP_PORT=8000
export OTLP_ENDPOINT=http://jaeger:4318
export S3_BUCKET_NAME=production-bucket
```

## 📊 API Endpoints

| Endpoint | Port | Method | Description | Example |
|----------|------|--------|-------------|---------|
| `/health` | 8000 | GET | Application health check | `curl http://localhost:8000/health` |
| `/api/devices` | 8000 | GET | List of 15 IoT devices with UUID, MAC, firmware | `curl http://localhost:8000/api/devices` |
| `/api/images` | 8000 | GET | List of container images | `curl http://localhost:8000/api/images` |
| `/metrics` | 8081 | GET | Prometheus metrics (separate port) | `curl http://localhost:8081/metrics` |

### Response Examples:

**Health Check:**
```json
{
  "status": "up"
}
```

**Devices (sample):**
```json
[
  {
    "UUID": "b0e42fe7-31a5-4894-a441-007e5256afea",
    "mac": "5F-33-CC-1F-43-82", 
    "firmware": "2.1.6"
  }
]
```

## 🔧 Development Commands

### Dependency Management
```bash
# Install/update dependencies
go mod tidy

# View dependencies
go list -m all

# Update dependency
go get -u github.com/gin-gonic/gin

# Add new dependency
go get github.com/new/package

# Remove unused dependencies
go mod tidy
```

### Working with Code
```bash
# Format code
go fmt ./...

# Check syntax and errors
go vet ./...

# Check imports
goimports -w .

# Linter (requires golangci-lint installation)
golangci-lint run
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Detailed coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test
go test -run TestSpecific ./...

# Benchmarks
go test -bench=. ./...
```

mc alias set local http://localhost:9000 admin devops123

mc cp thumbnail.png local/images/