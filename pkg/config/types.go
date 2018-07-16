package config

import (
	"reflect"
)

func decodeBool(name string, config interface{}) (bool, error) {
	var result bool
	switch t := reflect.ValueOf(config); t.Kind() {
	case reflect.Bool:
		result = t.Bool()
	default:
		return result, errorUnsupportedType(name, t.Kind())
	}
	return result, nil
}

func decodeString(name string, config interface{}) (string, error) {
	var result string
	switch t := reflect.ValueOf(config); t.Kind() {
	case reflect.String:
		return ApplyTemplate(t.String())
	default:
		return result, errorUnsupportedType(name, t.Kind())
	}
	return result, nil
}

func decodeStringSlice(name string, config interface{}) ([]string, error) {
	var result []string
	switch t := reflect.ValueOf(config); t.Kind() {
	case reflect.String:
		decoded, err := decodeString(name, t.String())
		if err != nil {
			return result, err
		}
		result = []string{decoded}
	case reflect.Slice:
		for _, v := range t.Interface().([]interface{}) {
			decoded, err := decodeString(name, v)
			if err != nil {
				return result, err
			}
			result = append(result, decoded)
		}
	default:
		return result, errorUnsupportedType(name, t.Kind())
	}
	return result, nil
}
