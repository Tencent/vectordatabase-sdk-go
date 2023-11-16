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

type VectorDBClient struct {
	DatabaseInterface
}

// DatabaseInterface database api
type DatabaseInterface interface {
	SdkClient
	CreateDatabase(ctx context.Context, name string, options ...*CreateDatabaseOption) (*CreateDatabaseResult, error)
	DropDatabase(ctx context.Context, name string, options ...*DropDatabaseOption) (*DropDatabaseResult, error)
	ListDatabase(ctx context.Context, options ...*ListDatabaseOption) (result *ListDatabaseResult, err error)
	CreateAIDatabase(ctx context.Context, name string, options ...*CreateAIDatabaseOption) (result *CreateAIDatabaseResult, err error)
	DropAIDatabase(ctx context.Context, name string, options ...*DropAIDatabaseOption) (result *DropAIDatabaseResult, err error)
	Database(name string) *Database
	AIDatabase(name string) *AIDatabase
}

// CollectionInterface collection api
type CollectionInterface interface {
	SdkClient
	CreateCollection(ctx context.Context, name string, shardNum, replicasNum uint32, description string,
		indexes Indexes, options ...*CreateCollectionOption) (*Collection, error)
	DescribeCollection(ctx context.Context, name string, options ...*DescribeCollectionOption) (result *DescribeCollectionResult, err error)
	DropCollection(ctx context.Context, name string, options ...*DropCollectionOption) (result *DropCollectionResult, err error)
	TruncateCollection(ctx context.Context, name string, options ...*TruncateCollectionOption) (result *TruncateCollectionResult, err error)
	ListCollection(ctx context.Context, options ...*ListCollectionOption) (result *ListCollectionResult, err error)
	Collection(name string) *Collection
}

type AICollectionInterface interface {
	SdkClient
	CreateCollection(ctx context.Context, name string, options ...*CreateAICollectionOption) (result *CreateAICollectionResult, err error)
	DescribeCollection(ctx context.Context, name string, options ...*DescribeAICollectionOption) (result *DescribeAICollectionResult, err error)
	DropCollection(ctx context.Context, name string, options ...*DropAICollectionOption) (result *DropAICollectionResult, err error)
	TruncateCollection(ctx context.Context, name string, options ...*TruncateAICollectionOption) (result *TruncateAICollectionResult, err error)
	ListCollection(ctx context.Context, options ...*ListAICollectionOption) (result *ListAICollectionResult, err error)
	Collection(name string) *AICollection
}

type AliasInterface interface {
	SdkClient
	SetAlias(ctx context.Context, collectionName, aliasName string, options ...*SetAliasOption) (result *SetAliasResult, err error)
	DeleteAlias(ctx context.Context, aliasName string, options ...*DeleteAliasOption) (result *DeleteAliasResult, err error)
}

type AIAliasInterface interface {
	SdkClient
	SetAlias(ctx context.Context, collectionName, aliasName string, options ...*SetAIAliasOption) (result *SetAIAliasResult, err error)
	DeleteAlias(ctx context.Context, aliasName string, options ...*DeleteAIAliasOption) (result *DeleteAIAliasResult, err error)
}

type IndexInterface interface {
	SdkClient
	IndexRebuild(ctx context.Context, collectionName string, options ...*IndexRebuildOption) (result *IndexReBuildResult, err error)
}

// DocumentInterface document api
type DocumentInterface interface {
	SdkClient
	Upsert(ctx context.Context, documents []Document, options ...*UpsertDocumentOption) (result *UpsertDocumentResult, err error)
	Query(ctx context.Context, documentIds []string, options ...*QueryDocumentOption) (result *QueryDocumentResult, err error)
	Search(ctx context.Context, vectors [][]float32, options ...*SearchDocumentOption) (result *SearchDocumentResult, err error)
	SearchById(ctx context.Context, documentIds []string, options ...*SearchDocumentOption) (*SearchDocumentResult, error)
	SearchByText(ctx context.Context, text map[string][]string, options ...*SearchDocumentOption) (*SearchDocumentResult, error)
	Delete(ctx context.Context, options ...*DeleteDocumentOption) (result *DeleteDocumentResult, err error)
	Update(ctx context.Context, options ...*UpdateDocumentOption) (*UpdateDocumentResult, error)
}

type AIDocumentInterface interface {
	SdkClient
	Query(ctx context.Context, options ...*QueryAIDocumentOption) (*QueryAIDocumentsResult, error)
	Search(ctx context.Context, text string, options ...*SearchAIDocumentOption) (*SearchAIDocumentResult, error)
	Delete(ctx context.Context, options ...*DeleteAIDocumentOption) (*DeleteAIDocumentResult, error)
	Update(ctx context.Context, options ...*UpdateAIDocumentOption) (*UpdateAIDocumentResult, error)
	Upload(ctx context.Context, localFilePath string, options ...*UploadAIDocumentOption) (*UploadAIDocumentResult, error)
	GetCosTmpSecret(ctx context.Context, localFilePath string, options ...*GetCosTmpSecretOption) (*GetCosTmpSecretResult, error)
}
