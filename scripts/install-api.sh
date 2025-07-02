#!/bin/bash
set -e

USER_NAME="ian-antking"
REPO="bear-print"
BINARY="bearprint-api"
INSTALL_DIR="/var/opt/bearprint"
SERVICE_SOURCE="bearprint-api/bearprint.service"
SERVICE_DEST="/etc/systemd/system/bearprint.service"

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

DOWNLOAD_URL="https://github.com/$USER_NAME/$REPO/releases/latest/download/${BINARY}-${OS}-${ARCH}"

echo "Downloading $BINARY for $OS/$ARCH..."
curl -sL "$DOWNLOAD_URL" -o "/tmp/$BINARY"

echo "Installing $BINARY to $INSTALL_DIR"
sudo mkdir -p "$INSTALL_DIR"
sudo install -m 755 "/tmp/$BINARY" "$INSTALL_DIR/$BINARY"

echo "Setting printer permissions..."
device_group=$(stat -c '%G' /dev/usb/lp0)
if ! groups $USER | grep -qw "$device_group"; then
    echo "Adding $USER to $device_group group..."
    sudo usermod -aG "$device_group" "$USER"
    echo "Added $USER to $device_group group. You may need to log out and back in."
fi

echo "Creating udev rule for printer permissions..."
sudo tee /etc/udev/rules.d/99-usb-printer.rules > /dev/null <<EOF
KERNEL=="lp0", MODE="0666"
EOF

sudo udevadm control --reload-rules
sudo udevadm trigger

echo "Copying service file to $SERVICE_DEST"
sudo cp "$SERVICE_SOURCE" "$SERVICE_DEST"

echo "Reloading systemd daemon and enabling service..."
sudo systemctl daemon-reload
sudo systemctl enable bearprint.service
sudo systemctl restart bearprint.service

echo "✅ BearPrint API system service installed and started!"

if groups $USER | grep -qw "$device_group"; then
    echo "You may now use BearPrint API without rebooting."
else
    echo "⚠️ Please log out and back in, or reboot, for group changes to take effect."
fi
