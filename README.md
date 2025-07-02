# ʕ•ᴥ•ʔ BearPrint

BearPrint is a tiny, networked thermal printer stack — perfect for receipts, notes, or weird little projects. Built to run on a Raspberry Pi zero.

## ✨ Features

- 🌐 Flask-based API to expose printing over your local network
- 🧾 CLI tool (`cat something.txt | bearprint`) to send jobs within your network

## 🚀 Getting Started

### Requirements

- Raspberry Pi Zero (or similar) with USB thermal printer (I'm using an Xprinter 80T)
- Python 3.9+
- `sudo` access to `/dev/usb/lp0` or equivalent

### Quickstart (Server)

#### 1. Clone the Repository Locally

Get the project files on your local development machine.

```bash
git clone https://github.com/ian-antking/bear-print.git)
cd bear-print/bearprint-server
```

#### 2. Deploy Files to the Pi

From your local machine, run the following command. This will sync the project files to `/opt/bearprint-server/` on the Raspberry Pi.

```bash
make deploy USER=your_pi_user HOST=your_pi_ip
```

#### 3. Install Dependencies on the Pi

SSH into your Raspberry Pi and run the `install` command from the deployment directory. This will install the necessary Python packages and set up the systemd service.

```bash
# SSH into the Pi
ssh your_pi_user@your_pi_ip

# Navigate to the directory and run install
cd /opt/bearprint-server
make install
```

### Installation (CLI Tool)

Run the following command in your terminal. It will automatically detect your operating system and architecture, then download and install the `bearprint` binary to `~/.local/bin`.

```bash
curl -sSL https://raw.githubusercontent.com/ian-antking/bear-print/main/scripts/install-cli.sh | bash
```

> **Note**: On macOS, you may need to manually approve the binary after installation. Navigate to `~/.local/bin` in Finder, right-click `bearprint`, and select "Open", Or run `xattr -d com.apple.quarantine ~/.local/bin/bearprint`.

#### Configuration

After installing, you must create a configuration file for the CLI tool to specify the server's address.

1. Create a file named `~/.bearprint/config` with the following content:

```ini
[default]
server_host = your_pi_ip
server_port = 8080
```

Replace `your_pi_ip` with the IP address of your Raspberry Pi.

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
