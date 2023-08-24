package engine

import (
	"context"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/internal/engine/api/document"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/internal/proto"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/model"
)

type implementerDocument struct {
	model.SdkClient
	databaseName   string
	collectionName string
}

func (i *implementerDocument) Upsert(ctx context.Context, documents []model.Document, buidIndex bool) (err error) {
	req := new(document.UpsertReq)
	req.Database = i.databaseName
	req.Collection = i.collectionName
	req.BuildIndex = buidIndex
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

	res := new(document.UpsertRes)
	err = i.Request(ctx, req, res)
	return
}

func (i *implementerDocument) Query(ctx context.Context, documentIds []string, retrieveVector bool) ([]model.Document, error) {
	req := new(document.QueryReq)
	req.Database = i.databaseName
	req.Collection = i.collectionName
	req.Query = &proto.QueryCond{
		DocumentIds:    documentIds,
		RetrieveVector: retrieveVector,
	}
	res := new(document.QueryRes)
	err := i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}
	var documents []model.Document
	for _, doc := range res.Documents {
		var d model.Document
		d.Id = doc.Id
		d.Vector = doc.Vector
		d.Fields = make(map[string]model.Field)

		for n, v := range doc.Fields {
			d.Fields[n] = model.Field{Val: v}
		}
		documents = append(documents, d)
	}
	return documents, nil
}

func (i *implementerDocument) Search(ctx context.Context, vectors [][]float32, filter *model.Filter, hnswParam *model.HNSWParam, retrieveVector bool, limit int) ([][]model.Document, error) {
	req := new(document.SearchReq)
	req.Database = i.databaseName
	req.Collection = i.collectionName
	req.Search = new(document.SearchCond)
	req.Search.Vectors = vectors
	if filter != nil {
		req.Search.Filter = filter.Cond()
	}
	req.Search.RetrieveVector = retrieveVector
	req.Search.Limit = uint32(limit)
	if hnswParam != nil {
		req.Search.Params = &proto.SearchParams{
			Ef: hnswParam.EfConstruction,
		}
	}

	res := new(document.SearchRes)
	err := i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}
	var documents [][]model.Document
	for _, result := range res.Documents {
		var vecDoc []model.Document
		for _, doc := range result {
			d := model.Document{
				Id:     doc.Id,
				Vector: doc.Vector,
				Score:  doc.Score,
				Fields: make(map[string]model.Field),
			}
			for n, v := range doc.Fields {
				d.Fields[n] = model.Field{Val: v}
			}
			vecDoc = append(vecDoc, d)
		}
		documents = append(documents, vecDoc)
	}

	return documents, nil
}

func (i *implementerDocument) SearchById(ctx context.Context, documentIds []string, filter *model.Filter, hnswParam *model.HNSWParam, retrieveVector bool, limit int) ([][]model.Document, error) {
	req := new(document.SearchReq)
	req.Database = i.databaseName
	req.Collection = i.collectionName
	req.Search = new(document.SearchCond)
	req.Search.DocumentIds = documentIds
	if filter != nil {
		req.Search.Filter = filter.Cond()
	}
	req.Search.RetrieveVector = retrieveVector
	req.Search.Limit = uint32(limit)
	if hnswParam != nil {
		req.Search.Params = &proto.SearchParams{
			Ef: hnswParam.EfConstruction,
		}
	}

	res := new(document.SearchRes)
	err := i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}
	var documents [][]model.Document
	for _, result := range res.Documents {
		var vecDoc []model.Document
		for _, doc := range result {
			d := model.Document{
				Id:     doc.Id,
				Vector: doc.Vector,
				Score:  doc.Score,
				Fields: make(map[string]model.Field),
			}
			for n, v := range doc.Fields {
				d.Fields[n] = model.Field{Val: v}
			}
			vecDoc = append(vecDoc, d)
		}
		documents = append(documents, vecDoc)
	}

	return documents, nil
}

func (i *implementerDocument) Delete(ctx context.Context, documentIds []string) (err error) {
	req := new(document.DeleteReq)
	req.Database = i.databaseName
	req.Collection = i.collectionName
	req.Query = &proto.QueryCond{
		DocumentIds: documentIds,
	}

	res := new(document.DeleteRes)
	err = i.Request(ctx, req, res)
	return
}
