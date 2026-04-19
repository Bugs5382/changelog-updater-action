# ISC License
#
# Copyright (c) 2026 Shane & Contributors
#
# Permission to use, copy, modify, and/or distribute this software for any
# purpose with or without fee is hereby granted, provided that the above
# copyright notice and this permission notice appear in all copies.
#
# THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
# WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
# MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
# ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
# WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
# ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
# OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

# syntax=docker/dockerfile:1.7

# ---------- Build stage ----------
FROM --platform=$BUILDPLATFORM golang:1.26.1-alpine AS build

# Build-time metadata (populated by CI; safe defaults for local builds).
ARG VERSION=local
ARG GITSHA=unknown
ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

# Cache modules separately so code-only changes don't re-download deps.
COPY go.mod go.sum* ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy the rest of the source.
COPY . .

# Produce a fully static, stripped binary for the requested target platform.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
        -trimpath \
        -ldflags "-s -w -X 'main.Version=${VERSION}' -X 'main.Gitsha=${GITSHA}'" \
        -o /out/changelog-updater-action \
        ./cmd/action

# ---------- Runtime stage ----------
FROM gcr.io/distroless/static-debian12:nonroot

# OCI labels — populated by CI, harmless when left as defaults.
ARG VERSION=local
ARG GITSHA=unknown
LABEL org.opencontainers.image.title="changelog-updater-action" \
      org.opencontainers.image.description="Updates CHANGELOG.md with release notes for GitHub Actions and GitLab CI." \
      org.opencontainers.image.source="https://github.com/Bugs5382/changelog-updater-action" \
      org.opencontainers.image.licenses="ISC" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.revision="${GITSHA}"

# GitHub Actions mounts the repo at /github/workspace and cd's into it.
# GitLab CI does the equivalent via CI_PROJECT_DIR. Defaulting WORKDIR to
# /github/workspace keeps `--path .` working out of the box on GitHub, and
# GitLab jobs can simply pass `--path "$CI_PROJECT_DIR"` (or `docker run -w`).
WORKDIR /github/workspace

COPY --from=build /out/changelog-updater-action /usr/local/bin/changelog-updater-action

# The binary parses its own flags via pflag, so make it the entrypoint and
# let callers append whatever flags they need. No shell, no wrapper.
ENTRYPOINT ["/usr/local/bin/changelog-updater-action"]