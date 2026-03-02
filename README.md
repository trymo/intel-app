# Intel App Mover (macOS)

A tiny Go utility that scans a macOS Applications folder for Intel (x86_64) apps and optionally moves them into a separate folder. Useful on Apple Silicon to identify Intel-only apps still running via Rosetta, or to organize mixed/universal installs.

## Features
- Detects Mach-O CPU slices (Intel x86_64 and ARM64) directly via the standard library.
- Dry-run by default; requires an explicit flag to move anything.
- Filters Intel-only apps by default, with an option to include any app that contains an Intel slice.
- Simple, fast, no external dependencies.

## Requirements
- macOS (scans `.app` bundles and Mach-O binaries).
- Go `1.20+` to build locally.
- Permission to move apps in the chosen directory (you may be prompted for admin rights when moving from `/Applications`).

## Install

Using Makefile:

```bash
make build
```

This produces the binary at `bin/intel-app-mover`.

Or directly with Go:

```bash
mkdir -p bin
go build -buildvcs=false -o bin/intel-app-mover ./cmd/intel-app-mover
```

## Usage

```text
bin/intel-app-mover [flags]
```

Flags:
- `-dir` Applications directory to scan. Default: `/Applications`
- `-execute` Perform moves (otherwise dry-run only)
- `-v` Verbose logging
- `-dest` Destination folder for moved apps. Default: `<dir>/intel-apps`
- `-intel-only` Move only Intel-only apps (no ARM64 slice). Default: `true`

### Examples
- Dry-run scan of `/Applications` (default):
  ```bash
  bin/intel-app-mover
  ```
- Move Intel-only apps from `/Applications` into `/Applications/intel-apps`:
  ```bash
  bin/intel-app-mover -execute
  ```
- Include any app that contains an Intel slice (even if universal):
  ```bash
  bin/intel-app-mover -execute -intel-only=false
  ```
- Scan a custom directory and choose a custom destination:
  ```bash
  bin/intel-app-mover -dir ~/Applications -dest ~/Applications/intel-archive -execute
  ```

## How It Works
- Reads each top-level `.app` bundle in `-dir` and extracts `CFBundleExecutable` from `Contents/Info.plist`.
- Inspects the corresponding Mach-O binary to detect Intel and/or ARM64 CPU slices.
- Lists matching apps during dry-run; when `-execute` is provided, moves them to `-dest` via `os.Rename`.

## Notes & Limitations
- Scans only the top-level of the specified directory (no recursion).
- Skips bundles missing `Info.plist` or the expected `Contents/MacOS/<executable>`.
- Moving apps out of `/Applications` can impact auto-updaters or launch services; confirm with your environment.
- You can undo a move by moving the `.app` bundle back to its original location.

## Development
- `make build` — compile to `bin/intel-app-mover`
- `make fmt` — format sources
- `make clean` — remove build artifacts

## License
GPL-3.0 — see `LICENSE` for full text.

