package torrent

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"gop2p/pkg/bencode"
	"io"
)

// TorrentInfo represents the structure of a .torrent file
type TorrentInfo struct {
	Announce      string      `bencode:"announce"`
	AnnounceList  [][]string  `bencode:"announce-list"`
	CreatedBy     string      `bencode:"created by"`
	CreationDate  int64       `bencode:"creation date"`
	Comment       string      `bencode:"comment"`
	Info          InfoDict    `bencode:"info"`
}

// InfoDict represents the "info" dictionary in a torrent file
// This is the section that gets hashed to create the Info Hash
type InfoDict struct {
	Name        string     `bencode:"name"`
	PieceLength int64      `bencode:"piece length"`
	Pieces      []byte     `bencode:"pieces"`
	Length      int64      `bencode:"length"`      // single-file mode
	Files       []FileInfo `bencode:"files"`       // multi-file mode
}

// FileInfo represents a single file in multi-file mode
type FileInfo struct {
	Length int64    `bencode:"length"`
	Path   []string `bencode:"path"`
}

// Mode represents whether torrent is single-file or multi-file
type Mode int

const (
	SingleFile Mode = iota
	MultiFile
)

// Mode returns the torrent mode (SingleFile or MultiFile)
func (t *TorrentInfo) Mode() Mode {
	if len(t.Info.Files) > 0 {
		return MultiFile
	}
	return SingleFile
}

// TotalSize returns total bytes to be downloaded
func (t *TorrentInfo) TotalSize() int64 {
	if t.Mode() == SingleFile {
		return t.Info.Length
	}

	total := int64(0)
	for _, f := range t.Info.Files {
		total += f.Length
	}
	return total
}

// NumPieces returns the number of pieces
func (t *TorrentInfo) NumPieces() int {
	// Each piece is represented by a 20-byte SHA1 hash
	return len(t.Info.Pieces) / 20
}

// GetPiece returns the SHA1 hash of piece at index
func (t *TorrentInfo) GetPiece(index int) []byte {
	if index < 0 || index >= t.NumPieces() {
		return nil
	}
	start := index * 20
	end := start + 20
	return t.Info.Pieces[start:end]
}

// GetPieceHex returns the hex-encoded SHA1 hash of piece at index
func (t *TorrentInfo) GetPieceHex(index int) string {
	piece := t.GetPiece(index)
	if piece == nil {
		return ""
	}
	return hex.EncodeToString(piece)
}

// InfoHash calculates and returns the Info Hash (SHA1 of bencode(info))
// Returns the hash as a 20-byte slice
func (t *TorrentInfo) InfoHash() ([]byte, error) {
	infoBencoded, err := bencode.Encode(t.Info)
	if err != nil {
		return nil, fmt.Errorf("failed to encode info dict: %w", err)
	}

	hasher := sha1.New()
	hasher.Write(infoBencoded)
	return hasher.Sum(nil), nil
}

// InfoHashHex returns the Info Hash as a 40-character hex string
func (t *TorrentInfo) InfoHashHex() (string, error) {
	hash, err := t.InfoHash()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash), nil
}

// InfoHashBytes returns the Info Hash as raw bytes for binary operations
// This is what gets used in peer connections and tracker requests
func (t *TorrentInfo) InfoHashBytes() ([]byte, error) {
	return t.InfoHash()
}

// ParseTorrent parses a torrent file from a reader
// Expects the reader to contain bencode-encoded torrent data
func ParseTorrent(r io.Reader) (*TorrentInfo, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read torrent file: %w", err)
	}

	var torrent TorrentInfo
	err = bencode.DecodeBytes(data, &torrent)
	if err != nil {
		return nil, fmt.Errorf("failed to decode torrent: %w", err)
	}

	return &torrent, nil
}

// EncodeTorrent encodes a TorrentInfo to bencode format
func EncodeTorrent(t *TorrentInfo) ([]byte, error) {
	return bencode.Encode(t)
}

// ValidateInfo checks if the torrent has required fields
func (t *TorrentInfo) ValidateInfo() error {
	if t.Announce == "" {
		return fmt.Errorf("announce field is required")
	}

	if t.Info.Name == "" {
		return fmt.Errorf("info.name field is required")
	}

	if t.Info.PieceLength <= 0 {
		return fmt.Errorf("info.piece length must be positive")
	}

	if len(t.Info.Pieces)%20 != 0 {
		return fmt.Errorf("pieces field must be multiple of 20 bytes")
	}

	// Check that exactly one of Length or Files is set
	hasLength := t.Info.Length > 0
	hasFiles := len(t.Info.Files) > 0

	if !hasLength && !hasFiles {
		return fmt.Errorf("either info.length (single-file) or info.files (multi-file) must be set")
	}

	if hasLength && hasFiles {
		return fmt.Errorf("cannot have both info.length and info.files")
	}

	if hasLength {
		if t.Info.Length < 0 {
			return fmt.Errorf("info.length must be non-negative in single-file mode")
		}
	}

	if hasFiles {
		for i, f := range t.Info.Files {
			if f.Length < 0 {
				return fmt.Errorf("file[%d].length must be non-negative", i)
			}
			if len(f.Path) == 0 {
				return fmt.Errorf("file[%d].path must not be empty", i)
			}
		}
	}

	return nil
}
