BIN := "./bin/calendar"
SCHEDULER_BIN := "./bin/calendar_scheduler"
SENDER_BIN := "./bin/calendar_sender"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

.PHONY: build run build-img run-img version test lint clean

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

install-lint-deps:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh \
	| sh -s -- -b $(shell go env GOPATH)/bin v2.1.6

lint: install-lint-deps
	golangci-lint run ./...

run: build
	$(BIN) --config ./configs/config.yaml

test:
	go test -race -v ./internal/...

GOOSE_MIGRATION_DIR?=./migrations
GOOSE_DRIVER?=postgres
GOOSE_DBSTRING?="postgres://calendar:calendar@localhost:5432/calendar"

export_env:
	@export GOOSE_MIGRATION_DIR GOOSE_DRIVER GOOSE_DBSTRING
	@echo goose settings: "$(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) $(GOOSE_DBSTRING)"

# goose -dir ./migrations/ create add_events_table sql
# goose -dir ./migrations/ postgres postgres://postgres:postgres@localhost:5432/postgres up
clean:
	rm -rf ./bin
	rm -rf ./logs/*
	docker compose down

migrate_new: export_env
	goose -dir $(GOOSE_MIGRATION_DIR) create $(name) go
	
migrate_up: export_env
	goose up

migrate_down: export_env
	goose down && $(MAKE) migrate_status

migrate_status: export_env
	goose status

up:
	docker compose down
	docker compose up -d
	
#grpc-----------------------------------------------------------------------
PROTO_DIR= api
OUT_DIR= api/grpcevents
PROTO_FILE= $(PROTO_DIR)/EventService.proto

generate:
# protoc --go_out=$(GRPC_PATH) --go-grpc_out=$(GRPC_PATH) --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative api/*.proto
# protoc -I $(GRPC_PATH)
	mkdir -p $(OUT_DIR)
	protoc --proto_path=$(PROTO_DIR) \
	--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
	--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
	$(PROTO_FILE)

#rabbitMQ-----------------------------------------------------------------------

build: build-scheduler build-sender
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

build-scheduler:
	go build -o $(SCHEDULER_BIN) ./cmd/scheduler

build-sender:
	go build -o $(SENDER_BIN) ./cmd/sender

run_all:
	docker compose up --build