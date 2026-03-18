.PHONY: test build validate

test:
	go test ./...

validate:
	go vet ./...
	staticcheck ./...

build: test
	go build -o ./sync-and-transcode .
