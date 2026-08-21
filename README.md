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

## Tests

![Go Tests](https://github.com/topfunky/learning-sync-and-transcode-music-files/actions/workflows/go.yml/badge.svg)

```
go test
```
