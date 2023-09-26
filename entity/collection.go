package entity

import (
	"time"
)

// Collection wrap the collection parameters and document interface to operating the document api
type Collection struct {
	DocumentInterface
	DatabaseName   string
	CollectionName string
	DocumentCount  int64
	Alias          []string
	ShardNum       uint32
	ReplicasNum    uint32
	Indexes        Indexes
	IndexStatus    IndexStatus
	Embedding      Embedding
	Description    string
	Size           uint64
	CreateTime     time.Time
}

type CollectionResult struct {
	AffectedCount int
}

type CreateCollectionOption struct {
	Embedding *Embedding
}

type DescribeCollectionOption struct{}

type DropCollectionOption struct{}

type TruncateCollectionOption struct{}

type ListCollectionOption struct{}
