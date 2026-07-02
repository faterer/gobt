# 读取和解析真实 Torrent 文件 - 完整指南

## 概述

这些示例展示了如何使用 `gop2p` 的 torrent 包来创建、读取和分析真实的 `.torrent` 文件。

## 示例程序

### 1️⃣ `read_real_torrent.go` - 创建和读取单文件 Torrent

这个程序展示了完整的工作流：
- 创建一个带有真实 SHA1 分片哈希的 Torrent
- 将其编码为 Bencode 格式
- 将其保存为 `.torrent` 文件
- 再读取并解析这个文件

**关键功能：**
- 基于实际文件内容生成 SHA1 哈希值
- 显示 Bencode 格式的原始内容
- 详细的元数据显示
- 磁力链接生成

**用法：**
```bash
go run examples/read_real_torrent.go
```

**输出示例：**
```
📄 文件内容大小: 272 字节
📦 分片大小: 50 字节
🔗 分片数: 6

✓ Torrent 元数据验证通过
✅ Torrent 文件已创建: example-file.torrent (0.49 KB)

【文件信息】(单文件模式)
  📄 文件名:     example-file.txt
  📊 大小:       272.0 B

【Info Hash】
  16进制格式: 9e8c387b8d00ff3ceff28f070e65f676a2d43c53
  二进制格式: [158 140 56 123 141 0 255 60...]

【磁力链接】
  magnet:?xt=urn:btih:9e8c387b8d00ff3ceff28f070e65f676a2d43c53&dn=example-file.txt&...
```

---

### 2️⃣ `analyze_torrent.go` - 高级 Torrent 分析工具

这个程序可以自动创建演示文件，也可以分析任何现有的 `.torrent` 文件。

**功能特性：**
- 自动创建多文件演示 Torrent
- 完整的 Torrent 文件分析
- 详细的元数据显示
- 文件列表和大小计算
- 下载时间估算
- 磁力链接生成

**用法：**

不带参数（自动创建演示文件）：
```bash
go run examples/analyze_torrent.go
```

分析现有的 Torrent 文件：
```bash
go run examples/analyze_torrent.go /path/to/your/file.torrent
```

**输出示例：**
```
【🌐 跟踪器信息】
  主要 Announce: http://tracker.example.com:6969/announce
  备用跟踪器列表:
    ├─ [1] http://tracker.opentrackr.org:1337/announce
    └─ [2] udp://tracker.leechers-paradise.org:6969/announce

【📋 元数据】
  创建者: gop2p v1.0
  创建时间: 2026-07-02 10:13:20 CST
  备注: 演示性多文件 Torrent

【📦 内容概览】
  名称: gop2p-project
  模式: 📂 多文件
  总大小: 235.0 B
  分片数: 4

【📁 文件详情】
  多文件模式 (3 个文件)
    1. README.md (83.0 B)
    2. docs/protocol.txt (79.0 B)
    3. src/main.go (73.0 B)

【🔗 分片哈希值 (SHA1)】
  总数: 4
  [  0] e0c9d41f76717b5b38afeccec1bafe534e455299
  [  1] 95670a4ae039eb060016d96daaf4053c62bc767c
  ...

【⚡ 下载统计】
  总大小: 235.0 B
  下载估算 (1Mbps): ~0.0 小时
  下载估算 (10Mbps): ~0.0 分钟
```

---

### 3️⃣ `parse_torrent_advanced.go` - Torrent 文件解析器

这是基于 torrent 包的解析工具，展示了如何读取和显示 Torrent 的所有信息。

**用法：**
```bash
go run examples/parse_torrent_advanced.go your-file.torrent
```

---

## 核心代码示例

### 创建一个 Torrent 文件

```go
package main

import (
    "gop2p/pkg/torrent"
    "os"
    "crypto/sha1"
)

func main() {
    // 生成分片哈希值
    fileContent := []byte("Hello, BitTorrent!")
    pieceLength := int64(10)
    
    var pieces []byte
    for i := int64(0); i < int64(len(fileContent)); i += pieceLength {
        end := i + pieceLength
        if end > int64(len(fileContent)) {
            end = int64(len(fileContent))
        }
        hash := sha1.Sum(fileContent[i:end])
        pieces = append(pieces, hash[:]...)
    }
    
    // 创建 Torrent 信息
    t := &torrent.TorrentInfo{
        Announce: "http://tracker.example.com:6969/announce",
        Info: torrent.InfoDict{
            Name:        "hello.txt",
            Length:      int64(len(fileContent)),
            PieceLength: pieceLength,
            Pieces:      pieces,
        },
    }
    
    // 验证
    if err := t.ValidateInfo(); err != nil {
        panic(err)
    }
    
    // 编码并保存
    encoded, _ := torrent.EncodeTorrent(t)
    os.WriteFile("hello.torrent", encoded, 0644)
    
    // 显示 Info Hash
    hash, _ := t.InfoHashHex()
    println("Info Hash:", hash)
}
```

### 读取和解析 Torrent 文件

