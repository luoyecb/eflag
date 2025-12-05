package eflag

import (
	"errors"
	"reflect"
	"strings"
)

const (
	COMMAND_FIELD_TAG_KEY   = "command"
	COMMAND_METHOD_NAME_KEY = "Command"
)

func runCommand(v interface{}) error {
	if !isStructPtr(v) {
		return errors.New("Must be a pointer to a struct type")
	}

	ReflectVisitStructField(v, func(rv reflect.Value, field reflect.StructField, value reflect.Value) bool {
		if field.Anonymous {
			return false
		}
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
		return callCommand(rm, value, runFlag)
	})
	return nil
}

func callCommand(method reflect.Value, value reflect.Value, runFlag string) (called bool) {
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
