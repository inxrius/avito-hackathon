.PHONY: fmt test vet race check example

fmt:
	gofmt -w ./cmd ./internal

test:
	go test ./...

vet:
	go vet ./...

race:
	go test -race ./...

check: vet test race

example:
	go run ./cmd/example
