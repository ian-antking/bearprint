.PHONY: deploy install run dev

deploy:
	rsync -avz --delete --rsync-path="sudo rsync" ./ $(USER)@$(HOST):/opt/print-server/

install:
	./scripts/install.sh

update: deploy
	ssh $(USER)@$(HOST) 'sudo systemctl restart bearprint.service'

run:
	PYTHONPATH=src gunicorn -c gunicorn.conf.py 'print_server.app:create_app()'

dev:
	poetry run python -m print_server
