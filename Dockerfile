FROM golang:1.23.4-alpine AS builder

# Install build dependencies
RUN apk add --no-cache make bash git

# Set build arguments
ARG VERSION=dev
ARG BUILD_DATE
ARG GIT_COMMIT

# Copy source code
COPY . /app
WORKDIR /app

# Build the application
RUN make build

# Create a non-root user
RUN adduser -D appuser

FROM alpine:3.18

# Install runtime dependencies
RUN apk add --no-cache ffmpeg

# Create app directory and user
RUN mkdir -p /app && chown appuser:appuser /app
WORKDIR /app

# Copy built binary
COPY --from=builder --chown=appuser:appuser /app/build/tesla-sentry-viewer .

# Copy configuration files if needed
COPY --from=builder --chown=appuser:appuser /app/config.yaml ./

# Set permissions
RUN chmod +x /app/tesla-sentry-viewer

# Switch to non-root user
USER appuser

# Set environment variables
ENV OUTPUT_DIR=/app/output
ENV FFMPEG_PATH=/usr/bin/ffmpeg
ENV MAX_WORKERS=4

# Health check
HEALTHCHECK --interval=30s --timeout=3s \
  CMD curl -f http://localhost:8150/health || exit 1

# Expose port
EXPOSE 8150

# Entrypoint
ENTRYPOINT ["/app/tesla-sentry-viewer"]
CMD ["--help"]
