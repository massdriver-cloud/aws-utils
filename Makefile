.PHONY: generate
generate:
	go run pkg/generator/gen.go

.PHONY: test
test:
	go test ./...