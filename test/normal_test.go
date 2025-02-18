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

package test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/document"
)

func TestDropDatabase(t *testing.T) {
	result, err := cli.DropDatabase(ctx, database)
	printErr(err)
	log.Printf("DropDatabase result: %+v", result)
}

func TestCreateDatabase(t *testing.T) {
	db, err := cli.CreateDatabase(ctx, database)
	printErr(err)
	log.Printf("create database success, %s", db.DatabaseName)
}

func TestExistsDatabase(t *testing.T) {
	dbExists, err := cli.ExistsDatabase(ctx, database)
	printErr(err)
	log.Printf("database %v exists: %v", database, dbExists)
}

func TestCreateDatabaseIfNotExist(t *testing.T) {
	db, err := cli.CreateDatabaseIfNotExists(ctx, database)
	printErr(err)
	log.Printf("create database if not exists success, %s", db.DatabaseName)
}

func TestListDatabase(t *testing.T) {
	dbList, err := cli.ListDatabase(ctx)
	printErr(err)
	log.Printf("base database ======================")
	for _, db := range dbList.Databases {
		log.Printf("database: %s, createTime: %s, dbType: %s, collectionNum: %v", db.DatabaseName,
			db.Info.CreateTime, db.Info.DbType, db.Info.Count)
	}

	log.Printf("AI database ======================")
	for _, db := range dbList.AIDatabases {
		log.Printf("database: %s, createTime: %s, dbType: %s, collectionViewNum: %v", db.DatabaseName,
			db.Info.CreateTime, db.Info.DbType, db.Info.Count)
	}
}

func TestDropCollection(t *testing.T) {
	db := cli.Database(database)

	// 删除collection，删除collection的同时，其中的数据也将被全部删除
	result, err := db.DropCollection(ctx, collectionName)
	printErr(err)
	log.Printf("drop collection result: %+v", result)
}

