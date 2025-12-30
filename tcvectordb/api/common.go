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

package api

import (
	"reflect"
)

type Meta struct{}

type CommonRes struct {
	Code int32  `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

const (
	EventualConsistency = "eventualConsistency"
	StrongConsistency   = "strongConsistency"
)

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

type IndexColumn struct {
	FieldName        string       `json:"fieldName,omitempty"`
	FieldType        string       `json:"fieldType,omitempty"`
	FieldElementType string       `json:"fieldElementType,omitempty"`
	IndexType        string       `json:"indexType,omitempty"`
	Dimension        uint32       `json:"dimension,omitempty"`
	MetricType       string       `json:"metricType,omitempty"`
	DiskSwapEnabled  *bool        `json:"diskSwapEnabled,omitempty"`
	IndexedCount     uint64       `json:"indexedCount,omitempty"`
	Params           *IndexParams `json:"params,omitempty"`
	AutoId           string       `json:"autoId,omitempty"`
}

type IndexParams struct {
	M              uint32  `protobuf:"varint,1,opt,name=M,proto3" json:"M,omitempty"`
	EfConstruction uint32  `protobuf:"varint,2,opt,name=efConstruction,proto3" json:"efConstruction,omitempty"`
	Nprobe         uint32  `protobuf:"varint,3,opt,name=nprobe,proto3" json:"nprobe,omitempty"`
	Nlist          uint32  `protobuf:"varint,4,opt,name=nlist,proto3" json:"nlist,omitempty"`
	Bits           *uint32 `protobuf:"varint,5,opt,name=bits,proto3" json:"bits,omitempty"`
}

// [ParsingProcess] holds the parameters for parsing files.
//
// Fields:
//   - ParsingType:  (Optional) The type of parsing files, which can be set to AlgorithmParsing
//     or VisionModelParsing (default to AlgorithmParsing).
type ParsingProcess struct {
	ParsingType string `json:"parsingType,omitempty"`
}
