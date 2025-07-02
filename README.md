# 🐻 BearPrint

BearPrint is a tiny, networked thermal printer stack — perfect for receipts, notes, or weird little projects. Built to run on a Raspberry Pi and scale with a UI in a k3s cluster.

## ✨ Features

- 📜 Simple Python script to send text to a thermal printer
- 🌐 Flask-based API to expose printing over your local network
- 🧾 CLI tool (`cat something.txt | thermal-print`) to send jobs from any device

## 🚀 Getting Started

### Requirements

- Raspberry Pi Zero (or similar) with USB thermal printer
- Python 3.9+
- `sudo` access to `/dev/usb/lp0` or equivalent

### Quickstart

```bash
# Clone the repo
git clone https://github.com/youruser/bearprint.git
cd bearprint/printer-server

# Install dependencies
pip install -r requirements.txt

# Run the server
cd bearprint-server
make dev
```

### Print a test message

```bash
curl -X POST http://your-pi-ip:8080/api/v1/print/text \
  -H "Content-Type: application/json" \
  -d '{"text": "Hello from BearPrint!"}'
```

## 🐾 Logo

```text
   ʕ•ᴥ•ʔ
 BearPrint
```

## 📜 License

MIT
