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
	fmt.Println("🔧 我们的Bencode解码器如何处理二进制数据")
	fmt.Println(strings.Repeat("═", 90))

	// 演示1: 解析简单字符串
	fmt.Println("\n【演示1】解析简单ASCII字符串: 5:hello")
	fmt.Println(strings.Repeat("─", 90))

	data1 := []byte("5:hello")
	decoder := bencode.NewDecoder(bytes.NewReader(data1))
	result, _ := decoder.DecodeString()

	fmt.Printf("输入字节:    %v\n", data1)
	fmt.Printf("十六进制:    %s\n", hex.EncodeToString(data1))
	fmt.Printf("解析过程:\n")
	fmt.Printf("  1. 读'5'(0x35) → 长度数字\n")
	fmt.Printf("  2. 读':'(0x3A) → 分隔符，停止\n")
	fmt.Printf("  3. 长度=5，准备读5个字节\n")
	fmt.Printf("  4. io.ReadFull() 读5个字节: %v\n", []byte("hello"))
	fmt.Printf("结果:        %s\n", result)

	// 演示2: 解析包含特殊字符的字符串
	fmt.Println("\n【演示2】解析包含特殊字符的字符串: 10:hello:eee5")
	fmt.Println(strings.Repeat("─", 90))

	data2 := []byte("10:hello:eee5")
	decoder = bencode.NewDecoder(bytes.NewReader(data2))
	result, _ = decoder.DecodeString()

	fmt.Printf("输入字节:    %v\n", data2)
	fmt.Printf("十六进制:    %s\n", hex.EncodeToString(data2))
	fmt.Printf("解析过程:\n")
	fmt.Printf("  1. 读'1','0'(0x31,0x30) → 长度数字\n")
	fmt.Printf("  2. 读':'(0x3A) → 分隔符，停止\n")
	fmt.Printf("  3. 长度=10，准备读10个字节\n")
	fmt.Printf("  4. io.ReadFull() 读10个字节\n")
	fmt.Printf("    └─ 即使中间有':' 'e' '5' 也照样读\n")
	fmt.Printf("    └─ 因为长度告诉我们就读这么多\n")
	fmt.Printf("结果:        %s\n", result)
	fmt.Printf("包含特殊字符: ':' 'e' '5' 都被作为数据读进来了\n")

	// 演示3: 解析纯二进制
	fmt.Println("\n【演示3】解析纯二进制数据: 5:\"\\x00\\xff\\xab\\xcd\\x12\"")
	fmt.Println(strings.Repeat("─", 90))

	binaryData := []byte{0x00, 0xFF, 0xAB, 0xCD, 0x12}
	bencodedBinary := append([]byte("5:"), binaryData...)

	decoder = bencode.NewDecoder(bytes.NewReader(bencodedBinary))
	result, _ = decoder.DecodeString()

	fmt.Printf("原始二进制数据: %v\n", binaryData)
	fmt.Printf("十六进制:      %s\n", hex.EncodeToString(binaryData))
	fmt.Printf("Bencode编码:   5:<5个二进制字节>\n")
	fmt.Printf("编码后十六进制: %s\n", hex.EncodeToString(bencodedBinary))
	fmt.Printf("\n解析过程:\n")
	fmt.Printf("  1. 读'5'(0x35) → 长度数字\n")
	fmt.Printf("  2. 读':'(0x3A) → 分隔符，停止\n")
	fmt.Printf("  3. 长度=5，准备读5个字节\n")
	fmt.Printf("  4. io.ReadFull() 读5个字节\n")
	fmt.Printf("    └─ 不关心这5个字节是什么\n")
	fmt.Printf("    └─ 二进制、ASCII、乱码都没关系\n")
	fmt.Printf("    └─ 因为长度告诉解析器\"就读这么多\"\n")
	fmt.Printf("\n解析结果（十六进制）: %s\n", hex.EncodeToString([]byte(result)))
	fmt.Printf("解析结果（原始值）:  %v\n", []byte(result))

	// 演示4: 解析完整的torrent结构
	fmt.Println("\n【演示4】解析包含SHA1哈希的字典")
	fmt.Println(strings.Repeat("─", 90))

	// 创建一个包含hash的数据
	fakeHash := []byte{0xBB, 0x0D, 0xF4, 0xA4, 0x29, 0x9E, 0x7B, 0x04,
		0xF6, 0x00, 0x21, 0x7E, 0x9E, 0x19, 0xE6, 0xEE,
		0x68, 0x77, 0xE7, 0x38} // 20字节SHA1

	// 手动构建Bencode: d4:hash20:<hash>e
	bencodedDict := []byte("d4:hash20:")
	bencodedDict = append(bencodedDict, fakeHash...)
	bencodedDict = append(bencodedDict, 'e')

	fmt.Printf("结构: {\"hash\": <20字节SHA1>}\n")
	fmt.Printf("Bencode编码:\n")
	fmt.Printf("  d                  ← 字典开始\n")
	fmt.Printf("  4:hash             ← key: \"hash\" (长度4)\n")
	fmt.Printf("  20:<20字节二进制>   ← value: 20字节数据\n")
	fmt.Printf("  e                  ← 字典结束\n")
	fmt.Printf("\n完整十六进制: %s\n", hex.EncodeToString(bencodedDict))

	decoder = bencode.NewDecoder(bytes.NewReader(bencodedDict))
	decoded, _ := decoder.Decode()

	result_map := decoded.(map[string]interface{})
	hashValue := result_map["hash"].(string)

	fmt.Printf("\n解析结果:\n")
	fmt.Printf("  key: \"hash\"\n")
	fmt.Printf("  value (hex): %s\n", hex.EncodeToString([]byte(hashValue)))
	fmt.Printf("  value (原始): %v\n", []byte(hashValue))

	// 演示5: 关键代码
	fmt.Println("\n【演示5】我们解码器中的关键代码")
	fmt.Println(strings.Repeat("─", 90))

	fmt.Println(`
在 pkg/bencode/decoder.go 中:

func (d *Decoder) DecodeString() (string, error) {
  // 第1步: 读取长度
  var lenStr string
  for {
    ch, _ := d.peek()
    if ch == ':' {
      break  // 看到':'就停止读长度
    }
    if ch >= '0' && ch <= '9' {
      b, _ := d.read()
      lenStr += string(b)  // 累积数字
    }
  }

  // 第2步: 解析长度为整数
  length, _ := strconv.Atoi(lenStr)  // 如: "20" → 20

  // 第3步: 消费':'
  d.read()  // 移过':'

  // 第4步: 关键！根据长度读数据
  strBytes := make([]byte, length)
  n, err := io.ReadFull(d.r, strBytes)  // ← 读length个字节
  
  // 这就是魔法！io.ReadFull会读指定数量的字节
  // 不管那些字节是什么，都照样读
  // 即使是:, e, d, 或任意二进制都没问题
  
  return string(strBytes), nil
}

关键点: io.ReadFull(d.r, strBytes)
  • 第1个参数: reader
  • 第2个参数: []byte切片（要读的字节数由切片长度决定）
  • 行为: 读指定数量的字节，不管内容是什么
  • 结果: 返回读到的字节

这就是为什么Bencode能安全处理二进制！
长度告诉io.ReadFull()要读多少字节
完全不需要检查内容
`)

	// 演示6: 对比
	fmt.Println("\n【演示6】为什么JSON不行？")
	fmt.Println(strings.Repeat("─", 90))

	fmt.Println(`
JSON处理字符串:

  JSON: {"hash": "..."}
  问题: 如果数据中有引号字符呢？
       需要转义为 \\\"
       
  后果: 
    • 复杂性增加
    • 二进制数据无法直接表示
    • 必须用base64编码等方式转换
    • 效率降低

Bencode处理字符串:

  Bencode: d4:hash20:<20字节>e
  优点: 
    • 长度明确标记
    • 不需要转义
    • 二进制数据直接存储
    • 完全二进制安全！

BitTorrent为什么用Bencode？
  ✓ 紧凑（比JSON小）
  ✓ 高效（解析简单）
  ✓ 二进制安全（SHA1可以直接存）
  ✓ 无歧义（长度完全确定）
`)

	fmt.Println("\n" + strings.Repeat("═", 90))
	fmt.Println("✨ 总结")
	fmt.Println(strings.Repeat("═", 90))

	fmt.Println(`
解码器处理二进制的秘密:

1. Bencode字符串: <长度>:<数据>
2. 解码器读长度
3. io.ReadFull() 根据长度读指定字节数
4. 无论字节是什么都照样读

结果:
  ✓ ASCII字符✓
  ✓ 特殊字符（d e : 等）✓
  ✓ 纯二进制（SHA1哈希）✓
  ✓ 混合内容✓

这就是Bencode的设计精妙之处！
`)

	fmt.Println(strings.Repeat("═", 90) + "\n")
}
