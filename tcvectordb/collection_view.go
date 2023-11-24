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

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb/api"
	collection_view "git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb/api/collection_view"
)

type CollectionViewInterface interface {
	SdkClient
	CreateCollectionView(ctx context.Context, name string, options ...*CreateCollectionViewOption) (result *CreateAICollectionResult, err error)
	DescribeCollectionView(ctx context.Context, name string, options ...*DescribeCollectionViewOption) (result *DescribeCollectionViewResult, err error)
	DropCollectionView(ctx context.Context, name string, options ...*DropCollectionViewOption) (result *DropCollectionViewResult, err error)
	TruncateCollectionView(ctx context.Context, name string, options ...*TruncateCollectionViewOption) (result *TruncateCollectionViewResult, err error)
	ListCollectionViews(ctx context.Context, options ...*ListCollectionViewsOption) (result *ListCollectionViewsResult, err error)
	CollectionView(name string) *CollectionView
}

var _ CollectionViewInterface = &implementerCollectionView{}

type implementerCollectionView struct {
	SdkClient
	database AIDatabase
}

// CreateCollectionView create a collectionView. It returns collection struct if err is nil.
// The parameter `name` must be a unique string, otherwise an error will be returned.
// The parameter `description` could be empty.
// You can set the index field in Indexes.
func (i *implementerCollectionView) CreateCollectionView(ctx context.Context, name string, options ...*CreateCollectionViewOption) (*CreateAICollectionResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(collection_view.CreateReq)
	req.Database = i.database.DatabaseName
	req.CollectionView = name

	if len(options) != 0 && options[0] != nil {
		option := options[0]
		req.Description = option.Description

		for _, v := range option.Indexes.FilterIndex {
			var column api.IndexColumn
			column.FieldName = v.FieldName
			column.FieldType = string(v.FieldType)
			column.IndexType = string(v.IndexType)
			req.Indexes = append(req.Indexes, &column)
		}

		if option.Embedding != nil {
			req.Embedding = new(collection_view.DocumentEmbedding)
			req.Embedding.Language = string(option.Embedding.Language)
			req.Embedding.EnableWordsEmbedding = option.Embedding.EnableWordsEmbedding
		}
		if option.SplitterPreprocess != nil {
			req.SplitterPreprocess = new(collection_view.SplitterPreprocess)
			req.SplitterPreprocess.AppendTitleToChunk = option.SplitterPreprocess.AppendTitleToChunk
			req.SplitterPreprocess.AppendKeywordsToChunk = option.SplitterPreprocess.AppendKeywordsToChunk
		}
	}

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
	result := new(CreateAICollectionResult)
	result.CollectionView = *coll

	return result, nil
}

// DescribeCollectionView get a collectionView detail.
// It returns the collectionView object to get collectionView parameters or operate document api
func (i *implementerCollectionView) DescribeCollectionView(ctx context.Context, name string, option ...*DescribeCollectionViewOption) (*DescribeCollectionViewResult, error) {
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
	result := new(DescribeCollectionViewResult)
	result.CollectionView = *coll
	return result, nil
}

// DropCollectionView drop a collectionView. If collectionView not exist, it return nil.
func (i *implementerCollectionView) DropCollectionView(ctx context.Context, name string, option ...*DropCollectionViewOption) (result *DropCollectionViewResult, err error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(collection_view.DropReq)
	req.Database = i.database.DatabaseName
	req.CollectionView = name

	res := new(collection_view.DropRes)
	err = i.Request(ctx, req, res)
	result = new(DropCollectionViewResult)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			return result, nil
		}
		return
	}
	result.AffectedCount = int(res.AffectedCount)
	return
}

func (i *implementerCollectionView) TruncateCollectionView(ctx context.Context, name string, option ...*TruncateCollectionViewOption) (result *TruncateCollectionViewResult, err error) {
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
	result = new(TruncateCollectionViewResult)
	result.AffectedCount = int(res.AffectedCount)
	return
}

// ListCollectionViews get collectionView list.
// It return the list of collectionView, each collectionView is as same as DescribeCollectionView return.
func (i *implementerCollectionView) ListCollectionViews(ctx context.Context, option ...*ListCollectionViewsOption) (*ListCollectionViewsResult, error) {
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
	result := new(ListCollectionViewsResult)
	for _, item := range res.CollectionViews {
		result.CollectionViews = append(result.CollectionViews, i.toCollectionView(item))
	}
	return result, nil
}

// CollectionView get a collectionView interface to operate the document api. It could not send http request to vectordb.
// If you want to show collectionView parameters, use DescribeCollectionView.
func (i *implementerCollectionView) CollectionView(name string) *CollectionView {
	coll := new(CollectionView)
	coll.DatabaseName = i.database.DatabaseName
	coll.CollectionName = name

	docImpl := new(implementerAIDocumentSet)
	docImpl.SdkClient = i.SdkClient
	docImpl.database = i.database
	docImpl.collectionView = *coll

	coll.AIDocumentSetInterface = docImpl
	return coll
}

func (i *implementerCollectionView) toCollectionView(item *collection_view.DescribeCollectionViewItem) *CollectionView {
	coll := new(CollectionView)
	coll.DatabaseName = i.database.DatabaseName
	coll.CollectionName = item.CollectionView
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

	docImpl := new(implementerAIDocumentSet)
	docImpl.SdkClient = i.SdkClient
	docImpl.database = i.database
	docImpl.collectionView = *coll
	coll.AIDocumentSetInterface = docImpl
	return coll
}
