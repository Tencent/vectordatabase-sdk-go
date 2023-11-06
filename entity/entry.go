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
	CreateDatabase(ctx context.Context, name string, option *CreateDatabaseOption) (*CreateDatabaseResult, error)
	DropDatabase(ctx context.Context, name string, option *DropDatabaseOption) (*DropDatabaseResult, error)
	ListDatabase(ctx context.Context, option *ListDatabaseOption) (result *ListDatabaseResult, err error)
	CreateAIDatabase(ctx context.Context, name string, option *CreateAIDatabaseOption) (result *CreateAIDatabaseResult, err error)
	DropAIDatabase(ctx context.Context, name string, option *DropAIDatabaseOption) (result *DropAIDatabaseResult, err error)
	Database(name string) *Database
	AIDatabase(name string) *AIDatabase
}

// CollectionInterface collection api
type CollectionInterface interface {
	SdkClient
	CreateCollection(ctx context.Context, name string, shardNum, replicasNum uint32, description string,
		indexes Indexes, option *CreateCollectionOption) (*Collection, error)
	DescribeCollection(ctx context.Context, name string, option *DescribeCollectionOption) (*DescribeCollectionResult, error)
	DropCollection(ctx context.Context, name string, option *DropCollectionOption) (*DropCollectionResult, error)
	TruncateCollection(ctx context.Context, name string, option *TruncateCollectionOption) (*TruncateCollectionResult, error)
	ListCollection(ctx context.Context, option *ListCollectionOption) (*ListCollectionResult, error)
	Collection(name string) *Collection
}

type AICollectionInterface interface {
	SdkClient
	CreateCollection(ctx context.Context, name string, option *CreateAICollectionOption) (*CreateAICollectionResult, error)
	DescribeCollection(ctx context.Context, name string, option *DescribeAICollectionOption) (*DescribeAICollectionResult, error)
	DropCollection(ctx context.Context, name string, option *DropAICollectionOption) (result *DropAICollectionResult, err error)
	TruncateCollection(ctx context.Context, name string, option *TruncateAICollectionOption) (result *TruncateAICollectionResult, err error)
	ListCollection(ctx context.Context, option *ListAICollectionOption) (*ListAICollectionResult, error)
	Collection(name string) *AICollection
}

type AliasInterface interface {
	SdkClient
	SetAlias(ctx context.Context, collectionName, aliasName string, option *SetAliasOption) (*SetAliasResult, error)
	DeleteAlias(ctx context.Context, aliasName string, option *DeleteAliasOption) (*DeleteAliasResult, error)
}

type AIAliasInterface interface {
	SdkClient
	SetAlias(ctx context.Context, collectionName, aliasName string, option *SetAIAliasOption) (*SetAIAliasResult, error)
	DeleteAlias(ctx context.Context, aliasName string, option *DeleteAIAliasOption) (*DeleteAIAliasResult, error)
}

type IndexInterface interface {
	SdkClient
	IndexRebuild(ctx context.Context, collectionName string, option *IndexRebuildOption) (*IndexReBuildResult, error)
}

// DocumentInterface document api
type DocumentInterface interface {
	SdkClient
	Upsert(ctx context.Context, documents []Document, option *UpsertDocumentOption) (result *UpsertDocumentResult, err error)
	Query(ctx context.Context, documentIds []string, option *QueryDocumentOption) (*QueryDocumentResult, error)
	Search(ctx context.Context, vectors [][]float32, option *SearchDocumentOption) (*SearchDocumentResult, error)
	SearchById(ctx context.Context, documentIds []string, option *SearchDocumentOption) (*SearchDocumentResult, error)
	SearchByText(ctx context.Context, text map[string][]string, option *SearchDocumentOption) (*SearchDocumentResult, error)
	Delete(ctx context.Context, option *DeleteDocumentOption) (result *DeleteDocumentResult, err error)
	Update(ctx context.Context, option *UpdateDocumentOption) (*UpdateDocumentResult, error)
}

type AIDocumentInterface interface {
	SdkClient
	Query(ctx context.Context, option *QueryAIDocumentOption) (*QueryAIDocumentsResult, error)
	Search(ctx context.Context, text string, option *SearchAIDocumentOption) (*SearchAIDocumentResult, error)
	Delete(ctx context.Context, option *DeleteAIDocumentOption) (*DeleteAIDocumentResult, error)
	Update(ctx context.Context, option *UpdateAIDocumentOption) (*UpdateAIDocumentResult, error)
	Upload(ctx context.Context, localFilePath string, option *UploadAIDocumentOption) (*UploadAIDocumentResult, error)
	GetCosTmpSecret(ctx context.Context, localFilePath string, option *GetCosTmpSecretOption) (*GetCosTmpSecretResult, error)
}
