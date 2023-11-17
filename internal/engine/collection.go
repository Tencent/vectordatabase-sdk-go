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

package engine

import (
	"context"
	"fmt"
	"strings"
	"time"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity/api"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity/api/collection"
)

var _ entity.CollectionInterface = &implementerCollection{}

type implementerCollection struct {
	entity.SdkClient
	database entity.Database
}

// CreateCollection create a collection. It returns collection struct if err is nil.
// The parameter `name` must be a unique string, otherwise an error will be returned.
// The parameter `shardNum`, `replicasNum` must bigger than 0, `description` could be empty.
// You can set the index field in entity.Indexes, the vectorIndex must be set one currently, and
// the filterIndex sets at least one primaryKey value.
func (i *implementerCollection) CreateCollection(ctx context.Context, name string, shardNum, replicasNum uint32,
	description string, indexes entity.Indexes, options ...*entity.CreateCollectionOption) (*entity.Collection, error) {
	if i.database.IsAIDatabase() {
		return nil, entity.AIDbTypeError
	}
	req := new(collection.CreateReq)
	req.Database = i.database.DatabaseName
	req.Collection = name
	req.ShardNum = shardNum
	req.ReplicaNum = replicasNum
	req.Description = description

	for _, v := range indexes.VectorIndex {
		var column api.IndexColumn
		column.FieldName = v.FieldName
		column.FieldType = string(v.FieldType)
		column.IndexType = string(v.IndexType)
		column.MetricType = string(v.MetricType)
		column.Dimension = v.Dimension

		optionParams(&column, v)

		req.Indexes = append(req.Indexes, &column)
	}
	for _, v := range indexes.FilterIndex {
		var column api.IndexColumn
		column.FieldName = v.FieldName
		column.FieldType = string(v.FieldType)
		if v.FieldType == entity.Array {
			column.FieldElementType = string(v.ElemType)
		}
		column.IndexType = string(v.IndexType)
		req.Indexes = append(req.Indexes, &column)
	}
	if len(options) != 0 && options[0] != nil {
		option := options[0]
		if option.Embedding != nil {
			req.Embedding.Field = option.Embedding.Field
			req.Embedding.VectorField = option.Embedding.VectorField
			req.Embedding.Model = string(option.Embedding.Model)
		}
	}

	res := new(collection.CreateRes)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return nil, err
	}

	coll := i.Collection(req.Collection)
	coll.ShardNum = req.ShardNum
	coll.ReplicasNum = req.ReplicaNum
	coll.Description = req.Description
	coll.Indexes = indexes

	return coll, nil
}

// DescribeCollection get a collection detail.
// It returns the collection object to get collecton parameters or operate document api
func (i *implementerCollection) DescribeCollection(ctx context.Context, name string, option ...*entity.DescribeCollectionOption) (*entity.DescribeCollectionResult, error) {
	if i.database.IsAIDatabase() {
		return nil, entity.AIDbTypeError
	}
	req := new(collection.DescribeReq)
	req.Database = i.database.DatabaseName
	req.Collection = name
	res := new(collection.DescribeRes)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return nil, err
	}
	if res.Collection == nil {
		return nil, fmt.Errorf("get collection %s failed", name)
	}
	coll := i.toCollection(res.Collection)
	result := new(entity.DescribeCollectionResult)
	result.Collection = *coll
	return result, nil
}

// DropCollection drop a collection. If collection not exist, it return nil.
func (i *implementerCollection) DropCollection(ctx context.Context, name string, option ...*entity.DropCollectionOption) (result *entity.DropCollectionResult, err error) {
	if i.database.IsAIDatabase() {
		return nil, entity.AIDbTypeError
	}
	req := new(collection.DropReq)
	req.Database = i.database.DatabaseName
	req.Collection = name

	res := new(collection.DropRes)
	err = i.Request(ctx, req, res)
	result = new(entity.DropCollectionResult)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			return result, nil
		}
		return
	}
	result.AffectedCount = res.AffectedCount
	return
}

func (i *implementerCollection) TruncateCollection(ctx context.Context, name string, option ...*entity.TruncateCollectionOption) (result *entity.TruncateCollectionResult, err error) {
	if i.database.IsAIDatabase() {
		return nil, entity.AIDbTypeError
	}
	req := new(collection.TruncateReq)
	req.Database = i.database.DatabaseName
	req.Collection = name

	res := new(collection.TruncateRes)
	err = i.Request(ctx, req, res)

	if err != nil {
		return
	}
	result = new(entity.TruncateCollectionResult)
	result.AffectedCount = res.AffectedCount
	return
}

// ListCollection get collection list.
// It return the list of collection, each collection same as DescribeCollection return.
func (i *implementerCollection) ListCollection(ctx context.Context, option ...*entity.ListCollectionOption) (*entity.ListCollectionResult, error) {
	if i.database.IsAIDatabase() {
		return nil, entity.AIDbTypeError
	}
	req := new(collection.ListReq)
	req.Database = i.database.DatabaseName
	res := new(collection.ListRes)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return nil, err
	}
	var collections []*entity.Collection
	for _, collection := range res.Collections {
		collections = append(collections, i.toCollection(collection))
	}
	result := new(entity.ListCollectionResult)
	result.Collections = collections
	return result, nil
}

