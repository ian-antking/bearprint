# 🐻 BearPrint

BearPrint is a tiny, networked thermal printer stack — perfect for receipts, notes, or weird little projects. Built to run on a Raspberry Pi and scale with a UI in a k3s cluster.

## ✨ Features

- 📜 Simple Python script to send text to a thermal printer
- 🌐 Flask-based API to expose printing over your local network
- 🧾 CLI tool (`cat something.txt | bearprint`) to send jobs from any device

## 🚀 Getting Started

### Requirements

- Raspberry Pi Zero (or similar) with USB thermal printer
- Python 3.9+
- `sudo` access to `/dev/usb/lp0` or equivalent

### Installation (CLI Tool)

Run the following command in your terminal. It will automatically detect your operating system and architecture, then download and install the `bearprint` binary to `~/.local/bin`.

```bash
curl -sSL https://raw.githubusercontent.com/ian-antking/bear-print/main/scripts/install-cli.sh | bash
```

> **Note**: On macOS, you may need to manually approve the binary after installation. Navigate to `~/.local/bin` in Finder, right-click `bearprint`, and select "Open". Or run `xattr -d com.apple.quarantine ~/.local/bin/bearprint
`

### Quickstart (Server)

```bash
# Clone the repo
git clone [https://github.com/ian-antking/bear-print.git](https://github.com/ian-antking/bear-print.git)
cd bear-print/bearprint-server

# Install dependencies
pip install -r requirements.txt

# Run the server
make dev
```

### Print a test message

Using the CLI:

```bash
echo "Hello from BearPrint!" | bearprint
```

Using cURL:

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
