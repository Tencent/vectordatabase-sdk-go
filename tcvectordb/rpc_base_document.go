package tcvectordb

import (
	"context"
	"fmt"

	"github.com/tencent/vectordatabase-sdk-go/tcvdbtext/encoder"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/olama"
)

var _ DocumentInterface = &rpcImplementerDocument{}
var _ FlatInterface = &rpcImplementerFlatDocument{}

type rpcImplementerDocument struct {
	SdkClient
	flat       FlatInterface
	rpcClient  olama.SearchEngineClient
	database   *Database
	collection *Collection
}

// [Upsert] upserts documents into a collection.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - documents: The list of the [Document] object or  map[string]interface{} to upsert. Maximum 1000.
//   - params: A pointer to a [UpsertDocumentParams] object that includes the other parameters for upserting documents' operation.
//     See [UpsertDocumentParams] for more information.
//
// Notes: The name of the database and the name of collection are from the fields of [rpcImplementerDocument].
//
// Returns a pointer to a [UpsertDocumentResult] object or an error.
func (r *rpcImplementerDocument) Upsert(ctx context.Context, documents interface{}, params ...*UpsertDocumentParams) (*UpsertDocumentResult, error) {
	return r.flat.Upsert(ctx, r.database.DatabaseName, r.collection.CollectionName, documents, params...)
}

// [Query] queries documents that satisfies the condition from the collection.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - documentIds: The list of the documents' ids, which are used for filtering documents.
//   - params: A pointer to a [QueryDocumentParams] object that includes the other parameters for querying documents' operation.
//     See [QueryDocumentParams] for more information.
//
// Notes: The name of the database and the name of collection are from the fields of [rpcImplementerDocument].
//
// Returns a pointer to a [QueryDocumentResult] object or an error.
func (r *rpcImplementerDocument) Query(ctx context.Context, documentIds []string, params ...*QueryDocumentParams) (*QueryDocumentResult, error) {
	return r.flat.Query(ctx, r.database.DatabaseName, r.collection.CollectionName, documentIds, params...)
}

// [Search] returns the most similar topK vectors by the given vectors.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - databaseName: The name of the database.
//   - collectionName: The name of the collection.
//   - vectors: The list of vectors to search. The maximum number of elements in the array is 20.
//   - params: A pointer to a [SearchDocumentParams] object that includes the other parameters for searching documents' operation.
//     See [SearchDocumentParams] for more information.
//
// Notes: The name of the database and the name of collection are from the fields of [rpcImplementerDocument].
//
// Returns a pointer to a [SearchDocumentResult] object or an error.
func (r *rpcImplementerDocument) Search(ctx context.Context, vectors [][]float32, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return r.flat.Search(ctx, r.database.DatabaseName, r.collection.CollectionName, vectors, params...)
}

// [SearchById] returns the most similar topK vectors by the given documentIds.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - documentIds: The list of the documents' ids, which are used for filtering documents.
//   - params: A pointer to a [SearchDocumentParams] object that includes the other parameters for searching documents' operation.
//     See [SearchDocumentParams] for more information.
//
// Notes: The name of the database and the name of collection are from the fields of [rpcImplementerDocument].
//
// Returns a pointer to a [SearchDocumentResult] object or an error.
func (r *rpcImplementerDocument) SearchById(ctx context.Context, documentIds []string, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return r.flat.SearchById(ctx, r.database.DatabaseName, r.collection.CollectionName, documentIds, params...)
}

// [SearchByText] returns the most similar topK vectors by the given text map.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - text: It is a map where the keys represent column names, and the values are lists of column values to be retrieved,
//     which are used to retrieve data similar to ones.
//   - params: A pointer to a [SearchDocumentParams] object that includes the other parameters for searching documents' operation.
//     See [SearchDocumentParams] for more information.
//
// Notes: The name of the database and the name of collection are from the fields of [rpcImplementerDocument].
//
// Returns a pointer to a [SearchDocumentResult] object or an error.
func (r *rpcImplementerDocument) SearchByText(ctx context.Context, text map[string][]string, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return r.flat.SearchByText(ctx, r.database.DatabaseName, r.collection.CollectionName, text, params...)
}

