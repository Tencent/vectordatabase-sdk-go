package ai_database

import "git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb/api"

type CreateReq struct {
	api.Meta `path:"/ai/database/create" tags:"Database" method:"Post" summary:"创建ai database，如果database已经存在返回成功"`
	Database string `json:"database,omitempty"`
}

type CreateRes struct {
	api.CommonRes
	Databases     []string `json:"databases,omitempty"`
	AffectedCount int32    `json:"affectedCount,omitempty"`
}

type DropReq struct {
	api.Meta `path:"/ai/database/drop" tags:"Database" method:"Post" summary:"删除ai database，并删除database中的所有collection以及数据，如果database不经存在返回成本"`
	Database string `json:"database,omitempty"`
}

// DropRes drop database response
type DropRes struct {
	api.CommonRes
	Databases     []string `protobuf:"bytes,4,rep,name=databases,proto3" json:"databases,omitempty"`
	AffectedCount int32    `protobuf:"varint,5,opt,name=affectedCount,proto3" json:"affectedCount,omitempty"`
}

type ListReq struct {
	api.Meta `path:"/database/list" tags:"Database" method:"Get" summary:"查询database列表"`
}

type ListRes struct {
	api.CommonRes
	Databases []string                 `json:"databases,omitempty"`
	Info      map[string]*DatabaseInfo `json:"info,omitempty"`
}

type DatabaseInfo struct {
	CreateTime string `json:"createTime,omitempty"`
	DbType     string `json:"dbType,omitempty"`
}
