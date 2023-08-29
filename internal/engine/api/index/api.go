package index

import (
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/internal/proto"
	"github.com/gogf/gf/v2/frame/g"
)

type RebuildReq struct {
	g.Meta `path:"/index/rebuild" tags:"Index" method:"Post" summary:"重建整个collection的所有索引"`
	proto.RebuildIndexRequest
}

type RebuildRes struct {
	proto.RebuildIndexResponse
}
