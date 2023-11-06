package entity

import "git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity/api/ai_document"

type QueryAIDocumentOption struct {
	DocumentIds  []string
	Filter       *Filter
	Limit        int64
	Offset       int64
	OutputFields []string
}

type QueryAIDocumentsResult struct {
	AffectedCount int
	Total         int
	Documents     []ai_document.QueryDocument
}

type SearchAIDocumentOption struct {
	Filter      *Filter
	ResultType  string
	ChunkExpand []int
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
	DocumentIds []string
	Filter      *Filter
}

type DeleteAIDocumentResult struct {
	AffectedCount uint64 `json:"affectedCount"`
}

type UpdateAIDocumentOption struct {
	QueryIds     []string
	QueryFilter  Filter
	UpdateFields map[string]interface{}
}

type UpdateAIDocumentResult struct {
	AffectedCount uint64 `json:"affectedCount"`
}

type GetCosTmpSecretOption struct {
	FileType FileType
}

type GetCosTmpSecretResult struct {
	CosEndpoint             string `json:"cosEndpoint"`
	CosRegion               string `json:"cosRegion,omitempty"`
	CosBucket               string `json:"cosBucket,omitempty"`
	UploadPath              string `json:"uploadPath"`
	TmpSecretID             string `json:"tmpSecretId"`
	TmpSecretKey            string `json:"tmpSecretKey"`
	SessionToken            string `json:"token"`
	MaxSupportContentLength int64  `json:"maxSupportContentLength"`
	FileId                  string `json:"fileId"`
}

type UploadAIDocumentOption struct {
	FileType FileType
	MetaData map[string]Field
}

type UploadAIDocumentResult struct {
	CosEndpoint string
	UploadPath  string
	FileId      string
}
