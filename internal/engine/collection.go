package engine

import (
	"context"
	"strings"

	"git.woa.com/cloud_nosql/vectordb/vectordb-sdk-go/internal/engine/api/collection"
	"git.woa.com/cloud_nosql/vectordb/vectordb-sdk-go/internal/proto"
	"git.woa.com/cloud_nosql/vectordb/vectordb-sdk-go/model"
)

type implementerCollection struct {
	model.SdkClient
	databaseName string
}

func (i *implementerCollection) CreateCollection(ctx context.Context, name string, shardNum, replicasNum uint32, description string, indexes model.Indexes) (*model.Collection, error) {
	req := new(collection.CreateReq)
	req.Database = i.databaseName
	req.Collection = name
	req.ShardNum = shardNum
	req.ReplicaNum = replicasNum
	req.Description = description

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

func (i *implementerCollection) DropCollection(ctx context.Context, collectionName string) (err error) {
	req := new(collection.DropReq)
	req.Database = i.databaseName
	req.Collection = collectionName

	res := new(collection.DropRes)
	err = i.Request(ctx, req, res)

	if err != nil && strings.Contains(err.Error(), "not exist") {
		return nil
	}
	return err
}

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

func (i *implementerCollection) toCollection(collectionRes *proto.CreateCollectionRequest) *model.Collection {
	coll := i.Collection(collectionRes.Collection)
	coll.ShardNum = collectionRes.ShardNum
	coll.ReplicasNum = collectionRes.ReplicaNum
	coll.Description = collectionRes.Description
	coll.CreateTime = collectionRes.CreateTime
	coll.Size = collectionRes.Size

	for _, index := range collectionRes.Indexes {
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
