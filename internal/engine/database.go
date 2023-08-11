package engine

import (
	"context"

	"vectordb-sdk-go/internal/client"
	"vectordb-sdk-go/internal/engine/api/database"
)

type implementerDatabase struct {
	client.SdkClient
}

func (i *implementerDatabase) CreateDatabase(ctx context.Context, name string) (*Database, error) {
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
	return
}

func (i *implementerDatabase) ListDatabase(ctx context.Context) (databases []*Database, err error) {
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

func (i *implementerDatabase) Database(name string) *Database {
	database := new(Database)
	collImpl := new(implementerCollection)
	collImpl.SdkClient = i.SdkClient
	collImpl.databaseName = name
	database.CollectionInterface = collImpl
	database.DatabaseName = name
	return database
}
