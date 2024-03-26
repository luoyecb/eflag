package eflag

import (
	"errors"
	"reflect"
	"strings"
)

func RunCommand(v interface{}) error {
	if !isStructPtr(v) {
		return errors.New("Must be a pointer to a struct type")
	}

	ReflectVisitStructField(v, func(rv reflect.Value, field reflect.StructField, value reflect.Value) bool {
		if field.Anonymous || field.Type.Kind() != reflect.Bool {
			return false
		}
		cmdStr, ok := field.Tag.Lookup("command")
		if !ok {
			return false
		}

		methodName := field.Name + "Command"
		runFlag := "true"
		if cmdStr != "" {
			parts := strings.Split(cmdStr, ",")
			if p1 := strings.TrimSpace(parts[0]); p1 != "" {
				methodName = p1
			}
			if len(parts) == 2 {
				runFlag = strings.TrimSpace(parts[1])
			}
		}

		if rm := rv.MethodByName(methodName); rm.IsValid() {
			val := value.Bool()
			if (val && runFlag == "true") || (!val && runFlag == "false") {
				rm.Call(nil)
				return true
			}
		}
		return false
	})
	return nil
}
