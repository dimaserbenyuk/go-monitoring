FROM golang:1.22.3-alpine AS build

WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

# Build with CGO disabled for static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /myapp .

FROM scratch

# Copy ca-certificates from build stage
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /myapp /myapp

ENTRYPOINT ["/myapp"]