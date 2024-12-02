package tcvectordb

import (
	"context"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/ai_document_set"
)

type AIDocumentSetInterface interface {
	// [Search] returns the most similar topK chunks from a specific documentSet(file).
	Search(ctx context.Context, param SearchAIDocumentSetParams) (*SearchAIDocumentSetResult, error)

	// [Delete] deletes the documentSet.
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

// [Search] returns the most similar topK chunks from a specific documentSet(file).
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - param: A pointer to a [SearchAIDocumentSetsParams] object that includes the other parameters for searching documentSets' operation.
//     See [SearchAIDocumentSetsParams] for more information.
//
// Notes: The name of the database, the name of collectionView and the name of documentSet are from
// the fields of [implementerAIDocumentSets].
//
// Returns a pointer to a [SearchAIDocumentSetResult] object or an error.
func (i *implementerAIDocumentSet) Search(ctx context.Context, param SearchAIDocumentSetParams) (*SearchAIDocumentSetResult, error) {
	return i.collectionView.Search(ctx, SearchAIDocumentSetsParams{
		Content:         param.Content,
		DocumentSetName: []string{i.documentSet.DocumentSetName},
		ExpandChunk:     param.ExpandChunk,
		RerankOption:    param.RerankOption,
	})
}

// [Delete] deletes the documentSet.
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//
// Notes: The name of the database, the name of collectionView and the name of documentSet are from
// the fields of [implementerAIDocumentSets].
//
// Returns a pointer to a [DeleteAIDocumentSetResult] object or an error.
func (i *implementerAIDocumentSet) Delete(ctx context.Context) (*DeleteAIDocumentSetResult, error) {
	return i.collectionView.DeleteByIds(ctx, i.documentSet.DocumentSetId)
}
