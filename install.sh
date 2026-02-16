#!/bin/bash
# HiveClaw Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/NivedChowdary/hiveclaw/main/install.sh | bash

set -e

VERSION="0.1.0"
REPO="NivedChowdary/hiveclaw"

echo "🐝 Installing HiveClaw v${VERSION}..."
echo ""

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
esac

echo "   OS: $OS"
echo "   Arch: $ARCH"

# Determine install directory
if [ -w "/usr/local/bin" ]; then
    INSTALL_DIR="/usr/local/bin"
else
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
fi

# For now, we need Go to build from source
# In future, will download pre-built binaries from releases

if ! command -v go &> /dev/null; then
    echo ""
    echo "📦 Go not found. Installing Go first..."
    
    GO_VERSION="1.22.0"
    GO_URL="https://go.dev/dl/go${GO_VERSION}.${OS}-${ARCH}.tar.gz"
    
    mkdir -p ~/go-install
    curl -fsSL "$GO_URL" | tar -C ~/go-install -xzf -
    export PATH=$PATH:~/go-install/go/bin
    
    echo "   Go installed to ~/go-install/go"
fi

echo ""
echo "📥 Cloning HiveClaw..."
TMPDIR=$(mktemp -d)
cd "$TMPDIR"
git clone --depth 1 https://github.com/${REPO}.git hiveclaw
cd hiveclaw

echo ""
echo "🔨 Building..."
go build -ldflags "-s -w" -o hiveclaw ./cmd/hiveclaw

echo ""
echo "📦 Installing to ${INSTALL_DIR}..."
mv hiveclaw "${INSTALL_DIR}/"
chmod +x "${INSTALL_DIR}/hiveclaw"

# Cleanup
cd /
rm -rf "$TMPDIR"

echo ""
echo "✅ HiveClaw installed successfully!"
echo ""
echo "Next steps:"
echo "  1. Run setup:    hiveclaw onboard"
echo "  2. Start server: hiveclaw start"
echo "  3. Open:         http://localhost:8080"
echo ""

# Add to PATH hint if needed
if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
    echo "💡 Add to PATH: export PATH=\"\$PATH:${INSTALL_DIR}\""
    echo ""
fi

echo "🐝 Welcome to the Hive!"
