package engine

import (
	"context"
	"strings"

	"git.woa.com/cloud_nosql/vectordb/vectordb-sdk-go/internal/engine/api/database"
	"git.woa.com/cloud_nosql/vectordb/vectordb-sdk-go/model"
)

type implementerDatabase struct {
	model.SdkClient
}

func VectorDB(sdkClient model.SdkClient) model.VectorDBClient {
	databaseImpl := new(implementerDatabase)
	databaseImpl.SdkClient = sdkClient
	return databaseImpl
}

func (i *implementerDatabase) CreateDatabase(ctx context.Context, name string) (*model.Database, error) {
	req := database.CreateReq{
		Database: name,
	}
	res := new(database.CreateRes)
	err := i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}
	return i.Database(name), err
}

func (i *implementerDatabase) DropDatabase(ctx context.Context, name string) (err error) {
	req := database.DropReq{Database: name}
	res := new(database.DropRes)
	err = i.Request(ctx, req, res)
	if err != nil && strings.Contains(err.Error(), "not exist") {
		return nil
	}
	return
}

func (i *implementerDatabase) ListDatabase(ctx context.Context) (databases []*model.Database, err error) {
	req := database.ListReq{}
	res := new(database.ListRes)
	err = i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}

	for _, v := range res.Databases {
		databases = append(databases, i.Database(v))
	}
	return
}

func (i *implementerDatabase) Database(name string) *model.Database {
	database := new(model.Database)
	collImpl := new(implementerCollection)
	collImpl.SdkClient = i.SdkClient
	collImpl.databaseName = name
	database.CollectionInterface = collImpl
	database.DatabaseName = name
	return database
}
