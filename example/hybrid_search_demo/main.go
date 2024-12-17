package main

import (
	"context"
	"log"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvdbtext/encoder"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
)

type Demo struct {
	client *tcvectordb.Client
}

func NewDemo(url, username, key string) (*Demo, error) {
	// cli, err := tcvectordb.NewRpcClient(url, username, key, &tcvectordb.ClientOption{
	// 	ReadConsistency: tcvectordb.EventualConsistency})
	cli, err := tcvectordb.NewClient(url, username, key, &tcvectordb.ClientOption{
		ReadConsistency: tcvectordb.EventualConsistency})
	if err != nil {
		return nil, err
	}
	// disable/enable http request log print
	// cli.Debug(false)
	return &Demo{client: cli}, nil
}

func (d *Demo) Clear(ctx context.Context, database string) error {
	log.Println("--------------------------- DropDatabase ---------------------------")
	result, err := d.client.DropDatabase(ctx, database)
	if err != nil {
		return err
	}
	log.Printf("drop database result: %+v", result)
	return nil
}

func (d *Demo) DeleteAndDrop(ctx context.Context, database, collection string) error {
	// 删除collection，删除collection的同时，其中的数据也将被全部删除
	log.Println("-------------------------- DropCollection --------------------------")
	colDropResult, err := d.client.Database(database).DropCollection(ctx, collection)
	if err != nil {
		return err
	}
	log.Printf("drop collection result: %+v", colDropResult)

	log.Println("--------------------------- DropDatabase ---------------------------")
	// 删除db，db下的所有collection都将被删除
	dbDropResult, err := d.client.DropDatabase(ctx, database)
	if err != nil {
		return err
	}
	log.Printf("drop database result: %+v", dbDropResult)
	return nil
}

func (d *Demo) CreateDBAndCollection(ctx context.Context, database, collection, alias string) error {
	// 创建DB--'book'
	log.Println("-------------------------- CreateDatabase --------------------------")
	db, err := d.client.CreateDatabase(ctx, database)
	if err != nil {
		return err
	}

	log.Println("--------------------------- ListDatabase ---------------------------")
	dbList, err := d.client.ListDatabase(ctx)
	if err != nil {
		return err
	}
	for _, db := range dbList.Databases {
		log.Printf("database: %s", db.DatabaseName)
	}

	log.Println("------------------------- CreateCollection -------------------------")
	// 创建 Collection

	// 第一步，设计索引（不是设计 Collection 的结构）
	// 1. 【重要的事】向量对应的文本字段不要建立索引，会浪费较大的内存，并且没有任何作用。
	// 2. 【必须的索引】：主键id、向量字段 vector 这两个字段目前是固定且必须的，参考下面的例子；如果使用稀疏向量，需要创建稀疏向量对应的索引
	// 3. 【其他索引】：检索时需作为条件查询的字段，比如要按书籍的作者进行过滤，这个时候 author 字段就需要建立索引，
	//     否则无法在查询的时候对 author 字段进行过滤，不需要过滤的字段无需加索引，会浪费内存；
	// 4.  向量数据库支持动态 Schema，写入数据时可以写入任何字段，无需提前定义，类似 MongoDB.
	// 5.  例子中创建一个书籍片段的索引，例如书籍片段的信息包括 {id, vector, segment, bookName, author, page},
	//     id 为主键需要全局唯一，segment 为文本片段, vector 字段需要建立向量索引，假如我们在查询的时候要查询指定书籍
	//     名称的内容，这个时候需要对 bookName 建立索引，其他字段没有条件查询的需要，无需建立索引。
	index := tcvectordb.Indexes{}
	index.VectorIndex = append(index.VectorIndex, tcvectordb.VectorIndex{
		FilterIndex: tcvectordb.FilterIndex{
			FieldName: "vector",
			FieldType: tcvectordb.Vector,
			IndexType: tcvectordb.HNSW,
		},
		Dimension:  3,
		MetricType: tcvectordb.IP,
		Params: &tcvectordb.HNSWParam{
			M:              16,
			EfConstruction: 200,
		},
	})
	index.SparseVectorIndex = append(index.SparseVectorIndex, tcvectordb.SparseVectorIndex{
		FieldName:  "sparse_vector",
		FieldType:  tcvectordb.SparseVector,
		IndexType:  tcvectordb.SPARSE_INVERTED,
		MetricType: tcvectordb.IP,
	})

	index.FilterIndex = append(index.FilterIndex, tcvectordb.FilterIndex{FieldName: "id", FieldType: tcvectordb.String, IndexType: tcvectordb.PRIMARY})

	// 第二步：创建 Collection
	// 创建collection耗时较长，需要调整客户端的timeout
	// 这里以三可用区实例作为参考，具体实例不同的规格所支持的shard和replicas区间不同，需要参考官方文档
	db.WithTimeout(time.Second * 30)
	_, err = db.CreateCollection(ctx, collection, 3, 1, "test collection", index)
	if err != nil {
		return err
	}

	log.Println("-------------------------- ListCollection --------------------------")
	// 列出所有 Collection
	collListRes, err := db.ListCollection(ctx)
	if err != nil {
		return err
	}
	for _, col := range collListRes.Collections {
		log.Printf("ListCollection: %+v", col)
	}

	log.Println("----------------------------- SetAlias -----------------------------")
	// 设置 Collection 的 alias
	_, err = db.SetAlias(ctx, collection, alias)
	if err != nil {
		return err
	}

	log.Println("------------------------ DescribeCollection ------------------------")
	// 查看 Collection 信息
	colRes, err := db.DescribeCollection(ctx, collection)
	if err != nil {
		return err
	}
	log.Printf("DescribeCollection: %+v", colRes)

	log.Println("---------------------------- DeleteAlias ---------------------------")
	// 删除 Collection 的 alias
	delAliasRes, err := db.DeleteAlias(ctx, alias)
	if err != nil {
		return err
	}
	log.Printf("DeleteAliasResult: %v", delAliasRes)
	return nil
}

