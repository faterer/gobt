# Project Plan (Week 1-12)

This document is the long-term execution guide for the gobt project.

## Goal

Build a practical BitTorrent client that can:
- Parse torrent metadata
- Communicate with trackers and DHT
- Discover peers
- Connect to peers and exchange protocol messages
- Download data pieces with verification
- Support resume and stable long-running sessions

## Progress Rule

- Update this file at the end of each week.
- Mark week status as: `Not Started`, `In Progress`, or `Done`.
- Add short notes for blockers and next actions.

## Weekly Plan

### Week 1: Project Foundation
Goal:
- Build project skeleton, module layout, and base testing workflow.

Deliverables:
- CLI entrypoint and minimal runtime flow
- Shared error handling and logging baseline
- Test command available in local environment

Acceptance:
- Project builds and runs successfully
- Baseline tests can be executed end-to-end

### Week 2: Bencode Encode/Decode
Goal:
- Implement bencode support for string, int, list, and dict.

Deliverables:
- Encoder and decoder
- Boundary/error handling tests

Acceptance:
- Round-trip encode/decode works for nested data
- Core bencode tests pass reliably

### Week 3: Torrent Model and Parser
Goal:
- Parse single-file and multi-file torrents with validation and info hash.

Deliverables:
- TorrentInfo model
- Torrent parser and validator
- Demo parsing utilities

Acceptance:
- Multiple sample torrents parse correctly
- Info hash generation is deterministic and valid

### Week 4: Tracker Communication
Goal:
- Build announce request and parse tracker responses.

Deliverables:
- Tracker client (HTTP)
- Request URL builder and response parser
- Tracker package tests

Acceptance:
- Can request peers from tracker
- Tracker tests pass

Status:
- Done

### Week 5: DHT Basics
Goal:
- Implement minimum DHT workflow for peer lookup.

Deliverables:
- Node ID and basic routing table behavior
- Core find_node/get_peers request flow

Acceptance:
- Can query bootstrap nodes and receive valid candidate nodes/peers

### Week 6: Peer Discovery Aggregation
Goal:
- Combine peer sources from Tracker and DHT.

Deliverables:
- Discovery manager
- Peer deduplication and health scoring

Acceptance:
- Discovery returns stable and deduplicated peer set

### Week 7: Peer Connection and Handshake
Goal:
- Connect to peers and complete BitTorrent handshake.

Deliverables:
- TCP connection lifecycle management
- Handshake marshal/unmarshal and timeout handling

Acceptance:
- Handshake succeeds with real peers in controlled tests

### Week 8: Message Protocol and State Machine
Goal:
- Implement core peer-wire messages and state transitions.

Deliverables:
- Message codec framework
- Handling for keep-alive, choke/unchoke, interested/not interested, have, bitfield

Acceptance:
- Message exchange is correct and peer state remains consistent

### Week 9: Piece Download Pipeline
Goal:
- Implement piece/block requesting, receiving, and retry logic.

Deliverables:
- Download scheduler
- Request queue, timeout retry, and reassembly flow

Acceptance:
- Small file can be downloaded repeatedly with stable success

### Week 10: Storage and Hash Verification
Goal:
- Persist data to disk and verify integrity.

Deliverables:
- File manager and piece writer
- SHA1 piece verifier and bad-piece re-request

Acceptance:
- Completed download passes hash checks
- Corrupted pieces are detected and re-downloaded

### Week 11: Resume and Stability
Goal:
- Support session recovery and long-running stability.

Deliverables:
- Resume metadata persistence
- Restart recovery flow
- Better resource cleanup and error recovery

Acceptance:
- Interrupted downloads can resume correctly
- Long-run tests show no major leaks/crashes

### Week 12: Performance and Release
Goal:
- Optimize throughput and finalize release.

Deliverables:
- Performance tuning (concurrency, request window, connection strategy)
- Release notes and user docs

Acceptance:
- Meets practical speed/success targets
- Repository ready for tagged release

## Weekly Tracking Table

| Week | Status | Notes | Last Update |
|---|---|---|---|
| Week 1 | Done | Foundation completed in early project phase | 2026-07-04 |
| Week 2 | Done | Bencode module and tests completed | 2026-07-04 |
| Week 3 | Done | Torrent parser and validation completed | 2026-07-04 |
| Week 4 | Done | Tracker communication implemented and tested | 2026-07-04 |
| Week 5 | Not Started | DHT implementation next | 2026-07-04 |
| Week 6 | Not Started | Depends on Week 5 outputs | 2026-07-04 |
| Week 7 | Not Started | Depends on discovery quality | 2026-07-04 |
| Week 8 | Not Started | Protocol/state machine phase | 2026-07-04 |
| Week 9 | Not Started | Piece pipeline phase | 2026-07-04 |
| Week 10 | Not Started | Storage and verification phase | 2026-07-04 |
| Week 11 | Not Started | Resume and stability phase | 2026-07-04 |
| Week 12 | Not Started | Performance and release phase | 2026-07-04 |

## Next Action

Start Week 5 with a 5-day execution breakdown:
- Day 1: DHT message format and bootstrap strategy
- Day 2: Node table and node distance logic
- Day 3: find_node flow
- Day 4: get_peers flow
- Day 5: tests and integration checkpoints
