// Enhanced std flag package.
// Bind command-line options to struct.
package eflag

import (
	"errors"
	"flag"
	"os"
	"reflect"
	"time"
)

var (
	defaultEFlag = NewEFlag()
)

// Parse parse command-line options to v.
func Parse(v interface{}) error {
	return defaultEFlag.Parse(v)
}

// EFlag
type EFlag struct {
	flagSet *flag.FlagSet
	config  *Config
}

// NewEFlag is the constructor of EFlag.
func NewEFlag(options ...EFlagOption) *EFlag {
	config := defaultConfig
	for _, opt := range options {
		opt(&config)
	}

	return &EFlag{
		flagSet: flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		config:  &config,
	}
}

// Parse parse command-line options to v.
func (e *EFlag) Parse(v interface{}) error {
	if e.flagSet.Parsed() || v == nil {
		return nil
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return errors.New("Must be a pointer")
	}
	rvp := rv // save pointer
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return errors.New("must be a struct")
	}

	rt := rv.Type()
	for i, j := 0, rt.NumField(); i < j; i++ {
		field := rt.Field(i)
		if field.Anonymous {
			continue
		}

		tagName := field.Tag.Get(e.config.TagName)
		if tagName != "" {
			defval := field.Tag.Get("default") // parse default value from tag

			var val flag.Value = NewValue(defval, rv.Field(i), e.config)
			if field.Type.Kind() == reflect.Bool {
				val = NewBoolValue(*(val.(*Value)))
			}
			val.Set(defval)

			// parse default value from default method
			rm := rvp.MethodByName(field.Name + "Default")
			if rm.IsValid() {
				results := rm.Call(nil)
				if len(results) > 0 {
					rv.Field(i).Set(results[0])
				}
			}

			e.flagSet.Var(val, tagName, field.Tag.Get("usage"))
		}
	}
	return e.flagSet.Parse(os.Args[1:])
}

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
