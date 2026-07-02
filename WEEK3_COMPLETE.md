# Week 3 - Torrent File Format

## 🎯 Overview

Week 3 focused on implementing the **torrent file format** package (`pkg/torrent/`), which is the core data structure for BitTorrent. This package handles parsing, creation, and validation of `.torrent` files using the Bencode format created in Week 2.

## 📦 What Was Implemented

### Core Package: `pkg/torrent/`

**File: `torrent.go` (247 lines)**

Key types:
- **`TorrentInfo`** - Main structure representing a complete .torrent file
- **`InfoDict`** - The "info" dictionary that gets hashed for the Info Hash
- **`FileInfo`** - Represents individual files in multi-file mode

Key methods:
- `Mode()` - Determine if torrent is single-file or multi-file
- `TotalSize()` - Calculate total bytes to download
- `NumPieces()` - Get number of pieces
- `GetPiece(index)` - Retrieve specific piece hash
- `GetPieceHex(index)` - Get piece hash as hex string
- `InfoHash()` - Calculate SHA1 hash of info dictionary
- `InfoHashHex()` - Get Info Hash as 40-character hex string
- `InfoHashBytes()` - Get Info Hash as raw 20 bytes
- `ParseTorrent(reader)` - Parse torrent from file/stream
- `EncodeTorrent(t)` - Encode torrent to bencode format
- `ValidateInfo()` - Comprehensive validation of all fields

**Comprehensive Tests: `torrent_test.go` (420 lines, 16 test cases)**

Test coverage (84.1%):
- ✅ Struct creation and initialization
- ✅ Single-file and multi-file modes
- ✅ Total size calculation
- ✅ Piece indexing and hashing
- ✅ Info Hash calculation (deterministic)
- ✅ Hex encoding/decoding
- ✅ Round-trip encode/decode validation
- ✅ Comprehensive field validation
- ✅ Error handling for invalid data

### Extended Bencode Package

**File: `pkg/bencode/struct.go` (300 lines)**

Added struct support to bencode package:
- `Encode(v interface{})` - Enhanced to support structs
- `DecodeBytes(data []byte, v interface{})` - Decode into struct
- `structToDict()` - Convert struct to map[string]interface{}
- `dictToStruct()` - Convert map back to struct
- Full support for bencode struct tags: `bencode:"fieldname"`
- Handles nested structs, pointers, slices, and interface{} types
- Proper reflection-based field mapping

**Tests: `pkg/bencode/struct_test.go` (350 lines, 13 test cases)**

Struct encoding/decoding tests covering:
- ✅ Simple struct encoding
- ✅ Struct decoding with tag mapping
- ✅ Binary fields in structs
- ✅ List fields ([]interface{})
- ✅ Nested struct support
- ✅ Pointer field handling
- ✅ Field tag filtering (bencode:"-")
- ✅ Zero value omission

### Example Programs

**`examples/create_torrent_advanced.go` (170 lines)**
- Creates both single-file and multi-file .torrent files
- Demonstrates TorrentInfo initialization
- Shows validation and Info Hash calculation
- Outputs real torrent files ready for sharing
- Shows announce list configuration

**`examples/parse_torrent_advanced.go` (210 lines)**
- Parses and displays detailed torrent information
- Pretty-prints metadata, file listings, and piece info
- Generates magnet links from torrent metadata
- Estimates download times
- Shows all tracker information

## 🔑 Key Features

### Binary-Safe Format Design
```go
// Torrent Info Hash Calculation
infoBytes, _ := bencode.Encode(t.Info)
hash := sha1.Sum(infoBytes)  // 20-byte binary hash
```

The format naturally handles:
- Mixed ASCII and binary data
- Piece hashes (20-byte SHA1)
- Tracker URLs (ASCII strings)
- File sizes (integers)
- All without special escaping

### Two Torrent Modes

**Single-File Mode:**
```go
InfoDict{
    Name:   "file.iso",
    Length: 3000000000,      // Total file size
    Pieces: []byte{...},      // 20-byte SHA1 hashes
}
```

**Multi-File Mode:**
```go
InfoDict{
    Name:  "folder",
    Files: []FileInfo{
        {Length: 512, Path: []string{"dir", "file1.txt"}},
        {Length: 1024, Path: []string{"dir", "file2.bin"}},
    },
}
```

### Validation

Comprehensive validation ensures:
- All required fields are present
- Piece hashes are 20-byte multiples
- File paths are non-empty
- Exactly one of Length (single) or Files (multi) is set
- Piece length is positive
- File sizes are non-negative

## 📊 Test Coverage

```
pkg/bencode   79.4% coverage
pkg/torrent   84.1% coverage  ✓
pkg/utils    100.0% coverage
```

Total: **126 test cases** covering:
- All encoding/decoding paths
- Edge cases (empty, large, special chars, binary)
- Round-trip validation
- Error handling
- Field validation

