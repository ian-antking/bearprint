import json
import os
import sys
import unittest

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import server  # noqa: E402


class ValidateItemsTest(unittest.TestCase):
    def test_accepts_every_item_type(self):
        items = [
            {"type": "text", "content": "hi", "align": "center"},
            {"type": "qrcode", "content": "https://example.com"},
            {"type": "blank", "count": 3},
            {"type": "line"},
            {"type": "cut"},
        ]
        self.assertEqual(server.validate_items(items), items)

    def test_rejects_bad_input(self):
        bad = [
            None,
            [],
            ["text"],
            [{"content": "no type"}],
            [{"type": "image"}],
            [{"type": "text"}],
            [{"type": "text", "content": "x", "align": "middle"}],
            [{"type": "text", "content": "x" * (server.MAX_CONTENT + 1)}],
            [{"type": "cut", "content": "unexpected"}],
            [{"type": "blank", "count": 0}],
            [{"type": "blank", "count": 11}],
            [{"type": "blank", "count": True}],
            [{"type": "cut"}] * (server.MAX_ITEMS + 1),
        ]
        for items in bad:
            with self.subTest(items=str(items)[:60]):
                with self.assertRaises(ValueError):
                    server.validate_items(items)


class CallToolTest(unittest.TestCase):
    def setUp(self):
        self._request = server.http_request
        self.requests = []

        def fake(method, path, body=None):
            self.requests.append((method, path, body))
            return 200, "ok"

        server.http_request = fake

    def tearDown(self):
        server.http_request = self._request

    def test_print_defaults_to_dry_run(self):
        out = server.call_tool("bearprint_print", {"items": [{"type": "cut"}]})
        self.assertTrue(out["dry_run"])
        self.assertFalse(out["printed"])
        self.assertEqual(self.requests, [])

    def test_only_an_explicit_false_prints(self):
        for value in (None, True, "false", 0):
            server.call_tool("bearprint_print", {"items": [{"type": "cut"}], "dry_run": value})
        self.assertEqual(self.requests, [])
        out = server.call_tool("bearprint_print", {"items": [{"type": "cut"}], "dry_run": False})
        self.assertTrue(out["printed"])
        self.assertEqual(self.requests, [("POST", "/api/v1/print", {"items": [{"type": "cut"}]})])

    def test_health(self):
        out = server.call_tool("bearprint_health", {})
        self.assertTrue(out["ok"])
        self.assertEqual(self.requests, [("GET", "/api/v1/health", None)])

    def test_unknown_tool(self):
        with self.assertRaises(ValueError):
            server.call_tool("nope", {})


class HandleTest(unittest.TestCase):
    def test_initialize_and_list_tools(self):
        init = server.handle({"jsonrpc": "2.0", "id": 1, "method": "initialize"})
        self.assertEqual(init["result"]["serverInfo"]["name"], "bearprint")
        listed = server.handle({"jsonrpc": "2.0", "id": 2, "method": "tools/list"})
        names = [t["name"] for t in listed["result"]["tools"]]
        self.assertEqual(names, ["bearprint_health", "bearprint_print"])

    def test_notifications_get_no_reply(self):
        self.assertIsNone(server.handle({"jsonrpc": "2.0", "method": "notifications/initialized"}))

    def test_validation_error_is_a_tool_error_not_a_crash(self):
        msg = {"jsonrpc": "2.0", "id": 3, "method": "tools/call",
               "params": {"name": "bearprint_print", "arguments": {"items": []}}}
        result = server.handle(msg)["result"]
        self.assertTrue(result["isError"])

    def test_unknown_method(self):
        reply = server.handle({"jsonrpc": "2.0", "id": 4, "method": "bogus"})
        self.assertEqual(reply["error"]["code"], -32601)


if __name__ == "__main__":
    unittest.main()
