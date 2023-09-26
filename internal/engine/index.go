package engine

import (
	"context"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/internal/engine/api/index"
)

var _ entity.IndexInterface = &implementerIndex{}

type implementerIndex struct {
	entity.SdkClient
	databaseName string
}

func (i *implementerIndex) IndexRebuild(ctx context.Context, collectionName string, option *entity.IndexRebuildOption) (*entity.IndexReBuildResult, error) {
	req := new(index.RebuildReq)
	req.Database = i.databaseName
	req.Collection = collectionName

	req.DropBeforeRebuild = option.DropBeforeRebuild
	req.Throttle = int32(option.Throttle)

	res := new(index.RebuildRes)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return nil, err
	}
	result := new(entity.IndexReBuildResult)
	result.TaskIds = res.TaskIds
	return result, nil
}
