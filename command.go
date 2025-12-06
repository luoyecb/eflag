package eflag

import (
	"reflect"
	"strings"

	"github.com/luoyecb/eflag/text"
)

type CommandMode uint

const (
	COMMAND_FIELD_TAG_KEY       = "command"
	COMMAND_METHOD_NAME_KEY     = "Command"
	COMMAND_SUB_COMMAND_TAG_KEY = "sub_command"

	SUM_COMMAND_INDEX = 1

	COMMAND_MODE_OPTION  CommandMode = iota // Option command, eg: go run main.go -COMMAND
	COMMAND_MODE_SUB_CMD                    // Sub command, eg: go run main.go SUB_COMMAND
)

type Command struct {
	Name       string
	MethodName string
	Mode       CommandMode
	Usage      string
	runFlag    string
	rv         reflect.Value
	value      reflect.Value
}

func (c *Command) UsageString() string {
	return c.Name + "    " + c.Usage
}

func formatCommandUsage(cmds []*Command) string {
	usages := make([]string, 0, len(cmds))
	for _, cmd := range cmds {
		usages = append(usages, cmd.UsageString())
	}

	align := text.NewAlignment("    ", "  ")
	return align.FormatLines(usages)
}

func (c *Command) ShouldRun(name string) bool {
	if c.Mode == COMMAND_MODE_SUB_CMD {
		return c.Name == name
	} else if c.Mode == COMMAND_MODE_OPTION {
		switch c.value.Kind() {
		case reflect.Bool:
			b := c.value.Bool()
			return (b && c.runFlag == "true") || (!b && c.runFlag == "false")
		case reflect.String:
			s := c.value.String()
			return (s == "" && c.runFlag == "empty") || (s != "" && c.runFlag == "notempty")
		}
	}
	return false
}

func (c *Command) Run() {
	method := c.rv.MethodByName(c.MethodName + COMMAND_METHOD_NAME_KEY)
	if method.IsValid() {
		method.Call(nil)
	}
}

// Format: methodName,runFlag
// Possible runFlag: true false
func parseCommand(cmdStr string, kind reflect.Kind, defaultMethodName string) (methodName, runFlag string) {
	parts := strings.Split(cmdStr, ",")
	if p1 := strings.TrimSpace(parts[0]); p1 != "" {
		methodName = p1
	}
	if len(parts) == 2 {
		runFlag = strings.TrimSpace(parts[1])
	}

	// default
	if methodName == "" {
		methodName = defaultMethodName
	}
	if runFlag == "" {
		runFlag = defaultRunFlag(kind)
	}
	return
}

func defaultRunFlag(kind reflect.Kind) string {
	if kind == reflect.Bool {
		return "true"
	} else if kind == reflect.String {
		return "notempty"
	}
	return ""
}
