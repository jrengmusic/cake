## What's New in v0.0.1

- Open Editor: press `o` with Ninja selected to open nvim in build directory
- Menu label changes dynamically: "Open IDE" for Xcode/VS, "Open Editor" for Ninja

## Platforms

| OS | Arch | Signed |
|----|------|--------|
| macOS | Intel (x86_64) | Yes (notarized) |
| macOS | Apple Silicon (arm64) | Yes (notarized) |
| Linux | x86_64 | - |
| Linux | arm64 | - |
| Windows | x86_64 | - |
| Windows | arm64 | - |

## Install

Download binary from this release, or:

```bash
go install github.com/jrengmusic/cake/cmd/cake@latest
```

## Keyboard

| Key | Action |
|-----|--------|
| g | Generate/Regenerate |
| b | Build |
| c | Clean |
| x | Clean All |
| o | Open IDE / Editor |
| / | Preferences |
