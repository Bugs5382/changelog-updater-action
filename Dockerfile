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
FROM golang:1.26.1-alpine

LABEL org.opencontainers.image.title="changelog-updater-action" \
      org.opencontainers.image.description="A Go Script that will update a CHANGELOG.md file with advanced settings and actions." \
      org.opencontainers.image.source="https://github.com/Bugs5382/changelog-updater-action" \
      org.opencontainers.image.licenses="ISC"

# Install curl to download the binary
RUN apk --no-cache add curl ca-certificates

# Get Binary
ARG ARCH=amd64
ENV BINARY_URL="https://github.com/Bugs5382/changelog-updater-action/releases/download/latest/changelog-updater-action-linux-${ARCH}"

# Download and set permissions
RUN curl -sL ${BINARY_URL} -o /usr/local/bin/changelog-updater-action && \
    chmod +x /usr/local/bin/changelog-updater-action

# The binary parses its own flags via pflag, so make it the entrypoint and
# let callers append whatever flags they need. No shell, no wrapper.
ENTRYPOINT ["/usr/local/bin/changelog-updater-action"]