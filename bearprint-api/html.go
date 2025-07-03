package main

import "net/http"

func writeHTMLHeader(w http.ResponseWriter) error {
    _, err := w.Write([]byte(`
        <!DOCTYPE html><html><head>
        <meta charset="UTF-8">
        <title>BearPrint API Docs</title>
        <style>
          body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
            margin: 4rem auto;
            max-width: 800px;
            padding: 0 1rem;
            background-color: #fff;
            color: #222;
          }
          pre {
            background-color: #2d2d2d;
            padding: 1rem;
            overflow-x: auto;
            border-radius: 5px;
          }
          table {
            border-collapse: collapse;
            width: 100%;
            margin-bottom: 1rem;
          }
          th, td {
            border: 1px solid #ddd;
            padding: 0.5rem;
            text-align: left;
          }
          th {
            background-color: #f4f4f4;
          }
        </style>
        </head><body>
    `))
    return err
}

func writeHTMLFooter(w http.ResponseWriter) error {
    _, err := w.Write([]byte(`</body></html>`))
    return err
}
