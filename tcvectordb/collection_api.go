package tcvectordb

import "time"

// Collection wrap the collection parameters and document interface to operating the document api
type Collection struct {
	DocumentInterface `json:"-"`
	IndexInterface    `json:"-"`
	DatabaseName      string
	CollectionName    string
	DocumentCount     int64
	Alias             []string
	ShardNum          uint32
	ReplicasNum       uint32
	Indexes           Indexes
	IndexStatus       IndexStatus
	Embedding         Embedding
	Description       string
	Size              uint64
	CreateTime        time.Time
}

func (c *Collection) Debug(v bool) {
	c.DocumentInterface.Debug(v)
}

func (c *Collection) WithTimeout(t time.Duration) {
	c.DocumentInterface.WithTimeout(t)
}

type CreateCollectionOption struct {
	Embedding *Embedding
}

type CreateCollectionResult struct{}

type DescribeCollectionOption struct {
}

type DescribeCollectionResult struct {
	Collection
}

type ListCollectionOption struct{}

type ListCollectionResult struct {
	Collections []*Collection
}

type DropCollectionOption struct{}

type DropCollectionResult struct {
	AffectedCount int
}

type TruncateCollectionOption struct{}

type TruncateCollectionResult struct {
	AffectedCount int
}