func TestCreateCollection(t *testing.T) {
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
			{FieldName: "id", FieldType: tcvectordb.String, IndexType: tcvectordb.PRIMARY, AutoId: "uuid"},
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

func TestExistsCollection(t *testing.T) {
	db := cli.Database(database)
	collExists, err := db.ExistsCollection(ctx, collectionName)
	printErr(err)
	log.Printf("collection %v exists: %v", collectionName, collExists)
}

func TestCreateCollectionIfNotExists(t *testing.T) {
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

	coll, err := db.CreateCollectionIfNotExists(ctx, collectionName, 3, 1, "test collection", index, param)
	printErr(err)
	log.Printf("create collection if not exists success: %v: %v", coll.DatabaseName, coll.CollectionName)
}

func TestListCollection(t *testing.T) {
	db := cli.Database(database)
	// 列出所有 Collection
	result, err := db.ListCollection(ctx)
	printErr(err)
	for _, col := range result.Collections {
		data, err := json.Marshal(col)
		printErr(err)
		log.Printf("%+v", string(data))
	}
}

func TestDescribeCollection(t *testing.T) {
	db := cli.Database(database)
	res, err := db.DescribeCollection(ctx, collectionName)
	printErr(err)
	log.Printf("DescribeCollection result: %+v", ToJson(res))
}

func TestUpsert(t *testing.T) {
	col := cli.Database(database).Collection(collectionName)

	buildIndex := true
	result, err := col.Upsert(ctx, []tcvectordb.Document{
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

func TestUpsertJson(t *testing.T) {
	col := cli.Database(database).Collection(collectionName)

	buildIndex := true
	result, err := col.Upsert(ctx, []map[string]interface{}{
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

func TestQuery(t *testing.T) {
	col := cli.Database(database).Collection(collectionName)
	option := &tcvectordb.QueryDocumentParams{
		Filter:       tcvectordb.NewFilter("page>21"),
		OutputFields: []string{"id", "page"},
		// RetrieveVector: true,
		Limit: 2,
		Sort: []document.SortRule{
			{
				FieldName: "page",
				Direction: "asc",
			},
		},
	}
	// documentId := []string{"0001", "0002", "0003", "0004", "0005"}
	result, err := col.Query(ctx, nil, option)
	printErr(err)
	log.Printf("total doc: %d", result.Total)
	for _, doc := range result.Documents {
		log.Printf("document: %+v", ToJson(doc))
	}
}

func TestSearch(t *testing.T) {
	col := cli.Database(database).Collection(collectionName)

	radius := float32(0.76)
	searchRes, err := col.Search(ctx, [][]float32{
		//{0.3123, 0.43, 0.213},
		{0.233, 0.12, 0.97},
	}, &tcvectordb.SearchDocumentParams{
		Params: &tcvectordb.SearchDocParams{
			Ef: 100,
		},
		Radius:         &radius,
		RetrieveVector: false,
		Limit:          2,
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

func TestSearchById(t *testing.T) {
	col := cli.Database(database).Collection(collectionName)

	filter := tcvectordb.NewFilter(`bookName="三国演义"`)
	documentId := []string{"0003"}
	searchRes, err := col.SearchById(ctx, documentId, &tcvectordb.SearchDocumentParams{
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

func TestUpdate(t *testing.T) {
	col := cli.Database(database).Collection(collectionName)

	result, err := col.Update(ctx, tcvectordb.UpdateDocumentParams{
		QueryIds:    []string{"0001", "0003"},
		QueryFilter: tcvectordb.NewFilter(`bookName="三国演义"`),
		UpdateFields: map[string]tcvectordb.Field{
			"page": {Val: 24},
		},
	})
	printErr(err)
	log.Printf("affect count: %d", result.AffectedCount)
	docs, err := col.Query(ctx, []string{"0003"})
	printErr(err)
	for _, doc := range docs.Documents {
		log.Printf("query document: %+v", doc)
	}
}

func TestDelete(t *testing.T) {
	col := cli.Database(database).Collection(collectionName)

	res, err := col.Delete(ctx, tcvectordb.DeleteDocumentParams{
		DocumentIds: []string{"0001", "0003"},
		Filter:      tcvectordb.NewFilter(`bookName="西游记"`),
	})
	printErr(err)
	log.Printf("Delete result: %+v", res)
}

func TestCount(t *testing.T) {
	col := cli.Database(database).Collection(collectionName)

	res, err := col.Count(ctx, tcvectordb.CountDocumentParams{
		CountFilter: tcvectordb.NewFilter(`bookName="西游记"`),
	})
	printErr(err)
	log.Printf("Count result: %+v", res)
}

func TestReupsertCollection(t *testing.T) {
	col := cli.Database(database).Collection(collectionName)
	testLen := int64(10)
	option := &tcvectordb.QueryDocumentParams{
		RetrieveVector: true,
		Limit:          testLen,
	}
	firstQuery, err := col.Query(ctx, nil, option)
	printErr(err)
	if firstQuery.Total == 0 {
		return
	}
	ids := make([]string, firstQuery.Total)
	for i := uint64(0); i < firstQuery.Total; i++ {
		ids[i] = firstQuery.Documents[i].Id
	}
	_, err = col.Delete(ctx, tcvectordb.DeleteDocumentParams{DocumentIds: ids})
	printErr(err)
	secondQuery, err := col.Query(ctx, ids, option)
	printErr(err)
	if secondQuery.Total != 0 {
		t.Fatalf("Delete by id failed")
	}
	_, err = col.Upsert(ctx, firstQuery.Documents)
	printErr(err)
	thirdQuery, err := col.Query(ctx, ids, option)
	printErr(err)
	if thirdQuery.Total == 0 {
		t.Fatalf("reupsert failed")
	}
}

func TestTruncateCollection(t *testing.T) {
	log.Println("wait collection indexes build successfully")
	time.Sleep(5 * time.Second)
	db := cli.Database(database)
	// 清空 Collection
	_, err := db.TruncateCollection(ctx, collectionName)
	printErr(err)
}

func TestJson(t *testing.T) {
	strT := "{\"databaseName\":\"go-sdk-test-db\",\"shardNum\":1846430633467240448}"
	temp := make(map[string]interface{}, 0)
	json.Unmarshal([]byte(strT), &temp)
	if v, ok := temp["shardNum"].(float64); ok {
		println(v, "float")
	}
	if v, ok := temp["shardNum"].(uint64); ok {
		println(v)
	}

	if v, ok := temp["shardNum"].(json.Number); ok {
		println(v)
	}

	temp1 := temp["shardNum"]
	println(fmt.Sprintf("%T, %v", temp1, temp1))
}

func TestEmeb(t *testing.T) {
	model := "model_bge"
	println(tcvectordb.EmbeddingModel(model))
}
