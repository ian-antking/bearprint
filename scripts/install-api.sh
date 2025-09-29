#!/bin/bash
set -e

USER_NAME="ian-antking"
REPO="bearprint"
BINARY="bearprint-api"
INSTALL_DIR="/var/opt/bearprint"
SERVICE_DEST="/etc/systemd/system/bearprint.service"
SERVICE_USER="bearprint"
CONFIG_PATH="/etc/bearprint/config.ini"

VERSION="bearprint-api-v0.1.8"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH_RAW=$(uname -m)
ARCH=""

case $ARCH_RAW in
  x86_64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  armv6l) ARCH="armv6" ;;
  *) echo "Unsupported architecture: $ARCH_RAW"; exit 1 ;;
esac

FILENAME="${BINARY}-${OS}-${ARCH}"
DOWNLOAD_URL="https://github.com/$USER_NAME/$REPO/releases/download/$VERSION/$FILENAME"

echo "Downloading $BINARY version $VERSION for $OS/$ARCH..."
curl -fsL "$DOWNLOAD_URL" -o "/tmp/$BINARY"

echo "Installing $BINARY to $INSTALL_DIR"
sudo mkdir -p "$INSTALL_DIR"
sudo install -m 755 "/tmp/$BINARY" "$INSTALL_DIR/$BINARY"

echo "Creating user and setting permissions..."
if ! id -u "$SERVICE_USER" >/dev/null 2>&1; then
    sudo useradd -r -s /bin/false "$SERVICE_USER"
fi
sudo chown -R "$SERVICE_USER":"$SERVICE_USER" "$INSTALL_DIR"

echo "Detecting connected USB printers..."
PRINTERS=($(ls /dev/usb/lp* 2>/dev/null || true))

if [ ${#PRINTERS[@]} -eq 0 ]; then
    echo "⚠️  No USB printers detected. Defaulting to /dev/usb/lp0"
    SELECTED_PRINTER="/dev/usb/lp0"
elif [ ${#PRINTERS[@]} -eq 1 ]; then
    SELECTED_PRINTER="${PRINTERS[0]}"
    echo "Found one printer: $SELECTED_PRINTER"
else
    echo "Multiple printers detected:"
    select p in "${PRINTERS[@]}"; do
        if [[ -n "$p" ]]; then
            SELECTED_PRINTER="$p"
            break
        fi
    done
fi

echo "Using printer device: $SELECTED_PRINTER"

echo "Creating config file at $CONFIG_PATH..."
sudo mkdir -p "$(dirname "$CONFIG_PATH")"
sudo tee "$CONFIG_PATH" > /dev/null <<EOF
[printer]
device = $SELECTED_PRINTER
EOF
sudo chown "$SERVICE_USER":"$SERVICE_USER" "$CONFIG_PATH"
sudo chmod 644 "$CONFIG_PATH"

echo "Setting printer permissions..."
if [ -e "$SELECTED_PRINTER" ]; then
    device_group=$(stat -c '%G' "$SELECTED_PRINTER")
    if ! groups "$SERVICE_USER" | grep -qw "$device_group"; then
        sudo usermod -aG "$device_group" "$SERVICE_USER"
    fi
fi

echo "Creating udev rule for printer permissions..."
sudo tee /etc/udev/rules.d/99-bearprint-printer.rules > /dev/null <<EOF
SUBSYSTEM=="usb", ATTRS{idVendor}=="0483", ATTRS{idProduct}=="5743", MODE="0660", GROUP="lp"
EOF

sudo udevadm control --reload-rules
sudo udevadm trigger

echo "Downloading service file to $SERVICE_DEST"
SERVICE_URL="https://raw.githubusercontent.com/$USER_NAME/$REPO/main/bearprint-api/bearprint.service"
sudo curl -fsL "$SERVICE_URL" -o "$SERVICE_DEST"

echo "Reloading systemd daemon and enabling service..."
sudo systemctl daemon-reload
sudo systemctl enable bearprint.service
sudo systemctl restart bearprint.service

echo "✅ BearPrint API system service installed and started!"

if [ -e "$SELECTED_PRINTER" ] && ! groups $USER | grep -qw "$device_group"; then
    echo "⚠️ Please log out and back in, or reboot, for group changes to take effect."
fi
