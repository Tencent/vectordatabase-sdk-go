package entity

import (
	"context"
)

type VectorDBClient interface {
	DatabaseInterface
}

// DatabaseInterface database api
type DatabaseInterface interface {
	SdkClient
	CreateDatabase(ctx context.Context, name string, option *CreateDatabaseOption) (*Database, error)
	DropDatabase(ctx context.Context, name string, option *DropDatabaseOption) (*DatabaseResult, error)
	ListDatabase(ctx context.Context, option *ListDatabaseOption) (databases []*Database, err error)
	Database(name string) *Database
}

// CollectionInterface collection api
type CollectionInterface interface {
	SdkClient
	CreateCollection(ctx context.Context, name string, shardNum, replicasNum uint32, description string, indexes Indexes, option *CreateCollectionOption) (*Collection, error)
	DescribeCollection(ctx context.Context, name string, option *DescribeCollectionOption) (*Collection, error)
	DropCollection(ctx context.Context, collectionName string, option *DropCollectionOption) (result *CollectionResult, err error)
	TruncateCollection(ctx context.Context, name string, option *TruncateCollectionOption) (result *CollectionResult, err error)
	ListCollection(ctx context.Context, option *ListCollectionOption) ([]*Collection, error)
	Collection(name string) *Collection
}

type AliasInterface interface {
	SdkClient
	SetAlias(ctx context.Context, collectionName, aliasName string, option *SetAliasOption) (*AliasResult, error)
	DeleteAlias(ctx context.Context, aliasName string, option *DeleteAliasOption) (*AliasResult, error)
	DescribeAlias(ctx context.Context, aliasName string, option *DescribeAliasOption) (*AliasResult, error)
	ListAlias(ctx context.Context, option *ListAliasOption) ([]*AliasResult, error)
}

// DocumentInterface document api
type DocumentInterface interface {
	SdkClient
	Upsert(ctx context.Context, documents []Document, option *UpsertDocumentOption) (result *DocumentResult, err error)
	Query(ctx context.Context, documentIds []string, option *QueryDocumentOption) (docs []Document, result *DocumentResult, err error)
	Search(ctx context.Context, vectors [][]float32, option *SearchDocumentOption) ([][]Document, error)
	SearchById(ctx context.Context, documentIds []string, option *SearchDocumentOption) ([][]Document, error)
	SearchByText(ctx context.Context, text map[string][]string, option *SearchDocumentOption) ([][]Document, error)
	Delete(ctx context.Context, option *DeleteDocumentOption) (result *DocumentResult, err error)
	Update(ctx context.Context, option *UpdateDocumentOption) (result *DocumentResult, err error)
}
