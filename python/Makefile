all: setup-local

.PHONY: setup-local
setup-local:
	poetry config virtualenvs.in-project true
	poetry install --no-root

.PHONY: run-local
run-local:
	poetry run python -m src

.PHONY: test
test:
	poetry run mypy ./
	@echo ""
	poetry run pytest -v tests/
