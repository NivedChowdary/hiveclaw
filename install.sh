#!/bin/bash
# HiveClaw Installer
# Usage: curl -fsSL https://hiveclaw.dev/install.sh | sh

set -e

VERSION="0.1.0"
REPO="nanilabs/hiveclaw"

echo "🐝 Installing HiveClaw v${VERSION}..."

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Determine install directory
INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
fi

# Download binary
BINARY_URL="https://github.com/${REPO}/releases/download/v${VERSION}/hiveclaw-${OS}-${ARCH}"
echo "Downloading from ${BINARY_URL}..."

if command -v curl &> /dev/null; then
    curl -fsSL "$BINARY_URL" -o "${INSTALL_DIR}/hiveclaw"
elif command -v wget &> /dev/null; then
    wget -q "$BINARY_URL" -O "${INSTALL_DIR}/hiveclaw"
else
    echo "Error: curl or wget required"
    exit 1
fi

chmod +x "${INSTALL_DIR}/hiveclaw"

# Verify installation
if [ -x "${INSTALL_DIR}/hiveclaw" ]; then
    echo ""
    echo "✅ HiveClaw installed successfully!"
    echo ""
    echo "Next steps:"
    echo "  1. Run setup wizard:  hiveclaw onboard"
    echo "  2. Start gateway:     hiveclaw start"
    echo "  3. Open dashboard:    http://localhost:8080"
    echo ""
    echo "🐝 Welcome to the Hive!"
else
    echo "❌ Installation failed"
    exit 1
fi

# Add to PATH if needed
if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
    echo ""
    echo "Note: Add ${INSTALL_DIR} to your PATH:"
    echo "  export PATH=\"\$PATH:${INSTALL_DIR}\""
fi
