# Default make

.PHONY: default
default:
	@make tests

# Generate files

.PHONY: gen
gen:
	@go generate ./...

# Lint

.PHONY: lint
lint:
	@golangci-lint run

# Test

TEST_FUNC=^.*$$
ifdef t
TEST_FUNC=$(t)
endif
TEST_PKG=./...
ifdef p
TEST_PKG=./$(p)
endif

.PHONY: test
test:
	@go test -v -timeout 30s -run $(TEST_FUNC) $(TEST_PKG)

.PHONY: tests
tests:
	@go test ./...

# Docs

.PHONY: docs
docs:
	@godoc -http=localhost:9995

