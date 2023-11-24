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

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb/api/document"
)

var _ DocumentInterface = &implementerDocument{}

// DocumentInterface document api
type DocumentInterface interface {
	SdkClient
	Upsert(ctx context.Context, documents []Document, params ...*UpsertDocumentParams) (result *UpsertDocumentResult, err error)
	Query(ctx context.Context, documentIds []string, params ...*QueryDocumentParams) (result *QueryDocumentResult, err error)
	Search(ctx context.Context, vectors [][]float32, params ...*SearchDocumentParams) (result *SearchDocumentResult, err error)
	SearchById(ctx context.Context, documentIds []string, params ...*SearchDocumentParams) (result *SearchDocumentResult, err error)
	SearchByText(ctx context.Context, text map[string][]string, params ...*SearchDocumentParams) (result *SearchDocumentResult, err error)
	Delete(ctx context.Context, param DeleteDocumentParams) (result *DeleteDocumentResult, err error)
	Update(ctx context.Context, param UpdateDocumentParams) (result *UpdateDocumentResult, err error)
}

type implementerDocument struct {
	SdkClient
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
func (i *implementerDocument) Upsert(ctx context.Context, documents []Document, params ...*UpsertDocumentParams) (result *UpsertDocumentResult, err error) {
	req := new(document.UpsertReq)
	req.Database = i.database.DatabaseName
	req.Collection = i.collection.CollectionName
	for _, doc := range documents {
		d := &document.Document{}
		d.Id = doc.Id
		d.Vector = doc.Vector
		d.Fields = make(map[string]interface{})
		for k, v := range doc.Fields {
			d.Fields[k] = v.Val
		}
		req.Documents = append(req.Documents, d)
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
	req := new(document.QueryReq)
	req.Database = i.database.DatabaseName
	req.Collection = i.collection.CollectionName
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
	Documents [][]Document
}

// Search search document topK by vector. The optional parameters filter will add the filter condition to search.
// The optional parameters hnswParam only be set with the HNSW vector index type.
func (i *implementerDocument) Search(ctx context.Context, vectors [][]float32, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return i.search(ctx, nil, vectors, nil, params...)
}

// Search search document topK by document ids. The optional parameters filter will add the filter condition to search.
// The optional parameters hnswParam only be set with the HNSW vector index type.
func (i *implementerDocument) SearchById(ctx context.Context, documentIds []string, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return i.search(ctx, documentIds, nil, nil, params...)
}

func (i *implementerDocument) SearchByText(ctx context.Context, text map[string][]string, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return i.search(ctx, nil, nil, text, params...)
}

func (i *implementerDocument) search(ctx context.Context, documentIds []string, vectors [][]float32, text map[string][]string, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	req := new(document.SearchReq)
	req.Database = i.database.DatabaseName
	req.Collection = i.collection.CollectionName
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
	result.Documents = documents
	return result, nil
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
	req := new(document.DeleteReq)
	req.Database = i.database.DatabaseName
	req.Collection = i.collection.CollectionName
	req.Query = &document.QueryCond{
		DocumentIds: param.DocumentIds,
		Filter:      param.Filter.Cond(),
	}

	res := new(document.DeleteRes)
	result = new(DeleteDocumentResult)
	err = i.Request(ctx, req, res)
	if err != nil {
		return
	}
	result.AffectedCount = res.AffectedCount
	return
}

type UpdateDocumentParams struct {
	QueryIds     []string
	QueryFilter  *Filter
	UpdateVector []float32
	UpdateFields map[string]Field
}

type UpdateDocumentResult struct {
	AffectedCount int
}

func (i *implementerDocument) Update(ctx context.Context, param UpdateDocumentParams) (*UpdateDocumentResult, error) {
	req := new(document.UpdateReq)
	req.Database = i.database.DatabaseName
	req.Collection = i.collection.CollectionName
	req.Query = new(document.QueryCond)

	req.Query.DocumentIds = param.QueryIds
	req.Query.Filter = param.QueryFilter.Cond()
	req.Update.Vector = param.UpdateVector
	req.Update.Fields = make(map[string]interface{})
	for k, v := range param.UpdateFields {
		req.Update.Fields[k] = v.Val
	}

	res := new(document.UpdateRes)
	result := new(UpdateDocumentResult)
	err := i.Request(ctx, req, res)
	if err != nil {
		return result, err
	}
	result.AffectedCount = int(res.AffectedCount)
	return result, nil
}

type Document struct {
	Id     string
	Vector []float32
	// omitempty when upsert
	Score  float32 `json:"_,omitempty"`
	Fields map[string]Field
}
