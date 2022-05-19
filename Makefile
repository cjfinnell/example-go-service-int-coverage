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
