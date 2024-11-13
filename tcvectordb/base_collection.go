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
	"strconv"
	"strings"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/collection"
)

var _ CollectionInterface = &implementerCollection{}

// CollectionInterface collection api
type CollectionInterface interface {
	SdkClient
	ExistsCollection(ctx context.Context, name string) (bool, error)
	CreateCollectionIfNotExists(ctx context.Context, name string, shardNum, replicasNum uint32, description string,
		indexes Indexes, params ...*CreateCollectionParams) (*Collection, error)
	CreateCollection(ctx context.Context, name string, shardNum, replicasNum uint32, description string,
		indexes Indexes, params ...*CreateCollectionParams) (*Collection, error)
	ListCollection(ctx context.Context) (result *ListCollectionResult, err error)
	DescribeCollection(ctx context.Context, name string) (result *DescribeCollectionResult, err error)
	DropCollection(ctx context.Context, name string) (result *DropCollectionResult, err error)
	TruncateCollection(ctx context.Context, name string) (result *TruncateCollectionResult, err error)
	Collection(name string) *Collection
}

type implementerCollection struct {
	SdkClient
	database *Database
}

type CreateCollectionParams struct {
	Embedding *Embedding
	TtlConfig *TtlConfig
}

type CreateCollectionResult struct {
	Collection
}

func (i *implementerCollection) ExistsCollection(ctx context.Context, name string) (bool, error) {
	res, err := i.DescribeCollection(ctx, name)
	if err != nil {
		if strings.Contains(err.Error(), strconv.Itoa(ERR_UNDEFINED_COLLECTION)) {
			return false, nil
		}
		return false, fmt.Errorf("get collection %s failed, err: %v", name, err.Error())
	}
	if res == nil {
		return false, fmt.Errorf("get collection %s failed", name)
	}
	return true, nil
}

func (i *implementerCollection) CreateCollectionIfNotExists(ctx context.Context, name string, shardNum, replicasNum uint32, description string,
	indexes Indexes, params ...*CreateCollectionParams) (*Collection, error) {
	res, err := i.DescribeCollection(ctx, name)
	if err != nil {
		if strings.Contains(err.Error(), strconv.Itoa(ERR_UNDEFINED_COLLECTION)) {
			return i.CreateCollection(ctx, name, shardNum, replicasNum, description, indexes, params...)
		}
		return nil, fmt.Errorf("get collection %s failed, err: %v", name, err.Error())
	}
	if res == nil {
		return nil, fmt.Errorf("get collection %s failed", name)
	}
	return &res.Collection, nil
}

