package tcvectordb

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/olama"
)

var _ CollectionInterface = &rpcImplementerCollection{}

type rpcImplementerCollection struct {
	SdkClient
	rpcClient olama.SearchEngineClient
	database  *Database
}

func (r *rpcImplementerCollection) ExistsCollection(ctx context.Context, name string) (bool, error) {
	res, err := r.DescribeCollection(ctx, name)
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

func (r *rpcImplementerCollection) CreateCollectionIfNotExists(ctx context.Context, name string, shardNum, replicasNum uint32, description string,
	indexes Indexes, params ...*CreateCollectionParams) (*Collection, error) {
	res, err := r.DescribeCollection(ctx, name)
	if err != nil {
		if strings.Contains(err.Error(), strconv.Itoa(ERR_UNDEFINED_COLLECTION)) {
			return r.CreateCollection(ctx, name, shardNum, replicasNum, description, indexes, params...)
		}
		return nil, fmt.Errorf("get collection %s failed, err: %v", name, err.Error())
	}
	if res == nil {
		return nil, fmt.Errorf("get collection %s failed", name)
	}
	return &res.Collection, nil
}

func (r *rpcImplementerCollection) CreateCollection(ctx context.Context, name string, shardNum, replicasNum uint32, description string, indexes Indexes, params ...*CreateCollectionParams) (*Collection, error) {
	if r.database.IsAIDatabase() {
		return nil, AIDbTypeError
	}
	req := &olama.CreateCollectionRequest{
		Database:    r.database.DatabaseName,
		Collection:  name,
		ShardNum:    shardNum,
		ReplicaNum:  replicasNum,
		Description: description,
		Indexes:     make(map[string]*olama.IndexColumn),
	}

	for _, v := range indexes.VectorIndex {
		column := &olama.IndexColumn{
			FieldName:  v.FieldName,
			FieldType:  string(v.FieldType),
			IndexType:  string(v.IndexType),
			MetricType: string(v.MetricType),
			Dimension:  v.Dimension,
		}
		optionRpcParams(column, v)
		req.Indexes[v.FieldName] = column
	}

	for _, v := range indexes.SparseVectorIndex {
		column := &olama.IndexColumn{
			FieldName:  v.FieldName,
			FieldType:  string(v.FieldType),
			IndexType:  string(v.IndexType),
			MetricType: string(v.MetricType),
		}
		req.Indexes[v.FieldName] = column
	}

	for _, v := range indexes.FilterIndex {
		column := &olama.IndexColumn{
			FieldName: v.FieldName,
			FieldType: string(v.FieldType),
			IndexType: string(v.IndexType),
		}
		if v.FieldType == Array {
			column.FieldElementType = string(String)
		}
		req.Indexes[v.FieldName] = column
	}
	if len(params) != 0 && params[0] != nil {
		param := params[0]
		if param.Embedding != nil {
			req.EmbeddingParams = &olama.EmbeddingParams{
				Field:       param.Embedding.Field,
				VectorField: param.Embedding.VectorField,
				ModelName:   string(param.Embedding.Model),
			}
			if param.Embedding.ModelName != "" {
				req.EmbeddingParams.ModelName = param.Embedding.ModelName
			}
		}
		if param.TtlConfig != nil {
			req.TtlConfig = &olama.TTLConfig{
				Enable:    param.TtlConfig.Enable,
				TimeField: param.TtlConfig.TimeField,
			}
		}
	}

	_, err := r.rpcClient.CreateCollection(ctx, req)
	if err != nil {
		return nil, err
	}

	coll := r.Collection(req.Collection)
	coll.ShardNum = req.ShardNum
	coll.ReplicasNum = req.ReplicaNum
	coll.Description = req.Description
	coll.Indexes = indexes
	return coll, nil
}

func (r *rpcImplementerCollection) ListCollection(ctx context.Context) (*ListCollectionResult, error) {
	if r.database.IsAIDatabase() {
		return nil, AIDbTypeError
	}
	req := &olama.ListCollectionsRequest{
		Database: r.database.DatabaseName,
	}
	res, err := r.rpcClient.ListCollections(ctx, req)
	if err != nil {
		return nil, err
	}
	var collections []*Collection
	for _, collection := range res.Collections {
		collections = append(collections, r.toCollection(collection))
	}
	result := &ListCollectionResult{
		Collections: collections,
	}
	return result, nil
}

func (r *rpcImplementerCollection) DescribeCollection(ctx context.Context, name string) (*DescribeCollectionResult, error) {
	if r.database.IsAIDatabase() {
		return nil, AIDbTypeError
	}
	req := &olama.DescribeCollectionRequest{
		Database:   r.database.DatabaseName,
		Collection: name,
	}
	res, err := r.rpcClient.DescribeCollection(ctx, req)
	if err != nil {
		return nil, err
	}
	if res.Collection == nil {
		return nil, fmt.Errorf("get collection %s failed", name)
	}
	coll := r.toCollection(res.Collection)
	result := &DescribeCollectionResult{
		Collection: *coll,
	}
	return result, nil
}

func (r *rpcImplementerCollection) DropCollection(ctx context.Context, name string) (*DropCollectionResult, error) {
	if r.database.IsAIDatabase() {
		return nil, AIDbTypeError
	}
	req := &olama.DropCollectionRequest{
		Database:   r.database.DatabaseName,
		Collection: name,
	}
	res, err := r.rpcClient.DropCollection(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			return &DropCollectionResult{}, nil
		}
		return nil, err
	}
	return &DropCollectionResult{AffectedCount: int(res.AffectedCount)}, nil
}

