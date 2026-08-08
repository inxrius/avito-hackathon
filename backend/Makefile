.PHONY: fmt test vet race check run example

fmt:
	gofmt -w ./cmd ./internal ./pkg

test:
	go test ./...

vet:
	go vet ./...

race:
	go test -race ./...

check: vet test

run:
	go run ./cmd/server

example:
	go run ./cmd/example
