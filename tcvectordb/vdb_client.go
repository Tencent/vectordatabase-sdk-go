package tcvectordb

import "context"

type VdbClient interface {
	Close()
	Debug(debugFlag bool)

	ExistsDatabase(ctx context.Context, databaseName string) (bool, error)
	CreateDatabaseIfNotExists(ctx context.Context, databaseName string) (*CreateDatabaseResult, error)
	CreateDatabase(ctx context.Context, databaseName string) (*CreateDatabaseResult, error)
	DropDatabase(ctx context.Context, databaseName string) (*DropDatabaseResult, error)
	ListDatabase(ctx context.Context) (result *ListDatabaseResult, err error)

	ExistsCollection(ctx context.Context, databaseName, collectionName string) (bool, error)
	CreateCollectionIfNotExists(ctx context.Context, databaseName, collectionName string, shardNum, replicasNum uint32,
		description string, indexes Indexes, params ...*CreateCollectionParams) (*Collection, error)
	CreateCollection(ctx context.Context, databaseName, collectionName string, shardNum, replicasNum uint32,
		description string, indexes Indexes, params ...*CreateCollectionParams) (*Collection, error)
	ListCollection(ctx context.Context, databaseName string) (result *ListCollectionResult, err error)
	DescribeCollection(ctx context.Context, databaseName, collectionName string) (result *DescribeCollectionResult, err error)
	DropCollection(ctx context.Context, databaseName, collectionName string) (result *DropCollectionResult, err error)
	TruncateCollection(ctx context.Context, databaseName, collectionName string) (result *TruncateCollectionResult, err error)

	Upsert(ctx context.Context, databaseName, collectionName string, documents interface{}, params ...*UpsertDocumentParams) (result *UpsertDocumentResult, err error)
	Query(ctx context.Context, databaseName, collectionName string, documentIds []string, params ...*QueryDocumentParams) (result *QueryDocumentResult, err error)
	Search(ctx context.Context, databaseName, collectionName string, vectors [][]float32, params ...*SearchDocumentParams) (result *SearchDocumentResult, err error)
	HybridSearch(ctx context.Context, databaseName, collectionName string, params HybridSearchDocumentParams) (result *SearchDocumentResult, err error)
	FullTextSearch(ctx context.Context, databaseName, collectionName string, params FullTextSearchParams) (result *SearchDocumentResult, err error)
	SearchById(ctx context.Context, databaseName, collectionName string, documentIds []string, params ...*SearchDocumentParams) (result *SearchDocumentResult, err error)
	SearchByText(ctx context.Context, databaseName, collectionName string, text map[string][]string, params ...*SearchDocumentParams) (result *SearchDocumentResult, err error)
	Delete(ctx context.Context, databaseName, collectionName string, param DeleteDocumentParams) (result *DeleteDocumentResult, err error)
	Update(ctx context.Context, databaseName, collectionName string, param UpdateDocumentParams) (result *UpdateDocumentResult, err error)
	Count(ctx context.Context, databaseName, collectionName string, params ...CountDocumentParams) (*CountDocumentResult, error)

	CreateUser(ctx context.Context, param CreateUserParams) error
	GrantToUser(ctx context.Context, param GrantToUserParams) error
	RevokeFromUser(ctx context.Context, param RevokeFromUserParams) error
	DescribeUser(ctx context.Context, param DescribeUserParams) (result *DescribeUserResult, err error)
	ListUser(ctx context.Context) (result *ListUserResult, err error)
	DropUser(ctx context.Context, param DropUserParams) error
	ChangePassword(ctx context.Context, param ChangePasswordParams) error

	UploadFile(ctx context.Context, databaseName, collectionName string, param UploadFileParams) (result *UploadFileResult, err error)
	GetImageUrl(ctx context.Context, databaseName, collectionName string, param GetImageUrlParams) (result *GetImageUrlResult, err error)

	RebuildIndex(ctx context.Context, databaseName, collectionName string, params ...*RebuildIndexParams) (result *RebuildIndexResult, err error)
	AddIndex(ctx context.Context, databaseName, collectionName string, params ...*AddIndexParams) (err error)
	DropIndex(ctx context.Context, databaseName, collectionName string, params DropIndexParams) (err error)
	ModifyVectorIndex(ctx context.Context, databaseName, collectionName string, param ModifyVectorIndexParam) (err error)

	Embedding(ctx context.Context, param EmbeddingParams) (result *EmbeddingResult, err error)
}
