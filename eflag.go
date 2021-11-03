package eflag

import (
	"errors"
	"flag"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	defaultStructTagName = "eflag"
	defaultItemSep       = "@" // for slice,map: item1@item2@item3
	defaultMapSep        = "=" // for map: key1=value1@key2=value2
	defaultFlagSet       = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	ErrUnsupportedType = errors.New("Unsupported type")
	ErrNotPointer      = errors.New("Parameter must be a pointer")
	ErrNotStruct       = errors.New("Parameter must be a struct pointer")
)

func SetTagName(s string) {
	if s != "" {
		defaultStructTagName = s
	}
}

func SetItemSep(s string) {
	if s != "" {
		defaultItemSep = s
	}
}

func SetMapSep(s string) {
	if s != "" {
		defaultMapSep = s
	}
}

func GetFlagSet() *flag.FlagSet {
	return defaultFlagSet
}

type Value struct {
	val  string
	rval reflect.Value
}

func NewValue(v string, rv reflect.Value) *Value {
	return &Value{
		val:  v,
		rval: rv,
	}
}

func (v *Value) String() string {
	return v.val
}

func (v *Value) Set(dval string) error {
	// Check time.Duration first
	if _, ok := v.rval.Interface().(time.Duration); ok {
		v.rval.SetInt(int64(ParseDuration(dval, 0)))
		return nil
	}

	val, err := ParseValue(v.rval.Type(), dval)
	if err != nil {
		return err
	}
	v.rval.Set(val)
	v.val = dval
	return nil
}

func IsInt(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	default:
		return false
	}
}

func IsUint(kind reflect.Kind) bool {
	switch kind {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}

func IsFloat(kind reflect.Kind) bool {
	switch kind {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func ParseValue(typ reflect.Type, strval string) (reflect.Value, error) {
	var val reflect.Value
	var kind = typ.Kind()
	items := strings.Split(strval, defaultItemSep)

	if kind == reflect.Map {
		kt, vt := typ.Key(), typ.Elem()
		kkind, vkind := kt.Kind(), vt.Kind()
		rmap := reflect.MakeMap(reflect.MapOf(kt, vt))
		for _, item := range items {
			strs := strings.Split(item, defaultMapSep)
			if len(strs) >= 2 {
				kval, err := ParseAtomValue(kkind, strs[0])
				if err != nil {
					return val, err
				}
				vval, err := ParseAtomValue(vkind, strs[1])
				if err != nil {
					return val, err
				}
				rmap.SetMapIndex(kval, vval)
			}
		}
		return rmap, nil
	} else if kind == reflect.Slice {
		rt := typ.Elem()
		skind := rt.Kind()
		slice := reflect.MakeSlice(reflect.SliceOf(rt), 0, len(items))
		for _, item := range items {
			if sval, err := ParseAtomValue(skind, item); err != nil {
				return val, err
			} else {
				slice = reflect.Append(slice, sval)
			}
		}
		return slice, nil
	}
	return ParseAtomValue(kind, strval)
}

func ParseAtomValue(kind reflect.Kind, strval string) (reflect.Value, error) {
	var iv int64
	var uv uint64
	var fv float64
	var val reflect.Value

	switch {
	case IsInt(kind):
		iv = ParseInt(strval, ParseBitSize(kind), 0)
	case IsUint(kind):
		uv = ParseUint(strval, ParseBitSize(kind), 0)
	case IsFloat(kind):
		fv = ParseFloat(strval, ParseBitSize(kind), float64(0))
	}

	switch {
	case kind == reflect.Bool:
		val = reflect.ValueOf(ParseBool(strval, false))
	case kind == reflect.String:
		val = reflect.ValueOf(strval)
	case kind == reflect.Int:
		val = reflect.ValueOf(int(iv))
	case kind == reflect.Int8:
		val = reflect.ValueOf(int8(iv))
	case kind == reflect.Int16:
		val = reflect.ValueOf(int16(iv))
	case kind == reflect.Int32:
		val = reflect.ValueOf(int32(iv))
	case kind == reflect.Int64:
		val = reflect.ValueOf(iv)
	case kind == reflect.Uint:
		val = reflect.ValueOf(uint(uv))
	case kind == reflect.Uint8:
		val = reflect.ValueOf(uint8(uv))
	case kind == reflect.Uint16:
		val = reflect.ValueOf(uint16(uv))
	case kind == reflect.Uint32:
		val = reflect.ValueOf(uint32(uv))
	case kind == reflect.Uint64:
		val = reflect.ValueOf(uv)
	case kind == reflect.Float32:
		val = reflect.ValueOf(float32(fv))
	case kind == reflect.Float64:
		val = reflect.ValueOf(fv)
	default:
		return val, ErrUnsupportedType
	}
	return val, nil
}

func ParseBitSize(kind reflect.Kind) int {
	switch kind {
	case reflect.Int, reflect.Uint:
		return 0
	case reflect.Int8, reflect.Uint8:
		return 8
	case reflect.Int16, reflect.Uint16:
		return 16
	case reflect.Int32, reflect.Uint32, reflect.Float32:
		return 32
	case reflect.Int64, reflect.Uint64, reflect.Float64:
		return 64
	}
	return 0
}

func ParseBool(s string, defval bool) bool {
	if bv, err := strconv.ParseBool(s); err != nil {
		return defval
	} else {
		return bv
	}
}

func ParseInt(s string, bitSize int, defval int64) int64 {
	if iv, err := strconv.ParseInt(s, 10, bitSize); err != nil {
		return defval
	} else {
		return iv
	}
}

func ParseUint(s string, bitSize int, defval uint64) uint64 {
	if uv, err := strconv.ParseUint(s, 10, bitSize); err != nil {
		return defval
	} else {
		return uv
	}
}

func ParseFloat(s string, bitSize int, defval float64) float64 {
	if fv, err := strconv.ParseFloat(s, bitSize); err != nil {
		return defval
	} else {
		return fv
	}
}

func ParseDuration(s string, defval time.Duration) time.Duration {
	if dv, err := time.ParseDuration(s); err != nil {
		return defval
	} else {
		return dv
	}
}

func Parse(v interface{}) error {
	if defaultFlagSet.Parsed() || v == nil {
		return nil
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return ErrNotPointer
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	rt := rv.Type()
	for i, j := 0, rt.NumField(); i < j; i++ {
		field := rt.Field(i)
		if field.Anonymous { // Ignore anonymous fields
			continue
		}
		tagStr := field.Tag.Get(defaultStructTagName)
		if tagStr != "" {
			bindFlag(rv.Field(i), strings.Split(tagStr, ","))
		}
	}

	return defaultFlagSet.Parse(os.Args[1:])
}

func bindFlag(v reflect.Value, tagslice []string) {
	var name = strings.TrimSpace(tagslice[0])
	var dval string
	var usage string

	if len(tagslice) > 1 {
		dval = tagslice[1]
	}
	if len(tagslice) > 2 {
		usage = tagslice[2]
	}

	val := NewValue(dval, v)
	val.Set(dval) // Set default value first

	defaultFlagSet.Var(val, name, usage)
}
