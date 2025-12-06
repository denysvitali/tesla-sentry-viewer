# Tesla Sentry Viewer

A comprehensive tool for processing and serving Tesla Sentry Mode videos with rich terminal UI.

## Features

- **Video Processing**: Merge multiple Tesla Sentry videos into single files
- **HTTP API**: Serve Sentry clips and metadata over a RESTful API
- **Advanced CLI**: Modern CLI with Cobra/Viper framework and Lipgloss terminal UI
- **Configuration**: Flexible configuration via YAML, environment variables, and CLI flags
- **Health Monitoring**: Built-in health check endpoints
- **Docker Support**: Containerized deployment with health checks
- **Rich Terminal Output**: Beautiful terminal tables, spinners, and styled output

## Installation

### From Source

```bash
git clone https://github.com/denysvitali/tesla-sentry-viewer.git
cd tesla-sentry-viewer
make build-cli
```

### Using Docker

```bash
docker build -t tesla-sentry-viewer .
docker run -p 8150:8150 -v /path/to/sentry/clips:/data tesla-sentry-viewer
```

## Usage

### New CLI Structure

The application now uses a modern CLI structure with Cobra and Viper:

```bash
# Show help
tesla-sentry-viewer --help

# Process videos from a Sentry event directory
tesla-sentry-viewer process /path/to/sentry/event/directory

# Start the API server
tesla-sentry-viewer server /path/to/sentry/clips/directory

# Show version information
tesla-sentry-viewer version
```

### CLI Mode (Video Processing)

```bash
# Process videos with debug logging
tesla-sentry-viewer process -D /path/to/sentry/event/directory

# Custom output directory
tesla-sentry-viewer process -o /custom/output /path/to/sentry/event/directory

# Custom ffmpeg path
tesla-sentry-viewer process --ffmpeg /usr/local/bin/ffmpeg /path/to/sentry/event/directory

# Using configuration file
tesla-sentry-viewer process --config /path/to/config.yaml /path/to/sentry/event/directory
```

### Server Mode (API)

```bash
# Start the API server
tesla-sentry-viewer server /path/to/sentry/clips/directory

# Custom listen address
tesla-sentry-viewer server -l 0.0.0.0:8080 /path/to/sentry/clips/directory

# Show version
tesla-sentry-viewer version
```

## Configuration

The application now supports multiple configuration methods in priority order:
1. CLI flags
2. Environment variables
3. Configuration file
4. Default values

### Configuration File

Create a YAML configuration file (e.g., `~/.tesla-sentry-viewer.yaml`):

```yaml
# Debug settings
debug: false

# File paths
output: "output"
ffmpeg: "ffmpeg"

# Server settings
server:
  listen: "127.0.0.1:8150"
  max_workers: 4

# Logging settings
log_level: "info"
```

### Environment Variables

Prefix all environment variables with `TESLA_SENTRY_`:

| Variable | CLI Flag | Default | Description |
|----------|----------|---------|-------------|
| `TESLA_SENTRY_DEBUG` | `-D, --debug` | `false` | Enable debug logging |
| `TESLA_SENTRY_OUTPUT` | `-o, --output` | `output` | Output directory for processed videos |
| `TESLA_SENTRY_FFMPEG` | `--ffmpeg` | `ffmpeg` | Path to ffmpeg binary |
| `TESLA_SENTRY_SERVER_LISTEN` | `-l, --listen` | `127.0.0.1:8150` | Server listen address |
| `TESLA_SENTRY_SERVER_MAX_WORKERS` | `--max-workers` | `4` | Maximum concurrent workers |
| `TESLA_SENTRY_LOG_LEVEL` | - | `info` | Log level (debug, info, warn, error) |

### CLI Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--debug` | `-D` | `false` | Enable debug logging |
| `--output` | `-o` | `output` | Output directory for processed videos |
| `--ffmpeg` | - | `ffmpeg` | Path to ffmpeg binary |
| `--listen` | `-l` | `127.0.0.1:8150` | Server listen address |
| `--max-workers` | - | `4` | Maximum concurrent workers |
| `--config` | - | - | Path to configuration file |

## API Endpoints

### Health Check
```
GET /health
```
Returns server health status and version information.

### Version
```
GET /version
```
Returns application version information.

### List Clips
```
GET /api/v1/clips
```
Returns a list of available Sentry clips.

### Get Clip Details
```
GET /api/v1/clips/:clip_id
```
Returns detailed information about a specific clip.

### Get Clip Thumbnail
```
GET /api/v1/clips/:clip_id/thumb
```
Returns the thumbnail image for a clip.

### Get Clip File
```
GET /api/v1/clips/:clip_id/:file_name
```
Returns a specific video file from a clip.

## Docker Usage

### Build
```bash
make docker-build
```

### Run
```bash
make docker-run
```

### Development
```bash
make docker-run-dev
```

## Releases

The project uses GoReleaser for automated releases. You can build releases using the following commands:

### Build Release

```bash
# Build a release (creates binaries for multiple platforms)
make release

# Build a snapshot release (for testing)
make release-snapshot

# Check GoReleaser configuration
make check
```

### Release Process

1. Update the version in `go.mod` if needed
2. Create a Git tag: `git tag v1.0.0`
3. Push the tag: `git push origin v1.0.0`
4. GoReleaser will automatically build and publish the release on GitHub

### GitHub Actions Integration

The project includes a GitHub Actions workflow (`.github/workflows/release.yml`) that automatically triggers GoReleaser when you push a tag starting with `v` (e.g., `v1.0.0`).

The workflow:
- Runs on Ubuntu latest
- Sets up Go 1.23.4
- Installs GoReleaser
- Builds binaries for all platforms
- Creates GitHub release with assets
- Uploads artifacts for debugging

### Available Binaries

GoReleaser builds binaries for the following platforms:
- Linux (amd64, arm64, armv7)
- Windows (amd64, arm64)
- macOS (amd64, arm64)

## Development

### Build
```bash
# Build CLI application
make build-cli

# Build server application
make build-server
```

### Test
```bash
make test
```

### Format
```bash
make format
```

### Lint
```bash
make lint
```

### Configuration Examples

See `config.example.yaml` for a complete configuration example.

## Error Codes

The API uses standardized error codes:

- `INVALID_INPUT`: Invalid input parameters
- `FILE_NOT_FOUND`: Requested file not found
- `DIRECTORY_NOT_FOUND`: Requested directory not found
- `PROCESSING_FAILED`: Video processing failed
- `INVALID_CLIP_ID`: Invalid clip ID format

## Performance Optimization

The application includes several performance optimizations:

- Concurrent video processing with configurable worker count
- Efficient file handling and memory management
- Optimized Docker image with minimal dependencies
- Health checks for monitoring

## Security

- Input validation and sanitization
- Proper error handling and logging
- Non-root user in Docker containers
- CORS configuration for API security
- Comprehensive security headers

## License

MIT