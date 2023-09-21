package index

import (
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/internal/engine/api"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/internal/proto"
)

type RebuildReq struct {
	api.Meta `path:"/index/rebuild" tags:"Index" method:"Post" summary:"重建整个collection的所有索引"`
	proto.RebuildIndexRequest
}

type RebuildRes struct {
	proto.RebuildIndexResponse
}
