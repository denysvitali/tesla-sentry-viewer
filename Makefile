IMAGE=dvitali/tesla-sentry-viewer
VERSION=$(shell ./get-version.sh)
TAG=$(VERSION)
BUILD_DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")

# Build targets
build:
	@echo "Building tesla-sentry-viewer version $(VERSION)"
	mkdir -p build
	CGO_ENABLED=0 \
	go build \
		-ldflags "-X main.Version=$(VERSION) -X 'main.BuildDate=$(BUILD_DATE)' -X 'main.GitCommit=$(GIT_COMMIT)'" \
		-o ./build/tesla-sentry-viewer ./cmd/server

build-server:
	@echo "Building tesla-sentry-viewer server version $(VERSION)"
	mkdir -p build
	CGO_ENABLED=0 \
	go build \
		-ldflags "-X main.Version=$(VERSION) -X 'main.BuildDate=$(BUILD_DATE)' -X 'main.GitCommit=$(GIT_COMMIT)'" \
		-o ./build/tesla-sentry-viewer-server ./cmd/server

# Docker targets
docker-build:
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		-t "$(IMAGE):$(TAG)" \
		.

docker-build-latest:
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		-t "$(IMAGE):latest" \
		.

docker-push:
	docker push "$(IMAGE):$(TAG)"

docker-push-latest:
	docker push "$(IMAGE):latest"

# Test targets
test:
	@echo "Running tests..."
	go test ./... -v

test-short:
	@echo "Running short tests..."
	go test ./... -short -v

test-coverage:
	@echo "Running tests with coverage..."
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

# Clean targets
clean:
	@echo "Cleaning build artifacts..."
	rm -rf build/
	rm -f coverage.out

# Run targets
run:
	@echo "Running tesla-sentry-viewer..."
	./build/tesla-sentry-viewer --help

docker-run:
	docker run \
		--rm \
		--name "tesla-sentry-viewer" \
		-v "/run/media/$$USER/TESLADRIVE/:/mnt:ro" \
		-p 8150:8150 \
		-e OUTPUT_DIR=/app/output \
		-e FFMPEG_PATH=/usr/bin/ffmpeg \
		-e MAX_WORKERS=4 \
		"$(IMAGE):$(TAG)" \
		-l "0.0.0.0:8150" \
		"/mnt/TeslaCam/SentryClips"

docker-run-dev:
	docker run \
		--rm \
		--name "tesla-sentry-viewer-dev" \
		-v "$(PWD):/app" \
		-v "/run/media/$$USER/TESLADRIVE/:/mnt:ro" \
		-p 8150:8150 \
		-e DEBUG=true \
		-e OUTPUT_DIR=/app/output \
		-e FFMPEG_PATH=/usr/bin/ffmpeg \
		-e MAX_WORKERS=2 \
		"$(IMAGE):$(TAG)" \
		-l "0.0.0.0:8150" \
		"/mnt/TeslaCam/SentryClips"

# Utility targets
format:
	@echo "Formatting code..."
	gofmt -w .

lint:
	@echo "Running linter..."
	golangci-lint run

# GoReleaser targets
release:
	@echo "Building release with GoReleaser..."
	goreleaser release --clean

release-snapshot:
	@echo "Building release snapshot with GoReleaser..."
	goreleaser release --snapshot --clean

release-dry-run:
	@echo "Dry run release with GoReleaser..."
	goreleaser release --dry-run --clean

# GoReleaser check
check:
	@echo "Checking GoReleaser configuration..."
	goreleaser check

.PHONY: build build-server docker-build docker-build-latest docker-push docker-push-latest test test-short test-coverage clean run docker-run docker-run-dev format lint release release-snapshot release-dry-run check