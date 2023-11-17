package entity

import (
	"time"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity/api/ai_collection"
)

// Collection wrap the collection parameters and document interface to operating the document api
type AICollection struct {
	AIDocumentInterface
	DatabaseName       string
	CollectionName     string
	Alias              []string
	AiConfig           AiConfig
	IndexedDocuments   uint64
	TotalDocuments     uint64
	UnIndexedDocuments uint64
	FilterIndexes      []FilterIndex
	Description        string
	CreateTime         time.Time
}

type CreateAICollectionOption struct {
	Description string
	Indexes     Indexes
	AiConfig    *AiConfig
}

type CreateAICollectionResult struct {
	AICollection
	AffectedCount int
}

type AiConfig struct {
	ExpectedFileNum    uint64                            `json:"expectedFileNum,omitempty"`
	AverageFileSize    uint64                            `json:"averageFileSize,omitempty"`
	Language           Language                          `json:"language,omitempty"`
	DocumentPreprocess *ai_collection.DocumentPreprocess `json:"documentPreprocess,omitempty"`
	// DocumentIndex      *ai_collection.DocumentIndex      `json:"documentIndex,omitempty"`
}

type DescribeAICollectionOption struct{}

type DescribeAICollectionResult struct {
	AICollection
}

type DropAICollectionOption struct{}

type DropAICollectionResult struct {
	AffectedCount int
}

type TruncateAICollectionOption struct{}

type TruncateAICollectionResult struct {
	AffectedCount int
}

type ListAICollectionOption struct{}

type ListAICollectionResult struct {
	Collections []*AICollection
}
