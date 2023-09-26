package engine

import (
	"context"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entry"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/internal/engine/api/index"
)

var _ entry.IndexInterface = &implementerIndex{}

type implementerIndex struct {
	entry.SdkClient
	databaseName string
}

func (i *implementerIndex) IndexRebuild(ctx context.Context, collectionName string, option *entry.IndexRebuildOption) (*entry.IndexReBuildResult, error) {
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
	result := new(entry.IndexReBuildResult)
	result.TaskIds = res.TaskIds
	return result, nil
}
