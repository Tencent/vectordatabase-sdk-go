package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvdbtext/encoder"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
)

type Demo struct {
	client *tcvectordb.RpcClient
	bm25   encoder.SparseEncoder
}

func NewDemo(url, username, key string) (*Demo, error) {
	cli, err := tcvectordb.NewRpcClient(url, username, key, &tcvectordb.ClientOption{
		ReadConsistency: tcvectordb.EventualConsistency})
	if err != nil {
		return nil, err
	}
	bm25, err := encoder.NewBM25Encoder(&encoder.BM25EncoderParams{Bm25Language: "zh"})
	if err != nil {
		return nil, err
	}
	// disable/enable rpc request log print
	cli.Debug(true)
	return &Demo{client: cli, bm25: bm25}, nil
}

func (d *Demo) DropDB(ctx context.Context, database string) error {
	log.Println("--------------------------- DropDatabase ---------------------------")
	result, err := d.client.DropDatabase(ctx, database)
	if err != nil {
		return err
	}
	log.Printf("drop database result: %+v", result)
	return nil
}

func (d *Demo) CreateDBAndCollection(ctx context.Context, database, collection string) error {

	log.Println("-------------------------- CreateDatabaseIfNotExists --------------------------")
	db, err := d.client.CreateDatabaseIfNotExists(ctx, database)
	if err != nil {
		return err
	}

	log.Println("------------------------- CreateCollection -------------------------")
	index := tcvectordb.Indexes{}
	index.FilterIndex = append(index.FilterIndex, tcvectordb.FilterIndex{FieldName: "id", FieldType: tcvectordb.String, IndexType: tcvectordb.PRIMARY})
	index.SparseVectorIndex = append(index.SparseVectorIndex, tcvectordb.SparseVectorIndex{
		FieldName:  "sparse_vector",
		FieldType:  tcvectordb.SparseVector,
		IndexType:  tcvectordb.SPARSE_INVERTED,
		MetricType: tcvectordb.IP,
	})
	// 第二步：创建 Collection
	// 创建collection耗时较长，需要调整客户端的timeout
	// 这里以三可用区实例作为参考，具体实例不同的规格所支持的shard和replicas区间不同，需要参考官方文档
	db.WithTimeout(time.Second * 30)
	_, err = db.CreateCollectionIfNotExists(ctx, collection, 1, 1, "test collection", index)
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

	return nil
}

func (d *Demo) UpsertData(ctx context.Context, database, collection string) error {

	log.Println("------------------------------ Upsert ------------------------------")

	texts := []string{
		"腾讯云向量数据库（Tencent Cloud VectorDB）是一款全托管的自研企业级分布式数据库服务，专用于存储、索引、检索、管理由深度神经网络或其他机器学习模型生成的大量多维嵌入向量。",
		"作为专门为处理输入向量查询而设计的数据库，它支持多种索引类型和相似度计算方法，单索引支持10亿级向量规模，高达百万级 QPS 及毫秒级查询延迟。",
		"不仅能为大模型提供外部知识库，提高大模型回答的准确性，还可广泛应用于推荐系统、NLP 服务、计算机视觉、智能客服等 AI 领域。",
		"腾讯云向量数据库（Tencent Cloud VectorDB）作为一种专门存储和检索向量数据的服务提供给用户， 在高性能、高可用、大规模、低成本、简单易用、稳定可靠等方面体现出显著优势。 ",
		"腾讯云向量数据库可以和大语言模型 LLM 配合使用。企业的私域数据在经过文本分割、向量化后，可以存储在腾讯云向量数据库中，构建起企业专属的外部知识库，从而在后续的检索任务中，为大模型提供提示信息，辅助大模型生成更加准确的答案。",
	}

	// 如需了解分词的情况，可参考下一行代码获取
	tokens := d.bm25.GetTokenizer().Tokenize(texts[0])
	fmt.Println("tokens: ", tokens)

	sparseVectors, err := d.bm25.EncodeTexts(texts)
	if err != nil {
		return err
	}

	documentList := make([]tcvectordb.Document, 0)
	for i := 0; i < 5; i++ {
		id := "000" + strconv.Itoa(i)
		documentList = append(documentList, tcvectordb.Document{
			Id:           id,
			SparseVector: sparseVectors[i],
			Fields: map[string]tcvectordb.Field{
				"text": {
					Val: texts[i],
				},
			},
		})
	}

	result, err := d.client.Upsert(ctx, database, collection, documentList)
	if err != nil {
		return err
	}
	log.Printf("UpsertResult: %+v", result)
	return nil
}

func (d *Demo) QueryAndSearchData(ctx context.Context, database, collection string) error {
	time.Sleep(2 * time.Second)
	log.Println("------------------------------ Query ------------------------------")
	documentIds := []string{"0000", "0001", "0002", "0003", "0004"}
	outputField := []string{"id", "sparse_vector"}

	result, err := d.client.Query(ctx, database, collection, documentIds, &tcvectordb.QueryDocumentParams{
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

	log.Println("------------------------------ FullTextSearch ------------------------------")

	query, err := d.bm25.EncodeQuery("什么是腾讯云向量数据库")
	if err != nil {
		return err
	}
	keywordSearch := &tcvectordb.FullTextSearchMatchOption{
		FieldName: "sparse_vector",
		Data:      query,
		// TerminateAfter: 4000,
		// CutoffFrequency: 0.1,
	}

	limit := 5
	searchRes, err := d.client.FullTextSearch(ctx, database, collection, tcvectordb.FullTextSearchParams{
		Match:        keywordSearch,
		Limit:        &limit,
		OutputFields: []string{"id", "text"},
	})
	if err != nil {
		return err
	}

	// 输出相似性检索结果，检索结果为二维数组，每一位为一组返回结果，分别对应search时指定的多个向量
	for i, item := range searchRes.Documents {
		log.Printf("SearchDocumentResult, index: %d ==================", i)
		for _, doc := range item {
			log.Printf("SearchDocument: %+v", doc)
		}
	}
	return nil
}

func printErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	database := "go-sdk-demo-db-fulltext-search"
	collectionName := "go-sdk-demo-col-fulltext-search"

	ctx := context.Background()
	testVdb, err := NewDemo("vdb http url or ip and port", "vdb username", "key get from web console")
	printErr(err)
	defer testVdb.client.Close()

	err = testVdb.CreateDBAndCollection(ctx, database, collectionName)
	printErr(err)
	err = testVdb.UpsertData(ctx, database, collectionName)
	printErr(err)
	err = testVdb.QueryAndSearchData(ctx, database, collectionName)
	printErr(err)
	err = testVdb.DropDB(ctx, database)
	printErr(err)
}
