# Build stage
FROM golang:1.21-alpine AS builder

# Build arguments
ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_DATE=unknown

# Install dependencies
RUN apk add --no-cache git make ca-certificates

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-w -s -X main.version=${VERSION} -X main.gitCommit=${GIT_COMMIT} -X main.buildDate=${BUILD_DATE}" \
    -o nrp-aws-kip \
    ./cmd/nrp-aws-kip

# Runtime stage
FROM alpine:3.19

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 nrpkip && \
    adduser -D -u 1000 -G nrpkip nrpkip

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/nrp-aws-kip /app/nrp-aws-kip

# Copy example config (optional)
COPY config.yaml.example /app/config.yaml.example

# Change ownership
RUN chown -R nrpkip:nrpkip /app

# Switch to non-root user
USER nrpkip

# Expose metrics port
EXPOSE 10255

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:10255/healthz || exit 1

# Run the binary
ENTRYPOINT ["/app/nrp-aws-kip"]