// [HybridSearch] retrieves both dense and sparse vectors to return the most similar topK vectors.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - params: A [HybridSearchDocumentParams] object that includes the other parameters for hybrid searching documents' operation.
//     See [HybridSearchDocumentParams] for more information.
//
// Notes: The name of the database and the name of collection are from the fields of [rpcImplementerDocument].
//
// Returns a pointer to a [SearchDocumentResult] object or an error.
func (r *rpcImplementerDocument) HybridSearch(ctx context.Context, params HybridSearchDocumentParams) (*SearchDocumentResult, error) {
	return r.flat.HybridSearch(ctx, r.database.DatabaseName, r.collection.CollectionName, params)
}

// [Delete] deletes documents by conditions.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - params: A [DeleteDocumentParams] object that includes the other parameters for deleting documents' operation.
//     See [DeleteDocumentParams] for more information.
//
// Notes: The name of the database and the name of collection are from the fields of [rpcImplementerDocument].
//
// Returns a pointer to a [DeleteDocumentResult] object or an error.
func (r *rpcImplementerDocument) Delete(ctx context.Context, param DeleteDocumentParams) (*DeleteDocumentResult, error) {
	return r.flat.Delete(ctx, r.database.DatabaseName, r.collection.CollectionName, param)
}

// [Update] updates documents by conditions.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - param: A [UpdateDocumentParams] object that includes the other parameters for updating documents' operation.
//     See [UpdateDocumentParams] for more information.
//
// Notes: The name of the database and the name of collection are from the fields of [rpcImplementerDocument].
//
// Returns a pointer to a [UpdateDocumentResult] object or an error.
func (r *rpcImplementerDocument) Update(ctx context.Context, param UpdateDocumentParams) (*UpdateDocumentResult, error) {
	return r.flat.Update(ctx, r.database.DatabaseName, r.collection.CollectionName, param)
}

// [Count] counts the number of documents in a collection that satisfy the specified filter conditions.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - databaseName: The name of the database.
//   - collectionName: The name of the collection.
//   - param: A [CountDocumentParams] object that includes the other parameters for counting documents' operation.
//     See [CountDocumentParams] for more information.
//
// Returns a pointer to a [CountDocumentResult] object or an error.
func (r *rpcImplementerDocument) Count(ctx context.Context, params ...CountDocumentParams) (*CountDocumentResult, error) {
	return r.flat.Count(ctx, r.database.DatabaseName, r.collection.CollectionName, params...)
}

type rpcImplementerFlatDocument struct {
	SdkClient
	rpcClient olama.SearchEngineClient
}

