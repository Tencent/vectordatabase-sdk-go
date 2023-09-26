package entry

import (
	"context"
	"time"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/model"
)

// CollectionInterface collection api
type CollectionInterface interface {
	SdkClient
	CreateCollection(ctx context.Context, name string, shardNum, replicasNum uint32, description string, indexes model.Indexes, option *CreateCollectionOption) (*Collection, error)
	DescribeCollection(ctx context.Context, name string, option *DescribeCollectionOption) (*Collection, error)
	DropCollection(ctx context.Context, collectionName string, option *DropCollectionOption) (result *CollectionResult, err error)
	TruncateCollection(ctx context.Context, name string, option *TruncateCollectionOption) (result *CollectionResult, err error)
	ListCollection(ctx context.Context, option *ListCollectionOption) ([]*Collection, error)
	Collection(name string) *Collection
}

// Collection wrap the collection parameters and document interface to operating the document api
type Collection struct {
	DocumentInterface
	DatabaseName   string
	CollectionName string
	DocumentCount  int64
	Alias          []string
	ShardNum       uint32
	ReplicasNum    uint32
	Indexes        model.Indexes
	IndexStatus    model.IndexStatus
	Embedding      model.Embedding
	Description    string
	Size           uint64
	CreateTime     time.Time
}

type CollectionResult struct {
	AffectedCount int
}

type CreateCollectionOption struct {
	Embedding *model.Embedding
}

type DescribeCollectionOption struct{}

type DropCollectionOption struct{}

type TruncateCollectionOption struct{}

type ListCollectionOption struct{}
