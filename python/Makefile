all: setup-local

.PHONY: setup-local
setup-local:
	poetry config virtualenvs.in-project true
	poetry install

.PHONY: run-local
run-local:
	poetry run python -m src
