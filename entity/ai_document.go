package entity

import "git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity/api/ai_document"

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
	Documents     []ai_document.QueryDocument
}

type SearchAIDocumentOption struct {
	FileName    string
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