// Collection get a collection interface to operate the document api. It could not send http request to vectordb.
// If you want to show collection parameters, use DescribeCollection.
func (i *implementerCollection) Collection(name string) *entity.Collection {
	coll := new(entity.Collection)
	coll.DatabaseName = i.database.DatabaseName
	coll.CollectionName = name

	docImpl := new(implementerDocument)
	docImpl.SdkClient = i.SdkClient
	docImpl.database = i.database
	docImpl.collection = *coll

	indexImpl := new(implementerIndex)
	indexImpl.SdkClient = i.SdkClient
	indexImpl.database = i.database
	indexImpl.collection = *coll

	coll.DocumentInterface = docImpl
	coll.IndexInterface = indexImpl

	return coll
}

func (i *implementerCollection) toCollection(collectionItem *collection.DescribeCollectionItem) *entity.Collection {
	coll := new(entity.Collection)
	coll.DatabaseName = i.database.DatabaseName
	coll.CollectionName = collectionItem.Collection
	coll.DocumentCount = collectionItem.DocumentCount
	coll.Alias = collectionItem.Alias
	coll.ShardNum = collectionItem.ShardNum
	coll.ReplicasNum = collectionItem.ReplicaNum
	coll.Description = collectionItem.Description
	coll.Size = collectionItem.Size
	if collectionItem.Embedding != nil {
		coll.Embedding.Field = collectionItem.Embedding.Field
		coll.Embedding.VectorField = collectionItem.Embedding.VectorField
		coll.Embedding.Model = entity.EmbeddingModel(collectionItem.Embedding.Model)
		coll.Embedding.Enabled = collectionItem.Embedding.Status == "enabled"
	}

	if collectionItem.IndexStatus != nil {
		coll.IndexStatus = entity.IndexStatus{
			Status: collectionItem.IndexStatus.Status,
		}
		coll.IndexStatus.StartTime, _ = time.Parse("2006-01-02 15:04:05", collectionItem.IndexStatus.StartTime)
	}
	coll.CreateTime, _ = time.Parse("2006-01-02 15:04:05", collectionItem.CreateTime)

	for _, index := range collectionItem.Indexes {
		switch index.FieldType {
		case string(entity.Vector):
			vector := entity.VectorIndex{}
			vector.FieldName = index.FieldName
			vector.FieldType = entity.FieldType(index.FieldType)
			vector.IndexType = entity.IndexType(index.IndexType)
			vector.Dimension = index.Dimension
			vector.MetricType = entity.MetricType(index.MetricType)
			vector.IndexedCount = index.IndexedCount
			switch vector.IndexType {
			case entity.HNSW:
				vector.Params = &entity.HNSWParam{M: index.Params.M, EfConstruction: index.Params.EfConstruction}
			case entity.IVF_FLAT:
				vector.Params = &entity.IVFFLATParams{NList: index.Params.Nlist}
			case entity.IVF_PQ:
				vector.Params = &entity.IVFPQParams{M: index.Params.M, NList: index.Params.Nlist}
			case entity.IVF_SQ8:
				vector.Params = &entity.IVFSQ8Params{NList: index.Params.Nlist}
			}
			coll.Indexes.VectorIndex = append(coll.Indexes.VectorIndex, vector)

		default:
			filter := entity.FilterIndex{}
			filter.FieldName = index.FieldName
			filter.FieldType = entity.FieldType(index.FieldType)
			filter.IndexType = entity.IndexType(index.IndexType)

			coll.Indexes.FilterIndex = append(coll.Indexes.FilterIndex, filter)
		}
	}

	docImpl := new(implementerDocument)
	docImpl.SdkClient = i.SdkClient
	docImpl.database = i.database
	docImpl.collection = *coll
	coll.DocumentInterface = docImpl
	return coll
}

// optionParams option index parameters
func optionParams(column *api.IndexColumn, v entity.VectorIndex) {
	column.Params = new(api.IndexParams)
	if v.IndexType == entity.HNSW {
		if param, ok := v.Params.(*entity.HNSWParam); ok && param != nil {
			column.Params.M = param.M
			column.Params.EfConstruction = param.EfConstruction
		}
	} else if v.IndexType == entity.IVF_FLAT {
		if param, ok := v.Params.(*entity.IVFFLATParams); ok && param != nil {
			column.Params.Nlist = param.NList
		}
	} else if v.IndexType == entity.IVF_SQ8 {
		if param, ok := v.Params.(*entity.IVFSQ8Params); ok && param != nil {
			column.Params.Nlist = param.NList
		}
	} else if v.IndexType == entity.IVF_PQ {
		if param, ok := v.Params.(*entity.IVFPQParams); ok && param != nil {
			column.Params.M = param.M
			column.Params.Nlist = param.NList
		}
	}
}
