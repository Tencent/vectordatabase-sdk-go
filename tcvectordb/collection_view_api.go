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
	AIDocumentSetInterface `json:"-"`
	DatabaseName           string                              `json:"databaseName"`
	CollectionViewName     string                              `json:"collectionViewName"`
	Alias                  []string                            `json:"alias"`
	Embedding              *collection_view.DocumentEmbedding  `json:"embedding"`
	SplitterPreprocess     *collection_view.SplitterPreprocess `json:"splitterPreprocess"`
	IndexedDocumentSets    uint64                              `json:"indexedDocumentSets"`
	TotalDocumentSets      uint64                              `json:"totalDocumentSets"`
	UnIndexedDocumentSets  uint64                              `json:"unIndexedDocumentSets"`
	FilterIndexes          []FilterIndex                       `json:"filterIndexes"`
	Description            string                              `json:"description"`
	CreateTime             time.Time                           `json:"createTime"`
}

type CreateCollectionViewParams struct {
	Description        string
	Indexes            Indexes                             `json:"indexes"`
	Embedding          *collection_view.DocumentEmbedding  `json:"embedding"`
	SplitterPreprocess *collection_view.SplitterPreprocess `json:"splitterPreprocess"`
}

type CreateAICollectionResult struct {
	CollectionView `json:"collectionView"`
	AffectedCount  int
}

type DescribeCollectionViewResult struct {
	CollectionView `json:"collectionView"`
}

type DropCollectionViewResult struct {
	AffectedCount int
}

type TruncateCollectionViewResult struct {
	AffectedCount int
}

type ListCollectionViewsResult struct {
	CollectionViews []*CollectionView `json:"collectionViews"`
}
