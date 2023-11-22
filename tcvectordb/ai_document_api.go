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

import "git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb/api/ai_document"

type QueryAIDocumentOption struct {
	FileName     string
	DocumentIds  []string
	Filter       *Filter
	Limit        int64
	Offset       int64
	OutputFields []string
}

type QueryAIDocumentsResult struct {
	AffectedCount int
	Total         int
	Documents     []ai_document.QueryDocumentSet
}

type GetAIDocumentOption struct {
	DocumentSetId   string
	DocumentSetName string
}

type GetAIDocumentResult struct {
	Count        uint64
	DocumentSets ai_document.GetDocumentSet `json:"documentSet"`
}

type SearchAIDocumentOption struct {
	FileName     string
	Filter       *Filter
	ResultType   string
	ChunkExpand  []int
	RerankOption *ai_document.RerankOption // 多路召回
	// MergeChunk  bool
	// Weights      SearchAIOptionWeight
	OutputFields []string
	Limit        int64
}

type SearchAIOptionWeight struct {
	ChunkSimilarity float64 `json:"chunkSimilarity,omitempty"`
	WordSimilarity  float64 `json:"wordSimilarity,omitempty"`
	WordBm25        float64 `json:"wordBm25,omitempty"`
}

type SearchAIDocumentResult struct {
	Documents []ai_document.SearchDocument
}

type DeleteAIDocumentOption struct {
	FileName    string
	DocumentIds []string
	Filter      *Filter
}

type DeleteAIDocumentResult struct {
	AffectedCount uint64 `json:"affectedCount"`
}

type UpdateAIDocumentOption struct {
	FileName     string
	QueryIds     []string
	QueryFilter  *Filter
	UpdateFields map[string]interface{}
}

type UpdateAIDocumentResult struct {
	AffectedCount uint64 `json:"affectedCount"`
}

type GetCosTmpSecretOption struct {
	FileType FileType
}

type GetCosTmpSecretResult struct {
	FileName                string `json:"fileName"`
	FileId                  string `json:"fileId"`
	CosEndpoint             string `json:"cosEndpoint"`
	CosRegion               string `json:"cosRegion,omitempty"`
	CosBucket               string `json:"cosBucket,omitempty"`
	UploadPath              string `json:"uploadPath"`
	TmpSecretID             string `json:"tmpSecretId"`
	TmpSecretKey            string `json:"tmpSecretKey"`
	SessionToken            string `json:"token"`
	MaxSupportContentLength int64  `json:"maxSupportContentLength"`
}

type UploadAIDocumentOption struct {
	FileType FileType
	MetaData map[string]Field
}

type UploadAIDocumentResult struct {
	FileName                string `json:"fileName"`
	FileId                  string `json:"fileId"`
	CosEndpoint             string `json:"cosEndpoint"`
	CosRegion               string `json:"cosRegion"`
	CosBucket               string `json:"cosBucket"`
	UploadPath              string `json:"uploadPath"`
	TmpSecretID             string `json:"tmpSecretID"`
	TmpSecretKey            string `json:"tmpSecretKey"`
	SessionToken            string `json:"sessionToken"`
	MaxSupportContentLength int64  `json:"maxSupportContentLength"`
}
