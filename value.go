package eflag

import (
	"reflect"
	"time"
)

// Value implemented flag.Value interface.
type Value struct {
	val    string
	rval   reflect.Value
	config *Config
}

// NewValue is the constructor of Value.
func NewValue(v string, rv reflect.Value, c *Config) *Value {
	return &Value{
		val:    v,
		rval:   rv,
		config: c,
	}
}

// String
func (v *Value) String() string {
	return v.val
}

// Set set new value.
func (v *Value) Set(sval string) error {
	// first check time.Duration
	if _, ok := v.rval.Interface().(time.Duration); ok {
		v.rval.SetInt(int64(ParseDuration(sval, 0)))
		return nil
	}

	if val, err := ParseValue(v.rval.Type(), sval, v.config.ItemSep, v.config.MapSep); err != nil {
		return err
	} else {
		v.rval.Set(val)
		v.val = sval
		return nil
	}
}

// BoolValue set flag as bool option.
type BoolValue struct {
	Value
}

// NewBoolValue is the constructor of BoolValue.
func NewBoolValue(v Value) *BoolValue {
	return &BoolValue{v}
}

func (b *BoolValue) IsBoolFlag() bool {
	return true
}
