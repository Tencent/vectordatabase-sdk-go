package database

import (
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/internal/proto"

	"github.com/gogf/gf/v2/frame/g"
)

// CreateReq create database request
type CreateReq struct {
	g.Meta   `path:"/database/create" tags:"Database" method:"Post" summary:"创建database，如果database已经存在返回成功"`
	Database string `json:"database,omitempty"`
}

// CreateRes create database response
type CreateRes struct {
	proto.DatabaseResponse
}

// DropReq drop database request
type DropReq struct {
	g.Meta   `path:"/database/drop" tags:"Database" method:"Post" summary:"删除database，并删除database中的所有collection以及数据，如果database不经存在返回成本"`
	Database string `json:"database,omitempty"`
}

// DropRes drop database response
type DropRes struct {
	proto.DatabaseResponse
}

// ListReq get database list request
type ListReq struct {
	g.Meta `path:"/database/list" tags:"Database" method:"Get" summary:"查询database列表"`
}

// ListRes get database list response
type ListRes struct {
	proto.DatabaseResponse
}
