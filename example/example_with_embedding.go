package main

import (
	"context"
	"log"
	"time"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb"
)

type EmbeddingDemo struct {
	Demo
}

func NewEmbeddingDemo(url, username, key string) (*EmbeddingDemo, error) {
	cli, err := tcvectordb.NewClient(url, username, key, &entity.ClientOption{ReadConsistency: entity.EventualConsistency})
	if err != nil {
		return nil, err
	}
	// disable/enable http request log print
	// cli.Debug(false)
	return &EmbeddingDemo{Demo: Demo{client: cli}}, nil
}

func (d *EmbeddingDemo) CreateDBAndCollection(ctx context.Context, database, collection, alias string) error {
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
	index := entity.Indexes{}
	index.VectorIndex = append(index.VectorIndex, entity.VectorIndex{
		FilterIndex: entity.FilterIndex{
			FieldName: "vector",
			FieldType: entity.Vector,
			IndexType: entity.HNSW,
		},
		Dimension:  768,
		MetricType: entity.COSINE,
		Params: &entity.HNSWParam{
			M:              16,
			EfConstruction: 200,
		},
	})
	index.FilterIndex = append(index.FilterIndex, entity.FilterIndex{FieldName: "id", FieldType: entity.String, IndexType: entity.PRIMARY})
	index.FilterIndex = append(index.FilterIndex, entity.FilterIndex{FieldName: "bookName", FieldType: entity.String, IndexType: entity.FILTER})
	index.FilterIndex = append(index.FilterIndex, entity.FilterIndex{FieldName: "page", FieldType: entity.Uint64, IndexType: entity.FILTER})

	ebd := &entity.Embedding{VectorField: "vector", Field: "text", Model: entity.BGE_BASE_ZH}
	// 第二步：创建 Collection
	// 创建支持 Embedding 的 Collection
	db.WithTimeout(time.Second * 30)
	_, err = db.CreateCollection(ctx, collection, 3, 2, "test collection", index, &entity.CreateCollectionOption{
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
	delAliasRes, err := db.DeleteAlias(ctx, alias, nil)
	if err != nil {
		return err
	}
	log.Printf("DeleteAliasResult: %v", delAliasRes)
	return nil
}

func (d *EmbeddingDemo) UpsertData(ctx context.Context, database, collection string) error {
	// 获取 Collection 对象
	coll := d.client.Database(database).Collection(collection)

	log.Println("------------------------------ Upsert ------------------------------")
	// upsert 写入数据，可能会有一定延迟
	// 1. 支持动态 Schema，除了 id、vector 字段必须写入，可以写入其他任意字段；
	// 2. upsert 会执行覆盖写，若文档id已存在，则新数据会直接覆盖原有数据(删除原有数据，再插入新数据)

	documentList := []entity.Document{
		{
			Id: "0001",
			Fields: map[string]entity.Field{
				"bookName": {Val: "西游记"},
				"author":   {Val: "吴承恩"},
				"page":     {Val: 21},
				"segment":  {Val: "富贵功名，前缘分定，为人切莫欺心。"},
				"text":     {Val: "富贵功名，前缘分定，为人切莫欺心。"},
			},
		},
		{
			Id: "0002",
			Fields: map[string]entity.Field{
				"bookName": {Val: "西游记"},
				"author":   {Val: "吴承恩"},
				"page":     {Val: 22},
				"segment":  {Val: "正大光明，忠良善果弥深。些些狂妄天加谴，眼前不遇待时临。"},
				"text":     {Val: "正大光明，忠良善果弥深。些些狂妄天加谴，眼前不遇待时临。"},
			},
		},
		{
			Id: "0003",
			Fields: map[string]entity.Field{
				"bookName": {Val: "三国演义"},
				"author":   {Val: "罗贯中"},
				"page":     {Val: 23},
				"segment":  {Val: "细作探知这个消息，飞报吕布。"},
				"text":     {Val: "细作探知这个消息，飞报吕布。"},
			},
		},
		{
			Id: "0004",
			Fields: map[string]entity.Field{
				"bookName": {Val: "三国演义"},
				"author":   {Val: "罗贯中"},
				"page":     {Val: 24},
				"segment":  {Val: "布大惊，与陈宫商议。宫曰：“闻刘玄德新领徐州，可往投之。”布从其言，竟投徐州来。有人报知玄德。"},
				"text":     {Val: "布大惊，与陈宫商议。宫曰：“闻刘玄德新领徐州，可往投之。”布从其言，竟投徐州来。有人报知玄德。"},
			},
		},
		{
			Id: "0005",
			Fields: map[string]entity.Field{
				"bookName": {Val: "三国演义"},
				"author":   {Val: "罗贯中"},
				"page":     {Val: 25},
				"segment":  {Val: "玄德曰：“布乃当今英勇之士，可出迎之。”糜竺曰：“吕布乃虎狼之徒，不可收留；收则伤人矣。"},
				"text":     {Val: "玄德曰：“布乃当今英勇之士，可出迎之。”糜竺曰：“吕布乃虎狼之徒，不可收留；收则伤人矣。"},
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

func (d *EmbeddingDemo) QueryData(ctx context.Context, database, collection string) error {
	// 获取 Collection 对象
	coll := d.client.Database(database).Collection(collection)

	log.Println("------------------------------ Query ------------------------------")
	// 查询
	// 1. query 用于查询数据
	// 2. 可以通过传入主键 id 列表或 filter 实现过滤数据的目的
	// 3. 如果没有主键 id 列表和 filter 则必须传入 limit 和 offset，类似 scan 的数据扫描功能
	// 4. 如果仅需要部分 field 的数据，可以指定 output_fields 用于指定返回数据包含哪些 field，不指定默认全部返回
	documentIds := []string{"0001", "0002", "0003", "0004", "0005"}
	filter := entity.NewFilter(`bookName="三国演义"`)
	outputField := []string{"id", "bookName"}

	result, err := coll.Query(ctx, documentIds, &entity.QueryDocumentOption{
		Filter:         filter,
		RetrieveVector: true,
		OutputFields:   outputField,
		Limit:          2,
		Offset:         1,
	})
	if err != nil {
		return err
	}
	log.Printf("QueryResult: total: %v, affect: %v", result.Total, result.AffectedCount)
	for _, doc := range result.Documents {
		log.Printf("QueryDocument: %+v", doc)
	}

	log.Println("---------------------------- SearchById ----------------------------")
	// searchById
	// 1. searchById 提供按 id 搜索的能力
	// 1. search 提供按照 vector 搜索的能力
	// 2. 支持通过 filter 过滤数据
	// 3. 如果仅需要部分 field 的数据，可以指定 output_fields 用于指定返回数据包含哪些 field，不指定默认全部返回
	// 4. limit 用于限制每个单元搜索条件的条数，如 vector 传入三组向量，limit 为 3，则 limit 限制的是每组向量返回 top 3 的相似度向量

	// 根据主键 id 查找 Top K 个相似性结果，向量数据库会根据ID 查找对应的向量，再根据向量进行TOP K 相似性检索
	searchResult, err := coll.SearchById(ctx, []string{"0003"}, &entity.SearchDocumentOption{
		Filter: filter,
		Params: &entity.SearchDocParams{Ef: 200},
		Limit:  2,
	})
	if err != nil {
		return err
	}
	for i, item := range searchResult.Documents {
		log.Printf("SearchDocumentResult, index: %d ==================", i)
		for _, doc := range item {
			log.Printf("SearchDocument: %+v", doc)
		}
	}

	log.Println("--------------------------- SearchByText ---------------------------")
	// 通过 embedding 文本搜索
	// 1. searchByText 提供基于 embedding 文本的搜索能力，会先将 embedding 内容做 Embedding 然后进行按向量搜索
	// 其他选项类似 search 接口

	// searchByText 返回类型为 Dict，接口查询过程中 embedding 可能会出现截断，如发生截断将会返回响应 warn 信息，如需确认是否截断可以
	// 使用 "warning" 作为 key 从 Dict 结果中获取警告信息，查询结果可以通过 "documents" 作为 key 从 Dict 结果中获取
	searchResult, err = coll.SearchByText(ctx, map[string][]string{"text": {"细作探知这个消息，飞报吕布。"}}, &entity.SearchDocumentOption{
		Params: &entity.SearchDocParams{Ef: 100}, // 若使用HNSW索引，则需要指定参数ef，ef越大，召回率越高，但也会影响检索速度
		Limit:  2,                                // 指定 Top K 的 K 值
	})
	if err != nil {
		return err
	}
	// 输出相似性检索结果，检索结果为二维数组，每一位为一组返回结果，分别对应search时指定的多个向量
	for i, item := range searchResult.Documents {
		log.Printf("SearchDocumentResult, index: %d ==================", i)
		for _, doc := range item {
			log.Printf("SearchDocument: %+v", doc)
		}
	}
	return nil
}
