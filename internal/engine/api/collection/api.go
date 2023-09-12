package collection

import (
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/internal/proto"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/model"

	"github.com/gogf/gf/v2/frame/g"
)

// CreateReq create collection request
type CreateReq struct {
	g.Meta `path:"/collection/create" tags:"Collection" method:"Post" summary:"创建collection"`
	proto.CreateCollectionRequest
	Embedding model.Embedding `json:"embedding"`
}

// CreateRes create collection response
type CreateRes struct {
	proto.CreateCollectionResponse
}

// DescribeReq get collection detail request
type DescribeReq struct {
	g.Meta `path:"/collection/describe" tags:"Collection" method:"Post" summary:"返回collection信息"`
	proto.DescribeCollectionRequest
}

// DescribeRes get collection detail response
type DescribeRes struct {
	proto.DescribeCollectionResponse
	Collection *DescribeCollectionItem `json:"collection"`
}

// DropReq delete collection request
type DropReq struct {
	g.Meta `path:"/collection/drop" tags:"Collection" method:"Post" summary:"删除collection，并删除collection中的所有文档，如果collectio不经存在返回失败"`
	proto.DropCollectionRequest
}

// DropReq delete collection response
type DropRes struct {
	proto.DropCollectionResponse
}

type ListReq struct {
	g.Meta `path:"/collection/list" tags:"Collection" method:"Post" summary:"列出指定database中的所有collection"`
	proto.ListCollectionsRequest
}

type ListRes struct {
	proto.ListCollectionsResponse
	Collections []*DescribeCollectionItem `json:"collections,omitempty"`
}

type DescribeCollectionItem struct {
	proto.CreateCollectionRequest
	Alias         []string     `json:"alias"`
	DocumentCount int64        `json:"documentCount,omitempty"`
	Embedding     EmbeddingRes `json:"embedding"`
}

type FlushReq struct {
	g.Meta `path:"/collection/truncate" tags:"Collection" method:"Post" summary:"清空 collection 中的所有数据和索引"`
	proto.FlushCollectionRequest
}

type FlushRes struct {
	proto.FlushCollectionResponse
}

type ModifyReq struct {
	g.Meta `path:"/collection/flush" tags:"Collection" method:"Post" summary:"清空 collection 中的所有数据和索引"`
	proto.UpdateCollectionRequest
}

type ModifyRes struct {
	proto.UpdateCollectionResponse
}

type EmbeddingRes struct {
	model.Embedding
	Status string
}
