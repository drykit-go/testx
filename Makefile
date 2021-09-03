# Default command

.PHONY: default
default:
	@make check

# Generate files

.PHONY: gen
gen:
	@echo "🛠  Building gen binary"
	@go build -o ./bin/gen ./cmd/gen/main.go
	@echo "✅ Done"
	@go generate ./...

# Check code

.PHONY: check
check:
	@make lint
	@make tests

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: tests
tests:
	@go test ./...

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

# Docs

.PHONY: docs
docs:
	@echo "\033[4mhttp://localhost:9995/pkg/github.com/drykit-go/testx/\033[0m"
	@godoc -http=localhost:9995
