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
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/collection_view"
)

var _ AICollectionViewInterface = &implementerCollectionView{}

type AICollectionViewInterface interface {
	SdkClient

	// [CreateCollectionView] creates and initializes a new collectionView.
	CreateCollectionView(ctx context.Context, name string, param CreateCollectionViewParams) (result *CreateAICollectionViewResult, err error)

	// [ListCollectionViews] retrieves the list of all collectionViews in the AI database.
	ListCollectionViews(ctx context.Context) (result *ListAICollectionViewsResult, err error)

	// [DescribeCollectionView] retrieves information about a specific collectionView. See [AICollectionView] for more information.
	DescribeCollectionView(ctx context.Context, name string) (result *DescribeAICollectionViewResult, err error)

	// [DropCollectionView] drops a specific collectionView.
	DropCollectionView(ctx context.Context, name string) (result *DropAICollectionViewResult, err error)

	// [TruncateCollectionView] clears all the data and indexes in the collectionView.
	TruncateCollectionView(ctx context.Context, name string) (result *TruncateAICollectionViewResult, err error)

	// [CollectionView] returns a pointer to a [AICollectionView] object. which includes the collectionView parameters
	// and some interfaces to operate the documentSet api.
	CollectionView(name string) *AICollectionView
}

// [AICollectionView] holds the collectionView parameters and some interfaces to operate the documentSet api.
//
// Fields:
//   - DatabaseName: The name of the database.
//   - CollectionViewName: The name of the collection.
//   - Alias: All aliases of the CollectionView.
//   - Embedding: A pointer to a [DocumentEmbedding] object, which includes the parameters for embedding.
//     See [DocumentEmbedding] for more information.
//   - SplitterPreprocess: A pointer to a [SplitterPreprocess] object, which includes the parameters
//     for splitting document chunks. See [SplitterPreprocess] for more information.
//   - ParsingProcess: A pointer to a [ParsingProcess] object, which includes the parameters
//     for parsing files. See [ParsingProcess] for more information.
//   - IndexedDocumentSets: The number of documentSets that have been processed.
//   - TotalDocumentSets: The total number of documentSets in this collectionView.
//   - UnIndexedDocumentSets: The number of documentSets that haven't been processed.
//   - FilterIndexes: A [Indexes] object that includes a list of the scalar filter index properties for the documentSets in a collectionView.
//   - Description: (Optional) The description of the collection.
//   - CreateTime: The create time of collectionView.
type AICollectionView struct {
	AIDocumentSetsInterface `json:"-"`
	DatabaseName            string                              `json:"databaseName"`
	CollectionViewName      string                              `json:"collectionViewName"`
	Alias                   []string                            `json:"alias"`
	Embedding               *collection_view.DocumentEmbedding  `json:"embedding"`
	SplitterPreprocess      *collection_view.SplitterPreprocess `json:"splitterPreprocess"`
	ParsingProcess          *api.ParsingProcess                 `json:"parsingProcess"`
	IndexedDocumentSets     uint64                              `json:"indexedDocumentSets"`
	TotalDocumentSets       uint64                              `json:"totalDocumentSets"`
	UnIndexedDocumentSets   uint64                              `json:"unIndexedDocumentSets"`
	FilterIndexes           []FilterIndex                       `json:"filterIndexes"`
	Description             string                              `json:"description"`
	CreateTime              time.Time                           `json:"createTime"`
}

type implementerCollectionView struct {
	SdkClient
	database *AIDatabase
}

// [CreateCollectionViewParams] holds the parameters for creating a new collectionView.
//
// Fields:
//   - Description: (Optional) The description of the collection.
//   - Indexes: A [Indexes] object that includes a list of the scalar field index properties for the documentSets in a collectionView.
//   - Embedding: A pointer to a [DocumentEmbedding] object, which includes the parameters for embedding.
//     See [DocumentEmbedding] for more information.
//   - SplitterPreprocess: A pointer to a [SplitterPreprocess] object, which includes the parameters
//     for splitter process. See [SplitterPreprocess] for more information.
//   - ParsingProcess: A pointer to a [ParsingProcess] object, which includes the parameters
//     for parsing parameters. See [ParsingProcess] for more information.
//   - ExpectedFileNum: Expected total number of documents.
//   - AverageFileSize: Estimate the average document size.
type CreateCollectionViewParams struct {
	Description        string
	Indexes            Indexes
	Embedding          *collection_view.DocumentEmbedding
	SplitterPreprocess *collection_view.SplitterPreprocess
	ParsingProcess     *api.ParsingProcess
	ExpectedFileNum    uint64
	AverageFileSize    uint64
}

