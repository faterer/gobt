package bencode

import (
	"bytes"
	"testing"
)

// TestStructEncoding tests encoding structs to bencode format
func TestStructEncoding(t *testing.T) {
	type Person struct {
		Name string `bencode:"name"`
		Age  int64  `bencode:"age"`
	}

	person := Person{
		Name: "Alice",
		Age:  30,
	}

	encoded, err := Encode(person)
	if err != nil {
		t.Fatalf("Encode() error: %v", err)
	}

	if len(encoded) == 0 {
		t.Error("Encode() produced empty result")
	}

	// Should be a dict starting with 'd'
	if encoded[0] != 'd' {
		t.Errorf("Expected dict (d), got %c", rune(encoded[0]))
	}
}

// TestStructDecoding tests decoding bencode to structs
func TestStructDecoding(t *testing.T) {
	type Book struct {
		Title  string `bencode:"title"`
		Pages  int64  `bencode:"pages"`
		Author string `bencode:"author"`
	}

	// Create proper bencode: d6:author5:Smith5:pagesi300e5:title4:Gone e
	// This encodes: {author: "Smith", pages: 300, title: "Gone"}
	data := []byte("d6:author5:Smith5:pagesi300e5:title4:Gonee")

	var book Book
	err := DecodeBytes(data, &book)
	if err != nil {
		t.Fatalf("DecodeBytes() error: %v", err)
	}

	if book.Title != "Gone" {
		t.Errorf("Title = %s, want Gone", book.Title)
	}

	if book.Pages != 300 {
		t.Errorf("Pages = %d, want 300", book.Pages)
	}

	if book.Author != "Smith" {
		t.Errorf("Author = %s, want Smith", book.Author)
	}
}

// TestStructWithBytes tests encoding/decoding structs with byte fields
func TestStructWithBytes(t *testing.T) {
	type Data struct {
		Name    string `bencode:"name"`
		Content []byte `bencode:"content"`
	}

	original := Data{
		Name:    "test",
		Content: []byte{0x00, 0x01, 0xFF, 0xFE},
	}

	// Encode
	encoded, err := Encode(original)
	if err != nil {
		t.Fatalf("Encode() error: %v", err)
	}

	// Decode
	var decoded Data
	err = DecodeBytes(encoded, &decoded)
	if err != nil {
		t.Fatalf("DecodeBytes() error: %v", err)
	}

	if decoded.Name != original.Name {
		t.Errorf("Name mismatch: %s vs %s", decoded.Name, original.Name)
	}

	if !bytes.Equal(decoded.Content, original.Content) {
		t.Errorf("Content mismatch: %v vs %v", decoded.Content, original.Content)
	}
}

// TestStructWithList tests structs containing list fields
func TestStructWithList(t *testing.T) {
	type Config struct {
		Name  string        `bencode:"name"`
		Hosts []interface{} `bencode:"hosts"`
	}

	original := Config{
		Name:  "server",
		Hosts: []interface{}{"localhost", "example.com"},
	}

	// Encode
	encoded, err := Encode(original)
	if err != nil {
		t.Fatalf("Encode() error: %v", err)
	}

	// Decode
	var decoded Config
	err = DecodeBytes(encoded, &decoded)
	if err != nil {
		t.Fatalf("DecodeBytes() error: %v", err)
	}

	if decoded.Name != original.Name {
		t.Errorf("Name mismatch")
	}

	if len(decoded.Hosts) != 2 {
		t.Errorf("Hosts count = %d, want 2", len(decoded.Hosts))
	}
}

// TestStructWithNestedStruct tests encoding/decoding nested structs
func TestStructWithNestedStruct(t *testing.T) {
	type Address struct {
		City   string `bencode:"city"`
		ZipCode string `bencode:"zip"`
	}

	type Person struct {
		Name    string  `bencode:"name"`
		Address Address `bencode:"address"`
	}

	original := Person{
		Name: "Bob",
		Address: Address{
			City:   "NYC",
			ZipCode: "10001",
		},
	}

	// Encode
	encoded, err := Encode(original)
	if err != nil {
		t.Fatalf("Encode() error: %v", err)
	}

	// Decode
	var decoded Person
	err = DecodeBytes(encoded, &decoded)
	if err != nil {
		t.Fatalf("DecodeBytes() error: %v", err)
	}

	if decoded.Name != original.Name {
		t.Errorf("Name mismatch")
	}

	if decoded.Address.City != original.Address.City {
		t.Errorf("City mismatch")
	}

	if decoded.Address.ZipCode != original.Address.ZipCode {
		t.Errorf("ZipCode mismatch")
	}
}

