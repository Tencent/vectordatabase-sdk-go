package tcvectordb

import (
	"context"
	"fmt"

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

func (r *rpcImplementerDocument) Upsert(ctx context.Context, documents interface{}, params ...*UpsertDocumentParams) (*UpsertDocumentResult, error) {
	return r.flat.Upsert(ctx, r.database.DatabaseName, r.collection.CollectionName, documents, params...)
}

func (r *rpcImplementerDocument) Query(ctx context.Context, documentIds []string, params ...*QueryDocumentParams) (*QueryDocumentResult, error) {
	return r.flat.Query(ctx, r.database.DatabaseName, r.collection.CollectionName, documentIds, params...)
}

func (r *rpcImplementerDocument) Search(ctx context.Context, vectors [][]float32, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return r.flat.Search(ctx, r.database.DatabaseName, r.collection.CollectionName, vectors, params...)
}

func (r *rpcImplementerDocument) SearchById(ctx context.Context, documentIds []string, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return r.flat.SearchById(ctx, r.database.DatabaseName, r.collection.CollectionName, documentIds, params...)
}

func (r *rpcImplementerDocument) SearchByText(ctx context.Context, text map[string][]string, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return r.flat.SearchByText(ctx, r.database.DatabaseName, r.collection.CollectionName, text, params...)
}

func (r *rpcImplementerDocument) Delete(ctx context.Context, param DeleteDocumentParams) (*DeleteDocumentResult, error) {
	return r.flat.Delete(ctx, r.database.DatabaseName, r.collection.CollectionName, param)
}

func (r *rpcImplementerDocument) Update(ctx context.Context, param UpdateDocumentParams) (*UpdateDocumentResult, error) {
	return r.flat.Update(ctx, r.database.DatabaseName, r.collection.CollectionName, param)
}

type rpcImplementerFlatDocument struct {
	SdkClient
	rpcClient olama.SearchEngineClient
}

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

func (r *rpcImplementerFlatDocument) Search(ctx context.Context, databaseName, collectionName string,
	vectors [][]float32, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return r.search(ctx, databaseName, collectionName, nil, vectors, nil, params...)
}

func (r *rpcImplementerFlatDocument) SearchById(ctx context.Context, databaseName, collectionName string,
	documentIds []string, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return r.search(ctx, databaseName, collectionName, documentIds, nil, nil, params...)
}

func (r *rpcImplementerFlatDocument) SearchByText(ctx context.Context, databaseName, collectionName string,
	text map[string][]string, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return r.search(ctx, databaseName, collectionName, nil, nil, text, params...)
}

func (r *rpcImplementerFlatDocument) Delete(ctx context.Context, databaseName, collectionName string,
	param DeleteDocumentParams) (*DeleteDocumentResult, error) {
	req := &olama.DeleteRequest{
		Database:   databaseName,
		Collection: collectionName,
		Query: &olama.QueryCond{
			DocumentIds: param.DocumentIds,
			Filter:      param.Filter.Cond(),
		},
	}
	res, err := r.rpcClient.Dele(ctx, req)
	if err != nil {
		return nil, err
	}
	return &DeleteDocumentResult{AffectedCount: int(res.AffectedCount)}, nil
}

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
				Radius: param.Params.Radius,
			}
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
