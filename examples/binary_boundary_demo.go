package main

import (
	"encoding/hex"
	"fmt"
	"strings"
)

func main() {
	fmt.Println("\n" + strings.Repeat("═", 90))
	fmt.Println("🔍 Bencode如何处理二进制数据 - 通过长度前缀")
	fmt.Println(strings.Repeat("═", 90))

	// 演示1: 简单字符串
	fmt.Println("\n【演示1】简单的ASCII字符串")
	fmt.Println(strings.Repeat("─", 90))

	fmt.Println(`
Bencode格式: <长度>:<数据>

例子: "hello"
  编码形式: 5:hello
  ├─ 5    = 长度（表示后面有5个字节）
  ├─ :    = 分隔符
  └─ hello = 5个字节的数据

内存中:
  [0x35] [0x3A] [0x68] [0x65] [0x6C] [0x6C] [0x6F]
   '5'    ':'    'h'    'e'    'l'    'l'    'o'
   ↑ 长度   ↑分隔   ↑────────────────────────────────────────
   告诉解析器  符     从这里开始读5个字节
                    ↓
                  读完5个字节后停止
                  继续解析下一部分
`)

	// 演示2: 二进制数据
	fmt.Println("\n【演示2】包含任意二进制的字符串")
	fmt.Println(strings.Repeat("─", 90))

	binaryData := []byte{0x00, 0xFF, 0xAB, 0xCD, 0x12}
	fmt.Printf(`
例子: 五个任意字节 %v
  十六进制: %s

编码形式: 5:<二进制>
  ├─ 5    = 长度（表示后面有5个字节）
  ├─ :    = 分隔符
  └─ <二进制数据> = 5个二进制字节

内存中:
  [0x35] [0x3A] [0x00] [0xFF] [0xAB] [0xCD] [0x12]
   '5'    ':'    二进制数据（无法显示为字符）
   ↑ 长度   ↑分隔   ↑────────────────────────────────────────
   告诉解析器  符     从这里开始读5个字节（不管是什么）
                    ↓
                  读完5个字节后停止
                  继续解析下一部分
`, binaryData, hex.EncodeToString(binaryData))

	// 演示3: 解析过程
	fmt.Println("\n【演示3】解析过程的伪代码")
	fmt.Println(strings.Repeat("─", 90))

	fmt.Println(`
编码的数据: [0x35] [0x3A] [0x00] [0xFF] [0xAB] [0xCD] [0x12]
           '5'    ':'    <5个任意字节>

解析器的工作流程:

步骤1: 读取长度
  ├─ 读取字节: 0x35 → 字符 '5' → 转换为数字 5
  ├─ 继续读取: 0x3A → 字符 ':' → 这是分隔符！停止
  └─ 得出: 后面应该有 5 个字节

步骤2: 根据长度读取数据
  ├─ 现在位置在 ':' 后面
  ├─ 读取接下来的 5 个字节:
  │   ├─ 第1个字节: 0x00
  │   ├─ 第2个字节: 0xFF
  │   ├─ 第3个字节: 0xAB
  │   ├─ 第4个字节: 0xCD
  │   └─ 第5个字节: 0x12
  └─ 完成！数据就是这5个字节

步骤3: 继续解析
  ├─ 现在位置在最后一个字节之后
  └─ 继续解析下一个元素
`)

	// 演示4: 实际的Bencode字典
	fmt.Println("\n【演示4】包含二进制数据的字典")
	fmt.Println(strings.Repeat("─", 90))

	fmt.Println(`
字典: {"name": "file", "hash": <20字节SHA1>}

编码过程:
  WriteRune('d')               → [0x64]       ('d'=字典开始)
  WriteString("4:hash")        → [0x34, 0x3A] ('4:')
  WriteString("20:")           → [0x32, 0x30, 0x3A] ('20:')
  Write(sha1_bytes)            → [0xBB, 0x0D, ...] (20个字节)
  WriteString("4:name")        → [0x34, 0x3A] ('4:')
  WriteString("4:file")        → [0x34, 0x3A] ('4:file')
  WriteRune('e')               → [0x65]       ('e'=字典结束)

结果 (简化版):
  [0x64] [0x34 0x3A 0x32 0x30 0x3A] [0xBB 0x0D ...20字节...] [0x34 0x3A 0x66 0x69 0x6C 0x65] [0x65]
   'd'   '4' ':'  '2'  '0' ':'    <20字节二进制>            '4' ':' 'f' 'i' 'l' 'e'    'e'

解析过程:
  看到 'd'        → 开始解析字典
  看到 '4:'       → 下个key长度4
  读4个字节 "hash" → key是"hash"
  看到 '20:'      → value长度20
  读20个字节 (二进制) → value是这20个字节
  看到 '4:'       → 下个key长度4
  读4个字节 "name" → key是"name"
  看到 '4:'       → value长度4
  读4个字节 "file" → value是"file"
  看到 'e'        → 字典结束
  
关键: 解析器永远知道要读多少字节，因为有长度前缀！
`)

	// 演示5: 为什么不会混淆
	fmt.Println("\n【演示5】为什么不会混淆边界")
	fmt.Println(strings.Repeat("─", 90))

	fmt.Println(`
想象有这样的编码数据:
  10:hello:eee5:world

这应该如何解析？

❌ 错误理解: 看到':' 就当分隔符
  "hello:eee5:world" 是一个字符串？不对！

✅ 正确解析:
  步骤1: 看到 '10:'
  ├─ 长度是 10
  └─ 后面读 10 个字节

  步骤2: 读 10 个字节
  ├─ 'h' 'e' 'l' 'l' 'o' ':' 'e' 'e' 'e' '5'
  └─ 这 10 个字节中的 ':' 和 '5' 只是数据的一部分！

  步骤3: 数据读完
  ├─ 现在位置在 ':' 之后
  └─ 下一个元素是 'world'

关键点:
  • 长度告诉解析器要读多少字节
  • 即使数据中有特殊字符（':' 'e' 'd' 等）也没关系
  • 因为解析器已经知道要读多少字节
  • 读完后就停止，不管数据里有什么
`)

	// 演示6: 具体的二进制例子
	fmt.Println("\n【演示6】完整的二进制例子")
	fmt.Println(strings.Repeat("─", 90))

	// 构造一个包含特殊字符的"字符串"
	specialData := []byte{
		0x64, 0x65, 0x6C, 0x69,  // 这看起来像 "deli"
		0x65, 0x3A, 0xFF, 0x00,  // 这看起来像 "e:___" (乱码)
		0xAB, 0xCD,
	}

	fmt.Printf(`
数据: %v
十六进制: %s
字符显示 (部分): %s

如果编码成Bencode:
  10:<上面这10个字节>

内存布局:
  [0x31] [0x30] [0x3A]
   '1'    '0'    ':'    ← 告诉解析器: 后面有10个字节
                        ↓
  [0x64] [0x65] [0x6C] [0x69] [0x65] [0x3A] [0xFF] [0x00] [0xAB] [0xCD]
   'd'    'e'    'l'    'i'    'e'    ':'    ??     ??     ??     ??

解析:
  ✓ 看到 '1' '0' ':' 知道长度是 10
  ✓ 读取接下来的 10 个字节（不管内容是什么）
  ✓ 完成！

关键: 即使数据中有 'd' 'e' ':' 这样的"特殊字符"
     解析器也完全不关心，因为长度告诉它
     "你要读的就这么多"
`, specialData, hex.EncodeToString(specialData), string(specialData))

	// 演示7: 对比
	fmt.Println("\n【演示7】对比：有长度 vs 没有长度")
	fmt.Println(strings.Repeat("─", 90))

	fmt.Println(`
如果没有长度前缀会怎样？

❌ 没有长度的格式: :hello
  解析器不知道 "hello" 有多长
  怎样才算结束？看到'e'吗？看到空格吗？
  如果数据中本身有'e'呢？就会出错！

❌ 另一个例子: :hell:world
  是 "hell:world" 吗？还是 "hell"？
  无法区分！

✅ 有长度的格式: 5:hello
  解析器知道：看到 '5:' 就读5个字节，然后停止
  无论数据中有什么都不关心，位置就在第5个字节后面

✅ 复杂例子: 10:hell:world:
  解析器知道：看到 '10:' 就读10个字节
  这10个字节 = 'hell:world:' （10个字符）
  读完就停止，继续解析下一部分
  
  即使 'hell:world:' 中有 ':' 也没问题
  因为长度告诉解析器"就这么多"
`)

	fmt.Println("\n" + strings.Repeat("═", 90))
	fmt.Println("✨ 结论")
	fmt.Println(strings.Repeat("═", 90))

	fmt.Println(`
Bencode的巧妙设计:

1. 字符串格式: <长度>:<数据>
   └─ 长度告诉解析器要读多少字节

2. 无论数据是什么都没关系:
   ├─ ASCII字符 ✓
   ├─ 任意二进制 ✓
   ├─ 包含特殊字符 ✓
   └─ 一切都可以！

3. 解析器永远知道边界在哪里:
   ├─ 看到 '5:' 就读5个字节
   ├─ 读完就停止
   └─ 完全不会混淆

4. 这样做的好处:
   ├─ 支持二进制安全（binary-safe）
   ├─ 没有歧义
   ├─ 紧凑高效
   └─ SHA1哈希、图片、任何二进制都行！

这就是为什么BitTorrent用Bencode而不用JSON！
JSON无法安全地处理任意二进制数据。
`)

	fmt.Println(strings.Repeat("═", 90) + "\n")
}
