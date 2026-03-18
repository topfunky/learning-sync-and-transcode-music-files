# Plan: Add `make build` and `make validate` Targets

## Goal

Extend the `Makefile` with two new targets: `build` and `validate`.

## Tasks

1. Add `make validate` target
   - Runs `go vet ./...`
   - Runs `staticcheck ./...` (install if not present: `go install honnef.co/go/tools/cmd/staticcheck@latest`)
   - Listed in `.PHONY`

2. Add `make build` target
   - Depends on `test` (tests must pass before building)
   - Runs `go build -o ./sync-and-transcode .`
   - Outputs binary `./sync-and-transcode` in project root
   - Listed in `.PHONY`

3. Add `./sync-and-transcode` binary to `.gitignore`

## Final Makefile Shape

```makefile
.PHONY: test build validate

test:
	go test ./...

validate:
	go vet ./...
	staticcheck ./...

build: test
	go build -o ./sync-and-transcode .
```

## Notes

- `staticcheck` must be installed separately; the plan assumes it is available on `PATH`.
- `make build` will fail fast if any test fails, preventing a broken binary.
- The binary name `sync-and-transcode` matches the project's primary purpose.
