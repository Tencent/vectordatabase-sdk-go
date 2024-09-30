package test

import (
	"log"
	"testing"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvdbtext/encoder"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
)

func Test_DropEmbeddingCollectionWithSparseVec(t *testing.T) {
	db := cli.Database(database)

	// 删除collection，删除collection的同时，其中的数据也将被全部删除
	result, err := db.DropCollection(ctx, embedCollWithSparseVec)
	printErr(err)
	log.Printf("drop collection result: %+v", result)
}

func Test_CreateEmbeddingCollectionWithSparseVec(t *testing.T) {
	db := cli.Database(database)

	// 设置embedding字段和模型
	param := &tcvectordb.CreateCollectionParams{
		Embedding: &tcvectordb.Embedding{
			Field:       "segment",
			VectorField: "vector",
			Model:       tcvectordb.BGE_BASE_ZH,
		},
	}

	index := tcvectordb.Indexes{
		// 指定embedding时，vector的维度可以不传，系统会使用embedding model的维度
		VectorIndex: []tcvectordb.VectorIndex{
			{
				FilterIndex: tcvectordb.FilterIndex{
					FieldName: "vector",
					FieldType: tcvectordb.Vector,
					IndexType: tcvectordb.HNSW,
				},
				MetricType: tcvectordb.COSINE,
				Params: &tcvectordb.HNSWParam{
					M:              16,
					EfConstruction: 200,
				},
			},
		},
		SparseVectorIndex: []tcvectordb.SparseVectorIndex{
			{
				FieldName:  "sparse_vector",
				FieldType:  tcvectordb.SparseVector,
				IndexType:  tcvectordb.SPARSE_INVERTED,
				MetricType: tcvectordb.IP,
			},
		},
		FilterIndex: []tcvectordb.FilterIndex{
			{
				FieldName: "id",
				FieldType: tcvectordb.String,
				IndexType: tcvectordb.PRIMARY,
			},
			{
				FieldName: "bookName",
				FieldType: tcvectordb.String,
				IndexType: tcvectordb.FILTER,
			},
			{
				FieldName: "page",
				FieldType: tcvectordb.Uint64,
				IndexType: tcvectordb.FILTER,
			},
		},
	}

	db.WithTimeout(time.Second * 30)
	_, err := db.CreateCollection(ctx, embedCollWithSparseVec, 3, 1, "desription doc", index, param)
	printErr(err)
}

func Test_DescribeEmbeddingCollectionWithSparseVec(t *testing.T) {
	db := cli.Database(database)
	res, err := db.DescribeCollection(ctx, embedCollWithSparseVec)
	printErr(err)
	log.Printf("DescribeCollection result: %+v", res)
}

func Test_ListEmbeddingCollectionWithSparseVec(t *testing.T) {
	db := cli.Database(database)
	res, err := db.ListCollection(ctx)
	printErr(err)
	log.Printf("ListCollection result: %+v", res)
}

func Test_UpsertDocsWithSparseVec(t *testing.T) {
	col := cli.Database(database).Collection(embedCollWithSparseVec)

	bm25, err := encoder.NewBM25Encoder(&encoder.BM25EncoderParams{Bm25Language: "zh"})
	if err != nil {
		log.Fatalf(err.Error())
	}

	segments := []string{
		"富贵功名，前缘分定，为人切莫欺心。",
		"正大光明，忠良善果弥深。些些狂妄天加谴，眼前不遇待时临。",
		"细作探知这个消息，飞报吕布。",
		"布大惊，与陈宫商议。宫曰：“闻刘玄德新领徐州，可往投之。”布从其言，竟投徐州来。有人报知玄德。",
		"玄德曰：“布乃当今英勇之士，可出迎之。”糜竺曰：“吕布乃虎狼之徒，不可收留；收则伤人矣。",
	}
	sparse_vectors, err := bm25.EncodeTexts(segments)
	if err != nil {
		log.Fatalf(err.Error())
	}

	res, err := col.Upsert(ctx, []tcvectordb.Document{
		{
			Id: "0001",
			Fields: map[string]tcvectordb.Field{
				"bookName": {Val: "西游记"},
				"author":   {Val: "吴承恩"},
				"page":     {Val: 21},
				"segment":  {Val: segments[0]},
			},
			SparseVector: sparse_vectors[0],
		},
		{
			Id: "0002",
			Fields: map[string]tcvectordb.Field{
				"bookName": {Val: "西游记"},
				"author":   {Val: "吴承恩"},
				"page":     {Val: 22},
				"segment":  {Val: segments[1]},
			},
			SparseVector: sparse_vectors[1],
		},
		{
			Id: "0003",
			Fields: map[string]tcvectordb.Field{
				"bookName": {Val: "三国演义"},
				"author":   {Val: "罗贯中"},
				"page":     {Val: 23},
				"segment":  {Val: segments[2]},
			},
			SparseVector: sparse_vectors[2],
		},
		{
			Id: "0004",
			Fields: map[string]tcvectordb.Field{
				"bookName": {Val: "三国演义"},
				"author":   {Val: "罗贯中"},
				"page":     {Val: 24},
				"segment":  {Val: segments[3]},
			},
			SparseVector: sparse_vectors[3],
		},
		{
			Id: "0005",
			Fields: map[string]tcvectordb.Field{
				"bookName": {Val: "三国演义"},
				"author":   {Val: "罗贯中"},
				"page":     {Val: 25},
				"segment":  {Val: segments[4]},
			},
			SparseVector: sparse_vectors[4],
		},
	}, nil)

	printErr(err)
	log.Printf("upsert result: %+v", res)
}

