# Examples - Bencode 和 Torrent 文件示例

## 📝 示例列表

### 1. Bencode 编码/解码 (`bencode_simple.go`)

演示如何使用 Bencode 库进行编码和解码：

```bash
# 运行
go run bencode_simple.go

# 输出示例
# === Bencode 编码示例 ===
#
# 1. 编码字符串
#    encoder.EncodeString("hello") => 5:hello
#
# 2. 编码整数
#    encoder.EncodeInteger(42) => i42e
# ...
```

**学到的内容：**
- 字符串编码：`5:hello` (5个字符的"hello")
- 整数编码：`i42e` (整数 42)
- 列表编码：`l...e` (Bencode 列表)
- 字典编码：`d...e` (Bencode 字典)
- 解码 Bencode 数据

---

### 2. Torrent 文件解析 (`parse_torrent.go`)

演示如何读取和解析 `.torrent` 文件：

```bash
# 运行（会自动在 examples 和 ../（根目录）查找 .torrent 文件）
go run parse_torrent.go init.go

# 输出示例显示：
# - Torrent 基本信息（名称、大小、分片数）
# - Tracker 信息（announce 地址）
# - 元数据（创建时间、创建者）
# - Info Hash（SHA1）
# - 文件列表（多文件模式）
# - 统计信息
```

**学到的内容：**
- Torrent 文件格式
- Info Hash 计算（SHA1）
- Tracker 地址处理
- 分片管理
- 元数据提取

---

## 📂 可用的 .torrent 文件

### 本地文件（examples 目录）
- **example-demo.torrent** - 简单演示文件（150 B，3 个分片）

### 从父目录加载
当从 examples 目录运行时，也会找到：
- **example-file.torrent** - 单文件模式（根目录）
- **demo-multifile.torrent** - 多文件模式（根目录）

---

## 🚀 快速开始

### 1. 编码和解码 Bencode

```bash
cd examples
go run bencode_simple.go
```

输出会显示各种数据类型的编码方式。

### 2. 解析 Torrent 文件

```bash
cd examples
go run parse_torrent.go init.go
```

会自动解析所有找到的 `.torrent` 文件。

---

## 💡 关键概念

### Bencode 格式

- **字符串**：`<长度>:<数据>` 例：`5:hello`
- **整数**：`i<数字>e` 例：`i42e`
- **列表**：`l<元素>e` 例：`l1:a1:be`
- **字典**：`d<键><值>e` 例：`d3:agei25ee`

### Torrent 文件结构

```
{
  "announce": "http://tracker.example.com:6969/announce",
  "info": {
    "name": "filename.txt",
    "length": 12345,      # 文件大小
    "piece length": 16384,
    "pieces": "<SHA1_0><SHA1_1>..."  # 分片哈希值
  },
  "creation date": 1234567890
}
```

**Info Hash** = SHA1(bencode(info_dict))

---

## 📚 更多学习资源

- 查看根目录的 `TORRENT_READING_GUIDE.md` - 详细完整指南
- 查看根目录的 `TORRENT_QUICK_REFERENCE.md` - 快速参考
- 查看 `pkg/bencode/` 包的源代码
- 查看 `pkg/torrent/` 包的源代码
