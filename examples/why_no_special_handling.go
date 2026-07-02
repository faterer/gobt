package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"gop2p/pkg/bencode"
	"strings"
)

func main() {
	fmt.Println("\n" + strings.Repeat("═", 90))
	fmt.Println("🔍 DecodeString() - 为什么不需要特殊处理二进制")
	fmt.Println(strings.Repeat("═", 90))

	// 演示1: ASCII字符串
	fmt.Println("\n【演示1】ASCII字符串")
	fmt.Println(strings.Repeat("─", 90))

	asciiData := []byte("5:hello")
	decoder := bencode.NewDecoder(bytes.NewReader(asciiData))
	result, _ := decoder.DecodeString()

	fmt.Printf("输入:        %s\n", asciiData)
	fmt.Printf("解析过程:\n")
	fmt.Printf("  1. 读长度   → '5'\n")
	fmt.Printf("  2. 遇到':'  → 停止\n")
	fmt.Printf("  3. 长度=5   → make([]byte, 5)\n")
	fmt.Printf("  4. io.ReadFull() 读5个字节\n")
	fmt.Printf("  5. 得到: %s\n", result)
	fmt.Printf("\n结果: %s (字符串)\n", result)

	// 演示2: 包含特殊字符的字符串
	fmt.Println("\n【演示2】包含特殊字符的字符串")
	fmt.Println(strings.Repeat("─", 90))

	specialData := []byte("15:hello:eee5e:d:")
	decoder = bencode.NewDecoder(bytes.NewReader(specialData))
	result, _ = decoder.DecodeString()

	fmt.Printf("输入:        %s\n", specialData)
	fmt.Printf("解析过程:\n")
	fmt.Printf("  1. 读长度   → '1' '5' = 15\n")
	fmt.Printf("  2. 遇到':'  → 停止\n")
	fmt.Printf("  3. 长度=15  → make([]byte, 15)\n")
	fmt.Printf("  4. io.ReadFull() 读15个字节\n")
	fmt.Printf("     └─ 这15个字节中有 ':' 'e' '5' 'd' 等\n")
	fmt.Printf("     └─ 都被当作数据读进来！\n")
	fmt.Printf("  5. 得到: %s\n", result)
	fmt.Printf("\n结果: %s (即使包含特殊字符)\n", result)

	// 演示3: 纯二进制
	fmt.Println("\n【演示3】纯二进制数据")
	fmt.Println(strings.Repeat("─", 90))

	binaryData := []byte{0x00, 0xFF, 0xAB, 0xCD, 0x12}
	bencodedBinary := append([]byte("5:"), binaryData...)

	decoder = bencode.NewDecoder(bytes.NewReader(bencodedBinary))
	result, _ = decoder.DecodeString()

	fmt.Printf("原始二进制:  %v\n", binaryData)
	fmt.Printf("十六进制:    %s\n", hex.EncodeToString(binaryData))
	fmt.Printf("Bencode:    5:<5字节>\n")
	fmt.Printf("完整十六进制: %s\n", hex.EncodeToString(bencodedBinary))
	fmt.Printf("\n解析过程:\n")
	fmt.Printf("  1. 读长度   → '5'\n")
	fmt.Printf("  2. 遇到':'  → 停止\n")
	fmt.Printf("  3. 长度=5   → make([]byte, 5)\n")
	fmt.Printf("  4. io.ReadFull() 读5个字节\n")
	fmt.Printf("     └─ [0x00, 0xFF, 0xAB, 0xCD, 0x12]\n")
	fmt.Printf("     └─ 不检查这些字节是什么！\n")
	fmt.Printf("  5. 得到: %v\n\n", []byte(result))

	// 演示4: SHA1哈希
	fmt.Println("【演示4】SHA1哈希（20字节二进制）")
	fmt.Println(strings.Repeat("─", 90))

	sha1Hash := []byte{0xBB, 0x0D, 0xF4, 0xA4, 0x29, 0x9E, 0x7B, 0x04,
		0xF6, 0x00, 0x21, 0x7E, 0x9E, 0x19, 0xE6, 0xEE,
		0x68, 0x77, 0xE7, 0x38}
	bencodedHash := append([]byte("20:"), sha1Hash...)

	decoder = bencode.NewDecoder(bytes.NewReader(bencodedHash))
	result, _ = decoder.DecodeString()

	fmt.Printf("SHA1哈希:    %v\n", sha1Hash)
	fmt.Printf("十六进制:    %s\n", hex.EncodeToString(sha1Hash))
	fmt.Printf("Bencode:    20:<20字节>\n")
	fmt.Printf("\n解析过程:\n")
	fmt.Printf("  1. 读长度   → '2' '0' = 20\n")
	fmt.Printf("  2. 遇到':'  → 停止\n")
	fmt.Printf("  3. 长度=20  → make([]byte, 20)\n")
	fmt.Printf("  4. io.ReadFull() 读20个字节\n")
	fmt.Printf("  5. 得到完整的SHA1哈希\n")
	fmt.Printf("\n结果（十六进制）: %s\n", hex.EncodeToString([]byte(result)))

	// 演示5: 为什么不需要特殊处理
	fmt.Println("\n" + strings.Repeat("═", 90))
	fmt.Println("💡 为什么不需要特殊处理二进制")
	fmt.Println(strings.Repeat("═", 90))

	fmt.Println(`
DecodeString() 的代码:

    func (d *Decoder) DecodeString() (string, error) {
        // 第1步: 读长度
        length := readLength()  // 如 "20"
        
        // 第2步: 创建缓冲区
        strBytes := make([]byte, length)  // 大小=长度
        
        // 第3步: 关键！
        io.ReadFull(d.r, strBytes)        // ← 读指定数量的字节
        
        return string(strBytes), nil
    }

关键分析:
  ✓ io.ReadFull(reader, buffer) 的行为:
    - 接收 reader （数据源）
    - 接收 buffer （[]byte缓冲区）
    - 从reader读取 len(buffer) 个字节
    - 填充到 buffer 中
    - 返回读到的字节数

  ✓ 重点：不检查字节值！
    - 无论字节是什么都照样读
    - 可以是 0x00（null字节）
    - 可以是 0xFF（最大值）
    - 可以是 0x3A（':'字符）
    - 一切都可以！

  ✓ 为什么不需要特殊处理：
    - 长度已经告诉我们要读多少
    - io.ReadFull() 就能处理任意字节
    - ASCII、二进制、混合都一样
    - 无需特殊代码！

二进制安全 = 长度前缀 + io.ReadFull()
这两样就够了！
`)

	// 演示6: 对比其他方式
	fmt.Println("\n【演示6】对比：为什么长度前缀很重要")
	fmt.Println(strings.Repeat("─", 90))

	fmt.Println(`
其他协议的问题:

方式1: 用终止符（如\0）
    字符串: hello\0
    问题: 如果字符串中本身有\0怎么办？
    结果: 无法区分！

方式2: 用特殊字符标记结束（如'\n'或'\r'）
    字符串: hello\n
    问题: 如果数据中有换行符呢？
    结果: 无法区分！

方式3: 用转义（如JSON）
    字符串: "hello\"world"
    问题: 需要转义，复杂低效
    结果: 体积增大，解析复杂

方式4: Bencode长度前缀 ✅
    字符串: 20:hello:e:d:任意内容
    优点: 长度告诉解析器要读多少
    结果: 无需转义，完全二进制安全！
`)

	fmt.Println(strings.Repeat("═", 90))
	fmt.Println("✨ 结论")
	fmt.Println(strings.Repeat("═", 90))

	fmt.Println(`
DecodeString() 就是 Bencode 的全部秘密：

1. 不需要特殊的"二进制模式"
2. 不需要检测字节值
3. 不需要转义处理
4. 只需要：
   ✓ 读长度
   ✓ make([]byte, length)
   ✓ io.ReadFull() 读字节
   
就这样！

结果：
  • 支持任意ASCII字符串 ✅
  • 支持特殊字符 ✅
  • 支持纯二进制数据 ✅
  • 支持混合内容 ✅
  • 无需转义 ✅
  • 最小体积 ✅
  • 最快速度 ✅
`)

	fmt.Println(strings.Repeat("═", 90) + "\n")
}
