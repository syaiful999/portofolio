GOPATH:=$(shell go env GOPATH)
export PATH:=$(GOPATH)/bin:$(PATH)

DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= your_username
DB_PASSWORD ?= your_password
DB_NAME ?= master_db
DB_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATIONS_PATH = database/migrations

.PHONY: init
init:
	@go get -u google.golang.org/protobuf/proto
	@go install github.com/golang/protobuf/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/asim/go-micro/cmd/protoc-gen-micro/v4@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest


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
	@go build -o master-service *.go

.PHONY: test
test:
	@go test -v ./... -cover

.PHONY: docker
docker:
	@docker build -t master-service:latest .

# ==================== Database Migrations ====================

.PHONY: migrate-up
migrate-up:
	@echo "Running all migrations..."
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

.PHONY: migrate-down
migrate-down:
	@echo "Rolling back last migration..."
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1

.PHONY: migrate-down-all
migrate-down-all:
	@echo "Rolling back all migrations..."
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down

.PHONY: migrate-reset
migrate-reset:
	@echo "Resetting database (down all + up all)..."
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" drop -f
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

.PHONY: migrate-version
migrate-version:
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" version

.PHONY: migrate-force
migrate-force:
	@echo "Forcing migration version to $(VERSION)..."
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" force $(VERSION)

.PHONY: migrate-create
migrate-create:
	@echo "Creating migration: $(NAME)..."
	@migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(NAME)

.PHONY: db-create
db-create:
	@echo "Creating database $(DB_NAME)..."
	@PGPASSWORD=$(DB_PASSWORD) createdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME)

.PHONY: db-drop
db-drop:
	@echo "Dropping database $(DB_NAME)..."
	@PGPASSWORD=$(DB_PASSWORD) dropdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME)

.PHONY: db-setup
db-setup: db-create migrate-up
	@echo "Database setup complete!"
