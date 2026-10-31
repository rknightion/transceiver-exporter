# ==============================================================================
# Dockerfile for transceiver-exporter
# Multi-stage build for minimal, secure container image
# ==============================================================================

# ------------------------------------------------------------------------------
# Stage 1: Build stage
# ------------------------------------------------------------------------------
FROM --platform=$BUILDPLATFORM golang:1.27.3-alpine@sha256:d200d2dd76b03a99b52651fd1325aff7dcd35422c35c01bac56571b0507f3deb AS builder

RUN apk upgrade --no-cache ca-certificates ca-certificates-bundle

# Set working directory
WORKDIR /src

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies (cached unless go.mod/go.sum change)
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build arguments for version injection
ARG VERSION=dev
ARG COMMIT=unknown

# Target platform arguments (automatically set by buildx)
ARG TARGETOS
ARG TARGETARCH

# Build the binary
# - CGO_ENABLED=0: Pure Go, static binary
# - -ldflags: Strip debug info and inject version
# - -trimpath: Remove file paths for reproducibility
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
    -ldflags="-w -s -X 'main.version=${VERSION}'" \
    -trimpath \
    -o /transceiver-exporter \
    .

# ------------------------------------------------------------------------------
# Stage 2: Production image
# Uses distroless for minimal attack surface (~2MB base)
# Includes: CA certs, tzdata, /tmp
# Note: Runs as root (required for ethtool module EEPROM access)
# ------------------------------------------------------------------------------
FROM gcr.io/distroless/static-debian13:latest@sha256:448810357cc8e67b7208988abdba0f5ec3d3aebee0a5b13eb624979792cffc0f

# OCI image labels
LABEL org.opencontainers.image.title="transceiver-exporter"
LABEL org.opencontainers.image.description="Prometheus exporter for pluggable transceivers on Linux"
LABEL org.opencontainers.image.source="https://github.com/rknightion/transceiver-exporter"
LABEL org.opencontainers.image.licenses="AGPL-3.0"

# Copy the binary from builder
COPY --from=builder /transceiver-exporter /transceiver-exporter

# Expose the metrics port
EXPOSE 9458

# Entrypoint
ENTRYPOINT ["/transceiver-exporter"]

# Default arguments
CMD ["--web.listen-address=[::]:9458"]
