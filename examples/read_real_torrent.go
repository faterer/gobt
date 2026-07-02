package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"os"
	"strings"
	"time"

	"gop2p/pkg/torrent"
)

// generateRealSHA1Pieces generates actual SHA1 hashes by "hashing" content
// In real usage, these would be SHA1 hashes of actual file content
func generateRealSHA1Pieces(content []byte, pieceLength int64) []byte {
	var pieces []byte
	
	for i := int64(0); i < int64(len(content)); i += pieceLength {
		end := i + pieceLength
		if end > int64(len(content)) {
			end = int64(len(content))
		}
		
		pieceData := content[i:end]
		hash := sha1.Sum(pieceData)
		pieces = append(pieces, hash[:]...)
	}
	
	return pieces
}

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    创建真实的 Torrent 文件示例                              ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════════════╝")

	// 模拟真实文件内容
	fileContent := []byte("这是一个示例文件的内容。\n" +
		"这个文件用于演示如何创建和解析真实的 Torrent 文件。\n" +
		"在实际应用中，这里会是真实文件的二进制内容。\n" +
		"Torrent 文件包含了文件的分片 SHA1 哈希值。\n" +
		"每个分片通常是 16KB 到 2MB。\n")

	pieceLength := int64(50) // 每片 50 字节（用于演示）

	// 生成真实的 SHA1 哈希
	pieces := generateRealSHA1Pieces(fileContent, pieceLength)

	fmt.Printf("\n📄 文件内容大小: %d 字节\n", len(fileContent))
	fmt.Printf("📦 分片大小: %d 字节\n", pieceLength)
	fmt.Printf("🔗 分片数: %d\n", len(pieces)/20)

	// 创建单文件 Torrent
	singleFileTorrent := &torrent.TorrentInfo{
		Announce: "http://tracker.example.com:6969/announce",
		AnnounceList: [][]string{
			{"http://backup1.example.com:6969/announce"},
			{"http://backup2.example.com:6969/announce"},
		},
		CreatedBy:    "gop2p v1.0 - Example",
		CreationDate: time.Now().Unix(),
		Comment:      "这是一个演示 Torrent 文件，包含真实的 SHA1 分片哈希",
		Info: torrent.InfoDict{
			Name:        "example-file.txt",
			Length:      int64(len(fileContent)),
			PieceLength: pieceLength,
			Pieces:      pieces,
		},
	}

	// 验证
	if err := singleFileTorrent.ValidateInfo(); err != nil {
		fmt.Printf("❌ 验证失败: %v\n", err)
		return
	}

	fmt.Println("\n✓ Torrent 元数据验证通过")

	// 编码为 bencode 格式
	encoded, err := torrent.EncodeTorrent(singleFileTorrent)
	if err != nil {
		fmt.Printf("❌ 编码失败: %v\n", err)
		return
	}

	// 保存为 .torrent 文件
	filename := "example-file.torrent"
	if err := os.WriteFile(filename, encoded, 0644); err != nil {
		fmt.Printf("❌ 写入文件失败: %v\n", err)
		return
	}

	fmt.Printf("\n✅ Torrent 文件已创建: %s (%.2f KB)\n", filename, float64(len(encoded))/1024)

	// 显示 bencode 原始格式（前 200 字节）
	fmt.Println("\n📋 Bencode 原始格式 (前 200 字节):")
	preview := encoded
	if len(preview) > 200 {
		preview = preview[:200]
	}
	fmt.Printf("%s...\n", string(preview))

	// 计算 Info Hash
	infoHash, _ := singleFileTorrent.InfoHashHex()
	fmt.Printf("\n🔐 Info Hash (Hex): %s\n", infoHash)

	// 现在读取刚创建的文件
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("现在读取刚创建的 Torrent 文件...")
	fmt.Println(strings.Repeat("=", 80))

	readTorrentFile(filename)
}

func readTorrentFile(filename string) {
	fmt.Printf("\n📂 读取文件: %s\n", filename)

	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("❌ 打开文件失败: %v\n", err)
		return
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, _ := file.Stat()
	fmt.Printf("📊 文件大小: %.2f KB\n", float64(fileInfo.Size())/1024)

	// 解析 Torrent 文件
	t, err := torrent.ParseTorrent(file)
	if err != nil {
		fmt.Printf("❌ 解析失败: %v\n", err)
		return
	}

	fmt.Println("\n✅ 解析成功！\n")

	// 显示详细信息
	displayTorrentDetails(t)
}

