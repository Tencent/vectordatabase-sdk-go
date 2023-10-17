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
	BuildIndex *bool
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

type UploadDocumentOption struct {
	FileType FileType
	MetaData map[string]string
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
