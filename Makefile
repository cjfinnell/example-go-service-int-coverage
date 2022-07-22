.PHONY: deps
deps:
	@go mod tidy
	@go mod vendor

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
