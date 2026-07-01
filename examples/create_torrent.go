package main

import (
	"crypto/sha1"
	"fmt"
	"gop2p/pkg/bencode"
	"io/ioutil"
	"strings"
	"time"
)

func main() {
	fmt.Println("🎬 Creating Sample Torrent Files...")
	fmt.Println(strings.Repeat("=", 60))

	// Create a simple torrent structure
	// Simulating Ubuntu 22.04 LTS ISO torrent

	// Generate fake piece hashes (20 bytes each for SHA1)
	pieceHashes := ""
	numPieces := 5
	for i := 0; i < numPieces; i++ {
		// Create a fake but valid SHA1 hash (20 bytes)
		hash := sha1.Sum([]byte(fmt.Sprintf("piece_%d", i)))
		pieceHashes += string(hash[:])
	}

	// Create the info dictionary (the part that gets hashed for info_hash)
	infoDict := map[string]interface{}{
		"name":         "ubuntu-22.04-desktop-amd64.iso",
		"length":       int64(3520000000), // 3.52 GB
		"piece length": int64(262144),    // 256 KB standard piece size
		"pieces":       pieceHashes,      // Binary SHA1 hashes
	}

	// Create the full torrent dictionary
	torrentDict := map[string]interface{}{
		"announce": "http://torrent.ubuntu.com:6969/announce",
		"announce-list": []interface{}{
			[]interface{}{
				"http://torrent.ubuntu.com:6969/announce",
				"http://ipv6.torrent.ubuntu.com:6969/announce",
			},
			[]interface{}{
				"http://torrent.ubuntulinux.nl:6969/announce",
			},
		},
		"info":            infoDict,
		"creation date":   time.Now().Unix(),
		"created by":      "gop2p/4.2.0",
		"comment":         "Official Ubuntu 22.04 LTS Desktop ISO",
	}

	// Encode to bencode
	encoder := bencode.NewEncoder()
	encoded, err := encoder.Encode(torrentDict)
	if err != nil {
		fmt.Printf("❌ Error encoding torrent: %v\n", err)
		return
	}

	// Save to file
	filename := "ubuntu-22.04.torrent"
	err = ioutil.WriteFile(filename, encoded, 0644)
	if err != nil {
		fmt.Printf("❌ Error writing file: %v\n", err)
		return
	}

	fmt.Printf("\n✅ Created: %s\n", filename)
	fmt.Printf("   Size:     %d bytes\n", len(encoded))
	fmt.Printf("   Pieces:   %d\n", numPieces)
	fmt.Printf("   Total:    3.52 GB (simulated)\n")

	// Also create a multi-file torrent example
	fmt.Println("\n📦 Creating Multi-File Torrent...")

	filesInfo := []interface{}{
		map[string]interface{}{
			"length": int64(1073741824),                    // 1 GB
			"path":   []interface{}{"images", "image1.iso"},
		},
		map[string]interface{}{
			"length": int64(2147483648),                    // 2 GB
			"path":   []interface{}{"images", "image2.iso"},
		},
		map[string]interface{}{
			"length": int64(524288000),                     // 500 MB
			"path":   []interface{}{"documents", "README.txt"},
		},
	}

	infoDictMulti := map[string]interface{}{
		"name":         "sample-collection",
		"files":        filesInfo,
		"piece length": int64(262144),
		"pieces":       pieceHashes,
	}

	torrentDictMulti := map[string]interface{}{
		"announce":        "http://tracker.example.com:6969/announce",
		"info":            infoDictMulti,
		"creation date":   time.Now().Unix(),
		"created by":      "gop2p/4.2.0",
		"comment":         "Multi-file torrent example",
	}

	encodedMulti, err := encoder.Encode(torrentDictMulti)
	if err != nil {
		fmt.Printf("❌ Error encoding multi-file torrent: %v\n", err)
		return
	}

	filenameMulti := "sample-collection.torrent"
	err = ioutil.WriteFile(filenameMulti, encodedMulti, 0644)
	if err != nil {
		fmt.Printf("❌ Error writing file: %v\n", err)
		return
	}

	fmt.Printf("\n✅ Created: %s\n", filenameMulti)
	fmt.Printf("   Size:     %d bytes\n", len(encodedMulti))
	fmt.Printf("   Files:    3\n")
	fmt.Printf("   Total:    3.57 GB (3 files)\n")

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("✨ Sample torrent files created successfully!")
	fmt.Println("\nNext, run the parser to view the content:")
	fmt.Println("  go run parse_torrent.go")
	fmt.Println(strings.Repeat("=", 60))
}
