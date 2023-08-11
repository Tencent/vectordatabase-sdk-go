package engine

import (
	"context"

	"git.woa.com/cloud_nosql/vectordb/vectordb-sdk-go/internal/client"
	"git.woa.com/cloud_nosql/vectordb/vectordb-sdk-go/model"
)

type DatabaseInterface interface {
	client.SdkClient
	CreateDatabase(ctx context.Context, name string) (*Database, error)
	DropDatabase(ctx context.Context, name string) error
	ListDatabase(ctx context.Context) (databases []*Database, err error)
	Database(name string) *Database
}

type CollectionInterface interface {
	client.SdkClient
	CreateCollection(ctx context.Context, name string, shardNum, replicasNum uint32, description string, indexes model.Indexes) (*Collection, error)
	DescribeCollection(ctx context.Context, name string) (*Collection, error)
	DropCollection(ctx context.Context, collectionName string) (err error)
	ListCollection(ctx context.Context) ([]*Collection, error)
	Collection(name string) *Collection
}

type DocumentInterface interface {
	client.SdkClient
	Upsert(ctx context.Context, documents []model.Document, buidIndex bool) (err error)
	Query(ctx context.Context, documentIds []string, retrieveVector bool) ([]model.Document, error)
	Search(ctx context.Context, vectors [][]float32, filter *model.Filter, hnswParam *model.HNSWParam, retrieveVector bool, limit int) ([][]model.Document, error)
	SearchById(ctx context.Context, documentIds []string, filter *model.Filter, hnswParam *model.HNSWParam, retrieveVector bool, limit int) ([][]model.Document, error)
	Delete(ctx context.Context, documentIds []string) (err error)
}

type VectorDBClient interface {
	DatabaseInterface
}

type Database struct {
	CollectionInterface
	DatabaseName string
}

type Collection struct {
	DocumentInterface
	DatabaseName   string
	CollectionName string
	ShardNum       uint32
	ReplicasNum    uint32
	Indexes        model.Indexes
	Description    string
	Size           uint64
	CreateTime     string
}

func VectorDB(sdkClient client.SdkClient) VectorDBClient {
	databaseImpl := new(implementerDatabase)
	databaseImpl.SdkClient = sdkClient
	return databaseImpl
}
