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
	"io"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb/api/ai_document_set"
)

type AIDocumentSet struct {
	ai_document_set.QueryDocumentSet
}

type QueryAIDocumentSetParams struct {
	Filter *Filter
	Limit  int64
	Offset int64
}

type QueryAIDocumentSetResult struct {
	Count     uint64
	Documents []AIDocumentSet
}

type GetAIDocumentSetParams struct {
	DocumentSetId   string
	DocumentSetName string
}

type GetAIDocumentSetResult struct {
	Count        uint64
	DocumentSets AIDocumentSet `json:"documentSet"`
}

type SearchAIDocumentSetParams struct {
	Content         string
	DocumentSetName []string
	ExpandChunk     []int                         // 搜索结果中，向前、向后补齐几个chunk的上下文
	RerankOption    *ai_document_set.RerankOption // 多路召回
	// MergeChunk  bool
	// Weights      SearchAIOptionWeight
	Filter *Filter
	Limit  int64
}

type SearchAIOptionWeight struct {
	ChunkSimilarity float64 `json:"chunkSimilarity,omitempty"`
	WordSimilarity  float64 `json:"wordSimilarity,omitempty"`
	WordBm25        float64 `json:"wordBm25,omitempty"`
}

type SearchAIDocumentSetResult struct {
	Documents []ai_document_set.SearchDocument
}

type DeleteAIDocumentSetParams struct {
	DocumentSetNames []string
	DocumentSetIds   []string
	Filter           *Filter
}

type DeleteAIDocumentSetResult struct {
	AffectedCount uint64 `json:"affectedCount"`
}

type UpdateAIDocumentSetParams struct {
	DocumentSetId   []string
	DocumentSetName []string
	Filter          *Filter
}

type UpdateAIDocumentSetResult struct {
	AffectedCount uint64 `json:"affectedCount"`
}

type GetCosTmpSecretParams struct {
	DocumentSetName string
}

type GetCosTmpSecretResult struct {
	DocumentSetId           string `json:"documentSetId"`
	DocumentSetName         string `json:"documentSetName"`
	CosEndpoint             string `json:"cosEndpoint"`
	CosRegion               string `json:"cosRegion,omitempty"`
	CosBucket               string `json:"cosBucket,omitempty"`
	UploadPath              string `json:"uploadPath"`
	TmpSecretID             string `json:"tmpSecretId"`
	TmpSecretKey            string `json:"tmpSecretKey"`
	SessionToken            string `json:"token"`
	MaxSupportContentLength int64  `json:"maxSupportContentLength"`
}

type LoadAndSplitTextParams struct {
	DocumentSetName string
	Reader          io.ReadCloser
	LocalFilePath   string
	MetaData        map[string]Field
}

type LoadAndSplitTextResult struct {
	DocumentSetId           string `json:"documentSetId"`
	DocumentSetName         string `json:"documentSetName"`
	CosEndpoint             string `json:"cosEndpoint"`
	CosRegion               string `json:"cosRegion"`
	CosBucket               string `json:"cosBucket"`
	UploadPath              string `json:"uploadPath"`
	TmpSecretID             string `json:"tmpSecretID"`
	TmpSecretKey            string `json:"tmpSecretKey"`
	SessionToken            string `json:"sessionToken"`
	MaxSupportContentLength int64  `json:"maxSupportContentLength"`
}
