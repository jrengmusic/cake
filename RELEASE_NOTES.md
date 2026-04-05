First public release.

## What is CAKE?

CMake project manager TUI. One keypress to generate, build, clean, or open IDE.
Menu shows only actions that will succeed based on actual build state.

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

## Requirements

- CMake in PATH
- Terminal 70x24 minimum

## Keyboard

| Key | Action |
|-----|--------|
| g | Generate/Regenerate |
| b | Build |
| c | Clean |
| x | Clean All |
| o | Open IDE |
| / | Preferences |
