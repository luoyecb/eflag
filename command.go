package eflag

import (
	"errors"
	"reflect"
	"strings"
)

const (
	COMMAND_FIELD_TAG_KEY       = "command"
	COMMAND_METHOD_NAME_KEY     = "Command"
	COMMAND_SUB_COMMAND_TAG_KEY = "sub_command"

	SUM_COMMAND_INDEX = 1
)

type CommandMode uint

const (
	COMMAND_MODE_OPTION  CommandMode = iota // Option command, eg: go run main.go -COMMAND
	COMMAND_MODE_SUB_CMD                    // Sub command, eg: go run main.go SUB_COMMAND
)

func runCommand(v interface{}, subCommandName string) error {
	if !isStructPtr(v) {
		return errors.New("Must be a pointer to a struct type")
	}
	if subCommandName == "" {
		ReflectVisitStructField(v, true, callCommand)
		return nil
	}

	ReflectVisitStructField(v, true, func(rv reflect.Value, field reflect.StructField, fieldValue reflect.Value) bool {
		tagStr, ok := field.Tag.Lookup(COMMAND_SUB_COMMAND_TAG_KEY)
		if !ok || tagStr != subCommandName {
			return false
		}
		method := rv.MethodByName(field.Name + COMMAND_METHOD_NAME_KEY)
		if method.IsValid() {
			method.Call(nil)
			return true
		}
		return false
	})
	return nil
}

func callCommand(rv reflect.Value, field reflect.StructField, fieldValue reflect.Value) bool {
	if !isReflectType(field.Type, reflect.Bool, reflect.String) {
		return false
	}
	cmdStr, ok := field.Tag.Lookup(COMMAND_FIELD_TAG_KEY)
	if !ok {
		return false
	}

	methodName, runFlag := parseCommand(cmdStr, field.Type.Kind(), field.Name)
	rm := rv.MethodByName(methodName + COMMAND_METHOD_NAME_KEY)
	if !rm.IsValid() {
		return false
	}
	return callMethod(rm, fieldValue, runFlag)
}

func callMethod(method reflect.Value, value reflect.Value, runFlag string) (called bool) {
	isRun := false
	switch value.Kind() {
	case reflect.Bool:
		b := value.Bool()
		isRun = (b && runFlag == "true") || (!b && runFlag == "false")
	case reflect.String:
		s := value.String()
		isRun = (s == "" && runFlag == "empty") || (s != "" && runFlag == "notempty")
	default:
	}
	if isRun {
		method.Call(nil)
	}
	return isRun
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
