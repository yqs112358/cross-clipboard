package maputil

import (
	"fmt"
	"reflect"
	"strings"
)

// ToMapString converts a struct to a map using the struct tags.
// ignore empty and `-` tag value
func ToMapString(in interface{}, tag string) (map[string]string, error) {
	out := make(map[string]string)

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ToMapString only accepts struct, got %T", v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		if tagv := fi.Tag.Get(tag); tagv != "" && tagv != "-" {
			h, _ := head(tagv, ",")
			valueField := v.Field(i)

			if valueField.Type().Kind() == reflect.Ptr {
				if valueField.IsNil() {
					continue
				}
				valueField = valueField.Elem()
			}

			out[h] = fmt.Sprintf("%v", valueField.Interface())
		}
	}
	return out, nil
}

func head(str, sep string) (head string, tail string) {
	idx := strings.Index(str, sep)
	if idx < 0 {
		return str, ""
	}
	return str[:idx], str[idx+len(sep):]
}