func displayTorrentDetails(t *torrent.TorrentInfo) {
	fmt.Println("╔════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                        TORRENT 文件详细信息                                ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════════════╝")

	// 基础信息
	fmt.Println("\n【元数据】")
	fmt.Printf("  📢 Announce:     %s\n", t.Announce)

	if len(t.AnnounceList) > 0 {
		fmt.Println("  📋 备用 Tracker:")
		for i, tier := range t.AnnounceList {
			for j, url := range tier {
				if j == 0 {
					fmt.Printf("     [等级 %d] %s\n", i+1, url)
				} else {
					fmt.Printf("             %s\n", url)
				}
			}
		}
	}

	if t.CreatedBy != "" {
		fmt.Printf("  👤 创建者:     %s\n", t.CreatedBy)
	}

	if t.CreationDate > 0 {
		createdTime := time.Unix(t.CreationDate, 0).Format("2006-01-02 15:04:05")
		fmt.Printf("  ⏰ 创建时间:     %s\n", createdTime)
	}

	if t.Comment != "" {
		fmt.Printf("  💬 备注:       %s\n", t.Comment)
	}

	// 内容信息
	fmt.Println("\n【内容信息】")
	fmt.Printf("  📁 名称:       %s\n", t.Info.Name)
	fmt.Printf("  📏 模式:       %s\n", getMode(t.Mode()))
	fmt.Printf("  💾 总大小:     %s\n", formatBytes(t.TotalSize()))
	fmt.Printf("  🧩 分片大小:   %s\n", formatBytes(t.Info.PieceLength))
	fmt.Printf("  🔢 分片数:     %d\n", t.NumPieces())

	// Info Hash
	infoHash, _ := t.InfoHashHex()
	infoHashBytes, _ := t.InfoHashBytes()
	fmt.Println("\n【Info Hash】")
	fmt.Printf("  16进制格式 (40 字符):  %s\n", infoHash)
	fmt.Printf("  二进制格式 (20 字节):  %v\n", infoHashBytes)

	// 磁力链接
	fmt.Println("\n【磁力链接】")
	magnet := fmt.Sprintf("magnet:?xt=urn:btih:%s&dn=%s", infoHash, t.Info.Name)
	if t.Announce != "" {
		magnet += "&tr=" + t.Announce
	}
	fmt.Printf("  %s\n", magnet)

	// 文件信息
	if t.Mode() == torrent.SingleFile {
		fmt.Println("\n【文件信息】(单文件模式)")
		fmt.Printf("  📄 文件名:     %s\n", t.Info.Name)
		fmt.Printf("  📊 大小:       %s\n", formatBytes(t.Info.Length))
	} else {
		fmt.Println("\n【文件列表】(多文件模式)")
		fmt.Printf("  📂 包含文件:   %d 个\n", len(t.Info.Files))
		for i, f := range t.Info.Files {
			pathStr := bytes.Join([][]byte{}, []byte("/"))
			for j, p := range f.Path {
				if j > 0 {
					pathStr = append(pathStr, '/')
				}
				pathStr = append(pathStr, []byte(p)...)
			}
			fmt.Printf("     %d. %s (%s)\n", i+1, string(pathStr), formatBytes(f.Length))
		}
	}

	// 分片哈希值
	fmt.Println("\n【分片哈希值】(SHA1)")
	numPieces := t.NumPieces()
	displayLimit := 5
	if numPieces <= displayLimit {
		displayLimit = numPieces
	}

	for i := 0; i < displayLimit; i++ {
		pieceHash := t.GetPieceHex(i)
		fmt.Printf("  分片 %3d: %s\n", i, pieceHash)
	}

	if numPieces > displayLimit {
		fmt.Printf("  ... 还有 %d 个分片\n", numPieces-displayLimit)
	}

	// 验证
	fmt.Println("\n【数据验证】")
	if err := t.ValidateInfo(); err != nil {
		fmt.Printf("  ❌ 验证失败: %v\n", err)
	} else {
		fmt.Println("  ✅ 所有字段有效")
	}

	// 统计信息
	fmt.Println("\n【下载统计】")
	fmt.Printf("  📥 总大小:     %s\n", formatBytes(t.TotalSize()))
	fmt.Printf("  ⚡ 以 1Mbps 速度下载需要: ~%.1f 小时\n",
		float64(t.TotalSize())/(125000*3600))
	fmt.Printf("  🔗 分片总数:   %d\n", t.NumPieces())
	fmt.Printf("  📦 平均分片大小: %.2f KB\n",
		float64(t.Info.PieceLength)/1024)

	fmt.Println("\n" + strings.Repeat("=", 80))
}

func formatBytes(bytes int64) string {
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

func getMode(mode torrent.Mode) string {
	if mode == torrent.SingleFile {
		return "单文件"
	}
	return "多文件"
}
