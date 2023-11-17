// Copyright (C) 2023 Tencent Cloud.
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the vectordb-sdk-java), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is furnished
// to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
// PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
// SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package tcvectordb

import (
	"context"
	"strings"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb/api/ai_database"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb/api/database"
)

var _ DatabaseInterface = &implementerDatabase{}

// DatabaseInterface database api
type DatabaseInterface interface {
	SdkClient
	CreateDatabase(ctx context.Context, name string, options ...*CreateDatabaseOption) (*CreateDatabaseResult, error)
	DropDatabase(ctx context.Context, name string, options ...*DropDatabaseOption) (*DropDatabaseResult, error)
	ListDatabase(ctx context.Context, options ...*ListDatabaseOption) (result *ListDatabaseResult, err error)
	CreateAIDatabase(ctx context.Context, name string, options ...*CreateAIDatabaseOption) (result *CreateAIDatabaseResult, err error)
	DropAIDatabase(ctx context.Context, name string, options ...*DropAIDatabaseOption) (result *DropAIDatabaseResult, err error)
	Database(name string) *Database
	AIDatabase(name string) *AIDatabase
}

type implementerDatabase struct {
	SdkClient
}

// CreateDatabase create database with database name. It returns error if name exist.
func (i *implementerDatabase) CreateDatabase(ctx context.Context, name string, option ...*CreateDatabaseOption) (result *CreateDatabaseResult, err error) {
	req := database.CreateReq{
		Database: name,
	}
	res := new(database.CreateRes)
	err = i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}
	result = new(CreateDatabaseResult)
	result.AffectedCount = res.AffectedCount
	result.Database = *(i.Database(name))
	return result, err
}

// CreateAIDatabase create ai database with database name. It returns error if name exist.
func (i *implementerDatabase) CreateAIDatabase(ctx context.Context, name string, option ...*CreateAIDatabaseOption) (result *CreateAIDatabaseResult, err error) {
	req := ai_database.CreateReq{
		Database: name,
	}
	res := new(ai_database.CreateRes)
	err = i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}
	result = new(CreateAIDatabaseResult)
	result.AffectedCount = res.AffectedCount
	result.AIDatabase = *(i.AIDatabase(name))
	return result, err
}

// DropDatabase drop database with database name. If database not exist, it return nil.
func (i *implementerDatabase) DropDatabase(ctx context.Context, name string, option ...*DropDatabaseOption) (result *DropDatabaseResult, err error) {
	result = new(DropDatabaseResult)

	req := database.DropReq{Database: name}
	res := new(database.DropRes)
	err = i.Request(ctx, req, res)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") || strings.Contains(err.Error(), "can not find database") {
			return result, nil
		}
		return
	}
	result.AffectedCount = int(res.AffectedCount)
	return
}

// DropAIDatabase drop ai database with database name. If database not exist, it return nil.
func (i *implementerDatabase) DropAIDatabase(ctx context.Context, name string, option ...*DropAIDatabaseOption) (result *DropAIDatabaseResult, err error) {
	result = new(DropAIDatabaseResult)

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

// ListDatabase get database list. It returns the database list to operate the collection.
func (i *implementerDatabase) ListDatabase(ctx context.Context, option ...*ListDatabaseOption) (result *ListDatabaseResult, err error) {
	req := database.ListReq{}
	res := new(database.ListRes)
	err = i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}

	result = new(ListDatabaseResult)
	for _, v := range res.Databases {
		if res.Info[v].DbType == AIDOCDbType {
			db := i.AIDatabase(v)
			db.Info.CreateTime = res.Info[v].CreateTime
			result.AIDatabases = append(result.AIDatabases, *db)
		} else {
			db := i.Database(v)
			db.Info.CreateTime = res.Info[v].CreateTime
			db.Info.DbType = res.Info[v].DbType
			result.Databases = append(result.Databases, *db)
		}
	}
	return result, nil
}

// Database get a database interface to operate collection.  It could not send http request to vectordb.
func (i *implementerDatabase) Database(name string) *Database {
	database := new(Database)
	database.DatabaseName = name
	database.Info = DatabaseItem{
		DbType: BASEDbType,
	}

	collImpl := new(implementerCollection)
	collImpl.SdkClient = i.SdkClient
	collImpl.database = *database

	aliasImpl := new(implementerAlias)
	aliasImpl.database = *database
	aliasImpl.SdkClient = i.SdkClient

	indexImpl := new(implementerIndex)
	indexImpl.database = *database
	indexImpl.SdkClient = i.SdkClient

	database.CollectionInterface = collImpl
	database.AliasInterface = aliasImpl
	database.IndexInterface = indexImpl
	return database
}

// Database get a ai_database interface to operate collection.  It could not send http request to vectordb.
func (i *implementerDatabase) AIDatabase(name string) *AIDatabase {
	database := new(AIDatabase)
	database.DatabaseName = name
	database.Info = DatabaseItem{
		DbType: AIDOCDbType,
	}

	collImpl := new(implementerAICollection)
	collImpl.SdkClient = i.SdkClient
	collImpl.database = *database

	aliasImpl := new(implementerAIAlias)
	aliasImpl.database = *database
	aliasImpl.SdkClient = i.SdkClient

	database.AICollectionInterface = collImpl
	database.AIAliasInterface = aliasImpl

	return database
}
