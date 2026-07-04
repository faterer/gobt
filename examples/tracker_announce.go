//go:build tracker_example

package main

import (
	"fmt"
	"gobt/pkg/tracker"
	"gobt/pkg/torrent"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run tracker_announce.go <torrent-file> <tracker-url>")
		return
	}

	torrentFile := os.Args[1]
	trackerURL := os.Args[2]

	file, err := os.Open(torrentFile)
	if err != nil {
		fmt.Printf("failed to open torrent: %v\n", err)
		return
	}
	defer file.Close()

	torrentInfo, err := torrent.ParseTorrent(file)
	if err != nil {
		fmt.Printf("failed to parse torrent: %v\n", err)
		return
	}

	if err := torrentInfo.ValidateInfo(); err != nil {
		fmt.Printf("invalid torrent: %v\n", err)
		return
	}

	client := tracker.NewClient()
	response, err := client.Announce(trackerURL, torrentInfo, "started")
	if err != nil {
		fmt.Printf("announce failed: %v\n", err)
		return
	}

	fmt.Println("Tracker announce succeeded")
	fmt.Printf("Interval: %d\n", response.Interval)
	fmt.Printf("Seeders: %d\n", response.Complete)
	fmt.Printf("Leechers: %d\n", response.Incomplete)
	fmt.Printf("Peers: %d\n", len(response.Peers))
	for i, peer := range response.Peers {
		fmt.Printf("%d. %s:%d\n", i+1, peer.IP, peer.Port)
	}
}