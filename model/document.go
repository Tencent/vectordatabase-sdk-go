package model

import (
	"github.com/gogf/gf/v2/util/gconv"
)

type Document struct {
	Id     string
	Vector []float32
	// omitempty when upsert
	Score  float32 `json:"_,omitempty"`
	Fields map[string]Field
}

type Field struct {
	Val interface{}
}

func (f Field) String() string {
	return gconv.String(f.Val)
}

func (f Field) Int() int64 {
	return gconv.Int64(f.Val)
}

func (f Field) Float() float64 {
	return gconv.Float64(f.Val)
}
