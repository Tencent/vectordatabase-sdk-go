package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvdbtext/encoder"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
)

type Demo struct {
	//client *tcvectordb.Client
	client *tcvectordb.RpcClient
	Bm25   encoder.SparseEncoder
}

func NewDemo(url, username, key string) (*Demo, error) {
	cli, err := tcvectordb.NewRpcClient(url, username, key, &tcvectordb.ClientOption{
		ReadConsistency: tcvectordb.EventualConsistency})
	// cli, err := tcvectordb.NewClient(url, username, key, &tcvectordb.ClientOption{
	// 	ReadConsistency: tcvectordb.EventualConsistency})
	if err != nil {
		return nil, err
	}
	// disable/enable http request log print
	// cli.Debug(false)

	bm25, err := encoder.NewBM25Encoder(&encoder.BM25EncoderParams{Bm25Language: "zh"})
	if err != nil {
		return nil, err
	}
	return &Demo{client: cli,
		Bm25: bm25}, nil
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

	log.Println("------------------------- CreateCollection with Embedding-------------------------")
	// 新建 Collection
	// 第一步，设计索引（不是设计表格的结构）
	// 1. 【重要的事】向量对应的文本字段不要建立索引，会浪费较大的内存，并且没有任何作用。
	// 2. 【必须的索引】：主键 id、向量字段 vector 这两个字段目前是固定且必须的，参考下面的例子；
	// 3. 【其他索引】：检索时需作为条件查询的字段，比如要按书籍的作者进行过滤，这个时候author字段就需要建立索引，
	//     否则无法在查询的时候对 author 字段进行过滤，不需要过滤的字段无需加索引，会浪费内存；
	// 4.  向量数据库支持动态 Schema，写入数据时可以写入任何字段，无需提前定义，类似 MongoDB.
	// 5.  例子中创建一个书籍片段的索引，例如书籍片段的信息包括 {id, vector, segment, bookName, page},
	//     id 为主键需要全局唯一，segment 为文本片段, vector 为 segment 的向量，vector 字段需要建立向量索引，假如我们在查询的时候要查询指定书籍
	//     名称的内容，这个时候需要对bookName建立索引，其他字段没有条件查询的需要，无需建立索引。
	// 6.  创建带 Embedding 的 collection 需要保证设置的 vector 索引的维度和 Embedding 所用模型生成向量维度一致，模型及维度关系：
	//     -----------------------------------------------------
	//             bge-base-zh                 ｜ 768
	//             bge-large-zh                ｜ 1024
	//             m3e-base                    ｜ 768
	//             text2vec-large-chinese      ｜ 1024
	//             e5-large-v2                 ｜ 1024
	//             multilingual-e5-base        ｜ 768
	//     -----------------------------------------------------
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
	index.FilterIndex = append(index.FilterIndex, tcvectordb.FilterIndex{FieldName: "bookName", FieldType: tcvectordb.String, IndexType: tcvectordb.FILTER})
	index.FilterIndex = append(index.FilterIndex, tcvectordb.FilterIndex{FieldName: "page", FieldType: tcvectordb.Uint64, IndexType: tcvectordb.FILTER})

	ebd := &tcvectordb.Embedding{VectorField: "vector", Field: "text", ModelName: "bge-base-zh"}

	// 第二步：创建 Collection
	// 创建支持 Embedding 的 Collection
	db.WithTimeout(time.Second * 30)
	_, err = db.CreateCollection(ctx, collection, 3, 1, "test collection", index, &tcvectordb.CreateCollectionParams{
		Embedding: ebd,
	})
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
	log.Println("------------------------------ Upsert ------------------------------")
	// upsert 写入数据，可能会有一定延迟
	// 1. 支持动态 Schema，除了 id、vector 字段必须写入，可以写入其他任意字段；
	// 2. upsert 会执行覆盖写，若文档id已存在，则新数据会直接覆盖原有数据(删除原有数据，再插入新数据)

	segments := []string{
		"1. 区域中间区域的普通人期望契约C级和B级灵兽，其中一人刚晋升为七品中期御兽师，略过了E级和D级灵兽蛋区域。 2. 御兽神殿的一男一女隐藏了一颗S级灵兽蛋在C级灵兽蛋中，他们对在这所偏远县城的御兽高中能否出现能契约S级灵兽的天骄表示怀疑。 3. 御兽神殿大祭司预言百年未被认可的S级灵兽蛋将在这所高中找到命定主人。 4. 一名清秀少年在人群中径直走向混在C级灵兽蛋区域中的那颗S级灵兽蛋。",
		"腾讯云向量数据库（Tencent Cloud VectorDB）是一款全托管的自研企业级分布式数据库服务，专用于存储、索引、检索、管理由深度神经网络或其他机器学习模型生成的大量多维嵌入向量。",
		"作为专门为处理输入向量查询而设计的数据库，它支持多种索引类型和相似度计算方法，单索引支持10亿级向量规模，高达百万级 QPS 及毫秒级查询延迟。",
		"不仅能为大模型提供外部知识库，提高大模型回答的准确性，还可广泛应用于推荐系统、NLP 服务、计算机视觉、智能客服等 AI 领域。",
		"腾讯云向量数据库（Tencent Cloud VectorDB）作为一种专门存储和检索向量数据的服务提供给用户， 在高性能、高可用、大规模、低成本、简单易用、稳定可靠等方面体现出显著优势。 ",
		"腾讯云向量数据库可以和大语言模型 LLM 配合使用。企业的私域数据在经过文本分割、向量化后，可以存储在腾讯云向量数据库中，构建起企业专属的外部知识库，从而在后续的检索任务中，为大模型提供提示信息，辅助大模型生成更加准确的答案。",
	}

	// 如需了解分词的情况，可参考下一行代码获取
	// tokens := d.Bm25.GetTokenizer().Tokenize(segments[0])
	// fmt.Println("tokens: ", tokens)

	sparseVectors, err := d.Bm25.EncodeTexts(segments)
	if err != nil {
		return err
	}

	documentList := make([]tcvectordb.Document, 0)
	for i := 0; i < 5; i++ {
		id := "000" + strconv.Itoa(i)
		documentList = append(documentList, tcvectordb.Document{
			Id: id,
			Fields: map[string]tcvectordb.Field{
				"text": {Val: segments[i]},
			},
			SparseVector: sparseVectors[i],
		})
	}

	result, err := d.client.Upsert(ctx, database, collection, documentList)
	if err != nil {
		return err
	}
	if result.EmbeddingExtraInfo != nil {
		log.Printf("UpsertResult: token used: %v", result.EmbeddingExtraInfo.TokenUsed)
	}
	return nil
}

