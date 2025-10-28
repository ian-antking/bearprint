#!/bin/bash
set -e

API_URL="${API_URL:-http://localhost:8080/api/v1/health}"

echo "Running BearPrint API smoke test against $API_URL..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$API_URL" || echo "000")

if [ "$HTTP_CODE" != "200" ]; then
    echo "❌ Smoke test failed: expected HTTP 200, got $HTTP_CODE"
    exit 1
fi

echo "✅ Smoke test passed: BearPrint API is healthy"