type CreateAICollectionViewResult struct {
	AICollectionView
	AffectedCount int
}

// [CreateCollectionView] creates and initializes a new collectionView.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - name: The name of the collectionView to create. CollectionView name must be 1-128 characters long,
//     start with an alphanumeric character,
//     and consist only of alphanumeric characters, numbers, '_' or '-'.
//   - params: A pointer to a [CreateCollectionViewParams] object that includes the other parameters for the collectionView.
//     See [CreateCollectionViewParams] for more information.
//
// Returns a pointer to a [CreateAICollectionViewResult] object or an error.
func (i *implementerCollectionView) CreateCollectionView(ctx context.Context, name string, param CreateCollectionViewParams) (*CreateAICollectionViewResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(collection_view.CreateReq)
	req.Database = i.database.DatabaseName
	req.CollectionView = name

	req.Description = param.Description

	for _, v := range param.Indexes.FilterIndex {
		var column api.IndexColumn
		column.FieldName = v.FieldName
		column.FieldType = string(v.FieldType)
		column.IndexType = string(v.IndexType)
		req.Indexes = append(req.Indexes, &column)
	}

	if param.Embedding != nil {
		req.Embedding = new(collection_view.DocumentEmbedding)
		req.Embedding.Language = string(param.Embedding.Language)
		req.Embedding.EnableWordsEmbedding = param.Embedding.EnableWordsEmbedding
	}
	if param.SplitterPreprocess != nil {
		req.SplitterPreprocess = new(collection_view.SplitterPreprocess)
		req.SplitterPreprocess.AppendTitleToChunk = param.SplitterPreprocess.AppendTitleToChunk
		req.SplitterPreprocess.AppendKeywordsToChunk = param.SplitterPreprocess.AppendKeywordsToChunk
	}
	if param.ParsingProcess != nil {
		req.ParsingProcess = param.ParsingProcess
	}
	req.AverageFileSize = param.AverageFileSize
	req.ExpectedFileNum = param.ExpectedFileNum

	res := new(collection_view.CreateRes)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return nil, err
	}

	coll := i.toCollectionView(&collection_view.DescribeCollectionViewItem{
		Database:           req.Database,
		CollectionView:     req.CollectionView,
		Description:        req.Description,
		Embedding:          req.Embedding,
		SplitterPreprocess: req.SplitterPreprocess,
		Indexes:            req.Indexes,
	})
	result := new(CreateAICollectionViewResult)
	result.AICollectionView = *coll

	return result, nil
}

type DescribeAICollectionViewResult struct {
	AICollectionView
}

// [ListCollectionViews] retrieves the list of all collectionViews in the AI database.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//
// Notes: The database name is from the field of [implementerCollectionView].
//
// Returns a pointer to a [ListAICollectionViewsResult] object or an error.
func (i *implementerCollectionView) ListCollectionViews(ctx context.Context) (*ListAICollectionViewsResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(collection_view.ListReq)
	req.Database = i.database.DatabaseName
	res := new(collection_view.ListRes)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return nil, err
	}
	result := new(ListAICollectionViewsResult)
	for _, item := range res.CollectionViews {
		result.CollectionViews = append(result.CollectionViews, i.toCollectionView(item))
	}
	return result, nil
}

// [DescribeCollectionView] retrieves information about a specific collectionView. See [AICollectionView] for more information.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - name: The name of the collectionView to describe.
//
// Notes: The database name is from the field of [implementerCollectionView].
//
// Returns a pointer to a [ListAICollectionViewsResult] object or an error.
func (i *implementerCollectionView) DescribeCollectionView(ctx context.Context, name string) (*DescribeAICollectionViewResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(collection_view.DescribeReq)
	req.Database = i.database.DatabaseName
	req.CollectionView = name
	res := new(collection_view.DescribeRes)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return nil, err
	}
	if res.CollectionView == nil {
		return nil, fmt.Errorf("get collectionView %s failed", name)
	}
	coll := i.toCollectionView(res.CollectionView)
	result := new(DescribeAICollectionViewResult)
	result.AICollectionView = *coll
	return result, nil
}

