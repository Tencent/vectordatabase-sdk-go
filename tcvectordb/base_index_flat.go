package tcvectordb

import (
	"context"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/index"
)

var _ FlatIndexInterface = &implementerFlatIndex{}

type FlatIndexInterface interface {
	SdkClient
	RebuildIndex(ctx context.Context, databaseName, collectionName string, params ...*RebuildIndexParams) (result *RebuildIndexResult, err error)
	AddIndex(ctx context.Context, databaseName, collectionName string, params ...*AddIndexParams) (err error)
}

type implementerFlatIndex struct {
	SdkClient
}

type RebuildIndexParams struct {
	DropBeforeRebuild bool
	Throttle          int
}

type AddIndexParams struct {
	FilterIndexs     []FilterIndex
	BuildExistedData *bool
}

func (i *implementerFlatIndex) RebuildIndex(ctx context.Context, databaseName, collectionName string, params ...*RebuildIndexParams) (*RebuildIndexResult, error) {
	req := new(index.RebuildReq)
	req.Database = databaseName
	req.Collection = collectionName

	if len(params) != 0 && params[0] != nil {
		param := params[0]
		req.DropBeforeRebuild = param.DropBeforeRebuild
		req.Throttle = int32(param.Throttle)
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
