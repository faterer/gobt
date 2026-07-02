package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"gop2p/pkg/torrent"
)

func main() {
	if len(os.Args) < 2 {
		// 如果没有提供参数，创建一个演示文件
		createDemoTorrent()
		fmt.Println("\n" + strings.Repeat("=", 80))
		fmt.Println("现在读取刚才创建的 Torrent 文件...")
		fmt.Println(strings.Repeat("=", 80) + "\n")
		
		analyzeTorrentFile("demo-multifile.torrent")
	} else {
		// 否则读取提供的文件
		analyzeTorrentFile(os.Args[1])
	}
}

func createDemoTorrent() {
	fmt.Println("╔════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                创建演示性的多文件 Torrent                                  ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════════════╝\n")

	// 模拟多个文件的内容
	files := map[string]string{
		"README.md": "# Go P2P BitTorrent 项目\n\n这是一个用 Go 实现的 P2P 文件分享系统。\n",
		"docs/protocol.txt": "BitTorrent 协议规范文档。\n包含了完整的消息格式和状态机。\n",
		"src/main.go": "package main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello P2P\")\n}\n",
	}

	pieceLength := int64(64) // 64 字节每片

	// 构建文件列表
	var filesList []torrent.FileInfo
	for path, content := range files {
		pathParts := strings.Split(path, "/")
		filesList = append(filesList, torrent.FileInfo{
			Length: int64(len(content)),
			Path:   pathParts,
		})
	}

	// 生成虚拟的分片哈希值
	totalSize := int64(0)
	for _, f := range filesList {
		totalSize += f.Length
	}
	pieces := make([]byte, (totalSize/pieceLength+1)*20)
	for i := 0; i < len(pieces); i += 20 {
		hash := sha1.Sum([]byte(fmt.Sprintf("demo-piece-%d", i/20)))
		copy(pieces[i:i+20], hash[:])
	}

	// 创建 Torrent
	t := &torrent.TorrentInfo{
		Announce: "http://tracker.example.com:6969/announce",
		AnnounceList: [][]string{
			{"http://tracker.opentrackr.org:1337/announce"},
			{"udp://tracker.leechers-paradise.org:6969/announce"},
		},
		CreatedBy:    "gop2p v1.0",
		CreationDate: time.Now().Unix(),
		Comment:      "演示性多文件 Torrent - 由 gop2p 创建",
		Info: torrent.InfoDict{
			Name:        "gop2p-project",
			PieceLength: pieceLength,
			Pieces:      pieces,
			Files:       filesList,
		},
	}

	// 验证并保存
	if err := t.ValidateInfo(); err != nil {
		fmt.Printf("❌ 验证失败: %v\n", err)
		return
	}

	encoded, err := torrent.EncodeTorrent(t)
	if err != nil {
		fmt.Printf("❌ 编码失败: %v\n", err)
		return
	}

	if err := os.WriteFile("demo-multifile.torrent", encoded, 0644); err != nil {
		fmt.Printf("❌ 写入失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 演示 Torrent 已创建: demo-multifile.torrent (%.2f KB)\n",
		float64(len(encoded))/1024)
}

func analyzeTorrentFile(filename string) {
	fmt.Printf("📂 分析文件: %s\n", filename)

	// 读取文件
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("❌ 无法读取文件: %v\n", err)
		return
	}

	fmt.Printf("📊 文件大小: %.2f KB\n\n", float64(len(data))/1024)

	// 解析
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("❌ 无法打开文件: %v\n", err)
		return
	}
	defer file.Close()

	t, err := torrent.ParseTorrent(file)
	if err != nil {
		fmt.Printf("❌ 解析失败: %v\n", err)
		return
	}

	// 验证
	if err := t.ValidateInfo(); err != nil {
		fmt.Printf("⚠️  验证警告: %v\n", err)
		return
	}

	fmt.Println("✅ 解析验证成功！\n")

	// 详细显示
	printTorrentAnalysis(t)
}

