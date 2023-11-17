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
	"log"
	"testing"
	"time"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb"
)

var (
	cli                 *entity.VectorDBClient
	ctx                 = context.Background()
	database            = "go-sdk-test-db"
	collectionName      = "go-sdk-test-coll"
	collectionAlias     = "go-sdk-test-alias"
	embeddingCollection = "go-sdk-test-emcoll"
)

func init() {
	// 初始化客户端
	var err error
	// cli, err = tcvectordb.NewClient("http://21.0.83.204:8100", "root", "RPo223wN2yXyUq16dmHcGyzXHaYfWCZWNMGwBC01", &entity.ClientOption{Timeout: 10 * time.Second})
	cli, err = tcvectordb.NewClient("http://lb-3fuz86n6-e8g7tor5zvbql29p.clb.ap-guangzhou.tencentclb.com:60000", "root", "Nfg5r1geFnuuR1uvkxaHqFjoWZwsm9FGr4I28NTK", &entity.ClientOption{Timeout: 10 * time.Second})
	// cli, err = tcvectordb.NewClient("http://21.0.83.222:8100", "root", "NGr06gxdHdS7U6QGOxYvQuEI5VscqoGHadEAvK45", &entity.ClientOption{Timeout: 10 * time.Second})
	// cli, err = tcvectordb.NewClient("http://lb-9wwd95re-yrqlnz5gavcf0b4j.clb.ap-guangzhou.tencentclb.com:20000", "root", "OYBR9chH4fC7f3RF8kKEImEtvuGCFrBhGMWFlZjI", &entity.ClientOption{Timeout: 10 * time.Second})
	if err != nil {
		panic(err)
	}
}

func TestDropDatabase(t *testing.T) {
	_, err := cli.DropDatabase(ctx, database)
	printErr(err)
}

func TestCreateDatabase(t *testing.T) {
	db, err := cli.CreateDatabase(ctx, database)
	printErr(err)
	t.Logf("create database success, %s", db.DatabaseName)
	dbList, err := cli.ListDatabase(ctx)
	printErr(err)

	for _, db := range dbList.Databases {
		t.Logf("database: %s, createTime: %s, dbType: %s", db.DatabaseName, db.Info.CreateTime, db.Info.DbType)
	}
}

func TestListDatabase(t *testing.T) {
	dbList, err := cli.ListDatabase(ctx)
	printErr(err)
	t.Logf("base database ======================")
	for _, db := range dbList.Databases {
		t.Logf("database: %s, createTime: %s, dbType: %s", db.DatabaseName, db.Info.CreateTime, db.Info.DbType)
	}

	t.Logf("AI database ======================")
	for _, db := range dbList.AIDatabases {
		t.Logf("database: %s, createTime: %s, dbType: %s", db.DatabaseName, db.Info.CreateTime, db.Info.DbType)
	}
}

func TestDropCollection(t *testing.T) {
	db := cli.Database(database)

	// 删除collection，删除collection的同时，其中的数据也将被全部删除
	{
		result, err := db.DropCollection(ctx, collectionName)
		printErr(err)
		t.Logf("drop collection result: %+v", result)
	}
}

