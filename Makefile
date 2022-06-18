IMAGE=dvitali/tesla-sentry-viewer
VERSION=$(shell ./get-version.sh)
TAG=$(VERSION)

build:
	CGO_ENABLED=0 \
	go build \
		-ldflags "-X main.Version=$(VERSION)" \
		-o ./build/tesla-sentry-viewer ./cmd/server

docker-build:
	docker build \
		-t "$(IMAGE):$(TAG)" \
		.

docker-push:
	docker push "$(IMAGE):$(TAG)"

docker-run:
	docker run \
		--rm \
		--name "tesla-sentry-viewer" \
		-v "/run/media/$$USER/TESLADRIVE/:/mnt:ro" \
		-p 8150:8150 \
		"$(IMAGE):$(TAG)" \
		-l "0.0.0.0:8150" \
		"/mnt/TeslaCam/SentryClips"

.PHONY: build