func printTorrentAnalysis(t *torrent.TorrentInfo) {
	fmt.Println(strings.Repeat("╔", 40))
	fmt.Println("║          TORRENT 文件完整分析          ║")
	fmt.Println(strings.Repeat("╚", 40))

	// 1. 跟踪器信息
	fmt.Println("\n【🌐 跟踪器信息】")
	fmt.Printf("  主要 Announce: %s\n", t.Announce)
	if len(t.AnnounceList) > 0 {
		fmt.Println("  备用跟踪器列表:")
		for i, tier := range t.AnnounceList {
			for j, url := range tier {
				prefix := "    ├─"
				if j == len(tier)-1 && i == len(t.AnnounceList)-1 {
					prefix = "    └─"
				}
				fmt.Printf("%s [%d] %s\n", prefix, i+1, url)
			}
		}
	}

	// 2. 元数据
	fmt.Println("\n【📋 元数据】")
	if t.CreatedBy != "" {
		fmt.Printf("  创建者: %s\n", t.CreatedBy)
	}
	if t.CreationDate > 0 {
		createdTime := time.Unix(t.CreationDate, 0).Format("2006-01-02 15:04:05 MST")
		fmt.Printf("  创建时间: %s\n", createdTime)
	}
	if t.Comment != "" {
		fmt.Printf("  备注: %s\n", t.Comment)
	}

	// 3. 内容概览
	fmt.Println("\n【📦 内容概览】")
	fmt.Printf("  名称: %s\n", t.Info.Name)
	fmt.Printf("  模式: %s\n", getModeStr(t.Mode()))
	fmt.Printf("  总大小: %s\n", formatSize(t.TotalSize()))
	fmt.Printf("  分片大小: %s (%d 字节)\n", formatSize(t.Info.PieceLength), t.Info.PieceLength)
	fmt.Printf("  分片数: %d\n", t.NumPieces())

	// 4. Info Hash
	fmt.Println("\n【🔐 Info Hash (唯一标识)】")
	infoHashHex, _ := t.InfoHashHex()
	infoHashBytes, _ := t.InfoHashBytes()
	fmt.Printf("  十六进制: %s\n", infoHashHex)
	fmt.Printf("  字节表示: %s\n", formatBytes(infoHashBytes))

	// 5. 磁力链接
	fmt.Println("\n【🧲 磁力链接】")
	magnetLink := fmt.Sprintf("magnet:?xt=urn:btih:%s", infoHashHex)
	magnetLink += fmt.Sprintf("&dn=%s", t.Info.Name)
	if t.Announce != "" {
		magnetLink += fmt.Sprintf("&tr=%s", t.Announce)
	}
	fmt.Printf("  %s\n", magnetLink)

	// 6. 文件信息
	fmt.Println("\n【📁 文件详情】")
	if t.Mode() == torrent.SingleFile {
		fmt.Printf("  单文件模式\n")
		fmt.Printf("  文件名: %s\n", t.Info.Name)
		fmt.Printf("  大小: %s\n", formatSize(t.Info.Length))
	} else {
		fmt.Printf("  多文件模式 (%d 个文件)\n", len(t.Info.Files))
		totalSize := int64(0)
		for i, f := range t.Info.Files {
			pathStr := strings.Join(f.Path, "/")
			fmt.Printf("    %d. %s (%s)\n", i+1, pathStr, formatSize(f.Length))
			totalSize += f.Length
		}
		if totalSize != t.TotalSize() {
			fmt.Printf("  ⚠️  警告: 文件大小总和 (%s) ≠ 声明大小 (%s)\n",
				formatSize(totalSize), formatSize(t.TotalSize()))
		}
	}

	// 7. 分片哈希值
	fmt.Println("\n【🔗 分片哈希值 (SHA1)】")
	fmt.Printf("  总数: %d\n", t.NumPieces())
	numShow := 3
	if t.NumPieces() < numShow {
		numShow = t.NumPieces()
	}
	fmt.Printf("  前 %d 个:\n", numShow)
	for i := 0; i < numShow; i++ {
		fmt.Printf("    [%3d] %s\n", i, t.GetPieceHex(i))
	}
	if t.NumPieces() > numShow {
		fmt.Printf("    ... 还有 %d 个\n", t.NumPieces()-numShow)
	}

	// 8. 下载信息
	fmt.Println("\n【⚡ 下载统计】")
	fmt.Printf("  总大小: %s\n", formatSize(t.TotalSize()))
	fmt.Printf("  平均分片: %.2f KB\n", float64(t.Info.PieceLength)/1024)
	fmt.Printf("  下载估算 (1Mbps): ~%.1f 小时\n",
		float64(t.TotalSize())/(125000*3600))
	fmt.Printf("  下载估算 (10Mbps): ~%.1f 分钟\n",
		float64(t.TotalSize())/(1250000*60))

	// 9. 协议信息
	fmt.Println("\n【🔧 协议信息】")
	fmt.Printf("  协议版本: BitTorrent v1\n")
	fmt.Printf("  哈希算法: SHA1\n")
	fmt.Printf("  编码格式: Bencode\n")

	fmt.Println("\n" + strings.Repeat("═", 80))
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

func formatBytes(data []byte) string {
	if len(data) == 0 {
		return "(empty)"
	}
	if len(data) <= 32 {
		return hex.EncodeToString(data)
	}
	return hex.EncodeToString(data[:16]) + "... (" + fmt.Sprintf("%d bytes", len(data)) + ")"
}

func getModeStr(mode torrent.Mode) string {
	if mode == torrent.SingleFile {
		return "📄 单文件"
	}
	return "📂 多文件"
}
