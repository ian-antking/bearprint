# ʕ•ᴥ•ʔ BearPrint

BearPrint is a tiny, networked thermal printer stack — perfect for receipts, notes, or weird little projects. Built to run on a Raspberry Pi zero.

## ✨ Features

- 🌐 API to expose printing over your local network
- 🧾 CLI tool (`cat something.txt | bearprint`) to send jobs within your network

## 🚀 Getting Started

### Requirements

- Raspberry Pi Zero (or similar) with a USB thermal printer (e.g., Xprinter 80T)
- A `systemd`-based Linux distribution (like Raspberry Pi OS)
- `sudo` access for installation

---

### Server Installation (on your Raspberry Pi)

Run the following command on your Raspberry Pi. It will automatically download the correct binary, install it as a `systemd` service, and set the necessary permissions.

```bash
curl -sSL https://raw.githubusercontent.com/ian-antking/bear-print/main/scripts/install-api.sh | bash
```

This single command installs and starts the server. The server will automatically run on boot.

---

### CLI Installation (on your other computers)

Run the following command on any Mac or Linux machine on your network. It will automatically detect the OS and architecture, then download and install the `bearprint` binary to `~/.local/bin`.

```bash
curl -sSL https://raw.githubusercontent.com/ian-antking/bear-print/main/scripts/install-cli.sh | bash
```

> **Note**: On macOS, you may need to manually approve the binary after installation. Navigate to `~/.local/bin` in Finder, right-click `bearprint`, and select "Open", or run `xattr -d com.apple.quarantine ~/.local/bin/bearprint`.

#### CLI Configuration

After installing the CLI, you must create a configuration file to point it to your server.

1. Create a file named `~/.bearprint/config` with the following content:

    ```ini
    [default]
    server_host = your_pi_ip
    server_port = 8080
    ```

Replace `your_pi_ip` with the IP address of your Raspberry Pi.

---

### Print a Test Message

Once the server is running on the Pi and the CLI is configured on another machine, you can send a print job.

**Using the CLI:**

```bash
echo "Hello from BearPrint!" | bearprint
```

**Using cURL:**

```bash
curl -X POST http://your-pi-ip:8080/api/v1/print \
  -H "Content-Type: application/json" \
  -d '{"items": [{"type": "text", "content": "Hello from cURL!"}]}'
```

## 🐾 Logo

```text
   ʕ•ᴥ•ʔ
 BearPrint
```

## 📜 License

MIT
