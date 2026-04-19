FROM alpine:3.20

# Install dependencies
RUN apk --no-cache add curl ca-certificates

# Set the default arch to amd64 but allow overrides
ARG TARGETARCH=amd64

# Download the binary based on the architecture
# Note: Ensure your release filenames match the naming convention below
RUN curl -sL "https://github.com/Bugs5382/changelog-updater-action/releases/download/latest/changelog-updater-action-linux-${TARGETARCH}" \
    -o /usr/local/bin/changelog-updater-action && \
    chmod +x /usr/local/bin/changelog-updater-action

ENTRYPOINT ["/usr/local/bin/changelog-updater-action"]