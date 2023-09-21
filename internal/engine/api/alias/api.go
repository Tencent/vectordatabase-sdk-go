package alias

import (
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/internal/engine/api"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/internal/proto"
)

type SetReq struct {
	api.Meta `path:"/alias/set" tags:"Alias" method:"Post" summary:"指定集合别名，新增/修改"`
	proto.SetAliasRequest
}

type SetRes struct {
	proto.UpdateAliasResponse
}

type DeleteReq struct {
	api.Meta `path:"/alias/delete" tags:"Alias" method:"Post" summary:"删除集合别名"`
	proto.DropAliasRequest
}

type DeleteRes struct {
	proto.UpdateAliasResponse
}

type DescribeReq struct {
	api.Meta `path:"/alias/describe" tags:"Alias" method:"Post" summary:"根据别名查找对应的集合信息"`
	proto.GetAliasRequest
}

type DescribeRes struct {
	proto.GetAliasResponse
}

type ListReq struct {
	api.Meta `path:"/alias/list" tags:"Alias" method:"Post" summary:"列举指定db下的所有别名信息"`
	Database string `json:"database"`
}

type ListRes struct {
	proto.GetAliasResponse
}
