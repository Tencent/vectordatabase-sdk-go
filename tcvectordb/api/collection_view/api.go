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

package collection_view

import "github.com/tencent/vectordatabase-sdk-go/tcvectordb/api"

// CreateReq create CollectionView request
type CreateReq struct {
	api.Meta       `path:"/ai/collectionView/create" tags:"ai" method:"Post" summary:"创建collection存储embedding文件集合"`
	Database       string `json:"database"`
	CollectionView string `json:"collectionView"`
	Description    string `json:"description,omitempty"`
	// ExpectedFileNum    uint64              `json:"expectedFileNum,omitempty"`
	// AverageFileSize    uint64              `json:"averageFileSize,omitempty"`
	Embedding          *DocumentEmbedding  `json:"embedding,omitempty"`
	SplitterPreprocess *SplitterPreprocess `json:"splitterPreprocess,omitempty"`
	ParsingProcess     *api.ParsingProcess `json:"parsingProcess,omitempty"`
	Indexes            []*api.IndexColumn  `json:"indexes,omitempty"`
	ExpectedFileNum    uint64              `json:"expectedFileNum,omitempty"`
	AverageFileSize    uint64              `json:"averageFileSize,omitempty"`
}

type DocumentEmbedding struct {
	Language             string `json:"language,omitempty"`
	EnableWordsEmbedding *bool  `json:"enableWordsEmbedding,omitempty"`
}

type CreateRes struct {
	api.CommonRes
	AffectedCount uint64 `json:"affectedCount"`
}

type SplitterPreprocess struct {
	// AppendTitleToChunk
	// false(default): no process
	// true: append the file paragraph title to chunk for embedding
	AppendTitleToChunk *bool `json:"appendTitleToChunk,omitempty"`
	// AppendKeywordsToChunk
	// false: no process
	// true(default): append the file keywords to chunk for embedding
	AppendKeywordsToChunk *bool `json:"appendKeywordsToChunk,omitempty"`
}

// DescribeReq get collectionView detail request
type DescribeReq struct {
	api.Meta       `path:"/ai/collectionView/describe" tags:"Collection" method:"Post" summary:"返回collection信息"`
	Database       string `json:"database"`
	CollectionView string `json:"collectionView"`
}

// DescribeRes get collectionView detail response
type DescribeRes struct {
	api.CommonRes
	CollectionView *DescribeCollectionViewItem `json:"collectionView"`
}

type DescribeCollectionViewItem struct {
	Database       string `json:"database"`
	CollectionView string `json:"collectionView"`
	Description    string `json:"description,omitempty"`
	// ExpectedFileNum    uint64              `json:"expectedFileNum,omitempty"`
	// AverageFileSize    uint64              `json:"averageFileSize,omitempty"`
	Embedding          *DocumentEmbedding  `json:"embedding,omitempty"`
	SplitterPreprocess *SplitterPreprocess `json:"splitterPreprocess,omitempty"`
	ParsingProcess     *api.ParsingProcess `json:"parsingProcess,omitempty"`
	Indexes            []*api.IndexColumn  `json:"indexes,omitempty"`

	CreateTime string   `json:"createTime"`
	Alias      []string `json:"alias,omitempty"`
	Status     *Status  `json:"stats"`
}

type Status struct {
	FailIndexedDocumentSets uint64 `json:"failIndexedDocumentSets"`
	IndexedDocumentSets     uint64 `json:"indexedDocumentSets"`
	TotalDocumentSets       uint64 `json:"totalDocumentSets"`
	UnIndexedDocumentSets   uint64 `json:"unIndexedDocumentSets"`
}

// DropReq delete collectionView request
type DropReq struct {
	api.Meta       `path:"/ai/collectionView/drop" tags:"Collection" method:"Post" summary:"删除collection，并删除collection中的所有文档，如果collectio不经存在返回失败"`
	Database       string `json:"database"`
	CollectionView string `json:"collectionView"`
}

// DropReq delete collectionView response
type DropRes struct {
	api.CommonRes
	AffectedCount uint64 `json:"affectedCount"`
}

type ListReq struct {
	api.Meta `path:"/ai/collectionView/list" tags:"Collection" method:"Post" summary:"列出指定database中的所有collectionView"`
	Database string `json:"database"`
}

type ListRes struct {
	api.CommonRes
	CollectionViews []*DescribeCollectionViewItem `json:"collectionViews"`
}

type TruncateReq struct {
	api.Meta       `path:"/ai/collectionView/truncate" tags:"Collection" method:"Post" summary:"清空 collection 中的所有数据和索引"`
	Database       string `json:"database"`
	CollectionView string `json:"collectionView"`
}

type TruncateRes struct {
	api.CommonRes
	AffectedCount uint64 `json:"affectedCount"`
}
