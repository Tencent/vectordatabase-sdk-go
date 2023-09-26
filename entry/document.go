package entry

import (
	"context"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/model"
)

// DocumentInterface document api
type DocumentInterface interface {
	SdkClient
	Upsert(ctx context.Context, documents []model.Document, option *UpsertDocumentOption) (result *DocumentResult, err error)
	Query(ctx context.Context, documentIds []string, option *QueryDocumentOption) (docs []model.Document, result *DocumentResult, err error)
	Search(ctx context.Context, vectors [][]float32, option *SearchDocumentOption) ([][]model.Document, error)
	SearchById(ctx context.Context, documentIds []string, option *SearchDocumentOption) ([][]model.Document, error)
	SearchByText(ctx context.Context, text map[string][]string, option *SearchDocumentOption) ([][]model.Document, error)
	Delete(ctx context.Context, option *DeleteDocumentOption) (result *DocumentResult, err error)
	Update(ctx context.Context, option *UpdateDocumentOption) (result *DocumentResult, err error)
}

type DocumentResult struct {
	AffectedCount int
	Total         int
}

type UpsertDocumentOption struct {
	BuildIndex bool
}

type QueryDocumentOption struct {
	Filter         *model.Filter
	RetrieveVector bool
	OutputFields   []string
	Offset         int64
	Limit          int64
}

type SearchDocumentOption struct {
	Filter         *model.Filter
	Params         *SearchDocParams
	RetrieveVector bool
	OutputFields   []string
	Limit          int64
}

type SearchDocParams struct {
	Nprobe uint32  `json:"nprobe,omitempty"` // 搜索时查找的聚类数量，使用索引默认值即可
	Ef     uint32  `json:"ef,omitempty"`     // HNSW
	Radius float32 `json:"radius,omitempty"` // 距离阈值,范围搜索时有效
}

type DeleteDocumentOption struct {
	DocumentIds []string
	Filter      *model.Filter
}

type UpdateDocumentOption struct {
	QueryIds     []string
	QueryFilter  *model.Filter
	UpdateVector []float32
	UpdateFields map[string]model.Field
}
