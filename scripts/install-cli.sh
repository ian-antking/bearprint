#!/bin/sh
set -e

USER_NAME="ian-antking"
REPO="bearprint"
BINARY="bearprint"
INSTALL_DIR="$HOME/.local/bin"

VERSION="bearprint-cli-v1.2.0"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
  x86_64)
    ARCH="amd64"
    ;;
  arm64 | aarch64)
    ARCH="arm64"
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

DOWNLOAD_URL="https://github.com/$USER_NAME/$REPO/releases/download/${VERSION}/${BINARY}-${OS}-${ARCH}"

echo "Downloading $BINARY version $VERSION for $OS/$ARCH..."
curl -sL "$DOWNLOAD_URL" -o "/tmp/$BINARY"

echo "Installing $BINARY to $INSTALL_DIR"
mkdir -p "$INSTALL_DIR"
install -m 755 "/tmp/$BINARY" "$INSTALL_DIR/$BINARY"

echo "$BINARY installed successfully!"
echo "Make sure '$INSTALL_DIR' is in your PATH."