func TestCreateCollection(t *testing.T) {
	db := cli.Database(database)

	index := entity.Indexes{
		VectorIndex: []entity.VectorIndex{
			{
				FilterIndex: entity.FilterIndex{
					FieldName: "vector",
					FieldType: entity.Vector,
					IndexType: entity.HNSW,
				},
				Dimension:  3,
				MetricType: entity.COSINE,
				Params: &entity.HNSWParam{
					M:              16,
					EfConstruction: 200,
				},
			},
		},
		FilterIndex: []entity.FilterIndex{
			{
				FieldName: "id",
				FieldType: entity.String,
				IndexType: entity.PRIMARY,
			},
			{
				FieldName: "bookName",
				FieldType: entity.String,
				IndexType: entity.FILTER,
			},
			{
				FieldName: "page",
				FieldType: entity.Uint64,
				IndexType: entity.FILTER,
			},
		},
	}

	db.WithTimeout(time.Second * 30)
	_, err := db.CreateCollection(ctx, collectionName, 1, 1, "test collection", index)
	printErr(err)

	// 列出所有 Collection
	result, err := db.ListCollection(ctx, nil)
	printErr(err)
	for _, col := range result.Collections {
		t.Logf("%+v", col)
	}

	// 设置 Collection 的 alias
	_, err = db.SetAlias(ctx, collectionName, collectionAlias, nil)
	printErr(err)

	// 查看 Collection 信息
	colRes, err := db.DescribeCollection(ctx, collectionName)
	printErr(err)
	t.Logf("%+v", colRes)

	// 删除 Collection 的 alias
	db.DeleteAlias(ctx, collectionAlias, nil)
}

func TestListCollection(t *testing.T) {
	db := cli.Database(database)
	// 列出所有 Collection
	result, err := db.ListCollection(ctx, nil)
	printErr(err)
	for _, col := range result.Collections {
		t.Logf("%+v", col)
	}
}

func TestCreateCollectionWithEmbedding(t *testing.T) {
	db := cli.Database(database)

	// 设置embedding字段和模型
	option := &entity.CreateCollectionOption{
		Embedding: &entity.Embedding{
			Field:       "segment",
			VectorField: "vector",
			Model:       entity.BGE_BASE_ZH,
		},
	}

	index := entity.Indexes{
		// 指定embedding时，vector的维度可以不传，系统会使用embedding model的维度
		VectorIndex: []entity.VectorIndex{
			{
				FilterIndex: entity.FilterIndex{
					FieldName: "vector",
					FieldType: entity.Vector,
					IndexType: entity.HNSW,
				},
				MetricType: entity.COSINE,
				Params: &entity.HNSWParam{
					M:              16,
					EfConstruction: 200,
				},
			},
		},
		FilterIndex: []entity.FilterIndex{
			{
				FieldName: "id",
				FieldType: entity.String,
				IndexType: entity.PRIMARY,
			},
			{
				FieldName: "bookName",
				FieldType: entity.String,
				IndexType: entity.FILTER,
			},
			{
				FieldName: "page",
				FieldType: entity.Uint64,
				IndexType: entity.FILTER,
			},
		},
	}

	db.WithTimeout(time.Second * 30)
	db.Debug(true)
	_, err := db.CreateCollection(ctx, embeddingCollection, 1, 1, "desription doc", index, option)
	printErr(err)

	col, err := db.DescribeCollection(ctx, embeddingCollection, nil)
	printErr(err)
	t.Logf("%+v", col)
}

