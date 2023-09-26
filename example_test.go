package main

import (
	"context"
	"log"
	"testing"
	"time"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entry"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/model"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb"
)

var (
	cli                 entry.VectorDBClient
	database            = "book"
	collectionName      = "book_segments"
	collectionAlias     = "book_segments_alias"
	embeddingCollection = "book_segments_em"
)

func init() {
	// 初始化客户端
	var err error
	cli, err = tcvectordb.NewClient("http://21.0.83.131:8100", "root", "214pFZBsOYegRCqwZ0IDkQJLkcazGnR3I8S4mroV", nil)
	if err != nil {
		panic(err)
	}
}

func TestClear(t *testing.T) {
	_, err := cli.DropDatabase(context.Background(), "book", nil)
	printErr(err)
}

func TestDeleteAndDrop(t *testing.T) {
	db := cli.Database(database)

	// 删除collection，删除collection的同时，其中的数据也将被全部删除
	{
		result, err := db.DropCollection(context.Background(), collectionName, nil)
		printErr(err)
		t.Logf("drop collection result: %+v", result)
	}

	// 删除db，db下的所有collection都将被删除
	{
		result, err := cli.DropDatabase(context.Background(), database, nil)
		printErr(err)
		t.Logf("drop collection result: %+v", result)
	}
}

func TestCreateDatabase(t *testing.T) {
	// 创建DB--'book'
	db, err := cli.CreateDatabase(context.Background(), database, nil)
	printErr(err)
	t.Logf("create database success, %s", db.DatabaseName)

	dbList, err := cli.ListDatabase(context.Background(), nil)
	printErr(err)

	for _, db := range dbList {
		t.Logf("database: %s", db.DatabaseName)
	}
}

func TestCreateCollection(t *testing.T) {
	db := cli.Database(database)

	// 创建 Collection

	// 第一步，设计索引（不是设计 Collection 的结构）
	// 1. 【重要的事】向量对应的文本字段不要建立索引，会浪费较大的内存，并且没有任何作用。
	// 2. 【必须的索引】：主键id、向量字段 vector 这两个字段目前是固定且必须的，参考下面的例子；
	// 3. 【其他索引】：检索时需作为条件查询的字段，比如要按书籍的作者进行过滤，这个时候 author 字段就需要建立索引，
	//     否则无法在查询的时候对 author 字段进行过滤，不需要过滤的字段无需加索引，会浪费内存；
	// 4.  向量数据库支持动态 Schema，写入数据时可以写入任何字段，无需提前定义，类似 MongoDB.
	// 5.  例子中创建一个书籍片段的索引，例如书籍片段的信息包括 {id, vector, segment, bookName, author, page},
	//     id 为主键需要全局唯一，segment 为文本片段, vector 字段需要建立向量索引，假如我们在查询的时候要查询指定书籍
	//     名称的内容，这个时候需要对 bookName 建立索引，其他字段没有条件查询的需要，无需建立索引。

	index := model.Indexes{
		VectorIndex: []model.VectorIndex{
			{
				FilterIndex: model.FilterIndex{
					FieldName: "vector",
					FieldType: model.Vector,
					IndexType: model.HNSW,
				},
				Dimension:  3,
				MetricType: model.COSINE,
				Params: &model.HNSWParam{
					M:              16,
					EfConstruction: 200,
				},
			},
		},
		FilterIndex: []model.FilterIndex{
			{
				FieldName: "id",
				FieldType: model.String,
				IndexType: model.PRIMARY,
			},
			{
				FieldName: "bookName",
				FieldType: model.String,
				IndexType: model.FILTER,
			},
			{
				FieldName: "page",
				FieldType: model.Uint64,
				IndexType: model.FILTER,
			},
		},
	}

	// 第二步：创建 Collection
	// 创建collection耗时较长，需要调整客户端的timeout
	db.WithTimeout(time.Second * 30)
	_, err := db.CreateCollection(context.Background(), collectionName, 3, 2, "test collection", index, nil)
	printErr(err)

	// 列出所有 Collection
	colList, err := db.ListCollection(context.Background(), nil)
	printErr(err)
	for _, col := range colList {
		t.Logf("%+v", col)
	}

	// 设置 Collection 的 alias
	_, err = db.SetAlias(context.Background(), collectionName, collectionAlias, nil)
	printErr(err)

	// 查看 Collection 信息
	colRes, err := db.DescribeCollection(context.Background(), collectionName, nil)
	printErr(err)
	t.Logf("%+v", colRes)

	// 删除 Collection 的 alias
	db.DeleteAlias(context.Background(), collectionAlias, nil)
}

