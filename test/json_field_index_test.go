package test

import (
	"log"
	"testing"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
)

func Test_CreateCollectionWithJsonFieldIndex(t *testing.T) {
	db, err := cli.CreateDatabaseIfNotExists(ctx, database)
	printErr(err)
	log.Printf("create database success, %s", db.DatabaseName)

	_, err = db.DropCollection(ctx, collectionName)
	printErr(err)

	index := tcvectordb.Indexes{}
	index.VectorIndex = append(index.VectorIndex, tcvectordb.VectorIndex{
		FilterIndex: tcvectordb.FilterIndex{
			FieldName: "vector",
			FieldType: tcvectordb.Vector,
			IndexType: tcvectordb.HNSW,
		},
		Dimension:  1,
		MetricType: tcvectordb.IP,
		Params: &tcvectordb.HNSWParam{
			M:              16,
			EfConstruction: 200,
		},
	})
	index.FilterIndex = append(index.FilterIndex, tcvectordb.FilterIndex{FieldName: "id", FieldType: tcvectordb.String, IndexType: tcvectordb.PRIMARY, AutoId: "uuid"})
	index.FilterIndex = append(index.FilterIndex, tcvectordb.FilterIndex{FieldName: "field_json", FieldType: tcvectordb.Json, IndexType: tcvectordb.FILTER})

	db.WithTimeout(time.Second * 30)

	params := tcvectordb.CreateCollectionParams{
		FilterIndexConfig: &tcvectordb.FilterIndexConfig{
			FilterAll: false,
			//FieldsWithoutIndex: []string{"field_json"},
		},
	}
	_, err = db.CreateCollection(ctx, collectionName, 3, 1, "test collection", index, &params)
	printErr(err)

	log.Println("------------------------ DescribeCollection ------------------------")
	// 查看 Collection 信息
	colRes, err := db.DescribeCollection(ctx, collectionName)
	printErr(err)
	log.Printf("DescribeCollection: %+v", ToJson(colRes))

}

func UpsertDataWithJsonField() {
	_, err := cli.Database(database).TruncateCollection(ctx, collectionName)
	printErr(err)

	log.Println("------------------------------ Upsert ------------------------------")
	documentList := []tcvectordb.Document{
		{

			Vector: []float32{0},
			Fields: map[string]tcvectordb.Field{
				"field_json": {Val: map[string]interface{}{
					"username": "Alice",
					"age":      28,
					"habit":    []string{"reading", "singing"},
				}},
			},
		},
		{
			Vector: []float32{0.1},
			Fields: map[string]tcvectordb.Field{
				"field_json": {Val: map[string]interface{}{
					"username": "Bob",
					"age":      25,
					"habit":    []string{"writing"},
				}},
			},
		},
		{
			Vector: []float32{0.2},
			Fields: map[string]tcvectordb.Field{
				"field_json": {Val: map[string]interface{}{
					"username": "Charlie",
					"age":      35,
					"habit":    []string{"singing", "drawing"},
				}},
			},
		},
		{
			Vector: []float32{0.3},
			Fields: map[string]tcvectordb.Field{
				"field_json": {Val: map[string]interface{}{
					"username": "David",
					"age":      28,
					"habit":    []string{"dancing", "drawing", "writing"},
				}},
			},
		},
	}
	result, err := cli.Upsert(ctx, database, collectionName, documentList)
	printErr(err)
	log.Printf("UpsertResult: %+v", result)
}

func QueryByFilter(filter *tcvectordb.Filter) {
	result, err := cli.Query(ctx, database, collectionName, []string{}, &tcvectordb.QueryDocumentParams{
		Filter:         filter,
		RetrieveVector: true,
		Limit:          10,
	})
	printErr(err)
	log.Printf("QueryResult: filter: %v, total: %v, affect: %v", filter.Cond(), result.Total, result.AffectedCount)
	for _, doc := range result.Documents {
		log.Printf("QueryDocument: %+v, field: %+v", doc.Id, doc.Fields["field_json"])
	}
}

