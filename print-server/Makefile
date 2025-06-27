.PHONY: deploy install run dev

deploy:
	rsync -avz --delete ./ $(USER)@$(HOST):~/print-server/

install:
	if ! grep -qxF 'export PATH="$HOME/.local/bin:$$PATH"' ~/.bashrc; then \
		echo 'export PATH="$HOME/.local/bin:$$PATH"' >> ~/.bashrc; \
		echo "Added ~/.local/bin to PATH in ~/.bashrc"; \
	fi
	sudo apt update
	sudo apt install -y python3-pip
	python3 -m pip install --break-system-packages --upgrade pip
	python3 -m pip install --break-system-packages -r requirements.txt

run:
	PYTHONPATH=src gunicorn -c gunicorn.conf.py 'print_server.app:create_app()'

dev:
	poetry run python -m print_server
