package main

import (
	"fmt"
	"gop2p/pkg/utils"
)

func main() {
	fmt.Println("=== gop2p BitTorrent Client ===")
	fmt.Printf("Version: %s\n", utils.Version())
	fmt.Println("Ready to download torrents!")
}
