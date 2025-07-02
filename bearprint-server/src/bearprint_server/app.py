from flask import Flask, request, jsonify, render_template, url_for
import markdown
from pathlib import Path
from .printer import Printer

def create_app(config: dict = None) -> Flask:
    app = Flask(__name__)

    if config:
        app.config.update(config)

    printer = Printer()

    @app.route("/")
    def _():
        md_path = Path(__file__).parent / "static" / "index.md"
        if not md_path.exists():
            return "Documentation not found", 404

        md_content = md_path.read_text(encoding="utf-8")
        html_body = markdown.markdown(
            md_content,
            extensions=['tables', 'fenced_code', 'codehilite'],
            extension_configs={
                'codehilite': {
                    'linenums': False,
                    'guess_lang': True
                }
            }
        )

        return render_template(
            "index.html",
            content=html_body,
        )

    @app.route("/api/v1/print", methods=["POST"])
    def __():
        job = request.get_json()

        if not isinstance(job, dict) or not isinstance(job.get("items"), list):
            return jsonify({"error": "Expected JSON with an 'items' list"}), 400

        try:
            printer.print_job(job["items"])
            return jsonify({"status": "printed"})
        except Exception as e:
            return jsonify({"error": str(e)}), 500

    return app

app = create_app()
