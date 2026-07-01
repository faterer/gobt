package main

import (
	"encoding/hex"
	"fmt"
	"gop2p/pkg/bencode"
	"strings"
)

func main() {
	fmt.Println("\n" + strings.Repeat("═", 80))
	fmt.Println("🔍 Bencode编码器输出演示 - 证明返回的就是二进制[]byte")
	fmt.Println(strings.Repeat("═", 80))

	// 演示1: 简单字符串
	fmt.Println("\n【演示1】编码字符串: \"hello\"")
	fmt.Println(strings.Repeat("─", 60))

	encoder := bencode.NewEncoder()
	encoded, _ := encoder.EncodeString("hello")

	fmt.Printf("Go代码:    encoder.EncodeString(\"hello\")\n")
	fmt.Printf("返回类型:   []byte (已经是二进制!)\n")
	fmt.Printf("长度:      %d 字节\n", len(encoded))
	fmt.Printf("十六进制:   %s\n", hex.EncodeToString(encoded))
	fmt.Printf("原始值:     %v\n", encoded)
	fmt.Printf("打印输出:   %s\n", string(encoded))
	fmt.Println()
	fmt.Println("解释: \"5:hello\" 中:")
	fmt.Println("  35 = ASCII '5'")
	fmt.Println("  3A = ASCII ':'")
	fmt.Println("  68 = ASCII 'h'")
	fmt.Println("  65 = ASCII 'e'")
	fmt.Println("  6C = ASCII 'l'")
	fmt.Println("  6C = ASCII 'l'")
	fmt.Println("  6F = ASCII 'o'")

	// 演示2: 包含特殊字符的数据
	fmt.Println("\n【演示2】编码二进制数据: bytes \"\\x00\\x01\\x02\"")
	fmt.Println(strings.Repeat("─", 60))

	binaryData := string([]byte{0x00, 0x01, 0x02})
	encoded, _ = encoder.EncodeString(binaryData)

	fmt.Printf("Go代码:    encoder.EncodeString(\"\\\\x00\\\\x01\\\\x02\")\n")
	fmt.Printf("返回类型:   []byte (二进制)\n")
	fmt.Printf("长度:      %d 字节\n", len(encoded))
	fmt.Printf("十六进制:   %s\n", hex.EncodeToString(encoded))
	fmt.Printf("原始值:     %v\n", encoded)
	fmt.Println()
	fmt.Println("解释: \"3:\\x00\\x01\\x02\" 中:")
	fmt.Println("  33 = ASCII '3'")
	fmt.Println("  3A = ASCII ':'")
	fmt.Println("  00 = 原始二进制字节")
	fmt.Println("  01 = 原始二进制字节")
	fmt.Println("  02 = 原始二进制字节")

	// 演示3: 整数
	fmt.Println("\n【演示3】编码整数: 42")
	fmt.Println(strings.Repeat("─", 60))

	encoded, _ = encoder.EncodeInteger(42)

	fmt.Printf("Go代码:    encoder.EncodeInteger(42)\n")
	fmt.Printf("返回类型:   []byte (二进制)\n")
	fmt.Printf("长度:      %d 字节\n", len(encoded))
	fmt.Printf("十六进制:   %s\n", hex.EncodeToString(encoded))
	fmt.Printf("打印输出:   %s\n", string(encoded))
	fmt.Println()
	fmt.Println("解释: \"i42e\" 中:")
	fmt.Println("  69 = ASCII 'i'")
	fmt.Println("  34 = ASCII '4'")
	fmt.Println("  32 = ASCII '2'")
	fmt.Println("  65 = ASCII 'e'")

	// 演示4: 完整的字典（包含SHA1哈希）
	fmt.Println("\n【演示4】编码包含哈希值的字典")
	fmt.Println(strings.Repeat("─", 60))

	// 模拟一个SHA1哈希（20字节的二进制数据）
	fakeHash := string([]byte{
		0xBB, 0x0D, 0xF4, 0xA4, 0x29, 0x9E, 0x7B, 0x04,
		0xF6, 0x00, 0x21, 0x7E, 0x9E, 0x19, 0xE6, 0xEE,
		0x68, 0x77, 0xE7, 0x38,
	})

	dict := map[string]interface{}{
		"name":   "test",
		"pieces": fakeHash,
	}

	encoded, _ = encoder.Encode(dict)

	fmt.Printf("Go代码:    encoder.Encode(map with SHA1 hash)\n")
	fmt.Printf("返回类型:   []byte (二进制)\n")
	fmt.Printf("长度:      %d 字节\n", len(encoded))
	fmt.Printf("十六进制:   %s\n", hex.EncodeToString(encoded))
	fmt.Println()
	fmt.Println("输出中包含:")
	fmt.Println("  • ASCII 可打印字符 (key names, structure)")
	fmt.Println("  • 原始二进制数据 (SHA1 hash)")

	// 演示5: 保存到文件时发生了什么
	fmt.Println("\n【演示5】保存过程")
	fmt.Println(strings.Repeat("─", 60))

	fmt.Println(`
写文件的代码:
  encoded, _ := encoder.Encode(data)
  ioutil.WriteFile("file.torrent", encoded, 0644)

步骤:
  1. encoder.Encode() 返回 []byte
  2. []byte 中包含：
     - ASCII 字符（看起来像文本）
     - 原始二进制数据（真正的二进制）
  3. WriteFile() 把 []byte 直接写到文件
  4. 没有任何转换！

所以最终文件 file.torrent 就是这个 []byte 的原始内容
`)

	// 演示6: 为什么有的部分看起来是文本
	fmt.Println("\n【演示6】为什么Bencode看起来像文本？")
	fmt.Println(strings.Repeat("─", 60))

	fmt.Println(`
原因:
  • Bencode 设计用 ASCII 字符来表示结构信息
  • 'l' = 列表开始，'e' = 列表结束
  • ':' = 分隔符，'i' = 整数开始
  • 所以前面全是 ASCII 字符

例子:
  d8:announce39:http://tracker.ubuntu.com:6969/announce ...
  ↑           ↑ ↑ ↑
  'd'        '8' ':' 字符串长度和URL都是ASCII
  字典        所有这些都可以作为文本查看
  开始

但是:
  ...6:pieces100:[二进制SHA1哈希]
                 ↑
                 这是原始二进制，无法显示为文本

所以 .torrent 文件:
  ✓ 前面看起来像文本（因为确实是ASCII）
  ✓ 后面是二进制（pieces的SHA1哈希）
  ✓ 整体仍然是二进制文件
`)

	fmt.Println("\n" + strings.Repeat("═", 80))
	fmt.Println("结论: Encoder 返回的 []byte 就已经是二进制了！")
	fmt.Println("      没有字符串→二进制的转换，因为本来就是二进制")
	fmt.Println(strings.Repeat("═", 80) + "\n")
}
