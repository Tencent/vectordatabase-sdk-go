package engine

import (
	"context"
	"strings"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity/api/ai_database"
)

var _ entity.AIDatabaseInterface = &implementerAIDatabase{}

type implementerAIDatabase struct {
	entity.SdkClient
}

func (i *implementerAIDatabase) CreateAIDatabase(ctx context.Context, name string, option *entity.CreateAIDatabaseOption) (result *entity.CreateAIDatabaseResult, err error) {
	req := ai_database.CreateReq{
		Database: name,
	}
	res := new(ai_database.CreateRes)
	err = i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}
	result = new(entity.CreateAIDatabaseResult)
	result.AffectedCount = res.AffectedCount
	result.AIDatabase = *(i.AIDatabase(name))
	return result, err
}

// DropDatabase drop database with database name. If database not exist, it return nil.
func (i *implementerAIDatabase) DropAIDatabase(ctx context.Context, name string, option *entity.DropAIDatabaseOption) (result *entity.DropAIDatabaseResult, err error) {
	result = new(entity.DropAIDatabaseResult)

	req := ai_database.DropReq{Database: name}
	res := new(ai_database.DropRes)
	err = i.Request(ctx, req, res)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			return result, nil
		}
		return
	}
	result.AffectedCount = res.AffectedCount
	return
}

// ListAIDatabase get ai_database list. It returns the ai_database list to operate the collection.
func (i *implementerAIDatabase) ListAIDatabase(ctx context.Context, option *entity.ListAIDatabaseOption) (result *entity.ListAIDatabaseResult, err error) {
	req := ai_database.ListReq{}
	res := new(ai_database.ListRes)
	err = i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}
	result = new(entity.ListAIDatabaseResult)
	for _, v := range res.Databases {
		if info := res.Info[v]; info != nil && info.DbType != entity.AIDOCDbType {
			continue
		}
		result.Databases = append(result.Databases, *(i.AIDatabase(v)))
	}
	return
}

// Database get a ai_database interface to operate collection.  It could not send http request to vectordb.
func (i *implementerAIDatabase) AIDatabase(name string) *entity.AIDatabase {
	database := new(entity.AIDatabase)
	collImpl := new(implementerAICollection)
	collImpl.SdkClient = i.SdkClient
	collImpl.databaseName = name
	database.AICollectionInterface = collImpl

	aliasImpl := new(implementerAIAlias)
	aliasImpl.databaseName = name
	aliasImpl.SdkClient = i.SdkClient
	database.AIAliasInterface = aliasImpl

	database.DatabaseName = name

	database.Info = entity.DatabaseItem{
		DbType: entity.AIDOCDbType,
	}

	return database
}
