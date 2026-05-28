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
PROTOC_GEN_GRPC := $(GOPATH)/bin/protoc-gen-go-grpc
PROTOC_GEN_GATEWAY := $(GOPATH)/bin/protoc-gen-grpc-gateway

proto:
	@echo "Generating micro services..."
# 	@protoc --proto_path=. --micro_out=. --go_out=.  --plugin=$(PROTOC_GEN_MICRO) pkg/auditlog/proto/*.proto
# 	@protoc --proto_path=. --micro_out=. --go_out=.  --plugin=$(PROTOC_GEN_MICRO) pkg/file/proto/*.proto
# 	@protoc --proto_path=. --micro_out=. --go_out=.  --plugin=$(PROTOC_GEN_MICRO) pkg/shift/proto/*.proto
# 	@protoc --proto_path=. --micro_out=. --go_out=.  --plugin=$(PROTOC_GEN_MICRO) pkg/shifttemplate/proto/*.proto
# 	@protoc --proto_path=. --micro_out=. --go_out=.  --plugin=$(PROTOC_GEN_MICRO) pkg/user/proto/*.proto
# 	@protoc --proto_path=. --micro_out=. --go_out=.  --plugin=$(PROTOC_GEN_MICRO) pkg/enum/proto/*.proto
# 	@protoc --proto_path=. --micro_out=. --go_out=.  --plugin=$(PROTOC_GEN_MICRO) pkg/reader/proto/*.proto
	@protoc --proto_path=. --micro_out=. --go_out=.  --plugin=$(PROTOC_GEN_MICRO) pkg/outsource/proto/*.proto
# 	@protoc --proto_path=. --micro_out=. --go_out=.  --plugin=$(PROTOC_GEN_MICRO) pkg/rfid/proto/*.proto
# 	@protoc --proto_path=. --micro_out=. --go_out=.  --plugin=$(PROTOC_GEN_MICRO) pkg/employee/proto/*.proto
# 	@protoc --proto_path=. --micro_out=. --go_out=.  --plugin=$(PROTOC_GEN_MICRO) pkg/city/proto/*.proto
# 	@protoc --proto_path=. --micro_out=. --go_out=.  --plugin=$(PROTOC_GEN_MICRO) pkg/province/proto/*.proto
# 	@protoc --proto_path=. --micro_out=. --go_out=.  --plugin=$(PROTOC_GEN_MICRO) pkg/kecamatan/proto/*.proto
# 	@protoc --proto_path=. --micro_out=. --go_out=.  --plugin=$(PROTOC_GEN_MICRO) pkg/kelurahan/proto/*.proto
# 	@protoc --proto_path=. --micro_out=. --go_out=.  --plugin=$(PROTOC_GEN_MICRO) pkg/logging/proto/*.proto
# 	@protoc --proto_path=. --micro_out=. --go_out=.  --plugin=$(PROTOC_GEN_MICRO) pkg/departmentdetail/proto/*.proto
# 	@protoc --proto_path=. --micro_out=. --go_out=.  --plugin=$(PROTOC_GEN_MICRO) pkg/section/proto/*.proto
# 	@protoc --proto_path=. --micro_out=. --go_out=:. --plugin=$(PROTOC_GEN_MICRO) pkg/jobdescription/proto/*.proto
# 	@protoc --proto_path=. --micro_out=. --go_out=:. --plugin=$(PROTOC_GEN_MICRO) pkg/jobfamily/proto/*.proto
# 	@protoc --proto_path=. --micro_out=. --go_out=:. --plugin=$(PROTOC_GEN_MICRO) pkg/jobcompetency/proto/*.proto
# 	@protoc --proto_path=. --micro_out=. --go_out=:. --plugin=$(PROTOC_GEN_MICRO) pkg/mail/proto/*.proto
# 	@protoc --proto_path=. --micro_out=. --go_out=:. --plugin=$(PROTOC_GEN_MICRO) pkg/role/proto/*.proto
# 	@protoc --proto_path=. --micro_out=. --go_out=:. --plugin=$(PROTOC_GEN_MICRO) pkg/menu/proto/*.proto


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
