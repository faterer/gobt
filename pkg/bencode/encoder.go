package bencode

import (
	"bytes"
	"fmt"
	"sort"
)

// Encoder handles bencode encoding
type Encoder struct {
	buf bytes.Buffer
}

// NewEncoder creates a new Encoder instance
func NewEncoder() *Encoder {
	return &Encoder{
		buf: bytes.Buffer{},
	}
}

// Encode encodes a value into bencode format
// Supports: int64, string, []interface{}, map[string]interface{}
func (e *Encoder) Encode(v interface{}) ([]byte, error) {
	e.buf.Reset()
	err := e.encode(v)
	if err != nil {
		return nil, err
	}
	return e.buf.Bytes(), nil
}

// encode recursively encodes a value
func (e *Encoder) encode(v interface{}) error {
	switch val := v.(type) {
	case int64:
		return e.encodeIntegerValue(val)
	case int:
		return e.encodeIntegerValue(int64(val))
	case string:
		return e.encodeStringValue(val)
	case []byte:
		return e.encodeByteStringValue(val)
	case []interface{}:
		return e.encodeListValue(val)
	case map[string]interface{}:
		return e.encodeDictValue(val)
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

// EncodeInteger encodes an integer to bencode format
// Format: i<number>e
// Example: i42e for 42, i-273e for -273
func (e *Encoder) EncodeInteger(n int64) ([]byte, error) {
	e.buf.Reset()
	err := e.encodeIntegerValue(n)
	if err != nil {
		return nil, err
	}
	return e.buf.Bytes(), nil
}

// encodeIntegerValue encodes an integer without resetting buffer
func (e *Encoder) encodeIntegerValue(n int64) error {
	e.buf.WriteRune('i')
	e.buf.WriteString(fmt.Sprintf("%d", n))
	e.buf.WriteRune('e')
	return nil
}

// EncodeString encodes a string to bencode format
// Format: <length>:<string>
// Example: 5:hello for "hello"
func (e *Encoder) EncodeString(s string) ([]byte, error) {
	e.buf.Reset()
	err := e.encodeStringValue(s)
	if err != nil {
		return nil, err
	}
	return e.buf.Bytes(), nil
}

// encodeStringValue encodes a string without resetting buffer
func (e *Encoder) encodeStringValue(s string) error {
	length := len(s)
	e.buf.WriteString(fmt.Sprintf("%d:", length))
	e.buf.WriteString(s)
	return nil
}

// encodeByteStringValue encodes a byte string without resetting buffer
func (e *Encoder) encodeByteStringValue(b []byte) error {
	length := len(b)
	e.buf.WriteString(fmt.Sprintf("%d:", length))
	e.buf.Write(b)
	return nil
}

// EncodeList encodes a list/slice to bencode format
// Format: l<items>e
// Example: li1e4:spame for [1, "spam"]
func (e *Encoder) EncodeList(items []interface{}) ([]byte, error) {
	e.buf.Reset()
	err := e.encodeListValue(items)
	if err != nil {
		return nil, err
	}
	return e.buf.Bytes(), nil
}

// encodeListValue encodes a list without resetting buffer
func (e *Encoder) encodeListValue(items []interface{}) error {
	e.buf.WriteRune('l')
	for _, item := range items {
		if err := e.encode(item); err != nil {
			return err
		}
	}
	e.buf.WriteRune('e')
	return nil
}

// EncodeDict encodes a dictionary to bencode format
// Format: d<key1><value1><key2><value2>...e
// Note: Keys must be in sorted order
// Example: d3:agei27e4:name3:Bobee for {"age": 27, "name": "Bob"}
func (e *Encoder) EncodeDict(dict map[string]interface{}) ([]byte, error) {
	e.buf.Reset()
	err := e.encodeDictValue(dict)
	if err != nil {
		return nil, err
	}
	return e.buf.Bytes(), nil
}

// encodeDictValue encodes a dictionary without resetting buffer
func (e *Encoder) encodeDictValue(dict map[string]interface{}) error {
	e.buf.WriteRune('d')

	// Sort keys alphabetically (required by bencode spec)
	keys := make([]string, 0, len(dict))
	for k := range dict {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Encode each key-value pair in sorted order
	for _, key := range keys {
		// Encode key (always a string)
		if err := e.encodeStringValue(key); err != nil {
			return err
		}
		// Encode value
		if err := e.encode(dict[key]); err != nil {
			return err
		}
	}

	e.buf.WriteRune('e')
	return nil
}
