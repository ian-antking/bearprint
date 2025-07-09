---
layout: home
---

# ʕ•ᴥ•ʔ BearPrint

**BearPrint** is a tiny local-first print server that lets you send text from your phone or laptop to a USB-connected thermal printer — no cloud required.

Perfect for to-do lists, shopping lists, quotes, notes to self, or tiny zines.

---

## ✨ Features

- Share text from Safari, Notes, or any app via the iOS share sheet
- Print ad hoc text via a shortcut
- Local network only — no accounts, no tracking
- Easy to run on a Raspberry Pi Zero

---

## 🛠️ Requirements

- Raspberry Pi Zero (or Zero W)
- USB thermal printer (e.g. [Xprinter 80T](https://a.aliexpress.com/_EQoGyOO))
- Optional: Ethernet + USB hub HAT ([Waveshare](https://www.amazon.co.uk/dp/B09K5DLR17))
- microSD card (8GB+)
- micro-USB power cable
- Ethernet cable (optional)

---

## 📦 Setup

To get started:

1. Flash Raspberry Pi OS Lite to your SD card
2. Clone this repo and install dependencies
3. Start the server with `bearprint serve`
4. [Install the iOS Shortcut](shortcut.html) and configure your server address

More details in the [Installation Guide »](install.html)

---

## 📲 iOS Shortcut

A Shortcuts-friendly way to send text to your BearPrint server.  
[View setup instructions »](shortcut.html)

---

## 📚 API

Want to integrate BearPrint into another app? Use the simple HTTP API.  
[See API reference »](api.html)

---

## 🧵 About

BearPrint is designed for hobbyists and tinkerers.  
It's open source, hackable, and intentionally simple.

[View on GitHub »](https://github.com/yourusername/bearprint)
