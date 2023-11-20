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
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb/api/ai_collection"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb/api/collection"
)

var _ AICollectionInterface = &implementerAICollection{}

type implementerAICollection struct {
	SdkClient
	database AIDatabase
}

// CreateCollection create a collection. It returns collection struct if err is nil.
// The parameter `name` must be a unique string, otherwise an error will be returned.
// The parameter `shardNum`, `replicasNum` must bigger than 0, `description` could be empty.
// You can set the index field in Indexes, the vectorIndex must be set one currently, and
// the filterIndex sets at least one primaryKey value.
func (i *implementerAICollection) CreateCollection(ctx context.Context, name string, options ...*CreateAICollectionOption) (*CreateAICollectionResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(ai_collection.CreateReq)
	req.Database = i.database.DatabaseName
	req.Collection = name

	if len(options) != 0 && options[0] != nil {
		option := options[0]
		req.Description = option.Description

		for _, v := range option.Indexes.FilterIndex {
			var column api.IndexColumn
			column.FieldName = v.FieldName
			column.FieldType = string(v.FieldType)
			column.IndexType = string(v.IndexType)
			req.Indexes = append(req.Indexes, column)
		}

		// defaultEnableWordsSimilarity := true
		if option.AiConfig != nil {
			req.ExpectedFileNum = option.AiConfig.ExpectedFileNum
			req.AverageFileSize = option.AiConfig.AverageFileSize
			req.Language = string(option.AiConfig.Language)
			if option.AiConfig.DocumentPreprocess != nil {
				req.DocumentPreprocess = option.AiConfig.DocumentPreprocess
			}
			req.EnableWordsEmbedding = option.AiConfig.EnableWordsEmbedding
			// if option.AiConfig.DocumentIndex != nil && option.AiConfig.DocumentIndex.EnableWordsSimilarity != nil {
			// 	req.DocumentIndex = option.AiConfig.DocumentIndex
			// } else {
			// 	req.DocumentIndex = &ai_collection.DocumentIndex{
			// 		EnableWordsSimilarity: &defaultEnableWordsSimilarity,
			// 	}
			// }
		}
	}

	res := new(collection.CreateRes)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return nil, err
	}

	coll := i.toCollection(&ai_collection.DescribeAICollectionItem{
		Database:           req.Database,
		Collection:         req.Collection,
		Language:           req.Language,
		Description:        req.Description,
		FilterIndexes:      req.Indexes,
		DocumentPreprocess: req.DocumentPreprocess,
		// DocumentIndex:      *req.DocumentIndex,
	})
	result := new(CreateAICollectionResult)
	result.AICollection = *coll

	return result, nil
}

// DescribeCollection get a collection detail.
// It returns the collection object to get collecton parameters or operate document api
func (i *implementerAICollection) DescribeCollection(ctx context.Context, name string, option ...*DescribeAICollectionOption) (*DescribeAICollectionResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(ai_collection.DescribeReq)
	req.Database = i.database.DatabaseName
	req.Collection = name
	res := new(ai_collection.DescribeRes)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return nil, err
	}
	if res.Collection == nil {
		return nil, fmt.Errorf("get collection %s failed", name)
	}
	coll := i.toCollection(res.Collection)
	result := new(DescribeAICollectionResult)
	result.AICollection = *coll
	return result, nil
}

// DropCollection drop a collection. If collection not exist, it return nil.
func (i *implementerAICollection) DropCollection(ctx context.Context, name string, option ...*DropAICollectionOption) (result *DropAICollectionResult, err error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(ai_collection.DropReq)
	req.Database = i.database.DatabaseName
	req.Collection = name

	res := new(ai_collection.DropRes)
	err = i.Request(ctx, req, res)
	result = new(DropAICollectionResult)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			return result, nil
		}
		return
	}
	result.AffectedCount = int(res.AffectedCount)
	return
}

func (i *implementerAICollection) TruncateCollection(ctx context.Context, name string, option ...*TruncateAICollectionOption) (result *TruncateAICollectionResult, err error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(ai_collection.TruncateReq)
	req.Database = i.database.DatabaseName
	req.Collection = name

	res := new(ai_collection.TruncateRes)
	err = i.Request(ctx, req, res)

	if err != nil {
		return
	}
	result = new(TruncateAICollectionResult)
	result.AffectedCount = int(res.AffectedCount)
	return
}

// ListCollection get collection list.
// It return the list of collection, each collection same as DescribeCollection return.
func (i *implementerAICollection) ListCollection(ctx context.Context, option ...*ListAICollectionOption) (*ListAICollectionResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(ai_collection.ListReq)
	req.Database = i.database.DatabaseName
	res := new(ai_collection.ListRes)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return nil, err
	}
	result := new(ListAICollectionResult)
	for _, item := range res.Collections {
		result.Collections = append(result.Collections, i.toCollection(&item))
	}
	return result, nil
}

// Collection get a collection interface to operate the document api. It could not send http request to vectordb.
// If you want to show collection parameters, use DescribeCollection.
func (i *implementerAICollection) Collection(name string) *AICollection {
	coll := new(AICollection)
	coll.DatabaseName = i.database.DatabaseName
	coll.CollectionName = name

	docImpl := new(implementerAIDocument)
	docImpl.SdkClient = i.SdkClient
	docImpl.database = i.database
	docImpl.collection = *coll

	coll.AIDocumentInterface = docImpl
	return coll
}

func (i *implementerAICollection) toCollection(item *ai_collection.DescribeAICollectionItem) *AICollection {
	coll := new(AICollection)
	coll.DatabaseName = i.database.DatabaseName
	coll.CollectionName = item.Collection
	coll.Description = item.Description
	coll.Alias = item.Alias
	coll.CreateTime, _ = time.Parse("2006-01-02 15:04:05", item.CreateTime)

	coll.AiConfig = AiConfig{
		ExpectedFileNum:    item.ExpectedFileNum,
		AverageFileSize:    item.AverageFileSize,
		Language:           Language(item.Language),
		DocumentPreprocess: item.DocumentPreprocess,
		// DocumentIndex:      &item.DocumentIndex,
	}
	if item.AiStatus != nil {
		coll.IndexedDocuments = item.AiStatus.IndexedDocuments
		coll.TotalDocuments = item.AiStatus.TotalDocuments
		coll.UnIndexedDocuments = item.AiStatus.UnIndexedDocuments
	}
	for _, index := range item.FilterIndexes {
		filter := FilterIndex{}
		filter.FieldName = index.FieldName
		filter.FieldType = FieldType(index.FieldType)
		filter.IndexType = IndexType(index.IndexType)

		coll.FilterIndexes = append(coll.FilterIndexes, filter)
	}

	docImpl := new(implementerAIDocument)
	docImpl.SdkClient = i.SdkClient
	docImpl.database = i.database
	docImpl.collection = *coll
	coll.AIDocumentInterface = docImpl
	return coll
}