func (d *Demo) UpsertData(ctx context.Context, database, collection string) error {
	// 获取 Collection 对象
	coll := d.client.Database(database).Collection(collection)

	log.Println("------------------------------ Upsert ------------------------------")

	documentList := []tcvectordb.Document{
		{
			Id:     "0001",
			Vector: []float32{0.2123, 0.21, 0.213},
			SparseVector: []encoder.SparseVecItem{
				{TermId: 1172076521, Score: 0.71296215},
				{TermId: 3434399993, Score: 0.71296215},
			},
		},
		{
			Id:     "0002",
			Vector: []float32{0.2123, 0.22, 0.213},
			SparseVector: []encoder.SparseVecItem{
				{TermId: 225043997, Score: 0.71296215},
				{TermId: 256049420, Score: 0.71296215},
			},
		},
		{
			Id:     "0003",
			Vector: []float32{0.2123, 0.23, 0.213},
			SparseVector: []encoder.SparseVecItem{
				{TermId: 4162843804, Score: 0.71296215},
				{TermId: 256049420, Score: 0.71296215},
			},
		},
		{
			Id:     "0004",
			Vector: []float32{0.2123, 0.24, 0.213},
			SparseVector: []encoder.SparseVecItem{
				{TermId: 3656406568, Score: 0.71296215},
				{TermId: 256049420, Score: 0.71296215},
			},
		},
		{
			Id:     "0005",
			Vector: []float32{0.2123, 0.25, 0.213},
			SparseVector: []encoder.SparseVecItem{
				{TermId: 4162843804, Score: 0.71296215},
				{TermId: 3656406568, Score: 0.71296215},
			},
		},
	}

	result, err := coll.Upsert(ctx, documentList)
	if err != nil {
		return err
	}
	log.Printf("UpsertResult: %+v", result)
	return nil
}

