package bencode

import (
	"bytes"
	"fmt"
	"reflect"
)

// Encode encodes any value to bencode format
// Supports int64, int, string, []interface{}, map[string]interface{}, and structs with "bencode" tags
func Encode(v interface{}) ([]byte, error) {
	encoder := NewEncoder()
	
	// If v is a struct, convert to map first
	if isStruct(v) {
		dictVal, err := structToDict(v)
		if err != nil {
			return nil, err
		}
		return encoder.Encode(dictVal)
	}
	
	return encoder.Encode(v)
}

// DecodeBytes decodes bencode data into a target value
// Supports int64, string, []interface{}, map[string]interface{}, and structs with "bencode" tags
func DecodeBytes(data []byte, v interface{}) error {
	reader := bytes.NewReader(data)
	decoder := NewDecoder(reader)
	decoded, err := decoder.Decode()
	if err != nil {
		return err
	}

	// If target is a struct, convert from map
	if isStruct(v) {
		dictVal, ok := decoded.(map[string]interface{})
		if !ok {
			return fmt.Errorf("expected dict for struct, got %T", decoded)
		}
		return dictToStruct(dictVal, v)
	}

	// Otherwise use reflection to assign
	return assignValue(decoded, v)
}

// isStruct checks if v is a struct or pointer to struct
func isStruct(v interface{}) bool {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val.Kind() == reflect.Struct
}

// structToDict converts a struct to map[string]interface{} using bencode tags
func structToDict(v interface{}) (map[string]interface{}, error) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %T", v)
	}

	result := make(map[string]interface{})
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		// Check for bencode tag
		tag, ok := field.Tag.Lookup("bencode")
		if !ok || tag == "-" {
			continue
		}

		// Skip zero values (omitempty behavior)
		if isZeroValue(fieldVal) {
			continue
		}

		// Convert field value to interface{}
		converted := fieldValueToInterface(fieldVal)
		if converted != nil {
			result[tag] = converted
		}
	}

	return result, nil
}

// dictToStruct converts map[string]interface{} to a struct using bencode tags
func dictToStruct(dictVal map[string]interface{}, v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("expected pointer to struct, got %T", v)
	}

	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("expected pointer to struct, got %T", v)
	}

	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag, ok := field.Tag.Lookup("bencode")
		if !ok || tag == "-" {
			continue
		}

		// Look for value in dict
		if dictValue, exists := dictVal[tag]; exists {
			fieldVal := val.Field(i)
			if err := assignValueToField(dictValue, fieldVal); err != nil {
				return fmt.Errorf("failed to assign field %s: %w", field.Name, err)
			}
		}
	}

	return nil
}

// fieldValueToInterface converts a reflect.Value to interface{}
func fieldValueToInterface(val reflect.Value) interface{} {
	if !val.IsValid() {
		return nil
	}

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(val.Uint())
	case reflect.String:
		return val.String()
	case reflect.Slice:
		// Return the slice even if empty/nil
		return sliceToInterface(val)
	case reflect.Array:
		return arrayToInterface(val)
	case reflect.Struct:
		// Convert nested struct to dict
		dict, _ := structToDict(val.Interface())
		return dict
	case reflect.Interface:
		// For interface{} type, recursively convert the underlying value
		if val.IsNil() {
			return nil
		}
		return fieldValueToInterface(val.Elem())
	case reflect.Ptr:
		if !val.IsNil() {
			return fieldValueToInterface(val.Elem())
		}
	}

	return nil
}

// sliceToInterface converts a slice to []interface{}
func sliceToInterface(val reflect.Value) interface{} {
	if val.Len() == 0 {
		return []interface{}{}
	}

	// Special handling for []byte
	if val.Type().Elem().Kind() == reflect.Uint8 {
		return val.Bytes()
	}

	result := make([]interface{}, val.Len())
	for i := 0; i < val.Len(); i++ {
		result[i] = fieldValueToInterface(val.Index(i))
	}
	return result
}

// arrayToInterface converts an array to []interface{}
func arrayToInterface(val reflect.Value) interface{} {
	result := make([]interface{}, val.Len())
	for i := 0; i < val.Len(); i++ {
		result[i] = fieldValueToInterface(val.Index(i))
	}
	return result
}

// assignValueToField assigns a bencode value to a struct field
func assignValueToField(src interface{}, fieldVal reflect.Value) error {
	if !fieldVal.CanSet() {
		return fmt.Errorf("field is not settable")
	}

	switch fieldVal.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if intVal, ok := src.(int64); ok {
			fieldVal.SetInt(intVal)
			return nil
		}
		return fmt.Errorf("expected int64, got %T", src)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if intVal, ok := src.(int64); ok {
			fieldVal.SetUint(uint64(intVal))
			return nil
		}
		return fmt.Errorf("expected int64, got %T", src)

	case reflect.String:
		if strVal, ok := src.(string); ok {
			fieldVal.SetString(strVal)
			return nil
		}
		return fmt.Errorf("expected string, got %T", src)

	case reflect.Slice:
		return assignToSlice(src, fieldVal)

	case reflect.Struct:
		if dictVal, ok := src.(map[string]interface{}); ok {
			return dictToStruct(dictVal, fieldVal.Addr().Interface())
		}
		return fmt.Errorf("expected dict for struct, got %T", src)

	case reflect.Interface:
		// For interface{} type, just set the value directly
		fieldVal.Set(reflect.ValueOf(src))
		return nil

	case reflect.Ptr:
		// For pointer fields, create a new instance and assign to it
		if src == nil {
			fieldVal.Set(reflect.Zero(fieldVal.Type()))
			return nil
		}
		ptrElem := fieldVal.Type().Elem()
		newVal := reflect.New(ptrElem)
		if err := assignValueToField(src, newVal.Elem()); err != nil {
			return err
		}
		fieldVal.Set(newVal)
		return nil

	default:
		return fmt.Errorf("unsupported field type: %v", fieldVal.Kind())
	}
}

// assignToSlice assigns a value to a slice field
func assignToSlice(src interface{}, fieldVal reflect.Value) error {
	// Special case for []byte
	if fieldVal.Type().Elem().Kind() == reflect.Uint8 {
		if bytesVal, ok := src.([]byte); ok {
			fieldVal.SetBytes(bytesVal)
			return nil
		}
		if strVal, ok := src.(string); ok {
			fieldVal.SetBytes([]byte(strVal))
			return nil
		}
		return fmt.Errorf("expected []byte or string, got %T", src)
	}

	// Handle []interface{}
	listVal, ok := src.([]interface{})
	if !ok {
		return fmt.Errorf("expected list, got %T", src)
	}

	// Create new slice
	slice := reflect.MakeSlice(fieldVal.Type(), len(listVal), len(listVal))
	for i, item := range listVal {
		if err := assignValueToField(item, slice.Index(i)); err != nil {
			return fmt.Errorf("list index %d: %w", i, err)
		}
	}

	fieldVal.Set(slice)
	return nil
}

// assignValue assigns a decoded value to a variable
func assignValue(src interface{}, dst interface{}) error {
	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Ptr {
		return fmt.Errorf("destination must be a pointer")
	}

	dstVal = dstVal.Elem()
	return assignValueToField(src, dstVal)
}

// isZeroValue checks if a value is a zero value
func isZeroValue(val reflect.Value) bool {
	if !val.IsValid() {
		return true
	}

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val.Uint() == 0
	case reflect.String:
		return val.String() == ""
	case reflect.Slice, reflect.Array, reflect.Map:
		return val.Len() == 0
	case reflect.Ptr:
		return val.IsNil()
	default:
		return false
	}
}
