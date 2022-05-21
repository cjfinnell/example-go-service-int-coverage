.PHONY: deps
deps:
	@echo "refreshing dependencies..."
	@go mod tidy
	@go mod vendor

.PHONY: mocks
mocks: .bin/mockery
	@echo "generating mocks..."
	@go generate ./...

.PHONY: run
run:
	@docker compose up --build --force-recreate

.PHONY: clean
clean:
	@docker compose rm -sf
	@rm -rf .bin/

.PHONY: lint
lint: .bin/golangci-lint
	@.bin/golangci-lint run

.PHONY: lint-fix
lint-fix: .bin/golangci-lint
	@.bin/golangci-lint run --fix

.bin:
	@mkdir .bin

.bin/golangci-lint: $(wildcard vendor/github.com/golangci/*/*.go) Makefile .bin
	@echo "building linter..."
	@cd vendor/github.com/golangci/golangci-lint/cmd/golangci-lint && go build -o $(shell git rev-parse --show-toplevel)/.bin/golangci-lint .

.bin/mockery: $(wildcard vendor/github.com/vektra/mockery/*/*.go) redis.go Makefile .bin
	@echo "building mock generator..."
	@cd vendor/github.com/vektra/mockery/v2 && go build -o $(shell git rev-parse --show-toplevel)/.bin/mockery .
