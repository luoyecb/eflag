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
	defaultTagName = "flag"
	defaultItemSep = "@" // item1@item2@item3
	defaultMapSep  = "=" // key1=value1@key2=value2
	defaultFlagSet = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
)

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

func (v *Value) Set(sval string) error {
	// First check the time.Duration
	if _, ok := v.rval.Interface().(time.Duration); ok {
		v.rval.SetInt(int64(ParseDuration(sval, 0)))
	} else {
		if val, err := ParseValue(v.rval.Type(), sval); err != nil {
			return err
		} else {
			v.rval.Set(val)
			v.val = sval
		}
	}
	return nil
}

func Parse(v interface{}) error {
	if defaultFlagSet.Parsed() || v == nil {
		return nil
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return errors.New("Parameter must be a pointer")
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return errors.New("Parameter must be a struct pointer")
	}

	rt := rv.Type()
	for i, j := 0, rt.NumField(); i < j; i++ {
		field := rt.Field(i)
		if field.Anonymous { // ignore anonymous field
			continue
		}
		if tagStr := field.Tag.Get(defaultTagName); tagStr != "" {
			bindFlag(rv.Field(i), strings.Split(tagStr, ","))
		}
	}
	return defaultFlagSet.Parse(os.Args[1:])
}

func bindFlag(v reflect.Value, tagslice []string) {
	var dval string
	if len(tagslice) > 1 {
		dval = tagslice[1]
	}

	var usage string
	if len(tagslice) > 2 {
		usage = tagslice[2]
	}

	val := NewValue(dval, v)
	val.Set(dval) // first set default value

	defaultFlagSet.Var(val, tagslice[0], usage)
}

func ParseValue(typ reflect.Type, strval string) (reflect.Value, error) {
	var val reflect.Value
	items := strings.Split(strval, defaultItemSep)
	switch typ.Kind() {
	case reflect.Map:
		kkind := typ.Key().Kind()
		vkind := typ.Elem().Kind()
		rmap := reflect.MakeMap(reflect.MapOf(typ.Key(), typ.Elem()))
		for _, item := range items {
			if elems := strings.Split(item, defaultMapSep); len(elems) >= 2 {
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
	case reflect.Slice:
		skind := typ.Elem().Kind()
		slice := reflect.MakeSlice(reflect.SliceOf(typ.Elem()), 0, len(items))
		for _, item := range items {
			sval, err := ParseAtomValue(skind, item)
			if err != nil {
				return val, err
			}
			slice = reflect.Append(slice, sval)
		}
		return slice, nil
	default:
		return ParseAtomValue(typ.Kind(), strval)
	}
}

func ParseAtomValue(kind reflect.Kind, strval string) (val reflect.Value, err error) {
	switch kind {
	case reflect.Bool:
		val = reflect.ValueOf(ParseBool(strval, false))
	case reflect.String:
		val = reflect.ValueOf(strval)
	case reflect.Int:
		val = reflect.ValueOf(int(ParseInt(strval, 0, 0)))
	case reflect.Int8:
		val = reflect.ValueOf(int8(ParseInt(strval, 8, 0)))
	case reflect.Int16:
		val = reflect.ValueOf(int16(ParseInt(strval, 16, 0)))
	case reflect.Int32:
		val = reflect.ValueOf(int32(ParseInt(strval, 32, 0)))
	case reflect.Int64:
		val = reflect.ValueOf(ParseInt(strval, 64, 0))
	case reflect.Uint:
		val = reflect.ValueOf(uint(ParseUint(strval, 0, 0)))
	case reflect.Uint8:
		val = reflect.ValueOf(uint8(ParseUint(strval, 8, 0)))
	case reflect.Uint16:
		val = reflect.ValueOf(uint16(ParseUint(strval, 16, 0)))
	case reflect.Uint32:
		val = reflect.ValueOf(uint32(ParseUint(strval, 32, 0)))
	case reflect.Uint64:
		val = reflect.ValueOf(ParseUint(strval, 64, 0))
	case reflect.Float32:
		val = reflect.ValueOf(float32(ParseFloat(strval, 32, float64(0))))
	case reflect.Float64:
		val = reflect.ValueOf(ParseFloat(strval, 64, float64(0)))
	default:
		err = errors.New("Unsupported kind type")
	}
	return
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

func SetTagName(s string) {
	defaultTagName = s
}

func SetItemSep(s string) {
	defaultItemSep = s
}

func SetMapSep(s string) {
	defaultMapSep = s
}