func TestUpsertDocument(t *testing.T) {
	col := cli.Database(database).Collection(collectionName)

	buildIndex := true
	_, err := col.Upsert(ctx, []entity.Document{
		{
			Id:     "0001",
			Vector: []float32{0.2123, 0.21, 0.213},
			Fields: map[string]entity.Field{
				"bookName": {Val: "西游记"},
				"author":   {Val: "吴承恩"},
				"page":     {Val: 21},
				"segment":  {Val: "富贵功名，前缘分定，为人切莫欺心。"},
			},
		},
		{
			Id:     "0002",
			Vector: []float32{0.2123, 0.22, 0.213},
			Fields: map[string]entity.Field{
				"bookName": {Val: "西游记"},
				"author":   {Val: "吴承恩"},
				"page":     {Val: 22},
				"segment":  {Val: "正大光明，忠良善果弥深。些些狂妄天加谴，眼前不遇待时临。"},
			},
		},
		{
			Id:     "0003",
			Vector: []float32{0.2123, 0.23, 0.213},
			Fields: map[string]entity.Field{
				"bookName": {Val: "三国演义"},
				"author":   {Val: "罗贯中"},
				"page":     {Val: 23},
				"segment":  {Val: "细作探知这个消息，飞报吕布。"},
			},
		},
		{
			Id:     "0004",
			Vector: []float32{0.2123, 0.24, 0.213},
			Fields: map[string]entity.Field{
				"bookName": {Val: "三国演义"},
				"author":   {Val: "罗贯中"},
				"page":     {Val: 24},
				"segment":  {Val: "布大惊，与陈宫商议。宫曰：“闻刘玄德新领徐州，可往投之。”布从其言，竟投徐州来。有人报知玄德。"},
			},
		},
		{
			Id:     "0005",
			Vector: []float32{0.2123, 0.25, 0.213},
			Fields: map[string]entity.Field{
				"bookName": {Val: "三国演义"},
				"author":   {Val: "罗贯中"},
				"page":     {Val: 25},
				"segment":  {Val: "玄德曰：“布乃当今英勇之士，可出迎之。”糜竺曰：“吕布乃虎狼之徒，不可收留；收则伤人矣。"},
			},
		},
	}, &entity.UpsertDocumentOption{BuildIndex: &buildIndex})

	printErr(err)
}

func TestQuery(t *testing.T) {
	// 查询
	// 1. query 用于查询数据
	// 2. 可以通过传入主键 id 列表或 filter 实现过滤数据的目的
	// 3. 如果没有主键 id 列表和 filter 则必须传入 limit 和 offset，类似 scan 的数据扫描功能
	// 4. 如果仅需要部分 field 的数据，可以指定 output_fields 用于指定返回数据包含哪些 field，不指定默认全部返回

	col := cli.Database(database).Collection(collectionName)
	col.Debug(true)
	option := &entity.QueryDocumentOption{
		// Filter: entity.NewFilter(`bookName="三国演义"`),
		// OutputFields:   []string{"id", "bookName"},
		// RetrieveVector: true,
		// Limit: 100,
	}
	col.Debug(true)
	result, err := col.Query(ctx, []string{"0001", "0002", "0003", "0004", "0005"}, option)
	printErr(err)
	t.Logf("total doc: %d", result.Total)
	for _, doc := range result.Documents {
		t.Logf("document: %+v", doc)
	}
}

func TestSearch(t *testing.T) {
	col := cli.Database(database).Collection(collectionName)

	// searchById
	// 1. searchById 提供按 id 搜索的能力
	// 1. search 提供按照 vector 搜索的能力
	// 2. 支持通过 filter 过滤数据
	// 3. 如果仅需要部分 field 的数据，可以指定 output_fields 用于指定返回数据包含哪些 field，不指定默认全部返回
	// 4. limit 用于限制每个单元搜索条件的条数，如 vector 传入三组向量，limit 为 3，则 limit 限制的是每组向量返回 top 3 的相似度向量

	// 根据主键 id 查找 Top K 个相似性结果，向量数据库会根据ID 查找对应的向量，再根据向量进行TOP K 相似性检索

	filter := entity.NewFilter(`bookName="三国演义"`)
	searchRes, err := col.SearchById(ctx, []string{"0003"}, &entity.SearchDocumentOption{
		Filter:         filter,                           // 过滤获取到结果
		Params:         &entity.SearchDocParams{Ef: 100}, // 若使用HNSW索引，则需要指定参数ef，ef越大，召回率越高，但也会影响检索速度
		RetrieveVector: false,                            // 是否需要返回向量字段，False：不返回，True：返回
		Limit:          2,                                // 指定 Top K 的 K 值
	})
	printErr(err)
	t.Log("SearchById-----------------")
	for i, docs := range searchRes.Documents {
		t.Logf("doc %d result: ", i)
		for _, doc := range docs {
			t.Logf("id: %s, vector: %v, field: %+v", doc.Id, doc.Vector, doc.Fields)
		}
	}

	// search
	// 1. search 提供按照 vector 搜索的能力
	// 其他选项类似 search 接口

	// 批量相似性查询，根据指定的多个向量查找多个 Top K 个相似性结果
	// 指定检索向量，最多指定20个
	searchRes, err = col.Search(ctx, [][]float32{
		{0.3123, 0.43, 0.213},
		{0.233, 0.12, 0.97},
	}, &entity.SearchDocumentOption{
		Params:         &entity.SearchDocParams{Ef: 100}, // 若使用HNSW索引，则需要指定参数ef，ef越大，召回率越高，但也会影响检索速度
		RetrieveVector: false,                            // 是否需要返回向量字段，False：不返回，True：返回
		Limit:          10,                               // 指定 Top K 的 K 值
	})
	printErr(err)
	t.Logf("search by vector-----------------")
	for i, docs := range searchRes.Documents {
		t.Logf("doc %d result: ", i)
		for _, doc := range docs {
			t.Logf("id: %s, vector: %v, field: %+v", doc.Id, doc.Vector, doc.Fields)
		}
	}
}

