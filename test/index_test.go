package test

import (
	"log"
	"testing"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
)

func TestAddIndexWithDefaultParam(t *testing.T) {
	indexPrepare()
	db := cli.Database(database)
	coll := db.Collection(collectionName)

	upsertDataBeforeAddIndex()

	addFilterIndexs := []tcvectordb.FilterIndex{
		{FieldName: "author", FieldType: tcvectordb.String, IndexType: tcvectordb.FILTER}}
	err := cli.AddIndex(ctx, database, collectionName, &tcvectordb.AddIndexParams{FilterIndexs: addFilterIndexs})
	printErr(err)

	time.Sleep(5 * time.Second)

	res, err := db.DescribeCollection(ctx, collectionName)
	printErr(err)
	log.Printf("DescribeCollection result: %+v", ToJson(res))

	option := &tcvectordb.QueryDocumentParams{
		Filter: tcvectordb.NewFilter("author=\"罗贯中\""),
		Limit:  100,
	}
	queryResult, queryErr := coll.Query(ctx, nil, option)
	printErr(queryErr)
	log.Printf("total doc: %d, should be 1", queryResult.Total)
	for _, doc := range queryResult.Documents {
		log.Printf("document: %+v", ToJson(doc))
	}
}

func TestAddIndexNoBuildExistedData(t *testing.T) {
	indexPrepare()
	db := cli.Database(database)
	coll := db.Collection(collectionName)

	upsertDataBeforeAddIndex()

	buildExistedData := false
	addFilterIndexs := []tcvectordb.FilterIndex{
		{FieldName: "author", FieldType: tcvectordb.String, IndexType: tcvectordb.FILTER}}
	err := cli.AddIndex(ctx, database, collectionName, &tcvectordb.AddIndexParams{FilterIndexs: addFilterIndexs,
		BuildExistedData: &buildExistedData})
	printErr(err)

	time.Sleep(5 * time.Second)

	res, err := db.DescribeCollection(ctx, collectionName)
	printErr(err)
	log.Printf("DescribeCollection result: %+v", ToJson(res))

	option := &tcvectordb.QueryDocumentParams{
		Filter: tcvectordb.NewFilter("author=\"罗贯中\""),
		Limit:  100,
	}
	queryResult, queryErr := coll.Query(ctx, nil, option)
	printErr(queryErr)
	log.Printf("total doc: %d, should be 0", queryResult.Total)
	for _, doc := range queryResult.Documents {
		log.Printf("document: %+v", ToJson(doc))
	}
}

func TestAddIndexBuildExistedData(t *testing.T) {
	indexPrepare()
	db := cli.Database(database)
	coll := db.Collection(collectionName)
	upsertDataBeforeAddIndex()

	buildExistedData := true
	addFilterIndexs := []tcvectordb.FilterIndex{
		{FieldName: "author", FieldType: tcvectordb.String, IndexType: tcvectordb.FILTER}}
	err := cli.AddIndex(ctx, database, collectionName, &tcvectordb.AddIndexParams{FilterIndexs: addFilterIndexs,
		BuildExistedData: &buildExistedData})
	printErr(err)

	time.Sleep(5 * time.Second)

	res, err := db.DescribeCollection(ctx, collectionName)
	printErr(err)
	log.Printf("DescribeCollection result: %+v", ToJson(res))

	option := &tcvectordb.QueryDocumentParams{
		Filter: tcvectordb.NewFilter("author=\"罗贯中\""),
		Limit:  100,
	}
	queryResult, queryErr := coll.Query(ctx, nil, option)
	printErr(queryErr)
	log.Printf("total doc: %d, should be 1", queryResult.Total)
	for _, doc := range queryResult.Documents {
		log.Printf("document: %+v", ToJson(doc))
	}
}

