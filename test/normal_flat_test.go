package test

import (
	"log"
	"testing"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
)

func TestDropFlatCaseDatabase(t *testing.T) {
	result, err := cli.DropDatabase(ctx, database)
	printErr(err)
	log.Printf("DropDatabase result: %+v", result)
}

func TestCreateFlatCaseDatabase(t *testing.T) {
	db, err := cli.CreateDatabase(ctx, database)
	printErr(err)
	log.Printf("create database success, %s", db.DatabaseName)
}

func TestCreateFlatCaseCollection(t *testing.T) {
	db := cli.Database(database)

	index := tcvectordb.Indexes{
		VectorIndex: []tcvectordb.VectorIndex{
			{
				FilterIndex: tcvectordb.FilterIndex{
					FieldName: "vector",
					FieldType: tcvectordb.Vector,
					IndexType: tcvectordb.HNSW,
				},
				Dimension:  3,
				MetricType: tcvectordb.COSINE,
				Params: &tcvectordb.HNSWParam{
					M:              16,
					EfConstruction: 200,
				},
			},
		},
		FilterIndex: []tcvectordb.FilterIndex{
			{FieldName: "id", FieldType: tcvectordb.String, IndexType: tcvectordb.PRIMARY},
			{FieldName: "bookName", FieldType: tcvectordb.String, IndexType: tcvectordb.FILTER},
			{FieldName: "page", FieldType: tcvectordb.Uint64, IndexType: tcvectordb.FILTER},
			{FieldName: "tag", FieldType: tcvectordb.Array, IndexType: tcvectordb.FILTER},
			{FieldName: "expire_at", FieldType: tcvectordb.Uint64, IndexType: tcvectordb.FILTER},
		},
	}

	db.WithTimeout(time.Second * 30)
	param := &tcvectordb.CreateCollectionParams{
		TtlConfig: &tcvectordb.TtlConfig{
			Enable:    true,
			TimeField: "expire_at",
		},
	}

	coll, err := db.CreateCollection(ctx, collectionName, 3, 1, "test collection", index, param)
	printErr(err)
	log.Printf("CreateCollection success: %v: %v", coll.DatabaseName, coll.CollectionName)
}

func TestFlatUpsert(t *testing.T) {
	buildIndex := true
	result, err := cli.Upsert(ctx, database, collectionName, []tcvectordb.Document{
		{
			Id:     "0001",
			Vector: []float32{0.2123, 0.21, 0.213},
			Fields: map[string]tcvectordb.Field{
				"bookName": {Val: "西游记"},
				"author":   {Val: "吴承恩"},
				"page":     {Val: 21},
				"segment":  {Val: "富贵功名，前缘分定，为人切莫欺心。"},
				"tag":      {Val: []string{"孙悟空", "猪八戒", "唐僧"}},
			},
		},
		{
			Id:     "0002",
			Vector: []float32{0.2123, 0.22, 0.213},
			Fields: map[string]tcvectordb.Field{
				"bookName": {Val: "西游记"},
				"author":   {Val: "吴承恩"},
				"page":     {Val: 22},
				"segment":  {Val: "正大光明，忠良善果弥深。些些狂妄天加谴，眼前不遇待时临。"},
				"tag":      {Val: []string{"孙悟空", "猪八戒", "唐僧"}},
			},
		},
		{
			Id:     "0003",
			Vector: []float32{0.2123, 0.23, 0.213},
			Fields: map[string]tcvectordb.Field{
				"bookName": {Val: "三国演义"},
				"author":   {Val: "罗贯中"},
				"page":     {Val: 23},
				"segment":  {Val: "细作探知这个消息，飞报吕布。"},
				"tag":      {Val: []string{"曹操", "诸葛亮", "刘备"}},
			},
		},
		{
			Id:     "0004",
			Vector: []float32{0.2123, 0.24, 0.213},
			Fields: map[string]tcvectordb.Field{
				"bookName": {Val: "三国演义"},
				"author":   {Val: "罗贯中"},
				"page":     {Val: 24},
				"segment":  {Val: "布大惊，与陈宫商议。宫曰：“闻刘玄德新领徐州，可往投之。”布从其言，竟投徐州来。有人报知玄德。"},
				"tag":      {Val: []string{"曹操", "诸葛亮", "刘备"}},
			},
		},
		{
			Id:     "0005",
			Vector: []float32{0.2123, 0.25, 0.213},
			Fields: map[string]tcvectordb.Field{
				"bookName": {Val: "三国演义"},
				"author":   {Val: "罗贯中"},
				"page":     {Val: 25},
				"segment":  {Val: "玄德曰：“布乃当今英勇之士，可出迎之。”糜竺曰：“吕布乃虎狼之徒，不可收留；收则伤人矣。"},
				"tag":      {Val: []string{"曹操", "诸葛亮", "刘备"}},
			},
		},
	}, &tcvectordb.UpsertDocumentParams{BuildIndex: &buildIndex})

	printErr(err)
	log.Printf("upsert result: %+v", result)
}

