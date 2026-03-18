# Copilot Instructions

## Project Overview

This Go project syncs and transcodes music files for use on a car's USB stick. It compares a source directory of music files with a destination directory and transcodes any missing files to `.mp3` format using FFmpeg. Files already in `.mp3` format are copied directly; files in `.aif`, `.wav`, or `.m4a` format are transcoded.

## Repository Structure

- `main.go` — Entry point; parses `-source` and `-destination` CLI flags
- `find_and_transcode_files.go` — Core logic: directory comparison, file copying, and transcoding
- `helper.go` — Utilities: ASCII normalization, file extension checks, string helpers
- `*_test.go` — Unit tests for each source file

## Build & Test

Requires **Go 1.22.2** and **FFmpeg** installed locally.

```bash
# Build
go build -v ./...

# Run all tests
go test

# Run tests with verbose output
go test -v
```

The CI workflow (`.github/workflows/go.yml`) installs FFmpeg automatically via `federicocarboni/setup-ffmpeg@v3.1`.

## Coding Conventions

- All public and package-level functions must have a GoDoc comment.
- Use `fmt.Fprintf(os.Stderr, ...)` for error output; use `fmt.Printf` for progress messages with emoji prefixes (e.g., `🔍`, `🔊`, `📂`, `❗️`).
- Non-ASCII characters in filenames are normalized to ASCII equivalents using `removeNonASCII` in `helper.go`. Add new character mappings there as needed.
- Supported source formats for transcoding: `.aif`, `.wav`, `.m4a`. Add new formats to the `extensions` slice in `isUntranscodedMusicFile` in `helper.go`.
- Hidden files (names starting with `._`) and non-music files (`.DS_Store`, `.txt`, etc.) are silently skipped.
- Use `github.com/stretchr/testify` for test assertions.
- Keep the `fileToTranscode` struct for pairing source and destination paths through the pipeline.

## Dependencies

- [`github.com/xfrr/goffmpeg`](https://github.com/xfrr/goffmpeg) — FFmpeg wrapper for transcoding
- [`github.com/stretchr/testify`](https://github.com/stretchr/testify) — Test assertions