// [Upsert] upserts documents into a collection.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - databaseName: The name of the database.
//   - collectionName: The name of the collection.
//   - documents: The list of the [Document] object or  map[string]interface{} to upsert. Maximum 1000.
//   - params: A pointer to a [UpsertDocumentParams] object that includes the other parameters for upserting documents' operation.
//     See [UpsertDocumentParams] for more information.
//
// Returns a pointer to a [UpsertDocumentResult] object or an error.
func (r *rpcImplementerFlatDocument) Upsert(ctx context.Context, databaseName, collectionName string,
	documents interface{}, params ...*UpsertDocumentParams) (*UpsertDocumentResult, error) {
	req := &olama.UpsertRequest{
		Database:   databaseName,
		Collection: collectionName,
	}

	if docs, ok := documents.([]Document); ok {
		for _, doc := range docs {
			d := &olama.Document{
				Id:     doc.Id,
				Vector: doc.Vector,
				Fields: make(map[string]*olama.Field),
			}

			d.SparseVector = make([]*olama.SparseVecItem, 0)
			for _, sv := range doc.SparseVector {
				d.SparseVector = append(d.SparseVector, &olama.SparseVecItem{
					TermId: sv.TermId,
					Score:  sv.Score,
				})
			}

			for k, v := range doc.Fields {
				d.Fields[k] = ConvertField2Grpc(&v)
			}
			req.Documents = append(req.Documents, d)
		}
	} else if docs, ok := documents.([]map[string]interface{}); ok {
		for _, doc := range docs {
			var sId string
			var aVector []float32
			if id, ok := doc["id"]; ok {
				if sId, ok = id.(string); ok {
					delete(doc, "id")
				} else {
					return nil, fmt.Errorf("upsert failed, because of incorrect id field type, which must be string")
				}
			}
			if vector, ok := doc["vector"]; ok {
				if aVector, ok = vector.([]float32); ok {
					delete(doc, "vector")
				} else {
					return nil, fmt.Errorf("upsert failed, because of incorrect vector field type, which must be []float32")
				}
			}

			d := &olama.Document{
				Id:     sId,
				Vector: aVector,
				Fields: make(map[string]*olama.Field),
			}

			if sparseVector, ok := doc["sparse_vector"]; ok {
				if aSparseVector, ok := sparseVector.([][]interface{}); ok {
					d.SparseVector = make([]*olama.SparseVecItem, 0)
					for _, sv := range aSparseVector {
						svItem, err := ConvSliceInterface2SparseVecItem(sv)
						if err != nil {
							return nil, fmt.Errorf("upsert failed. doc's sparse_vector data is incorrect. doc id is %v. err: %v", d.Id, err.Error())
						}
						d.SparseVector = append(d.SparseVector, &olama.SparseVecItem{
							TermId: svItem.TermId,
							Score:  svItem.Score})
					}
					delete(doc, "sparse_vector")
				} else {
					return nil, fmt.Errorf("upsert failed, because of incorrect sparse_vector field type, which must be [][]interface{}")
				}
			}

			for k, v := range doc {
				d.Fields[k] = ConvertField2Grpc(&Field{Val: v})
			}
			req.Documents = append(req.Documents, d)
		}
	} else {
		return nil, fmt.Errorf("upsert failed, because of incorrect documents type, which must be []Document or []map[string]interface{}")
	}

	if len(params) != 0 && params[0] != nil {
		param := params[0]
		if param.BuildIndex != nil {
			req.BuildIndex = *param.BuildIndex
		} else {
			req.BuildIndex = true
		}
	} else {
		req.BuildIndex = true
	}

	res, err := r.rpcClient.Upsert(ctx, req)
	if err != nil {
		return nil, err
	}
	return &UpsertDocumentResult{AffectedCount: int(res.AffectedCount)}, nil
}

// [Query] queries documents that satisfies the condition from the collection.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - databaseName: The name of the database.
//   - collectionName: The name of the collection.
//   - documentIds: The list of the documents' ids, which are used for filtering documents.
//   - params: A pointer to a [QueryDocumentParams] object that includes the other parameters for querying documents' operation.
//     See [QueryDocumentParams] for more information.
//
// Returns a pointer to a [QueryDocumentResult] object or an error.
func (r *rpcImplementerFlatDocument) Query(ctx context.Context, databaseName, collectionName string,
	documentIds []string, params ...*QueryDocumentParams) (*QueryDocumentResult, error) {
	req := &olama.QueryRequest{
		Database:   databaseName,
		Collection: collectionName,
		Query: &olama.QueryCond{
			DocumentIds: documentIds,
		},
		ReadConsistency: string(r.SdkClient.Options().ReadConsistency),
	}
	if len(params) != 0 && params[0] != nil {
		param := params[0]
		req.Query.Filter = param.Filter.Cond()
		req.Query.RetrieveVector = param.RetrieveVector
		req.Query.OutputFields = param.OutputFields
		req.Query.Offset = param.Offset
		req.Query.Limit = param.Limit
	}
	res, err := r.rpcClient.Query(ctx, req)
	if err != nil {
		return nil, err
	}
	result := &QueryDocumentResult{}
	var documents []Document
	for _, doc := range res.Documents {
		var d Document
		d.Id = doc.Id
		d.Vector = doc.Vector
		d.SparseVector = make([]encoder.SparseVecItem, 0)
		for _, sv := range doc.SparseVector {
			d.SparseVector = append(d.SparseVector, encoder.SparseVecItem{
				TermId: sv.TermId,
				Score:  sv.Score,
			})
		}
		d.Fields = make(map[string]Field)

		for n, v := range doc.Fields {
			d.Fields[n] = *ConvertGrpc2Field(v)
		}
		documents = append(documents, d)
	}
	result.Documents = documents
	result.AffectedCount = len(documents)
	result.Total = res.Count
	return result, nil
}

