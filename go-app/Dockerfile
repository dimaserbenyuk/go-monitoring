# Multi-platform Go build
FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS build

# Install build dependencies
RUN apk --no-cache add ca-certificates git tzdata

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build arguments for multi-platform
ARG TARGETOS
ARG TARGETARCH

# Build the application
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -a -installsuffix cgo \
    -ldflags='-w -s -extldflags "-static"' \
    -o /app/go-monitoring .

# Final stage - minimal image
FROM scratch

# Copy timezone data
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo

# Copy ca-certificates for HTTPS
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary
COPY --from=build /app/go-monitoring /go-monitoring

# Expose ports
EXPOSE 8000 8081

# Add labels for better container management
LABEL org.opencontainers.image.title="Go Monitoring App"
LABEL org.opencontainers.image.description="Monitoring service with MinIO and PostgreSQL integration"
LABEL org.opencontainers.image.version="1.0.0"
LABEL org.opencontainers.image.source="https://github.com/dimaserbenyuk/go-monitoring"

ENTRYPOINT ["/go-monitoring"]