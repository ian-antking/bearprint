#!/bin/bash
set -e

# ----------------------------
# BearPrint Install Script
# ----------------------------

USER_NAME="ian-antking"
REPO="bearprint"
BINARY="bearprint-api"
INSTALL_DIR="/var/opt/bearprint"
SERVICE_DEST="/etc/systemd/system/bearprint.service"
SERVICE_USER="bearprint"
CONFIG_PATH="/etc/bearprint/config.ini"
VERSION="bearprint-api-v0.1.9"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH_RAW=$(uname -m)
ARCH=""

echo "⚠️ You may be prompted for your password to install BearPrint..."
sudo -v

case $ARCH_RAW in
  x86_64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  armv6l) ARCH="armv6" ;;
  *) echo "Unsupported architecture: $ARCH_RAW"; exit 1 ;;
esac

# ----------------------------
# Create system user first
# ----------------------------
if ! id -u "$SERVICE_USER" >/dev/null 2>&1; then
    echo "Creating system user '$SERVICE_USER'..."
    sudo useradd -r -s /bin/false "$SERVICE_USER"
fi

# ----------------------------
# Detect connected USB printers
# ----------------------------
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
        else
            echo "Invalid selection, try again."
        fi
    done
fi

echo "Using printer device: $SELECTED_PRINTER"

echo "Enter a name for this printer (default: bearprint):"
read -r PRINTER_NAME < /dev/tty
PRINTER_NAME=${PRINTER_NAME:-bearprint}
echo "Using printer name: $PRINTER_NAME"

# ----------------------------
# Create config file
# ----------------------------
echo "Creating config file at $CONFIG_PATH..."
sudo mkdir -p "$(dirname "$CONFIG_PATH")"
sudo tee "$CONFIG_PATH" > /dev/null <<EOF
[printer]
device = $SELECTED_PRINTER
name = $PRINTER_NAME
EOF
sudo chown "$SERVICE_USER":"$SERVICE_USER" "$CONFIG_PATH"
sudo chmod 644 "$CONFIG_PATH"

# ----------------------------
# Download and install binary
# ----------------------------
FILENAME="${BINARY}-${OS}-${ARCH}"
DOWNLOAD_URL="https://github.com/$USER_NAME/$REPO/releases/download/$VERSION/$FILENAME"
TMP_BINARY="/tmp/$BINARY"
INSTALL_TMP="$INSTALL_DIR/$BINARY.tmp"

echo "Downloading $BINARY version $VERSION for $OS/$ARCH..."
curl -fsL "$DOWNLOAD_URL" -o "$TMP_BINARY"

# Safety check: ensure binary is non-empty
if [ ! -s "$TMP_BINARY" ]; then
    echo "❌ Download failed or binary is empty, aborting installation."
    echo "Attempted URL: $DOWNLOAD_URL"
    head -n 5 "$TMP_BINARY"
    rm -f "$TMP_BINARY"
    exit 1
fi

echo "Installing $BINARY to $INSTALL_DIR..."
sudo mkdir -p "$INSTALL_DIR"

# Copy to temporary file in install dir first
sudo cp "$TMP_BINARY" "$INSTALL_TMP"
sudo chmod 755 "$INSTALL_TMP"
sudo chown "$SERVICE_USER":"$SERVICE_USER" "$INSTALL_TMP"

# Atomically replace old binary
sudo mv "$INSTALL_TMP" "$INSTALL_DIR/$BINARY"

# Clean up temporary download
rm -f "$TMP_BINARY"

# ----------------------------
# Printer permissions
# ----------------------------
if [ -e "$SELECTED_PRINTER" ]; then
    device_group=$(stat -c '%G' "$SELECTED_PRINTER")
    if ! groups "$SERVICE_USER" | grep -qw "$device_group"; then
        echo "Adding user '$SERVICE_USER' to printer group '$device_group'..."
        sudo usermod -aG "$device_group" "$SERVICE_USER"
    fi
fi

# ----------------------------
# Create generic udev rule
# ----------------------------
echo "Creating udev rule for printer permissions..."
sudo tee /etc/udev/rules.d/99-bearprint-printer.rules > /dev/null <<EOF
SUBSYSTEM=="usb", KERNEL=="lp*", MODE="0660", GROUP="$device_group"
EOF

sudo udevadm control --reload-rules
sudo udevadm trigger

# ----------------------------
# Install systemd service
# ----------------------------
echo "Downloading service file to $SERVICE_DEST..."
SERVICE_URL="https://raw.githubusercontent.com/$USER_NAME/$REPO/main/bearprint-api/bearprint.service"
TEMP_SERVICE="/tmp/bearprint.service"

sudo curl -fsL "$SERVICE_URL" -o "$TEMP_SERVICE"

# Safety check: ensure service file is non-empty
if [ ! -s "$TEMP_SERVICE" ]; then
    echo "❌ Service file download failed or file is empty, aborting installation."
    rm -f "$TEMP_SERVICE"
    exit 1
fi

# Move valid service file into place
sudo mv "$TEMP_SERVICE" "$SERVICE_DEST"

echo "Reloading systemd daemon and enabling service..."
sudo systemctl daemon-reload
sudo systemctl enable bearprint.service
sudo systemctl restart bearprint.service

echo "✅ BearPrint API system service installed and started!"

# ----------------------------
# Reminder for group changes
# ----------------------------
if [ -e "$SELECTED_PRINTER" ] && ! groups $(whoami) | grep -qw "$device_group"; then
    echo "⚠️ Please log out and back in, or reboot, for group changes to take effect."
fi
