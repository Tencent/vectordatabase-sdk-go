package tcvectordb

import (
	"context"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/olama"
)

var _ FlatIndexInterface = &rpcImplementerFlatIndex{}

type rpcImplementerFlatIndex struct {
	SdkClient
	rpcClient olama.SearchEngineClient
}

func (r *rpcImplementerFlatIndex) RebuildIndex(ctx context.Context, databaseName, collectionName string, params ...*RebuildIndexParams) (*RebuildIndexResult, error) {
	req := &olama.RebuildIndexRequest{
		Database:   databaseName,
		Collection: collectionName,
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

func (r *rpcImplementerFlatIndex) AddIndex(ctx context.Context, databaseName, collectionName string, params ...*AddIndexParams) error {
	req := &olama.AddIndexRequest{
		Database:   databaseName,
		Collection: collectionName,
	}
	defaultBuildExistedData := true
	if len(params) != 0 && params[0] != nil {
		param := params[0]
		req.Indexes = make(map[string]*olama.IndexColumn, len(param.FilterIndexs))
		for _, index := range param.FilterIndexs {
			req.Indexes[index.FieldName] = &olama.IndexColumn{
				FieldName: index.FieldName,
				FieldType: string(index.FieldType),
				IndexType: string(index.IndexType),
			}
		}
		if param.BuildExistedData == nil {
			req.BuildExistedData = defaultBuildExistedData
		} else {
			req.BuildExistedData = *param.BuildExistedData
		}
	} else {
		req.BuildExistedData = defaultBuildExistedData
	}

	_, err := r.rpcClient.AddIndex(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
