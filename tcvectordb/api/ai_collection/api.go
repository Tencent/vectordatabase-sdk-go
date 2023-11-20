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

package ai_collection

import "git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb/api"

// CreateReq create collection request
type CreateReq struct {
	api.Meta             `path:"/ai/collection/create" tags:"ai" method:"Post" summary:"创建collection存储embedding文件集合"`
	Database             string              `json:"database"`
	Collection           string              `json:"collection"`
	Description          string              `json:"description"`
	ExpectedFileNum      uint64              `json:"expectedFileNum"`
	AverageFileSize      uint64              `json:"averageFileSize"`
	Language             string              `json:"language"`
	DocumentPreprocess   *DocumentPreprocess `json:"documentPreprocess"`
	EnableWordsEmbedding *bool               `json:"enableWordsEmbedding,omitempty"`
	// DocumentIndex      *DocumentIndex      `json:"document_index"`
	Indexes []api.IndexColumn `json:"indexes"`
}

type CreateRes struct {
	api.CommonRes
	AffectedCount int `json:"affectedCount"`
}

type DocumentPreprocess struct {
	// AppendTitleToChunk
	// "0", "": no process
	// "1": append the file paragraph title to chunk for embedding
	AppendTitleToChunk string `json:"appendTitleToChunk"`
	// AppendKeywordsToChunk
	// "0", "": no process
	// "1": append the file keywords to chunk for embedding
	AppendKeywordsToChunk string `json:"appendKeywordsToChunk"`
}

type DocumentIndex struct {
	EnableWordsSimilarity *bool `json:"enableWordsSimilarity"`
}

// DescribeReq get collection detail request
type DescribeReq struct {
	api.Meta   `path:"/ai/collection/describe" tags:"Collection" method:"Post" summary:"返回collection信息"`
	Database   string `json:"database"`
	Collection string `json:"collection"`
}

// DescribeRes get collection detail response
type DescribeRes struct {
	api.CommonRes
	Collection *DescribeAICollectionItem `json:"collection"`
}

// DropReq delete collection request
type DropReq struct {
	api.Meta   `path:"/ai/collection/drop" tags:"Collection" method:"Post" summary:"删除collection，并删除collection中的所有文档，如果collectio不经存在返回失败"`
	Database   string `json:"database"`
	Collection string `json:"collection"`
}

// DropReq delete collection response
type DropRes struct {
	api.CommonRes
	AffectedCount int32 `json:"affectedCount"`
}

type ListReq struct {
	api.Meta `path:"/ai/collection/list" tags:"Collection" method:"Post" summary:"列出指定database中的所有collection"`
	Database string `json:"database"`
}

type ListRes struct {
	api.CommonRes
	Collections []DescribeAICollectionItem `json:"collections"`
}

type TruncateReq struct {
	api.Meta   `path:"/ai/collection/truncate" tags:"Collection" method:"Post" summary:"清空 collection 中的所有数据和索引"`
	Database   string `json:"database"`
	Collection string `json:"collection"`
}

type TruncateRes struct {
	api.CommonRes
	AffectedCount int32 `json:"affectedCount"`
}

type DescribeAICollectionItem struct {
	Database           string              `json:"database"`
	Collection         string              `json:"collection"`
	Language           string              `json:"language"`
	ExpectedFileNum    uint64              `json:"expectedFileNum"`
	AverageFileSize    uint64              `json:"averageFileSize"`
	CreateTime         string              `json:"createTime"`
	Description        string              `json:"description"`
	FilterIndexes      []api.IndexColumn   `json:"indexes"`
	Alias              []string            `json:"alias"`
	AiStatus           *AiStatus           `json:"aiStatus"`
	DocumentPreprocess *DocumentPreprocess `json:"documentPreprocess"`
	// DocumentIndex      DocumentIndex      `json:"document_index"`
}

type AiStatus struct {
	IndexedDocuments   uint64 `json:"indexedDocuments"`
	TotalDocuments     uint64 `json:"totalDocuments"`
	UnIndexedDocuments uint64 `json:"unIndexedDocuments"`
}
