package tcvectordb

import (
	"context"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/olama"
)

var _ IndexInterface = &rpcImplementerIndex{}

type rpcImplementerIndex struct {
	SdkClient
	rpcClient  olama.SearchEngineClient
	database   *Database
	collection *Collection
}

func (r *rpcImplementerIndex) RebuildIndex(ctx context.Context, params ...*RebuildIndexParams) (*RebuildIndexResult, error) {
	if r.database.IsAIDatabase() {
		return nil, AIDbTypeError
	}
	req := &olama.RebuildIndexRequest{
		Database:   r.database.DatabaseName,
		Collection: r.collection.CollectionName,
	}
	if len(params) != 0 && params[0] != nil {
		param := params[0]
		req.DropBeforeRebuild = param.DropBeforeRebuild
		req.Throttle = int32(param.Throttle)
	}
	res, err := r.rpcClient.RebuildIndex(ctx, req)
	if err != nil {
		return nil, err
	}
	return &RebuildIndexResult{TaskIds: res.TaskIds}, nil
}
