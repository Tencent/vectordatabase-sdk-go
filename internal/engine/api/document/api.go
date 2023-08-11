package document

import (
	"vectordb-sdk-go/internal/proto"

	"github.com/gogf/gf/v2/frame/g"
)

type UpsertReq struct {
	g.Meta `path:"/document/upsert" tags:"Document" method:"Post" summary:"插入一条文档数据"`
	proto.UpsertRequest
	Document Document
}

type UpsertRes struct {
	proto.UpsertResponse
}

type SearchReq struct {
	g.Meta `path:"/document/search" tags:"Document" method:"Post" summary:"向量查询接口，支持向量检索以及向量+标量混合检索"`
	proto.SearchRequest
	Search *SearchCond `json:"search,omitempty"`
}

type SearchRes struct {
	proto.SearchResponse
}

type SearchCond struct {
	proto.SearchCond
	Vectors [][]float32 `json:"vectors,omitempty"`
}

type QueryReq struct {
	g.Meta `path:"/document/query" tags:"Document" method:"Post" summary:"标量查询接口，当前仅支持主键id查询"`
	proto.QueryRequest
}

type QueryRes struct {
	proto.QueryResponse
}

type DeleteReq struct {
	g.Meta `path:"/document/delete" tags:"Document" method:"Post" summary:"删除指定id的文档,flat 索引不支持删除"`
	proto.DeleteRequest
}

type DeleteRes struct {
	proto.DeleteResponse
}

type Document struct {
	proto.Document
	Fields map[string]interface{}
}
