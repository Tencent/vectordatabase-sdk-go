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
	"fmt"
	"strings"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/ai_database"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/database"
)

var _ DatabaseInterface = &implementerDatabase{}

// DatabaseInterface provides apis of a database.
type DatabaseInterface interface {
	SdkClient

	// [ExistsDatabase] checks the existence of a specific database.
	ExistsDatabase(ctx context.Context, name string) (bool, error)

	// [CreateDatabaseIfNotExists] creates a database if it doesn't exist.
	CreateDatabaseIfNotExists(ctx context.Context, name string) (*CreateDatabaseResult, error)

	// [CreateDatabase] creates a new database with user-defined name, and its database type is BASE_DB.
	CreateDatabase(ctx context.Context, name string) (*CreateDatabaseResult, error)

	// [DropDatabase] drops a specific database.
	DropDatabase(ctx context.Context, name string) (*DropDatabaseResult, error)

	// [ListDatabase] retrieves the list of all Database and the list of all AIDatabases in a vectordb.
	ListDatabase(ctx context.Context) (result *ListDatabaseResult, err error)

	// [CreateAIDatabase] creates a new AI database with user-defined name, and its database type is AI_DB.
	CreateAIDatabase(ctx context.Context, name string) (result *CreateAIDatabaseResult, err error)

	// [DropAIDatabase] drops the AI database.
	DropAIDatabase(ctx context.Context, name string) (result *DropAIDatabaseResult, err error)

	// [Database] returns a pointer to a [Database] object.
	Database(name string) *Database

	// [AIDatabase] returns a pointer to an [AIDatabase] object.
	AIDatabase(name string) *AIDatabase
}

type implementerDatabase struct {
	SdkClient
}

type CreateDatabaseResult struct {
	Database
	AffectedCount int
}

// [ExistsDatabase] checks the existence of a specific database.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - name: The name of the [database] to check.
//
// Notes: It returns true if the database exists.
//
// Returns a boolean variable indicating whether the database exists or an error.
func (i *implementerDatabase) ExistsDatabase(ctx context.Context, name string) (bool, error) {
	dbList, err := i.ListDatabase(ctx)
	if err != nil {
		return false, fmt.Errorf("judging whether the database exists failed. err is %v", err.Error())
	}
	for _, db := range dbList.Databases {
		if db.DatabaseName == name {
			return true, nil
		}
	}
	for _, db := range dbList.AIDatabases {
		if db.DatabaseName == name {
			return true, nil
		}
	}
	return false, nil
}

// [CreateDatabaseIfNotExists] creates a database if it doesn't exist.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - name: The name of the [database] to create.
//
// Returns a pointer to a [CreateDatabaseResult] object or an error.
func (i *implementerDatabase) CreateDatabaseIfNotExists(ctx context.Context, name string) (*CreateDatabaseResult, error) {
	dbList, err := i.ListDatabase(ctx)
	if err != nil {
		return nil, fmt.Errorf("judging whether the database exists failed. err is %v", err.Error())
	}
	for _, db := range dbList.Databases {
		if db.DatabaseName == name {
			result := new(CreateDatabaseResult)
			result.Database = *(i.Database(name))
			return result, err
		}
	}

	return i.CreateDatabase(ctx, name)
}

