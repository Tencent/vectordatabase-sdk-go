package main

import (
	"context"
	"log"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
)

type IVFRabitQDemo struct {
	//client *tcvectordb.RpcClient
	client *tcvectordb.RpcClient
}

func NewIVFRabitQDemo(url, username, key string) (*IVFRabitQDemo, error) {
	// cli, err := tcvectordb.NewRpcClient(url, username, key, &tcvectordb.ClientOption{
	// 	ReadConsistency: tcvectordb.EventualConsistency})
	cli, err := tcvectordb.NewRpcClient(url, username, key, &tcvectordb.ClientOption{
		ReadConsistency: tcvectordb.EventualConsistency})
	if err != nil {
		return nil, err
	}
	// disable/enable http request log print
	// cli.Debug(false)
	return &IVFRabitQDemo{client: cli}, nil
}

func (d *IVFRabitQDemo) DropDatabase(ctx context.Context, database string) error {
	log.Println("--------------------------- DropDatabase ---------------------------")
	result, err := d.client.DropDatabase(ctx, database)
	if err != nil {
		return err
	}
	log.Printf("drop database result: %+v", result)
	return nil
}

func (d *IVFRabitQDemo) CreateDBAndCollection(ctx context.Context, database, collection string) error {
	log.Println("-------------------------- CreateDatabaseIfNotExists --------------------------")
	db, err := d.client.CreateDatabaseIfNotExists(ctx, database)
	if err != nil {
		return err
	}
	log.Println("------------------------- CreateCollection -------------------------")
	bits := uint32(9)
	index := tcvectordb.Indexes{
		VectorIndex: []tcvectordb.VectorIndex{
			{
				FilterIndex: tcvectordb.FilterIndex{
					FieldName: "vector",
					FieldType: tcvectordb.Vector,
					IndexType: tcvectordb.IVF_RABITQ,
				},
				Dimension:  768,
				MetricType: tcvectordb.COSINE,
				Params: &tcvectordb.IVFRabitQParams{
					NList: 1,
					Bits:  &bits,
				},
			},
		},
		FilterIndex: []tcvectordb.FilterIndex{
			{FieldName: "id", FieldType: tcvectordb.String, IndexType: tcvectordb.PRIMARY, AutoId: "uuid"},
			{FieldName: "bookName", FieldType: tcvectordb.String, IndexType: tcvectordb.FILTER},
			{FieldName: "page", FieldType: tcvectordb.Uint64, IndexType: tcvectordb.FILTER},
			{FieldName: "tag", FieldType: tcvectordb.Array, IndexType: tcvectordb.FILTER},
			{FieldName: "double_field", FieldType: tcvectordb.Double, IndexType: tcvectordb.FILTER},
			{FieldName: "int64_field", FieldType: tcvectordb.Int64, IndexType: tcvectordb.FILTER},
		},
	}

	db.WithTimeout(time.Second * 30)
	param := &tcvectordb.CreateCollectionParams{
		Embedding: &tcvectordb.Embedding{
			VectorField: "vector",
			Field:       "text",
			ModelName:   "bge-base-zh",
		},
	}
	coll, err := db.CreateCollectionIfNotExists(ctx, collection, 1, 0, "test collection", index, param)
	if err != nil {
		return err
	}
	log.Printf("create collection if not exists success: %v: %v", coll.DatabaseName, coll.CollectionName)
	return nil
}

func (d *IVFRabitQDemo) UpsertData(ctx context.Context, database, collection string) error {
	time.Sleep(time.Second * 5)
	log.Println("------------------------------ Upsert ------------------------------")
	// upsert 写入数据，可能会有一定延迟
	// 1. 支持动态 Schema，除了 id、vector 字段必须写入，可以写入其他任意字段；
	// 2. upsert 会执行覆盖写，若文档id已存在，则新数据会直接覆盖原有数据(删除原有数据，再插入新数据)

	buildExistedData := false
	documentList := []tcvectordb.Document{
		{
			Fields: map[string]tcvectordb.Field{
				"bookName":     {Val: "西游记"},
				"author":       {Val: "吴承恩"},
				"page":         {Val: 21},
				"segment":      {Val: "富贵功名，前缘分定，为人切莫欺心。"},
				"text":         {Val: "富贵功名，前缘分定，为人切莫欺心。"},
				"double_field": {Val: 0.1},
				"int64_field":  {Val: -1},
			},
		},
		{
			Fields: map[string]tcvectordb.Field{
				"bookName":     {Val: "西游记"},
				"author":       {Val: "吴承恩"},
				"page":         {Val: 22},
				"segment":      {Val: "正大光明，忠良善果弥深。些些狂妄天加谴，眼前不遇待时临。"},
				"text":         {Val: "正大光明，忠良善果弥深。些些狂妄天加谴，眼前不遇待时临。"},
				"double_field": {Val: 0.2},
				"int64_field":  {Val: -2},
			},
		},
		{
			Fields: map[string]tcvectordb.Field{
				"bookName":     {Val: "三国演义"},
				"author":       {Val: "罗贯中"},
				"page":         {Val: 23},
				"segment":      {Val: "细作探知这个消息，飞报吕布。"},
				"text":         {Val: "细作探知这个消息，飞报吕布。"},
				"double_field": {Val: 0.3},
				"int64_field":  {Val: -3},
			},
		},
		{
			Fields: map[string]tcvectordb.Field{
				"bookName":     {Val: "三国演义"},
				"author":       {Val: "罗贯中"},
				"page":         {Val: 24},
				"segment":      {Val: "布大惊，与陈宫商议。宫曰：“闻刘玄德新领徐州，可往投之。”布从其言，竟投徐州来。有人报知玄德。"},
				"text":         {Val: "布大惊，与陈宫商议。宫曰：“闻刘玄德新领徐州，可往投之。”布从其言，竟投徐州来。有人报知玄德。"},
				"double_field": {Val: 0.4},
				"int64_field":  {Val: -4},
			},
		},
		{

			Fields: map[string]tcvectordb.Field{
				"bookName":     {Val: "三国演义"},
				"author":       {Val: "罗贯中"},
				"page":         {Val: 25},
				"segment":      {Val: "玄德曰：“布乃当今英勇之士，可出迎之。”糜竺曰：“吕布乃虎狼之徒，不可收留；收则伤人矣。"},
				"text":         {Val: "玄德曰：“布乃当今英勇之士，可出迎之。”糜竺曰：“吕布乃虎狼之徒，不可收留；收则伤人矣。"},
				"double_field": {Val: 0.5},
				"int64_field":  {Val: -5},
			},
		},
	}

	docs := make([]tcvectordb.Document, 0)
	for i := 0; i < 10; i++ {
		docs = append(docs, documentList...)
	}
	result, err := d.client.Upsert(ctx, database, collection, docs, &tcvectordb.UpsertDocumentParams{
		BuildIndex: &buildExistedData,
	})
	if err != nil {
		return err
	}
	log.Printf("Upsert affected count: %+v", result.AffectedCount)

	rebuildIndexResult, err := d.client.RebuildIndex(ctx, database, collection)
	if err != nil {
		return err
	}
	log.Printf("RebuildIndexResult: %+v", rebuildIndexResult)
	time.Sleep(time.Second * 5)
	return nil
}

