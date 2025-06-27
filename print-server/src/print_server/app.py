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

    @app.route("/v1/print", methods=["POST"])
    def print_route():
        job = request.get_json()

        if not isinstance(job, dict) or not isinstance(job.get("items"), list):
            return jsonify({"error": "Expected JSON with an 'items' list"}), 400

        try:
            printer = Printer()
            printer.print_job(job["items"])
            return jsonify({"status": "printed"})
        except Exception as e:
            return jsonify({"error": str(e)}), 500

    return app