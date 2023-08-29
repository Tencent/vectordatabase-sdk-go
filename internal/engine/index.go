package engine

import (
	"context"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/internal/engine/api/index"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/model"
)

type implementerIndex struct {
	model.SdkClient
	databaseName string
}

func (i *implementerIndex) IndexRebuild(ctx context.Context, collectionName string, dropBeforeRebuild bool, throttle int) error {
	req := new(index.RebuildReq)
	req.Database = i.databaseName
	req.Collection = collectionName
	req.DropBeforeRebuild = dropBeforeRebuild
	req.Throttle = int32(throttle)

	res := new(index.RebuildRes)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return err
	}
	return nil
}
