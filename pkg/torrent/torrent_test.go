package torrent

import (
	"bytes"
	"testing"
)

func TestNewTorrent(t *testing.T) {
	torrent := &TorrentInfo{
		Announce: "http://tracker.example.com:6969/announce",
		Info: InfoDict{
			Name:        "test.txt",
			Length:      1024,
			PieceLength: 16384,
			Pieces:      make([]byte, 20), // 1 piece
		},
	}

	if torrent.Announce != "http://tracker.example.com:6969/announce" {
		t.Errorf("Expected announce URL, got %s", torrent.Announce)
	}

	if torrent.Info.Name != "test.txt" {
		t.Errorf("Expected name 'test.txt', got %s", torrent.Info.Name)
	}
}

func TestMode(t *testing.T) {
	// Single-file mode
	singleFile := &TorrentInfo{
		Announce: "http://tracker.example.com:6969/announce",
		Info: InfoDict{
			Name:        "test.txt",
			Length:      1024,
			PieceLength: 16384,
			Pieces:      make([]byte, 20),
		},
	}

	if singleFile.Mode() != SingleFile {
		t.Error("Expected SingleFile mode")
	}

	// Multi-file mode
	multiFile := &TorrentInfo{
		Announce: "http://tracker.example.com:6969/announce",
		Info: InfoDict{
			Name:        "my-folder",
			PieceLength: 16384,
			Pieces:      make([]byte, 20),
			Files: []FileInfo{
				{Length: 512, Path: []string{"file1.txt"}},
				{Length: 512, Path: []string{"subfolder", "file2.txt"}},
			},
		},
	}

	if multiFile.Mode() != MultiFile {
		t.Error("Expected MultiFile mode")
	}
}

func TestTotalSize(t *testing.T) {
	tests := []struct {
		name        string
		torrent     *TorrentInfo
		expectedLen int64
	}{
		{
			name: "single-file mode",
			torrent: &TorrentInfo{
				Announce: "http://tracker.example.com:6969/announce",
				Info: InfoDict{
					Name:        "test.txt",
					Length:      2048,
					PieceLength: 16384,
					Pieces:      make([]byte, 20),
				},
			},
			expectedLen: 2048,
		},
		{
			name: "multi-file mode",
			torrent: &TorrentInfo{
				Announce: "http://tracker.example.com:6969/announce",
				Info: InfoDict{
					Name:        "folder",
					PieceLength: 16384,
					Pieces:      make([]byte, 20),
					Files: []FileInfo{
						{Length: 512, Path: []string{"file1.txt"}},
						{Length: 512, Path: []string{"file2.txt"}},
						{Length: 1024, Path: []string{"subfolder", "file3.txt"}},
					},
				},
			},
			expectedLen: 2048,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.torrent.TotalSize(); got != tt.expectedLen {
				t.Errorf("TotalSize() = %d, want %d", got, tt.expectedLen)
			}
		})
	}
}

func TestNumPieces(t *testing.T) {
	tests := []struct {
		name          string
		piecesData    []byte
		expectedCount int
	}{
		{
			name:          "1 piece",
			piecesData:    make([]byte, 20),
			expectedCount: 1,
		},
		{
			name:          "3 pieces",
			piecesData:    make([]byte, 60),
			expectedCount: 3,
		},
		{
			name:          "0 pieces",
			piecesData:    make([]byte, 0),
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			torrent := &TorrentInfo{
				Announce: "http://tracker.example.com:6969/announce",
				Info: InfoDict{
					Name:        "test.txt",
					Length:      1024,
					PieceLength: 16384,
					Pieces:      tt.piecesData,
				},
			}

			if got := torrent.NumPieces(); got != tt.expectedCount {
				t.Errorf("NumPieces() = %d, want %d", got, tt.expectedCount)
			}
		})
	}
}

