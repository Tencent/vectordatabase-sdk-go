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
	"reflect"

	"github.com/tencent/vectordatabase-sdk-go/tcvdbtext/encoder"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/document"
)

var _ DocumentInterface = &implementerDocument{}
var _ FlatInterface = &implementerFlatDocument{}

// DocumentInterface document api
type DocumentInterface interface {
	SdkClient

	// [Upsert] upserts documents into a collection.
	Upsert(ctx context.Context, documents interface{}, params ...*UpsertDocumentParams) (result *UpsertDocumentResult, err error)

	// [Query] queries documents that satisfies the condition from the collection.
	Query(ctx context.Context, documentIds []string, params ...*QueryDocumentParams) (result *QueryDocumentResult, err error)

	// [Search] returns the most similar topK vectors by the given vectors.
	// Search is a Batch API.
	Search(ctx context.Context, vectors [][]float32, params ...*SearchDocumentParams) (result *SearchDocumentResult, err error)

	// [HybridSearch] retrieves both dense and sparse vectors to return the most similar topK vectors.
	HybridSearch(ctx context.Context, params HybridSearchDocumentParams) (result *SearchDocumentResult, err error)

	// [SearchById] returns the most similar topK vectors by the given documentIds.
	SearchById(ctx context.Context, documentIds []string, params ...*SearchDocumentParams) (result *SearchDocumentResult, err error)

	// [SearchByText] returns the most similar topK vectors by the given text map.
	// The texts will be firstly embedded into vectors using the embedding model of the collection on the server.
	SearchByText(ctx context.Context, text map[string][]string, params ...*SearchDocumentParams) (result *SearchDocumentResult, err error)

	// [Delete] deletes documents by conditions.
	Delete(ctx context.Context, param DeleteDocumentParams) (result *DeleteDocumentResult, err error)

	// [Update] updates documents by conditions.
	Update(ctx context.Context, param UpdateDocumentParams) (result *UpdateDocumentResult, err error)

	// [Count] counts the number of documents in a collection that satisfy the specified filter conditions.
	Count(ctx context.Context, params ...CountDocumentParams) (*CountDocumentResult, error)
}

type FlatInterface interface {
	// [Upsert] upserts documents into a collection.
	Upsert(ctx context.Context, databaseName, collectionName string, documents interface{}, params ...*UpsertDocumentParams) (result *UpsertDocumentResult, err error)

	// [Query] queries documents that satisfies the condition from the collection.
	Query(ctx context.Context, databaseName, collectionName string, documentIds []string, params ...*QueryDocumentParams) (result *QueryDocumentResult, err error)

	// [Search] returns the most similar topK vectors by the given vectors.
	// Search is a Batch API.
	Search(ctx context.Context, databaseName, collectionName string, vectors [][]float32, params ...*SearchDocumentParams) (result *SearchDocumentResult, err error)

	// [HybridSearch] retrieves both dense and sparse vectors to return the most similar topK vectors.
	HybridSearch(ctx context.Context, databaseName, collectionName string, params HybridSearchDocumentParams) (result *SearchDocumentResult, err error)

	// [SearchById] returns the most similar topK vectors by the given documentIds.
	SearchById(ctx context.Context, databaseName, collectionName string, documentIds []string, params ...*SearchDocumentParams) (result *SearchDocumentResult, err error)

	// [SearchByText] returns the most similar topK vectors by the given text map.
	// The texts will be firstly embedded into vectors using the embedding model of the collection on the server.
	SearchByText(ctx context.Context, databaseName, collectionName string, text map[string][]string, params ...*SearchDocumentParams) (result *SearchDocumentResult, err error)

	// [Delete] deletes documents by conditions.
	Delete(ctx context.Context, databaseName, collectionName string, param DeleteDocumentParams) (result *DeleteDocumentResult, err error)

	// [Update] updates documents by conditions.
	Update(ctx context.Context, databaseName, collectionName string, param UpdateDocumentParams) (result *UpdateDocumentResult, err error)

	// [Count] counts the number of documents in a collection that satisfy the specified filter conditions.
	Count(ctx context.Context, databaseName, collectionName string,
		params ...CountDocumentParams) (*CountDocumentResult, error)
}

