package api

import (
	"reflect"
)

type Meta struct{}

func Path(s interface{}) string {
	reflectType := reflect.TypeOf(s)
	if reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}
	field, ok := reflectType.FieldByName("Meta")
	if !ok {
		return ""
	}
	return field.Tag.Get("path")
}

func Method(s interface{}) string {
	reflectType := reflect.TypeOf(s)
	if reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}
	field, ok := reflectType.FieldByName("Meta")
	if !ok {
		return ""
	}
	return field.Tag.Get("method")
}
