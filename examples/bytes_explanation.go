package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("\n" + strings.Repeat("═", 90))
	fmt.Println("📊 Go中[]byte的本质 - 为什么Encoder输出就已经是二进制")
	fmt.Println(strings.Repeat("═", 90))

	// 演示：[] byte并不是"字符串编码"，而就是原始字节
	fmt.Println("\n【概念1】[]byte就是一个字节数组，不是字符串")
	fmt.Println(strings.Repeat("─", 90))

	fmt.Println(`
在Go中:
  • string = "不可变"的UTF-8字符序列（概念层面）
  • []byte = 可变的原始字节序列（低层面）
  
关键区别:
  ┌─────────────────────────────────────────────────────┐
  │ string "hello"                                      │
  │ ↓ 转换为 []byte                                     │
  │ []byte{104, 101, 108, 108, 111}  // 原始字节序列    │
  │ ↑ 这就是在内存中的样子！                            │
  └─────────────────────────────────────────────────────┘
`)

	// 具体演示
	fmt.Println("\n【演示】具体例子")
	fmt.Println(strings.Repeat("─", 90))

	// 方式1：直接创建[]byte
	fmt.Println("\n方式1: 直接创建[]byte（可以包含任意字节）")
	bytes1 := []byte{0x68, 0x65, 0x6C, 0x6C, 0x6F}  // 'hello'
	fmt.Printf("  []byte{0x68, 0x65, 0x6C, 0x6C, 0x6F}\n")
	fmt.Printf("  = %v\n", bytes1)
	fmt.Printf("  作为字符串显示: %s\n", bytes1)
	fmt.Printf("  十六进制: %x\n\n", bytes1)

	// 方式2：字符串转[]byte
	fmt.Println("方式2: 字符串转[]byte")
	str := "hello"
	bytes2 := []byte(str)
	fmt.Printf("  []byte(\"hello\")\n")
	fmt.Printf("  = %v\n", bytes2)
	fmt.Printf("  十六进制: %x\n\n", bytes2)

	// 方式3：包含非ASCII字节
	fmt.Println("方式3: []byte可以包含任意字节值")
	bytes3 := []byte{0xFF, 0xAB, 0xCD, 0x00, 0x01}  // 无法显示为文本
	fmt.Printf("  []byte{0xFF, 0xAB, 0xCD, 0x00, 0x01}\n")
	fmt.Printf("  = %v\n", bytes3)
	fmt.Printf("  作为字符串显示: %s (乱码)\n", bytes3)
	fmt.Printf("  十六进制: %x\n", bytes3)

	// 现在解释Encoder的逻辑
	fmt.Println("\n" + strings.Repeat("═", 90))
	fmt.Println("【关键】Bencode Encoder的逻辑")
	fmt.Println(strings.Repeat("═", 90))

	fmt.Println(`
Encoder.Encode() 做的事情:

代码:
  func (e *Encoder) Encode(v interface{}) ([]byte, error) {
    e.buf.Reset()
    err := e.encode(v)
    if err != nil {
      return nil, err
    }
    return e.buf.Bytes(), nil  // ← 返回[]byte!
  }

过程:
  1. 创建一个 bytes.Buffer （内部存储）
  2. 写入字节到buffer:
     - WriteRune('d')      ← 写入ASCII 'd' (0x64)
     - WriteString("key")  ← 写入ASCII 'k','e','y' 
     - WriteRune('e')      ← 写入ASCII 'e' (0x65)
     - 也可以写入原始[]byte（如SHA1哈希）
  3. 返回 buffer.Bytes() 作为 []byte

关键点:
  ✓ WriteRune() 写入ASCII字符 → 看起来像文本
  ✓ WriteString() 写入字符串 → 看起来像文本
  ✓ buffer.Write() 写入[]byte → 原始二进制
  ✓ 最后返回的 []byte → 上述的混合体
`)

	// 视觉化展示完整过程
	fmt.Println("\n" + strings.Repeat("═", 90))
	fmt.Println("【完整过程】从Go数据到.torrent文件")
	fmt.Println(strings.Repeat("═", 90))

	fmt.Println(`
步骤1: Go数据结构
┌─────────────────────────────────────────────────────┐
│ map[string]interface{}{                             │
│   "name": "test",                                   │
│   "pieces": [20字节的SHA1哈希]                      │
│ }                                                   │
└─────────────────────────────────────────────────────┘
                        ↓
            encoder.Encode(data)

步骤2: Encoder构建[]byte
┌─────────────────────────────────────────────────────┐
│ bytes.Buffer (可变的字节数组)                        │
│                                                     │
│ WriteRune('d')        → [0x64]                      │
│ WriteString("name")   → [0x64, 0x6E, 0x61, 0x6D...]│
│ WriteString("test")   → [...更多ASCII字符...]       │
│ WriteString("pieces") → [...更多ASCII字符...]       │
│ Write([20字节哈希])   → [...混合ASCII和二进制...]   │
│ WriteRune('e')        → [..., 0x65]                 │
│                                                     │
│ → 最后调用 buffer.Bytes()                          │
│ → 返回 []byte                                       │
└─────────────────────────────────────────────────────┘
                        ↓
           []byte (混合ASCII和二进制)
           示例: [0x64, 0x6E... 0xFF 0xAB... 0x65]

步骤3: 保存到文件
┌─────────────────────────────────────────────────────┐
│ ioutil.WriteFile("file.torrent", encoded, 0644)    │
│                                                     │
│ 这里的 encoded 就是上面的 []byte                   │
│ WriteFile直接把[]byte写到文件                       │
│ 没有转换！就是原始字节                              │
└─────────────────────────────────────────────────────┘
                        ↓
           file.torrent (二进制文件)
           内容 = 完全相同的[]byte内容
`)

	// 关键悟点
	fmt.Println("\n" + strings.Repeat("═", 90))
	fmt.Println("💡 关键悟点")
	fmt.Println(strings.Repeat("═", 90))

	fmt.Println(`
❌ 错误理解:
  "Encoder输出字符串" → "WriteFile转换为二进制"
  这是错的！

✅ 正确理解:
  Encoder从一开始就输出[]byte（二进制）
  ↓
  这个[]byte中：
    • 有ASCII字符（看起来像文本）
    • 有原始二进制数据（SHA1等）
  ↓
  WriteFile把这个[]byte直接写到文件
  ↓
  完成！

类比:
  ┌─────────────────────────────────────────────┐
  │ 你有一个"混合包"：                          │
  │  • 前部分：普通信件（ASCII可见）            │
  │  • 后部分：密封的礼物（二进制不可见）       │
  │                                             │
  │ 寄邮件时，邮递员不区分：                    │
  │  他就是把整个包（按原样）寄出去             │
  │  不会"转换"任何东西                         │
  │                                             │
  │ 同样，WriteFile也不做转换：                │
  │  就是把[]byte按原样写到文件                │
  └─────────────────────────────────────────────┘
`)

	// 代码展示
	fmt.Println("\n" + strings.Repeat("═", 90))
	fmt.Println("【代码】Encoder的简化版本")
	fmt.Println(strings.Repeat("═", 90))

	fmt.Println(`
type SimpleEncoder struct {
  buf bytes.Buffer
}

func (e *SimpleEncoder) Encode(v interface{}) []byte {
  e.buf.Reset()
  
  // 写入ASCII结构信息
  e.buf.WriteRune('d')            // 字典开始
  e.buf.WriteString("pieces")     // key名
  e.buf.WriteString("20:")         // value长度
  
  // 写入原始二进制数据
  e.buf.Write(shaHash)            // 20字节的SHA1
  
  e.buf.WriteRune('e')            // 字典结束
  
  // 关键：直接返回[]byte
  // buffer.Bytes() 返回的就是[]byte
  // 包含：ASCII字符 + 二进制数据 的混合
  return e.buf.Bytes()
}

// 使用
data := map[string]interface{}{"pieces": sha1Hash}
encoded := encoder.Encode(data)   // 返回[]byte

// 保存
ioutil.WriteFile("file.torrent", encoded, 0644)
// 这里没有任何转换，encoded 的原始内容
// 直接写到文件
`)

	fmt.Println("\n" + strings.Repeat("═", 90))
	fmt.Println("✨ 结论")
	fmt.Println(strings.Repeat("═", 90))

	fmt.Println(`
1. Encoder.Encode() 返回 []byte
   → []byte 就是一个字节数组，不是字符串

2. 这个[]byte包含：
   → ASCII字符（看起来像文本）
   → 原始二进制（SHA1等）

3. 保存到文件时：
   → WriteFile直接把[]byte写到文件
   → 没有转换！

4. 为什么有人会困惑：
   → 因为前面的ASCII部分可以显示为文本
   → 但整体仍然是二进制文件
   → .torrent文件根本不能用文本编辑器打开

5. 这就是Bencode的聪明之处：
   → 用ASCII字符表示结构
   → 但整体仍然是紧凑的二进制格式
   → 同时又支持任意二进制数据
`)

	fmt.Println(strings.Repeat("═", 90) + "\n")
}
