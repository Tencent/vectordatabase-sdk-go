package tcvectordb

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/olama"
)

var _ DatabaseInterface = &rpcImplementerDatabase{}

type rpcImplementerDatabase struct {
	SdkClient
	httpImplementer DatabaseInterface
	rpcClient       olama.SearchEngineClient
}

func (r *rpcImplementerDatabase) ExistsDatabase(ctx context.Context, name string) (bool, error) {
	dbList, err := r.ListDatabase(ctx)
	if err != nil {
		return false, fmt.Errorf("judging whether the database exists failed. err is %v", err.Error())
	}
	for _, db := range dbList.Databases {
		if db.DatabaseName == name {
			return true, nil
		}
	}
	return false, nil
}

func (r *rpcImplementerDatabase) CreateDatabaseIfNotExists(ctx context.Context, name string) (*CreateDatabaseResult, error) {
	dbList, err := r.ListDatabase(ctx)
	if err != nil {
		return nil, fmt.Errorf("judging whether the database exists failed. err is %v", err.Error())
	}
	for _, db := range dbList.Databases {
		if db.DatabaseName == name {
			result := new(CreateDatabaseResult)
			result.Database = *(r.Database(name))
			return result, err
		}
	}
	return r.CreateDatabase(ctx, name)
}

func (r *rpcImplementerDatabase) CreateDatabase(ctx context.Context, name string) (*CreateDatabaseResult, error) {
	req := &olama.DatabaseRequest{
		Database: name,
		DbType:   olama.DataType_BASE,
	}
	res, err := r.rpcClient.CreateDatabase(ctx, req)
	if err != nil {
		return nil, err
	}
	result := new(CreateDatabaseResult)
	result.AffectedCount = int(res.AffectedCount)
	result.Database = *(r.Database(name))
	return result, err
}

func (r *rpcImplementerDatabase) DropDatabase(ctx context.Context, name string) (*DropDatabaseResult, error) {
	result := new(DropDatabaseResult)
	req := &olama.DatabaseRequest{
		Database: name,
		DbType:   olama.DataType_BASE,
	}
	res, err := r.rpcClient.DropDatabase(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") || strings.Contains(err.Error(), "can not find database") {
			return result, nil
		}
		return result, err
	}
	result.AffectedCount = int(res.AffectedCount)
	return result, err
}

func (r *rpcImplementerDatabase) ListDatabase(ctx context.Context) (result *ListDatabaseResult, err error) {
	req := &olama.DatabaseRequest{}
	res, err := r.rpcClient.ListDatabases(ctx, req)
	if err != nil {
		return nil, err
	}
	result = new(ListDatabaseResult)
	for _, v := range res.Databases {
		if res.Info[v].DbType == olama.DataType_AI_DOC {
			db := r.AIDatabase(v)
			db.Info.CreateTime = strconv.FormatInt(res.Info[v].CreateTime, 10)
			result.AIDatabases = append(result.AIDatabases, *db)
		} else {
			db := r.Database(v)
			db.Info.CreateTime = strconv.FormatInt(res.Info[v].CreateTime, 10)
			db.Info.DbType = ConvertDbType(res.Info[v].DbType)
			result.Databases = append(result.Databases, *db)
		}
	}
	return result, nil
}

func (r *rpcImplementerDatabase) CreateAIDatabase(ctx context.Context, name string) (result *CreateAIDatabaseResult, err error) {
	return r.httpImplementer.CreateAIDatabase(ctx, name)
}

func (r *rpcImplementerDatabase) DropAIDatabase(ctx context.Context, name string) (*DropAIDatabaseResult, error) {
	return r.httpImplementer.DropAIDatabase(ctx, name)
}

func (r *rpcImplementerDatabase) Database(name string) *Database {
	database := &Database{
		DatabaseName: name,
		Info: DatabaseItem{
			DbType: DbTypeBase,
		},
	}
	collImpl := &rpcImplementerCollection{
		SdkClient: r.SdkClient,
		rpcClient: r.rpcClient,
		database:  database,
	}
	aliasImpl := &rpcImplementerAlias{
		SdkClient: r.SdkClient,
		rpcClient: r.rpcClient,
		database:  database,
	}
	indexImpl := &rpcImplementerIndex{
		SdkClient: r.SdkClient,
		rpcClient: r.rpcClient,
		database:  database,
	}
	database.CollectionInterface = collImpl
	database.AliasInterface = aliasImpl
	database.IndexInterface = indexImpl
	return database
}

func (r *rpcImplementerDatabase) AIDatabase(name string) *AIDatabase {
	return r.httpImplementer.AIDatabase(name)
}