func TestCreateCollectionWithEmbedding(t *testing.T) {
	db := cli.Database(database)

	// 设置embedding字段和模型
	option := &entry.CreateCollectionOption{
		Embedding: &model.Embedding{
			Field:       "segment",
			VectorField: "vector",
			Model:       model.BGE_BASE_ZH,
		},
	}

	index := model.Indexes{
		// 指定embedding时，vector的维度可以不传，系统会使用embedding model的维度
		VectorIndex: []model.VectorIndex{
			{
				FilterIndex: model.FilterIndex{
					FieldName: "vector",
					FieldType: model.Vector,
					IndexType: model.HNSW,
				},
				MetricType: model.COSINE,
				Params: &model.HNSWParam{
					M:              16,
					EfConstruction: 200,
				},
			},
		},
		FilterIndex: []model.FilterIndex{
			{
				FieldName: "id",
				FieldType: model.String,
				IndexType: model.PRIMARY,
			},
			{
				FieldName: "bookName",
				FieldType: model.String,
				IndexType: model.FILTER,
			},
			{
				FieldName: "page",
				FieldType: model.Uint64,
				IndexType: model.FILTER,
			},
		},
	}

	db.WithTimeout(time.Second * 30)
	_, err := db.CreateCollection(context.Background(), embeddingCollection, 3, 2, "desription doc", index, option)
	printErr(err)

	col, err := db.DescribeCollection(context.Background(), embeddingCollection, nil)
	printErr(err)
	t.Logf("%+v", col)
}

func TestUpsertDocument(t *testing.T) {
	col := cli.Database(database).Collection(collectionName)

	_, err := col.Upsert(context.Background(), []model.Document{
		{
			Id:     "0001",
			Vector: []float32{0.2123, 0.21, 0.213},
			Fields: map[string]model.Field{
				"bookName": {Val: "西游记"},
				"author":   {Val: "吴承恩"},
				"page":     {Val: 21},
				"segment":  {Val: "富贵功名，前缘分定，为人切莫欺心。"},
			},
		},
		{
			Id:     "0002",
			Vector: []float32{0.2123, 0.22, 0.213},
			Fields: map[string]model.Field{
				"bookName": {Val: "西游记"},
				"author":   {Val: "吴承恩"},
				"page":     {Val: 22},
				"segment":  {Val: "正大光明，忠良善果弥深。些些狂妄天加谴，眼前不遇待时临。"},
			},
		},
		{
			Id:     "0003",
			Vector: []float32{0.2123, 0.23, 0.213},
			Fields: map[string]model.Field{
				"bookName": {Val: "三国演义"},
				"author":   {Val: "罗贯中"},
				"page":     {Val: 23},
				"segment":  {Val: "细作探知这个消息，飞报吕布。"},
			},
		},
		{
			Id:     "0004",
			Vector: []float32{0.2123, 0.24, 0.213},
			Fields: map[string]model.Field{
				"bookName": {Val: "三国演义"},
				"author":   {Val: "罗贯中"},
				"page":     {Val: 24},
				"segment":  {Val: "布大惊，与陈宫商议。宫曰：“闻刘玄德新领徐州，可往投之。”布从其言，竟投徐州来。有人报知玄德。"},
			},
		},
		{
			Id:     "0005",
			Vector: []float32{0.2123, 0.25, 0.213},
			Fields: map[string]model.Field{
				"bookName": {Val: "三国演义"},
				"author":   {Val: "罗贯中"},
				"page":     {Val: 25},
				"segment":  {Val: "玄德曰：“布乃当今英勇之士，可出迎之。”糜竺曰：“吕布乃虎狼之徒，不可收留；收则伤人矣。"},
			},
		},
	}, &entry.UpsertDocumentOption{BuildIndex: true})

	printErr(err)
}

