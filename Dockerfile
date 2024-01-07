# syntax = docker/dockerfile:1-experimental

# ==================================
# Base image
# ===================================
FROM --platform=$BUILDPLATFORM golang:1.21-alpine3.19 as base

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

# ==================================
# Builder image
# ===================================
FROM base as builder

ARG BUILDTIME
ARG VERSION
ARG REVISION
ARG GC_FLAGS

ARG TARGETOS
ARG TARGETARCH

RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -gcflags "${GC_FLAGS}" \
    -ldflags="-w -s \
        -X main.version=${VERSION} \
        -X main.gitCommit=${REVISION} \
        -X main.buildDate=${BUILDTIME} \
        -extldflags '-static'" \
    -o /go/bin/github-notion-stars-sync cmd/main.go

# ===================================
# production image
# ===================================
FROM alpine:3.19 as production

ARG UID=1000
ARG GID=1000

# hadolint ignore=DL3018
RUN apk add --no-cache curl ca-certificates && \
    addgroup -g ${GID} app && \
    adduser -D -u ${UID} -G app app

COPY --from=builder --chown=app:app /go/bin/github-stars-notion-sync /usr/local/bin/github-stars-notion-sync

ENTRYPOINT [ "/usr/local/bin/github-stars-notion-sync" ]

USER app

LABEL org.opencontainers.image.title "Github Stars Notion Sync"
LABEL org.opencontainers.image.description "CLI tool to sync Github Stars to a Notion Database"
LABEL org.opencontainers.image.authors "Bruno Paz"
LABEL org.opencontainers.image.url "https://github.com/brpaz/github-stars-notion-sync"
