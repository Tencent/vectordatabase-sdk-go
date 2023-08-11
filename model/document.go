package model

import (
	"fmt"
	"reflect"
	"strconv"
)

type Document struct {
	Id     string
	Vector []float32
	Score  float32
	Fields map[string]Field
}

type Field struct {
	Val interface{}
}

func (f *Field) String() string {
	return fmt.Sprintf("%v", f.Val)
}

func (f *Field) Int() int64 {
	switch v := f.Val.(type) {
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(v).Int()
	case uint, uint8, uint16, uint32, uint64:
		return int64(reflect.ValueOf(v).Uint())
	case string:
		n, _ := strconv.ParseInt(v, 10, 64)
		return n
	case float32, float64:
		return int64(reflect.ValueOf(v).Float())
	}
	return 0
}

func (f *Field) Float() float64 {
	switch v := f.Val.(type) {
	case int, int8, int16, int32, int64:
		return float64(reflect.ValueOf(v).Int())
	case uint, uint8, uint16, uint32, uint64:
		return float64(reflect.ValueOf(v).Uint())
	case string:
		n, _ := strconv.ParseFloat(v, 64)
		return n
	case float32, float64:
		return reflect.ValueOf(v).Float()
	}
	return 0
}
