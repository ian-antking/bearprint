from flask import Flask

def create_app(config: dict = None) -> Flask:
    app = Flask(__name__)

    if config:
        app.config.update(config)

    @app.route("/")
    def index():
        return "Hello from BearPrint!"

    return app