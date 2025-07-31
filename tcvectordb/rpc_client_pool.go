package tcvectordb

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

type RpcClientPool struct {
	FlatInterface
	clients []*RpcClient
	index   int64 // Changed to int64 for atomic operations
	mux     sync.Mutex

	url      string
	username string
	key      string
	option   *ClientOption
}

func NewRpcClientPool(url, username, key string, option *ClientOption) (VdbClient, error) {
	if option == nil {
		option = &defaultOption
	}
	if option.RpcPoolSize == 0 {
		option.RpcPoolSize = defaultRpcPoolSize
	}
	clients := make([]*RpcClient, option.RpcPoolSize)
	for i := 0; i < option.RpcPoolSize; i++ {
		client, err := NewRpcClient(url, username, key, option)
		if err != nil {
			for j := 0; j < i; j++ {
				clients[i].Close()
			}
			return nil, fmt.Errorf("new rpc client for client pool failed. err: %v", err.Error())
		}
		clients[i] = client
	}
	return &RpcClientPool{
		clients:  clients,
		url:      url,
		username: username,
		key:      key,
		option:   option,
	}, nil
}

func (pool *RpcClientPool) getRpcClient() (*RpcClient, error) {

	// 直接使用模运算，避免大数字
	currentIndex := atomic.AddInt64(&pool.index, 1) % int64(pool.option.RpcPoolSize)
	if currentIndex < 0 {
		currentIndex = 0 // 处理负数情况
	}
	client := pool.clients[currentIndex]

	//println("get rpc client from pool, which index is ", currentIndex)
	state := client.cc.GetState().String()
	var err error
	if state == "SHUTDOWN" || state == "INVALID_STATE" {
		// Only lock when we need to replace a client
		pool.mux.Lock()
		// Double-check the state to avoid race conditions
		if client.cc.GetState().String() == "SHUTDOWN" || client.cc.GetState().String() == "INVALID_STATE" {
			oldClient := pool.clients[currentIndex]
			client, err = NewRpcClient(pool.url, pool.username, pool.key, pool.option)
			if err != nil {
				pool.mux.Unlock()
				return nil, err
			}
			pool.clients[currentIndex] = client
			//println("new rpc client from pool, which index is ", currentIndex)
			oldClient.Close()
		}
		pool.mux.Unlock()
	}

	return client, nil
}

func (pool *RpcClientPool) Close() {
	for _, client := range pool.clients {
		client.Close()
	}
}

func (pool *RpcClientPool) Debug(v bool) {
	for _, client := range pool.clients {
		client.httpImplementer.Debug(v)
		client.debug = v
	}
}

func (pool *RpcClientPool) ExistsDatabase(ctx context.Context, databaseName string) (bool, error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return false, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.ExistsDatabase(ctx, databaseName)
}

func (pool *RpcClientPool) CreateDatabaseIfNotExists(ctx context.Context, databaseName string) (*CreateDatabaseResult, error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return nil, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.CreateDatabaseIfNotExists(ctx, databaseName)
}

func (pool *RpcClientPool) CreateDatabase(ctx context.Context, databaseName string) (*CreateDatabaseResult, error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return nil, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.CreateDatabase(ctx, databaseName)
}

func (pool *RpcClientPool) DropDatabase(ctx context.Context, databaseName string) (*DropDatabaseResult, error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return nil, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.DropDatabase(ctx, databaseName)
}

func (pool *RpcClientPool) ListDatabase(ctx context.Context) (result *ListDatabaseResult, err error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return nil, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.ListDatabase(ctx)
}

func (pool *RpcClientPool) ExistsCollection(ctx context.Context, databaseName, collectionName string) (bool, error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return false, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	db := client.Database(databaseName)
	return db.ExistsCollection(ctx, collectionName)
}

func (pool *RpcClientPool) CreateCollectionIfNotExists(ctx context.Context, databaseName, collectionName string, shardNum,
	replicasNum uint32, description string,
	indexes Indexes, params ...*CreateCollectionParams) (*Collection, error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return nil, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	db := client.Database(databaseName)
	return db.CreateCollectionIfNotExists(ctx, collectionName, shardNum, replicasNum, description, indexes, params...)
}

func (pool *RpcClientPool) CreateCollection(ctx context.Context, databaseName, collectionName string, shardNum, replicasNum uint32, description string,
	indexes Indexes, params ...*CreateCollectionParams) (*Collection, error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return nil, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	db := client.Database(databaseName)
	return db.CreateCollection(ctx, collectionName, shardNum, replicasNum, description, indexes, params...)
}

func (pool *RpcClientPool) ListCollection(ctx context.Context, databaseName string) (result *ListCollectionResult, err error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return nil, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	db := client.Database(databaseName)
	return db.ListCollection(ctx)
}

func (pool *RpcClientPool) DescribeCollection(ctx context.Context, databaseName, collectionName string) (result *DescribeCollectionResult, err error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return nil, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	db := client.Database(databaseName)
	return db.DescribeCollection(ctx, collectionName)
}

func (pool *RpcClientPool) DropCollection(ctx context.Context, databaseName, collectionName string) (result *DropCollectionResult, err error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return nil, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	db := client.Database(databaseName)
	return db.DropCollection(ctx, collectionName)
}

func (pool *RpcClientPool) TruncateCollection(ctx context.Context, databaseName, collectionName string) (result *TruncateCollectionResult, err error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return nil, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	db := client.Database(databaseName)
	return db.TruncateCollection(ctx, collectionName)
}

