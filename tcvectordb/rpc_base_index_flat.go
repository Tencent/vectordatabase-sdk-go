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

// [RebuildIndex] rebuilds all indexes under the specified collection.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - databaseName: The name of the database.
//   - collectionName: The name of the collection.
//   - params: A pointer to a [RebuildIndexParams] object that includes the other parameters for the rebuilding indexes operation.
//     See [RebuildIndexParams] for more information.
//
// Returns a pointer to a [RebuildIndexResult] object or an error.
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

// [AddIndex] adds scalar field index to an existing collection.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - databaseName: The name of the database.
//   - collectionName: The name of the collection.
//   - params: A pointer to a [AddIndexParams] object that includes the other parameters for the adding scalar field index operation.
//     See [AddIndexParams] for more information.
//
// Returns an error if the addition fails.
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
