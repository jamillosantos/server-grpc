.PHONY: lint
lint:
	go tool golangci-lint run

.PHONY: lint-fix
lint-fix:
	go tool golangci-lint run --fix

.PHONY: test
test:
	go test -v ./...