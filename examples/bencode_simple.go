//go:build bencode_example

package main

import (
	"bytes"
	"fmt"
	"gobt/pkg/bencode"
)

func main() {
	fmt.Println("\n=== Bencode 编码示例 ===\n")

	encoder := bencode.NewEncoder()

	// 示例1: 编码字符串
	fmt.Println("1. 编码字符串")
	result, _ := encoder.EncodeString("hello")
	fmt.Printf("   encoder.EncodeString(\"hello\") => %s\n", string(result))

	// 示例2: 编码整数
	fmt.Println("\n2. 编码整数")
	result, _ = encoder.EncodeInteger(42)
	fmt.Printf("   encoder.EncodeInteger(42) => %s\n", string(result))

	// 示例3: 编码列表
	fmt.Println("\n3. 编码列表")
	list := []interface{}{"a", "b", "c"}
	result, _ = encoder.Encode(list)
	fmt.Printf("   encoder.Encode([\"a\", \"b\", \"c\"]) => %s\n", string(result))

	// 示例4: 编码字典
	fmt.Println("\n4. 编码字典")
	dict := map[string]interface{}{
		"name": "gobt",
		"version": 1,
	}
	result, _ = encoder.Encode(dict)
	fmt.Printf("   encoder.Encode(map) => %s\n", string(result))

	// 示例5: 解码
	fmt.Println("\n5. 解码示例")
	bencodedData := []byte("d4:name4:gobt7:versioni1ee")
	decoder := bencode.NewDecoder(bytes.NewReader(bencodedData))
	decoded, _ := decoder.Decode()
	fmt.Printf("   decoder.Decode() => %v\n", decoded)
}
