from flask import Flask, request, jsonify
from .printer import Printer

def create_app(config: dict = None) -> Flask:
    app = Flask(__name__)

    if config:
        app.config.update(config)

    printer = Printer()

    @app.route("/")
    def index():
        return "Hello from BearPrint!"

    @app.route("/print/text", methods=["POST"])
    def print_route():
        data = request.json
        text = data.get("text", "")
        if not text:
            return jsonify({"error": "No text provided"}), 400
        printer.print_text(text)
        return jsonify({"status": "printed"})

    return app