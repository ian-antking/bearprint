# ʕ•ᴥ•ʔ BearPrint

BearPrint is a tiny, networked thermal printer stack — perfect for receipts, notes, or weird little projects. Built to run on a Raspberry Pi zero.

## ✨ Features

- 🌐 API to expose printing over your local network
- 🧾 CLI tool (`cat something.txt | bearprint`) to send jobs within your network

## 🚀 Getting Started

### Requirements

To build a BearPrint server, you’ll need:

- 🐻 **Raspberry Pi Zero (or Zero W)**  
  - You may prefer the vanilla Pi Zero if you're connecting via Ethernet.
- 🌐 **Waveshare Ethernet + USB Hub HAT** *(optional — only needed for wired networking)*  
  - [Amazon UK link](https://www.amazon.co.uk/dp/B09K5DLR17)
- 🖨️ **USB Thermal Printer** (e.g. Xprinter 80T with auto cutter)  
  - [AliExpress link](https://a.aliexpress.com/_EQoGyOO)
- 💾 **Micro SD card**  
  - Doesn’t need to be large — 8GB+ is fine.
- 🔌 **Micro-USB power cable**
- 🧵 **Ethernet cable**

> The software is designed for maker-style setups and open source tinkering. No cloud connection required.

---

### Server Installation (on your Raspberry Pi)

Run the following command on your Raspberry Pi. It will automatically download the correct binary, install it as a `systemd` service, and set the necessary permissions.

```bash
curl -sSL https://raw.githubusercontent.com/ian-antking/bearprint/main/scripts/install-api.sh | bash
```

This single command installs and starts the server. The server will automatically run on boot.

---

### CLI Installation (on your other computers)

Run the following command on any Mac or Linux machine on your network. It will automatically detect the OS and architecture, then download and install the `bearprint` binary to `~/.local/bin`.

```bash
curl -sSL https://raw.githubusercontent.com/ian-antking/bearprint/main/scripts/install-cli.sh | bash
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

### iOS Shortcut (optional)

Want to print directly from Safari or Notes on your iPhone or iPad? Use the BearPrint iOS Shortcut:

👉 [Install the BearPrint Shortcut](https://www.icloud.com/shortcuts/243fe324569f40ed856e326eb42bfc5f)

This shortcut lets you:

- Send article text from Safari or Notes app share sheet
- Send ad hoc text by invoking the shortcut directly
- Automatically formats text with blank lines and a final cut

> **Note**: After installing, tap the shortcut's "..." icon and edit it to add your BearPrint server URL.

## 🐾 Logo

```text
   ʕ•ᴥ•ʔ
 BearPrint
```

## 📜 License

MIT
