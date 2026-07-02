package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"gop2p/pkg/bencode"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// TorrentInfo represents the parsed torrent file
type TorrentInfo struct {
	Announce     string
	AnnounceList [][]string
	Comment      string
	CreationDate time.Time
	CreatedBy    string
	Name         string
	Length       int64
	PieceLength  int64
	NumPieces    int
	InfoHash     string
	Files        []FileInfo
}

// FileInfo represents a file in a multi-file torrent
type FileInfo struct {
	Path   []string
	Length int64
}

// ParseTorrentFile parses a .torrent file and returns its information
func ParseTorrentFile(filename string) (*TorrentInfo, error) {
	// Read file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	// Decode bencode
	decoder := bencode.NewDecoder(bytes.NewReader(data))
	decoded, err := decoder.Decode()
	if err != nil {
		return nil, fmt.Errorf("failed to decode bencode: %v", err)
	}

	// Convert to dict
	torrentDict, ok := decoded.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid torrent format: expected dict at top level")
	}

	info := &TorrentInfo{}

	// Parse announce
	if announce, ok := torrentDict["announce"].(string); ok {
		info.Announce = announce
	}

	// Parse announce-list
	if announceList, ok := torrentDict["announce-list"].([]interface{}); ok {
		for _, tier := range announceList {
			if tierList, ok := tier.([]interface{}); ok {
				var tierURLs []string
				for _, url := range tierList {
					if urlStr, ok := url.(string); ok {
						tierURLs = append(tierURLs, urlStr)
					}
				}
				if len(tierURLs) > 0 {
					info.AnnounceList = append(info.AnnounceList, tierURLs)
				}
			}
		}
	}

	// Parse comment
	if comment, ok := torrentDict["comment"].(string); ok {
		info.Comment = comment
	}

	// Parse creation date
	if creationDate, ok := torrentDict["creation date"].(int64); ok {
		info.CreationDate = time.Unix(creationDate, 0)
	}

	// Parse created by
	if createdBy, ok := torrentDict["created by"].(string); ok {
		info.CreatedBy = createdBy
	}

	// Parse info dict
	infoDict, ok := torrentDict["info"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid torrent format: missing 'info' field")
	}

	// Calculate info hash
	encoder := bencode.NewEncoder()
	infoEncoded, err := encoder.Encode(infoDict)
	if err != nil {
		return nil, fmt.Errorf("failed to encode info: %v", err)
	}
	infoHash := sha1.Sum(infoEncoded)
	info.InfoHash = hex.EncodeToString(infoHash[:])

	// Parse name
	if name, ok := infoDict["name"].(string); ok {
		info.Name = name
	}

	// Parse piece length
	if pieceLength, ok := infoDict["piece length"].(int64); ok {
		info.PieceLength = pieceLength
	}

	// Parse pieces
	if pieces, ok := infoDict["pieces"].(string); ok {
		// Each piece hash is 20 bytes (SHA1)
		info.NumPieces = len(pieces) / 20
	}

	// Parse single file or multi-file mode
	if length, ok := infoDict["length"].(int64); ok {
		// Single file mode
		info.Length = length
	} else if files, ok := infoDict["files"].([]interface{}); ok {
		// Multi-file mode
		var totalLength int64
		for _, fileEntry := range files {
			if fileDict, ok := fileEntry.(map[string]interface{}); ok {
				var fileLen int64
				if len, ok := fileDict["length"].(int64); ok {
					fileLen = len
					totalLength += len
				}

				var filePath []string
				if pathList, ok := fileDict["path"].([]interface{}); ok {
					for _, pathPart := range pathList {
						if part, ok := pathPart.(string); ok {
							filePath = append(filePath, part)
						}
					}
				}

				if len(filePath) > 0 {
					info.Files = append(info.Files, FileInfo{
						Path:   filePath,
						Length: fileLen,
					})
				}
			}
		}
		info.Length = totalLength
	}

	return info, nil
}

// PrintTorrentInfo prints the parsed torrent information
func PrintTorrentInfo(info *TorrentInfo) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📋 TORRENT FILE INFORMATION")
	fmt.Println(strings.Repeat("=", 80))

	// Basic info
	fmt.Printf("\n📌 Basic Information:\n")
	fmt.Printf("  Name:              %s\n", info.Name)
	fmt.Printf("  Total Size:        %s\n", formatBytes(info.Length))
	fmt.Printf("  Piece Length:      %s\n", formatBytes(info.PieceLength))
	fmt.Printf("  Number of Pieces:  %d\n", info.NumPieces)

	// Trackers
	fmt.Printf("\n📡 Tracker Information:\n")
	fmt.Printf("  Announce:          %s\n", info.Announce)
	if len(info.AnnounceList) > 0 {
		fmt.Printf("  Announce List:     (%d tiers)\n", len(info.AnnounceList))
		for i, tier := range info.AnnounceList {
			fmt.Printf("    Tier %d:\n", i+1)
			for _, url := range tier {
				fmt.Printf("      - %s\n", url)
			}
		}
	}

	// Metadata
	fmt.Printf("\n📝 Metadata:\n")
	if info.Comment != "" {
		fmt.Printf("  Comment:           %s\n", info.Comment)
	}
	if !info.CreationDate.IsZero() {
		fmt.Printf("  Creation Date:     %s\n", info.CreationDate.Format("2006-01-02 15:04:05"))
	}
	if info.CreatedBy != "" {
		fmt.Printf("  Created By:        %s\n", info.CreatedBy)
	}

	// Info hash
	fmt.Printf("\n🔐 Info Hash:\n")
	fmt.Printf("  %s\n", info.InfoHash)

	// Files
	if len(info.Files) > 0 {
		fmt.Printf("\n📁 Files (%d):\n", len(info.Files))
		for i, file := range info.Files {
			fmt.Printf("  %d. %s (%s)\n", i+1, strings.Join(file.Path, "/"), formatBytes(file.Length))
		}
	}

	// Statistics
	fmt.Printf("\n📊 Statistics:\n")
	fmt.Printf("  Total Data Size:   %s\n", formatBytes(info.Length))
	expectedSize := int64(info.NumPieces) * info.PieceLength
	fmt.Printf("  Expected Size:     %s\n", formatBytes(expectedSize))
	if info.Length > 0 {
		overheadPercent := float64(expectedSize-info.Length) / float64(info.Length) * 100
		if overheadPercent > 0 {
			fmt.Printf("  Overhead:          %.2f%%\n", overheadPercent)
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
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
	fmt.Println("\n🚀 Bencode Torrent Parser Example")
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

		info, err := ParseTorrentFile(torrentFile)
		if err != nil {
			fmt.Printf("❌ Error parsing torrent: %v\n", err)
			continue
		}

		PrintTorrentInfo(info)
	}

	fmt.Printf("\n✅ Successfully parsed %d torrent file(s)\n\n", len(torrentFiles))
}