func TestUpdateAndDelete(t *testing.T) {
	col := cli.Database(database).Collection(collectionName)

	// update
	// 1. update 提供基于 [主键查询] 和 [Filter 过滤] 的部分字段更新或者非索引字段新增

	// filter 限制仅会更新 id = "0003"
	result, err := col.Update(ctx, &entity.UpdateDocumentOption{
		QueryIds:    []string{"0001", "0003"},
		QueryFilter: entity.NewFilter(`bookName="三国演义"`),
		UpdateFields: map[string]entity.Field{
			"page": {Val: 24},
		},
	})
	printErr(err)
	t.Logf("affect count: %d", result.AffectedCount)
	docs, err := col.Query(ctx, []string{"0003"}, nil)
	printErr(err)
	for _, doc := range docs.Documents {
		t.Logf("query document is: %+v", doc.Fields)
	}

	// delete
	// 1. delete 提供基于 [主键查询] 和 [Filter 过滤] 的数据删除能力
	// 2. 删除功能会受限于 collection 的索引类型，部分索引类型不支持删除操作

	// filter 限制只会删除 id="0001" 成功
	col.Delete(ctx, &entity.DeleteDocumentOption{
		DocumentIds: []string{"0001", "0003"},
		Filter:      entity.NewFilter(`bookName="西游记"`),
	})
}

func TestBuildIndex(t *testing.T) {
	coll := cli.Database(database).Collection(collectionName)
	// 索引重建，重建期间不支持写入
	_, err := coll.RebuildIndex(ctx, &entity.RebuildIndexOption{Throttle: 1})
	printErr(err)
}

func TestTruncateCollection(t *testing.T) {
	db := cli.Database(database)
	// 清空 Collection
	_, err := db.TruncateCollection(ctx, collectionName, nil)
	printErr(err)
}

func TestUpsertEmbedding(t *testing.T) {
	col := cli.Database(database).Collection(embeddingCollection)

	_, err := col.Upsert(ctx, []entity.Document{
		{
			Id: "0001",
			Fields: map[string]entity.Field{
				"bookName": {Val: "西游记"},
				"author":   {Val: "吴承恩"},
				"page":     {Val: 21},
				"segment":  {Val: "富贵功名，前缘分定，为人切莫欺心。"},
			},
		},
		{
			Id: "0002",
			Fields: map[string]entity.Field{
				"bookName": {Val: "西游记"},
				"author":   {Val: "吴承恩"},
				"page":     {Val: 22},
				"segment":  {Val: "正大光明，忠良善果弥深。些些狂妄天加谴，眼前不遇待时临。"},
			},
		},
		{
			Id: "0003",
			Fields: map[string]entity.Field{
				"bookName": {Val: "三国演义"},
				"author":   {Val: "罗贯中"},
				"page":     {Val: 23},
				"segment":  {Val: "细作探知这个消息，飞报吕布。"},
			},
		},
		{
			Id: "0004",
			Fields: map[string]entity.Field{
				"bookName": {Val: "三国演义"},
				"author":   {Val: "罗贯中"},
				"page":     {Val: 24},
				"segment":  {Val: "布大惊，与陈宫商议。宫曰：“闻刘玄德新领徐州，可往投之。”布从其言，竟投徐州来。有人报知玄德。"},
			},
		},
		{
			Id: "0005",
			Fields: map[string]entity.Field{
				"bookName": {Val: "三国演义"},
				"author":   {Val: "罗贯中"},
				"page":     {Val: 25},
				"segment":  {Val: "玄德曰：“布乃当今英勇之士，可出迎之。”糜竺曰：“吕布乃虎狼之徒，不可收留；收则伤人矣。"},
			},
		},
	}, nil)

	printErr(err)
}

