package main

import (
	"crypto/sha1"
	"fmt"
	"gop2p/pkg/bencode"
	"io/ioutil"
	"time"
)

// Create a simple demo .torrent file for testing
func init() {
	// Generate piece hashes
	pieceHashes := ""
	for i := 0; i < 3; i++ {
		hash := sha1.Sum([]byte(fmt.Sprintf("demo_piece_%d", i)))
		pieceHashes += string(hash[:])
	}

	// Create torrent structure
	infoDict := map[string]interface{}{
		"name":         "demo-example.txt",
		"length":       int64(150),
		"piece length": int64(50),
		"pieces":       pieceHashes,
	}

	torrentDict := map[string]interface{}{
		"announce":      "http://tracker.example.com:6969/announce",
		"info":          infoDict,
		"creation date": time.Now().Unix(),
		"created by":    "gop2p example",
	}

	// Encode to bencode
	encoder := bencode.NewEncoder()
	encoded, err := encoder.Encode(torrentDict)
	if err != nil {
		return
	}

	// Write to file
	ioutil.WriteFile("example-demo.torrent", encoded, 0644)
}
