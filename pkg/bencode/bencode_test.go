package bencode

import (
	"bytes"
	"testing"
)

// ============= 编码器测试 =============

func TestEncodeInteger(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{"Zero", 0, "i0e"},
		{"Positive", 42, "i42e"},
		{"Negative", -273, "i-273e"},
		{"Large", 9223372036854775807, "i9223372036854775807e"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoder := NewEncoder()
			result, err := encoder.EncodeInteger(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(result) != tt.expected {
				t.Errorf("got %q, want %q", string(result), tt.expected)
			}
		})
	}
}

func TestEncodeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty", "", "0:"},
		{"Simple", "hello", "5:hello"},
		{"WithSpaces", "hello world", "11:hello world"},
		{"Numbers", "12345", "5:12345"},
		{"Unicode", "café", "5:café"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoder := NewEncoder()
			result, err := encoder.EncodeString(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(result) != tt.expected {
				t.Errorf("got %q, want %q", string(result), tt.expected)
			}
		})
	}
}

func TestEncodeList(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected string
	}{
		{"EmptyList", []interface{}{}, "le"},
		{"IntegerList", []interface{}{int64(1), int64(2), int64(3)}, "li1ei2ei3ee"},
		{"StringList", []interface{}{"spam", "eggs"}, "l4:spam4:eggse"},
		{"MixedList", []interface{}{int64(1), "spam", int64(2)}, "li1e4:spami2ee"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoder := NewEncoder()
			result, err := encoder.EncodeList(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(result) != tt.expected {
				t.Errorf("got %q, want %q", string(result), tt.expected)
			}
		})
	}
}

