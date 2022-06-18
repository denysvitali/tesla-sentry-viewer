FROM golang:1.18-alpine AS builder
RUN apk add --no-cache make bash
ARG VERSION=dev
COPY / /app
WORKDIR /app
RUN make build

FROM scratch
COPY --from=builder /app/build/tesla-sentry-viewer /
ENTRYPOINT ["/tesla-sentry-viewer"]