func TestGetPiece(t *testing.T) {
	// Create pieces data with known values
	pieces := make([]byte, 60) // 3 pieces
	// First piece: all 0x01
	for i := 0; i < 20; i++ {
		pieces[i] = 0x01
	}
	// Second piece: all 0x02
	for i := 20; i < 40; i++ {
		pieces[i] = 0x02
	}
	// Third piece: all 0x03
	for i := 40; i < 60; i++ {
		pieces[i] = 0x03
	}

	torrent := &TorrentInfo{
		Announce: "http://tracker.example.com:6969/announce",
		Info: InfoDict{
			Name:        "test.txt",
			Length:      1024,
			PieceLength: 16384,
			Pieces:      pieces,
		},
	}

	tests := []struct {
		index    int
		expected byte
	}{
		{0, 0x01},
		{1, 0x02},
		{2, 0x03},
	}

	for _, tt := range tests {
		piece := torrent.GetPiece(tt.index)
		if piece == nil || piece[0] != tt.expected {
			t.Errorf("GetPiece(%d) first byte = 0x%02x, want 0x%02x", tt.index, piece[0], tt.expected)
		}
	}

	// Test out-of-bounds
	if torrent.GetPiece(10) != nil {
		t.Error("GetPiece(10) should return nil")
	}

	if torrent.GetPiece(-1) != nil {
		t.Error("GetPiece(-1) should return nil")
	}
}

func TestGetPieceHex(t *testing.T) {
	pieces := make([]byte, 20)
	pieces[0] = 0xAB
	pieces[1] = 0xCD

	torrent := &TorrentInfo{
		Announce: "http://tracker.example.com:6969/announce",
		Info: InfoDict{
			Name:        "test.txt",
			Length:      1024,
			PieceLength: 16384,
			Pieces:      pieces,
		},
	}

	hex := torrent.GetPieceHex(0)
	if hex[:4] != "abcd" {
		t.Errorf("GetPieceHex(0) starts with %s, want abcd", hex[:4])
	}

	if len(hex) != 40 {
		t.Errorf("GetPieceHex(0) length = %d, want 40", len(hex))
	}
}

func TestInfoHash(t *testing.T) {
	torrent := &TorrentInfo{
		Announce: "http://tracker.example.com:6969/announce",
		Info: InfoDict{
			Name:        "test.txt",
			Length:      1024,
			PieceLength: 16384,
			Pieces:      make([]byte, 20),
		},
	}

	hash, err := torrent.InfoHash()
	if err != nil {
		t.Fatalf("InfoHash() error: %v", err)
	}

	if len(hash) != 20 {
		t.Errorf("InfoHash length = %d, want 20", len(hash))
	}

	// Same torrent should produce same hash
	hash2, err := torrent.InfoHash()
	if err != nil {
		t.Fatalf("Second InfoHash() error: %v", err)
	}

	if !bytes.Equal(hash, hash2) {
		t.Error("InfoHash should be deterministic")
	}

	// Different torrent should produce different hash
	torrent2 := &TorrentInfo{
		Announce: "http://tracker.example.com:6969/announce",
		Info: InfoDict{
			Name:        "different.txt",
			Length:      1024,
			PieceLength: 16384,
			Pieces:      make([]byte, 20),
		},
	}

	hash3, err := torrent2.InfoHash()
	if err != nil {
		t.Fatalf("Second torrent InfoHash() error: %v", err)
	}

	if bytes.Equal(hash, hash3) {
		t.Error("Different torrents should have different hashes")
	}
}

func TestInfoHashHex(t *testing.T) {
	torrent := &TorrentInfo{
		Announce: "http://tracker.example.com:6969/announce",
		Info: InfoDict{
			Name:        "test.txt",
			Length:      1024,
			PieceLength: 16384,
			Pieces:      make([]byte, 20),
		},
	}

	hex, err := torrent.InfoHashHex()
	if err != nil {
		t.Fatalf("InfoHashHex() error: %v", err)
	}

	if len(hex) != 40 {
		t.Errorf("InfoHashHex length = %d, want 40", len(hex))
	}

	// Check it's valid hex
	for _, ch := range hex {
		if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f')) {
			t.Errorf("InfoHashHex contains non-hex character: %c", ch)
		}
	}
}

