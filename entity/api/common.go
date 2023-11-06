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
	FieldName    string       `protobuf:"bytes,1,opt,name=fieldName,proto3" json:"fieldName,omitempty"`
	FieldType    string       `protobuf:"bytes,2,opt,name=fieldType,proto3" json:"fieldType,omitempty"`
	IndexType    string       `protobuf:"bytes,3,opt,name=indexType,proto3" json:"indexType,omitempty"`
	Dimension    uint32       `protobuf:"varint,4,opt,name=dimension,proto3" json:"dimension,omitempty"`
	MetricType   string       `protobuf:"bytes,5,opt,name=metricType,proto3" json:"metricType,omitempty"`
	IndexedCount uint64       `protobuf:"varint,6,opt,name=indexedCount,proto3" json:"indexedCount,omitempty"`
	Params       *IndexParams `protobuf:"bytes,8,opt,name=params,proto3" json:"params,omitempty"`
}

type IndexParams struct {
	M              uint32 `protobuf:"varint,1,opt,name=M,proto3" json:"M,omitempty"`
	EfConstruction uint32 `protobuf:"varint,2,opt,name=efConstruction,proto3" json:"efConstruction,omitempty"`
	Nprobe         uint32 `protobuf:"varint,3,opt,name=nprobe,proto3" json:"nprobe,omitempty"`
	Nlist          uint32 `protobuf:"varint,4,opt,name=nlist,proto3" json:"nlist,omitempty"`
}
