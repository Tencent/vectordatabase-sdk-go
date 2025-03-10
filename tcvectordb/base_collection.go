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
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/collection"
)

var _ CollectionInterface = &implementerCollection{}

// [CollectionInterface] provides apis of a collection.
type CollectionInterface interface {
	SdkClient

	// [ExistsCollection] checks the existence of a specific collection.
	ExistsCollection(ctx context.Context, name string) (bool, error)

	// [CreateCollectionIfNotExists] creates and initializes a new [Collection] if it doesn't exist.
	CreateCollectionIfNotExists(ctx context.Context, name string, shardNum, replicasNum uint32, description string,
		indexes Indexes, params ...*CreateCollectionParams) (*Collection, error)

	// [CreateCollection] creates and initializes a new [Collection].
	CreateCollection(ctx context.Context, name string, shardNum, replicasNum uint32, description string,
		indexes Indexes, params ...*CreateCollectionParams) (*Collection, error)

	// [ListCollection] retrieves the list of all collections in the database.
	ListCollection(ctx context.Context) (result *ListCollectionResult, err error)

	// [DescribeCollection] retrieves information about a specific [Collection]. See [Collection] for more information.
	DescribeCollection(ctx context.Context, name string) (result *DescribeCollectionResult, err error)

	// [DropCollection] drops a specific collection.
	DropCollection(ctx context.Context, name string) (result *DropCollectionResult, err error)

	// [TruncateCollection] clears all the data and indexes in the Collection.
	TruncateCollection(ctx context.Context, name string) (result *TruncateCollectionResult, err error)

	// [Collection] returns a pointer to a [Collection] object. which includes the collection parameters
	// and some interfaces to operate the document/index api.
	Collection(name string) *Collection
}

type implementerCollection struct {
	SdkClient
	database *Database
}

// [CreateCollectionParams] holds the parameters for creating a new collection.
//
// Fields:
//   - Embedding: An optional embedding for embedding text when upsert documents.
//   - TtlConfig: TTL configuration. When TtlConfig.Enable is set to True and TtlConfig.TimeField
//     is set to expire_at, it means that TTL (Time to Live) is enabled.
//     In this case, the document will be automatically removed after 60 minites when the time specified
//     in the expire_at field is reached.
type CreateCollectionParams struct {
	Embedding         *Embedding
	TtlConfig         *TtlConfig
	FilterIndexConfig *FilterIndexConfig
}

type CreateCollectionResult struct {
	Collection
}

// [ExistsCollection] checks the existence of a specific collection.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - name: The name of the collection to check.
//
// Notes: It returns true if the collection exists.
//
// Returns a boolean variable indicating whether the collection exists or an error.
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

