package collection

import (
	"vectordb-sdk-go/internal/proto"

	"github.com/gogf/gf/v2/frame/g"
)

type CreateReq struct {
	g.Meta `path:"/collection/createCollection" tags:"Collection" method:"Post" summary:"创建collection"`
	proto.CreateCollectionRequest
}

type CreateRes struct {
	proto.CreateCollectionResponse
}

type DescribeReq struct {
	g.Meta `path:"/collection/describeCollection" tags:"Collection" method:"Post" summary:"返回collection信息"`
	proto.DescribeCollectionRequest
}

type DescribeRes struct {
	proto.DescribeCollectionResponse
}

type DropReq struct {
	g.Meta `path:"/collection/dropCollection" tags:"Collection" method:"Post" summary:"删除collection，并删除collection中的所有文档，如果collectio不经存在返回失败"`
	proto.DropCollectionRequest
}

type DropRes struct {
	proto.DropCollectionResponse
}

type ListReq struct {
	g.Meta `path:"/collection/listCollections" tags:"Collection" method:"Post" summary:"列出指定database中的所有collection"`
	proto.ListCollectionsRequest
}

type ListRes struct {
	proto.ListCollectionsResponse
}
