GRPC_SOURCE_DIR = schema/grpc/
GRPC_OUTPUT_DIR = ../../interfaces/grpc/
k8s_version = v1.24


setup:
	asdf install
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo 'All runtimes installed'
	@echo ''
	ctlptl create cluster kind --registry=ctlptl-registry --kubernetes-version=${k8s_version}

teardown: clean-setup clean-docker

clean-setup:
	ctlptl delete cluster kind
	@echo ''
	docker rm -f ctlptl-registry
	@echo 'ctlptl-registry deleted'
	@echo ''

clean-docker:
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

