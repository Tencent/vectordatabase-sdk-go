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

package database

import "github.com/tencent/vectordatabase-sdk-go/tcvectordb/api"

// CreateReq create database request
type CreateReq struct {
	api.Meta `path:"/database/create" tags:"Database" method:"Post" summary:"创建database，如果database已经存在返回成功"`
	Database string `json:"database,omitempty"`
}

// CreateRes create database response
type CreateRes struct {
	Code          int32    `json:"code,omitempty"`
	Msg           string   `json:"msg,omitempty"`
	Redirect      string   `json:"redirect,omitempty"`
	Databases     []string `json:"databases,omitempty"`
	AffectedCount int      `json:"affectedCount,omitempty"`
}

// DropReq drop database request
type DropReq struct {
	api.Meta `path:"/database/drop" tags:"Database" method:"Post" summary:"删除database，并删除database中的所有collection以及数据，如果database不经存在返回成本"`
	Database string `json:"database,omitempty"`
}

// DropRes drop database response
type DropRes struct {
	Code          int32    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg           string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Redirect      string   `protobuf:"bytes,3,opt,name=redirect,proto3" json:"redirect,omitempty"`
	Databases     []string `protobuf:"bytes,4,rep,name=databases,proto3" json:"databases,omitempty"`
	AffectedCount int32    `protobuf:"varint,5,opt,name=affectedCount,proto3" json:"affectedCount,omitempty"`
}

// ListReq get database list request
type ListReq struct {
	api.Meta `path:"/database/list" tags:"Database" method:"Get" summary:"查询database列表"`
}

// ListRes get database list response
type ListRes struct {
	Code          int32                    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg           string                   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Redirect      string                   `protobuf:"bytes,3,opt,name=redirect,proto3" json:"redirect,omitempty"`
	Databases     []string                 `protobuf:"bytes,4,rep,name=databases,proto3" json:"databases,omitempty"`
	AffectedCount int32                    `protobuf:"varint,5,opt,name=affectedCount,proto3" json:"affectedCount,omitempty"`
	Info          map[string]*DatabaseInfo `json:"info,omitempty"`
}

type DatabaseInfo struct {
	CreateTime string `json:"createTime,omitempty"`
	DbType     string `json:"dbType,omitempty"`
}