func (pool *RpcClientPool) Upsert(ctx context.Context, databaseName, collectionName string,
	documents interface{}, params ...*UpsertDocumentParams) (result *UpsertDocumentResult, err error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return nil, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.Upsert(ctx, databaseName, collectionName, documents, params...)
}

func (pool *RpcClientPool) Query(ctx context.Context, databaseName, collectionName string,
	documentIds []string, params ...*QueryDocumentParams) (result *QueryDocumentResult, err error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return nil, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.Query(ctx, databaseName, collectionName, documentIds, params...)
}

func (pool *RpcClientPool) Search(ctx context.Context, databaseName, collectionName string,
	vectors [][]float32, params ...*SearchDocumentParams) (result *SearchDocumentResult, err error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return nil, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.Search(ctx, databaseName, collectionName, vectors, params...)
}

func (pool *RpcClientPool) HybridSearch(ctx context.Context, databaseName, collectionName string,
	params HybridSearchDocumentParams) (result *SearchDocumentResult, err error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return nil, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.HybridSearch(ctx, databaseName, collectionName, params)
}

func (pool *RpcClientPool) FullTextSearch(ctx context.Context, databaseName, collectionName string,
	params FullTextSearchParams) (result *SearchDocumentResult, err error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return nil, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.FullTextSearch(ctx, databaseName, collectionName, params)
}

func (pool *RpcClientPool) SearchById(ctx context.Context, databaseName, collectionName string,
	documentIds []string, params ...*SearchDocumentParams) (result *SearchDocumentResult, err error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return nil, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.SearchById(ctx, databaseName, collectionName, documentIds, params...)
}

func (pool *RpcClientPool) SearchByText(ctx context.Context, databaseName, collectionName string,
	text map[string][]string, params ...*SearchDocumentParams) (result *SearchDocumentResult, err error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return nil, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.SearchByText(ctx, databaseName, collectionName, text, params...)
}

func (pool *RpcClientPool) Delete(ctx context.Context, databaseName, collectionName string,
	param DeleteDocumentParams) (result *DeleteDocumentResult, err error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return nil, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.Delete(ctx, databaseName, collectionName, param)
}

func (pool *RpcClientPool) Update(ctx context.Context, databaseName, collectionName string,
	param UpdateDocumentParams) (result *UpdateDocumentResult, err error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return nil, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.Update(ctx, databaseName, collectionName, param)
}

func (pool *RpcClientPool) Count(ctx context.Context, databaseName, collectionName string,
	params ...CountDocumentParams) (*CountDocumentResult, error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return nil, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.Count(ctx, databaseName, collectionName, params...)
}

func (pool *RpcClientPool) CreateUser(ctx context.Context, param CreateUserParams) error {
	client, err := pool.getRpcClient()
	if err != nil {
		return fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.CreateUser(ctx, param)
}

func (pool *RpcClientPool) GrantToUser(ctx context.Context, param GrantToUserParams) error {
	client, err := pool.getRpcClient()
	if err != nil {
		return fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.GrantToUser(ctx, param)
}

func (pool *RpcClientPool) RevokeFromUser(ctx context.Context, param RevokeFromUserParams) error {
	client, err := pool.getRpcClient()
	if err != nil {
		return fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.RevokeFromUser(ctx, param)
}

func (pool *RpcClientPool) DescribeUser(ctx context.Context, param DescribeUserParams) (
	result *DescribeUserResult, err error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return nil, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.DescribeUser(ctx, param)
}

func (pool *RpcClientPool) ListUser(ctx context.Context) (result *ListUserResult, err error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return nil, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.ListUser(ctx)
}

func (pool *RpcClientPool) DropUser(ctx context.Context, param DropUserParams) error {
	client, err := pool.getRpcClient()
	if err != nil {
		return fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.DropUser(ctx, param)
}

func (pool *RpcClientPool) ChangePassword(ctx context.Context, param ChangePasswordParams) error {
	client, err := pool.getRpcClient()
	if err != nil {
		return fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.ChangePassword(ctx, param)
}

func (pool *RpcClientPool) UploadFile(ctx context.Context, databaseName, collectionName string,
	param UploadFileParams) (result *UploadFileResult, err error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return nil, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.UploadFile(ctx, databaseName, collectionName, param)
}

func (pool *RpcClientPool) GetImageUrl(ctx context.Context, databaseName, collectionName string,
	param GetImageUrlParams) (result *GetImageUrlResult, err error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return nil, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.GetImageUrl(ctx, databaseName, collectionName, param)
}

func (pool *RpcClientPool) RebuildIndex(ctx context.Context, databaseName, collectionName string,
	params ...*RebuildIndexParams) (result *RebuildIndexResult, err error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return nil, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.RebuildIndex(ctx, databaseName, collectionName, params...)
}

func (pool *RpcClientPool) AddIndex(ctx context.Context, databaseName, collectionName string,
	params ...*AddIndexParams) (err error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.AddIndex(ctx, databaseName, collectionName, params...)
}

func (pool *RpcClientPool) DropIndex(ctx context.Context, databaseName, collectionName string,
	params DropIndexParams) (err error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.DropIndex(ctx, databaseName, collectionName, params)
}

func (pool *RpcClientPool) ModifyVectorIndex(ctx context.Context, databaseName, collectionName string,
	param ModifyVectorIndexParam) (err error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.ModifyVectorIndex(ctx, databaseName, collectionName, param)
}

func (pool *RpcClientPool) Embedding(ctx context.Context, param EmbeddingParams) (result *EmbeddingResult, err error) {
	client, err := pool.getRpcClient()
	if err != nil {
		return nil, fmt.Errorf("get rpc client failed. err: %v", err.Error())
	}
	return client.Embedding(ctx, param)
}
