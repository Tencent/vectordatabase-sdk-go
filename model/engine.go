package model

import (
	"context"
	"time"
)

type SdkClient interface {
	Close()
	Request(ctx context.Context, req, res interface{}) error
	WithTimeout(d time.Duration)
	Debug(v bool)
}

type DatabaseInterface interface {
	SdkClient
	CreateDatabase(ctx context.Context, name string) (*Database, error)
	DropDatabase(ctx context.Context, name string) error
	ListDatabase(ctx context.Context) (databases []*Database, err error)
	Database(name string) *Database
}

type CollectionInterface interface {
	SdkClient
	CreateCollection(ctx context.Context, name string, shardNum, replicasNum uint32, description string, indexes Indexes) (*Collection, error)
	DescribeCollection(ctx context.Context, name string) (*Collection, error)
	DropCollection(ctx context.Context, collectionName string) (err error)
	ListCollection(ctx context.Context) ([]*Collection, error)
	Collection(name string) *Collection
}

type DocumentInterface interface {
	SdkClient
	Upsert(ctx context.Context, documents []Document, buidIndex bool) (err error)
	Query(ctx context.Context, documentIds []string, retrieveVector bool) ([]Document, error)
	Search(ctx context.Context, vectors [][]float32, filter *Filter, hnswParam *HNSWParam, retrieveVector bool, limit int) ([][]Document, error)
	SearchById(ctx context.Context, documentIds []string, filter *Filter, hnswParam *HNSWParam, retrieveVector bool, limit int) ([][]Document, error)
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
	Indexes        Indexes
	Description    string
	Size           uint64
	CreateTime     string
}
