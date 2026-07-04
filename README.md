# gobt

A Go learning project for building a BitTorrent client step by step.

## Current Status

Implemented modules:
- `pkg/bencode`: bencode encode/decode
- `pkg/torrent`: torrent metadata model, parse/encode, validation, info-hash
- `pkg/tracker`: HTTP announce request/response handling
- `pkg/utils`: shared helpers and versioning

In progress roadmap:
- DHT basics
- Peer discovery aggregation
- Peer handshake and message protocol
- Piece download pipeline and verification

## Quick Start

```bash
# clone
 git clone https://github.com/faterer/gobt.git
 cd gobt

# test all packages
 go test ./...

# run app entry
 go run ./cmd

# run examples
 cd examples
 go run parse_torrent.go init.go
 # optional tagged examples
 go run -tags bencode_example bencode_simple.go
 go run -tags tracker_example tracker_announce.go <torrent-file> <tracker-url>
```

## Project Layout

```text
gobt/
├── cmd/                # app entry
├── pkg/                # core modules
│   ├── bencode/
│   ├── torrent/
│   ├── tracker/
│   └── utils/
├── examples/           # runnable examples
├── docs/               # active docs
│   └── archive/        # historical docs
├── project_plan.md     # Week 1-12 execution plan
└── go.mod
```

## Documentation

Active docs:
- `project_plan.md`: execution and weekly tracking baseline
- `docs/INDEX.md`: document navigation
- `docs/GETTING_STARTED.md`: setup and first steps
- `docs/ROADMAP.md`: high-level implementation roadmap
- `docs/ARCHITECTURE.md`: architecture notes
- `docs/QUICKREF.md`: protocol and implementation quick reference

Archived docs are kept in `docs/archive/` for historical context.

## Validation Commands

```bash
# full test suite
 go test ./...

# tracker module tests
 go test ./pkg/tracker -v

# coverage profile (optional)
 go test -coverprofile=coverage ./...
```

## License

MIT
