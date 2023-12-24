build:
	@go build -o bin/storer cmd/main.go

build-arm:
	@env GOOS=linux GOARCH=arm GOARM=7 go build -o bin/storer cmd/main.go

run: build
	@./bin/storer

test:
	@go test -v ./...