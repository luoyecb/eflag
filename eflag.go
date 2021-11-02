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
	// It should be checked first
	if _, ok := v.rval.Interface().(time.Duration); ok {
		v.rval.SetInt(int64(ParseDuration(dval, 0)))
		return nil
	}

	switch kind := v.rval.Kind(); kind {
	case reflect.Bool:
		v.rval.SetBool(ParseBool(dval, false))
	case reflect.String:
		v.rval.SetString(dval)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.rval.SetInt(ParseInt(dval, ParseBitSize(kind), 0))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.rval.SetUint(ParseUint(dval, ParseBitSize(kind), 0))
	case reflect.Float32, reflect.Float64:
		v.rval.SetFloat(ParseFloat(dval, ParseBitSize(kind), float64(0)))
	default:
		return ErrUnsupportedType
	}

	v.val = dval
	return nil
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
	val.Set(dval) // Set default value

	defaultFlagSet.Var(val, name, usage)
}
