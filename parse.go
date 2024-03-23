package eflag

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// ParseValue convert string to the specified type based on typ.
func ParseValue(typ reflect.Type, strval, itemSep, mapSep string) (val reflect.Value, err error) {
	items := strings.Split(strval, itemSep)
	switch typ.Kind() {
	case reflect.Map:
		val, err = ParseMap(typ, items, mapSep)
	case reflect.Slice:
		val, err = ParseSlice(typ, items)
	default:
		val, err = ParseAtomValue(typ.Kind(), strval)
	}
	return
}

// ParseMap convert items to map.
func ParseMap(typ reflect.Type, items []string, mapSep string) (reflect.Value, error) {
	var val reflect.Value
	rmap := reflect.MakeMap(reflect.MapOf(typ.Key(), typ.Elem()))

	kkind := typ.Key().Kind()
	vkind := typ.Elem().Kind()
	for _, item := range items {
		if elems := strings.Split(item, mapSep); len(elems) >= 2 {
			kval, err := ParseAtomValue(kkind, elems[0])
			if err != nil {
				return val, err
			}
			vval, err := ParseAtomValue(vkind, elems[1])
			if err != nil {
				return val, err
			}
			rmap.SetMapIndex(kval, vval)
		}
	}
	return rmap, nil
}

// ParseSlice convert items to slice.
func ParseSlice(typ reflect.Type, items []string) (reflect.Value, error) {
	var val reflect.Value
	slice := reflect.MakeSlice(reflect.SliceOf(typ.Elem()), 0, len(items))

	kind := typ.Elem().Kind()
	for _, item := range items {
		ival, err := ParseAtomValue(kind, item)
		if err != nil {
			return val, err
		}
		slice = reflect.Append(slice, ival)
	}
	return slice, nil
}

// ParseAtomValue convert string to the specified type based on kind.
func ParseAtomValue(kind reflect.Kind, strval string) (val reflect.Value, err error) {
	var v interface{}
	switch kind {
	case reflect.Bool:
		v = ParseBool(strval, false)
	case reflect.String:
		v = strval
	case reflect.Int:
		v = int(ParseInt(strval, 0, 0))
	case reflect.Int8:
		v = int8(ParseInt(strval, 8, 0))
	case reflect.Int16:
		v = int16(ParseInt(strval, 16, 0))
	case reflect.Int32:
		v = int32(ParseInt(strval, 32, 0))
	case reflect.Int64:
		v = ParseInt(strval, 64, 0)
	case reflect.Uint:
		v = uint(ParseUint(strval, 0, 0))
	case reflect.Uint8:
		v = uint8(ParseUint(strval, 8, 0))
	case reflect.Uint16:
		v = uint16(ParseUint(strval, 16, 0))
	case reflect.Uint32:
		v = uint32(ParseUint(strval, 32, 0))
	case reflect.Uint64:
		v = ParseUint(strval, 64, 0)
	case reflect.Float32:
		v = float32(ParseFloat(strval, 32, 0.0))
	case reflect.Float64:
		v = ParseFloat(strval, 64, 0.0)
	default:
		err = errors.New("Unsupported kind type")
	}
	val = reflect.ValueOf(v)
	return
}

// ParseBool convert string to boolean.
// When the conversion fails, defval specifies the default value.
func ParseBool(s string, defval bool) bool {
	if v, err := strconv.ParseBool(s); err == nil {
		return v
	}
	return defval
}

// ParseInt convert string to int64.
// When the conversion fails, defval specifies the default value.
func ParseInt(s string, bitSize int, defval int64) int64 {
	if v, err := strconv.ParseInt(s, 10, bitSize); err == nil {
		return v
	}
	return defval
}

// ParseUint convert string to uint64.
// When the conversion fails, defval specifies the default value.
func ParseUint(s string, bitSize int, defval uint64) uint64 {
	if v, err := strconv.ParseUint(s, 10, bitSize); err == nil {
		return v
	}
	return defval
}

// ParseFloat convert string to float64.
// When the conversion fails, defval specifies the default value.
func ParseFloat(s string, bitSize int, defval float64) float64 {
	if v, err := strconv.ParseFloat(s, bitSize); err == nil {
		return v
	}
	return defval
}

// ParseDuration convert string to time.Duration.
// When the conversion fails, defval specifies the default value.
func ParseDuration(s string, defval time.Duration) time.Duration {
	if v, err := time.ParseDuration(s); err == nil {
		return v
	}
	return defval
}