func TestEncodeTorrent(t *testing.T) {
	torrent := &TorrentInfo{
		Announce: "http://tracker.example.com:6969/announce",
		Info: InfoDict{
			Name:        "test.txt",
			Length:      1024,
			PieceLength: 16384,
			Pieces:      make([]byte, 20),
		},
	}

	encoded, err := EncodeTorrent(torrent)
	if err != nil {
		t.Fatalf("EncodeTorrent() error: %v", err)
	}

	if len(encoded) == 0 {
		t.Error("EncodeTorrent produced empty result")
	}

	// Should start with 'd' (bencode dict)
	if encoded[0] != 'd' {
		t.Errorf("EncodeTorrent() first byte = 0x%02x, want 0x%02x (d)", encoded[0], 'd')
	}
}

func TestParseTorrent(t *testing.T) {
	// Create and encode a torrent
	original := &TorrentInfo{
		Announce: "http://tracker.example.com:6969/announce",
		Info: InfoDict{
			Name:        "test.txt",
			Length:      1024,
			PieceLength: 16384,
			Pieces:      make([]byte, 20),
		},
	}

	encoded, err := EncodeTorrent(original)
	if err != nil {
		t.Fatalf("EncodeTorrent() error: %v", err)
	}

	// Parse it back
	parsed, err := ParseTorrent(bytes.NewReader(encoded))
	if err != nil {
		t.Fatalf("ParseTorrent() error: %v", err)
	}

	// Verify fields
	if parsed.Announce != original.Announce {
		t.Errorf("Announce = %s, want %s", parsed.Announce, original.Announce)
	}

	if parsed.Info.Name != original.Info.Name {
		t.Errorf("Info.Name = %s, want %s", parsed.Info.Name, original.Info.Name)
	}

	if parsed.Info.Length != original.Info.Length {
		t.Errorf("Info.Length = %d, want %d", parsed.Info.Length, original.Info.Length)
	}

	if parsed.Info.PieceLength != original.Info.PieceLength {
		t.Errorf("Info.PieceLength = %d, want %d", parsed.Info.PieceLength, original.Info.PieceLength)
	}
}

func TestValidateInfo(t *testing.T) {
	tests := []struct {
		name      string
		torrent   *TorrentInfo
		shouldErr bool
		errMsg    string
	}{
		{
			name: "valid single-file",
			torrent: &TorrentInfo{
				Announce: "http://tracker.example.com:6969/announce",
				Info: InfoDict{
					Name:        "test.txt",
					Length:      1024,
					PieceLength: 16384,
					Pieces:      make([]byte, 20),
				},
			},
			shouldErr: false,
		},
		{
			name: "valid multi-file",
			torrent: &TorrentInfo{
				Announce: "http://tracker.example.com:6969/announce",
				Info: InfoDict{
					Name:        "folder",
					PieceLength: 16384,
					Pieces:      make([]byte, 20),
					Files: []FileInfo{
						{Length: 512, Path: []string{"file1.txt"}},
					},
				},
			},
			shouldErr: false,
		},
		{
			name: "missing announce",
			torrent: &TorrentInfo{
				Info: InfoDict{
					Name:        "test.txt",
					Length:      1024,
					PieceLength: 16384,
					Pieces:      make([]byte, 20),
				},
			},
			shouldErr: true,
			errMsg:    "announce",
		},
		{
			name: "missing name",
			torrent: &TorrentInfo{
				Announce: "http://tracker.example.com:6969/announce",
				Info: InfoDict{
					Length:      1024,
					PieceLength: 16384,
					Pieces:      make([]byte, 20),
				},
			},
			shouldErr: true,
			errMsg:    "name",
		},
		{
			name: "invalid piece length",
			torrent: &TorrentInfo{
				Announce: "http://tracker.example.com:6969/announce",
				Info: InfoDict{
					Name:        "test.txt",
					Length:      1024,
					PieceLength: 0,
					Pieces:      make([]byte, 20),
				},
			},
			shouldErr: true,
			errMsg:    "piece length",
		},
		{
			name: "invalid pieces length",
			torrent: &TorrentInfo{
				Announce: "http://tracker.example.com:6969/announce",
				Info: InfoDict{
					Name:        "test.txt",
					Length:      1024,
					PieceLength: 16384,
					Pieces:      make([]byte, 21), // Not multiple of 20
				},
			},
			shouldErr: true,
			errMsg:    "multiple of 20",
		},
		{
			name: "negative length in single-file",
			torrent: &TorrentInfo{
				Announce: "http://tracker.example.com:6969/announce",
				Info: InfoDict{
					Name:        "test.txt",
					Length:      -1,
					PieceLength: 16384,
					Pieces:      make([]byte, 20),
				},
			},
			shouldErr: true,
			errMsg:    "length",
		},
		{
			name: "empty files list",
			torrent: &TorrentInfo{
				Announce: "http://tracker.example.com:6969/announce",
				Info: InfoDict{
					Name:        "folder",
					PieceLength: 16384,
					Pieces:      make([]byte, 20),
					Files:       []FileInfo{},
				},
			},
			shouldErr: true,
			errMsg:    "either",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.torrent.ValidateInfo()
			if (err != nil) != tt.shouldErr {
				t.Fatalf("ValidateInfo() error = %v, shouldErr = %v", err, tt.shouldErr)
			}
			if tt.shouldErr && err != nil && !bytes.Contains([]byte(err.Error()), []byte(tt.errMsg)) {
				t.Errorf("ValidateInfo() error message should contain %q, got %q", tt.errMsg, err.Error())
			}
		})
	}
}

