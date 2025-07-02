# 🐻 BearPrint CLI

`bearprint` is a command-line tool for sending print jobs to a BearPrint server. It reads data from standard input, making it easy to pipe text or other content directly to your thermal printer.

## Installation

Run the following command in your terminal. It will automatically detect your operating system and architecture, then download and install the `bearprint` binary to `~/.local/bin`.

```bash
curl -sSL https://raw.githubusercontent.com/ian-antking/bear-print/main/scripts/install-cli.sh | bash
```

> **Note**: On macOS, you may need to manually approve the binary after installation. Navigate to `~/.local/bin` in Finder, right-click `bearprint`, and select "Open", Or run `xattr -d com.apple.quarantine ~/.local/bin/bearprint`.

## Configuration

The CLI can be configured using a file or command-line flags.

### Config File (Recommended)

For persistent configuration, create a file at `~/.bearprint/config`:

1. **Create the directory**:

  ```bash
  mkdir -p ~/.bearprint
  ```

1. **Create the config file**:

  ```ini
  [default]
  server_host = your_pi_ip
  server_port = 8080
  ```

Replace `your_pi_ip` with the IP address of your Raspberry Pi.

## Usage

The CLI reads from standard input, allowing you to pipe content from other commands.

### Printing Text

By default, each line of input is sent as a separate text item.

```bash
# Print a single line
echo "Hello, world!" | bearprint

# Print the contents of a file
cat receipt.txt | bearprint
```

### Printing QR Codes

Use the `-q` or `--qrcode` flag to treat the entire input as a single item to be encoded as a QR code.

```bash
echo "[https://example.com](https://example.com)" | bearprint --qrcode
```

### Command-Line Flags

You can override the settings from the config file for a single command by using flags.

```bash
# Connect to a different server for one command
echo "Temporary print job" | bearprint --host 10.0.1.50 --port 9999
```

### Getting Help

To see all available command-line flags and usage information, use the `-h` or `--help` flag.

```bash
$ bearprint -h
Usage of bearprint:
  -host string
     Server host
  -port string
     Server port
  -q Treat input as a single QR code
  -qrcode
     Treat input as a single QR code
```
