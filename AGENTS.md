# Agent Guidelines

## Project Overview

**sync-and-transcode-music-files** is a CLI tool written in Go that syncs music from a source directory to a destination directory. It:
- Compares source and destination to find files not yet transcoded
- Transcodes `.aif`, `.wav`, and `.m4a` files to `.mp3` using ffmpeg (via `goffmpeg`)
- Copies existing `.mp3` files as-is
- Removes duplicate files from the destination, keeping the best-quality MP3
- Supports `--dry-run` to preview deletions without removing anything

## Essential Commands

### Version control

Generate a summary and commit to VCS with `jj`:

```bash
jj status
jj diff
jj commit -m "feat: transcode wav files"
```

Never run `git` for any reason. Only use `jj`.

Trust that the commit was effective; do not run any follow up commands unless the commit fails.

### Testing
```bash
make test       # Run all tests (go test ./...)
go test -v      # Run with verbose output
```

### Development
```bash
go run . --source=<dir> --destination=<dir>
go run . --source=<dir> --destination=<dir> --dry-run
```

## Code Organization

The codebase uses multiple files organized by responsibility:
- **`main.go`**: Entry point; parses CLI flags (`--source`, `--destination`, `--dry-run`) and calls top-level functions
- **`find_and_transcode_files.go`**: Directory comparison, transcoding via ffmpeg, and MP3 file copying
- **`remove_duplicate_files.go`**: Duplicate detection and deletion; prefers MP3 over other formats, largest file as quality proxy
- **`helper.go`**: Shared utilities — `isUntranscodedMusicFile`, `removeNonASCII`, `stringInSlice`
- **`*_test.go`**: One test file per source file; uses `testify/assert`

## Architecture & Data Flow

1. `compareDirectories` walks source and destination, returning files present in source but not yet in destination (matched by transcoded filename)
1. Each missing file is either:
   - Transcoded to `.mp3` if it is `.aif`, `.wav`, or `.m4a`
   - Copied directly if it is already `.mp3`
1. `removeDuplicateFiles` scans the destination, groups files by base path (path without extension), and for each group keeps the best MP3 (largest file wins as a bit-rate proxy), deleting the rest

### Key Types

```go
type fileToTranscode struct {
    sourcePath      string  // relative path in source directory
    destinationPath string  // relative path in destination (with .mp3 extension)
}
```

## Code Patterns & Conventions

### Markdown

- Write all ordered lists starting with `1.`

### File Structure
- Multiple `.go` files organized by responsibility
- All files are in the same package (`main`)
- Each file has a clear, focused purpose
- Related functionality grouped together logically

### Naming Conventions
- **Types**: PascalCase (`fileToTranscode`)
- **Constants**: camelCase (`maxASCIIIndex`)
- **Functions**: camelCase, action-oriented (`findAndTranscodeFiles`, `removeDuplicateFiles`, `selectPreferredFile`)

### Go Idioms Used
- Error handling: `if err := ...; err != nil { return err }` pattern
- Errors printed to stderr: `fmt.Fprintf(os.Stderr, "❗️ Error: %v\n", err)`
- Defer for cleanup: `defer file.Close()`
- Return error as last value: standard Go convention
- Directory creation: `os.MkdirAll(path, 0755)` before writing files

## Testing Approach

### Test Patterns
- One `*_test.go` file per source file, named to match (e.g., `helper_test.go`)
- Uses `testify/assert` for readable assertions
- Tests create isolated temp directories: `os.MkdirTemp("", "test")`
- Tests clean up with `defer os.RemoveAll(tempDir)`
- All tests must pass before considering work complete: `make test`

### Important: Test Always for User-Facing Changes
- Any new file-handling logic must include tests with temp directories
- Duplicate detection/deletion logic must be tested with `--dry-run` behavior verified
- All tests must pass: `make test`

## Code Quality & Review Checklist

Before submitting changes:
- [ ] All tests pass: `make test`
- [ ] Changes follow Conventional Commits format
- [ ] New features have corresponding tests
- [ ] Non-MP3 music files handled: `.aif`, `.wav`, `.m4a`
- [ ] Hidden files (prefixed `._`) are skipped
- [ ] Non-ASCII characters in filenames are sanitized via `removeNonASCII`
- [ ] Errors printed to stderr, not swallowed silently
- [ ] `--dry-run` flag honored (no deletions when enabled)

## Documentation

- **README.md**: User-facing documentation (installation, usage, features)
- **AGENTS.md** (this file): Developer guidelines for working in the codebase
