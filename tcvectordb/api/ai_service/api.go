package ai_service

import "github.com/tencent/vectordatabase-sdk-go/tcvectordb/api"

type ModelParams struct {
	RetrieveDenseVector  *bool `json:"retrieveDenseVector,omitempty"`
	RetrieveSparseVector *bool `json:"retrieveSparseVector,omitempty"`
}

type EmbeddingReq struct {
	api.Meta    `path:"/ai/service/embedding" tags:"AI Service" method:"Post" summary:"embedding"`
	Model       string       `json:"model,omitempty"`
	ModelParams *ModelParams `json:"modelParams,omitempty"`
	DataType    string       `json:"dataType,omitempty"`
	Data        []string     `json:"data,omitempty"`
}

type EmbeddingRes struct {
	api.CommonRes
	TokenUsed    int64                `json:"tokenUsed,omitempty"`
	DenseVector  [][]float32          `json:"denseVector,omitempty"`
	SparseVector []map[string]float32 `json:"sparseVector,omitempty"`
}
