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

package alias

import "git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity/api"

type SetReq struct {
	api.Meta   `path:"/alias/set" tags:"Alias" method:"Post" summary:"指定集合别名，新增/修改"`
	Database   string `protobuf:"bytes,1,opt,name=database,proto3" json:"database,omitempty"`
	Collection string `protobuf:"bytes,2,opt,name=collection,proto3" json:"collection,omitempty"`
	Alias      string `protobuf:"bytes,3,opt,name=alias,proto3" json:"alias,omitempty"`
}

type SetRes struct {
	Code          int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg           string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Redirect      string `protobuf:"bytes,3,opt,name=redirect,proto3" json:"redirect,omitempty"`
	AffectedCount int32  `protobuf:"varint,4,opt,name=affectedCount,proto3" json:"affectedCount,omitempty"`
}

type DeleteReq struct {
	api.Meta `path:"/alias/delete" tags:"Alias" method:"Post" summary:"删除集合别名"`
	Database string `protobuf:"bytes,1,opt,name=database,proto3" json:"database,omitempty"`
	Alias    string `protobuf:"bytes,2,opt,name=alias,proto3" json:"alias,omitempty"`
}

type DeleteRes struct {
	Code          int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg           string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Redirect      string `protobuf:"bytes,3,opt,name=redirect,proto3" json:"redirect,omitempty"`
	AffectedCount int32  `protobuf:"varint,4,opt,name=affectedCount,proto3" json:"affectedCount,omitempty"`
}

type DescribeReq struct {
	api.Meta `path:"/alias/describe" tags:"Alias" method:"Post" summary:"根据别名查找对应的集合信息"`
	Database string `protobuf:"bytes,1,opt,name=database,proto3" json:"database,omitempty"`
	Alias    string `protobuf:"bytes,2,opt,name=alias,proto3" json:"alias,omitempty"`
}

type DescribeRes struct {
	Code     int32        `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg      string       `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Redirect string       `protobuf:"bytes,3,opt,name=redirect,proto3" json:"redirect,omitempty"`
	Aliases  []*AliasItem `protobuf:"bytes,4,rep,name=aliases,proto3" json:"aliases,omitempty"`
}

type AliasItem struct {
	Alias      string `protobuf:"bytes,1,opt,name=alias,proto3" json:"alias,omitempty"`
	Collection string `protobuf:"bytes,2,opt,name=collection,proto3" json:"collection,omitempty"`
}

type ListReq struct {
	api.Meta `path:"/alias/list" tags:"Alias" method:"Post" summary:"列举指定db下的所有别名信息"`
	Database string `json:"database"`
}

type ListRes struct {
	Code     int32        `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg      string       `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Redirect string       `protobuf:"bytes,3,opt,name=redirect,proto3" json:"redirect,omitempty"`
	Aliases  []*AliasItem `protobuf:"bytes,4,rep,name=aliases,proto3" json:"aliases,omitempty"`
}