func TestAddIndexString(t *testing.T) {
	indexPrepare()
	db := cli.Database(database)
	coll := db.Collection(collectionName)
	upsertDataBeforeAddIndex()

	buildExistedData := true
	addFilterIndexs := []tcvectordb.FilterIndex{
		{FieldName: "author", FieldType: tcvectordb.String, IndexType: tcvectordb.FILTER}}
	err := cli.AddIndex(ctx, database, collectionName, &tcvectordb.AddIndexParams{FilterIndexs: addFilterIndexs,
		BuildExistedData: &buildExistedData})
	printErr(err)

	time.Sleep(5 * time.Second)

	res, err := db.DescribeCollection(ctx, collectionName)
	printErr(err)
	log.Printf("DescribeCollection result: %+v", ToJson(res))

	upsertDataAfterAddIndex()

	option := &tcvectordb.QueryDocumentParams{
		Filter: tcvectordb.NewFilter("author=\"罗贯中\""),
		Limit:  100,
	}
	queryResult, queryErr := coll.Query(ctx, nil, option)
	printErr(queryErr)
	log.Printf("total doc: %d, should be 2", queryResult.Total)
	for _, doc := range queryResult.Documents {
		log.Printf("document: %+v", ToJson(doc))
	}
}

func TestAddIndexWithWrongFilterType(t *testing.T) {
	indexPrepare()
	db := cli.Database(database)
	coll := db.Collection(collectionName)
	upsertDataBeforeAddIndex()

	buildExistedData := true
	addFilterIndexs := []tcvectordb.FilterIndex{
		{FieldName: "author", FieldType: tcvectordb.Uint64, IndexType: tcvectordb.FILTER}}
	err := cli.AddIndex(ctx, database, collectionName, &tcvectordb.AddIndexParams{FilterIndexs: addFilterIndexs,
		BuildExistedData: &buildExistedData})
	printErr(err)

	time.Sleep(5 * time.Second)

	res, err := db.DescribeCollection(ctx, collectionName)
	printErr(err)
	log.Printf("DescribeCollection result: %+v", ToJson(res))

	upsertDataAfterAddIndex()

	option := &tcvectordb.QueryDocumentParams{
		Filter: tcvectordb.NewFilter("author=\"罗贯中\""),
		Limit:  100,
	}
	queryResult, queryErr := coll.Query(ctx, nil, option)
	printErr(queryErr)
	log.Printf("total doc: %d, should be 2", queryResult.Total)
	for _, doc := range queryResult.Documents {
		log.Printf("document: %+v", ToJson(doc))
	}
}

func TestAddIndexUint64(t *testing.T) {
	indexPrepare()
	db := cli.Database(database)
	coll := db.Collection(collectionName)
	upsertDataBeforeAddIndex()

	buildExistedData := true
	addFilterIndexs := []tcvectordb.FilterIndex{
		{FieldName: "page", FieldType: tcvectordb.Uint64, IndexType: tcvectordb.FILTER}}
	err := cli.AddIndex(ctx, database, collectionName, &tcvectordb.AddIndexParams{FilterIndexs: addFilterIndexs,
		BuildExistedData: &buildExistedData})
	printErr(err)

	time.Sleep(5 * time.Second)

	res, err := db.DescribeCollection(ctx, collectionName)
	printErr(err)
	log.Printf("DescribeCollection result: %+v", ToJson(res))

	upsertDataAfterAddIndex()

	option := &tcvectordb.QueryDocumentParams{
		Filter: tcvectordb.NewFilter("page=25"),
		Limit:  100,
	}
	queryResult, queryErr := coll.Query(ctx, nil, option)
	printErr(queryErr)
	log.Printf("total doc: %d, should be 3", queryResult.Total)
	for _, doc := range queryResult.Documents {
		log.Printf("document: %+v", ToJson(doc))
	}
}
func TestAddIndexArray(t *testing.T) {
	indexPrepare()
	db := cli.Database(database)
	coll := db.Collection(collectionName)
	upsertDataBeforeAddIndex()

	buildExistedData := true
	addFilterIndexs := []tcvectordb.FilterIndex{
		{FieldName: "tag", FieldType: tcvectordb.Array, IndexType: tcvectordb.FILTER}}
	err := cli.AddIndex(ctx, database, collectionName, &tcvectordb.AddIndexParams{FilterIndexs: addFilterIndexs,
		BuildExistedData: &buildExistedData})
	printErr(err)

	time.Sleep(5 * time.Second)

	res, err := db.DescribeCollection(ctx, collectionName)
	printErr(err)
	log.Printf("DescribeCollection result: %+v", ToJson(res))

	upsertDataAfterAddIndex()

	option := &tcvectordb.QueryDocumentParams{
		Filter: tcvectordb.NewFilter(tcvectordb.Include("tag", []string{"贾宝玉"})),
		Limit:  100,
	}
	queryResult, queryErr := coll.Query(ctx, nil, option)
	printErr(queryErr)
	log.Printf("total doc: %d, should be 2", queryResult.Total)
	for _, doc := range queryResult.Documents {
		log.Printf("document: %+v", ToJson(doc))
	}
}

