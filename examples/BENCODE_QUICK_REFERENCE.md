# Bencode 快速参考卡

## 核心公式

```
Bencode = 结构（ASCII） + 数据（二进制）
        = 紧凑 + 安全 + 高效
```

---

## 四种数据类型

### 1. 整数（Integer）
```
格式: i<number>e

例子:
  42     → i42e
  -273   → i-273e
  0      → i0e

在内存中:
  i42e = [0x69, 0x34, 0x32, 0x65]
       = ['i', '4', '2', 'e']
```

### 2. 字符串（String）- 最重要的一个
```
格式: <length>:<data>

例子:
  "hello"           → 5:hello
  特殊字符           → 11:hello:world
  纯二进制          → 5:[0xFF,0xAB,...]
  SHA1哈希          → 20:[20字节]

关键特性: 长度前缀！
  ✓ 支持任意字节
  ✓ 支持特殊字符
  ✓ 支持纯二进制
```

### 3. 列表（List）
```
格式: l<items>e

例子:
  [1, "hello"]  → li1e5:helloe
  嵌套列表      → lli1ei2ee
```

### 4. 字典（Dictionary）
```
格式: d<key><value>...e
      （keys必须按字母顺序排序！）

例子:
  {"name": "test"}        → d4:name4:teste
  {"hash": [20字节]}      → d4:hash20:[20字节]e
  {"age": 27, "name": "Bob"} → d3:agei27e4:name3:Bobee
                                 (age排在name前面)
```

---

## Encoder 工作流程

```
Go数据结构
    ↓
bytes.Buffer (内存缓冲区)
    ├─ WriteRune('d')       → [0x64]
    ├─ WriteString("key")   → [0x6B, 0x65, 0x79]
    ├─ WriteRune(':')       → [0x3A]
    ├─ WriteString("value") → [...]
    ├─ Write(binaryData)    → [0xFF, 0xAB, ...]
    └─ WriteRune('e')       → [0x65]
    ↓
[]byte (混合ASCII和二进制)
    ↓
WriteFile() 直接写到文件
    ↓
.torrent 文件
```

---

## Decoder 工作流程

```
读取 <长度>
    ↓
解析长度为整数
    ↓
io.ReadFull(reader, buffer)
    ├─ buffer大小 = 长度
    ├─ 读指定数量的字节
    └─ 不检查字节值！
    ↓
返回数据 (可以是任意字节)
```

---

## 关键代码片段

### Encoder（核心）
```go
// 全是字节操作
buf.WriteRune('d')       // 字符 → ASCII码
buf.WriteString("key")   // 字符串 → ASCII字节
buf.Write(binaryData)    // 二进制 → 直接写
return buf.Bytes()       // 返回[]byte
```

### Decoder（核心）
```go
// 读长度
length := parseLength()  // 如"20"→20

// 魔法！根据长度读字节
io.ReadFull(reader, make([]byte, length))
// 不管内容是什么都照样读
```

---

## 例子：.torrent文件的pieces字段

### 编码
```go
pieces := []string{
    sha1_hash_1,  // 20字节
    sha1_hash_2,  // 20字节
    ...
}

// Encoder处理
d4:hash1000:...e
    └─ 4:hash    = "hash"
    └─ 1000:     = 长度1000（50个×20）
    └─ [1000字节] = 全是二进制SHA1
```

### 解析
```go
// 看到"1000:"
length := 1000

// io.ReadFull读1000个字节
// 得到50个SHA1哈希
// 完成！
```

---

## 为什么Bencode完美

| 特性 | Bencode | JSON | XML |
|------|---------|------|-----|
| 二进制安全 | ✅ | ❌ | ❌ |
| 紧凑度 | ✅ | △ | ❌ |
| 无歧义 | ✅ | △ | △ |
| 实现简单 | ✅ | △ | ❌ |

---

## 三句话总结

1. **格式** - `<长度>:<数据>` 确定边界
2. **原理** - 字符→ASCII码，二进制→直接写
3. **好处** - 安全紧凑，无需转义

---

## 常见疑惑解答

### Q: "为什么torrent是二进制文件？"
A: 因为它包含二进制数据（SHA1哈希）。虽然前面有ASCII结构，但整体是二进制。

### Q: "为什么不用JSON?"
A: JSON无法安全表示二进制，需要base64编码导致体积增大。

### Q: "长度前缀怎么处理的？"
A: Decoder读长度数字，直到看到':'，然后根据长度读字节。

### Q: "如果数据中有':'怎么办？"
A: 没关系！长度已经告诉解析器要读多少字节，数据中的所有字符都被看作数据。

---

## 学习资源

- [BINARY_SAFETY_GUIDE.md](BINARY_SAFETY_GUIDE.md) - 完整技术指南
- [bencode_output_demo.go](bencode_output_demo.go) - 十六进制演示
- [decoder_binary_demo.go](decoder_binary_demo.go) - 解码器演示
- [pkg/bencode/encoder.go](../pkg/bencode/encoder.go) - 实现
- [pkg/bencode/decoder.go](../pkg/bencode/decoder.go) - 实现

---

**记住：Bencode的核心秘密就是长度前缀！** 🔑
