.PHONY: deps
deps:
	@echo "refreshing dependencies..."
	@go mod tidy
	@go mod vendor

.PHONY: mocks
mocks: .bin/mockery
	@echo "generating mocks..."
	@go generate ./...

.PHONY: unit-test
unit-test: .bin/gotestsum
	@.bin/gotestsum -- -cover -coverprofile coverage.txt -short ./...

.PHONY: test
test:
	@docker compose up -d --build --force-recreate redis
	@docker build --target test -t rediswrapper-int-test .
	@docker run -it --network testnet --name rediswrapper-int-test rediswrapper-int-test
	@docker cp rediswrapper-int-test:/test/coverage.txt .
	@docker cp rediswrapper-int-test:/test/results.xml .

.PHONY: _int-test
_int-test: .bin/gotestsum
	@.bin/gotestsum --junitfile results.xml -- -cover -coverprofile coverage.txt ./...

.PHONY: fuzz
fuzz:
	@docker compose up -d --build --force-recreate redis
	@docker build --target fuzz -t rediswrapper-fuzz .
	@docker run -it --network testnet --name rediswrapper-fuzz rediswrapper-fuzz
	@docker cp rediswrapper-int-test:/fuzz/testdata .

.PHONY: load
load:
	@docker compose up -d --build --force-recreate
	-@docker run -i --network testnet --name k6 loadimpact/k6 run --out json=/tmp/stats.json - < load.js
	@docker cp k6:/tmp/stats.json .
	@docker compose down

.PHONY: run
run:
	@docker compose up --build --force-recreate

.PHONY: clean
clean:
	@docker compose rm -sf
	-@docker rm -f rediswrapper-int-test rediswrapper-fuzz k6
	-@docker network rm testnet
	@rm -rf .bin/ coverage.txt results.xml stats.json
	@go clean -testcache

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

.bin/gotestsum: $(wildcard vendor/gotest.bin/*/*.go) Makefile .bin
	@echo "building test runner..."
	@cd vendor/gotest.tools/gotestsum && go build -o $(shell git rev-parse --show-toplevel)/.bin/gotestsum .
