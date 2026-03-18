# Plan: Update AGENTS.md

## Goal

Rewrite AGENTS.md to accurately reflect the sync-and-transcode-music-files project — its structure, commands, conventions, and files — replacing all content inherited from the Tiny Timer project.

## Steps

1. Replace the Project Overview section to describe this project: a CLI tool that syncs music files from a source directory to a destination directory, transcoding non-MP3 formats (`.aif`, `.wav`, `.m4a`) to MP3 using ffmpeg via `goffmpeg`, and removing duplicates.

1. Replace the Essential Commands section with commands that apply to this project:
   - Version control (`jj`) — same workflow, keep as-is
   - Testing: `go test ./...` or `make test`
   - Running: `go run . --source=<dir> --destination=<dir> [--dry-run]`
   - No build/install targets exist in the current Makefile; omit or note absence

1. Replace the Code Organization section to reflect actual files:
   - `main.go` — entry point, CLI flags (`--source`, `--destination`, `--dry-run`)
   - `find_and_transcode_files.go` — directory comparison, transcoding, file copy
   - `remove_duplicate_files.go` — duplicate detection and deletion
   - `helper.go` — shared utilities: `isUntranscodedMusicFile`, `removeNonASCII`, `stringInSlice`
   - `*_test.go` files — one per source file, using `testify/assert`

1. Replace the Architecture section with a description of the main data flow:
   1. `compareDirectories` finds files in source not yet in destination (by transcoded filename)
   1. Each missing file is either transcoded (`.aif`, `.wav`, `.m4a` → `.mp3`) or copied (`.mp3`)
   1. `removeDuplicateFiles` scans destination, groups by base path, keeps the best MP3, deletes the rest

1. Replace Naming Conventions to match this codebase:
   - Types: PascalCase (`fileToTranscode`)
   - Functions: camelCase, action-oriented (`findAndTranscodeFiles`, `removeDuplicateFiles`, `selectPreferredFile`)
   - Constants: camelCase (`maxASCIIIndex`)
   - Error handling: `if err := ...; err != nil { return err }` pattern
   - Defer for cleanup: `defer file.Close()`

1. Replace the Testing section:
   - One `*_test.go` file per source file
   - Uses `testify/assert`
   - Tests use `os.MkdirTemp` for isolated temp directories
   - Tests clean up with `defer os.RemoveAll(tempDir)`
   - All tests must pass: `make test`

1. Replace the Code Quality Checklist:
   - [ ] All tests pass: `make test`
   - [ ] New features have corresponding tests
   - [ ] Non-MP3 music files handled: `.aif`, `.wav`, `.m4a`
   - [ ] Hidden files (prefixed `._`) are skipped
   - [ ] Non-ASCII characters in filenames are sanitized via `removeNonASCII`
   - [ ] Error handling present; errors printed to stderr, not swallowed
   - [ ] `--dry-run` flag honored (no deletions in dry-run mode)
   - [ ] Changes follow Conventional Commits format

1. Update the Documentation section to include only README.md (JOURNAL.md does not exist in this project).

## Files to Change

- `AGENTS.md` — full rewrite, same section structure, updated content