// [Search] returns the most similar topK vectors by the given vectors.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - databaseName: The name of the database.
//   - collectionName: The name of the collection.
//   - vectors: The list of vectors to search. The maximum number of elements in the array is 20.
//   - params: A pointer to a [SearchDocumentParams] object that includes the other parameters for searching documents' operation.
//     See [SearchDocumentParams] for more information.
//
// Returns a pointer to a [SearchDocumentResult] object or an error.
func (r *rpcImplementerFlatDocument) Search(ctx context.Context, databaseName, collectionName string,
	vectors [][]float32, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return r.search(ctx, databaseName, collectionName, nil, vectors, nil, params...)
}

// [SearchById] returns the most similar topK vectors by the given documentIds.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - databaseName: The name of the database.
//   - collectionName: The name of the collection.
//   - documentIds: The list of the documents' ids, which are used for filtering documents.
//   - params: A pointer to a [SearchDocumentParams] object that includes the other parameters for searching documents' operation.
//     See [SearchDocumentParams] for more information.
//
// Returns a pointer to a [SearchDocumentResult] object or an error.
func (r *rpcImplementerFlatDocument) SearchById(ctx context.Context, databaseName, collectionName string,
	documentIds []string, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return r.search(ctx, databaseName, collectionName, documentIds, nil, nil, params...)
}

// [SearchByText] returns the most similar topK vectors by the given text map.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - databaseName: The name of the database.
//   - collectionName: The name of the collection.
//   - text: It is a map where the keys represent column names, and the values are lists of column values to be retrieved,
//     which are used to retrieve data similar to ones.
//   - params: A pointer to a [SearchDocumentParams] object that includes the other parameters for searching documents' operation.
//     See [SearchDocumentParams] for more information.
//
// Returns a pointer to a [SearchDocumentResult] object or an error.
func (r *rpcImplementerFlatDocument) SearchByText(ctx context.Context, databaseName, collectionName string,
	text map[string][]string, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return r.search(ctx, databaseName, collectionName, nil, nil, text, params...)
}