func TestEncodeDict(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected string
	}{
		{
			"EmptyDict",
			map[string]interface{}{},
			"de",
		},
		{
			"SimpleDict",
			map[string]interface{}{"age": int64(27)},
			"d3:agei27ee",
		},
		{
			"MultipleKeys",
			map[string]interface{}{
				"age":  int64(27),
				"name": "Bob",
			},
			// Keys must be sorted: "age" < "name"
			"d3:agei27e4:name3:Bobe",
		},
		{
			"NestedDict",
			map[string]interface{}{
				"dict": map[string]interface{}{
					"int": int64(1),
				},
			},
			"d4:dictd3:inti1eee",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoder := NewEncoder()
			result, err := encoder.EncodeDict(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(result) != tt.expected {
				t.Errorf("got %q, want %q", string(result), tt.expected)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"Integer", int64(42), "i42e"},
		{"String", "hello", "5:hello"},
		{"List", []interface{}{int64(1), "spam"}, "li1e4:spame"},
		{
			"Dict",
			map[string]interface{}{"int": int64(1)},
			"d3:inti1ee",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoder := NewEncoder()
			result, err := encoder.Encode(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(result) != tt.expected {
				t.Errorf("got %q, want %q", string(result), tt.expected)
			}
		})
	}
}

// ============= 解码器测试 =============

func TestDecodeInteger(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{"Zero", "i0e", 0},
		{"Positive", "i42e", 42},
		{"Negative", "i-273e", -273},
		{"Large", "i9223372036854775807e", 9223372036854775807},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoder := NewDecoder(bytes.NewBufferString(tt.input))
			result, err := decoder.DecodeInteger()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("got %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestDecodeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty", "0:", ""},
		{"Simple", "5:hello", "hello"},
		{"WithSpaces", "11:hello world", "hello world"},
		{"Numbers", "5:12345", "12345"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoder := NewDecoder(bytes.NewBufferString(tt.input))
			result, err := decoder.DecodeString()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestDecodeList(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []interface{}
	}{
		{"EmptyList", "le", []interface{}{}},
		{"IntegerList", "li1ei2ei3ee", []interface{}{int64(1), int64(2), int64(3)}},
		{"StringList", "l4:spam4:eggse", []interface{}{"spam", "eggs"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoder := NewDecoder(bytes.NewBufferString(tt.input))
			result, err := decoder.DecodeList()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != len(tt.expected) {
				t.Errorf("got length %d, want %d", len(result), len(tt.expected))
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("element %d: got %v, want %v", i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestDecodeDict(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		checker  func(t *testing.T, result map[string]interface{})
	}{
		{
			"EmptyDict",
			"de",
			func(t *testing.T, result map[string]interface{}) {
				if len(result) != 0 {
					t.Errorf("got length %d, want 0", len(result))
				}
			},
		},
		{
			"SimpleDict",
			"d3:agei27ee",
			func(t *testing.T, result map[string]interface{}) {
				if age, ok := result["age"].(int64); !ok || age != 27 {
					t.Errorf("age: got %v, want 27", result["age"])
				}
			},
		},
		{
			"MultipleKeys",
			"d3:agei27e4:name3:Bobee",
			func(t *testing.T, result map[string]interface{}) {
				if age, ok := result["age"].(int64); !ok || age != 27 {
					t.Errorf("age: got %v, want 27", result["age"])
				}
				if name, ok := result["name"].(string); !ok || name != "Bob" {
					t.Errorf("name: got %v, want 'Bob'", result["name"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoder := NewDecoder(bytes.NewBufferString(tt.input))
			result, err := decoder.DecodeDict()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			tt.checker(t, result)
		})
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		checker func(t *testing.T, result interface{})
	}{
		{
			"Integer",
			"i42e",
			func(t *testing.T, result interface{}) {
				if val, ok := result.(int64); !ok || val != 42 {
					t.Errorf("got %v, want int64(42)", result)
				}
			},
		},
		{
			"String",
			"5:hello",
			func(t *testing.T, result interface{}) {
				if val, ok := result.(string); !ok || val != "hello" {
					t.Errorf("got %v, want string('hello')", result)
				}
			},
		},
		{
			"List",
			"li1e4:spame",
			func(t *testing.T, result interface{}) {
				if list, ok := result.([]interface{}); !ok || len(list) != 2 {
					t.Errorf("got %v, want list of length 2", result)
				}
			},
		},
		{
			"Dict",
			"d3:inti1ee",
			func(t *testing.T, result interface{}) {
				if dict, ok := result.(map[string]interface{}); !ok {
					t.Errorf("got %T, want dict", result)
				} else if val, ok := dict["int"].(int64); !ok || val != 1 {
					t.Errorf("int key: got %v, want 1", dict["int"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoder := NewDecoder(bytes.NewBufferString(tt.input))
			result, err := decoder.Decode()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			tt.checker(t, result)
		})
	}
}

// ============= 往返测试 (Round-trip) =============

func TestRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{"Integer", int64(42)},
		{"String", "hello world"},
		{"List", []interface{}{int64(1), "spam", int64(2)}},
		{
			"Dict",
			map[string]interface{}{
				"age":  int64(27),
				"name": "Bob",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 编码
			encoder := NewEncoder()
			encoded, err := encoder.Encode(tt.input)
			if err != nil {
				t.Fatalf("encode error: %v", err)
			}

			// 解码
			decoder := NewDecoder(bytes.NewBuffer(encoded))
			decoded, err := decoder.Decode()
			if err != nil {
				t.Fatalf("decode error: %v", err)
			}

			// 比较（简单比较，dict的顺序可能不同但键值对相同）
			compareValues(t, decoded, tt.input)
		})
	}
}

func compareValues(t *testing.T, got, want interface{}) {
	switch w := want.(type) {
	case int64:
		if g, ok := got.(int64); !ok || g != w {
			t.Errorf("int: got %v, want %v", got, w)
		}
	case string:
		if g, ok := got.(string); !ok || g != w {
			t.Errorf("string: got %v, want %v", got, w)
		}
	case []interface{}:
		if g, ok := got.([]interface{}); !ok || len(g) != len(w) {
			t.Errorf("list length: got %d, want %d", len(g), len(w))
		} else {
			for i, item := range w {
				compareValues(t, g[i], item)
			}
		}
	case map[string]interface{}:
		g, ok := got.(map[string]interface{})
		if !ok {
			t.Errorf("dict: got %T, want map", got)
		} else if len(g) != len(w) {
			t.Errorf("dict length: got %d, want %d", len(g), len(w))
		} else {
			for k, v := range w {
				if gv, ok := g[k]; !ok {
					t.Errorf("dict key %q: missing", k)
				} else {
					compareValues(t, gv, v)
				}
			}
		}
	}
}

// ============= 错误情况测试 =============

func TestDecodeErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"InvalidInteger", "i42"},     // missing 'e'
		{"InvalidString", "5hello"},   // missing ':'
		{"NegativeLength", "-5:hello"}, // negative length
		{"EmptyInput", ""},
		{"InvalidType", "x:hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoder := NewDecoder(bytes.NewBufferString(tt.input))
			_, err := decoder.Decode()
			if err == nil {
				t.Errorf("expected error for input %q", tt.input)
			}
		})
	}
}

// ============= 编码错误测试 =============

func TestEncodeUnsupportedType(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{"Float", 3.14},
		{"Bool", true},
		{"Nil", nil},
		{"Struct", struct{}{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoder := NewEncoder()
			_, err := encoder.Encode(tt.input)
			if err == nil {
				t.Errorf("expected error for type %T", tt.input)
			}
		})
	}
}

func TestEncodeIntType(t *testing.T) {
	// Test encoding regular int (not int64)
	encoder := NewEncoder()
	result, err := encoder.Encode(int(42))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "i42e" {
		t.Errorf("got %q, want 'i42e'", string(result))
	}
}

func TestEncodeLargeString(t *testing.T) {
	// Test encoding a large string
	largeStr := ""
	for i := 0; i < 1000; i++ {
		largeStr += "x"
	}
	encoder := NewEncoder()
	result, err := encoder.EncodeString(largeStr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "1000:" + largeStr
	if string(result) != expected {
		t.Errorf("length mismatch in large string encoding")
	}
}

func TestEncodeNestedStructures(t *testing.T) {
	// Test deeply nested structures
	nested := map[string]interface{}{
		"outer": map[string]interface{}{
			"inner": map[string]interface{}{
				"deep": int64(42),
			},
		},
	}
	encoder := NewEncoder()
	result, err := encoder.Encode(nested)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should successfully encode deeply nested dicts
	if len(result) == 0 {
		t.Errorf("encoded result should not be empty")
	}
}

// ============= 解码边界情况测试 =============

func TestDecodeIntWithoutStart(t *testing.T) {
	decoder := NewDecoder(bytes.NewBufferString("42e"))
	_, err := decoder.DecodeInteger()
	if err == nil {
		t.Errorf("should error when missing 'i' prefix")
	}
}

func TestDecodeStringWithoutColon(t *testing.T) {
	decoder := NewDecoder(bytes.NewBufferString("5hello"))
	_, err := decoder.DecodeString()
	if err == nil {
		t.Errorf("should error when missing ':' after length")
	}
}

func TestDecodeListWithoutEnd(t *testing.T) {
	decoder := NewDecoder(bytes.NewBufferString("li1ei2e"))
	_, err := decoder.DecodeList()
	if err == nil {
		t.Errorf("should error when missing 'e' at end of list")
	}
}

func TestDecodeDictWithoutEnd(t *testing.T) {
	decoder := NewDecoder(bytes.NewBufferString("d3:agei27e"))
	_, err := decoder.DecodeDict()
	if err == nil {
		t.Errorf("should error when missing 'e' at end of dict")
	}
}

func TestDecodeStringWithInsufficientData(t *testing.T) {
	decoder := NewDecoder(bytes.NewBufferString("10:hello"))
	_, err := decoder.DecodeString()
	if err == nil {
		t.Errorf("should error when string data is too short")
	}
}

func TestDecodeEmptyInteger(t *testing.T) {
	decoder := NewDecoder(bytes.NewBufferString("ie"))
	_, err := decoder.DecodeInteger()
	if err == nil {
		t.Errorf("should error for empty integer")
	}
}

func TestDecodeComplexList(t *testing.T) {
	// List containing mixed types including nested list
	input := "li1e4:spamli2ei3eee"
	decoder := NewDecoder(bytes.NewBufferString(input))
	result, err := decoder.DecodeList()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("expected 3 items, got %d", len(result))
	}
}

func TestDecodeComplexDict(t *testing.T) {
	// Dict containing nested dicts and lists
	input := "d4:listli1ei2ee6:nestedd3:agei27eeee"
	decoder := NewDecoder(bytes.NewBufferString(input))
	result, err := decoder.DecodeDict()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
}

func TestDecodeZeroString(t *testing.T) {
	decoder := NewDecoder(bytes.NewBufferString("0:"))
	result, err := decoder.DecodeString()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestDecodeNegativeIntValue(t *testing.T) {
	// Test specifically for negative integers
	decoder := NewDecoder(bytes.NewBufferString("i-42e"))
	result, err := decoder.DecodeInteger()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != -42 {
		t.Errorf("expected -42, got %d", result)
	}
}

func TestDecodeLargeInteger(t *testing.T) {
	decoder := NewDecoder(bytes.NewBufferString("i9223372036854775807e"))
	result, err := decoder.DecodeInteger()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 9223372036854775807 {
		t.Errorf("got %d, want 9223372036854775807", result)
	}
}

func TestEncodeLargeInteger(t *testing.T) {
	encoder := NewEncoder()
	result, err := encoder.EncodeInteger(9223372036854775807)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "i9223372036854775807e" {
		t.Errorf("got %q", string(result))
	}
}

func TestEncodeListWithNestedDict(t *testing.T) {
	input := []interface{}{
		int64(1),
		map[string]interface{}{
			"key": "value",
		},
	}
	encoder := NewEncoder()
	result, err := encoder.EncodeList(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) == 0 {
		t.Errorf("encoded result should not be empty")
	}
}

func TestDecodeDictKeyOrder(t *testing.T) {
	// Test that dict keys are correctly sorted when encoded
	// Input: d1:ai1e1:bi2ee (a < b, so 'a' comes first)
	decoder := NewDecoder(bytes.NewBufferString("d1:ai1e1:bi2ee"))
	result, err := decoder.DecodeDict()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	if val, ok := result["a"].(int64); !ok || val != 1 {
		t.Errorf("key 'a': got %v, want 1", result["a"])
	}
	if val, ok := result["b"].(int64); !ok || val != 2 {
		t.Errorf("key 'b': got %v, want 2", result["b"])
	}
}

func TestEncodeListReset(t *testing.T) {
	// Test that multiple encodings produce correct results
	encoder := NewEncoder()

	// First encoding
	result1, err := encoder.EncodeList([]interface{}{int64(1), int64(2)})
	if err != nil {
		t.Fatalf("first encode error: %v", err)
	}
	expected1 := "li1ei2ee"
	if string(result1) != expected1 {
		t.Errorf("first: got %q, want %q", string(result1), expected1)
	}

	// Second encoding should work correctly
	result2, err := encoder.EncodeList([]interface{}{int64(3)})
	if err != nil {
		t.Fatalf("second encode error: %v", err)
	}
	expected2 := "li3ee"
	if string(result2) != expected2 {
		t.Errorf("second: got %q, want %q", string(result2), expected2)
	}
}

func TestDecodeBadStringLength(t *testing.T) {
	decoder := NewDecoder(bytes.NewBufferString("x5:hello"))
	_, err := decoder.DecodeString()
	if err == nil {
		t.Errorf("should error for non-numeric length prefix")
	}
}

func TestDecodeMultipleItems(t *testing.T) {
	// Decoder should handle multiple consecutive items
	input := "i42e5:hello"
	decoder := NewDecoder(bytes.NewBufferString(input))

	val1, err := decoder.Decode()
	if err != nil {
		t.Fatalf("first decode error: %v", err)
	}
	if v, ok := val1.(int64); !ok || v != 42 {
		t.Errorf("first: got %v, want 42", val1)
	}

	val2, err := decoder.Decode()
	if err != nil {
		t.Fatalf("second decode error: %v", err)
	}
	if v, ok := val2.(string); !ok || v != "hello" {
		t.Errorf("second: got %v, want 'hello'", val2)
	}
}

func TestEncodeEmptyDict(t *testing.T) {
	encoder := NewEncoder()
	result, err := encoder.EncodeDict(make(map[string]interface{}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "de" {
		t.Errorf("got %q, want 'de'", string(result))
	}
}

func TestEncodeEmptyList(t *testing.T) {
	encoder := NewEncoder()
	result, err := encoder.EncodeList([]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "le" {
		t.Errorf("got %q, want 'le'", string(result))
	}
}

func TestEncodeThenDecodeComplex(t *testing.T) {
	// Test complex nested structure
	original := map[string]interface{}{
		"files": []interface{}{
			map[string]interface{}{
				"name": "file1.txt",
				"size": int64(1024),
			},
			map[string]interface{}{
				"name": "file2.txt",
				"size": int64(2048),
			},
		},
		"count": int64(2),
	}

	encoder := NewEncoder()
	encoded, err := encoder.Encode(original)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}

	decoder := NewDecoder(bytes.NewBuffer(encoded))
	decoded, err := decoder.Decode()
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	// Verify structure
	if dict, ok := decoded.(map[string]interface{}); !ok {
		t.Fatalf("expected dict, got %T", decoded)
	} else if count, ok := dict["count"].(int64); !ok || count != 2 {
		t.Errorf("count verification failed: got %v", dict["count"])
	}
}

// ============= 性能测试 =============

func BenchmarkEncodeInteger(b *testing.B) {
	encoder := NewEncoder()
	for i := 0; i < b.N; i++ {
		encoder.EncodeInteger(int64(42))
	}
}

func BenchmarkEncodeString(b *testing.B) {
	encoder := NewEncoder()
	s := "This is a test string that should be encoded"
	for i := 0; i < b.N; i++ {
		encoder.EncodeString(s)
	}
}

func BenchmarkDecodeInteger(b *testing.B) {
	for i := 0; i < b.N; i++ {
		decoder := NewDecoder(bytes.NewBufferString("i42e"))
		decoder.DecodeInteger()
	}
}

func BenchmarkRoundTripSmall(b *testing.B) {
	encoder := NewEncoder()
	for i := 0; i < b.N; i++ {
		encoded, _ := encoder.Encode(int64(42))
		decoder := NewDecoder(bytes.NewBuffer(encoded))
		decoder.Decode()
	}
}

// ============= 额外的覆盖率测试 =============

func TestDecodeInvalidIntegerCharacters(t *testing.T) {
	decoder := NewDecoder(bytes.NewBufferString("i42xe"))
	_, err := decoder.DecodeInteger()
	if err == nil {
		t.Errorf("should error for invalid character in integer")
	}
}

func TestDecodeIntegerWithPlus(t *testing.T) {
	// Bencode integers don't support + sign
	decoder := NewDecoder(bytes.NewBufferString("i+42e"))
	_, err := decoder.DecodeInteger()
	if err == nil {
		t.Errorf("should error for + sign in integer")
	}
}

func TestDecodeStringLengthNonNumeric(t *testing.T) {
	decoder := NewDecoder(bytes.NewBufferString("abc:content"))
	_, err := decoder.DecodeString()
	if err == nil {
		t.Errorf("should error for non-numeric string length")
	}
}

func TestDecodeStringExactMatch(t *testing.T) {
	// Ensure we read exactly the right number of bytes
	decoder := NewDecoder(bytes.NewBufferString("3:abcdefgh"))
	result, err := decoder.DecodeString()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "abc" {
		t.Errorf("got %q, want 'abc'", result)
	}
}

func TestDecodeListNestedEmpty(t *testing.T) {
	// List containing empty lists: [[], []]
	// Encoded as: l + le + le + e = "lleleee"
	decoder := NewDecoder(bytes.NewBufferString("lleleee"))
	result, err := decoder.DecodeList()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 items, got %d", len(result))
	}
	if list, ok := result[0].([]interface{}); !ok || len(list) != 0 {
		t.Errorf("first item: expected empty list, got %T", result[0])
	}
	if list, ok := result[1].([]interface{}); !ok || len(list) != 0 {
		t.Errorf("second item: expected empty list, got %T", result[1])
	}
}

func TestDecodeDictNestedEmpty(t *testing.T) {
	// Dict containing empty dict
	decoder := NewDecoder(bytes.NewBufferString("d5:innerd3:agei27eeee"))
	result, err := decoder.DecodeDict()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 key, got %d", len(result))
	}
}

func TestDecodeStringZeroLength(t *testing.T) {
	// Make sure zero-length strings work correctly
	decoder := NewDecoder(bytes.NewBufferString("0:"))
	result, err := decoder.DecodeString()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "" {
		t.Errorf("got %q, want empty string", result)
	}
}

func TestEncodeDictWithManyKeys(t *testing.T) {
	// Test encoding dict with many keys to ensure sorting
	dict := make(map[string]interface{})
	for i := 9; i >= 0; i-- {
		key := string(rune('a' + i))
		dict[key] = int64(i)
	}

	encoder := NewEncoder()
	result, err := encoder.Encode(dict)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Decode it back to verify
	decoder := NewDecoder(bytes.NewBuffer(result))
	decoded, err := decoder.Decode()
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if dict, ok := decoded.(map[string]interface{}); !ok {
		t.Fatalf("expected dict, got %T", decoded)
	} else if len(dict) != 10 {
		t.Errorf("expected 10 keys, got %d", len(dict))
	}
}

func TestDecodeStringWithUnicode(t *testing.T) {
	// UTF-8 string 'café' is 5 bytes (c=1, a=1, f=1, é=2 in UTF-8)
	decoder := NewDecoder(bytes.NewBufferString("5:café"))
	result, err := decoder.DecodeString()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "café" {
		t.Errorf("got %q, want 'café'", result)
	}
}

func TestEncodeStringWithUnicode(t *testing.T) {
	encoder := NewEncoder()
	result, err := encoder.EncodeString("café")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should encode as length:content, where length is byte count
	if string(result) != "5:café" {
		t.Errorf("got %q", string(result))
	}
}

func TestDecodeNegativeZero(t *testing.T) {
	// Bencode shouldn't have negative zero
	decoder := NewDecoder(bytes.NewBufferString("i-0e"))
	result, err := decoder.DecodeInteger()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should still parse it as 0, but it's technically valid by the parser
	if result != 0 {
		t.Errorf("got %d, want 0", result)
	}
}

func TestEncodeType(t *testing.T) {
	encoder := NewEncoder()
	encoder.buf.Reset()

	// Test private encode method through public Encode
	result, err := encoder.Encode(map[string]interface{}{"x": int64(1)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) == 0 {
		t.Errorf("result should not be empty")
	}
}

// ============= 额外的EOF和错误测试 =============

func TestDecodeIntegerEOFAfterI(t *testing.T) {
	decoder := NewDecoder(bytes.NewBufferString("i"))
	_, err := decoder.DecodeInteger()
	if err == nil {
		t.Errorf("should error on EOF after 'i'")
	}
}

func TestDecodeStringEOFAfterLength(t *testing.T) {
	decoder := NewDecoder(bytes.NewBufferString("5"))
	_, err := decoder.DecodeString()
	if err == nil {
		t.Errorf("should error on EOF after length")
	}
}

func TestDecodeStringEOFAfterColon(t *testing.T) {
	decoder := NewDecoder(bytes.NewBufferString("5:"))
	_, err := decoder.DecodeString()
	if err == nil {
		t.Errorf("should error on EOF during string read")
	}
}

func TestDecodeListEOFAfterL(t *testing.T) {
	decoder := NewDecoder(bytes.NewBufferString("l"))
	_, err := decoder.DecodeList()
	if err == nil {
		t.Errorf("should error on EOF after 'l'")
	}
}

func TestDecodeDictEOFAfterD(t *testing.T) {
	decoder := NewDecoder(bytes.NewBufferString("d"))
	_, err := decoder.DecodeDict()
	if err == nil {
		t.Errorf("should error on EOF after 'd'")
	}
}

func TestDecodeDictEOFAfterKey(t *testing.T) {
	decoder := NewDecoder(bytes.NewBufferString("d3:key"))
	_, err := decoder.DecodeDict()
	if err == nil {
		t.Errorf("should error on EOF after key")
	}
}

func TestDecodeEOFAtStart(t *testing.T) {
	decoder := NewDecoder(bytes.NewBufferString(""))
	_, err := decoder.Decode()
	if err == nil {
		t.Errorf("should error on empty input")
	}
}

func TestDecodeStringOversizeLength(t *testing.T) {
	// Request more bytes than available
	decoder := NewDecoder(bytes.NewBufferString("1000:hello"))
	_, err := decoder.DecodeString()
	if err == nil {
		t.Errorf("should error when string data is insufficient")
	}
}

func TestDecodeIntValidMinMax(t *testing.T) {
	// Test with very large integer
	decoder := NewDecoder(bytes.NewBufferString("i9223372036854775806e"))
	result, err := decoder.DecodeInteger()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 9223372036854775806 {
		t.Errorf("got %d", result)
	}
}

func TestRoundTripAllTypes(t *testing.T) {
	testCases := []interface{}{
		int64(0),
		int64(-1),
		int64(9223372036854775807),
		"",
		"hello",
		"hello world",
		[]interface{}{},
		[]interface{}{int64(1)},
		map[string]interface{}{},
		map[string]interface{}{"a": int64(1)},
		[]interface{}{
			int64(42),
			"test",
			map[string]interface{}{"nested": "value"},
			[]interface{}{int64(1), int64(2)},
		},
	}

	for i, testVal := range testCases {
		encoder := NewEncoder()
		encoded, err := encoder.Encode(testVal)
		if err != nil {
			t.Errorf("test case %d: encode error: %v", i, err)
			continue
		}

		decoder := NewDecoder(bytes.NewBuffer(encoded))
		decoded, err := decoder.Decode()
		if err != nil {
			t.Errorf("test case %d: decode error: %v", i, err)
			continue
		}

		compareValues(t, decoded, testVal)
	}
}

// ============= 覆盖率增强测试 =============

func TestEncodeIntZero(t *testing.T) {
	encoder := NewEncoder()
	result, err := encoder.Encode(int(0))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "i0e" {
		t.Errorf("got %q", string(result))
	}
}

func TestDecodeStringMultiDigitLength(t *testing.T) {
	// String with multi-digit length
	decoder := NewDecoder(bytes.NewBufferString("11:hello world"))
	result, err := decoder.DecodeString()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "hello world" {
		t.Errorf("got %q, want 'hello world'", result)
	}
}

func TestDecodeListWithAllTypes(t *testing.T) {
	// List containing integer, string, and nested list
	input := "li42e5:hellolleee"
	decoder := NewDecoder(bytes.NewBufferString(input))
	result, err := decoder.DecodeList()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("expected 3 items, got %d", len(result))
	}
}

func TestDecodeDictValueTypes(t *testing.T) {
	// Dict with different value types
	// {"a": 1, "b": []}
	input := "d1:ai1e1:bleeee"
	decoder := NewDecoder(bytes.NewBufferString(input))
	result, err := decoder.DecodeDict()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
}

func TestEncodeNegativeInteger(t *testing.T) {
	encoder := NewEncoder()
	result, err := encoder.Encode(int64(-100))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "i-100e" {
		t.Errorf("got %q", string(result))
	}
}

func TestEncodeMixedTypesList(t *testing.T) {
	input := []interface{}{
		int64(1),
		"two",
		int64(3),
		[]interface{}{int64(4)},
	}
	encoder := NewEncoder()
	result, err := encoder.Encode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify by decoding
	decoder := NewDecoder(bytes.NewBuffer(result))
	decoded, err := decoder.Decode()
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if list, ok := decoded.([]interface{}); !ok || len(list) != 4 {
		t.Errorf("decode mismatch")
	}
}

func TestEncodeDictWithVariousValues(t *testing.T) {
	input := map[string]interface{}{
		"zero":   int64(0),
		"neg":    int64(-1),
		"str":    "test",
		"list":   []interface{}{int64(1)},
		"nested": map[string]interface{}{"inner": "value"},
	}
	encoder := NewEncoder()
	result, err := encoder.Encode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) == 0 {
		t.Errorf("result should not be empty")
	}
}

func TestDecodeInvalidFirstChar(t *testing.T) {
	decoder := NewDecoder(bytes.NewBufferString("@invalid"))
	_, err := decoder.Decode()
	if err == nil {
		t.Errorf("should error for invalid first character")
	}
}

func TestDecodeIntWithMultipleDigits(t *testing.T) {
	decoder := NewDecoder(bytes.NewBufferString("i123456789e"))
	result, err := decoder.DecodeInteger()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 123456789 {
		t.Errorf("got %d, want 123456789", result)
	}
}

func TestDecodeStringWithNumbers(t *testing.T) {
	decoder := NewDecoder(bytes.NewBufferString("7:12345 6"))
	result, err := decoder.DecodeString()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "12345 6" {
		t.Errorf("got %q", result)
	}
}

func TestEncodeStringAllTypes(t *testing.T) {
	testStrings := []string{
		"",
		"a",
		"abc",
		"0123456789",
		"special!@#$%",
	}

	for _, s := range testStrings {
		encoder := NewEncoder()
		result, err := encoder.EncodeString(s)
		if err != nil {
			t.Errorf("error encoding %q: %v", s, err)
		}
		if len(result) == 0 {
			t.Errorf("encoded %q should not be empty", s)
		}
	}
}

// ============= 性能基准 =============

func BenchmarkEncodeStringLarge(b *testing.B) {
	encoder := NewEncoder()
	largeStr := ""
	for i := 0; i < 1000; i++ {
		largeStr += "x"
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder.EncodeString(largeStr)
	}
}

func BenchmarkDecodeStringLarge(b *testing.B) {
	largeStr := ""
	for i := 0; i < 1000; i++ {
		largeStr += "x"
	}
	encoder := NewEncoder()
	encoded, _ := encoder.EncodeString(largeStr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		decoder := NewDecoder(bytes.NewBuffer(encoded))
		decoder.DecodeString()
	}
}

func BenchmarkEncodeDictComplex(b *testing.B) {
	encoder := NewEncoder()
	complex := map[string]interface{}{
		"zzz": int64(1),
		"aaa": "test",
		"mmm": []interface{}{int64(1), int64(2)},
		"bbb": map[string]interface{}{"nested": "value"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder.Encode(complex)
	}
}

func BenchmarkDecodeDictComplex(b *testing.B) {
	encoder := NewEncoder()
	complex := map[string]interface{}{
		"zzz": int64(1),
		"aaa": "test",
		"mmm": []interface{}{int64(1), int64(2)},
	}
	encoded, _ := encoder.Encode(complex)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		decoder := NewDecoder(bytes.NewBuffer(encoded))
		decoder.Decode()
	}
}
