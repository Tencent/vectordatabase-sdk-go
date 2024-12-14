package tcvectordb

import (
	"context"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/olama"
)

var _ IndexInterface = &rpcImplementerIndex{}

type rpcImplementerIndex struct {
	SdkClient
	rpcClient  olama.SearchEngineClient
	flat       FlatIndexInterface
	database   *Database
	collection *Collection
}

func (r *rpcImplementerIndex) RebuildIndex(ctx context.Context, params ...*RebuildIndexParams) (*RebuildIndexResult, error) {
	return r.flat.RebuildIndex(ctx, r.database.DatabaseName, r.collection.CollectionName, params...)
}

func (r *rpcImplementerIndex) AddIndex(ctx context.Context, params ...*AddIndexParams) error {
	return r.flat.AddIndex(ctx, r.database.DatabaseName, r.collection.CollectionName, params...)
}

// [ModifyVectorIndex] modifies vector indexes to an existing collection.
func (r *rpcImplementerIndex) ModifyVectorIndex(ctx context.Context, param ModifyVectorIndexParam) error {
	return r.flat.ModifyVectorIndex(ctx, r.database.DatabaseName, r.collection.CollectionName, param)
}
