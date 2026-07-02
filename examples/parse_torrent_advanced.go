package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"gop2p/pkg/torrent"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: parse_torrent_advanced <torrent-file>")
		fmt.Println("\nExample:")
		fmt.Println("  parse_torrent_advanced ubuntu-22.04.torrent")
		os.Exit(1)
	}

	filename := os.Args[1]

	// Open and parse torrent file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	t, err := torrent.ParseTorrent(file)
	if err != nil {
		fmt.Printf("Error parsing torrent: %v\n", err)
		os.Exit(1)
	}

	// Validate
	if err := t.ValidateInfo(); err != nil {
		fmt.Printf("Validation error: %v\n", err)
		os.Exit(1)
	}

	// Display header
	fmt.Println("╔════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                          TORRENT FILE INFORMATION                          ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════════════╝")

	// Basic info
	fmt.Println("\n=== METADATA ===")
	fmt.Printf("Announce URL: %s\n", t.Announce)

	if len(t.AnnounceList) > 0 {
		fmt.Println("Backup Trackers:")
		for i, tier := range t.AnnounceList {
			for j, url := range tier {
				if j == 0 {
					fmt.Printf("  Tier %d: %s\n", i+1, url)
				} else {
					fmt.Printf("           %s\n", url)
				}
			}
		}
	}

	if t.CreatedBy != "" {
		fmt.Printf("Created By: %s\n", t.CreatedBy)
	}

	if t.CreationDate > 0 {
		fmt.Printf("Creation Date: %s\n", formatTimestamp(t.CreationDate))
	}

	if t.Comment != "" {
		fmt.Printf("Comment: %s\n", t.Comment)
	}

	// Info dict
	fmt.Println("\n=== CONTENT INFORMATION ===")
	fmt.Printf("Name: %s\n", t.Info.Name)
	fmt.Printf("Mode: %s\n", getModeString(t.Mode()))
	fmt.Printf("Total Size: %s\n", formatSize(t.TotalSize()))
	fmt.Printf("Piece Length: %s (%d bytes)\n", formatSize(t.Info.PieceLength), t.Info.PieceLength)
	fmt.Printf("Number of Pieces: %d\n", t.NumPieces())

	// Info Hash
	infoHash, _ := t.InfoHashHex()
	infoHashBytes, _ := t.InfoHashBytes()
	fmt.Println("\n=== INFO HASH ===")
	fmt.Printf("Hex Format (40 chars): %s\n", infoHash)
	fmt.Printf("Raw Format (20 bytes): %v\n", infoHashBytes)

	// Files
	if t.Mode() == torrent.SingleFile {
		fmt.Println("\n=== FILE (SINGLE-FILE MODE) ===")
		fmt.Printf("Filename: %s\n", t.Info.Name)
		fmt.Printf("Size: %s\n", formatSize(t.Info.Length))
	} else {
		fmt.Println("\n=== FILES (MULTI-FILE MODE) ===")
		fmt.Printf("Total Files: %d\n", len(t.Info.Files))
		totalSize := int64(0)
		for i, f := range t.Info.Files {
			pathStr := filepath.Join(f.Path...)
			fmt.Printf("\n  %d. %s\n", i+1, pathStr)
			fmt.Printf("     Size: %s\n", formatSize(f.Length))
			totalSize += f.Length
		}
		if totalSize != t.TotalSize() {
			fmt.Printf("\nWarning: Sum of file sizes (%s) != declared total size (%s)\n",
				formatSize(totalSize), formatSize(t.TotalSize()))
		}
	}

	// Piece hashes
	fmt.Println("\n=== PIECE HASHES (SHA1) ===")
	fmt.Printf("First 5 pieces (of %d):\n", t.NumPieces())
	for i := 0; i < 5 && i < t.NumPieces(); i++ {
		pieceHash := t.GetPieceHex(i)
		fmt.Printf("  Piece %3d: %s\n", i, pieceHash)
	}
	if t.NumPieces() > 5 {
		fmt.Printf("  ... (%d more pieces)\n", t.NumPieces()-5)
	}

	// Magnet link
	fmt.Println("\n=== MAGNET LINK ===")
	magnetLink := generateMagnetLink(t, infoHash)
	fmt.Printf("%s\n", magnetLink)

	// Summary stats
	fmt.Println("\n=== SUMMARY ===")
	fmt.Printf("Download Speed Estimate (at 1Mbps): %.1f hours\n",
		float64(t.TotalSize())/(125000*3600))

	fmt.Println("\n✓ Torrent file parsed successfully!")
}

func formatSize(bytes int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	size := float64(bytes)

	for _, unit := range units {
		if size < 1024.0 {
			if size < 10 && unit != "B" {
				return fmt.Sprintf("%.2f %s", size, unit)
			}
			return fmt.Sprintf("%.1f %s", size, unit)
		}
		size /= 1024.0
	}

	return fmt.Sprintf("%.2f PB", size)
}

func formatTimestamp(unix int64) string {
	// Simple format for demonstration
	return fmt.Sprintf("%d (Unix timestamp)", unix)
}

func getModeString(mode torrent.Mode) string {
	if mode == torrent.SingleFile {
		return "Single File"
	}
	return "Multi-File"
}

func generateMagnetLink(t *torrent.TorrentInfo, infoHash string) string {
	link := fmt.Sprintf("magnet:?xt=urn:btih:%s&dn=%s", infoHash, t.Info.Name)

	// Add tracker URLs
	if t.Announce != "" {
		link += fmt.Sprintf("&tr=%s", t.Announce)
	}

	for _, tier := range t.AnnounceList {
		for _, url := range tier {
			link += fmt.Sprintf("&tr=%s", url)
		}
	}

	return link
}
