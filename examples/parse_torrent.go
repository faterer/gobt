package main

import (
	"fmt"
	"gobt/pkg/torrent"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)


// ParseTorrentFile parses a .torrent file and returns its information
func ParseTorrentFile(filename string) (*torrent.TorrentInfo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Use the torrent package to parse
	return torrent.ParseTorrent(file)
}


// PrintTorrentInfo prints the parsed torrent information
func PrintTorrentInfo(t *torrent.TorrentInfo) {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("📋 TORRENT FILE INFORMATION")
	fmt.Println(strings.Repeat("=", 80))

	// Basic information
	fmt.Println("\n📌 Basic Information:")
	fmt.Printf("  Name:              %s\n", t.Info.Name)
	fmt.Printf("  Total Size:        %s\n", formatBytes(t.TotalSize()))
	fmt.Printf("  Piece Length:      %s\n", formatBytes(t.Info.PieceLength))
	fmt.Printf("  Number of Pieces:  %d\n", t.NumPieces())

	// Tracker information
	fmt.Println("\n📡 Tracker Information:")
	fmt.Printf("  Announce:          %s\n", t.Announce)
	if len(t.AnnounceList) > 0 {
		fmt.Printf("  Announce List:     (%d tiers)\n", len(t.AnnounceList))
		for i, tier := range t.AnnounceList {
			fmt.Printf("    Tier %d:\n", i+1)
			for _, url := range tier {
				fmt.Printf("      - %s\n", url)
			}
		}
	}

	// Metadata
	fmt.Println("\n📝 Metadata:")
	if t.CreatedBy != "" {
		fmt.Printf("  Created By:        %s\n", t.CreatedBy)
	}
	if t.CreationDate > 0 {
		createdTime := time.Unix(t.CreationDate, 0)
		fmt.Printf("  Creation Date:     %s\n", createdTime.Format("2006-01-02 15:04:05"))
	}
	if t.Comment != "" {
		fmt.Printf("  Comment:           %s\n", t.Comment)
	}

	// Info Hash
	fmt.Println("\n🔐 Info Hash:")
	infoHashHex, err := t.InfoHashHex()
	if err == nil {
		fmt.Printf("  %s\n", infoHashHex)
	}

	// Files (for multi-file torrents)
	if t.Mode() == torrent.MultiFile && len(t.Info.Files) > 0 {
		fmt.Printf("\n📁 Files (%d):\n", len(t.Info.Files))
		for i, f := range t.Info.Files {
			pathStr := strings.Join(f.Path, "/")
			fmt.Printf("  %d. %s (%s)\n", i+1, pathStr, formatBytes(f.Length))
		}
	}

	// Statistics
	fmt.Println("\n📊 Statistics:")
	fmt.Printf("  Total Data Size:   %s\n", formatBytes(t.TotalSize()))
	expectedSize := int64(t.NumPieces()) * t.Info.PieceLength
	fmt.Printf("  Expected Size:     %s\n", formatBytes(expectedSize))
	if t.TotalSize() > 0 {
		overhead := float64(expectedSize-t.TotalSize()) / float64(t.TotalSize()) * 100
		if overhead > 0 {
			fmt.Printf("  Overhead:          %.2f%%\n", overhead)
		}
	}

	// Mode
	mode := "Single-file"
	if t.Mode() == torrent.MultiFile {
		mode = "Multi-file"
	}
	fmt.Printf("  Mode:              %s\n", mode)

	fmt.Println(strings.Repeat("=", 80))
}

// formatBytes converts bytes to human-readable format
func formatBytes(bytes int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	size := float64(bytes)
	unitIdx := 0

	for size >= 1024 && unitIdx < len(units)-1 {
		size /= 1024
		unitIdx++
	}

	if unitIdx == 0 {
		return fmt.Sprintf("%d %s", int64(size), units[unitIdx])
	}
	return fmt.Sprintf("%.2f %s", size, units[unitIdx])
}

func main() {
	fmt.Println("\n🚀 Torrent File Parser Example")
	fmt.Println(strings.Repeat("=", 60))

	// Find torrent files - search in current directory and parent directory
	torrentFiles, err := filepath.Glob("*.torrent")
	if err != nil {
		fmt.Printf("Error searching for torrent files: %v\n", err)
		return
	}

	// If no files in current directory, search in parent directory
	if len(torrentFiles) == 0 {
		parentFiles, err := filepath.Glob("../*.torrent")
		if err == nil && len(parentFiles) > 0 {
			torrentFiles = parentFiles
		}
	}

	if len(torrentFiles) == 0 {
		fmt.Println("\n❌ No .torrent files found!")
		fmt.Println("\nMake sure you have .torrent files in:")
		fmt.Println("  • Current directory (./)")
		fmt.Println("  • Parent directory (../)")
		return
	}

	// Parse each torrent file
	sort.Strings(torrentFiles)

	for _, torrentFile := range torrentFiles {
		fmt.Printf("\n📂 Parsing: %s\n", torrentFile)

		t, err := ParseTorrentFile(torrentFile)
		if err != nil {
			fmt.Printf("❌ Error parsing torrent: %v\n", err)
			continue
		}

		PrintTorrentInfo(t)
	}

	fmt.Printf("\n✅ Successfully parsed %d torrent file(s)\n\n", len(torrentFiles))
}