func TestQueryEmbedding(t *testing.T) {
	col := cli.Database(database).Collection(embeddingCollection)

	option := &entity.QueryDocumentOption{
		Filter:         entity.NewFilter(`bookName="三国演义"`),
		OutputFields:   []string{"id", "bookName", "segment"},
		RetrieveVector: false,
		Limit:          2,
		Offset:         1,
	}
	col.Debug(true)
	docs, err := col.Query(ctx, []string{"0001", "0002", "0003"}, option)
	printErr(err)
	t.Logf("total doc: %d", docs.Total)
	for _, doc := range docs.Documents {
		t.Logf("%+v", doc)
	}
}

func TestSearchEmbedding(t *testing.T) {
	col := cli.Database(database).Collection(embeddingCollection)
	filter := entity.NewFilter(`bookName="三国演义"`)
	searchRes, err := col.SearchById(ctx, []string{"0003"}, &entity.SearchDocumentOption{
		Filter:         filter,                           // 过滤获取到结果
		Params:         &entity.SearchDocParams{Ef: 100}, // 若使用HNSW索引，则需要指定参数ef，ef越大，召回率越高，但也会影响检索速度
		RetrieveVector: false,                            // 是否需要返回向量字段，False：不返回，True：返回
		Limit:          2,                                // 指定 Top K 的 K 值
	})
	printErr(err)
	t.Log("SearchById-----------------")
	for i, docs := range searchRes.Documents {
		t.Logf("doc %d result: ", i)
		for _, doc := range docs {
			t.Logf("id: %s, vector: %v, field: %+v", doc.Id, doc.Vector, doc.Fields)
		}
	}

	// searchByText
	// 1. searchByText 提供按 文本 搜索的能力，需要开启embedding服务，传入embedding字段的文本值
	// 2. 支持通过 filter 过滤数据
	// 3. 如果仅需要部分 field 的数据，可以指定 output_fields 用于指定返回数据包含哪些 field，不指定默认全部返回
	// 4. limit 用于限制每个单元搜索条件的条数，如 vector 传入三组向量，limit 为 3，则 limit 限制的是每组向量返回 top 3 的相似度向量

	searchRes, err = col.SearchByText(ctx, map[string][]string{"segment": {"吕布"}}, &entity.SearchDocumentOption{
		Params:         &entity.SearchDocParams{Ef: 100}, // 若使用HNSW索引，则需要指定参数ef，ef越大，召回率越高，但也会影响检索速度
		RetrieveVector: false,                            // 是否需要返回向量字段，False：不返回，True：返回
		Limit:          2,                                // 指定 Top K 的 K 值
	})
	printErr(err)
	t.Log("searchByText-----------------")
	for i, docs := range searchRes.Documents {
		t.Logf("doc %d result: ", i)
		for _, doc := range docs {
			t.Logf("id: %s, vector: %v, field: %v", doc.Id, doc.Vector, doc.Fields)
		}
	}
}

func printErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
