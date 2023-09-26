package entry

import (
	"context"
	"time"
)

// DatabaseInterface database api
type DatabaseInterface interface {
	SdkClient
	CreateDatabase(ctx context.Context, name string, option *CreateDatabaseOption) (*Database, error)
	DropDatabase(ctx context.Context, name string, option *DropDatabaseOption) (*DatabaseResult, error)
	ListDatabase(ctx context.Context, option *ListDatabaseOption) (databases []*Database, err error)
	Database(name string) *Database
}

// Database wrap the database parameters and collection interface to operating the collection api
type Database struct {
	CollectionInterface
	AliasInterface
	IndexInterface
	DatabaseName string
}

func (d *Database) Debug(v bool) {
	d.CollectionInterface.Debug(v)
}

func (d *Database) WithTimeout(t time.Duration) {
	d.CollectionInterface.WithTimeout(t)
}

type DatabaseResult struct {
	AffectedCount int
}

type CreateDatabaseOption struct{}

type DropDatabaseOption struct{}

type ListDatabaseOption struct{}
