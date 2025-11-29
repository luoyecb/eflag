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
		if field.Anonymous {
			return false
		}
		if !isReflectType(field.Type, reflect.Bool, reflect.String) {
			return false
		}
		cmdStr, ok := field.Tag.Lookup("command")
		if !ok {
			return false
		}

		fieldKind := field.Type.Kind()

		methodName := field.Name + "Command"
		runFlag := ""
		if cmdStr != "" {
			// format: methodName,runFlag
			parts := strings.Split(cmdStr, ",")
			if p1 := strings.TrimSpace(parts[0]); p1 != "" {
				methodName = p1
			}
			if len(parts) == 2 {
				runFlag = strings.TrimSpace(parts[1])
			}
		}
		if runFlag == "" {
			// default runFlag
			if fieldKind == reflect.Bool {
				runFlag = "true"
			} else if fieldKind == reflect.String {
				runFlag = "notempty"
			}
		}

		// call command method
		rm := rv.MethodByName(methodName)
		if !rm.IsValid() {
			return false
		}
		if fieldKind == reflect.Bool {
			val := value.Bool()
			if (val && runFlag == "true") || (!val && runFlag == "false") {
				rm.Call(nil)
				return true
			}
		} else if fieldKind == reflect.String {
			sval := value.String()
			if (sval == "" && runFlag == "empty") || (sval != "" && runFlag == "notempty") {
				rm.Call(nil)
				return true
			}
		}
		return false
	})
	return nil
}