func TestUpsertJsonsWithSparseVec(t *testing.T) {
	col := cli.Database(database).Collection(embedCollWithSparseVec)

	bm25, err := encoder.NewBM25Encoder(&encoder.BM25EncoderParams{Bm25Language: "zh"})
	if err != nil {
		log.Fatalf(err.Error())
	}

	segments := []string{
		"布大惊，与陈宫商议。",
	}
	sparse_vectors, err := bm25.EncodeTexts(segments)
	if err != nil {
		log.Fatalf(err.Error())
	}

	println(ToJson(sparse_vectors))

	buildIndex := true
	result, err := col.Upsert(ctx, []map[string]interface{}{
		{
			"id":       "0006",
			"bookName": "三国演义",
			"author":   "罗贯中",
			"page":     24,
			"segment":  segments[0],
			"tag":      []string{"曹操", "诸葛亮", "刘备"},
			"sparse_vector": [][]interface{}{
				{4056151844, 0.7627807},
				{4293601821, 0.7627807},
				{1834252564, 0.7627807},
			},
		},
	}, &tcvectordb.UpsertDocumentParams{BuildIndex: &buildIndex})

	printErr(err)
	log.Printf("upsert result: %+v", result)
}

func TestQueryWithSparseVec(t *testing.T) {
	col := cli.Database(database).Collection(embedCollWithSparseVec)
	option := &tcvectordb.QueryDocumentParams{
		// Filter: tcvectordb.NewFilter(tcvectordb.Include("tag", []string{"曹操", "刘备"})),
		OutputFields: []string{"id", "sparse_vector", "segment"},
		// RetrieveVector: true,
		Limit: 1000,
	}
	documentId := []string{"0001", "0002", "0003", "0004", "0005", "0006", "0007"}
	result, err := col.Query(ctx, documentId, option)
	printErr(err)
	log.Printf("total doc: %d", result.Total)
	for _, doc := range result.Documents {
		log.Printf("document: %+v", ToJson(doc))
	}
}

func TestQueryDirectlyWithSparseVec(t *testing.T) {
	option := &tcvectordb.QueryDocumentParams{
		// Filter: tcvectordb.NewFilter(tcvectordb.Include("tag", []string{"曹操", "刘备"})),
		OutputFields: []string{"id", "sparse_vector", "segment"},
		// RetrieveVector: true,
		Limit: 10,
	}
	documentId := []string{"0001", "0002", "0003", "0004", "0005", "0006", "0007"}
	result, err := cli.Query(ctx, database, embedCollWithSparseVec, documentId, option)
	printErr(err)
	log.Printf("total doc: %d", result.Total)
	for _, doc := range result.Documents {
		log.Printf("document: %+v", ToJson(doc))
	}
}

func TestUpdateWithSparseVec(t *testing.T) {
	col := cli.Database(database).Collection(embedCollWithSparseVec)

	bm25, err := encoder.NewBM25Encoder(&encoder.BM25EncoderParams{Bm25Language: "zh"})
	if err != nil {
		log.Fatalf(err.Error())
	}

	segments := []string{
		"布大惊，与陈宫商议。宫曰：“闻刘玄德新领徐州，可往投之。”布从其言，竟投徐州来。有人报知玄德。",
	}
	sparse_vectors, err := bm25.EncodeTexts(segments)
	if err != nil {
		log.Fatalf(err.Error())
	}

	result, err := col.Update(ctx, tcvectordb.UpdateDocumentParams{
		QueryIds: []string{"0001"},
		//QueryFilter: tcvectordb.NewFilter(`bookName="三国演义"`),
		UpdateFields: map[string]tcvectordb.Field{
			"segment": {Val: segments[0]},
		},
		UpdateSparseVec: sparse_vectors[0],
	})
	printErr(err)
	log.Printf("affect count: %d", result.AffectedCount)
	docs, err := col.Query(ctx, []string{"0001"}, &tcvectordb.QueryDocumentParams{
		OutputFields: []string{"id", "sparse_vector", "segment"},
		Limit:        10,
	})
	printErr(err)
	for _, doc := range docs.Documents {
		log.Printf("query document: %+v", doc.SparseVector)
	}
}

func TestHybridSearchWithSparseVec(t *testing.T) {
	col := cli.Database(database).Collection(embedCollWithSparseVec)

	bm25, err := encoder.NewBM25Encoder(&encoder.BM25EncoderParams{Bm25Language: "zh"})
	if err != nil {
		log.Fatalf(err.Error())
	}

	annSearch := &tcvectordb.AnnParam{
		Data: make([]float32, 768),
	}

	sparseVec, err := bm25.EncodeQuery("刘玄德")
	if err != nil {
		log.Fatalf(err.Error())
	}

	keywordSearch := &tcvectordb.MatchOption{
		Data: sparseVec,
	}

	limit := 3
	searchRes, err := col.HybridSearch(ctx, tcvectordb.HybridSearchDocumentParams{
		AnnParams: []*tcvectordb.AnnParam{annSearch},
		Match:     []*tcvectordb.MatchOption{keywordSearch},
		Rerank: &tcvectordb.RerankOption{
			Method:    "weighted",
			FieldList: []string{"vector", "sparse_vector"},
			Weight:    []float32{0.6, 0.4},
		},
		Limit:        &limit,
		OutputFields: []string{"id", "sparse_vector", "segment"},
	})
	printErr(err)
	log.Printf("search by vector-----------------")
	for i, docs := range searchRes.Documents {
		log.Printf("doc %d result: ", i)
		for _, doc := range docs {
			log.Printf("document: %+v", doc)
		}
	}
}