func (d *IVFRabitQDemo) QueryAndSearch(ctx context.Context, database, collection string) error {
	log.Println("------------------------------ Query ------------------------------")
	// 查询
	// 1. query 用于查询数据
	// 2. 可以通过传入主键 id 列表或 filter 实现过滤数据的目的
	// 3. 如果没有主键 id 列表和 filter 则必须传入 limit 和 offset，类似 scan 的数据扫描功能
	// 4. 如果仅需要部分 field 的数据，可以指定 output_fields 用于指定返回数据包含哪些 field，不指定默认全部返回
	filter := tcvectordb.NewFilter(`bookName="三国演义" and double_field>0.3 and int64_field<-2`)
	outputField := []string{"id", "bookName", "double_field", "int64_field"}

	result, err := d.client.Query(ctx, database, collection, nil, &tcvectordb.QueryDocumentParams{
		Filter:       filter,
		OutputFields: outputField,
		Limit:        2,
		Offset:       1,
	})
	if err != nil {
		return err
	}
	log.Printf("QueryResult: total: %v, affect: %v", result.Total, result.AffectedCount)
	for _, doc := range result.Documents {
		log.Printf("QueryDocument: %+v", doc)
	}

	log.Println("--------------------------- SearchByText ---------------------------")
	// 通过 embedding 文本搜索
	// 1. searchByText 提供基于 embedding 文本的搜索能力，会先将 embedding 内容做 Embedding 然后进行按向量搜索
	// 其他选项类似 search 接口

	// searchByText 返回类型为 Dict，接口查询过程中 embedding 可能会出现截断，如发生截断将会返回响应 warn 信息，如需确认是否截断可以
	// 使用 "warning" 作为 key 从 Dict 结果中获取警告信息，查询结果可以通过 "documents" 作为 key 从 Dict 结果中获取
	searchResult, err := d.client.SearchByText(ctx, database, collection, map[string][]string{"text": {"细作探知这个消息，飞报吕布。"}}, &tcvectordb.SearchDocumentParams{
		Params: &tcvectordb.SearchDocParams{Nprobe: 1}, // 若使用HNSW索引，则需要指定参数ef，ef越大，召回率越高，但也会影响检索速度
		Limit:  2,                                      // 指定 Top K 的 K 值
		Filter: filter,
	})
	if err != nil {
		return err
	}
	if searchResult.EmbeddingExtraInfo != nil {
		log.Printf("SearchByTextResult: token used: %v", searchResult.EmbeddingExtraInfo.TokenUsed)
	}
	// 输出相似性检索结果，检索结果为二维数组，每一位为一组返回结果，分别对应search时指定的多个向量
	for i, item := range searchResult.Documents {
		log.Printf("SearchDocumentResult, batch search index: %d ==================", i)
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
	database := "go-sdk-ivf-rabitq-demo-db"
	collectionName := "go-sdk-ivf-rabitq-demo-col"
	ctx := context.Background()
	testVdb, err := NewIVFRabitQDemo("vdb http url or ip and port", "vdb username", "key get from web console")
	printErr(err)
	defer testVdb.client.Close()

	err = testVdb.DropDatabase(ctx, database)
	printErr(err)
	err = testVdb.CreateDBAndCollection(ctx, database, collectionName)
	printErr(err)
	err = testVdb.UpsertData(ctx, database, collectionName)
	printErr(err)
	err = testVdb.QueryAndSearch(ctx, database, collectionName)
	printErr(err)
	// err = testVdb.DropDatabase(ctx, database)
	// printErr(err)
}
