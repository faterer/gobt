# 快速参考：Torrent 文件读取

## 🚀 快速开始

### 读取一个现有的 .torrent 文件

```go
package main

import (
    "fmt"
    "os"
    "gop2p/pkg/torrent"
)

func main() {
    // 打开文件
    file, err := os.Open("my-file.torrent")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    // 解析
    t, err := torrent.ParseTorrent(file)
    if err != nil {
        panic(err)
    }

    // 显示基本信息
    fmt.Printf("📁 名称: %s\n", t.Info.Name)
    fmt.Printf("💾 大小: %d 字节\n", t.TotalSize())
    fmt.Printf("📦 分片数: %d\n", t.NumPieces())
    fmt.Printf("🔐 Info Hash: %s\n", must(t.InfoHashHex()))
}

func must(s string, err error) string {
    if err != nil {
        panic(err)
    }
    return s
}
```

### 创建一个 .torrent 文件

```go
package main

import (
    "crypto/sha1"
    "gop2p/pkg/torrent"
    "os"
)

func main() {
    // 模拟文件内容和分片哈希
    fileContent := []byte("Hello, World!")
    
    // 生成分片哈希值（用真实内容的 SHA1）
    var pieces []byte
    hash := sha1.Sum(fileContent)
    pieces = append(pieces, hash[:]...)
    
    // 创建 Torrent 元数据
    t := &torrent.TorrentInfo{
        Announce: "http://tracker.example.com:6969/announce",
        Info: torrent.InfoDict{
            Name:        "hello.txt",
            Length:      int64(len(fileContent)),
            PieceLength: 16384,
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
}
```

---

## 📚 常用方法

| 方法 | 说明 | 返回值 |
|------|------|--------|
| `t.Mode()` | 判断是单/多文件 | `SingleFile` 或 `MultiFile` |
| `t.TotalSize()` | 获取总大小 | `int64` |
| `t.NumPieces()` | 获取分片数 | `int` |
| `t.GetPiece(i)` | 获取第 i 片的哈希 | `[]byte` (20 bytes) |
| `t.GetPieceHex(i)` | 获取第 i 片的十六进制哈希 | `string` (40 chars) |
| `t.InfoHash()` | 获取 Info Hash (二进制) | `[]byte, error` |
| `t.InfoHashHex()` | 获取 Info Hash (十六进制) | `string, error` |
| `t.ValidateInfo()` | 验证 Torrent 数据 | `error` |
| `torrent.ParseTorrent(reader)` | 解析 Torrent 文件 | `*TorrentInfo, error` |
| `torrent.EncodeTorrent(t)` | 编码为 Bencode | `[]byte, error` |

---

## 🔧 常见任务

### 任务 1：获取 Info Hash

```go
hash, _ := t.InfoHashHex()
fmt.Printf("Info Hash: %s\n", hash)
// 输出: Info Hash: a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6
```

### 任务 2：生成磁力链接

```go
hash, _ := t.InfoHashHex()
magnet := fmt.Sprintf("magnet:?xt=urn:btih:%s&dn=%s", hash, t.Info.Name)
```

### 任务 3：列出所有文件

```go
if t.Mode() == torrent.SingleFile {
    fmt.Println(t.Info.Name, "-", t.Info.Length, "bytes")
} else {
    for _, f := range t.Info.Files {
        path := strings.Join(f.Path, "/")
        fmt.Println(path, "-", f.Length, "bytes")
    }
}
```

### 任务 4：验证 Torrent 有效性

```go
if err := t.ValidateInfo(); err != nil {
    fmt.Printf("❌ Torrent 无效: %v\n", err)
} else {
    fmt.Println("✅ Torrent 有效")
}
```

### 任务 5：计算下载时间

```go
totalBytes := t.TotalSize()
bytesPerSecond := int64(1000000)  // 1 Mbps = 125 KB/s
seconds := totalBytes / bytesPerSecond
hours := float64(seconds) / 3600

fmt.Printf("下载需要: %.2f 小时\n", hours)
```

---

## 🏃 完整示例程序

