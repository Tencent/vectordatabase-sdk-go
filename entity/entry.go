// Copyright (C) 2023 Tencent Cloud.
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the vectordb-sdk-java), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is furnished
// to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
// PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
// SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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
