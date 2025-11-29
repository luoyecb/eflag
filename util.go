package eflag

import (
	"reflect"
)

func ReflectVisitStructField(v interface{}, fn func(vType reflect.Value, field reflect.StructField, value reflect.Value) bool) {
	if v == nil || fn == nil {
		return
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
		if rv.Kind() != reflect.Struct {
			return
		}
	} else if rv.Kind() != reflect.Struct {
		return
	}

	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		if fn(reflect.ValueOf(v), rt.Field(i), rv.Field(i)) {
			break
		}
	}
}

func isReflectType(typ reflect.Type, expected ...reflect.Kind) bool {
	kind := typ.Kind()
	for _, k := range expected {
		if kind == k {
			return true
		}
	}
	return false
}

func isStructPtr(v interface{}) bool {
	if v == nil {
		return false
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return false
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return false
	}
	return true
}
