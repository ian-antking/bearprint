#!/bin/bash

set -e

# Install system-wide gunicorn if not already installed
if ! command -v gunicorn >/dev/null 2>&1; then
    echo "Installing gunicorn system-wide..."
    sudo apt update
    sudo apt install gunicorn -y
fi

# Install Python dependencies locally (still useful)
python3 -m pip install --break-system-packages --upgrade pip
python3 -m pip install --break-system-packages -r requirements.txt

# Set printer permissions
device_group=$(stat -c '%G' /dev/usb/lp0)
if ! groups $USER | grep -qw "$device_group"; then
    sudo usermod -aG "$device_group" "$USER"
    echo "Added $USER to $device_group group. Please log out and back in."
fi

# udev rule for printer permissions
sudo tee /etc/udev/rules.d/99-usb-printer.rules > /dev/null <<EOF
KERNEL=="lp0", MODE="0666"
EOF

sudo udevadm control --reload-rules
sudo udevadm trigger

# Install systemd system service (adjust path to service file if needed)
sudo cp bearprint.service /etc/systemd/system/bearprint.service

# Reload systemd, enable and start service
sudo systemctl daemon-reload
sudo systemctl enable bearprint.service
sudo systemctl restart bearprint.service

echo "✅ BearPrint system service installed and started!"

# Recommend reboot if groups were changed
if groups $USER | grep -qw "$device_group"; then
    echo "You may now use BearPrint without rebooting."
else
    echo "⚠️  Please log out and back in, or reboot, for group changes to take effect."
fi
