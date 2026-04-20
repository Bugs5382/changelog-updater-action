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
FROM alpine:3.20

# Install dependencies
RUN apk --no-cache add curl ca-certificates

# Set the default arch to amd64 but allow overrides
ARG TARGETARCH=amd64

# Download the binary based on the architecture
# Note: Ensure your release filenames match the naming convention below
RUN curl -sL "https://github.com/Bugs5382/changelog-updater-action/releases/latest/download/changelog-updater-action-linux-${TARGETARCH}" \
    -o /usr/local/bin/changelog-updater-action && \
    chmod +x /usr/local/bin/changelog-updater-action

ENTRYPOINT ["/usr/local/bin/changelog-updater-action"]