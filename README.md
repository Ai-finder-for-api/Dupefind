# dupefind

A small, dependency-free Go CLI that finds duplicate files by content (SHA-256).
Files are first grouped by size so only plausible duplicates get hashed.

## Build & test
    go build -o dupefind .
    go test ./...

## Usage
    ./dupefind [-min bytes] <dir> [dir...]
