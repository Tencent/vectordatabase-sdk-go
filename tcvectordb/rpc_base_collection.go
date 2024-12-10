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

// [ExistsCollection] checks the existence of a specific collection.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - name: The name of the [collection] to check.
//
// Notes: It returns true if the collection exists.
//
// Returns a boolean variable indicating whether the collection exists or an error.
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

// [CreateCollection] creates a collection if it doesn't exist.
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

// [ListCollection] retrieves the list of all collections in the database.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//
// Notes: The database name is from the field of [implementerCollection].
//
// Returns a pointer to a [ListCollectionResult] object or an error.
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

// [DescribeCollection] retrieves information about a specific [Collection]. See [Collection] for more information.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - name: The name of the collection.
//
// Notes: The database name is from the field of [rpcImplementerCollection].
//
// Returns a pointer to a [DescribeCollectionResult] object or an error.
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

// [TruncateCollection] clears all the data and indexes in the Collection.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - name: The name of the collection to truncate.
//
// Returns a pointer to a [TruncateCollectionResult] object or an error.
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

// [Collection] returns a pointer to a [Collection] object. which includes the collection parameters
// and some interfaces to operate the document/index api.
//
// Parameters:
//   - name: The name of the collection.
//
// Notes:  If you want to show collection parameters, use DescribeCollection.
//
// Returns a pointer to a [Collection] object.
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
