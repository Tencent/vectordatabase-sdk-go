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
	Upsert(ctx context.Context, documents interface{}, params ...*UpsertDocumentParams) (result *UpsertDocumentResult, err error)
	Query(ctx context.Context, documentIds []string, params ...*QueryDocumentParams) (result *QueryDocumentResult, err error)
	Search(ctx context.Context, vectors [][]float32, params ...*SearchDocumentParams) (result *SearchDocumentResult, err error)
	HybridSearch(ctx context.Context, params HybridSearchDocumentParams) (result *SearchDocumentResult, err error)
	SearchById(ctx context.Context, documentIds []string, params ...*SearchDocumentParams) (result *SearchDocumentResult, err error)
	SearchByText(ctx context.Context, text map[string][]string, params ...*SearchDocumentParams) (result *SearchDocumentResult, err error)
	Delete(ctx context.Context, param DeleteDocumentParams) (result *DeleteDocumentResult, err error)
	Update(ctx context.Context, param UpdateDocumentParams) (result *UpdateDocumentResult, err error)
}

type FlatInterface interface {
	Upsert(ctx context.Context, databaseName, collectionName string, documents interface{}, params ...*UpsertDocumentParams) (result *UpsertDocumentResult, err error)
	Query(ctx context.Context, databaseName, collectionName string, documentIds []string, params ...*QueryDocumentParams) (result *QueryDocumentResult, err error)
	Search(ctx context.Context, databaseName, collectionName string, vectors [][]float32, params ...*SearchDocumentParams) (result *SearchDocumentResult, err error)
	HybridSearch(ctx context.Context, databaseName, collectionName string, params HybridSearchDocumentParams) (result *SearchDocumentResult, err error)
	SearchById(ctx context.Context, databaseName, collectionName string, documentIds []string, params ...*SearchDocumentParams) (result *SearchDocumentResult, err error)
	SearchByText(ctx context.Context, databaseName, collectionName string, text map[string][]string, params ...*SearchDocumentParams) (result *SearchDocumentResult, err error)
	Delete(ctx context.Context, databaseName, collectionName string, param DeleteDocumentParams) (result *DeleteDocumentResult, err error)
	Update(ctx context.Context, databaseName, collectionName string, param UpdateDocumentParams) (result *UpdateDocumentResult, err error)
}

type implementerDocument struct {
	SdkClient
	flat       FlatInterface
	database   *Database
	collection *Collection
}

type UpsertDocumentParams struct {
	BuildIndex *bool
}

type UpsertDocumentResult struct {
	AffectedCount int
}

// Upsert upsert documents into collection. Support for repeated insertion
func (i *implementerDocument) Upsert(ctx context.Context, documents interface{}, params ...*UpsertDocumentParams) (result *UpsertDocumentResult, err error) {
	return i.flat.Upsert(ctx, i.database.DatabaseName, i.collection.CollectionName, documents, params...)
}

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

// Query query the document by document ids.
// The parameters retrieveVector set true, will return the vector field, but will reduce the api speed.
func (i *implementerDocument) Query(ctx context.Context, documentIds []string, params ...*QueryDocumentParams) (*QueryDocumentResult, error) {
	return i.flat.Query(ctx, i.database.DatabaseName, i.collection.CollectionName, documentIds, params...)
}

type SearchDocumentParams struct {
	Filter         *Filter
	Params         *SearchDocParams
	RetrieveVector bool
	OutputFields   []string
	Limit          int64
}

type SearchDocParams struct {
	Nprobe uint32  `json:"nprobe,omitempty"` // 搜索时查找的聚类数量，使用索引默认值即可
	Ef     uint32  `json:"ef,omitempty"`     // HNSW
	Radius float32 `json:"radius,omitempty"` // 距离阈值,范围搜索时有效
}

type SearchDocumentResult struct {
	Warning   string
	Documents [][]Document
}

// Search search document topK by vector. The optional parameters filter will add the filter condition to search.
// The optional parameters hnswParam only be set with the HNSW vector index type.
func (i *implementerDocument) Search(ctx context.Context, vectors [][]float32, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return i.flat.Search(ctx, i.database.DatabaseName, i.collection.CollectionName, vectors, params...)
}

// Search search document topK by document ids. The optional parameters filter will add the filter condition to search.
// The optional parameters hnswParam only be set with the HNSW vector index type.
func (i *implementerDocument) SearchById(ctx context.Context, documentIds []string, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return i.flat.SearchById(ctx, i.database.DatabaseName, i.collection.CollectionName, documentIds, params...)
}

func (i *implementerDocument) SearchByText(ctx context.Context, text map[string][]string, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return i.flat.SearchByText(ctx, i.database.DatabaseName, i.collection.CollectionName, text, params...)
}

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
type RerankOption struct {
	Method    RerankMethod
	FieldList []string
	Weight    []float32
	RrfK      int32
}
type MatchOption struct {
	FieldName string
	Data      interface{}
	Limit     *int
}

type AnnParam struct {
	FieldName string
	Data      interface{}
	Params    *SearchDocParams
	Limit     *int
}

func (i *implementerDocument) HybridSearch(ctx context.Context, params HybridSearchDocumentParams) (*SearchDocumentResult, error) {
	return i.flat.HybridSearch(ctx, i.database.DatabaseName, i.collection.CollectionName, params)
}

type DeleteDocumentParams struct {
	DocumentIds []string
	Filter      *Filter
}

type DeleteDocumentResult struct {
	AffectedCount int
}

// Delete delete document by document ids
func (i *implementerDocument) Delete(ctx context.Context, param DeleteDocumentParams) (result *DeleteDocumentResult, err error) {
	return i.flat.Delete(ctx, i.database.DatabaseName, i.collection.CollectionName, param)
}

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

func (i *implementerDocument) Update(ctx context.Context, param UpdateDocumentParams) (*UpdateDocumentResult, error) {
	return i.flat.Update(ctx, i.database.DatabaseName, i.collection.CollectionName, param)
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

func (i *implementerFlatDocument) Query(ctx context.Context, databaseName, collectionName string, documentIds []string, params ...*QueryDocumentParams) (*QueryDocumentResult, error) {
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

func (i *implementerFlatDocument) Search(ctx context.Context, databaseName, collectionName string,
	vectors [][]float32, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return i.search(ctx, databaseName, collectionName, nil, vectors, nil, params...)
}

func (i *implementerFlatDocument) SearchById(ctx context.Context, databaseName, collectionName string,
	documentIds []string, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return i.search(ctx, databaseName, collectionName, documentIds, nil, nil, params...)
}

func (i *implementerFlatDocument) SearchByText(ctx context.Context, databaseName, collectionName string,
	text map[string][]string, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return i.search(ctx, databaseName, collectionName, nil, nil, text, params...)
}

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

func (i *implementerFlatDocument) Delete(ctx context.Context, databaseName, collectionName string,
	param DeleteDocumentParams) (*DeleteDocumentResult, error) {
	req := new(document.DeleteReq)
	req.Database = databaseName
	req.Collection = collectionName
	req.Query = &document.QueryCond{
		DocumentIds: param.DocumentIds,
		Filter:      param.Filter.Cond(),
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
