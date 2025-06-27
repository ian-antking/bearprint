.PHONY: deploy install run dev

deploy:
	rsync -avz --delete ./ $(USER)@$(HOST):~/print-server/

install:
	./scripts/install.sh

run:
	PYTHONPATH=src gunicorn -c gunicorn.conf.py 'print_server.app:create_app()'

dev:
	poetry run python -m print_server