func DeleteByFilter(filter *tcvectordb.Filter) {
	result, err := cli.Delete(ctx, database, collectionName, tcvectordb.DeleteDocumentParams{
		Filter: filter,
		Limit:  10,
	})
	printErr(err)
	log.Printf("DeleteResult: filter: %v, affect: %v", filter.Cond(), result.AffectedCount)
}

func SearchByFilter(filter *tcvectordb.Filter) {
	result, err := cli.Search(ctx, database, collectionName, [][]float32{{0.1}}, &tcvectordb.SearchDocumentParams{
		Filter: filter,
		Limit:  10,
	})
	printErr(err)
	log.Printf("SearchResult: filter: %v", filter.Cond())
	for _, docs := range result.Documents {
		for _, doc := range docs {
			log.Printf("SearchDocument: %+v, field: %+v", doc.Id, doc.Fields["field_json"])
		}
	}
}

func HybridSearchByFilter(filter *tcvectordb.Filter) {
	annParams := make([]*tcvectordb.AnnParam, 0)
	annParam := &tcvectordb.AnnParam{
		FieldName: "vector",
		Data:      []float32{0.1},
	}
	annParams = append(annParams, annParam)

	limit := 10
	result, err := cli.HybridSearch(ctx, database, collectionName, tcvectordb.HybridSearchDocumentParams{
		AnnParams: annParams,
		Filter:    filter,
		Limit:     &limit,
	})
	printErr(err)
	log.Printf("HybridSearchResult: filter: %v", filter.Cond())
	for _, docs := range result.Documents {
		for _, doc := range docs {
			log.Printf("HybridSearchDocument: %+v, field: %+v", doc.Id, doc.Fields["field_json"])
		}
	}
}

func Test_FilterJsonStringHybridSearch(t *testing.T) {
	UpsertDataWithJsonField()
	log.Println("------------------------------ HybridSearch in ------------------------------")
	filter := tcvectordb.NewFilter(`field_json.username in ("Bob")`)
	HybridSearchByFilter(filter)

	log.Println("------------------------------ Query not in ------------------------------")
	filter = tcvectordb.NewFilter(`field_json.username not in ("Bob")`)
	HybridSearchByFilter(filter)

	log.Println("------------------------------ Query = ------------------------------")
	filter = tcvectordb.NewFilter(`field_json.username="Bob"`)
	HybridSearchByFilter(filter)

	log.Println("------------------------------ Query != ------------------------------")
	filter = tcvectordb.NewFilter(`field_json.username!="Bob"`)
	HybridSearchByFilter(filter)
}

func Test_FilterJsonUint64HybridSearch(t *testing.T) {
	UpsertDataWithJsonField()
	log.Println("------------------------------ Query in ------------------------------")
	filter := tcvectordb.NewFilter(`field_json.age in (28)`)
	HybridSearchByFilter(filter)

	log.Println("------------------------------ Query not in ------------------------------")
	filter = tcvectordb.NewFilter(`field_json.age not in (28)`)
	HybridSearchByFilter(filter)

	log.Println("------------------------------ Query > ------------------------------")
	filter = tcvectordb.NewFilter(`field_json.age > 28`)
	HybridSearchByFilter(filter)

	log.Println("------------------------------ Query >= ------------------------------")
	filter = tcvectordb.NewFilter(`field_json.age >= 28`)
	HybridSearchByFilter(filter)

	log.Println("------------------------------ Query < ------------------------------")
	filter = tcvectordb.NewFilter(`field_json.age < 28`)
	HybridSearchByFilter(filter)

	log.Println("------------------------------ Query <= ------------------------------")
	filter = tcvectordb.NewFilter(`field_json.age <= 28`)
	HybridSearchByFilter(filter)
}

