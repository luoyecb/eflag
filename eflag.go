// Enhanced std flag package.
// Bind command-line options to struct.
package eflag

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
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
	return defaultEFlag.ParseAndRunCommand(v)
}

// EFlag
type EFlag struct {
	flagSet *flag.FlagSet
	config  *Config

	errOutput strings.Builder

	commandMode CommandMode
	commandName string
	commandList []*Command
}

// NewEFlag is the constructor of EFlag.
func NewEFlag(commandMode CommandMode, options ...EFlagOption) *EFlag {
	config := defaultConfig
	for _, opt := range options {
		opt(&config)
	}

	eFlag := &EFlag{
		config:      &config,
		commandMode: commandMode,
	}

	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagSet.Usage = eFlag.Usage
	flagSet.SetOutput(&eFlag.errOutput)

	eFlag.flagSet = flagSet
	return eFlag
}

// Parse parse command-line options to v.
func (e *EFlag) Parse(v interface{}) error {
	if e.flagSet.Parsed() {
		return nil
	}
	if !isStructPtr(v) {
		return errors.New("Must be a pointer to a struct type")
	}

	ReflectVisitStructField(v, true, e.parse)
	err := e.flagSet.Parse(e.checkCommandMode(true))
	if err == nil {
		e.setArgs(v)
	}
	return err
}

func (e *EFlag) isMode(mode CommandMode) bool {
	return e.commandMode == mode
}

func (e *EFlag) checkCommandMode(exitOnError bool) (args []string) {
	if e.isMode(COMMAND_MODE_OPTION) {
		return os.Args[1:]
	} else if e.isMode(COMMAND_MODE_SUB_CMD) && len(os.Args) > SUM_COMMAND_INDEX {
		if name := os.Args[SUM_COMMAND_INDEX]; name[0] != '-' {
			e.commandName = name
			return os.Args[SUM_COMMAND_INDEX+1:]
		} else if exitOnError {
			fmt.Fprintf(os.Stderr, "Not valid sub command format\n")
			os.Exit(1)
		}
	}
	return
}

func (e *EFlag) parse(rv reflect.Value, field reflect.StructField, fieldValue reflect.Value) (ret bool) {
	e.parseCommand(rv, field, fieldValue)
	// parse flag
	tagName := field.Tag.Get(e.config.TagName)
	if tagName == "" {
		return
	}

	val := e.parseDefault(rv, field, fieldValue)
	usage := field.Tag.Get("usage")
	e.flagSet.Var(val, tagName, usage)

	// parse short flag
	tagNameShort := field.Tag.Get(e.config.TagNameShort)
	if tagNameShort != "" {
		cval := val
		e.flagSet.Var(cval, tagNameShort, fmt.Sprintf("%s(same as %s)", usage, tagName))
	}
	return
}

func (e *EFlag) parseCommand(rv reflect.Value, field reflect.StructField, fieldValue reflect.Value) {
	if e.isMode(COMMAND_MODE_SUB_CMD) {
		// parse sub command
		tagStr := field.Tag.Get(COMMAND_SUB_COMMAND_TAG_KEY)
		if tagStr == "" {
			return
		}
		e.commandList = append(e.commandList, &Command{
			Name:       tagStr,
			MethodName: field.Name,
			Mode:       COMMAND_MODE_SUB_CMD,
			Usage:      field.Tag.Get("usage"),
			rv:         rv,
		})
	} else if e.isMode(COMMAND_MODE_OPTION) {
		// parse option command
		cmdStr, ok := field.Tag.Lookup(COMMAND_FIELD_TAG_KEY)
		if !ok || !isReflectType(field.Type, reflect.Bool, reflect.String) {
			return
		}
		methodName, runFlag := parseCommand(cmdStr, field.Type.Kind(), field.Name)
		e.commandList = append(e.commandList, &Command{
			Name:       methodName,
			MethodName: methodName,
			Mode:       COMMAND_MODE_OPTION,
			Usage:      field.Tag.Get("usage"),
			runFlag:    runFlag,
			value:      fieldValue,
			rv:         rv,
		})
	}
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

func (e *EFlag) RunCommand() {
	var currentCommand *Command
	for _, cmd := range e.commandList {
		if cmd.ShouldRun(e.commandName) {
			currentCommand = cmd
			break
		}
	}
	if currentCommand != nil {
		currentCommand.Run()
	} else if e.isMode(COMMAND_MODE_SUB_CMD) {
		fmt.Fprintf(os.Stderr, "Not support sub command\n")
		e.Usage()
		os.Exit(1)
	}
}

func (e *EFlag) ParseAndRunCommand(v interface{}) error {
	err := e.Parse(v)
	if err == nil {
		e.RunCommand()
	}
	return err
}

func (e *EFlag) Usage() {
	binName := e.flagSet.Name()
	if e.errOutput.Len() != 0 {
		e.errOutput.WriteByte('\n')
	}
	e.errOutput.WriteString(fmt.Sprintf("Usage of %s:\n", binName))
	if e.isMode(COMMAND_MODE_SUB_CMD) && len(e.commandList) > 0 {
		e.errOutput.WriteString(fmt.Sprintf("%s {SUB_COMMAND} {OPTION}\n", binName))
		e.errOutput.WriteString("SUB_COMMAND is\n")
		e.errOutput.WriteString(formatCommandUsage(e.commandList))
		e.errOutput.WriteString("\nOPTION is\n")
	}
	e.flagSet.PrintDefaults()
	fmt.Print(e.errOutput.String())
}
