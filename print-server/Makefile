.PHONY: run

run:
	poetry run gunicorn -c gunicorn.conf.py 'print_server.app:create_app()'