```go
package main

import (
    "fmt"
    "os"
    "gop2p/pkg/torrent"
)

func main() {
    // 打开文件
    file, _ := os.Open("hello.torrent")
    defer file.Close()
    
    // 解析
    t, _ := torrent.ParseTorrent(file)
    
    // 显示信息
    fmt.Printf("Name: %s\n", t.Info.Name)
    fmt.Printf("Size: %d bytes\n", t.TotalSize())
    fmt.Printf("Pieces: %d\n", t.NumPieces())
    fmt.Printf("Announce: %s\n", t.Announce)
    
    hash, _ := t.InfoHashHex()
    fmt.Printf("Info Hash: %s\n", hash)
    
    // 获取特定分片的哈希值
    for i := 0; i < t.NumPieces(); i++ {
        fmt.Printf("Piece %d: %s\n", i, t.GetPieceHex(i))
    }
}
```

---

## 关键概念

### Info Hash

Info Hash 是 Torrent 的唯一标识符，通过对 info 字典进行 Bencode 编码后计算 SHA1 得到：

```
Info Hash = SHA1(bencode(info_dict))
```

- **长度**：20 字节（二进制）或 40 个十六进制字符
- **用途**：在 Tracker 通信、DHT 查询和 Peer 连接中使用
- **重要性**：一旦内容改变，Info Hash 就会改变

### 分片哈希值（Pieces）

- 每个分片都有一个 20 字节的 SHA1 哈希值
- 用于验证下载的数据完整性
- 存储在 `info.pieces` 字段中（连续的 20 字节块）

### 单文件 vs 多文件

**单文件模式：**
```go
Info: torrent.InfoDict{
    Name:   "file.iso",
    Length: 1024,              // 文件总大小
    Pieces: []byte{...},        // 分片哈希
}
```

**多文件模式：**
```go
Info: torrent.InfoDict{
    Name:  "folder",
    Files: []torrent.FileInfo{
        {Length: 512, Path: []string{"dir", "file1.txt"}},
        {Length: 512, Path: []string{"file2.txt"}},
    },
    Pieces: []byte{...},
}
```

### Tracker Announce

- **HTTP Tracker**：`http://tracker.example.com:6969/announce`
- **UDP Tracker**：`udp://tracker.example.com:6969/announce`
- **AnnounceList**：备用 Tracker 列表，按优先级分层组织

---

## 实际应用

### 场景 1：验证下载的 Torrent 文件

```go
file, _ := os.Open("downloaded.torrent")
t, _ := torrent.ParseTorrent(file)

if err := t.ValidateInfo(); err != nil {
    fmt.Println("❌ Torrent 文件损坏或无效")
} else {
    fmt.Println("✅ Torrent 文件有效")
}
```

### 场景 2：提取磁力链接

```go
hash, _ := t.InfoHashHex()
magnet := fmt.Sprintf("magnet:?xt=urn:btih:%s&dn=%s&tr=%s",
    hash,
    t.Info.Name,
    t.Announce,
)
fmt.Println(magnet)
```

### 场景 3：计算下载所需时间

```go
totalSize := t.TotalSize()
speedBytesPerSecond := int64(1000000) // 1Mbps = 125KB/s

secondsNeeded := totalSize / speedBytesPerSecond
hoursNeeded := float64(secondsNeeded) / 3600

fmt.Printf("下载需要时间: %.1f 小时\n", hoursNeeded)
```

### 场景 4：列出所有文件

```go
if t.Mode() == torrent.SingleFile {
    fmt.Printf("File: %s (%d bytes)\n", t.Info.Name, t.Info.Length)
} else {
    for i, f := range t.Info.Files {
        pathStr := strings.Join(f.Path, "/")
        fmt.Printf("[%d] %s (%d bytes)\n", i+1, pathStr, f.Length)
    }
}
```

---

## 测试生成的 Torrent 文件

运行示例程序后，会生成以下文件：

1. **example-file.torrent** - 单文件 Torrent
2. **demo-multifile.torrent** - 多文件 Torrent

你可以使用标准的 BitTorrent 客户端来打开这些文件：
- Transmission
- qBittorrent
- Deluge
- 或任何其他支持 BitTorrent 的客户端

---

## 常见错误和解决方案

### "invalid character in string length"
这通常意味着 Bencode 数据格式不正确。确保：
- 数据实际上是 Bencode 格式
- 文件没有被损坏

### "failed to encode info dict"
确保所有字段都是支持的类型（int64、string、[]byte、[]interface{}）。

### Info Hash 不匹配
检查：
- `info.pieces` 是否包含正确的 SHA1 哈希值
- 是否修改了 info 字典中的任何字段
- Bencode 编码是否一致

---

## 下一步

现在你有了完整的 Torrent 文件读写能力，下一步是：
1. **Week 4**：实现 Tracker 通信（HTTP/UDP announce 请求）
2. **Week 5**：实现 DHT 节点发现
3. **Week 6**：实现 Peer Wire 协议（实际文件下载）

这些都将使用 Info Hash 来与网络进行通信。

---

## 相关文件

- `pkg/torrent/torrent.go` - Torrent 包实现
- `pkg/torrent/torrent_test.go` - 测试用例
- `pkg/bencode/struct.go` - 结构体序列化支持
- `examples/read_real_torrent.go` - 本文档的第一个示例
- `examples/analyze_torrent.go` - 本文档的第二个示例
- `examples/parse_torrent_advanced.go` - 本文档的第三个示例