func TestParams(t *testing.T) {
	indexPrepare()
	db := cli.Database(database)
	upsertDataBeforeAddIndex()

	err := cli.AddIndex(ctx, database, collectionName)
	printErr(err)

	time.Sleep(5 * time.Second)

	res, err := db.DescribeCollection(ctx, collectionName)
	printErr(err)
	log.Printf("DescribeCollection result: %+v", ToJson(res))
}

func indexPrepare() {
	db, err := cli.CreateDatabaseIfNotExists(ctx, database)
	printErr(err)
	log.Printf("create database if not exists success, %s", db.DatabaseName)

	result, err := db.DropCollection(ctx, collectionName)
	printErr(err)
	log.Printf("drop collection result: %+v", result)

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

func upsertDataBeforeAddIndex() {
	buildIndex := true
	upsertResult, upsertErr := cli.Upsert(ctx, database, collectionName, []map[string]interface{}{
		{
			"id":       "0001",
			"vector":   []float32{0.2123, 0.25, 0.213},
			"bookName": "红楼梦",
			"author":   "曹雪芹",
			"page":     25,
			"segment":  "黛玉葬花",
			"tag":      []string{"贾宝玉", "林黛玉", "薛宝钗"},
		},
		{
			"id":       "0002",
			"vector":   []float32{0.2123, 0.24, 0.213},
			"bookName": "三国演义",
			"author":   "罗贯中",
			"page":     24,
			"segment":  "布大惊，与陈宫商议。宫曰：“闻刘玄德新领徐州，可往投之。”布从其言，竟投徐州来。有人报知玄德。",
			"tag":      []string{"曹操", "诸葛亮", "刘备"},
		},
	}, &tcvectordb.UpsertDocumentParams{BuildIndex: &buildIndex})
	printErr(upsertErr)
	log.Printf("upsert result: %+v", upsertResult)
}

func upsertDataAfterAddIndex() {
	buildIndex := true
	upsertResult, upsertErr := cli.Upsert(ctx, database, collectionName, []map[string]interface{}{
		{
			"id":       "0003",
			"vector":   []float32{0.2123, 0.25, 0.213},
			"bookName": "红楼梦",
			"author":   "曹雪芹",
			"page":     25,
			"segment":  "刘姥姥进大观园",
			"tag":      []string{"刘姥姥", "贾宝玉"},
		},
		{
			"id":       "0004",
			"vector":   []float32{0.2123, 0.25, 0.213},
			"bookName": "三国演义",
			"author":   "罗贯中",
			"page":     25,
			"segment":  "玄德曰：“布乃当今英勇之士，可出迎之。”糜竺曰：“吕布乃虎狼之徒，不可收留；收则伤人矣。",
			"tag":      []string{"曹操", "诸葛亮", "刘备"},
		},
	}, &tcvectordb.UpsertDocumentParams{BuildIndex: &buildIndex})
	printErr(upsertErr)
	log.Printf("upsert result: %+v", upsertResult)
}

func TestBuildIndex(t *testing.T) {
	coll := cli.Database(database).Collection(collectionName)
	// 索引重建，重建期间不支持写入
	_, err := coll.RebuildIndex(ctx, &tcvectordb.RebuildIndexParams{Throttle: 1})
	printErr(err)
}
