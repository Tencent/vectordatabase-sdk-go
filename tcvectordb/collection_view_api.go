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

	collection_view "git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb/api/collection_view"
)

// CollectionView wrap the collectionView parameters and document interface to operating the document api
type CollectionView struct {
	AIDocumentSetInterface
	DatabaseName          string
	CollectionName        string
	Alias                 []string
	Embedding             *collection_view.DocumentEmbedding  `json:"embedding,omitempty"`
	SplitterPreprocess    *collection_view.SplitterPreprocess `json:"splitterPreprocess,omitempty"`
	IndexedDocumentSets   uint64
	TotalDocumentSets     uint64
	UnIndexedDocumentSets uint64
	FilterIndexes         []FilterIndex
	Description           string
	CreateTime            time.Time
}

type CreateCollectionViewOption struct {
	Description        string
	Indexes            Indexes
	Embedding          *collection_view.DocumentEmbedding  `json:"embedding,omitempty"`
	SplitterPreprocess *collection_view.SplitterPreprocess `json:"splitterPreprocess,omitempty"`
}

type CreateAICollectionResult struct {
	CollectionView
	AffectedCount int
}

type DescribeCollectionViewOption struct{}

type DescribeCollectionViewResult struct {
	CollectionView
}

type DropCollectionViewOption struct{}

type DropCollectionViewResult struct {
	AffectedCount int
}

type TruncateCollectionViewOption struct{}

type TruncateCollectionViewResult struct {
	AffectedCount int
}

type ListCollectionViewsOption struct{}

type ListCollectionViewsResult struct {
	CollectionViews []*CollectionView
}
