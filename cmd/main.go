package main

import (
	"fmt"
	"gobt/pkg/utils"
)

func main() {
	fmt.Println("=== gobt BitTorrent Client ===")
	fmt.Printf("Version: %s\n", utils.Version())
	fmt.Println("Ready to download torrents!")
}
