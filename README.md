# 🐻 BearPrint

BearPrint is a tiny, networked thermal printer stack — perfect for receipts, notes, or weird little projects. Built to run on a Raspberry Pi and scale with a UI in a k3s cluster.

## ✨ Features

- 📜 Simple Python script to send text to a thermal printer
- 🌐 Flask-based API to expose printing over your local network
- 🧾 CLI tool (`cat something.txt | thermal-print`) to send jobs from any device
- ⚛️ Preact frontend for formatting and sending print jobs
- 🔐 Planned authentication & OAuth support
- 📦 Monorepo layout for easy management

## 🧱 Architecture

```mermaid
graph TD
  CLI[CLI Tool]
  API[Flask API<br/>(Print Server)]
  PRINTER[Thermal Printer]

  CLI -->|HTTP| API
  API -->|USB or Serial| PRINTER
```

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
python app.py
```

### Print a test message

```bash
curl -X POST http://your-pi-ip:5000/v1/print/text \
  -H "Content-Type: application/json" \
  -d '{"text": "Hello from BearPrint!"}'
```

## 🧪 Project Structure

```text
bearprint/
├── printer-server/      # Flask API backend
├── printer-ui/          # Preact frontend (WIP)
├── cli-tool/            # Simple shell/Node-based CLI (WIP)
├── shared/              # Shared constants/types
└── README.md            # You're here!
```

## 💡 Future Ideas

- `POST /v1/print/image`
- `POST /v1/print/composite`
- Authenticated dashboard with print history
- QR code printing
- Emoji/art templates

## 🐾 Logo

```text
   ʕ•ᴥ•ʔ
 BearPrint
```

## 📜 License

MIT
