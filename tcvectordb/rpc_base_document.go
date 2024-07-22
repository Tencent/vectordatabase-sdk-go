package tcvectordb

import (
	"context"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/olama"
)

type rpcImplementerDocument struct {
	SdkClient
	rpcClient  olama.SearchEngineClient
	database   *Database
	collection *Collection
}

func (r *rpcImplementerDocument) Upsert(ctx context.Context, documents []Document, params ...*UpsertDocumentParams) (*UpsertDocumentResult, error) {
	req := &olama.UpsertRequest{
		Database:   r.database.DatabaseName,
		Collection: r.collection.CollectionName,
	}
	for _, doc := range documents {
		d := &olama.Document{
			Id:     doc.Id,
			Vector: doc.Vector,
			Fields: make(map[string]*olama.Field),
		}
		for k, v := range doc.Fields {
			d.Fields[k] = ConvertField(&v)
		}
		req.Documents = append(req.Documents, d)
	}
	if len(params) != 0 && params[0] != nil {
		param := params[0]
		if param.BuildIndex != nil {
			req.BuildIndex = *param.BuildIndex
		}
	}
	res, err := r.rpcClient.Upsert(ctx, req)
	if err != nil {
		return nil, err
	}
	return &UpsertDocumentResult{AffectedCount: int(res.AffectedCount)}, nil
}

func (r *rpcImplementerDocument) Query(ctx context.Context, documentIds []string, params ...*QueryDocumentParams) (*QueryDocumentResult, error) {
	req := &olama.QueryRequest{
		Database:   r.database.DatabaseName,
		Collection: r.collection.CollectionName,
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
			d.Fields[n] = Field{Val: v}
		}
		documents = append(documents, d)
	}
	result.Documents = documents
	result.AffectedCount = len(documents)
	result.Total = res.Count
	return result, nil
}

func (r *rpcImplementerDocument) Search(ctx context.Context, vectors [][]float32, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return r.search(ctx, nil, vectors, nil, params...)
}

func (r *rpcImplementerDocument) SearchById(ctx context.Context, documentIds []string, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return r.search(ctx, documentIds, nil, nil, params...)
}

func (r *rpcImplementerDocument) SearchByText(ctx context.Context, text map[string][]string, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	return r.search(ctx, nil, nil, text, params...)
}

func (r *rpcImplementerDocument) Delete(ctx context.Context, param DeleteDocumentParams) (*DeleteDocumentResult, error) {
	req := &olama.DeleteRequest{
		Database:   r.database.DatabaseName,
		Collection: r.collection.CollectionName,
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

func (r *rpcImplementerDocument) Update(ctx context.Context, param UpdateDocumentParams) (*UpdateDocumentResult, error) {
	req := &olama.UpdateRequest{
		Database:   r.database.DatabaseName,
		Collection: r.collection.CollectionName,
		Query: &olama.QueryCond{
			DocumentIds: param.QueryIds,
			Filter:      param.QueryFilter.Cond(),
		},
		Update: &olama.Document{
			Vector: param.UpdateVector,
			Fields: make(map[string]*olama.Field),
		},
	}
	for k, v := range param.UpdateFields {
		req.Update.Fields[k] = ConvertField(&v)
	}
	res, err := r.rpcClient.Update(ctx, req)
	if err != nil {
		return nil, err
	}
	return &UpdateDocumentResult{AffectedCount: int(res.AffectedCount)}, nil
}

func (r *rpcImplementerDocument) search(ctx context.Context, documentIds []string, vectors [][]float32, text map[string][]string, params ...*SearchDocumentParams) (*SearchDocumentResult, error) {
	req := &olama.SearchRequest{
		Database:        r.database.DatabaseName,
		Collection:      r.collection.CollectionName,
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
				d.Fields[n] = Field{Val: v}
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
