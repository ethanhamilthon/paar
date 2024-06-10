build:
	@go build -o bin/paar cmd/main.go

run:build
	@bin/paar

test:
	@go test ./... 

testv:
	@go test -v ./...

testcover:
	@go test -cover ./...

testrace:
	@go test -race ./...