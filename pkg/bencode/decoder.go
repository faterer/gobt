package bencode

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// Decoder handles bencode decoding
type Decoder struct {
	r      *bufio.Reader
	lastCh byte // for peeking
	eof    bool
}

// NewDecoder creates a new Decoder instance
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r: bufio.NewReader(r),
	}
}

// Decode decodes a bencode value
// Returns interface{} which can be int64, string, []interface{}, or map[string]interface{}
func (d *Decoder) Decode() (interface{}, error) {
	ch, err := d.peek()
	if err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	switch {
	case ch == 'i':
		return d.DecodeInteger()
	case ch >= '0' && ch <= '9':
		return d.DecodeString()
	case ch == 'l':
		return d.DecodeList()
	case ch == 'd':
		return d.DecodeDict()
	default:
		return nil, fmt.Errorf("invalid bencode: unknown type starting with '%c'", ch)
	}
}

// DecodeInteger decodes a bencode integer
// Format: i<number>e
// Example: i42e for 42
func (d *Decoder) DecodeInteger() (int64, error) {
	// Consume 'i'
	if ch, err := d.read(); err != nil || ch != 'i' {
		return 0, fmt.Errorf("expected 'i', got '%c'", ch)
	}

	// Read digits (possibly with negative sign)
	var numStr string
	for {
		ch, err := d.peek()
		if err != nil {
			return 0, fmt.Errorf("unexpected EOF while reading integer")
		}

		if ch == 'e' {
			break
		}

		if ch == '-' || (ch >= '0' && ch <= '9') {
			b, _ := d.read()
			numStr += string(b)
		} else {
			return 0, fmt.Errorf("invalid character in integer: '%c'", ch)
		}
	}

	// Consume 'e'
	if ch, err := d.read(); err != nil || ch != 'e' {
		return 0, fmt.Errorf("expected 'e', got '%c'", ch)
	}

	// Parse number
	num, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer: %s", numStr)
	}

	return num, nil
}

// DecodeString decodes a bencode string
// Format: <length>:<string>
// Example: 5:hello for "hello"
func (d *Decoder) DecodeString() (string, error) {
	// Read length
	var lenStr string
	for {
		ch, err := d.peek()
		if err != nil {
			return "", fmt.Errorf("unexpected EOF while reading string length")
		}

		if ch == ':' {
			break
		}

		if ch >= '0' && ch <= '9' {
			b, _ := d.read()
			lenStr += string(b)
		} else {
			return "", fmt.Errorf("invalid character in string length: '%c'", ch)
		}
	}

	// Parse length
	length, err := strconv.Atoi(lenStr)
	if err != nil || length < 0 {
		return "", fmt.Errorf("invalid string length: %s", lenStr)
	}

	// Consume ':'
	if ch, err := d.read(); err != nil || ch != ':' {
		return "", fmt.Errorf("expected ':', got '%c'", ch)
	}

	// Read string data
	strBytes := make([]byte, length)
	n, err := io.ReadFull(d.r, strBytes)
	if err != nil || n != length {
		return "", fmt.Errorf("failed to read string data: expected %d bytes", length)
	}

	return string(strBytes), nil
}

// DecodeList decodes a bencode list
// Format: l<items>e
// Example: li1e4:spame for [1, "spam"]
func (d *Decoder) DecodeList() ([]interface{}, error) {
	// Consume 'l'
	if ch, err := d.read(); err != nil || ch != 'l' {
		return nil, fmt.Errorf("expected 'l', got '%c'", ch)
	}

	var items []interface{}

	for {
		ch, err := d.peek()
		if err != nil {
			return nil, fmt.Errorf("unexpected EOF while reading list")
		}

		if ch == 'e' {
			break
		}

		// Recursively decode each item
		item, err := d.Decode()
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	// Consume 'e'
	if ch, err := d.read(); err != nil || ch != 'e' {
		return nil, fmt.Errorf("expected 'e', got '%c'", ch)
	}

	return items, nil
}

// DecodeDict decodes a bencode dictionary
// Format: d<key><value>...e
// Note: Keys are expected to be in sorted order
func (d *Decoder) DecodeDict() (map[string]interface{}, error) {
	// Consume 'd'
	if ch, err := d.read(); err != nil || ch != 'd' {
		return nil, fmt.Errorf("expected 'd', got '%c'", ch)
	}

	dict := make(map[string]interface{})

	for {
		ch, err := d.peek()
		if err != nil {
			return nil, fmt.Errorf("unexpected EOF while reading dict")
		}

		if ch == 'e' {
			break
		}

		// Keys must be strings in bencode
		key, err := d.DecodeString()
		if err != nil {
			return nil, fmt.Errorf("invalid dict key: %w", err)
		}

		// Decode value
		value, err := d.Decode()
		if err != nil {
			return nil, err
		}

		dict[key] = value
	}

	// Consume 'e'
	if ch, err := d.read(); err != nil || ch != 'e' {
		return nil, fmt.Errorf("expected 'e', got '%c'", ch)
	}

	return dict, nil
}

// peek reads the next byte without consuming it
func (d *Decoder) peek() (byte, error) {
	if d.eof {
		return 0, io.EOF
	}

	// Use buffered reader's ability to peek
	bytes, err := d.r.Peek(1)
	if err != nil {
		if err == io.EOF {
			d.eof = true
		}
		return 0, err
	}

	if len(bytes) == 0 {
		d.eof = true
		return 0, io.EOF
	}

	return bytes[0], nil
}

// read reads and consumes the next byte
func (d *Decoder) read() (byte, error) {
	if d.eof {
		return 0, io.EOF
	}

	b, err := d.r.ReadByte()
	if err != nil {
		if err == io.EOF {
			d.eof = true
		}
		return 0, err
	}

	return b, nil
}
