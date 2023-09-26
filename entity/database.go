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
