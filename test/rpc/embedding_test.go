package test

import (
	"log"
	"testing"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
)

func TestCreateCollectionWithEmbedding(t *testing.T) {
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
	_, err := db.CreateCollection(ctx, embeddingCollection, 1, 0, "desription doc", index, param)
	printErr(err)

	col, err := db.DescribeCollection(ctx, embeddingCollection)
	printErr(err)
	log.Printf("%+v", col)
}

func TestUpsertEmbedding(t *testing.T) {
	col := cli.Database(database).Collection(embeddingCollection)

	buildIndex := true
	res, err := col.Upsert(ctx, []tcvectordb.Document{
		{
			Id: "0001",
			Fields: map[string]tcvectordb.Field{
				"bookName": {Val: "西游记"},
				"author":   {Val: "吴承恩"},
				"page":     {Val: 21},
				"segment":  {Val: "富贵功名，前缘分定，为人切莫欺心。"},
			},
		},
		{
			Id: "0002",
			Fields: map[string]tcvectordb.Field{
				"bookName": {Val: "西游记"},
				"author":   {Val: "吴承恩"},
				"page":     {Val: 22},
				"segment":  {Val: "正大光明，忠良善果弥深。些些狂妄天加谴，眼前不遇待时临。"},
			},
		},
		{
			Id: "0003",
			Fields: map[string]tcvectordb.Field{
				"bookName": {Val: "三国演义"},
				"author":   {Val: "罗贯中"},
				"page":     {Val: 23},
				"segment":  {Val: "细作探知这个消息，飞报吕布。"},
			},
		},
		{
			Id: "0004",
			Fields: map[string]tcvectordb.Field{
				"bookName": {Val: "三国演义"},
				"author":   {Val: "罗贯中"},
				"page":     {Val: 24},
				"segment":  {Val: "布大惊，与陈宫商议。宫曰：“闻刘玄德新领徐州，可往投之。”布从其言，竟投徐州来。有人报知玄德。"},
			},
		},
		{
			Id: "0005",
			Fields: map[string]tcvectordb.Field{
				"bookName": {Val: "三国演义"},
				"author":   {Val: "罗贯中"},
				"page":     {Val: 25},
				"segment":  {Val: "玄德曰：“布乃当今英勇之士，可出迎之。”糜竺曰：“吕布乃虎狼之徒，不可收留；收则伤人矣。"},
			},
		},
	}, &tcvectordb.UpsertDocumentParams{BuildIndex: &buildIndex})

	printErr(err)
	log.Printf("upsert result: %+v", res)
}

func TestQueryEmbedding(t *testing.T) {
	col := cli.Database(database).Collection(embeddingCollection)

	param := &tcvectordb.QueryDocumentParams{
		Filter:         tcvectordb.NewFilter(`bookName="三国演义"`),
		OutputFields:   []string{"id", "bookName", "segment"},
		RetrieveVector: false,
		Limit:          2,
		Offset:         1,
	}
	docs, err := col.Query(ctx, nil, param)
	printErr(err)
	log.Printf("total doc: %d", docs.Total)
	for _, doc := range docs.Documents {
		log.Printf("%+v", doc)
	}
}

func TestSearchEmbedding(t *testing.T) {
	col := cli.Database(database).Collection(embeddingCollection)

	searchRes, err := col.SearchByText(ctx, map[string][]string{"segment": {"吕布"}}, &tcvectordb.SearchDocumentParams{
		Params:         &tcvectordb.SearchDocParams{Ef: 100}, // 若使用HNSW索引，则需要指定参数ef，ef越大，召回率越高，但也会影响检索速度
		RetrieveVector: false,                                // 是否需要返回向量字段，False：不返回，True：返回
		Limit:          2,                                    // 指定 Top K 的 K 值
	})
	printErr(err)
	t.Log("searchByText-----------------")
	for i, docs := range searchRes.Documents {
		log.Printf("doc %d result: ", i)
		for _, doc := range docs {
			log.Printf("document: %+v", doc)
		}
	}
}
