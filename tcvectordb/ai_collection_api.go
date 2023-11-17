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
	"time"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb/api/ai_collection"
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