func TestRoundTripEncodeDecode(t *testing.T) {
	original := &TorrentInfo{
		Announce: "http://tracker.example.com:6969/announce",
		AnnounceList: [][]string{
			{"http://backup1.example.com:6969/announce"},
			{"http://backup2.example.com:6969/announce"},
		},
		CreatedBy:    "gop2p",
		CreationDate: 1234567890,
		Comment:      "Test torrent",
		Info: InfoDict{
			Name:        "test-folder",
			PieceLength: 16384,
			Pieces:      make([]byte, 40), // 2 pieces
			Files: []FileInfo{
				{Length: 1024, Path: []string{"file1.txt"}},
				{Length: 2048, Path: []string{"subfolder", "file2.bin"}},
			},
		},
	}

	// Encode
	encoded, err := EncodeTorrent(original)
	if err != nil {
		t.Fatalf("EncodeTorrent() error: %v", err)
	}

	// Decode
	decoded, err := ParseTorrent(bytes.NewReader(encoded))
	if err != nil {
		t.Fatalf("ParseTorrent() error: %v", err)
	}

	// Verify all fields
	if decoded.Announce != original.Announce {
		t.Errorf("Announce mismatch")
	}

	if decoded.CreatedBy != original.CreatedBy {
		t.Errorf("CreatedBy mismatch")
	}

	if decoded.CreationDate != original.CreationDate {
		t.Errorf("CreationDate mismatch")
	}

	if decoded.Comment != original.Comment {
		t.Errorf("Comment mismatch")
	}

	if decoded.Info.Name != original.Info.Name {
		t.Errorf("Info.Name mismatch")
	}

	if decoded.Info.PieceLength != original.Info.PieceLength {
		t.Errorf("Info.PieceLength mismatch")
	}

	if len(decoded.Info.Files) != len(original.Info.Files) {
		t.Errorf("Files count mismatch: %d vs %d", len(decoded.Info.Files), len(original.Info.Files))
	}

	for i := range decoded.Info.Files {
		if decoded.Info.Files[i].Length != original.Info.Files[i].Length {
			t.Errorf("File[%d].Length mismatch", i)
		}
	}

	// Verify InfoHash is consistent
	hash1, _ := original.InfoHash()
	hash2, _ := decoded.InfoHash()
	if !bytes.Equal(hash1, hash2) {
		t.Error("InfoHash mismatch after decode")
	}
}
