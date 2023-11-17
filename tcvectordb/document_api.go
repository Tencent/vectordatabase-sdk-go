// Copyright (C) 2023 Tencent Cloud.
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the vectordb-sdk-java), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is furnished
// to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
// PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
// SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package tcvectordb

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
	BuildIndex *bool
}

type UpsertDocumentResult struct {
	AffectedCount int
}

type QueryDocumentOption struct {
	Filter         *Filter
	RetrieveVector bool
	OutputFields   []string
	Offset         int64
	Limit          int64
}

type QueryDocumentResult struct {
	Documents     []Document
	AffectedCount int
	Total         uint64
}

type SearchDocumentOption struct {
	Filter         *Filter
	Params         *SearchDocParams
	RetrieveVector bool
	OutputFields   []string
	Limit          int64
}

type SearchDocumentResult struct {
	Documents [][]Document
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

type DeleteDocumentResult struct {
	AffectedCount int
}

type UpdateDocumentOption struct {
	QueryIds     []string
	QueryFilter  *Filter
	UpdateVector []float32
	UpdateFields map[string]Field
}

type UpdateDocumentResult struct {
	AffectedCount int
}

type Document struct {
	Id     string
	Vector []float32
	// omitempty when upsert
	Score  float32 `json:"_,omitempty"`
	Fields map[string]Field
}

type Field struct {
	Val interface{} `json:"val,omitempty"`
}

func (f Field) String() string {
	return fmt.Sprintf("%v", f.Val)
}

func (f Field) StringArray() []string {
	t := reflect.TypeOf(f.Val)
	if t.Kind() != reflect.Slice && t.Kind() != reflect.Array {
		return nil
	}
	v := reflect.ValueOf(f.Val)
	res := make([]string, v.Len())
	for i := 0; i < v.Len(); i++ {
		res[i], _ = v.Index(i).Interface().(string)
	}
	return res
}

func (f Field) Uint64Array() []uint64 {
	t := reflect.TypeOf(f.Val)
	if t.Kind() != reflect.Slice && t.Kind() != reflect.Array {
		return nil
	}
	v := reflect.ValueOf(f.Val)
	res := make([]uint64, v.Len())
	for i := 0; i < v.Len(); i++ {
		switch v.Index(i).Kind() {
		case reflect.Uint, reflect.Uint64:
			res[i] = v.Index(i).Uint()
		case reflect.Int, reflect.Int64:
			res[i] = uint64(v.Index(i).Int())
		}
	}
	return res
}

func (f Field) Uint64() uint64 {
	switch v := f.Val.(type) {
	case int, int8, int16, int32, int64:
		return uint64(reflect.ValueOf(v).Int())
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(v).Uint()
	case string:
		n, _ := strconv.ParseUint(v, 10, 64)
		return n
	case float32, float64:
		return uint64(reflect.ValueOf(v).Float())
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

func (f Field) Type() FieldType {
	switch f.Val.(type) {
	case int, int8, int16, int32, int64:
		return Uint64
	case uint, uint8, uint16, uint32, uint64:
		return Uint64
	case string:
		return String
	case []string, []uint64, []int64, []int, []uint:
		return Array
	}
	return ""
}