// [HybridSearch] retrieves both dense and sparse vectors to return the most similar topK vectors.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - databaseName: The name of the database.
//   - collectionName: The name of the collection.
//   - params: A [HybridSearchDocumentParams] object that includes the other parameters for hybrid searching documents' operation.
//     See [HybridSearchDocumentParams] for more information.
//
// Returns a pointer to a [SearchDocumentResult] object or an error.
func (r *rpcImplementerFlatDocument) HybridSearch(ctx context.Context, databaseName, collectionName string,
	params HybridSearchDocumentParams) (*SearchDocumentResult, error) {
	req := &olama.SearchRequest{
		Database:        databaseName,
		Collection:      collectionName,
		ReadConsistency: string(r.SdkClient.Options().ReadConsistency),
		Search:          &olama.SearchCond{},
	}

	req.Search.Ann = make([]*olama.AnnData, 0)
	req.Search.Sparse = make([]*olama.SparseData, 0)

	for i, annParam := range params.AnnParams {
		fieldName := "vector"
		if annParam.FieldName != "" {
			fieldName = annParam.FieldName
		}
		req.Search.Ann = append(req.Search.Ann, &olama.AnnData{
			FieldName: fieldName,
		})
		if annParam.Limit != nil {
			req.Search.Ann[i].Limit = uint32(*annParam.Limit)
		}

		if vec, ok := annParam.Data.([]float32); ok {
			req.Search.Ann[i].Data = make([]*olama.VectorArray, 0, len(req.Search.Vectors))
			req.Search.Ann[i].Data = append(req.Search.Ann[i].Data, &olama.VectorArray{Vector: vec})
		} else if text, ok := annParam.Data.(string); ok {
			req.Search.Ann[i].EmbeddingItems = append(req.Search.EmbeddingItems, text)
		} else {
			return nil, fmt.Errorf("hybridSearch failed, because of AnnParam.Data field type, " +
				"which must be []float32 or string")
		}

		if annParam.Params != nil {
			req.Search.Ann[i].Params = new(olama.SearchParams)
			req.Search.Ann[i].Params.Nprobe = annParam.Params.Nprobe
			req.Search.Ann[i].Params.Ef = annParam.Params.Ef
			req.Search.Ann[i].Params.Radius = annParam.Params.Radius
		}
		break
	}

	for i, matchParam := range params.Match {
		fieldName := "sparse_vector"
		if matchParam.FieldName != "" {
			fieldName = matchParam.FieldName
		}
		req.Search.Sparse = append(req.Search.Sparse, &olama.SparseData{
			FieldName: fieldName,
		})
		if matchParam.Limit != nil {
			req.Search.Sparse[i].Limit = uint32(*matchParam.Limit)
		}

		sparseVectorArray := make([]*olama.SparseVectorArray, 0)

		if svs, ok := matchParam.Data.([]encoder.SparseVecItem); ok {
			data := make([]*olama.SparseVecItem, 0)
			for _, sv := range svs {
				data = append(data, &olama.SparseVecItem{
					TermId: sv.TermId,
					Score:  sv.Score,
				})
			}
			sparseVectorArray = append(sparseVectorArray, &olama.SparseVectorArray{
				SpVector: data,
			})
			req.Search.Sparse[i].Data = sparseVectorArray
		} else {
			return nil, fmt.Errorf("hybridSearch failed, because of Match.Data field type, " +
				"which must be []encoder.SparseVecItem")
		}

		req.Search.Sparse[i].Params = new(olama.SparseSearchParams)
		req.Search.Sparse[i].Params.TerminateAfter = matchParam.TerminateAfter
		req.Search.Sparse[i].Params.CutoffFrequency = matchParam.CutoffFrequency
		break
	}

	if params.Rerank != nil {
		req.Search.RerankParams = new(olama.RerankParams)
		if len(params.Rerank.FieldList) != len(params.Rerank.Weight) {
			return nil, fmt.Errorf("the length of fieldlist should be equal with the length of weights")
		}
		req.Search.RerankParams.Method = string(params.Rerank.Method)
		req.Search.RerankParams.Weights = make(map[string]float32, 0)
		for i, fieldName := range params.Rerank.FieldList {
			req.Search.RerankParams.Weights[fieldName] = params.Rerank.Weight[i]
		}
		req.Search.RerankParams.RrfK = params.Rerank.RrfK
	}

	req.Search.Filter = params.Filter.Cond()
	req.Search.RetrieveVector = params.RetrieveVector
	req.Search.Outputfields = params.OutputFields
	req.Search.Limit = uint32(*params.Limit)

	res, err := r.rpcClient.HybridSearch(ctx, req)
	if err != nil {
		return nil, err
	}
	var documents [][]Document
	for _, result := range res.Results {
		var vecDoc []Document
		for _, doc := range result.Documents {
			d := Document{
				Id:     doc.Id,
				Vector: doc.Vector,
				Score:  doc.Score,
				Fields: make(map[string]Field),
			}

			d.SparseVector = make([]encoder.SparseVecItem, 0)
			for _, sv := range doc.SparseVector {
				d.SparseVector = append(d.SparseVector, encoder.SparseVecItem{
					TermId: sv.TermId,
					Score:  sv.Score,
				})
			}

			for n, v := range doc.Fields {
				d.Fields[n] = *ConvertGrpc2Field(v)
			}
			vecDoc = append(vecDoc, d)
		}
		documents = append(documents, vecDoc)
	}
	result := &SearchDocumentResult{
		Warning:   res.Warning,
		Documents: documents,
	}
	return result, nil
}

