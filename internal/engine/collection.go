package engine

import (
	"context"
	"strings"
	"time"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/internal/engine/api/collection"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/internal/proto"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/model"
)

type implementerCollection struct {
	model.SdkClient
	databaseName string
}

// CreateCollection create a collection. It returns collection struct if err is nil.
// The parameter `name` must be a unique string, otherwise an error will be returned.
// The parameter `shardNum`, `replicasNum` must bigger than 0, `description` could be empty.
// You can set the index field in model.Indexes, the vectorIndex must be set one currently, and
// the filterIndex sets at least one primaryKey value.
func (i *implementerCollection) CreateCollection(ctx context.Context, name string, shardNum, replicasNum uint32, description string, indexes model.Indexes, embedding *model.Embedding) (*model.Collection, error) {
	req := new(collection.CreateReq)
	req.Database = i.databaseName
	req.Collection = name
	req.ShardNum = shardNum
	req.ReplicaNum = replicasNum
	req.Description = description
	if embedding != nil {
		req.EmbeddingParams = &proto.EmbeddingParams{
			TextField:   embedding.TextField,
			VectorField: embedding.VectorField,
			Model:       embedding.Model,
		}
	}

	for _, v := range indexes.VectorIndex {
		var column proto.IndexColumn
		column.FieldName = v.FieldName
		column.FieldType = string(v.FieldType)
		column.IndexType = string(v.IndexType)
		column.MetricType = string(v.MetricType)
		column.Dimension = v.Dimension

		if v.IndexType == model.HNSW {
			column.Params = new(proto.IndexParams)
			column.Params.M = v.HNSWParam.M
			column.Params.EfConstruction = v.HNSWParam.EfConstruction
		}
		req.Indexes = append(req.Indexes, &column)
	}
	for _, v := range indexes.FilterIndex {
		var column proto.IndexColumn
		column.FieldName = v.FieldName
		column.FieldType = string(v.FieldType)
		column.IndexType = string(v.IndexType)
		req.Indexes = append(req.Indexes, &column)
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

// DescribeCollection get a collection detail. It returns the collection object to get collecton parameters or operate document api
func (i *implementerCollection) DescribeCollection(ctx context.Context, name string) (*model.Collection, error) {
	req := new(collection.DescribeReq)
	req.Database = i.databaseName
	req.Collection = name
	res := new(collection.DescribeRes)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return nil, err
	}
	coll := i.toCollection(res.Collection)

	return coll, nil
}

// DropCollection drop a collection. If collection not exist, it return nil.
func (i *implementerCollection) DropCollection(ctx context.Context, name string) (err error) {
	req := new(collection.DropReq)
	req.Database = i.databaseName
	req.Collection = name

	res := new(collection.DropRes)
	err = i.Request(ctx, req, res)

	if err != nil && strings.Contains(err.Error(), "not exist") {
		return nil
	}
	return err
}

func (i *implementerCollection) FlushCollection(ctx context.Context, name string) (affectedCount int, err error) {
	req := new(collection.FlushReq)
	req.Database = i.databaseName
	req.Collection = name

	res := new(collection.FlushRes)
	err = i.Request(ctx, req, res)

	if err != nil {
		return 0, err
	}
	return int(res.AffectedCount), nil
}

// ListCollection get collection list. It return the list of collection, each collection same as DescribeCollection return.
func (i *implementerCollection) ListCollection(ctx context.Context) ([]*model.Collection, error) {
	req := new(collection.ListReq)
	req.Database = i.databaseName
	res := new(collection.ListRes)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return nil, err
	}
	var collections []*model.Collection
	for _, collection := range res.Collections {
		collections = append(collections, i.toCollection(collection))
	}
	return collections, nil
}

func (i *implementerCollection) ModifyCollection(ctx context.Context, name, alias string) error {
	req := new(collection.ModifyReq)
	req.Database = i.databaseName
	res := new(collection.ModifyRes)
	err := i.Request(ctx, req, &res)
	return err
}

// Collection get a collection interface to operate the document api. It could not send http request to vectordb.
// If you want to show collection parameters, use DescribeCollection.
func (i *implementerCollection) Collection(name string) *model.Collection {
	coll := new(model.Collection)
	docImpl := new(implementerDocument)
	docImpl.SdkClient = i.SdkClient
	docImpl.databaseName = i.databaseName
	docImpl.collectionName = name
	coll.DocumentInterface = docImpl
	coll.DatabaseName = i.databaseName
	coll.CollectionName = name
	return coll
}

func (i *implementerCollection) toCollection(collectionItem *collection.DescribeCollectionItem) *model.Collection {
	coll := i.Collection(collectionItem.Collection)
	coll.DocumentCount = collectionItem.DocumentCount
	coll.Alias = collectionItem.AliasList
	coll.ShardNum = collectionItem.ShardNum
	coll.ReplicasNum = collectionItem.ReplicaNum
	coll.Description = collectionItem.Description
	coll.Size = collectionItem.Size
	if collectionItem.EmbeddingParams != nil {
		coll.Embedding = model.Embedding{
			TextField:   collectionItem.EmbeddingParams.TextField,
			VectorField: collectionItem.EmbeddingParams.VectorField,
			Model:       collectionItem.EmbeddingParams.Model,
			Enabled:     collectionItem.EmbeddingParams.Status == "enabled",
		}
	}
	if collectionItem.IndexStatus != nil {
		coll.IndexStatus = model.IndexStatus{
			Status: collectionItem.IndexStatus.Status,
		}
		coll.IndexStatus.StartTime, _ = time.Parse("2006-01-02 15:04:05", collectionItem.IndexStatus.StartTime)
	}
	coll.CreateTime, _ = time.Parse("2006-01-02 15:04:05", collectionItem.CreateTime)

	for _, index := range collectionItem.Indexes {
		switch index.FieldType {
		case string(model.Vector):
			vector := model.VectorIndex{}
			vector.FieldName = index.FieldName
			vector.FieldType = model.FieldType(index.FieldType)
			vector.IndexType = model.IndexType(index.IndexType)
			vector.Dimension = index.Dimension
			vector.MetricType = model.MetricType(index.MetricType)
			vector.IndexedCount = index.IndexedCount
			if vector.IndexType == model.HNSW {
				vector.HNSWParam.M = index.Params.M
				vector.HNSWParam.EfConstruction = index.Params.EfConstruction
			}
			coll.Indexes.VectorIndex = append(coll.Indexes.VectorIndex, vector)

		case string(model.FILTER):
			filter := model.FilterIndex{}
			filter.FieldName = index.FieldName
			filter.FieldType = model.FieldType(index.FieldType)
			filter.IndexType = model.IndexType(index.IndexType)

			coll.Indexes.FilterIndex = append(coll.Indexes.FilterIndex, filter)
		}
	}
	return coll
}
