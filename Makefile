all: setup-local

.PHONY: add-plugins
add-plugins:
	- asdf plugin-add golang https://github.com/asdf-community/asdf-golang.git
	- asdf plugin-add python https://github.com/danhper/asdf-python.git
	- asdf plugin-add poetry https://github.com/asdf-community/asdf-poetry.git

.PHONY: setup-local
setup-local: add-plugins
	asdf install
