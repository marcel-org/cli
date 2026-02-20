#!/bin/bash

set -e

REPO="marcel-org/cli"
INSTALL_DIR="$HOME/.local/bin"

echo "Installing Marcel CLI..."

mkdir -p "$INSTALL_DIR"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
    linux|darwin) ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

BINARY_NAME="marcel-${OS}-${ARCH}"

if command -v go >/dev/null 2>&1; then
    echo "Building from source..."
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"

    if command -v git >/dev/null 2>&1; then
        git clone "https://github.com/${REPO}.git" .
        go build -o marcel
    else
        echo "Error: git is required to build from source"
        exit 1
    fi
else
    echo "Go not found. Please install Go or use a pre-built binary."
    exit 1
fi

mv marcel "$INSTALL_DIR/marcel"
chmod +x "$INSTALL_DIR/marcel"

rm -rf "$TEMP_DIR" 2>/dev/null || true

if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
    echo ""
    echo "Add $INSTALL_DIR to your PATH by adding this to your shell config:"
    echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
fi

echo ""
echo "Marcel CLI installed successfully!"
echo ""
echo "Configuration:"
echo "  1. Get your token from Marcel web app settings"
echo "  2. Set environment variable: export MARCEL_TOKEN=\"your_token_here\""
echo "  3. Run: marcel"
echo ""
