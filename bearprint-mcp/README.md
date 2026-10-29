# ʕ•ᴥ•ʔ BearPrint MCP Server (optional)

An optional, self-contained [Model Context Protocol](https://modelcontextprotocol.io) server that lets an AI assistant (for example ChatGPT) print through a BearPrint server. Nothing else in this repository depends on it.

It speaks MCP over stdio (newline-delimited JSON-RPC 2.0), uses only the Python 3 standard library, and opens no network port. It calls the BearPrint HTTP API, so it can run on any machine that can reach your BearPrint server, usually the Pi itself.

## Tools

| Tool | What it does |
|------|--------------|
| `bearprint_health` | Calls `GET /api/v1/health` and reports whether the server is up. |
| `bearprint_print` | Validates a list of items and, only when `dry_run` is `false`, sends them to `POST /api/v1/print`. |

The item types are the ones the API accepts (see [`bearprint-api/README.md`](../bearprint-api/README.md)): `text`, `qrcode`, `blank`, `line` and `cut`, with optional `align` (`left`, `center`, `right`) and `count` (`blank` only).

Extra limits the MCP server applies: at most 50 items per job, at most 500 characters of `content`, and `count` from 1 to 10.

### Safety

- `dry_run` defaults to `true`. A job is only printed when `dry_run` is the boolean `false`; anything else (missing, `null`, the string `"false"`) is treated as a dry run.
- Unknown fields and item types are rejected before anything reaches the printer.
- There is no authentication in the server itself. **Anyone who can call the tool can print**, so limit who can reach it (see the tunnel section below).

## Run it

```bash
BEARPRINT_URL=http://your-pi-ip:8080 python3 bearprint-mcp/server.py
```

`BEARPRINT_URL` defaults to `http://localhost:8080`, which is right when the MCP server runs on the same Pi as the BearPrint server.

A quick smoke test over stdio (this lists the tools and runs a health check):

```bash
printf '%s\n' \
  '{"jsonrpc":"2.0","id":1,"method":"initialize"}' \
  '{"jsonrpc":"2.0","id":2,"method":"tools/list"}' \
  '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"bearprint_health","arguments":{}}}' \
  | BEARPRINT_URL=http://your-pi-ip:8080 python3 bearprint-mcp/server.py
```

Run the tests (no network needed):

```bash
cd bearprint-mcp && python3 -m unittest -v
```

## Using it from ChatGPT with an OpenAI Secure MCP Tunnel

A [Secure MCP Tunnel](https://developers.openai.com/api/docs/guides/secure-mcp-tunnels) lets ChatGPT reach a private MCP server without opening an inbound port: the tunnel client on your Pi polls OpenAI outbound and starts the stdio server locally. This section is a summary of a setup that worked end to end on a Raspberry Pi 5 (linux/arm64); the official guide is the source of truth.

1. **Create the tunnel** in the OpenAI Platform and note its tunnel ID and the organization ID that owns it.
2. **Create a runtime API key** in that same organization with Tunnels **Read** and **Use**. A key that can read the tunnel record but lacks Use can show up as polling failures.
3. **Install `tunnel-client`** from the [releases page](https://github.com/openai/tunnel-client/releases) (pick the build for your OS and architecture, and check it against the published SHA256 sums).
4. **Store the key** in an env file that contains only that key, readable only by you:

   ```bash
   install -m 600 /dev/null ~/.config/bearprint-mcp.env
   printf 'CONTROL_PLANE_API_KEY=%s\n' "$YOUR_KEY" > ~/.config/bearprint-mcp.env
   ```

5. **Create the profile** from [`examples/tunnel-client-profile.yaml`](examples/tunnel-client-profile.yaml): fill in the tunnel ID, organization ID and the path to `server.py`, and save it as `~/.config/tunnel-client/bearprint.yaml`.
6. **Run the documented check:** `tunnel-client doctor --profile bearprint --explain`. One failing check is expected while the service is already running: the health listener cannot bind a port that is in use.
7. **Run it as a service** with [`examples/bearprint-mcp-tunnel.service`](examples/bearprint-mcp-tunnel.service) as a systemd user unit, and run `loginctl enable-linger $USER` so it starts at boot.
8. **Check it is polling:** `curl -s http://127.0.0.1:8790/readyz` should say ready, and `commands_poll_last_successful_timestamp_seconds` on `/metrics` should be non-zero and recent.
9. **Add it in ChatGPT** (Plugins, add a custom MCP server, connection type *Tunnel*, select the tunnel). Authentication can be *No authentication*, because access is controlled by who can use the tunnel. Then call `bearprint_health` first.

### Gotchas

- **One key per env file.** `tunnel-client` prefers a `CONTROL_PLANE_API_KEY` environment variable over the profile's `api_key: "env:..."` reference. If the service's environment already contains a key for another organization, that key is sent no matter what the profile names, and you get `mismatched_organization` (with `organization_id` set) or `tunnel_active_organization_required` (without it).
- **`organization_id` is required** in the profile's `control_plane` section for tunnels that need an active organization context, and it must be the organization the key belongs to.
- **The service's logs are in the system journal.** For a user unit, `journalctl --user -u <unit>` may show nothing; use `journalctl _SYSTEMD_USER_UNIT=<unit> -o cat` instead.
- **Role changes take time to propagate.** If you just granted Tunnels access, give it a while before assuming the key is wrong.

## Files

| File | Purpose |
|------|---------|
| `server.py` | The MCP server. |
| `test_server.py` | Unit tests (stdlib `unittest`). |
| `examples/tunnel-client-profile.yaml` | `tunnel-client` profile template with placeholders. |
| `examples/bearprint-mcp-tunnel.service` | systemd user unit template. |
