package model

import (
	"context"
	"time"
)

// SdkClient the http client interface
type SdkClient interface {
	Close()
	Request(ctx context.Context, req, res interface{}) error
	WithTimeout(d time.Duration)
	Debug(v bool)
}

// DatabaseInterface database api
type DatabaseInterface interface {
	SdkClient
	CreateDatabase(ctx context.Context, name string) (*Database, error)
	DropDatabase(ctx context.Context, name string) error
	ListDatabase(ctx context.Context) (databases []*Database, err error)
	Database(name string) *Database
}

// CollectionInterface collection api
type CollectionInterface interface {
	SdkClient
	CreateCollection(ctx context.Context, name string, shardNum, replicasNum uint32, description string, indexes Indexes, embedding *Embedding) (*Collection, error)
	DescribeCollection(ctx context.Context, name string) (*Collection, error)
	DropCollection(ctx context.Context, collectionName string) (err error)
	TruncateCollection(ctx context.Context, name string) (affectedCount int, err error)
	ListCollection(ctx context.Context) ([]*Collection, error)
	Collection(name string) *Collection
}

type AliasInterface interface {
	SdkClient
	AliasSet(ctx context.Context, collectionName, aliasName string) (int, error)
	AliasDelete(ctx context.Context, aliasName string) (int, error)
	AliasDescribe(ctx context.Context, aliasName string) (*Alias, error)
	AliasList(ctx context.Context) ([]*Alias, error)
}

type IndexInterface interface {
	SdkClient
	IndexRebuild(ctx context.Context, collectionName string, dropBeforeRebuild bool, throttle int) error
}

// DocumentInterface document api
type DocumentInterface interface {
	SdkClient
	Upsert(ctx context.Context, documents []Document, buidIndex bool) (err error)
	Query(ctx context.Context, documentIds []string, filter *Filter, readConsistency string, retrieveVector bool, outputFields []string, offset, limit int64) (docs []Document, count uint64, err error)
	Search(ctx context.Context, vectors [][]float32, filter *Filter, readConsistency ReadConsistency, searchParam *SearchParams, retrieveVector bool, outputFields []string, limit int) ([][]Document, error)
	SearchById(ctx context.Context, documentIds []string, filter *Filter, readConsistency ReadConsistency, searchParam *SearchParams, retrieveVector bool, outputFields []string, limit int) ([][]Document, error)
	SearchByText(ctx context.Context, text map[string][]string, filter *Filter, readConsistency ReadConsistency, searchParam *SearchParams, retrieveVector bool, outputFields []string, limit int) ([][]Document, error)
	Delete(ctx context.Context, documentIds []string, filter *Filter) (err error)
	Update(ctx context.Context, queryIds []string, queryFilter *Filter, updateVector []float32, updateFields map[string]Field) (uint64, error)
}

type VectorDBClient interface {
	DatabaseInterface
}

// Database wrap the database parameters and collection interface to operating the collection api
type Database struct {
	CollectionInterface
	AliasInterface
	IndexInterface
	DatabaseName string
}

func (d *Database) Debug(v bool) {
	d.CollectionInterface.Debug(v)
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
	Indexes        Indexes
	IndexStatus    IndexStatus
	Embedding      Embedding
	Description    string
	Size           uint64
	CreateTime     time.Time
}
