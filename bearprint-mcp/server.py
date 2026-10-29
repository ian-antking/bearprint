#!/usr/bin/env python3
"""Minimal MCP server (stdio, JSON-RPC 2.0) wrapping the BearPrint HTTP API.

Designed to run beside `tunnel-client --mcp.command "python3 server.py"`, so
nothing listens on a network port. Default is dry-run: print_job only returns
the validated job unless dry_run is explicitly false.
"""
import json
import os
import sys
import urllib.error
import urllib.request

BEARPRINT_URL = os.environ.get("BEARPRINT_URL", "http://localhost:8080")
PROTOCOL_VERSION = "2025-06-18"
TIMEOUT = 10

ITEM_SCHEMAS = {
    "text": {"required": ["content"], "optional": ["align"]},
    "qrcode": {"required": ["content"], "optional": ["align"]},
    "blank": {"required": [], "optional": ["count"]},
    "line": {"required": [], "optional": []},
    "cut": {"required": [], "optional": []},
}
ALIGNS = {"left", "center", "right"}
MAX_ITEMS = 50
MAX_CONTENT = 500


def validate_items(items):
    if not isinstance(items, list) or not items:
        raise ValueError("items must be a non-empty list")
    if len(items) > MAX_ITEMS:
        raise ValueError(f"items may contain at most {MAX_ITEMS} entries")
    for i, item in enumerate(items):
        if not isinstance(item, dict) or "type" not in item:
            raise ValueError(f"item {i} must be an object with a type")
        kind = item["type"]
        if kind not in ITEM_SCHEMAS:
            raise ValueError(f"item {i} has unknown type {kind!r}")
        allowed = set(["type"] + ITEM_SCHEMAS[kind]["required"] + ITEM_SCHEMAS[kind]["optional"])
        extra = set(item) - allowed
        if extra:
            raise ValueError(f"item {i} has unexpected fields {sorted(extra)}")
        for field in ITEM_SCHEMAS[kind]["required"]:
            if field not in item:
                raise ValueError(f"item {i} ({kind}) needs {field}")
        if "content" in item:
            if not isinstance(item["content"], str) or len(item["content"]) > MAX_CONTENT:
                raise ValueError(f"item {i} content must be a string up to {MAX_CONTENT} chars")
        if "align" in item and item["align"] not in ALIGNS:
            raise ValueError(f"item {i} align must be left, center or right")
        if "count" in item:
            c = item["count"]
            if not isinstance(c, int) or isinstance(c, bool) or not 1 <= c <= 10:
                raise ValueError(f"item {i} count must be an integer 1-10")
    return items


def http_request(method, path, body=None):
    data = json.dumps(body).encode() if body is not None else None
    req = urllib.request.Request(
        BEARPRINT_URL + path,
        data=data,
        method=method,
        headers={"Content-Type": "application/json"},
    )
    try:
        with urllib.request.urlopen(req, timeout=TIMEOUT) as resp:
            return resp.status, resp.read().decode("utf-8", "replace")
    except urllib.error.HTTPError as e:
        return e.code, e.read().decode("utf-8", "replace")


TOOLS = [
    {
        "name": "bearprint_health",
        "description": "Check that the BearPrint printer server is running.",
        "inputSchema": {"type": "object", "properties": {}, "additionalProperties": False},
    },
    {
        "name": "bearprint_print",
        "description": (
            "Print a job on the BearPrint thermal printer. Items are printed in order. "
            "Defaults to dry_run=true, which only validates and returns the job. "
            "Set dry_run to false to actually print."
        ),
        "inputSchema": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "minItems": 1,
                    "maxItems": MAX_ITEMS,
                    "items": {
                        "type": "object",
                        "properties": {
                            "type": {"enum": list(ITEM_SCHEMAS)},
                            "content": {"type": "string", "maxLength": MAX_CONTENT},
                            "align": {"enum": sorted(ALIGNS)},
                            "count": {"type": "integer", "minimum": 1, "maximum": 10},
                        },
                        "required": ["type"],
                        "additionalProperties": False,
                    },
                },
                "dry_run": {"type": "boolean", "default": True},
            },
            "required": ["items"],
            "additionalProperties": False,
        },
    },
]


def call_tool(name, args):
    if name == "bearprint_health":
        status, body = http_request("GET", "/api/v1/health")
        return {"ok": status == 200, "status": status, "body": body}
    if name == "bearprint_print":
        items = validate_items(args.get("items"))
        dry_run = args.get("dry_run", True) is not False
        job = {"items": items}
        if dry_run:
            return {"dry_run": True, "job": job, "printed": False}
        status, body = http_request("POST", "/api/v1/print", job)
        return {"dry_run": False, "printed": status == 200, "status": status, "body": body}
    raise ValueError(f"unknown tool {name!r}")


def handle(msg):
    method = msg.get("method")
    mid = msg.get("id")
    if mid is None:  # notification, no reply
        return None
    try:
        if method == "initialize":
            result = {
                "protocolVersion": PROTOCOL_VERSION,
                "capabilities": {"tools": {}},
                "serverInfo": {"name": "bearprint", "version": "0.1.0"},
            }
        elif method == "tools/list":
            result = {"tools": TOOLS}
        elif method == "tools/call":
            params = msg.get("params", {})
            try:
                out = call_tool(params.get("name"), params.get("arguments") or {})
                result = {"content": [{"type": "text", "text": json.dumps(out)}], "isError": False}
            except ValueError as e:
                result = {"content": [{"type": "text", "text": str(e)}], "isError": True}
            except OSError as e:
                result = {"content": [{"type": "text", "text": f"BearPrint unreachable: {e}"}], "isError": True}
        elif method == "ping":
            result = {}
        else:
            return {"jsonrpc": "2.0", "id": mid,
                    "error": {"code": -32601, "message": f"method not found: {method}"}}
        return {"jsonrpc": "2.0", "id": mid, "result": result}
    except Exception as e:  # keep the server alive on bad input
        return {"jsonrpc": "2.0", "id": mid, "error": {"code": -32603, "message": str(e)}}


def main():
    for line in sys.stdin:
        line = line.strip()
        if not line:
            continue
        try:
            msg = json.loads(line)
        except json.JSONDecodeError:
            continue
        reply = handle(msg)
        if reply is not None:
            sys.stdout.write(json.dumps(reply) + "\n")
            sys.stdout.flush()


if __name__ == "__main__":
    main()