// [CreateDatabase] creates a new database with user-defined name, and its database type is BASE_DB.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - name: The name of the [database] to drop.
//
// Notes: It returns error if the database exist.
//
// Returns a pointer to a [CreateDatabaseResult] object or an error.
func (i *implementerDatabase) CreateDatabase(ctx context.Context, name string) (result *CreateDatabaseResult, err error) {
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

type CreateAIDatabaseResult struct {
	AIDatabase
	AffectedCount int32
}

// [CreateAIDatabase] creates a new AI database with user-defined name, and its database type is AI_DB.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - name: The name of the [database] to create.
//
// Notes: It returns error if the database exist.
//
// Returns a pointer to a [CreateAIDatabaseResult] object or an error.
func (i *implementerDatabase) CreateAIDatabase(ctx context.Context, name string) (result *CreateAIDatabaseResult, err error) {
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

type DropDatabaseResult struct {
	AffectedCount int
}

// [DropDatabase] drops a specific database.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - name: The name of the [database] to drop.
//
// Notes: If the database doesn't exist, it returns 0 for DropDatabaseResult.AffectedCount.
//
// Returns a pointer to a [DropDatabaseResult] object or an error.
func (i *implementerDatabase) DropDatabase(ctx context.Context, name string) (result *DropDatabaseResult, err error) {
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

type DropAIDatabaseResult struct {
	AffectedCount int32
}

// [DropAIDatabase] drops a specific AI database. Its database type is AI_DB
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - name: The name of the AI [database] to drop.
//
// Notes: If the database doesn't exist, it returns 0 for DropAIDatabaseResult.AffectedCount.
//
// Returns a pointer to a [DropAIDatabaseResult] object or an error.
func (i *implementerDatabase) DropAIDatabase(ctx context.Context, name string) (result *DropAIDatabaseResult, err error) {
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

type ListDatabaseResult struct {
	Databases   []Database
	AIDatabases []AIDatabase
}

// [ListDatabase] retrieves the list of all Database and the list of all AIDatabases in a vectordb.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//
// Returns a pointer to a [ListDatabaseResult] object or an error. See [ListDatabaseResult] for more information.
func (i *implementerDatabase) ListDatabase(ctx context.Context) (result *ListDatabaseResult, err error) {
	req := database.ListReq{}
	res := new(database.ListRes)
	err = i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}

	result = new(ListDatabaseResult)
	for _, v := range res.Databases {
		if res.Info[v].DbType == AIDOCDbType || res.Info[v].DbType == DbTypeAI {
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

// [Database] returns a pointer to a [Database] object.
//
// Parameters:
//   - name: The name of the [database], which database type is BASE_DB.

// Returns a pointer to a [Database] object, which includes some interfaces to operate collection, index, and alias.
func (i *implementerDatabase) Database(name string) *Database {
	database := new(Database)
	database.DatabaseName = name
	database.Info = DatabaseItem{
		DbType: DbTypeBase,
	}

	collImpl := new(implementerCollection)
	collImpl.SdkClient = i.SdkClient
	collImpl.database = database

	aliasImpl := new(implementerAlias)
	aliasImpl.database = database
	aliasImpl.SdkClient = i.SdkClient

	indexImpl := new(implementerIndex)
	indexImpl.database = database
	indexImpl.SdkClient = i.SdkClient

	database.CollectionInterface = collImpl
	database.AliasInterface = aliasImpl
	database.IndexInterface = indexImpl
	return database
}

// [AIDatabase] returns a pointer to a [AIDatabase] object.
//
// Parameters:
//   - name: The name of the [database], which database type is AI_DB.
//
// Returns a pointer to a [AIDatabase] object, which includes some interfaces to operate collectionView and alias.
func (i *implementerDatabase) AIDatabase(name string) *AIDatabase {
	database := new(AIDatabase)
	database.DatabaseName = name
	database.Info = DatabaseItem{
		DbType: DbTypeAI,
	}

	collImpl := new(implementerCollectionView)
	collImpl.SdkClient = i.SdkClient
	collImpl.database = database

	aliasImpl := new(implementerAIAlias)
	aliasImpl.database = database
	aliasImpl.SdkClient = i.SdkClient

	database.AICollectionViewInterface = collImpl
	database.AIAliasInterface = aliasImpl

	return database
}

// Database wraps the database parameters,
// collection interface to operate the collection api,
// alias interface to operate the alias api,
// and index interface to operate the index api.
type Database struct {
	CollectionInterface `json:"-"`
	AliasInterface      `json:"-"`
	IndexInterface      `json:"-"`
	DatabaseName        string
	Info                DatabaseItem
}

// IsAIDatabase  checks if the database is an AI database, according to the database type AI_DOC or AI_DB.
// Deprecated: the database type AI_DOC is no longer in use.
func (d Database) IsAIDatabase() bool {
	return d.Info.DbType == AIDOCDbType || d.Info.DbType == DbTypeAI
}

type DatabaseItem struct {
	CreateTime string `json:"createTime,omitempty"`
	DbType     string `json:"dbType,omitempty"`
}

// Debug sets the debug mode for the SdkClient, which prints the request and response of a network call.
func (d *Database) Debug(v bool) {
	d.CollectionInterface.Debug(v)
}

// WithTimeout sets the timeout duration for network operations for the SdkClient.
func (d *Database) WithTimeout(t time.Duration) {
	d.CollectionInterface.WithTimeout(t)
}

// AIDatabase wrap the database parameters and collection interface to operating the ai_collection api
type AIDatabase struct {
	AICollectionViewInterface
	AIAliasInterface
	DatabaseName string
	Info         DatabaseItem
}

func (d AIDatabase) IsAIDatabase() bool {
	return d.Info.DbType == AIDOCDbType || d.Info.DbType == DbTypeAI
}

func (d *AIDatabase) Debug(v bool) {
	d.AICollectionViewInterface.Debug(v)
}

func (d *AIDatabase) WithTimeout(t time.Duration) {
	d.AICollectionViewInterface.WithTimeout(t)
}
