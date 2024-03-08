package utinterface

import (
	"encoding/json"
	"reflect"
	"strconv"
	"time"
)

func IsNil(value interface{}) (res bool) {
	return (value == nil || (reflect.TypeOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil()))
}

func ToString(value interface{}) (res string) {
	if !IsNil(value) {
		val := reflect.ValueOf(value)
		switch val.Kind() {
		case reflect.String:
			res = val.String()

		case reflect.Ptr:
			res = ToString(reflect.Indirect(val))

		default:
			switch valx := value.(type) {
			case []byte:
				res = string(valx)

			case time.Time:
				res = valx.Format(time.RFC3339Nano)

			default:
				byt, err := json.Marshal(value)
				if err == nil {
					res = string(byt)
				}
			}
		}
	}
	return
}

func ToInt(value interface{}, def int64) int64 {
	r, err := strconv.ParseInt(ToString(value), 10, 64)
	if err != nil {
		r = def
	}
	return r
}