func Test_FilterJsonString(t *testing.T) {
	UpsertDataWithJsonField()
	log.Println("------------------------------ Query in ------------------------------")
	filter := tcvectordb.NewFilter(`field_json.username in ("Bob")`)
	QueryByFilter(filter)
	SearchByFilter(filter)
	DeleteByFilter(filter)
	QueryByFilter(filter)
	SearchByFilter(filter)

	UpsertDataWithJsonField()
	log.Println("------------------------------ Query not in ------------------------------")
	filter = tcvectordb.NewFilter(`field_json.username not in ("Bob")`)
	QueryByFilter(filter)
	SearchByFilter(filter)
	DeleteByFilter(filter)
	QueryByFilter(filter)
	SearchByFilter(filter)

	UpsertDataWithJsonField()
	log.Println("------------------------------ Query = ------------------------------")
	filter = tcvectordb.NewFilter(`field_json.username="Bob"`)
	QueryByFilter(filter)
	SearchByFilter(filter)
	DeleteByFilter(filter)
	QueryByFilter(filter)
	SearchByFilter(filter)

	UpsertDataWithJsonField()
	log.Println("------------------------------ Query != ------------------------------")
	filter = tcvectordb.NewFilter(`field_json.username!="Bob"`)
	QueryByFilter(filter)
	SearchByFilter(filter)
	DeleteByFilter(filter)
	QueryByFilter(filter)
	SearchByFilter(filter)
}

func Test_FilterJsonUint64(t *testing.T) {
	UpsertDataWithJsonField()
	log.Println("------------------------------ Query in ------------------------------")
	filter := tcvectordb.NewFilter(`field_json.age in (28)`)
	QueryByFilter(filter)
	SearchByFilter(filter)
	DeleteByFilter(filter)
	QueryByFilter(filter)
	SearchByFilter(filter)

	UpsertDataWithJsonField()
	log.Println("------------------------------ Query not in ------------------------------")
	filter = tcvectordb.NewFilter(`field_json.age not in (28)`)
	QueryByFilter(filter)
	SearchByFilter(filter)
	DeleteByFilter(filter)
	QueryByFilter(filter)
	SearchByFilter(filter)

	UpsertDataWithJsonField()
	log.Println("------------------------------ Query > ------------------------------")
	filter = tcvectordb.NewFilter(`field_json.age > 28`)
	QueryByFilter(filter)
	SearchByFilter(filter)
	DeleteByFilter(filter)
	QueryByFilter(filter)
	SearchByFilter(filter)

	UpsertDataWithJsonField()
	log.Println("------------------------------ Query >= ------------------------------")
	filter = tcvectordb.NewFilter(`field_json.age >= 28`)
	QueryByFilter(filter)
	SearchByFilter(filter)
	DeleteByFilter(filter)
	QueryByFilter(filter)
	SearchByFilter(filter)

	UpsertDataWithJsonField()
	log.Println("------------------------------ Query < ------------------------------")
	filter = tcvectordb.NewFilter(`field_json.age < 28`)
	QueryByFilter(filter)
	SearchByFilter(filter)
	DeleteByFilter(filter)
	QueryByFilter(filter)
	SearchByFilter(filter)

	UpsertDataWithJsonField()
	log.Println("------------------------------ Query <= ------------------------------")
	filter = tcvectordb.NewFilter(`field_json.age <= 28`)
	QueryByFilter(filter)
	SearchByFilter(filter)
	DeleteByFilter(filter)
	QueryByFilter(filter)
	SearchByFilter(filter)
}

func Test_FilterJsonArray(t *testing.T) {
	UpsertDataWithJsonField()
	log.Println("------------------------------ query include  ------------------------------")
	filter := tcvectordb.NewFilter(`field_json.habit include  ("dancing", "writing")`)
	QueryByFilter(filter)
	SearchByFilter(filter)
	HybridSearchByFilter(filter)

	log.Println("------------------------------ query exclude  ------------------------------")
	filter = tcvectordb.NewFilter(`field_json.habit exclude  ("dancing", "writing")`)
	QueryByFilter(filter)
	SearchByFilter(filter)
	HybridSearchByFilter(filter)

	log.Println("------------------------------ query include all  ------------------------------")
	filter = tcvectordb.NewFilter(`field_json.habit include all  ("dancing", "writing")`)
	QueryByFilter(filter)
	SearchByFilter(filter)
	HybridSearchByFilter(filter)
}

