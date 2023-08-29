package alias

import (
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/internal/proto"
	"github.com/gogf/gf/v2/frame/g"
)

type SetReq struct {
	g.Meta `path:"/alias/set" tags:"Alias" method:"Post" summary:"指定集合别名，新增/修改"`
	proto.SetAliasRequest
}

type SetRes struct {
	proto.UpdateAliasResponse
}

type DropReq struct {
	g.Meta `path:"/alias/drop" tags:"Alias" method:"Post" summary:"删除集合别名"`
	proto.DropAliasRequest
}

type DropRes struct {
	proto.UpdateAliasResponse
}

type DescribeReq struct {
	g.Meta `path:"/alias/describe" tags:"Alias" method:"Post" summary:"根据别名查找对应的集合信息"`
	proto.GetAliasRequest
}

type DescribeRes struct {
	proto.GetAliasResponse
}

type ListReq struct {
	g.Meta   `path:"/alias/drop" tags:"Alias" method:"Post" summary:"列举指定db下的所有别名信息"`
	Database string `json:"database"`
}

type ListRes struct {
	proto.GetAliasResponse
}