// [Delete] deletes documents by conditions.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - databaseName: The name of the database.
//   - collectionName: The name of the collection.
//   - params: A [DeleteDocumentParams] object that includes the other parameters for deleting documents' operation.
//     See [DeleteDocumentParams] for more information.
//
// Returns a pointer to a [DeleteDocumentResult] object or an error.
func (r *rpcImplementerFlatDocument) Delete(ctx context.Context, databaseName, collectionName string,
	param DeleteDocumentParams) (*DeleteDocumentResult, error) {
	req := &olama.DeleteRequest{
		Database:   databaseName,
		Collection: collectionName,
		Query: &olama.QueryCond{
			DocumentIds: param.DocumentIds,
			Filter:      param.Filter.Cond(),
			Limit:       param.Limit,
		},
	}
	res, err := r.rpcClient.Dele(ctx, req)
	if err != nil {
		return nil, err
	}
	return &DeleteDocumentResult{AffectedCount: int(res.AffectedCount)}, nil
}

// [Update] updates documents by conditions.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - databaseName: The name of the database.
//   - collectionName: The name of the collection.
//   - param: A [UpdateDocumentParams] object that includes the other parameters for updating documents' operation.
//     See [UpdateDocumentParams] for more information.
//
// Returns a pointer to a [UpdateDocumentResult] object or an error.
func (r *rpcImplementerFlatDocument) Update(ctx context.Context, databaseName, collectionName string,
	param UpdateDocumentParams) (*UpdateDocumentResult, error) {
	req := &olama.UpdateRequest{
		Database:   databaseName,
		Collection: collectionName,
		Query: &olama.QueryCond{
			DocumentIds: param.QueryIds,
			Filter:      param.QueryFilter.Cond(),
		},
		Update: &olama.Document{
			Vector: param.UpdateVector,
			Fields: make(map[string]*olama.Field),
		},
	}

	req.Update.SparseVector = make([]*olama.SparseVecItem, 0)
	for _, sv := range param.UpdateSparseVec {
		req.Update.SparseVector = append(req.Update.SparseVector, &olama.SparseVecItem{
			TermId: sv.TermId,
			Score:  sv.Score,
		})
	}

	if updatefields, ok := param.UpdateFields.(map[string]Field); ok {
		for k, v := range updatefields {
			req.Update.Fields[k] = ConvertField2Grpc(&v)
		}
	} else if updatefields, ok := param.UpdateFields.(map[string]interface{}); ok {
		if vector, ok := updatefields["vector"]; ok {
			if aVector, okV := vector.([]float32); okV {
				req.Update.Vector = aVector
				delete(updatefields, "vector")
			} else {
				return nil, fmt.Errorf("update failed, because of incorrect vector field type, " +
					"which must be []float32")
			}
		}

		if sparseVector, ok := updatefields["sparse_vector"]; ok {
			if aSparseVector, ok := sparseVector.([][]interface{}); ok {
				req.Update.SparseVector = make([]*olama.SparseVecItem, 0)
				for _, sv := range aSparseVector {
					if len(sv) != 2 {
						continue
					}
					req.Update.SparseVector = append(req.Update.SparseVector, &olama.SparseVecItem{
						TermId: int64(sv[0].(uint64)),
						Score:  float32(sv[1].(float64)),
					})
				}
				delete(updatefields, "sparse_vector")
			} else {
				return nil, fmt.Errorf("update failed, because of incorrect sparse_vector field type, which must be [][]interface{}")
			}
		}

		for k, v := range updatefields {
			req.Update.Fields[k] = ConvertField2Grpc(&Field{Val: v})
		}
	} else {
		return nil, fmt.Errorf("update failed, because of incorrect UpdateDocumentParams.UpdateFields field type, " +
			"which must be map[string]Field or map[string]interface{}")
	}

	res, err := r.rpcClient.Update(ctx, req)
	if err != nil {
		return nil, err
	}
	return &UpdateDocumentResult{AffectedCount: int(res.AffectedCount)}, nil
}

