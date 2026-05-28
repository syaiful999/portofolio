GOPATH:=$(shell go env GOPATH)
export PATH:=$(GOPATH)/bin:$(PATH)

.PHONY: init
init:
	@go get -u google.golang.org/protobuf/proto
	@go install github.com/golang/protobuf/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/asim/go-micro/cmd/protoc-gen-micro/v4@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest


.PHONY: proto
PROTOC_GEN_MICRO := $(GOPATH)/bin/protoc-gen-micro

proto:
	@echo "Generating micro services..."
	@protoc --proto_path=. --micro_out=. --go_out=. --plugin=$(PROTOC_GEN_MICRO) pkg/enum/proto/*.proto
	@protoc --proto_path=. --micro_out=. --go_out=. --plugin=$(PROTOC_GEN_MICRO) pkg/user/proto/*.proto


.PHONY: update
update:
	@go get -u

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: build
build:
	@go build -o auth *.go

.PHONY: test
test:
	@go test -v ./... -cover

.PHONY: docker
docker:
	@docker build -t auth:latest .
