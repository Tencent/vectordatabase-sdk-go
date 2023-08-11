package database

import (
	"vectordb-sdk-go/internal/proto"

	"github.com/gogf/gf/v2/frame/g"
)

type CreateReq struct {
	g.Meta   `path:"/database/create" tags:"Database" method:"Post" summary:"创建database，如果database已经存在返回成功"`
	Database string `json:"database,omitempty"`
}

type CreateRes struct {
	proto.DatabaseResponse
}

type DropReq struct {
	g.Meta   `path:"/database/drop" tags:"Database" method:"Post" summary:"删除database，并删除database中的所有collection以及数据，如果database不经存在返回成本"`
	Database string `json:"database,omitempty"`
}

type DropRes struct {
	proto.DatabaseResponse
}

type ListReq struct {
	g.Meta `path:"/database/list" tags:"Database" method:"Get" summary:"查询database列表"`
}

type ListRes struct {
	proto.DatabaseResponse
}