type DropAICollectionViewResult struct {
	AffectedCount int
}

// [DropCollectionView] drops a specific collectionView.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - name: The name of the collectionView to drop.
//
// Notes: The database name is from the field of [implementerCollectionView].
//
// Returns a pointer to a [DropAICollectionViewResult] object or an error.
func (i *implementerCollectionView) DropCollectionView(ctx context.Context, name string) (result *DropAICollectionViewResult, err error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(collection_view.DropReq)
	req.Database = i.database.DatabaseName
	req.CollectionView = name

	res := new(collection_view.DropRes)
	err = i.Request(ctx, req, res)
	result = new(DropAICollectionViewResult)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			return result, nil
		}
		return
	}
	result.AffectedCount = int(res.AffectedCount)
	return
}

type TruncateAICollectionViewResult struct {
	AffectedCount int
}

// [TruncateCollectionView] clears all the data and indexes in the collectionView.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - name: The name of the collectionView to truncate.
//
// Notes: The database name is from the field of [implementerCollectionView].
//
// Returns a pointer to a [TruncateAICollectionViewResult] object or an error.
func (i *implementerCollectionView) TruncateCollectionView(ctx context.Context, name string) (result *TruncateAICollectionViewResult, err error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(collection_view.TruncateReq)
	req.Database = i.database.DatabaseName
	req.CollectionView = name

	res := new(collection_view.TruncateRes)
	err = i.Request(ctx, req, res)

	if err != nil {
		return
	}
	result = new(TruncateAICollectionViewResult)
	result.AffectedCount = int(res.AffectedCount)
	return
}

type ListAICollectionViewsResult struct {
	CollectionViews []*AICollectionView `json:"collectionViews"`
}

// [CollectionView] returns a pointer to a [AICollectionView] object. which includes the collectionView parameters
// and some interfaces to operate the documentSet api.
//
// Parameters:
//   - name: The name of the collectionView to truncate.
//
// Notes: The database name is from the field of [implementerCollectionView].
//
// Returns a pointer to a [AICollectionView] object.
func (i *implementerCollectionView) CollectionView(name string) *AICollectionView {
	coll := new(AICollectionView)
	coll.DatabaseName = i.database.DatabaseName
	coll.CollectionViewName = name

	docImpl := new(implementerAIDocumentSets)
	docImpl.SdkClient = i.SdkClient
	docImpl.database = i.database
	docImpl.collectionView = coll

	coll.AIDocumentSetsInterface = docImpl
	return coll
}

func (i *implementerCollectionView) toCollectionView(item *collection_view.DescribeCollectionViewItem) *AICollectionView {
	coll := new(AICollectionView)
	coll.DatabaseName = i.database.DatabaseName
	coll.CollectionViewName = item.CollectionView
	coll.Description = item.Description
	coll.Alias = item.Alias
	coll.CreateTime, _ = time.Parse("2006-01-02 15:04:05", item.CreateTime)

	if item.Embedding != nil {
		coll.Embedding = &collection_view.DocumentEmbedding{
			Language:             item.Embedding.Language,
			EnableWordsEmbedding: item.Embedding.EnableWordsEmbedding,
		}
	}
	if item.SplitterPreprocess != nil {
		coll.SplitterPreprocess = &collection_view.SplitterPreprocess{
			AppendTitleToChunk:    item.SplitterPreprocess.AppendTitleToChunk,
			AppendKeywordsToChunk: item.SplitterPreprocess.AppendKeywordsToChunk,
		}
	}
	if item.ParsingProcess != nil {
		coll.ParsingProcess = item.ParsingProcess
	}

	if item.Status != nil {
		coll.IndexedDocumentSets = item.Status.IndexedDocumentSets
		coll.TotalDocumentSets = item.Status.TotalDocumentSets
		coll.UnIndexedDocumentSets = item.Status.UnIndexedDocumentSets
	}
	for _, index := range item.Indexes {
		filter := FilterIndex{}
		filter.FieldName = index.FieldName
		filter.FieldType = FieldType(index.FieldType)
		filter.IndexType = IndexType(index.IndexType)

		coll.FilterIndexes = append(coll.FilterIndexes, filter)
	}

	docImpl := new(implementerAIDocumentSets)
	docImpl.SdkClient = i.SdkClient
	docImpl.database = i.database
	docImpl.collectionView = coll
	coll.AIDocumentSetsInterface = docImpl
	return coll
}