运行以下命令来看实际演示：

### 示例 1：创建和读取单文件 Torrent
```bash
go run examples/read_real_torrent.go
```

输出：
- 创建 `example-file.torrent` 文件
- 显示 Bencode 原始格式
- 解析并显示所有元数据
- 生成磁力链接

### 示例 2：分析任何 Torrent 文件
```bash
# 自动创建演示文件并分析
go run examples/analyze_torrent.go

# 或分析现有文件
go run examples/analyze_torrent.go your-file.torrent
```

输出：
- 跟踪器信息
- 完整的元数据
- 文件列表（带大小）
- 分片哈希值（前几个）
- 下载估算
- 协议信息

### 示例 3：通用 Torrent 解析器
```bash
go run examples/parse_torrent_advanced.go my-file.torrent
```

---

## 📋 Torrent 文件结构

### 示例 Bencode 格式

```
d
  8:announce      40:http://tracker.example.com:6969/announce
  13:announce-list
    l
      l40:http://backup1.example.com:6969/announcee
      l40:http://backup2.example.com:6969/announcee
    e
  7:comment  69:这是一个演示 Torrent 文件...
  10:created by  20:gop2p v1.0 - Example
  13:creation date  i1782958346e
  4:info
    d
      6:length  i272e
      4:name    16:example-file.txt
      12:piece length  i50e
      6:pieces  120:<20-byte-SHA1-hash>...<20-byte-SHA1-hash>
    e
e
```

---

## 🔑 关键概念

### Info Hash
```
Info Hash = SHA1(bencode(info_dict))
```
- 20 字节二进制格式
- 40 个十六进制字符格式
- 唯一标识 Torrent 的内容

### 分片（Pieces）
- 每个分片是一个完整下载单元
- 每个分片有 20 字节的 SHA1 哈希值
- 用于验证下载的数据完整性

### 单文件 vs 多文件

**单文件：**
```
info.length = 1024        // 文件大小
info.pieces = [hash1, hash2, ...]
```

**多文件：**
```
info.files = [
  {length: 512, path: ["dir1", "file1.txt"]},
  {length: 256, path: ["dir2", "file2.txt"]},
]
info.pieces = [hash1, hash2, ...]
```

---

## ⚠️ 常见错误

### "failed to encode info dict: unsupported type"
**原因**：字段类型不支持

**解决**：确保所有字段都是下列类型：
- `string`
- `int64`
- `[]byte`
- `[]interface{}`
- `map[string]interface{}`
- 嵌套的结构体（带 bencode 标签）

### "invalid character in string length"
**原因**：文件不是有效的 Bencode 格式

**解决**：确保文件是真正的 .torrent 文件，未被损坏

### "pieces field must be multiple of 20 bytes"
**原因**：分片哈希值不是 20 字节的倍数

**解决**：确保 `pieces` 字段的长度 % 20 == 0

---

## 📁 生成的文件

运行示例程序后，会生成：

| 文件 | 说明 | 大小 |
|------|------|------|
| `example-file.torrent` | 单文件演示 | ~0.5 KB |
| `demo-multifile.torrent` | 多文件演示 | ~0.5 KB |

这些都是真实有效的 .torrent 文件，可以用任何 BitTorrent 客户端打开。

---

## 下一步

现在你可以：
1. ✅ 读取 .torrent 文件
2. ✅ 创建 .torrent 文件
3. ✅ 获取 Info Hash
4. ✅ 生成磁力链接

接下来学习：
- **Week 4**：Tracker 通信（使用 Info Hash 向 Tracker 发起 announce 请求）
- **Week 5**：DHT 发现（P2P 方式查找 Peers）
- **Week 6**：Peer Wire 协议（实际下载文件）

---

## 更多资源

- [完整指南](TORRENT_READING_GUIDE.md) - 详细的 Torrent 读取指南
- [Week 3 文档](WEEK3_COMPLETE.md) - Torrent 包实现细节
- [Bencode 说明](BINARY_SAFETY_GUIDE.md) - 二进制安全编码
