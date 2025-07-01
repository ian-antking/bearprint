# ʕ•ᴥ•ʔ BearPrint API

## POST /api/v1/print

Send a JSON object with an `items` array describing the print job. Each item is an object with a `"type"` field and optional other fields depending on the type.

### Supported item types

| Type     | Description                                      | Fields                                                                  | Example                                                     |
|----------|--------------------------------------------------|-------------------------------------------------------------------------|-------------------------------------------------------------|
| `text`   | Prints text, optionally aligned                  | `content` (string), `align` (optional: `"left"`, `"center"`, `"right"`) | `{ "type": "text", "content": "Hello", "align": "center" }` |
| `qrcode` | Prints a QR code                                 | `content` (string)  `align` (optional: `"left"`, `"center"`, `"right"`) | `{ "type": "qrcode", "content": "https://example.com" }`    |
| `blank`  | Prints blank lines                               | `count` (optional, integer, default 1)                                  | `{ "type": "blank", "count": 3 }`                           |
| `line`   | Prints a horizontal line (dashes)                | None                                                                    | `{ "type": "line" }`                                        |
| `cut`    | Cuts the paper                                   | None                                                                    | `{ "type": "cut" }`                                         |

---

### Example request body

```json
{
  "items": [
    { "type": "text", "content": "Hello from BearPrint!", "align": "center" },
    { "type": "qrcode", "content": "https://example.com" },
    { "type": "line" },
    { "type": "blank", "count": 2 },
    { "type": "cut" }
  ]
}
```

---

### Example curl command

```bash
curl -X POST -H "Content-Type: application/json" \
    -d '{
          "items": [
            { "type": "text", "content": "Hello from BearPrint!", "align": "center" },
            { "type": "qrcode", "content": "https://example.com" },
            { "type": "line" },
            { "type": "blank", "count": 2 },
            { "type": "cut" }
          ]
        }' \
    http://localhost:5000/v1/print
Response
200 OK
```

```json
{ "status": "printed" }
```

400 Bad Request

```json
{ "error": "Expected JSON with an 'items' list" }
```

500 Internal Server Error

```json
{ "error": "<error message>" }
```
