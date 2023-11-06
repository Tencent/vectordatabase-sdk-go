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

package entity

import (
	"time"
)

// Database wrap the database parameters and collection interface to operating the collection api
type Database struct {
	CollectionInterface
	AliasInterface
	IndexInterface
	DatabaseName string
	Info         DatabaseItem
}

func (d Database) IsAIDatabase() bool {
	return d.Info.DbType == AIDOCDbType
}

type DatabaseItem struct {
	CreateTime string `json:"createTime,omitempty"`
	DbType     string `json:"dbType,omitempty"`
}

func (d *Database) Debug(v bool) {
	d.CollectionInterface.Debug(v)
}

func (d *Database) WithTimeout(t time.Duration) {
	d.CollectionInterface.WithTimeout(t)
}

type CreateDatabaseOption struct{}

type CreateDatabaseResult struct {
	Database
	AffectedCount int
}

type DropDatabaseOption struct{}

type DropDatabaseResult struct {
	AffectedCount int
}

type ListDatabaseOption struct{}

type ListDatabaseResult struct {
	Databases   []Database
	AIDatabases []AIDatabase
}
