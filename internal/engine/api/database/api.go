package database

import (
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/internal/engine/api"
)

// CreateReq create database request
type CreateReq struct {
	api.Meta `path:"/database/create" tags:"Database" method:"Post" summary:"创建database，如果database已经存在返回成功"`
	Database string `json:"database,omitempty"`
}

// CreateRes create database response
type CreateRes struct {
	Code          int32    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg           string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Redirect      string   `protobuf:"bytes,3,opt,name=redirect,proto3" json:"redirect,omitempty"`
	Databases     []string `protobuf:"bytes,4,rep,name=databases,proto3" json:"databases,omitempty"`
	AffectedCount int32    `protobuf:"varint,5,opt,name=affectedCount,proto3" json:"affectedCount,omitempty"`
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
	Code          int32    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg           string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Redirect      string   `protobuf:"bytes,3,opt,name=redirect,proto3" json:"redirect,omitempty"`
	Databases     []string `protobuf:"bytes,4,rep,name=databases,proto3" json:"databases,omitempty"`
	AffectedCount int32    `protobuf:"varint,5,opt,name=affectedCount,proto3" json:"affectedCount,omitempty"`
}
