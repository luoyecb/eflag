// Enhanced std flag package.
// Bind command-line options to struct.
package eflag

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
)

var (
	defaultEFlag = NewEFlag(COMMAND_MODE_OPTION)
	// defaultEFlag = NewEFlag(COMMAND_MODE_SUB_CMD)
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

	commandMode    CommandMode
	subCommandName string
}

// NewEFlag is the constructor of EFlag.
func NewEFlag(commandMode CommandMode, options ...EFlagOption) *EFlag {
	config := defaultConfig
	for _, opt := range options {
		opt(&config)
	}

	return &EFlag{
		flagSet:     flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		config:      &config,
		commandMode: commandMode,
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
	ReflectVisitStructField(v, true, e.parse)

	err := e.flagSet.Parse(e.checkCommandMode(true))
	if err == nil {
		e.setArgs(v)
	}
	return err
}

func (e *EFlag) checkCommandMode(exitOnError bool) (args []string) {
	if e.commandMode == COMMAND_MODE_SUB_CMD {
		isErr := true
		if len(os.Args) > SUM_COMMAND_INDEX {
			if name := os.Args[SUM_COMMAND_INDEX]; name[0] != '-' {
				isErr = false
				e.subCommandName = name
				args = os.Args[SUM_COMMAND_INDEX+1:]
			}
		}
		if isErr && exitOnError {
			fmt.Fprintf(os.Stderr, "Not a valid sub command format\n")
			os.Exit(1)
		}
	} else if e.commandMode == COMMAND_MODE_OPTION {
		args = os.Args[1:]
	}
	return
}

func (e *EFlag) parse(rv reflect.Value, field reflect.StructField, fieldValue reflect.Value) (ret bool) {
	tagName := field.Tag.Get(e.config.TagName)
	if tagName == "" {
		return
	}

	val := e.parseDefault(rv, field, fieldValue)
	usage := field.Tag.Get("usage")
	e.flagSet.Var(val, tagName, usage)

	// parse short tag
	tagNameShort := field.Tag.Get(e.config.TagNameShort)
	if tagNameShort != "" {
		cval := val
		e.flagSet.Var(cval, tagNameShort, fmt.Sprintf("%s(same as %s)", usage, tagName))
	}
	return
}

func (e *EFlag) parseDefault(rv reflect.Value, field reflect.StructField, fieldValue reflect.Value) flag.Value {
	defaultValue := field.Tag.Get("default") // from struct tag

	var flagValue flag.Value

	val := NewValue(defaultValue, fieldValue, e.config)
	flagValue = val
	if fieldValue.Kind() == reflect.Bool {
		flagValue = NewBoolValue(*val)
	}

	flagValue.Set(defaultValue)
	// from default method
	rm := rv.MethodByName(field.Name + "Default")
	if rm.IsValid() {
		if results := rm.Call(nil); len(results) > 0 {
			fieldValue.Set(results[0])
		}
	}
	return flagValue
}

func (e *EFlag) setArgs(v interface{}) {
	if e.flagSet.NArg() > 0 {
		elem := reflect.ValueOf(v).Elem()
		structField, ok := elem.Type().FieldByName("Args")
		if !ok || !isStringSlice(structField) {
			return
		}
		// set args
		elem.FieldByName("Args").Set(reflect.ValueOf(e.flagSet.Args()))
	}
}

func (e *EFlag) RunCommand() error {
	return runCommand(e.input, e.subCommandName)
}
