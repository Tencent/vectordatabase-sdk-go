package entity

import (
	"fmt"
	"reflect"
	"strconv"
)

type DocumentResult struct {
	AffectedCount int
	Total         int
}

type UpsertDocumentOption struct {
	BuildIndex bool
}

type QueryDocumentOption struct {
	Filter         *Filter
	RetrieveVector bool
	OutputFields   []string
	Offset         int64
	Limit          int64
}

type SearchDocumentOption struct {
	Filter         *Filter
	Params         *SearchDocParams
	RetrieveVector bool
	OutputFields   []string
	Limit          int64
}

type SearchDocParams struct {
	Nprobe uint32  `json:"nprobe,omitempty"` // 搜索时查找的聚类数量，使用索引默认值即可
	Ef     uint32  `json:"ef,omitempty"`     // HNSW
	Radius float32 `json:"radius,omitempty"` // 距离阈值,范围搜索时有效
}

type DeleteDocumentOption struct {
	DocumentIds []string
	Filter      *Filter
}

type UpdateDocumentOption struct {
	QueryIds     []string
	QueryFilter  *Filter
	UpdateVector []float32
	UpdateFields map[string]Field
}

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
	return fmt.Sprintf("%v", f.Val)
}

func (f Field) Int() int64 {
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

func (f Field) Float() float64 {
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
