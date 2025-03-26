package tcvectordb

import (
	"context"
	"log"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/index"
)

var _ FlatIndexInterface = &implementerFlatIndex{}

type FlatIndexInterface interface {
	SdkClient

	// [RebuildIndex] rebuilds all indexes under the specified collection.
	RebuildIndex(ctx context.Context, databaseName, collectionName string, params ...*RebuildIndexParams) (result *RebuildIndexResult, err error)

	// [AddIndex] adds scalar field index to an existing collection.
	AddIndex(ctx context.Context, databaseName, collectionName string, params ...*AddIndexParams) (err error)

	// [DropIndex] drops scalar field index to an existing collection.
	DropIndex(ctx context.Context, databaseName, collectionName string, params DropIndexParams) (err error)

	// [ModifyVectorIndex] modifies vector indexes to an existing collection.
	ModifyVectorIndex(ctx context.Context, databaseName, collectionName string, param ModifyVectorIndexParam) (err error)
}

type implementerFlatIndex struct {
	SdkClient
}

// [RebuildIndexParams] holds the parameters for rebuilding indexes in a collection.
//
// Fields:
//   - DropBeforeRebuild: Whether to delete the old index before rebuilding the new index (defaults to false).
//     true: first delete the old index and then rebuild the index.
//     false: after creating the new index, then delete the old index.
//   - Throttle: (Optional) The number of CPU cores for building an index on a single node. only support 1 currently.
//   - UnLimitedCPU: (Optional) Using all CPU cores for building an index on a single node (defaults to false).
type RebuildIndexParams struct {
	DropBeforeRebuild bool
	Throttle          int
	UnLimitedCPU      bool
}

// [AddIndexParams] holds the parameters for adding scalar field index in a collection.
//
// Fields:
//   - FilterIndexs: Whether to delete the old index before rebuilding the new index.
//     true: first delete the old index and then rebuild the index.
//     false: after creating the new index, then delete the old index.
//   - BuildExistedData: (Optional) Whether scan historical data and build index (defaults to true).
//     If there is no need to scan historical data, you can set this to false.
type AddIndexParams struct {
	FilterIndexs     []FilterIndex
	BuildExistedData *bool
}

type ModifyVectorIndexParam struct {
	VectorIndexes []ModifyVectorIndex
	RebuildRules  *index.RebuildRules
}

type ModifyVectorIndex struct {
	FieldName  string
	MetricType MetricType
	Params     IndexParams
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
func (i *implementerFlatIndex) RebuildIndex(ctx context.Context, databaseName, collectionName string, params ...*RebuildIndexParams) (*RebuildIndexResult, error) {
	req := new(index.RebuildReq)
	req.Database = databaseName
	req.Collection = collectionName

	if len(params) != 0 && params[0] != nil {
		param := params[0]
		req.DropBeforeRebuild = param.DropBeforeRebuild
		req.Throttle = int32(param.Throttle)
		if req.Throttle == 0 {
			// UnLimitedCPU is true and  Throttle is 0, mean unlimited cpu to build
			if !param.UnLimitedCPU {
				req.Throttle = 1
			}
		}

	}

	res := new(index.RebuildRes)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return nil, err
	}
	result := new(RebuildIndexResult)
	result.TaskIds = res.TaskIds
	return result, nil
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
func (i *implementerFlatIndex) AddIndex(ctx context.Context, databaseName, collectionName string, params ...*AddIndexParams) error {
	req := new(index.AddReq)
	req.Database = databaseName
	req.Collection = collectionName
	if len(params) != 0 && params[0] != nil {
		param := params[0]
		for _, index := range param.FilterIndexs {
			req.Indexes = append(req.Indexes, &api.IndexColumn{
				FieldName: index.FieldName,
				FieldType: string(index.FieldType),
				IndexType: string(index.IndexType),
			})
		}
		req.BuildExistedData = param.BuildExistedData
	}

	res := new(index.AddRes)
	err := i.Request(ctx, req, res)
	if err != nil {
		return err
	}
	return nil
}

type DropIndexParams struct {
	FieldNames []string
}

func (i *implementerFlatIndex) DropIndex(ctx context.Context, databaseName, collectionName string, params DropIndexParams) error {
	req := new(index.DropReq)
	req.Database = databaseName
	req.Collection = collectionName
	req.FieldNames = params.FieldNames
	res := new(index.DropRes)
	err := i.Request(ctx, req, res)
	return err
}

// [ModifyVectorIndex] modifies vector indexes to an existing collection.
func (i *implementerFlatIndex) ModifyVectorIndex(ctx context.Context, databaseName, collectionName string, param ModifyVectorIndexParam) error {
	req := new(index.ModifyVectorIndexReq)
	req.Database = databaseName
	req.Collection = collectionName

	for _, v := range param.VectorIndexes {
		var column api.IndexColumn
		column.FieldName = v.FieldName
		column.MetricType = string(v.MetricType)

		optionParamsFromIndexParams(&column, v.Params)

		req.VectorIndexes = append(req.VectorIndexes, &column)
	}
	req.RebuildRules = param.RebuildRules

	res := new(index.ModifyVectorIndexRes)
	err := i.Request(ctx, req, res)
	if err != nil {
		return err
	}
	log.Println("[Warning]", res.Msg)
	return nil
}
