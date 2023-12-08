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
	"context"
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
)

var (
	cli                 *tcvectordb.Client
	ctx                 = context.Background()
	database            = "go-sdk-test-db"
	collectionName      = "go-sdk-test-coll"
	collectionAlias     = "go-sdk-test-alias"
	embeddingCollection = "go-sdk-test-emcoll"
)

func init() {
	// 初始化客户端
	var err error
	cli, err = tcvectordb.NewClient("", "root", "", &tcvectordb.ClientOption{Timeout: 10 * time.Second})
	if err != nil {
		panic(err)
	}
}

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

func TestListDatabase(t *testing.T) {
	dbList, err := cli.ListDatabase(ctx)
	printErr(err)
	log.Printf("base database ======================")
	for _, db := range dbList.Databases {
		log.Printf("database: %s, createTime: %s, dbType: %s", db.DatabaseName, db.Info.CreateTime, db.Info.DbType)
	}

	log.Printf("AI database ======================")
	for _, db := range dbList.AIDatabases {
		log.Printf("database: %s, createTime: %s, dbType: %s", db.DatabaseName, db.Info.CreateTime, db.Info.DbType)
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
			{FieldName: "id", FieldType: tcvectordb.String, IndexType: tcvectordb.PRIMARY},
			{FieldName: "bookName", FieldType: tcvectordb.String, IndexType: tcvectordb.FILTER},
			{FieldName: "page", FieldType: tcvectordb.Uint64, IndexType: tcvectordb.FILTER},
			{FieldName: "tag", FieldType: tcvectordb.Array, IndexType: tcvectordb.FILTER},
		},
	}

	db.WithTimeout(time.Second * 30)
	coll, err := db.CreateCollection(ctx, collectionName, 1, 0, "test collection", index)
	printErr(err)
	log.Printf("CreateCollection success: %v: %v", coll.DatabaseName, coll.CollectionName)
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
	log.Printf("DescribeCollection result: %+v", res)
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

func TestQuery(t *testing.T) {
	col := cli.Database(database).Collection(collectionName)
	option := &tcvectordb.QueryDocumentParams{
		// Filter: tcvectordb.NewFilter(tcvectordb.Include("tag", []string{"曹操", "刘备"})),
		// OutputFields:   []string{"id", "bookName"},
		// RetrieveVector: true,
		Limit: 100,
	}
	// documentId := []string{"0001", "0002", "0003", "0004", "0005"}
	result, err := col.Query(ctx, nil, option)
	printErr(err)
	log.Printf("total doc: %d", result.Total)
	for _, doc := range result.Documents {
		log.Printf("document: %+v", doc)
	}
}

func TestSearch(t *testing.T) {
	col := cli.Database(database).Collection(collectionName)

	searchRes, err := col.Search(ctx, [][]float32{
		{0.3123, 0.43, 0.213},
		{0.233, 0.12, 0.97},
	}, &tcvectordb.SearchDocumentParams{
		Params:         &tcvectordb.SearchDocParams{Ef: 100},
		RetrieveVector: false,
		Limit:          10,
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

func TestBuildIndex(t *testing.T) {
	coll := cli.Database(database).Collection(collectionName)
	// 索引重建，重建期间不支持写入
	_, err := coll.RebuildIndex(ctx, &tcvectordb.RebuildIndexParams{Throttle: 1})
	printErr(err)
}

func TestTruncateCollection(t *testing.T) {
	db := cli.Database(database)
	// 清空 Collection
	_, err := db.TruncateCollection(ctx, collectionName)
	printErr(err)
}

func printErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
