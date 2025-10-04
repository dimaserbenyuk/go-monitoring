# Go Client - Load Tester

Go Client is a load testing tool designed to test the performance of the go-app monitoring service.

## What it does

- **Load Testing**: Creates virtual clients that send HTTP requests to test endpoints
- **Gradual Scaling**: Slowly increases the number of concurrent clients from 0 to maxClients
- **Multiple Endpoints**: Tests different API endpoints randomly:
  - `/health` - Health check endpoint
  - `/api/devices` - Device listing endpoint  
  - `/api/images` - Image processing endpoint
- **Metrics Collection**: Collects performance metrics (response time, status codes)
- **Prometheus Export**: Exposes metrics on port 8082 for monitoring

## How it works

1. **Startup**: Waits 5 seconds for target service to be ready
2. **Client Scaling**: Gradually increases virtual clients every 500ms (default)
3. **Request Sending**: Each client sends requests to random endpoints
4. **Metrics Recording**: Records response times and status codes
5. **Connection Reuse**: Uses HTTP connection pooling for efficiency

## Usage

### Basic usage (test localhost:8000):
```bash
./client
```

### Custom configuration:
```bash
./client -maxClients=20 -baseURL=http://localhost:8000 -scaleInterval=1000
```

### Parameters:
- `-maxClients`: Maximum number of virtual clients (default: 10)
- `-scaleInterval`: Time between scaling up clients in milliseconds (default: 500)
- `-randomSleep`: Random sleep between requests in microseconds (default: 1000)
- `-baseURL`: Base URL of the target server (default: http://localhost:8000)

## Metrics

Client exposes its own metrics on `http://localhost:8082/metrics`:

- `tester_request_duration_seconds`: Request response time with labels:
  - `path`: The endpoint being tested
  - `status`: HTTP response status code

## Example Workflow

1. Start the go-app server:
```bash
cd ../go-app && ./app
```

2. Run load test:
```bash
./client -maxClients=10 -scaleInterval=500
```

3. Monitor metrics:
```bash
curl http://localhost:8082/metrics | grep tester_request
```

## Docker Support

Build Docker image:
```bash
docker build -t go-client .
```

Run with Docker:
```bash
docker run --network=host go-client -baseURL=http://localhost:8000
```

## Integration with Monitoring

The client integrates well with Prometheus and Grafana for visualization:
- Client metrics: `localhost:8082/metrics`
- Server metrics: `localhost:8081/metrics` (from go-app)

This allows you to compare client-side and server-side performance metrics.

curl -s http://localhost:8081/metrics | grep myapp_request_duration_seconds
# HELP myapp_request_duration_seconds Duration of the request.
# TYPE myapp_request_duration_seconds summary
myapp_request_duration_seconds{op="db",quantile="0.9"} 0.003777834
myapp_request_duration_seconds{op="db",quantile="0.99"} 0.012852958
myapp_request_duration_seconds_sum{op="db"} 0.037874084999999995
myapp_request_duration_seconds_count{op="db"} 11
myapp_request_duration_seconds{op="devices",quantile="0.9"} 0.000572959
myapp_request_duration_seconds{op="devices",quantile="0.99"} 0.001122083
myapp_request_duration_seconds_sum{op="devices"} 0.0024880430000000005
myapp_request_duration_seconds_count{op="devices"} 15
myapp_request_duration_seconds{op="health",quantile="0.9"} 0.000122375
myapp_request_duration_seconds{op="health",quantile="0.99"} 0.000777125
myapp_request_duration_seconds_sum{op="health"} 0.001481083
myapp_request_duration_seconds_count{op="health"} 13
myapp_request_duration_seconds{op="s3",quantile="0.9"} 0.015397709
myapp_request_duration_seconds{op="s3",quantile="0.99"} 0.017274833
myapp_request_duration_seconds_sum{op="s3"} 0.091543292
myapp_request_duration_seconds_count{op="s3"} 11

~
❯ curl -s http://localhost:8082/metrics | grep tester_request
# HELP tester_request_duration_seconds Duration of the request.
# TYPE tester_request_duration_seconds summary
tester_request_duration_seconds{path="http://localhost:8000/api/devices",status="200",quantile="0.9"} 0.002500625
tester_request_duration_seconds{path="http://localhost:8000/api/devices",status="200",quantile="0.99"} 0.029524
tester_request_duration_seconds_sum{path="http://localhost:8000/api/devices",status="200"} 0.06785970499999999
tester_request_duration_seconds_count{path="http://localhost:8000/api/devices",status="200"} 35
tester_request_duration_seconds{path="http://localhost:8000/api/images",status="200",quantile="0.9"} 0.016649292
tester_request_duration_seconds{path="http://localhost:8000/api/images",status="200",quantile="0.99"} 0.031311542
tester_request_duration_seconds_sum{path="http://localhost:8000/api/images",status="200"} 0.28472445799999996
tester_request_duration_seconds_count{path="http://localhost:8000/api/images",status="200"} 27
tester_request_duration_seconds{path="http://localhost:8000/health",status="200",quantile="0.9"} 0.004420792
tester_request_duration_seconds{path="http://localhost:8000/health",status="200",quantile="0.99"} 0.015404458
tester_request_duration_seconds_sum{path="http://localhost:8000/health",status="200"} 0.059923917
tester_request_duration_seconds_count{path="http://localhost:8000/health",status="200"} 29

go build -o client && ./client -maxClients=3 -scaleInterval=2000

./client -stats