// [CreateCollectionIfNotExists] creates a collection if it doesn't exist.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - name: The name of the collection. Collection name must be 1-128 characters long,
//     start with an alphanumeric character,
//     and consist only of alphanumeric characters, numbers, '_' or '-'.
//   - shardNum: The shard number of the collection, which must bigger than 0.
//   - replicasNum: The replicas number of the collection.
//   - description: (Optional) The description of the collection.
//   - index: A [Indexes] object that includes a list of the index properties for the documents in a collection. The vectorIndex
//     must be set one currently, and the filterIndex sets at least one primaryKey called "id".
//   - params: A pointer to a [CreateCollectionParams] object that includes the other parameters for the collection.
//     See [CreateCollectionParams] for more information.
//
// Returns a pointer to a [Collection] object or an error.
func (i *implementerCollection) CreateCollectionIfNotExists(ctx context.Context, name string, shardNum, replicasNum uint32,
	description string, indexes Indexes, params ...*CreateCollectionParams) (*Collection, error) {
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

// [CreateCollection] creates and initializes a new [Collection].
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - name: The name of the collection. Collection name must be 1-128 characters long,
//     start with an alphanumeric character,
//     and consist only of alphanumeric characters, numbers, '_' or '-'.
//   - shardNum: The shard number of the collection, which must bigger than 0.
//   - replicasNum: The replicas number of the collection.
//   - description: (Optional) The description of the collection.
//   - index: A [Indexes] object that includes a list of the index properties for the documents in a collection. The vectorIndex
//     must be set one currently, and the filterIndex sets at least one primaryKey called "id".
//   - params: A pointer to a [CreateCollectionParams] object that includes the other parameters for the collection.
//     See [CreateCollectionParams] for more information.
//
// Returns a pointer to a [Collection] object or an error.
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
		column.AutoId = v.AutoId
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
		if param.FilterIndexConfig != nil {
			req.FilterIndexConfig = new(collection.FilterIndexConfig)
			req.FilterIndexConfig.FilterAll = param.FilterIndexConfig.FilterAll
			req.FilterIndexConfig.FieldsWithoutIndex = param.FilterIndexConfig.FieldsWithoutIndex
			req.FilterIndexConfig.MaxStrLen = param.FilterIndexConfig.MaxStrLen
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

// [ListCollection] retrieves the list of all collections in the database.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//
// Notes: The database name is from the field of [implementerCollection].
//
// Returns a pointer to a [ListCollectionResult] object or an error.
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

// [DescribeCollection] retrieves information about a specific [Collection]. See [Collection] for more information.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - name: The name of the collection.
//
// Notes: The database name is from the field of [implementerCollection].
//
// Returns a pointer to a [DescribeCollectionResult] object or an error.
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

// [DropCollection] drops a specific collection.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - name: The name of the collection to drop.
//
// Notes: If the collection doesn't exist, it returns 0 for DropCollectionResult.AffectedCount.
//
// Returns a pointer to a [DropCollectionResult] object or an error.
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

// [TruncateCollection] clears all the data and indexes in the Collection.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - name: The name of the collection to truncate.
//
// Returns a pointer to a [TruncateCollectionResult] object or an error.
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

// [Collection] returns a pointer to a [Collection] object. which includes the collection parameters
// and some interfaces to operate the document/index api.
//
// Parameters:
//   - name: The name of the collection.
//
// Notes:  If you want to show collection parameters, use DescribeCollection.
//
// Returns a pointer to a [Collection] object.
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
	if collectionItem.FilterIndexConfig != nil {
		coll.FilterIndexConfig = new(FilterIndexConfig)
		coll.FilterIndexConfig.FilterAll = collectionItem.FilterIndexConfig.FilterAll
		coll.FilterIndexConfig.FieldsWithoutIndex = collectionItem.FilterIndexConfig.FieldsWithoutIndex
		coll.FilterIndexConfig.MaxStrLen = collectionItem.FilterIndexConfig.MaxStrLen
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
			filter.AutoId = index.AutoId
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

func optionParamsFromIndexParams(column *api.IndexColumn, v IndexParams) {
	column.Params = new(api.IndexParams)
	switch param := v.(type) {
	case *HNSWParam:
		if param != nil {
			column.Params.M = param.M
			column.Params.EfConstruction = param.EfConstruction
		}
	case *IVFFLATParams:
		if param != nil {
			column.Params.Nlist = param.NList
		}
	case *IVFSQParams:
		if param != nil {
			column.Params.Nlist = param.NList
		}
	case *IVFPQParams:
		if param != nil {
			column.Params.M = param.M
			column.Params.Nlist = param.NList
		}
	default:
		log.Printf("[Warning] unknown type: %v", reflect.TypeOf(v))
	}

}

// [Collection] holds the collection parameters and some interfaces to operate the document/index api.
//
// Fields:
//   - DatabaseName: The name of the database.
//   - CollectionName: The name of the collection.
//   - DocumentCount: The number of documents in the Collection.
//   - Alias: All aliases of the Collection.
//   - ShardNum: The shard number of the collection, which must bigger than 0.
//   - ReplicasNum: The replicas number of the collection.
//   - Indexes: A [Indexes] object that includes a list of the index properties for the documents in a collection.
//   - IndexStatus: The status of index. which has four states:
//     ready, which indicates that the current collection's index is ready to use.
//     training, which indicates that the current Collection is undergoing data training, like training the model to generate vector data.
//     building, which indicates that the current Collection is rebuilding the index, like storing the generated vector data into a new index.
//     failed, which indicates that the index rebuilding has failed, which may affect collection read and write operations.
//   - Embedding: An optional embedding for embedding text when upsert documents.
//   - Description: The description of the collection.
//   - CreateTime: The create time of collection.
//   - TtlConfig: TTL configuration. When TtlConfig.Enable is set to True and TtlConfig.TimeField
//     is set to expire_at, it means that TTL (Time to Live) is enabled.
//     In this case, the document will be automatically removed after 60 minites when the time specified
//     in the expire_at field is reached.
type Collection struct {
	DocumentInterface `json:"-"`
	IndexInterface    `json:"-"`
	DatabaseName      string             `json:"databaseName"`
	CollectionName    string             `json:"collectionName"`
	DocumentCount     int64              `json:"documentCount"`
	Alias             []string           `json:"alias"`
	ShardNum          uint32             `json:"shardNum"`
	ReplicasNum       uint32             `json:"replicasNum"`
	Indexes           Indexes            `json:"indexes"`
	IndexStatus       IndexStatus        `json:"indexStatus"`
	Embedding         Embedding          `json:"embedding"`
	Description       string             `json:"description"`
	Size              uint64             `json:"size"`
	CreateTime        time.Time          `json:"createTime"`
	TtlConfig         *TtlConfig         `json:"ttlConfig,omitempty"`
	FilterIndexConfig *FilterIndexConfig `json:"filterIndexConfig,omitempty"`
}

func (c *Collection) Debug(v bool) {
	c.DocumentInterface.Debug(v)
}

func (c *Collection) WithTimeout(t time.Duration) {
	c.DocumentInterface.WithTimeout(t)
}

type Embedding struct {
	Field       string `json:"field,omitempty"`
	VectorField string `json:"vectorField,omitempty"`
	// Deprecated: Use ModelName instead.
	Model     EmbeddingModel `json:"model,omitempty"`
	ModelName string         `json:"modelName,omitempty"`
	Enabled   bool           `json:"enabled,omitempty"` // 返回数据
}

type IndexStatus struct {
	Status    string
	StartTime time.Time
}

type TtlConfig struct {
	Enable    bool   `json:"enable"`
	TimeField string `json:"timeField,omitempty"`
}

type FilterIndexConfig struct {
	FilterAll          bool     `json:"filterAll"`
	FieldsWithoutIndex []string `json:"fieldsWithoutIndex,omitempty"`
	MaxStrLen          *uint32  `json:"maxStrLen,omitempty"`
}