// CreateCollection create a collection. It returns collection struct if err is nil.
// The parameter `name` must be a unique string, otherwise an error will be returned.
// The parameter `shardNum`, `replicasNum` must bigger than 0, `description` could be empty.
// You can set the index field in Indexes, the vectorIndex must be set one currently, and
// the filterIndex sets at least one primaryKey value.
func (i *implementerCollection) CreateCollection(ctx context.Context, name string, shardNum, replicasNum uint32,
	description string, indexes Indexes, params ...*CreateCollectionParams) (*Collection, error) {
	if i.database.IsAIDatabase() {
		return nil, AIDbTypeError
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

	for _, v := range indexes.SparseVectorIndex {
		var column api.IndexColumn
		column.FieldName = v.FieldName
		column.FieldType = string(v.FieldType)
		column.IndexType = string(v.IndexType)
		column.MetricType = string(v.MetricType)

		req.Indexes = append(req.Indexes, &column)
	}

	for _, v := range indexes.FilterIndex {
		var column api.IndexColumn
		column.FieldName = v.FieldName
		column.FieldType = string(v.FieldType)
		if v.FieldType == Array {
			column.FieldElementType = string(v.ElemType)
		}
		column.IndexType = string(v.IndexType)
		req.Indexes = append(req.Indexes, &column)
	}
	if len(params) != 0 && params[0] != nil {
		param := params[0]
		if param.Embedding != nil {
			req.Embedding.Field = param.Embedding.Field
			req.Embedding.VectorField = param.Embedding.VectorField
			req.Embedding.Model = string(param.Embedding.Model)
			if param.Embedding.ModelName != "" {
				req.Embedding.Model = param.Embedding.ModelName
			}
		}
		if param.TtlConfig != nil {
			req.TtlConfig = new(collection.TtlConfig)
			req.TtlConfig.Enable = param.TtlConfig.Enable
			req.TtlConfig.TimeField = param.TtlConfig.TimeField
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

type ListCollectionResult struct {
	Collections []*Collection
}

// ListCollection get collection list.
// It return the list of collection, each collection same as DescribeCollection return.
func (i *implementerCollection) ListCollection(ctx context.Context) (*ListCollectionResult, error) {
	if i.database.IsAIDatabase() {
		return nil, AIDbTypeError
	}
	req := new(collection.ListReq)
	req.Database = i.database.DatabaseName
	res := new(collection.ListRes)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return nil, err
	}
	var collections []*Collection
	for _, collection := range res.Collections {
		collections = append(collections, i.toCollection(collection))
	}
	result := new(ListCollectionResult)
	result.Collections = collections
	return result, nil
}

type DescribeCollectionResult struct {
	Collection
}

// DescribeCollection get a collection detail.
// It returns the collection object to get collecton parameters or operate document api
func (i *implementerCollection) DescribeCollection(ctx context.Context, name string) (*DescribeCollectionResult, error) {
	if i.database.IsAIDatabase() {
		return nil, AIDbTypeError
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
	result := new(DescribeCollectionResult)
	result.Collection = *coll
	return result, nil
}

type DropCollectionResult struct {
	AffectedCount int
}

// DropCollection drop a collection. If collection not exist, it return nil.
func (i *implementerCollection) DropCollection(ctx context.Context, name string) (result *DropCollectionResult, err error) {
	if i.database.IsAIDatabase() {
		return nil, AIDbTypeError
	}
	req := new(collection.DropReq)
	req.Database = i.database.DatabaseName
	req.Collection = name

	res := new(collection.DropRes)
	err = i.Request(ctx, req, res)
	result = new(DropCollectionResult)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			return result, nil
		}
		return
	}
	result.AffectedCount = res.AffectedCount
	return
}

type TruncateCollectionResult struct {
	AffectedCount int
}

func (i *implementerCollection) TruncateCollection(ctx context.Context, name string) (result *TruncateCollectionResult, err error) {
	if i.database.IsAIDatabase() {
		return nil, AIDbTypeError
	}
	req := new(collection.TruncateReq)
	req.Database = i.database.DatabaseName
	req.Collection = name

	res := new(collection.TruncateRes)
	err = i.Request(ctx, req, res)

	if err != nil {
		return
	}
	result = new(TruncateCollectionResult)
	result.AffectedCount = res.AffectedCount
	return
}

// Collection get a collection interface to operate the document api. It could not send http request to vectordb.
// If you want to show collection parameters, use DescribeCollection.
func (i *implementerCollection) Collection(name string) *Collection {
	coll := new(Collection)
	coll.DatabaseName = i.database.DatabaseName
	coll.CollectionName = name

	flatImpl := new(implementerFlatDocument)
	flatImpl.SdkClient = i.SdkClient

	flatIndexImpl := new(implementerFlatIndex)
	flatIndexImpl.SdkClient = i.SdkClient

	docImpl := new(implementerDocument)
	docImpl.SdkClient = i.SdkClient
	docImpl.database = i.database
	docImpl.collection = coll
	docImpl.flat = flatImpl

	indexImpl := new(implementerIndex)
	indexImpl.SdkClient = i.SdkClient
	indexImpl.database = i.database
	indexImpl.collection = coll
	indexImpl.flat = flatIndexImpl

	coll.DocumentInterface = docImpl
	coll.IndexInterface = indexImpl

	return coll
}

func (i *implementerCollection) toCollection(collectionItem *collection.DescribeCollectionItem) *Collection {
	coll := new(Collection)
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
		coll.Embedding.Model = EmbeddingModel(collectionItem.Embedding.Model)
		coll.Embedding.ModelName = collectionItem.Embedding.Model
		coll.Embedding.Enabled = collectionItem.Embedding.Status == "enabled"
	}
	if collectionItem.TtlConfig != nil {
		coll.TtlConfig = new(TtlConfig)
		coll.TtlConfig.Enable = collectionItem.TtlConfig.Enable
		coll.TtlConfig.TimeField = collectionItem.TtlConfig.TimeField
	}

	if collectionItem.IndexStatus != nil {
		coll.IndexStatus = IndexStatus{
			Status: collectionItem.IndexStatus.Status,
		}
		coll.IndexStatus.StartTime, _ = time.Parse("2006-01-02 15:04:05", collectionItem.IndexStatus.StartTime)
	}
	coll.CreateTime, _ = time.Parse("2006-01-02 15:04:05", collectionItem.CreateTime)

	for _, index := range collectionItem.Indexes {
		if index == nil {
			continue
		}
		switch index.FieldType {
		case string(Vector):
			vector := VectorIndex{}
			vector.FieldName = index.FieldName
			vector.FieldType = FieldType(index.FieldType)
			vector.IndexType = IndexType(index.IndexType)
			vector.Dimension = index.Dimension
			vector.MetricType = MetricType(index.MetricType)
			vector.IndexedCount = index.IndexedCount

			if index.Params != nil {
				switch vector.IndexType {
				case HNSW:
					vector.Params = &HNSWParam{M: index.Params.M, EfConstruction: index.Params.EfConstruction}
				case IVF_FLAT:
					vector.Params = &IVFFLATParams{NList: index.Params.Nlist}
				case IVF_PQ:
					vector.Params = &IVFPQParams{M: index.Params.M, NList: index.Params.Nlist}
				case IVF_SQ4, IVF_SQ8, IVF_SQ16:
					vector.Params = &IVFSQParams{NList: index.Params.Nlist}
				}
			}
			coll.Indexes.VectorIndex = append(coll.Indexes.VectorIndex, vector)

		case string(SparseVector):
			vector := SparseVectorIndex{}
			vector.FieldName = index.FieldName
			vector.FieldType = FieldType(index.FieldType)
			vector.IndexType = IndexType(index.IndexType)
			vector.MetricType = MetricType(index.MetricType)
			coll.Indexes.SparseVectorIndex = append(coll.Indexes.SparseVectorIndex, vector)

		case string(Array):
			filter := FilterIndex{}
			filter.FieldName = index.FieldName
			filter.FieldType = FieldType(index.FieldType)
			filter.IndexType = IndexType(index.IndexType)
			filter.ElemType = FieldType(index.FieldElementType)
			coll.Indexes.FilterIndex = append(coll.Indexes.FilterIndex, filter)

		default:
			filter := FilterIndex{}
			filter.FieldName = index.FieldName
			filter.FieldType = FieldType(index.FieldType)
			filter.IndexType = IndexType(index.IndexType)
			coll.Indexes.FilterIndex = append(coll.Indexes.FilterIndex, filter)
		}
	}

	flatImpl := new(implementerFlatDocument)
	flatImpl.SdkClient = i.SdkClient

	flatIdexImpl := new(implementerFlatIndex)
	flatIdexImpl.SdkClient = i.SdkClient

	docImpl := new(implementerDocument)
	docImpl.SdkClient = i.SdkClient
	docImpl.database = i.database
	docImpl.collection = coll
	docImpl.flat = flatImpl
	coll.DocumentInterface = docImpl

	indexImpl := new(implementerIndex)
	indexImpl.SdkClient = i.SdkClient
	indexImpl.database = i.database
	indexImpl.collection = coll
	indexImpl.flat = flatIdexImpl
	coll.IndexInterface = indexImpl
	return coll
}

// optionParams param index parameters
func optionParams(column *api.IndexColumn, v VectorIndex) {
	column.Params = new(api.IndexParams)
	switch v.IndexType {
	case HNSW:
		if param, ok := v.Params.(*HNSWParam); ok && param != nil {
			column.Params.M = param.M
			column.Params.EfConstruction = param.EfConstruction
		}
	case IVF_FLAT:
		if param, ok := v.Params.(*IVFFLATParams); ok && param != nil {
			column.Params.Nlist = param.NList
		}
	case IVF_SQ4, IVF_SQ8, IVF_SQ16:
		if param, ok := v.Params.(*IVFSQParams); ok && param != nil {
			column.Params.Nlist = param.NList
		}
	case IVF_PQ:
		if param, ok := v.Params.(*IVFPQParams); ok && param != nil {
			column.Params.M = param.M
			column.Params.Nlist = param.NList
		}
	}
}

// Collection wrap the collection parameters and document interface to operating the document api
type Collection struct {
	DocumentInterface `json:"-"`
	IndexInterface    `json:"-"`
	DatabaseName      string      `json:"databaseName"`
	CollectionName    string      `json:"collectionName"`
	DocumentCount     int64       `json:"documentCount"`
	Alias             []string    `json:"alias"`
	ShardNum          uint32      `json:"shardNum"`
	ReplicasNum       uint32      `json:"replicasNum"`
	Indexes           Indexes     `json:"indexes"`
	IndexStatus       IndexStatus `json:"indexStatus"`
	Embedding         Embedding   `json:"embedding"`
	Description       string      `json:"description"`
	Size              uint64      `json:"size"`
	CreateTime        time.Time   `json:"createTime"`
	TtlConfig         *TtlConfig  `json:"ttlConfig,omitempty"`
}

func (c *Collection) Debug(v bool) {
	c.DocumentInterface.Debug(v)
}

func (c *Collection) WithTimeout(t time.Duration) {
	c.DocumentInterface.WithTimeout(t)
}

type Embedding struct {
	Field       string         `json:"field,omitempty"`
	VectorField string         `json:"vectorField,omitempty"`
	Model       EmbeddingModel `json:"model,omitempty"`
	ModelName   string         `json:"modelName,omitempty"` // 如果设置了ModelName，则使用ModelName；如果没设置，则使用Model
	Enabled     bool           `json:"enabled,omitempty"`   // 返回数据
}

type IndexStatus struct {
	Status    string
	StartTime time.Time
}

type TtlConfig struct {
	Enable    bool   `json:"enable"`
	TimeField string `json:"timeField,omitempty"`
}
