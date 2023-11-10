GRPC_SOURCE_DIR = schema/grpc/
GRPC_OUTPUT_DIR = ../../interfaces/grpc/

K8S_VERSION = 1.23

add-plugin:
	- asdf plugin-add golang https://github.com/asdf-community/asdf-golang.git
	- asdf plugin-add golangci-lint https://github.com/hypnoglow/asdf-golangci-lint.git
	- asdf plugin-add protoc https://github.com/paxosglobal/asdf-protoc.git
	- asdf plugin-add ctlptl https://github.com/ezcater/asdf-ctlptl.git
	- asdf plugin-add kind https://github.com/johnlayton/asdf-kind.git
	- asdf plugin-add kubectl https://github.com/asdf-community/asdf-kubectl.git

install-dependency: add-plugin
	asdf install
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

setup: install-dependency
	ctlptl create cluster kind --registry=ctlptl-registry --kubernetes-version=$(K8S_VERSION)

teardown: clean-setup clean-docker

clean-setup:
	ctlptl delete cluster kind
	@echo ''
	docker rm -f ctlptl-registry
	@echo 'ctlptl-registry deleted'
	@echo ''

clean-docker:
	docker image prune -f
	docker volume prune -f
	docker network prune -f
	docker builder prune -f

test:
	go mod tidy
	go fmt ./...
	go vet ./...
	golangci-lint run
	go test -v ./...

gen:
	export PATH="$(PATH):$$(go env GOPATH)/bin" && \
    cd $(GRPC_SOURCE_DIR) && \
    protoc --go_out=$(GRPC_OUTPUT_DIR) \
      --go_opt=paths=source_relative \
      --go-grpc_out=$(GRPC_OUTPUT_DIR) \
      --go-grpc_opt=paths=source_relative \
      *.proto
