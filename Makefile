.PHONY: tests
tests:
	@go test ./...

.PHONY: gen
gen:
	@go generate ./...
