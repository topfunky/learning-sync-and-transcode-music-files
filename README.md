# Learning: Sync and transcode music files

My car's built-in computer only plays mp3 files. This project looks at source and destination files, then transcodes missing music files to mp3 so they can be played from the car's USB stick.

## Usage

Coming soon.

```bash
make build

./sync-and-transcode-music-files -source `/Volumes/A/music-source/` -destination `/Volumes/B/music-mp3/`
```

Use `-dry-run` to preview all changes (transcode, copy, delete) without writing or deleting anything:

```bash
./sync-and-transcode-music-files -source `/Volumes/A/music-source/` -destination `/Volumes/B/music-mp3/` -dry-run
```

## Performance

Transcoding runs in parallel across all CPU cores (`runtime.NumCPU()` workers); MP3 copies are sequential. Progress is shown as `[n/total]` per file. On a multi-core machine, expect a wall-clock speedup of roughly `min(NumCPU, fileCount)`× for transcode-heavy libraries (measured 3.6× faster on 20 files).

## Tests

![Go Tests](https://github.com/topfunky/learning-sync-and-transcode-music-files/actions/workflows/go.yml/badge.svg)

```
go test
```