func (d *Demo) QueryData(ctx context.Context, database, collection string) error {
	// 获取 Collection 对象
	coll := d.client.Database(database).Collection(collection)

	log.Println("------------------------------ Query ------------------------------")
	// 查询
	// 1. query 用于查询数据
	// 2. 可以通过传入主键 id 列表或 filter 实现过滤数据的目的
	// 3. 如果没有主键 id 列表和 filter 则必须传入 limit 和 offset，类似 scan 的数据扫描功能
	// 4. 如果仅需要部分 field 的数据，可以指定 output_fields 用于指定返回数据包含哪些 field，不指定默认全部返回
	documentIds := []string{"0000", "0001", "0002", "0003", "0004"}
	outputField := []string{"id", "sparse_vector"}

	result, err := coll.Query(ctx, documentIds, &tcvectordb.QueryDocumentParams{
		RetrieveVector: false,
		OutputFields:   outputField,
		Limit:          5,
		Offset:         0,
	})
	if err != nil {
		return err
	}
	log.Printf("QueryResult: total: %v, affect: %v", result.Total, result.AffectedCount)
	for _, doc := range result.Documents {
		log.Printf("QueryDocument: %+v", doc)
	}

	log.Println("------------------------------ hybridSearch ------------------------------")

	annSearch := &tcvectordb.AnnParam{
		FieldName: "vector",
		Data:      []float32{0.2123, 0.21, 0.213},
	}

	keywordSearch := &tcvectordb.MatchOption{
		FieldName: "sparse_vector",
		Data: []encoder.SparseVecItem{
			{TermId: 1172076521, Score: 0.71296215},
			{TermId: 3434399993, Score: 0.71296215},
		},
	}

	limit := 2
	searchRes, err := coll.HybridSearch(ctx, tcvectordb.HybridSearchDocumentParams{
		AnnParams: []*tcvectordb.AnnParam{annSearch},
		Match:     []*tcvectordb.MatchOption{keywordSearch},
		// rerank也支持rrf，使用方式见下
		// Rerank: &tcvectordb.RerankOption{
		// 	Method:    tcvectordb.RerankRrf,
		// 	RrfK: 1,
		// },
		Rerank: &tcvectordb.RerankOption{
			Method:    tcvectordb.RerankWeighted,
			FieldList: []string{"vector", "sparse_vector"},
			Weight:    []float32{0.1, 0.9},
		},
		Limit:        &limit,
		OutputFields: []string{"id", "sparse_vector"},
	})
	if err != nil {
		return err
	}

	// 输出相似性检索结果，检索结果为二维数组，每一位为一组返回结果，分别对应search时指定的多个向量
	for i, item := range searchRes.Documents {
		log.Printf("HybridSearchDocumentResult, index: %d ==================", i)
		for _, doc := range item {
			log.Printf("HybridSearchDocument: %+v", doc)
		}
	}
	return nil
}

func (d *Demo) UpdateAndDeleteCollection(ctx context.Context, database, collection string) error {
	// 获取 Collection 对象
	db := d.client.Database(database)
	coll := db.Collection(collection)

	log.Println("------------------------------ Update ------------------------------")

	documentId := []string{"0002"}
	bm25, err := encoder.NewBM25Encoder(&encoder.BM25EncoderParams{Bm25Language: "zh"})
	if err != nil {
		log.Fatalf(err.Error())
	}

	segments := []string{
		"腾讯云向量数据库（Tencent Cloud VectorDB）是一款全托管的自研企业级分布式数据库服务，专用于存储、索引、检索、管理由深度神经网络或其他机器学习模型生成的大量多维嵌入向量。",
	}

	sparse_vectors, err := bm25.EncodeTexts(segments)
	if err != nil {
		log.Fatalf(err.Error())
	}

	result, err := coll.Update(ctx, tcvectordb.UpdateDocumentParams{
		QueryIds:        documentId,
		UpdateSparseVec: sparse_vectors[0],
	})
	if err != nil {
		return err
	}
	log.Printf("UpdateResult: %+v", result)

	log.Println("------------------------ TruncateCollection ------------------------")
	// truncate_collection
	// 清空 Collection
	time.Sleep(time.Second * 5)
	truncateRes, err := db.TruncateCollection(ctx, collection)
	if err != nil {
		return err
	}
	log.Printf("TruncateResult: %+v", truncateRes)
	return nil
}

func printErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	database := "go-sdk-demo-db"
	collectionName := "go-sdk-demo-col-sparsevec"
	collectionAlias := "go-sdk-demo-col-sparsevec-alias"

	ctx := context.Background()
	testVdb, err := NewDemo("vdb http url or ip and port", "root", "key get from web console")
	printErr(err)
	err = testVdb.Clear(ctx, database)
	printErr(err)
	err = testVdb.CreateDBAndCollection(ctx, database, collectionName, collectionAlias)
	printErr(err)
	err = testVdb.UpsertData(ctx, database, collectionName)
	printErr(err)
	err = testVdb.QueryData(ctx, database, collectionName)
	printErr(err)
	err = testVdb.UpdateAndDeleteCollection(ctx, database, collectionName)
	printErr(err)
	err = testVdb.DeleteAndDrop(ctx, database, collectionName)
	printErr(err)
}
