package eflag

import (
	"reflect"
)

func ReflectVisitStructField(v interface{}, ignoreAnonymous bool, fn func(vType reflect.Value, field reflect.StructField, fieldValue reflect.Value) bool) {
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

	rawRv := reflect.ValueOf(v)
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if field.Anonymous && ignoreAnonymous {
			continue
		}
		if fn(rawRv, field, rv.Field(i)) {
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

func isStringSlice(field reflect.StructField) bool {
	if field.Type.Kind() != reflect.Slice {
		return false
	}
	if field.Type.Elem().Kind() != reflect.String {
		return false
	}
	return true
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