func (d *Demo) HybridSearchWithEmbedding(ctx context.Context, database, collection string) error {
	log.Println("------------------------------ hybridSearch ------------------------------")

	searchText := "1. 区域中间区域的普通人期望契约C级和B级灵兽，其中一人刚晋升为七品中期御兽师，略过了E级和D级灵兽蛋区域。 2. 御兽神殿的一男一女隐藏了一颗S级灵兽蛋在C级灵兽蛋中，他们对在这所偏远县城的御兽高中能否出现能契约S级灵兽的天骄表示怀疑。 3. 御兽神殿大祭司预言百年未被认可的S级灵兽蛋将在这所高中找到命定主人。 4. 一名清秀少年在人群中径直走向混在C级灵兽蛋区域中的那颗S级灵兽蛋。"
	annSearch := &tcvectordb.AnnParam{
		FieldName: "text",
		Data:      searchText,
	}

	sparseVec, err := d.Bm25.EncodeQuery(searchText)
	if err != nil {
		return err
	}
	keywordSearch := &tcvectordb.MatchOption{
		FieldName:       "sparse_vector",
		Data:            sparseVec,
		TerminateAfter:  4000,
		CutoffFrequency: 0.1,
	}

	limit := 2
	searchRes, err := d.client.HybridSearch(ctx, database, collection, tcvectordb.HybridSearchDocumentParams{
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
	if searchRes.EmbeddingExtraInfo != nil {
		log.Printf("HybridSearchDocumentResult: token used: %v", searchRes.EmbeddingExtraInfo.TokenUsed)
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
	testVdb, err := NewDemo("vdb http url or ip and port", "vdb username", "key get from web console")
	printErr(err)
	defer testVdb.client.Close()

	err = testVdb.Clear(ctx, database)
	printErr(err)
	err = testVdb.CreateDBAndCollection(ctx, database, collectionName, collectionAlias)
	printErr(err)
	err = testVdb.UpsertData(ctx, database, collectionName)
	printErr(err)
	err = testVdb.HybridSearchWithEmbedding(ctx, database, collectionName)
	printErr(err)
	err = testVdb.DeleteAndDrop(ctx, database, collectionName)
	printErr(err)
}
