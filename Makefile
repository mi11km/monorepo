GRPC_SOURCE_DIR = schema/grpc/
GRPC_OUTPUT_DIR = ../../interfaces/grpc/

setup:
	asdf install
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

test:
	go mod tidy
	go fmt ./...
	go vet ./...
	golangci-lint run
	go test -v ./...

gen:
	export PATH="$(PATH):$(go env GOPATH)/bin" && \
    cd $(GRPC_SOURCE_DIR) && \
    protoc --go_out=$(GRPC_OUTPUT_DIR) \
      --go_opt=paths=source_relative \
      --go-grpc_out=$(GRPC_OUTPUT_DIR) \
      --go-grpc_opt=paths=source_relative \
      *.proto