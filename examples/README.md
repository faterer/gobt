# Bencode Torrent Parser - 使用示例

这个示例展示了如何使用我们实现的 Bencode 编解码器来创建和解析真实的 .torrent 文件。

## 📁 文件结构

```
examples/
├── create_torrent.go      # 生成示例torrent文件
├── parse_torrent.go       # 解析并显示torrent文件信息
├── ubuntu-22.04.torrent   # 单文件torrent示例（自动生成）
└── sample-collection.torrent  # 多文件torrent示例（自动生成）
```

## 🚀 快速开始

### 步骤1: 生成示例Torrent文件

```bash
cd examples
go run create_torrent.go
```

这会创建两个示例torrent文件：

1. **ubuntu-22.04.torrent** - 单文件模式
   - 模拟Ubuntu ISO文件 (3.52 GB)
   - 包含多层Tracker列表
   - 包含元数据信息

2. **sample-collection.torrent** - 多文件模式
   - 包含3个文件
   - 总大小 3.57 GB
   - 展示嵌套的文件列表结构

### 步骤2: 解析Torrent文件

```bash
go run parse_torrent.go
```

程序会自动扫描目录中的所有 `.torrent` 文件，逐个解析并显示详细信息。

## 📊 输出示例

### 单文件Torrent (ubuntu-22.04.torrent)

```
🚀 Bencode Torrent Parser Example
============================================================

📂 Parsing: ubuntu-22.04.torrent

================================================================================
📋 TORRENT FILE INFORMATION
================================================================================

📌 Basic Information:
  Name:              ubuntu-22.04-desktop-amd64.iso
  Total Size:        3.28 GB
  Piece Length:      256.00 KB
  Number of Pieces:  5

📡 Tracker Information:
  Announce:          http://torrent.ubuntu.com:6969/announce
  Announce List:     (2 tiers)
    Tier 1:
      - http://torrent.ubuntu.com:6969/announce
      - http://ipv6.torrent.ubuntu.com:6969/announce
    Tier 2:
      - http://torrent.ubuntulinux.nl:6969/announce

📝 Metadata:
  Comment:           Official Ubuntu 22.04 LTS Desktop ISO
  Creation Date:     2026-07-01 15:45:34
  Created By:        gop2p/4.2.0

🔐 Info Hash:
  d7c96f1308375f4a873db6e5414580b90864e1f7

📊 Statistics:
  Total Data Size:   3.28 GB
  Expected Size:     1.25 MB
```

### 多文件Torrent (sample-collection.torrent)

```
📂 Parsing: sample-collection.torrent

📋 TORRENT FILE INFORMATION

📌 Basic Information:
  Name:              sample-collection
  Total Size:        3.49 GB
  Piece Length:      256.00 KB
  Number of Pieces:  5

📡 Tracker Information:
  Announce:          http://tracker.example.com:6969/announce

📁 Files (3):
  1. images/image1.iso (1.00 GB)
  2. images/image2.iso (2.00 GB)
  3. documents/README.txt (500.00 MB)

📝 Metadata:
  Comment:           Multi-file torrent example
  Creation Date:     2026-07-01 15:45:34
  Created By:        gop2p/4.2.0

🔐 Info Hash:
  4305e9990fc0961e513024a0089f5a38f20ceab1
```

## 🔍 深入理解

### Bencode格式示例

生成的torrent文件使用Bencode编码。原始内容看起来像：

```
d8:announce39:http://torrent.ubuntu.com:6969/announce
13:announce-listll39:http://torrent.ubuntu.com:6969/announce
44:http://ipv6.torrent.ubuntu.com:6969/announceel43:http://torrent.ubuntulinux.nl:6969/announceee
7:comment37:Official Ubuntu 22.04 LTS Desktop ISO
10:created by11:gop2p/4.2.0
13:creation datei1782891934e
4:infod
  6:lengthi3520000000e
  4:name30:ubuntu-22.04-desktop-amd64.iso
  12:piece lengthi262144e
  6:pieces100:[20字节 x 5]
e
```

### 关键概念

1. **Info Hash**
   - SHA1(bencode(info部分))
   - 用于唯一标识torrent
   - 在BitTorrent协议中必须的

2. **Announce**
   - Tracker服务器地址
   - 用于报告下载进度
   - 支持多层备用列表

3. **Pieces**
   - 文件分块的SHA1哈希
   - 每个哈希20字节
   - 用于验证数据完整性

4. **单文件 vs 多文件**
   - 单文件：info中有"length"字段
   - 多文件：info中有"files"字段（嵌套列表）

## 🧪 完整工作流程

```
1. 创建数据结构 (Go maps)
   ↓
2. 使用Bencode编码器编码为字节
   ↓
3. 保存为.torrent文件
   ↓
4. 使用Bencode解码器读取文件
   ↓
5. 计算Info Hash (SHA1)
   ↓
6. 显示解析结果
```

## 💡 学习要点

### 我们实现的Bencode编解码做了什么？

1. **编码器** (`encoder.go`)
   - 支持4种数据类型：int64, string, list, dict
   - 自动处理嵌套结构
   - 字典key自动排序（Bencode规范要求）

2. **解码器** (`decoder.go`)
   - 递归解析复杂结构
   - 完整的错误处理
   - 支持UTF-8字符串

### 为什么这很重要？

- BitTorrent协议**必须**使用Bencode
- 所有.torrent文件都用Bencode编码
- Tracker通信使用Bencode
- 了解编解码是理解P2P协议的基础

## 📝 修改示例

如果想修改生成的torrent文件，编辑 `create_torrent.go`：

```go
// 修改文件大小
"length": int64(5368709120), // 5 GB

// 修改piece大小
"piece length": int64(131072), // 128 KB

// 添加tracker
"announce-list": []interface{}{
    []interface{}{
        "http://new-tracker.example.com:6969/announce",
    },
}

// 修改文件列表（多文件模式）
"files": []interface{}{
    map[string]interface{}{
        "length": int64(1000000000),
        "path": []interface{}{"myfile.iso"},
    },
}
```

## 🔧 调试技巧

### 查看原始Bencode

```bash
# 查看torrent文件的原始Bencode内容
hexdump -C ubuntu-22.04.torrent | head -20
```

### 比较两个torrent文件

```bash
# 检查Info Hash是否相同（这决定了torrent的唯一性）
cd examples
go run parse_torrent.go | grep "Info Hash"
```

## 🚀 下一步

这个示例展示了如何：
1. ✅ 创建Bencode格式的数据
2. ✅ 编码为二进制.torrent文件
3. ✅ 解析torrent文件
4. ✅ 计算Info Hash
5. ✅ 显示torrent的各种信息

在实际的BitTorrent客户端中，你会：
- 连接到Tracker报告下载进度
- 从其他Peer下载文件块
- 验证块的SHA1哈希
- 管理下载队列和带宽

## 📖 参考资源

- BitTorrent规范: https://www.bittorrent.org/beps/bep_0003.html
- Bencode规范: 在规范中的第二部分
- Info Hash计算: 必须对info字典进行Bencode编码后SHA1哈希

---

**提示**: 这些示例使用虚拟的文件数据和哈希值。在实际应用中，你会使用真实文件的SHA1哈希。
