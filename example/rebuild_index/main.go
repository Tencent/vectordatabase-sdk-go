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
	// cli.Debug(false)
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
	index.VectorIndex = append(index.VectorIndex, tcvectordb.VectorIndex{
		FilterIndex: tcvectordb.FilterIndex{
			FieldName: "vector",
			FieldType: tcvectordb.Vector,
			IndexType: tcvectordb.HNSW,
		},
		Dimension:  768,
		MetricType: tcvectordb.COSINE,
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

	ebd := &tcvectordb.Embedding{VectorField: "vector", Field: "text", ModelName: "bge-base-zh"}

	db.WithTimeout(time.Second * 30)
	_, err = db.CreateCollectionIfNotExists(ctx, collection, 3, 1, "test collection", index, &tcvectordb.CreateCollectionParams{
		Embedding: ebd,
	})
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

func (d *Demo) QueryData(ctx context.Context, database, collection string) error {
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
	return nil
}

func (d *Demo) RebuildIndex(ctx context.Context, database, collection string) error {
	log.Println("------------------------------ Rebuild vector ------------------------------")

	_, err := d.client.RebuildIndex(ctx, database, collection, &tcvectordb.RebuildIndexParams{
		FieldName:         "vector",
		DropBeforeRebuild: true,
		Throttle:          1,
	})
	if err != nil {
		return err
	}
	log.Println("------------------------------ Rebuild sparse vector ------------------------------")
	_, err = d.client.RebuildIndex(ctx, database, collection, &tcvectordb.RebuildIndexParams{
		FieldName:         "sparse_vector",
		DropBeforeRebuild: true,
		Throttle:          1,
	})
	if err != nil {
		return err
	}
	return nil
}

func printErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	database := "go-demo-db-rebuild"
	collectionName := "go-demo-col-rebuild"

	ctx := context.Background()
	testVdb, err := NewDemo("vdb http url or ip and port", "vdb username", "key get from web console")
	printErr(err)
	defer testVdb.client.Close()

	err = testVdb.CreateDBAndCollection(ctx, database, collectionName)
	printErr(err)
	err = testVdb.UpsertData(ctx, database, collectionName)
	printErr(err)
	err = testVdb.QueryData(ctx, database, collectionName)
	printErr(err)
	err = testVdb.RebuildIndex(ctx, database, collectionName)
	printErr(err)
	err = testVdb.DropDB(ctx, database)
	printErr(err)
}
