package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gop2p/pkg/torrent"
)

func main() {
	// Example 1: Create a single-file torrent
	createSingleFileTorrent()

	// Example 2: Create a multi-file torrent
	createMultiFileTorrent()
}

func createSingleFileTorrent() {
	fmt.Println("\n=== Creating Single-File Torrent ===\n")

	// Create torrent metadata for a single file
	torrentInfo := &torrent.TorrentInfo{
		Announce:     "http://tracker.ubuntu.com:6969/announce",
		CreatedBy:    "gop2p/1.0",
		CreationDate: time.Now().Unix(),
		Comment:      "Ubuntu 22.04 ISO",
		Info: torrent.InfoDict{
			Name:        "ubuntu-22.04-desktop-amd64.iso",
			Length:      3000000000, // 3GB
			PieceLength: 262144,     // 256KB pieces
			Pieces: generateSHA1Pieces(3000000000 / 262144),
		},
	}

	// Validate
	if err := torrentInfo.ValidateInfo(); err != nil {
		fmt.Printf("Validation error: %v\n", err)
		return
	}

	// Display info
	fmt.Printf("Torrent: %s\n", torrentInfo.Info.Name)
	fmt.Printf("Size: %.2f MB\n", float64(torrentInfo.TotalSize())/1024/1024)
	fmt.Printf("Piece Length: %d KB\n", torrentInfo.Info.PieceLength/1024)
	fmt.Printf("Number of Pieces: %d\n", torrentInfo.NumPieces())

	// Calculate Info Hash
	hash, _ := torrentInfo.InfoHashHex()
	fmt.Printf("Info Hash: %s\n", hash)

	// Encode to bencode
	encoded, err := torrent.EncodeTorrent(torrentInfo)
	if err != nil {
		fmt.Printf("Encoding error: %v\n", err)
		return
	}

	// Write to file
	filename := "ubuntu-22.04.torrent"
	if err := os.WriteFile(filename, encoded, 0644); err != nil {
		fmt.Printf("Write error: %v\n", err)
		return
	}

	fmt.Printf("✓ Torrent saved to: %s (%.2f KB)\n", filename, float64(len(encoded))/1024)
}

func createMultiFileTorrent() {
	fmt.Println("\n=== Creating Multi-File Torrent ===\n")

	// Create torrent metadata for multiple files
	torrentInfo := &torrent.TorrentInfo{
		Announce: "http://tracker.openbittorrent.com:80/announce",
		AnnounceList: [][]string{
			{"http://tracker.opentrackr.org:1337/announce"},
			{"http://tracker.openbittorrent.com:80/announce"},
		},
		CreatedBy:    "gop2p/1.0",
		CreationDate: time.Now().Unix(),
		Comment:      "Go BitTorrent P2P implementation demo",
		Info: torrent.InfoDict{
			Name:        "gop2p-collection",
			PieceLength: 65536, // 64KB pieces
			Pieces:      generateSHA1Pieces(10),
			Files: []torrent.FileInfo{
				{
					Length: 1024 * 1024,           // 1MB
					Path:   []string{"README.md"},
				},
				{
					Length: 512 * 1024,                        // 512KB
					Path:   []string{"docs", "PROTOCOL.md"},
				},
				{
					Length: 2 * 1024 * 1024,              // 2MB
					Path:   []string{"docs", "SPEC.txt"},
				},
				{
					Length: 5 * 1024 * 1024,             // 5MB
					Path:   []string{"examples", "demo.go"},
				},
				{
					Length: 1024,                              // 1KB
					Path:   []string{"examples", "config.json"},
				},
			},
		},
	}

	// Validate
	if err := torrentInfo.ValidateInfo(); err != nil {
		fmt.Printf("Validation error: %v\n", err)
		return
	}

	// Display info
	fmt.Printf("Torrent: %s\n", torrentInfo.Info.Name)
	fmt.Printf("Total Size: %.2f MB\n", float64(torrentInfo.TotalSize())/1024/1024)
	fmt.Printf("Piece Length: %d KB\n", torrentInfo.Info.PieceLength/1024)
	fmt.Printf("Number of Pieces: %d\n", torrentInfo.NumPieces())
	fmt.Printf("Number of Files: %d\n", len(torrentInfo.Info.Files))

	fmt.Println("\nFiles:")
	for i, f := range torrentInfo.Info.Files {
		pathStr := filepath.Join(f.Path...)
		fmt.Printf("  %d. %s (%.2f KB)\n", i+1, pathStr, float64(f.Length)/1024)
	}

	// Calculate Info Hash
	hash, _ := torrentInfo.InfoHashHex()
	fmt.Printf("\nInfo Hash: %s\n", hash)

	// Encode to bencode
	encoded, err := torrent.EncodeTorrent(torrentInfo)
	if err != nil {
		fmt.Printf("Encoding error: %v\n", err)
		return
	}

	// Write to file
	filename := "gop2p-collection.torrent"
	if err := os.WriteFile(filename, encoded, 0644); err != nil {
		fmt.Printf("Write error: %v\n", err)
		return
	}

	fmt.Printf("✓ Torrent saved to: %s (%.2f KB)\n", filename, float64(len(encoded))/1024)
}

// generateSHA1Pieces generates fake SHA1 hashes for demo purposes
// In real usage, these would be actual SHA1 hashes of file pieces
func generateSHA1Pieces(count int64) []byte {
	pieces := make([]byte, 0, count*20)
	for i := int64(0); i < count; i++ {
		// Generate a pseudo-random 20-byte hash
		hash := make([]byte, 20)
		for j := 0; j < 20; j++ {
			hash[j] = byte((i*31 + int64(j)*17) % 256)
		}
		pieces = append(pieces, hash...)
	}
	return pieces
}