func Test_FilterJsonArrayInclude(t *testing.T) {
	UpsertDataWithJsonField()
	log.Println("------------------------------ delete include ------------------------------")
	filter := tcvectordb.NewFilter(`field_json.habit include ("dancing")`)
	QueryByFilter(filter)
	SearchByFilter(filter)
	DeleteByFilter(filter)
	QueryByFilter(filter)
	SearchByFilter(filter)
}

func Test_FilterJsonArrayExclude(t *testing.T) {
	UpsertDataWithJsonField()
	log.Println("------------------------------ delete exclude ------------------------------")
	filter := tcvectordb.NewFilter(`field_json.habit exclude ("dancing")`)
	QueryByFilter(filter)
	SearchByFilter(filter)
	DeleteByFilter(filter)
	QueryByFilter(filter)
	SearchByFilter(filter)
}

func Test_FilterJsonArrayIncludeAll(t *testing.T) {
	UpsertDataWithJsonField()
	log.Println("------------------------------ delete include all ------------------------------")
	filter := tcvectordb.NewFilter(`field_json.habit include all ("dancing", "writing")`)
	QueryByFilter(filter)
	SearchByFilter(filter)
	DeleteByFilter(filter)
	QueryByFilter(filter)
	SearchByFilter(filter)
}

func Test_FilterJsonAnd(t *testing.T) {
	UpsertDataWithJsonField()
	log.Println("------------------------------ delete a.b and a.c ------------------------------")
	filter := tcvectordb.NewFilter(`field_json.username="Bob"`).And(`field_json.age=25`)
	QueryByFilter(filter)
	SearchByFilter(filter)
	HybridSearchByFilter(filter)
	DeleteByFilter(filter)
	QueryByFilter(filter)
	SearchByFilter(filter)
}

func Test_FilterJsonOr(t *testing.T) {
	UpsertDataWithJsonField()
	log.Println("------------------------------ delete a.b or a.c ------------------------------")
	filter := tcvectordb.NewFilter(`field_json.username="Bob"`).Or(`field_json.age=28`)
	QueryByFilter(filter)
	SearchByFilter(filter)
	HybridSearchByFilter(filter)
	DeleteByFilter(filter)
	QueryByFilter(filter)
	SearchByFilter(filter)
}

func Test_FilterJsonNot(t *testing.T) {
	UpsertDataWithJsonField()
	log.Println("------------------------------ delete a.b and not a.c ------------------------------")
	filter := tcvectordb.NewFilter(`field_json.username="Bob"`).AndNot(`field_json.age=28`)
	QueryByFilter(filter)
	SearchByFilter(filter)
	HybridSearchByFilter(filter)
	DeleteByFilter(filter)
	QueryByFilter(filter)
	SearchByFilter(filter)

	UpsertDataWithJsonField()
	log.Println("------------------------------ delete a.b or not a.c ------------------------------")
	filter = tcvectordb.NewFilter(`field_json.username="Bob"`).OrNot(`field_json.age=28`)
	QueryByFilter(filter)
	SearchByFilter(filter)
	HybridSearchByFilter(filter)
	DeleteByFilter(filter)
	QueryByFilter(filter)
	SearchByFilter(filter)
}

func Test_FieldJsonUpdate(t *testing.T) {

	res, err := cli.Update(ctx, database, collectionName, tcvectordb.UpdateDocumentParams{
		QueryIds: []string{"E893FBD2-8B90-2404-B983-C944353FE7AE"},
		UpdateFields: map[string]tcvectordb.Field{
			"field_json": {Val: map[string]interface{}{
				"name":    "John",
				"age":     70,
				"hobbies": []string{"reading", "sports", "music"},
			}},
		},
	})
	printErr(err)
	log.Printf("UpdateResult: affect: %v", res.AffectedCount)
}
