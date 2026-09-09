# syntax=docker/dockerfile:1

# -------------------------------------------------------------------
# Stage 1: Build binary using official Go Alpine
# -------------------------------------------------------------------
FROM golang:1.23-alpine AS builder

WORKDIR /src

# Install build dependencies, CA certificates, and tzdata
RUN apk add --no-cache ca-certificates tzdata git

# Create unprivileged user for security
RUN adduser -D -g '' -u 1001 appuser

# Download modules (cached if go.mod and go.sum don't change)
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked binary with stripped symbols
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -trimpath \
    -o /bin/song-moodboard-api .

# -------------------------------------------------------------------
# Stage 2: Minimal, secure runtime image (~15MB)
# -------------------------------------------------------------------
FROM alpine:3.20

WORKDIR /app

# Copy system certs and timezone
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy compiled binary
COPY --from=builder /bin/song-moodboard-api /app/song-moodboard-api

# Install media extraction and processing tools (FFmpeg, Python3, Node.js, yt-dlp)
RUN apk add --no-cache ca-certificates tzdata ffmpeg python3 py3-pip nodejs curl && \
    pip install --break-system-packages -U yt-dlp && \
    mkdir -p /tmp/downloads && \
    chown -R appuser:appuser /tmp/downloads

# Run as non-root user
USER appuser:appuser

EXPOSE 8080

ENTRYPOINT ["/app/song-moodboard-api"]