func (r *rpcImplementerCollection) TruncateCollection(ctx context.Context, name string) (*TruncateCollectionResult, error) {
	if r.database.IsAIDatabase() {
		return nil, AIDbTypeError
	}
	req := &olama.TruncateCollectionRequest{
		Database:   r.database.DatabaseName,
		Collection: name,
	}
	res, err := r.rpcClient.TruncateCollection(ctx, req)
	if err != nil {
		return nil, err
	}
	return &TruncateCollectionResult{AffectedCount: int(res.AffectedCount)}, nil
}

func (r *rpcImplementerCollection) Collection(name string) *Collection {
	coll := &Collection{
		DatabaseName:   r.database.DatabaseName,
		CollectionName: name,
	}
	flatImpl := &rpcImplementerFlatDocument{
		SdkClient: r.SdkClient,
		rpcClient: r.rpcClient,
	}
	flatIndexImpl := &rpcImplementerFlatIndex{
		SdkClient: r.SdkClient,
		rpcClient: r.rpcClient,
	}
	docImpl := &rpcImplementerDocument{
		SdkClient:  r.SdkClient,
		flat:       flatImpl,
		rpcClient:  r.rpcClient,
		database:   r.database,
		collection: coll,
	}
	indexImpl := &rpcImplementerIndex{
		SdkClient:  r.SdkClient,
		rpcClient:  r.rpcClient,
		flat:       flatIndexImpl,
		database:   r.database,
		collection: coll,
	}
	coll.DocumentInterface = docImpl
	coll.IndexInterface = indexImpl
	return coll
}

func (r *rpcImplementerCollection) toCollection(collectionItem *olama.CreateCollectionRequest) *Collection {
	coll := &Collection{
		DatabaseName:   r.database.DatabaseName,
		CollectionName: collectionItem.Collection,
		DocumentCount:  int64(collectionItem.Size),
		Alias:          collectionItem.AliasList,
		ShardNum:       collectionItem.ShardNum,
		ReplicasNum:    collectionItem.ReplicaNum,
		Description:    collectionItem.Description,
		Size:           collectionItem.Size,
	}
	if collectionItem.EmbeddingParams != nil {
		coll.Embedding.Field = collectionItem.EmbeddingParams.Field
		coll.Embedding.VectorField = collectionItem.EmbeddingParams.VectorField
		coll.Embedding.Model = EmbeddingModel(collectionItem.EmbeddingParams.ModelName)
		coll.Embedding.ModelName = collectionItem.EmbeddingParams.ModelName
		coll.Embedding.Enabled = false
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
			vector.IndexedCount = collectionItem.Size

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
	flatImpl := &rpcImplementerFlatDocument{
		SdkClient: r.SdkClient,
		rpcClient: r.rpcClient,
	}
	flatIndexImpl := &rpcImplementerFlatIndex{
		SdkClient: r.SdkClient,
		rpcClient: r.rpcClient,
	}
	docImpl := &rpcImplementerDocument{
		SdkClient:  r.SdkClient,
		flat:       flatImpl,
		rpcClient:  r.rpcClient,
		database:   r.database,
		collection: coll,
	}
	coll.DocumentInterface = docImpl
	indexImpl := &rpcImplementerIndex{
		r.SdkClient,
		r.rpcClient,
		flatIndexImpl,
		r.database,
		coll,
	}
	coll.IndexInterface = indexImpl
	return coll
}

func optionRpcParams(column *olama.IndexColumn, v VectorIndex) {
	column.Params = new(olama.IndexParams)
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