func TestFlatUpsertJson(t *testing.T) {

	buildIndex := true
	result, err := cli.Upsert(ctx, database, collectionName, []map[string]interface{}{
		{
			"id":       "11",
			"vector":   []float32{0.2123, 0.25, 0.213},
			"bookName": "三国演义",
			"author":   "罗贯中",
			"page":     25,
			"segment":  "玄德曰：“布乃当今英勇之士，可出迎之。”糜竺曰：“吕布乃虎狼之徒，不可收留；收则伤人矣。",
			"tag":      []string{"曹操", "诸葛亮", "刘备"},
		},
		{
			"id":       "12",
			"vector":   []float32{0.2123, 0.24, 0.213},
			"bookName": "三国演义",
			"author":   "罗贯中",
			"page":     24,
			"segment":  "布大惊，与陈宫商议。宫曰：“闻刘玄德新领徐州，可往投之。”布从其言，竟投徐州来。有人报知玄德。",
			"tag":      []string{"曹操", "诸葛亮", "刘备"},
		},
	}, &tcvectordb.UpsertDocumentParams{BuildIndex: &buildIndex})

	printErr(err)
	log.Printf("upsert result: %+v", result)
}

func TestFlatQuery(t *testing.T) {
	option := &tcvectordb.QueryDocumentParams{
		// Filter: tcvectordb.NewFilter(tcvectordb.Include("tag", []string{"曹操", "刘备"})),
		// OutputFields:   []string{"id", "bookName"},
		// RetrieveVector: true,
		Limit: 100,
	}
	// documentId := []string{"0001", "0002", "0003", "0004", "0005"}
	result, err := cli.Query(ctx, database, collectionName, nil, option)
	printErr(err)
	log.Printf("total doc: %d", result.Total)
	for _, doc := range result.Documents {
		log.Printf("document: %+v", doc)
	}
}

func TestFlatSearchById(t *testing.T) {
	filter := tcvectordb.NewFilter(`bookName="三国演义"`)
	documentId := []string{"0003"}
	searchRes, err := cli.SearchById(ctx, database, collectionName, documentId, &tcvectordb.SearchDocumentParams{
		Filter:         filter,
		Params:         &tcvectordb.SearchDocParams{Ef: 100},
		RetrieveVector: false,
		Limit:          2,
	})
	printErr(err)
	t.Log("SearchById-----------------")
	for i, docs := range searchRes.Documents {
		log.Printf("doc %d result: ", i)
		for _, doc := range docs {
			log.Printf("document: %+v", doc)
		}
	}
}

func TestFlatUpdate(t *testing.T) {
	result, err := cli.Update(ctx, database, collectionName, tcvectordb.UpdateDocumentParams{
		QueryIds:    []string{"0001", "0003"},
		QueryFilter: tcvectordb.NewFilter(`bookName="三国演义"`),
		UpdateFields: map[string]tcvectordb.Field{
			"page": {Val: 24},
		},
	})
	printErr(err)
	log.Printf("affect count: %d", result.AffectedCount)
	docs, err := cli.Query(ctx, database, collectionName, []string{"0003"})
	printErr(err)
	for _, doc := range docs.Documents {
		log.Printf("query document: %+v", doc)
	}
}

func TestFlatUpdateJson(t *testing.T) {
	docs, err := cli.Query(ctx, database, collectionName, []string{"0003"})
	printErr(err)
	for _, doc := range docs.Documents {
		log.Printf("before updating, query document: %+v", ToJson(doc))
	}

	result, err := cli.Update(ctx, database, collectionName, tcvectordb.UpdateDocumentParams{
		QueryIds:    []string{"0001", "0003"},
		QueryFilter: tcvectordb.NewFilter(`bookName="三国演义"`),
		UpdateFields: map[string]interface{}{
			"page":   24,
			"vector": []float32{0.2123, 0.25, 0.213},
		},
	})
	printErr(err)
	log.Printf("affect count: %d", result.AffectedCount)
	docs, err = cli.Query(ctx, database, collectionName, []string{"0003"})
	printErr(err)
	for _, doc := range docs.Documents {
		log.Printf("after updating, query document: %+v", ToJson(doc))
	}
}

func TestFlatDelete(t *testing.T) {
	res, err := cli.Delete(ctx, database, collectionName, tcvectordb.DeleteDocumentParams{
		DocumentIds: []string{"0001", "0003"},
		Filter:      tcvectordb.NewFilter(`bookName="西游记"`),
	})
	printErr(err)
	log.Printf("Delete result: %+v", res)
}

func TestFlatCount(t *testing.T) {
	res, err := cli.Count(ctx, database, collectionName, tcvectordb.CountDocumentParams{
		CountFilter: tcvectordb.NewFilter(`bookName="西游记"`),
	})
	printErr(err)
	log.Printf("Count result: %+v", res)
}