## 🧠 Design Insights

### Why No Special Binary Handling?

Bencode's binary safety comes from **length prefixes**:
```
<length>:<data>
```

This means:
- `20:` followed by 20 arbitrary bytes = binary piece hash
- `8:` followed by "announce" = ASCII string
- Single decoder handles both without checking content

### Info Hash Significance

```
Info Hash = SHA1(bencode(info_dict))
```

The Info Hash:
- Uniquely identifies the torrent content
- Used in tracker requests and DHT
- Is 20 bytes (binary) or 40 characters (hex)
- Must remain constant after creation

### Struct Tag Mapping

Bencode tags map struct fields to dictionary keys:
```go
type InfoDict struct {
    Name        string `bencode:"name"`
    PieceLength int64  `bencode:"piece length"`
    Pieces      []byte `bencode:"pieces"`
}
```

Allows clean serialization to/from torrent files.

## 🚀 Example Usage

### Creating a Torrent

```go
import "gop2p/pkg/torrent"

t := &torrent.TorrentInfo{
    Announce: "http://tracker.example.com/announce",
    Info: torrent.InfoDict{
        Name:        "myfile.zip",
        Length:      1024 * 1024,  // 1MB
        PieceLength: 65536,         // 64KB pieces
        Pieces:      realSHA1Hashes, // 20-byte hashes
    },
}

// Validate
if err := t.ValidateInfo(); err != nil {
    log.Fatal(err)
}

// Encode to file
encoded, _ := torrent.EncodeTorrent(t)
os.WriteFile("myfile.torrent", encoded, 0644)

// Get Info Hash
hash, _ := t.InfoHashHex()
fmt.Println("Magnet: magnet:?xt=urn:btih:" + hash)
```

### Parsing a Torrent

```go
file, _ := os.Open("myfile.torrent")
t, _ := torrent.ParseTorrent(file)

fmt.Printf("Name: %s\n", t.Info.Name)
fmt.Printf("Size: %d bytes\n", t.TotalSize())
fmt.Printf("Pieces: %d\n", t.NumPieces())
fmt.Printf("Info Hash: %s\n", t.InfoHashHex())
```

## 📁 File Structure

```
pkg/torrent/
├── torrent.go          (247 lines)  - Core implementation
└── torrent_test.go     (420 lines)  - Comprehensive tests (16 cases)

pkg/bencode/
├── struct.go           (300 lines)  - Struct encoding/decoding
├── struct_test.go      (350 lines)  - Struct tests (13 cases)
├── encoder.go          (180 lines)  - Updated with []byte support
├── decoder.go          (220 lines)  - Unchanged
└── bencode_test.go     (1200+ lines) - Original tests

examples/
├── create_torrent_advanced.go    (170 lines)  - Create torrents demo
├── parse_torrent_advanced.go     (210 lines)  - Parse torrents demo
└── [other Week 1-2 examples]
```

## 🔄 Round-Trip Validation

The test suite validates that:
1. **Encode → Decode**: Encoding and decoding preserves all data
2. **Info Hash Consistency**: Same data always produces same hash
3. **Field Mapping**: Bencode tags correctly map to/from struct fields
4. **Type Preservation**: All field types survive round-trip

Example:
```go
original := &TorrentInfo{...}
encoded, _ := EncodeTorrent(original)
decoded, _ := ParseTorrent(bytes.NewReader(encoded))
// ✓ decoded equals original (Info Hash identical)
```

## 💡 Key Learnings

1. **Bencode elegance**: Length prefixes solve binary safety perfectly
2. **Struct serialization**: Using reflect enables clean struct↔bencode mapping
3. **Info Hash importance**: It's the unique identifier for all torrent operations
4. **Validation matters**: Checking all constraints prevents downstream errors
5. **Binary formats**: Understanding byte-level semantics is crucial

## 🎓 What This Enables

This Week 3 implementation is the foundation for:
- **Week 4**: Tracker communication (announce requests with Info Hash)
- **Week 5**: DHT peer discovery (using Info Hash)
- **Week 6**: Peer wire protocol (file downloads)

Without a solid torrent file format package, none of the networking layers can work.

## ✅ Completion Status

- ✅ TorrentInfo struct with all required fields
- ✅ Single-file and multi-file mode support
- ✅ Info Hash calculation (SHA1 of bencode(info))
- ✅ Comprehensive validation
- ✅ Bencode struct encoding/decoding
- ✅ 84.1% test coverage (16 test cases)
- ✅ Example programs (create and parse)
- ✅ All tests passing
- ✅ Git commits logged

**Status: COMPLETE ✓**

---

**Next: Week 4 - Tracker Communication**
- Implement HTTP/UDP tracker announce requests
- Parse tracker responses (peer list)
- Handle seeders/leechers reporting