// [search] returns the most similar topK vectors by the given conditions.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - databaseName: The name of the database.
//   - collectionName: The name of the collection.
//   - documentIds: The list of the documents' ids, which are used for filtering documents.
//   - vectors: The list of vectors to search. The maximum number of elements in the array is 20.
//     Only one of the fields, vectors or documentIds, needs to be configured.
//   - text: It is a map where the keys represent column names, and the values are lists of column values to be retrieved,
//     which are used to retrieve data similar to ones.
//   - params: A pointer to a [SearchDocumentParams] object that includes the other parameters for searching documents' operation.
//     See [SearchDocumentParams] for more information.
//
// Returns a pointer to a [SearchDocumentResult] object or an error.
func (r *rpcImplementerFlatDocument) search(ctx context.Context, databaseName, collectionName string,
	documentIds []string, vectors [][]float32, text map[string][]string, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	req := &olama.SearchRequest{
		Database:        databaseName,
		Collection:      collectionName,
		ReadConsistency: string(r.SdkClient.Options().ReadConsistency),
		Search:          &olama.SearchCond{},
	}
	req.Search.DocumentIds = documentIds
	vectorArray := make([]*olama.VectorArray, 0, len(req.Search.Vectors))
	for _, vector := range vectors {
		vectorArray = append(vectorArray, &olama.VectorArray{Vector: vector})
	}
	req.Search.Vectors = vectorArray
	for _, v := range text {
		req.Search.EmbeddingItems = v
	}
	if len(params) != 0 && params[0] != nil {
		param := params[0]
		req.Search.Filter = param.Filter.Cond()
		req.Search.RetrieveVector = param.RetrieveVector
		req.Search.Outputfields = param.OutputFields
		req.Search.Limit = uint32(param.Limit)
		if param.Params != nil {
			req.Search.Params = &olama.SearchParams{
				Nprobe: param.Params.Nprobe,
				Ef:     param.Params.Ef,
			}

		}
		if param.Radius != nil {
			if req.Search.Params == nil {
				req.Search.Params = new(olama.SearchParams)
			}
			req.Search.Range = true
			req.Search.Params.Radius = *param.Radius
		}
	}
	res, err := r.rpcClient.Search(ctx, req)
	if err != nil {
		return nil, err
	}
	var documents [][]Document
	for _, result := range res.Results {
		var vecDoc []Document
		for _, doc := range result.Documents {
			d := Document{
				Id:     doc.Id,
				Vector: doc.Vector,
				Score:  doc.Score,
				Fields: make(map[string]Field),
			}
			for n, v := range doc.Fields {
				d.Fields[n] = *ConvertGrpc2Field(v)
			}
			vecDoc = append(vecDoc, d)
		}
		documents = append(documents, vecDoc)
	}
	result := &SearchDocumentResult{
		Warning:   res.Warning,
		Documents: documents,
	}
	return result, nil
}

// [Count] counts the number of documents in a collection that satisfy the specified filter conditions.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - databaseName: The name of the database.
//   - collectionName: The name of the collection.
//   - param: A [CountDocumentParams] object that includes the other parameters for counting documents' operation.
//     See [CountDocumentParams] for more information.
//
// Returns a pointer to a [CountDocumentResult] object or an error.
func (r *rpcImplementerFlatDocument) Count(ctx context.Context, databaseName, collectionName string,
	params ...CountDocumentParams) (*CountDocumentResult, error) {
	req := &olama.CountRequest{
		Database:   databaseName,
		Collection: collectionName,
	}

	if len(params) != 0 {
		param := params[0]
		req.Query = &olama.QueryCond{
			Filter: param.CountFilter.Cond(),
		}
	}
	res, err := r.rpcClient.Count(ctx, req)
	if err != nil {
		return nil, err
	}
	return &CountDocumentResult{Count: res.Count}, nil
}
