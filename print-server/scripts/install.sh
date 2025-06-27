#!/bin/bash

if ! grep -qxF 'export PATH="$HOME/.local/bin:$PATH"' ~/.bashrc; then
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
    echo "Added ~/.local/bin to PATH in ~/.bashrc"
fi

sudo apt update
sudo apt install python3-pip -y

python3 -m pip install --break-system-packages --upgrade pip
python3 -m pip install --break-system-packages -r requirements.txt

device_group=$(stat -c '%G' /dev/usb/lp0)
if ! groups $USER | grep -qw "$device_group"; then
    sudo usermod -aG "$device_group" "$USER"
    echo "Added $USER to $device_group group. Please log out and back in."
fi

sudo tee /etc/udev/rules.d/99-usb-printer.rules > /dev/null <<EOF
KERNEL=="lp0", MODE="0666"
EOF

sudo udevadm control --reload-rules
sudo udevadm trigger

sudo reboot