func TestQuery(t *testing.T) {
	// 查询
	// 1. query 用于查询数据
	// 2. 可以通过传入主键 id 列表或 filter 实现过滤数据的目的
	// 3. 如果没有主键 id 列表和 filter 则必须传入 limit 和 offset，类似 scan 的数据扫描功能
	// 4. 如果仅需要部分 field 的数据，可以指定 output_fields 用于指定返回数据包含哪些 field，不指定默认全部返回

	col := cli.Database(database).Collection(collectionName)
	option := &entry.QueryDocumentOption{
		Filter:         model.NewFilter(`bookName="三国演义"`),
		OutputFields:   []string{"id", "bookName"},
		RetrieveVector: true,
		Limit:          2,
		Offset:         1,
	}
	col.Debug(true)
	docs, result, err := col.Query(context.Background(), []string{"0001", "0002", "0003", "0004", "0005"}, option)
	printErr(err)
	t.Logf("total doc: %d", result.Total)
	for _, doc := range docs {
		t.Logf("id: %s, vector: %v, field: %+v", doc.Id, doc.Vector, doc.Fields)
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

	filter := model.NewFilter(`bookName="三国演义"`)
	searchRes, err := col.SearchById(context.Background(), []string{"0003"}, &entry.SearchDocumentOption{
		Filter:         filter,                          // 过滤获取到结果
		Params:         &entry.SearchDocParams{Ef: 100}, // 若使用HNSW索引，则需要指定参数ef，ef越大，召回率越高，但也会影响检索速度
		RetrieveVector: false,                           // 是否需要返回向量字段，False：不返回，True：返回
		Limit:          2,                               // 指定 Top K 的 K 值
	})
	printErr(err)
	t.Log("SearchById-----------------")
	for i, docs := range searchRes {
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
	searchRes, err = col.Search(context.Background(), [][]float32{
		{0.3123, 0.43, 0.213},
		{0.233, 0.12, 0.97},
	}, &entry.SearchDocumentOption{
		Params:         &entry.SearchDocParams{Ef: 100}, // 若使用HNSW索引，则需要指定参数ef，ef越大，召回率越高，但也会影响检索速度
		RetrieveVector: false,                           // 是否需要返回向量字段，False：不返回，True：返回
		Limit:          10,                              // 指定 Top K 的 K 值
	})
	printErr(err)
	t.Logf("search by vector-----------------")
	for i, docs := range searchRes {
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
	result, err := col.Update(context.Background(), &entry.UpdateDocumentOption{
		QueryIds:    []string{"0001", "0003"},
		QueryFilter: model.NewFilter(`bookName="三国演义"`),
		UpdateFields: map[string]model.Field{
			"page": {Val: 24},
		},
	})
	printErr(err)
	t.Logf("affect count: %d", result.AffectedCount)
	docs, _, err := col.Query(context.Background(), []string{"0003"}, nil)
	printErr(err)
	for _, doc := range docs {
		t.Logf("query document is: %+v", doc.Fields)
	}

	// delete
	// 1. delete 提供基于 [主键查询] 和 [Filter 过滤] 的数据删除能力
	// 2. 删除功能会受限于 collection 的索引类型，部分索引类型不支持删除操作

	// filter 限制只会删除 id="0001" 成功
	col.Delete(context.Background(), &entry.DeleteDocumentOption{
		DocumentIds: []string{"0001", "0003"},
		Filter:      model.NewFilter(`bookName="西游记"`),
	})
}

func TestBuildIndex(t *testing.T) {
	db := cli.Database(database)
	// 索引重建，重建期间不支持写入
	_, err := db.IndexRebuild(context.Background(), collectionName, &entry.IndexRebuildOption{Throttle: 1})
	printErr(err)
}

func TestTruncateCollection(t *testing.T) {
	db := cli.Database(database)
	// 清空 Collection
	_, err := db.TruncateCollection(context.Background(), collectionName, nil)
	printErr(err)
}

func TestUpsertEmbedding(t *testing.T) {
	col := cli.Database(database).Collection(embeddingCollection)

	_, err := col.Upsert(context.Background(), []model.Document{
		{
			Id: "0001",
			Fields: map[string]model.Field{
				"bookName": {Val: "西游记"},
				"author":   {Val: "吴承恩"},
				"page":     {Val: 21},
				"segment":  {Val: "富贵功名，前缘分定，为人切莫欺心。"},
			},
		},
		{
			Id: "0002",
			Fields: map[string]model.Field{
				"bookName": {Val: "西游记"},
				"author":   {Val: "吴承恩"},
				"page":     {Val: 22},
				"segment":  {Val: "正大光明，忠良善果弥深。些些狂妄天加谴，眼前不遇待时临。"},
			},
		},
		{
			Id: "0003",
			Fields: map[string]model.Field{
				"bookName": {Val: "三国演义"},
				"author":   {Val: "罗贯中"},
				"page":     {Val: 23},
				"segment":  {Val: "细作探知这个消息，飞报吕布。"},
			},
		},
		{
			Id: "0004",
			Fields: map[string]model.Field{
				"bookName": {Val: "三国演义"},
				"author":   {Val: "罗贯中"},
				"page":     {Val: 24},
				"segment":  {Val: "布大惊，与陈宫商议。宫曰：“闻刘玄德新领徐州，可往投之。”布从其言，竟投徐州来。有人报知玄德。"},
			},
		},
		{
			Id: "0005",
			Fields: map[string]model.Field{
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

	option := &entry.QueryDocumentOption{
		Filter:         model.NewFilter(`bookName="三国演义"`),
		OutputFields:   []string{"id", "bookName"},
		RetrieveVector: false,
		Limit:          2,
		Offset:         1,
	}
	docs, result, err := col.Query(context.Background(), []string{"0001", "0002", "0003"}, option)
	printErr(err)
	t.Logf("total doc: %d", result.Total)
	for _, doc := range docs {
		t.Logf("%+v", doc)
	}
}

func TestSearchEmbedding(t *testing.T) {
	col := cli.Database(database).Collection(embeddingCollection)
	filter := model.NewFilter(`bookName="三国演义"`)
	searchRes, err := col.SearchById(context.Background(), []string{"0003"}, &entry.SearchDocumentOption{
		Filter:         filter,                          // 过滤获取到结果
		Params:         &entry.SearchDocParams{Ef: 100}, // 若使用HNSW索引，则需要指定参数ef，ef越大，召回率越高，但也会影响检索速度
		RetrieveVector: false,                           // 是否需要返回向量字段，False：不返回，True：返回
		Limit:          2,                               // 指定 Top K 的 K 值
	})
	printErr(err)
	t.Log("SearchById-----------------")
	for i, docs := range searchRes {
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

	searchRes, err = col.SearchByText(context.Background(), map[string][]string{"segment": {"吕布"}}, &entry.SearchDocumentOption{
		Params:         &entry.SearchDocParams{Ef: 100}, // 若使用HNSW索引，则需要指定参数ef，ef越大，召回率越高，但也会影响检索速度
		RetrieveVector: false,                           // 是否需要返回向量字段，False：不返回，True：返回
		Limit:          2,                               // 指定 Top K 的 K 值
	})
	printErr(err)
	t.Log("searchByText-----------------")
	for i, docs := range searchRes {
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