// TestStructWithIgnoredFields tests that fields without bencode tag are ignored
func TestStructWithIgnoredFields(t *testing.T) {
	type Data struct {
		Public  string `bencode:"public"`
		Private string // No bencode tag
	}

	original := Data{
		Public:  "visible",
		Private: "hidden",
	}

	// Encode
	encoded, err := Encode(original)
	if err != nil {
		t.Fatalf("Encode() error: %v", err)
	}

	// The encoded data should not contain "hidden"
	if bytes.Contains(encoded, []byte("hidden")) {
		t.Error("Private field should not be encoded")
	}
}

// TestStructWithZeroValues tests that zero values are skipped
func TestStructWithZeroValues(t *testing.T) {
	type Data struct {
		Name   string `bencode:"name"`
		Count  int64  `bencode:"count"`
		Active int64  `bencode:"active"`
	}

	data := Data{
		Name:   "test",
		Count:  0, // Zero value, should be skipped
		Active: 1,
	}

	encoded, err := Encode(data)
	if err != nil {
		t.Fatalf("Encode() error: %v", err)
	}

	// Decode and check
	var decoded Data
	err = DecodeBytes(encoded, &decoded)
	if err != nil {
		t.Fatalf("DecodeBytes() error: %v", err)
	}

	if decoded.Name != data.Name {
		t.Errorf("Name mismatch")
	}

	if decoded.Active != data.Active {
		t.Errorf("Active mismatch")
	}
}

// TestEncodeTopLevelBytes tests encoding plain bytes
func TestEncodeTopLevelBytes(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04}

	encoded, err := Encode(data)
	if err != nil {
		t.Fatalf("Encode() error: %v", err)
	}

	// Should be length:data format
	if !bytes.HasPrefix(encoded, []byte("4:")) {
		t.Error("Expected length prefix for bytes")
	}
}

// TestStructFieldTypes tests various field types in structs
func TestStructFieldTypes(t *testing.T) {
	type AllTypes struct {
		Str     string `bencode:"str"`
		Int     int64  `bencode:"int"`
		IntNeg  int64  `bencode:"intneg"`
		Bytes   []byte `bencode:"bytes"`
		List    []interface{} `bencode:"list"`
	}

	original := AllTypes{
		Str:     "test",
		Int:     42,
		IntNeg:  -10,
		Bytes:   []byte{0xFF, 0xFE},
		List:    []interface{}{"item1", int64(100)},
	}

	// Encode
	encoded, err := Encode(original)
	if err != nil {
		t.Fatalf("Encode() error: %v", err)
	}

	// Decode
	var decoded AllTypes
	err = DecodeBytes(encoded, &decoded)
	if err != nil {
		t.Fatalf("DecodeBytes() error: %v", err)
	}

	if decoded.Str != original.Str {
		t.Errorf("Str mismatch")
	}

	if decoded.Int != original.Int {
		t.Errorf("Int mismatch")
	}

	if decoded.IntNeg != original.IntNeg {
		t.Errorf("IntNeg mismatch: %d vs %d", decoded.IntNeg, original.IntNeg)
	}

	if !bytes.Equal(decoded.Bytes, original.Bytes) {
		t.Errorf("Bytes mismatch")
	}

	if len(decoded.List) != len(original.List) {
		t.Errorf("List length mismatch")
	}
}

// TestStructPointerFields tests structs with pointer fields
func TestStructPointerFields(t *testing.T) {
	type Inner struct {
		Value string `bencode:"value"`
	}

	type Outer struct {
		Name  string `bencode:"name"`
		Inner *Inner `bencode:"inner"`
	}

	inner := Inner{Value: "test"}
	outer := Outer{
		Name:  "outer",
		Inner: &inner,
	}

	// Encode
	encoded, err := Encode(outer)
	if err != nil {
		t.Fatalf("Encode() error: %v", err)
	}

	// Decode
	var decoded Outer
	err = DecodeBytes(encoded, &decoded)
	if err != nil {
		t.Fatalf("DecodeBytes() error: %v", err)
	}

	if decoded.Name != outer.Name {
		t.Errorf("Name mismatch")
	}

	if decoded.Inner == nil {
		t.Error("Inner should not be nil")
	} else if decoded.Inner.Value != inner.Value {
		t.Errorf("Inner.Value mismatch")
	}
}

// TestErrorCases tests error handling
func TestStructErrorCases(t *testing.T) {
	type Data struct {
		Name string `bencode:"name"`
	}

	// Test decoding non-dict into struct
	err := DecodeBytes([]byte("4:test"), &Data{})
	if err == nil {
		t.Error("Should error when decoding string as dict")
	}

	// Test decoding list into non-pointer
	var notPointer Data
	err = DecodeBytes([]byte("d4:name4:teste"), notPointer)
	if err == nil {
		t.Error("Should error when target is not pointer")
	}
}
