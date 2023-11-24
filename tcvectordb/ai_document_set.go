package tcvectordb

import (
	"context"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb/api/ai_document_set"
)

type AIDocumentSetInterface interface {
	Search(ctx context.Context, param SearchAIDocumentSetParams) (*SearchAIDocumentSetResult, error)
	Delete(ctx context.Context) (*DeleteAIDocumentSetResult, error)
}

type implementerAIDocumentSet struct {
	SdkClient
	database       *AIDatabase
	collectionView *AICollectionView
	documentSet    *AIDocumentSet
}

type SearchAIDocumentSetParams struct {
	Content      string                        `json:"content"`
	ExpandChunk  []int                         `json:"expandChunk"` // 搜索结果中，向前、向后补齐几个chunk的上下文
	RerankOption *ai_document_set.RerankOption `json:"rerankOption"`
}

func (i *implementerAIDocumentSet) Search(ctx context.Context, param SearchAIDocumentSetParams) (*SearchAIDocumentSetResult, error) {
	return i.collectionView.Search(ctx, SearchAIDocumentSetsParams{
		Content:         param.Content,
		DocumentSetName: []string{i.documentSet.DocumentSetName},
		ExpandChunk:     param.ExpandChunk,
		RerankOption:    param.RerankOption,
	})
}

func (i *implementerAIDocumentSet) Delete(ctx context.Context) (*DeleteAIDocumentSetResult, error) {
	return i.collectionView.DeleteByIds(ctx, i.documentSet.DocumentSetId)
}