type implementerDocument struct {
	SdkClient
	flat       FlatInterface
	database   *Database
	collection *Collection
}

// [UpsertDocumentParams] holds the parameters for upserting documents to a collection.
//
// Fields:
//   - BuildIndex:  (Optional) if BuildIndex is true, the upserted documents' indexes will be built immediately,
//     which will affect the performance of upsert.
type UpsertDocumentParams struct {
	BuildIndex *bool
}

type UpsertDocumentResult struct {
	AffectedCount int
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
// Notes: The name of the database and the name of collection are from the fields of [implementerDocument].
//
// Returns a pointer to a [UpsertDocumentResult] object or an error.
func (i *implementerDocument) Upsert(ctx context.Context, documents interface{}, params ...*UpsertDocumentParams) (result *UpsertDocumentResult, err error) {
	return i.flat.Upsert(ctx, i.database.DatabaseName, i.collection.CollectionName, documents, params...)
}

// [QueryDocumentParams] holds the parameters for querying documents to a collection.
//
// Fields:
//   - Filter:  (Optional) Filter documents by [Filter] conditions before returning the result.
//   - RetrieveVector: (Optional) Specify whether to return vector values or not (default to false).
//     If RetrieveVector is true, the vector values will be returned.
//   - OutputFields: (Optional) Return columns specified by the list of column names.
//   - Offset: (Optional) Skip a specified number of documents in the query result set.
//   - Limit: (Optional) Limit the number of documents returned (default to 1).
type QueryDocumentParams struct {
	Filter         *Filter
	RetrieveVector bool
	OutputFields   []string
	Offset         int64
	Limit          int64
}

type QueryDocumentResult struct {
	Documents     []Document
	AffectedCount int
	Total         uint64
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
// Notes: The name of the database and the name of collection are from the fields of [implementerDocument].
//
// Returns a pointer to a [QueryDocumentResult] object or an error.
func (i *implementerDocument) Query(ctx context.Context, documentIds []string, params ...*QueryDocumentParams) (*QueryDocumentResult, error) {
	return i.flat.Query(ctx, i.database.DatabaseName, i.collection.CollectionName, documentIds, params...)
}

// [SearchDocumentParams] holds the parameters for searching documents to a collection.
//
// Fields:
//   - Filter:  (Optional) Filter documents by [Filter] conditions before searching the results.
//   - Params: A pointer to a [SearchDocParams] object that includes the other parameters for searching documents' operation.
//     See [SearchDocParams] for more information.
//   - RetrieveVector: (Optional) Specify whether to return vector values or not (default to false).
//     If RetrieveVector is true, the vector values will be returned.
//   - OutputFields: (Optional) Return columns specified by the list of column names.
//   - Limit: (Required) Limit the number of documents returned (default to 1).
type SearchDocumentParams struct {
	Filter         *Filter
	Params         *SearchDocParams
	RetrieveVector bool
	OutputFields   []string
	Limit          int64
}

// [SearchDocParams] holds the parameters for searching documents to a collection.
//
// Fields:
//   - Nprobe: (Optional)  IVF type index requires configuration parameter nprobe to specify the number
//     of vectors to be accessed. Valid range is [1, nlist], and nlist is defined by creating collection.
//   - Ef: (Optional) HNSW type index requires configuration parameter ef to specify the number
//     of vectors to be accessed (default to 10). Valid range is [1, 32768]
//   - Radius: (Optional) Specifies the radius range for similarity retrieval.
type SearchDocParams struct {
	Nprobe uint32  `json:"nprobe,omitempty"` // 搜索时查找的聚类数量，使用索引默认值即可
	Ef     uint32  `json:"ef,omitempty"`     // HNSW
	Radius float32 `json:"radius,omitempty"` // 距离阈值,范围搜索时有效
}

type SearchDocumentResult struct {
	Warning   string
	Documents [][]Document
}

// [Search] returns the most similar topK vectors by the given vectors.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - vectors: The list of vectors to search. The maximum number of elements in the array is 20.
//   - params: A pointer to a [SearchDocumentParams] object that includes the other parameters for searching documents' operation.
//     See [SearchDocumentParams] for more information.
//
// Notes: The name of the database and the name of collection are from the fields of [implementerDocument].
//
// Returns a pointer to a [SearchDocumentResult] object or an error.
func (i *implementerDocument) Search(ctx context.Context, vectors [][]float32, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return i.flat.Search(ctx, i.database.DatabaseName, i.collection.CollectionName, vectors, params...)
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
// Notes: The name of the database and the name of collection are from the fields of [implementerDocument].
//
// Returns a pointer to a [SearchDocumentResult] object or an error.
func (i *implementerDocument) SearchById(ctx context.Context, documentIds []string, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return i.flat.SearchById(ctx, i.database.DatabaseName, i.collection.CollectionName, documentIds, params...)
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
// Notes: The name of the database and the name of collection are from the fields of [implementerDocument].
//
// Returns a pointer to a [SearchDocumentResult] object or an error.
func (i *implementerDocument) SearchByText(ctx context.Context, text map[string][]string, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return i.flat.SearchByText(ctx, i.database.DatabaseName, i.collection.CollectionName, text, params...)
}

// [HybridSearchDocumentParams] holds the parameters for hybrid searching documents to a collection.
//
// Fields:
//   - Filter:  (Optional) Filter documents by [Filter] conditions before hybrid searching the results.
//   - Params: A pointer to a [SearchDocParams] object that includes the other parameters for hybrid searching documents' operation.
//     See [SearchDocParams] for more information.
//   - RetrieveVector: (Optional) Specify whether to return vector values or not (default to false).
//     If RetrieveVector is true, the vector values will be returned.
//   - OutputFields: (Optional) Return columns specified by the list of column names.
//   - Limit: (Required) Limit the number of documents returned (default to 1).
//   - AnnParams: The list of [AnnParam] pointers for vectors retrieval configuration.
//     See [AnnParam] for more information.
//   - Rerank: A pointer to a [RerankOption] object for re-ranking configuration in retrieval.
//     See [RerankOption] for more information.
//   - Match: The list of [MatchOption] pointers for sparse vectors retrieval configuration.
//     See [MatchOption] for more information.
type HybridSearchDocumentParams struct {
	Filter         *Filter
	Params         *SearchDocParams
	RetrieveVector bool
	OutputFields   []string
	Limit          *int

	AnnParams []*AnnParam
	Rerank    *RerankOption
	Match     []*MatchOption
}

// [RerankOption] holds the parameters for re-ranking configuration in retrieval.
//
// Fields:
//   - Method: The parameter method specifies the method for Rerank.
//     It can take the following enumerated values: RerankWeighted(weighted), RerankRrf(rrf)
//   - FieldList: Lists the fields used for weighted calculation. For example, FieldList:
//     []string{"vector", "sparse_vector"} represents weighted calculation for dense and sparse vectors.
//   - Weight: Sorts based on a weighted combination of scores from different fields.
//   - RrfK: If Method is RerankRrf, you should config the RrfK, which is used to calculate the reciprocal rank score.
//     It adjusts the scoring formula to control the distribution of rank scores (default to 60).
type RerankOption struct {
	Method    RerankMethod
	FieldList []string
	Weight    []float32
	RrfK      int32
}

// [MatchOption] holds the parameters for sparse vectors retrieval configuration.
//
// Fields:
//   - FieldName: The field name for sparse vector retrieval, for example: sparse_vector.
//   - Data: The sparse vectors to retrieve, supporting only sparse vectors for one sentence.
//   - Limit: The number of results returned from sparse vector retrieval.
//   - TerminateAfter: (Optional) Threshold for early termination of keyword search,
//     used to improve search efficiency.
//   - CutoffFrequency: (Optional) CutoffFrequency specifies a positive integer limit, ranging from [1, +∞].
//     If the term frequency is less than the cutoffFrequency, the term will be ignored during retrieval.
//     It also supports decimal values within the range [0,1].
type MatchOption struct {
	FieldName       string
	Data            interface{}
	Limit           *int
	TerminateAfter  uint32
	CutoffFrequency float64
}

// [AnnParam] holds the parameters for vectors hybrid retrieval configuration.
//
// Fields:
//   - FieldName: The field name for retrieval, and you can set vector or id.
//   - Data: The vectors to retrieve, supporting only vectors for one document.
//   - Params: A pointer to a [SearchDocParams] object that includes the other parameters for searching documents' operation.
//     See [SearchDocParams] for more information.
//   - Limit: The number of results returned from vector retrieval.
type AnnParam struct {
	FieldName string
	Data      interface{}
	Params    *SearchDocParams
	Limit     *int
}

// [HybridSearch] retrieves both dense and sparse vectors to return the most similar topK vectors.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - params: A [HybridSearchDocumentParams] object that includes the other parameters for hybrid searching documents' operation.
//     See [HybridSearchDocumentParams] for more information.
//
// Notes: The name of the database and the name of collection are from the fields of [implementerDocument].
//
// Returns a pointer to a [SearchDocumentResult] object or an error.
func (i *implementerDocument) HybridSearch(ctx context.Context, params HybridSearchDocumentParams) (*SearchDocumentResult, error) {
	return i.flat.HybridSearch(ctx, i.database.DatabaseName, i.collection.CollectionName, params)
}

// [DeleteDocumentParams] holds the parameters for deleting documents to a collection.
//
// Fields:
//   - DocumentIds: The list of the documents' ids to delete.  The maximum size of the array is 20.
//   - Filter:  (Optional) Filter documents by [Filter] conditions to delete.
//   - Limit: (Optional) Limit the number of documents deleted.
type DeleteDocumentParams struct {
	DocumentIds []string
	Filter      *Filter
	Limit       int64
}

type DeleteDocumentResult struct {
	AffectedCount int
}

// [Delete] deletes documents by conditions.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - params: A [DeleteDocumentParams] object that includes the other parameters for deleting documents' operation.
//     See [DeleteDocumentParams] for more information.
//
// Notes: The name of the database and the name of collection are from the fields of [implementerDocument].
//
// Returns a pointer to a [DeleteDocumentResult] object or an error.
func (i *implementerDocument) Delete(ctx context.Context, param DeleteDocumentParams) (result *DeleteDocumentResult, err error) {
	return i.flat.Delete(ctx, i.database.DatabaseName, i.collection.CollectionName, param)
}

// [UpdateDocumentParams] holds the parameters for updating documents to a collection.
//
// Fields:
//   - QueryIds: The list of the documents' ids to update.  The maximum size of the array is 20.
//   - QueryFilter: (Optional) Filter documents by [Filter] conditions to update.
//   - UpdateVector: The values with which you want to update the vector, and the updated documents are queried by QueryIds and QueryFilter.
//   - UpdateSparseVec: The sparse values with which you want to update the vector, and the updated documents are queried by QueryIds and QueryFilter.
//   - UpdateFields: Update documents' fields by this value, and the updated documents are queried by QueryIds and QueryFilter.
type UpdateDocumentParams struct {
	QueryIds        []string
	QueryFilter     *Filter
	UpdateVector    []float32
	UpdateSparseVec []encoder.SparseVecItem
	UpdateFields    interface{}
}

type UpdateDocumentResult struct {
	AffectedCount int
}

// [Update] updates documents by conditions.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - param: A [UpdateDocumentParams] object that includes the other parameters for updating documents' operation.
//     See [UpdateDocumentParams] for more information.
//
// Notes: The name of the database and the name of collection are from the fields of [implementerDocument].
//
// Returns a pointer to a [UpdateDocumentResult] object or an error.
func (i *implementerDocument) Update(ctx context.Context, param UpdateDocumentParams) (*UpdateDocumentResult, error) {
	return i.flat.Update(ctx, i.database.DatabaseName, i.collection.CollectionName, param)
}

// [Count] counts the number of documents in a collection that satisfy the specified filter conditions.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - param: A [CountDocumentParams] object that includes the other parameters for counting documents' operation.
//     See [CountDocumentParams] for more information.
//
// Notes: The name of the database and the name of collection are from the fields of [implementerDocument].
//
// Returns a pointer to a [CountDocumentResult] object or an error.
func (i *implementerDocument) Count(ctx context.Context, params ...CountDocumentParams) (*CountDocumentResult, error) {
	return i.flat.Count(ctx, i.database.DatabaseName, i.collection.CollectionName, params...)
}

type Document struct {
	Id           string                  `json:"id"`
	Vector       []float32               `json:"vector"`
	SparseVector []encoder.SparseVecItem `json:"sparse_vector"`
	// omitempty when upsert
	Score  float32 `json:"score"`
	Fields map[string]Field
}

type implementerFlatDocument struct {
	SdkClient
}

// [Upsert] upserts documents into a collection.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - db: The name of the database.
//   - coll: The name of the collection.
//   - documents: The list of the [Document] object or  map[string]interface{} to upsert. Maximum 1000.
//   - params: A pointer to a [UpsertDocumentParams] object that includes the other parameters for upserting documents' operation.
//     See [UpsertDocumentParams] for more information.
//
// Returns a pointer to a [UpsertDocumentResult] object or an error.
func (i *implementerFlatDocument) Upsert(ctx context.Context, db, coll string, documents interface{}, params ...*UpsertDocumentParams) (result *UpsertDocumentResult, err error) {
	req := new(document.UpsertReq)
	req.Database = db
	req.Collection = coll

	if docs, ok := documents.([]Document); ok {
		for _, doc := range docs {
			d := &document.Document{}
			d.Id = doc.Id
			d.Vector = doc.Vector

			d.SparseVector = make([][]interface{}, 0)
			for _, sv := range doc.SparseVector {
				d.SparseVector = append(d.SparseVector, []interface{}{sv.TermId, sv.Score})
			}

			d.Fields = make(map[string]interface{})
			for k, v := range doc.Fields {
				d.Fields[k] = v.Val
			}
			req.Documents = append(req.Documents, d)
		}
	} else if docs, ok := documents.([]map[string]interface{}); ok {
		for _, doc := range docs {
			d := &document.Document{}
			if id, ok := doc["id"]; ok {
				if sId, ok := id.(string); ok {
					d.Id = sId
					delete(doc, "id")
				} else {
					return nil, fmt.Errorf("upsert failed, because of incorrect id field type, which must be string")
				}
			}
			if vector, ok := doc["vector"]; ok {
				if aVector, ok := vector.([]float32); ok {
					d.Vector = aVector
					delete(doc, "vector")
				} else {
					return nil, fmt.Errorf("upsert failed, because of incorrect vector field type, which must be []float32")
				}
			}
			if sparseVector, ok := doc["sparse_vector"]; ok {
				if aSparseVector, ok := sparseVector.([][]interface{}); ok {
					d.SparseVector = make([][]interface{}, 0)
					for _, sv := range aSparseVector {
						svItem, err := ConvSliceInterface2SparseVecItem(sv)
						if err != nil {
							return nil, fmt.Errorf("upsert failed. doc's sparse_vector data is incorrect. doc id is %v. err: %v", d.Id, err.Error())
						}
						d.SparseVector = append(d.SparseVector, []interface{}{svItem.TermId, svItem.Score})
					}
					delete(doc, "sparse_vector")
				} else {
					return nil, fmt.Errorf("upsert failed, because of incorrect sparse_vector field type, which must be [][]interface{}")
				}
			}

			d.Fields = make(map[string]interface{})
			for k, v := range doc {
				d.Fields[k] = v
			}
			req.Documents = append(req.Documents, d)
		}
	} else {
		return nil, fmt.Errorf("upsert failed, because of incorrect documents type, which must be []Document or []map[string]interface{}")
	}

	if len(params) != 0 && params[0] != nil {
		param := params[0]
		if param.BuildIndex != nil {
			req.BuildIndex = param.BuildIndex
		}
	}

	res := new(document.UpsertRes)
	result = new(UpsertDocumentResult)
	err = i.Request(ctx, req, res)
	if err != nil {
		return
	}
	result.AffectedCount = int(res.AffectedCount)
	return
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
func (i *implementerFlatDocument) Query(ctx context.Context, databaseName, collectionName string,
	documentIds []string, params ...*QueryDocumentParams) (*QueryDocumentResult, error) {
	req := new(document.QueryReq)
	req.Database = databaseName
	req.Collection = collectionName
	req.Query = &document.QueryCond{
		DocumentIds: documentIds,
	}
	req.ReadConsistency = string(i.SdkClient.Options().ReadConsistency)
	if len(params) != 0 && params[0] != nil {
		param := params[0]
		req.Query.Filter = param.Filter.Cond()
		req.Query.RetrieveVector = param.RetrieveVector
		req.Query.OutputFields = param.OutputFields
		req.Query.Offset = param.Offset
		req.Query.Limit = param.Limit
	}

	res := new(document.QueryRes)
	err := i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}

	result := new(QueryDocumentResult)
	var documents []Document
	for _, doc := range res.Documents {
		var d Document
		d.Id = doc.Id
		d.Vector = doc.Vector

		d.SparseVector = make([]encoder.SparseVecItem, 0)
		for _, sv := range doc.SparseVector {
			svItem, err := ConvSliceInterface2SparseVecItem(sv)
			if err != nil {
				return nil, fmt.Errorf("query failed. doc's sparse_vector data is incorrect. doc id is %v. err: %v", d.Id, err.Error())
			}
			d.SparseVector = append(d.SparseVector, *svItem)
		}

		d.Fields = make(map[string]Field)

		for n, v := range doc.Fields {
			d.Fields[n] = Field{Val: v}
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
func (i *implementerFlatDocument) Search(ctx context.Context, databaseName, collectionName string,
	vectors [][]float32, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return i.search(ctx, databaseName, collectionName, nil, vectors, nil, params...)
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
func (i *implementerFlatDocument) SearchById(ctx context.Context, databaseName, collectionName string,
	documentIds []string, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return i.search(ctx, databaseName, collectionName, documentIds, nil, nil, params...)
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
func (i *implementerFlatDocument) SearchByText(ctx context.Context, databaseName, collectionName string,
	text map[string][]string, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return i.search(ctx, databaseName, collectionName, nil, nil, text, params...)
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
func (i *implementerFlatDocument) search(ctx context.Context, databaseName, collectionName string,
	documentIds []string, vectors [][]float32, text map[string][]string, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	req := new(document.SearchReq)
	req.Database = databaseName
	req.Collection = collectionName
	req.ReadConsistency = string(i.SdkClient.Options().ReadConsistency)
	req.Search = new(document.SearchCond)
	req.Search.DocumentIds = documentIds
	req.Search.Vectors = vectors
	for _, v := range text {
		req.Search.EmbeddingItems = v
	}

	if len(params) != 0 && params[0] != nil {
		param := params[0]
		req.Search.Filter = param.Filter.Cond()
		req.Search.RetrieveVector = param.RetrieveVector
		req.Search.OutputFields = param.OutputFields
		req.Search.Limit = param.Limit

		if param.Params != nil {
			req.Search.Params = new(document.SearchParams)
			req.Search.Params.Nprobe = param.Params.Nprobe
			req.Search.Params.Ef = param.Params.Ef
			req.Search.Params.Radius = param.Params.Radius
		}
	}

	res := new(document.SearchRes)
	err := i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}
	var documents [][]Document
	for _, result := range res.Documents {
		var vecDoc []Document
		for _, doc := range result {
			d := Document{
				Id:     doc.Id,
				Vector: doc.Vector,
				Score:  doc.Score,
				Fields: make(map[string]Field),
			}
			for n, v := range doc.Fields {
				d.Fields[n] = Field{Val: v}
			}
			vecDoc = append(vecDoc, d)
		}
		documents = append(documents, vecDoc)
	}
	result := new(SearchDocumentResult)
	result.Warning = res.Warning
	result.Documents = documents
	return result, nil
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
func (i *implementerFlatDocument) HybridSearch(ctx context.Context, databaseName, collectionName string,
	params HybridSearchDocumentParams) (*SearchDocumentResult, error) {
	req := new(document.HybridSearchReq)
	req.Database = databaseName
	req.Collection = collectionName
	req.ReadConsistency = string(i.SdkClient.Options().ReadConsistency)
	req.Search = new(document.HybridSearchCond)
	req.Search.AnnParams = make([]*document.AnnParam, 0)
	req.Search.Match = make([]*document.MatchOption, 0)

	for i, annParam := range params.AnnParams {
		fieldName := "vector"
		if annParam.FieldName != "" {
			fieldName = annParam.FieldName
		}
		req.Search.AnnParams = append(req.Search.AnnParams, &document.AnnParam{
			FieldName: fieldName,
			Limit:     annParam.Limit,
		})

		req.Search.AnnParams[i].Data = make([]interface{}, 0)
		if vec, ok := annParam.Data.([]float32); ok {
			req.Search.AnnParams[i].Data = append(req.Search.AnnParams[i].Data, vec)
		} else if text, ok := annParam.Data.(string); ok {
			req.Search.AnnParams[i].Data = append(req.Search.AnnParams[i].Data, text)
		} else {
			return nil, fmt.Errorf("hybridSearch failed, because of AnnParam.Data field type, " +
				"which must be []float32")
		}

		if annParam.Params != nil {
			req.Search.AnnParams[i].Params = new(document.SearchParams)
			req.Search.AnnParams[i].Params.Nprobe = annParam.Params.Nprobe
			req.Search.AnnParams[i].Params.Ef = annParam.Params.Ef
			req.Search.AnnParams[i].Params.Radius = annParam.Params.Radius
		}
		break
	}

	for i, matchParam := range params.Match {
		fieldName := "sparse_vector"
		if matchParam.FieldName != "" {
			fieldName = matchParam.FieldName
		}
		req.Search.Match = append(req.Search.Match, &document.MatchOption{
			FieldName: fieldName,
		})

		req.Search.Match[i].Data = make([][][]interface{}, 0)
		if svs, ok := matchParam.Data.([]encoder.SparseVecItem); ok {
			sparseVector := make([][]interface{}, 0)
			for _, svItem := range svs {
				sparseVector = append(sparseVector, []interface{}{svItem.TermId, svItem.Score})
			}
			req.Search.Match[i].Data = append(req.Search.Match[i].Data, sparseVector)
		} else {
			return nil, fmt.Errorf("hybridSearch failed, because of Match.Data field type, " +
				"which must be []encoder.SparseVecItem")
		}

		if matchParam.Limit != nil {
			req.Search.Match[i].Limit = *matchParam.Limit
		}
		req.Search.Match[i].TerminateAfter = matchParam.TerminateAfter
		req.Search.Match[i].CutoffFrequency = matchParam.CutoffFrequency
		break
	}

	if params.Rerank != nil {
		req.Search.Rerank = new(document.RerankOption)
		req.Search.Rerank.FieldList = params.Rerank.FieldList
		req.Search.Rerank.Method = string(params.Rerank.Method)
		req.Search.Rerank.Weight = params.Rerank.Weight
		req.Search.Rerank.RrfK = params.Rerank.RrfK
	}

	req.Search.Filter = params.Filter.Cond()
	req.Search.RetrieveVector = params.RetrieveVector
	req.Search.OutputFields = params.OutputFields
	req.Search.Limit = params.Limit

	res := new(document.SearchRes)
	err := i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}
	var documents [][]Document
	for _, result := range res.Documents {
		var vecDoc []Document
		for _, doc := range result {
			d := Document{
				Id:     doc.Id,
				Vector: doc.Vector,
				Score:  doc.Score,
				Fields: make(map[string]Field),
			}

			d.SparseVector = make([]encoder.SparseVecItem, 0)
			for _, sv := range doc.SparseVector {
				svItem, err := ConvSliceInterface2SparseVecItem(sv)
				if err != nil {
					return nil, fmt.Errorf("the search response's doc sparse_vector data is incorrect. doc id is %v. err: %v", d.Id, err.Error())
				}
				d.SparseVector = append(d.SparseVector, *svItem)
			}

			for n, v := range doc.Fields {
				d.Fields[n] = Field{Val: v}
			}
			vecDoc = append(vecDoc, d)
		}
		documents = append(documents, vecDoc)
	}
	result := new(SearchDocumentResult)
	result.Warning = res.Warning
	result.Documents = documents
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
func (i *implementerFlatDocument) Delete(ctx context.Context, databaseName, collectionName string,
	param DeleteDocumentParams) (*DeleteDocumentResult, error) {
	req := new(document.DeleteReq)
	req.Database = databaseName
	req.Collection = collectionName
	req.Query = &document.QueryCond{
		DocumentIds: param.DocumentIds,
		Filter:      param.Filter.Cond(),
		Limit:       param.Limit,
	}

	res := new(document.DeleteRes)
	result := new(DeleteDocumentResult)
	err := i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}
	result.AffectedCount = res.AffectedCount
	return result, nil
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
func (i *implementerFlatDocument) Update(ctx context.Context, databaseName, collectionName string,
	param UpdateDocumentParams) (*UpdateDocumentResult, error) {
	req := new(document.UpdateReq)
	req.Database = databaseName
	req.Collection = collectionName
	req.Query = new(document.QueryCond)

	req.Query.DocumentIds = param.QueryIds
	req.Query.Filter = param.QueryFilter.Cond()
	req.Update.Vector = param.UpdateVector
	req.Update.SparseVector = make([][]interface{}, 0)
	for _, sv := range param.UpdateSparseVec {
		req.Update.SparseVector = append(req.Update.SparseVector, []interface{}{sv.TermId, sv.Score})
	}
	req.Update.Fields = make(map[string]interface{}, 0)

	if updatefields, ok := param.UpdateFields.(map[string]Field); ok {
		for k, v := range updatefields {
			req.Update.Fields[k] = v.Val
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
		if vector, ok := updatefields["sparse_vector"]; ok {
			if aSparseVector, okV := vector.([][]interface{}); okV {
				req.Update.SparseVector = aSparseVector
				delete(updatefields, "sparse_vector")
			} else {
				return nil, fmt.Errorf("update failed, because of incorrect sparse_vector field type, which must be [][]interface{}")
			}
		}
		for k, v := range updatefields {
			req.Update.Fields[k] = v
		}
	} else if param.UpdateFields != nil {
		return nil, fmt.Errorf("update failed, because of incorrect UpdateDocumentParams.UpdateFields field type, " +
			"which must be map[string]Field or map[string]interface{}")
	}

	res := new(document.UpdateRes)
	result := new(UpdateDocumentResult)
	err := i.Request(ctx, req, res)
	if err != nil {
		return result, err
	}
	result.AffectedCount = res.AffectedCount
	return result, nil
}

// [CountDocumentParams] holds the parameters for counting the number of documents to a collection based on the [Filter] conditions.
//
// Fields:
//   - CountFilter: (Optional) Filter documents by [Filter] conditions to count.
type CountDocumentParams struct {
	CountFilter *Filter
}

// [CountDocumentResult] holds the results for counting the number of documents to a collection based on the [Filter] conditions.
//
// Fields:
//   - Count: The number of documents to a collection based on the [Filter].
type CountDocumentResult struct {
	Count uint64
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
func (i *implementerFlatDocument) Count(ctx context.Context, databaseName, collectionName string,
	params ...CountDocumentParams) (*CountDocumentResult, error) {
	req := new(document.CountReq)
	req.Database = databaseName
	req.Collection = collectionName
	req.Query = new(document.CountQueryCond)

	if len(params) != 0 {
		param := params[0]
		req.Query.Filter = param.CountFilter.Cond()
	}

	res := new(document.CountRes)
	result := new(CountDocumentResult)
	err := i.Request(ctx, req, res)
	if err != nil {
		return result, err
	}
	result.Count = res.Count
	return result, nil
}

func ConvSliceInterface2SparseVecItem(sv []interface{}) (*encoder.SparseVecItem, error) {

	svItem := new(encoder.SparseVecItem)
	if len(sv) != 2 {
		return nil, fmt.Errorf("incorrect sparse_vector data %v, which the length of token sparse_vector must is 2, but current is %v",
			sv, len(sv))
	}

	switch v := sv[0].(type) {
	case int, int8, int16, int32, int64:
		svItem.TermId = int64(reflect.ValueOf(v).Int())
	case uint, uint8, uint16, uint32, uint64:
		svItem.TermId = int64(reflect.ValueOf(v).Uint())
	case float32, float64:
		svItem.TermId = int64(reflect.ValueOf(v).Float())
	default:
		return nil, fmt.Errorf("incorrect sparse_vector data %v, which first item datatype must be int64", sv)
	}

	switch v := sv[1].(type) {
	case float32, float64:
		svItem.Score = float32(reflect.ValueOf(v).Float())
	default:
		return nil, fmt.Errorf("incorrect sparse_vector data %v, which second item datatype must be float32/float64", sv)
	}

	return svItem, nil
}
