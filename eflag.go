// Enhanced std flag package.
// Bind command-line options to struct.
package eflag

import (
	"errors"
	"flag"
	"fmt"
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

func ParseAndRunCommand(v interface{}) error {
	err := defaultEFlag.Parse(v)
	if err != nil {
		return err
	}
	return defaultEFlag.RunCommand()
}

// EFlag
type EFlag struct {
	flagSet *flag.FlagSet
	config  *Config
	input   interface{}
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
	if e.flagSet.Parsed() {
		return nil
	}
	if !isStructPtr(v) {
		return errors.New("Must be a pointer to a struct type")
	}

	e.input = v
	ReflectVisitStructField(v, func(rv reflect.Value, field reflect.StructField, value reflect.Value) bool {
		if field.Anonymous {
			return false
		}
		tagName := field.Tag.Get(e.config.TagName)
		if tagName == "" {
			return false
		}

		defval := field.Tag.Get("default") // parse default value from tag

		var val flag.Value = NewValue(defval, value, e.config)
		if field.Type.Kind() == reflect.Bool {
			val = NewBoolValue(*(val.(*Value)))
		}
		val.Set(defval)

		// parse default value from default method
		if rm := rv.MethodByName(field.Name + "Default"); rm.IsValid() {
			results := rm.Call(nil)
			if len(results) > 0 {
				value.Set(results[0])
			}
		}

		usage := field.Tag.Get("usage")
		e.flagSet.Var(val, tagName, usage)

		// parse short tag
		tagNameShort := field.Tag.Get(e.config.TagNameShort)
		if tagNameShort != "" {
			cval := val
			e.flagSet.Var(cval, tagNameShort, fmt.Sprintf("%s(same as %s)", usage, tagName))
		}

		return false
	})

	err := e.flagSet.Parse(os.Args[1:])
	if err == nil {
		e.setArgs(v)
	}
	return err
}

func (e *EFlag) setArgs(v interface{}) {
	if e.flagSet.NArg() == 0 {
		return
	}
	elem := reflect.ValueOf(v).Elem()

	// check type
	rt := elem.Type()
	structField, ok := rt.FieldByName("Args")
	if !ok {
		return
	}
	if structField.Type.Kind() != reflect.Slice || structField.Type.Elem().Kind() != reflect.String {
		return
	}
	// set value
	elem.FieldByName("Args").Set(reflect.ValueOf(e.flagSet.Args()))
}

func (e *EFlag) RunCommand() error {
	return RunCommand(e.input)
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
