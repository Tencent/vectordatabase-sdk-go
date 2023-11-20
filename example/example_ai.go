package example

import (
	"context"
	"log"
	"time"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb/api/ai_collection"
)

type AIDemo struct {
	client *tcvectordb.Client
}

func NewAIDemo(url, username, key string) (*AIDemo, error) {
	cli, err := tcvectordb.NewClient(url, username, key, &tcvectordb.ClientOption{ReadConsistency: tcvectordb.EventualConsistency})
	if err != nil {
		return nil, err
	}
	// disable/enable http request log print
	// cli.Debug(false)
	return &AIDemo{client: cli}, nil
}

func (d *AIDemo) Clear(ctx context.Context, database string) error {
	log.Println("--------------------------- DropDatabase ---------------------------")
	result, err := d.client.DropAIDatabase(ctx, database)
	if err != nil {
		return err
	}
	log.Printf("drop database result: %+v", result)
	return nil
}

func (d *AIDemo) DeleteAndDrop(ctx context.Context, database, collection, fileName string) error {
	// 删除collection，删除collection的同时，其中的数据也将被全部删除
	log.Println("-------------------------- Delete Document --------------------------")
	cdocDelResult, err := d.client.AIDatabase(database).Collection(collection).Delete(ctx, &tcvectordb.DeleteAIDocumentOption{
		FileName: fileName,
	})
	if err != nil {
		return err
	}
	log.Printf("delete document result: %+v", cdocDelResult)

	// 删除collection，删除collection的同时，其中的数据也将被全部删除
	log.Println("-------------------------- DropCollection --------------------------")
	colDropResult, err := d.client.AIDatabase(database).DropCollection(ctx, collection)
	if err != nil {
		return err
	}
	log.Printf("drop collection result: %+v", colDropResult)

	log.Println("--------------------------- DropDatabase ---------------------------")
	// 删除db，db下的所有collection都将被删除
	dbDropResult, err := d.client.DropAIDatabase(ctx, database)
	if err != nil {
		return err
	}
	log.Printf("drop database result: %+v", dbDropResult)
	return nil
}

func (d *AIDemo) CreateAIDatabase(ctx context.Context, database string) error {
	log.Println("-------------------------- CreateDatabase --------------------------")
	_, err := d.client.CreateAIDatabase(ctx, database)
	if err != nil {
		return err
	}
	log.Println("--------------------------- ListDatabase ---------------------------")
	dbList, err := d.client.ListDatabase(ctx)
	if err != nil {
		return err
	}

	log.Printf("base database ======================")
	for _, db := range dbList.Databases {
		log.Printf("database: %s, createTime: %s, dbType: %s", db.DatabaseName, db.Info.CreateTime, db.Info.DbType)
	}

	log.Printf("ai database ======================")
	for _, db := range dbList.AIDatabases {
		log.Printf("database: %s, createTime: %s, dbType: %s", db.DatabaseName, db.Info.CreateTime, db.Info.DbType)
	}
	return nil
}

func (d *AIDemo) CreateCollection(ctx context.Context, database, collection string) error {
	db := d.client.AIDatabase(database)

	log.Println("------------------------- CreateCollection -------------------------")
	index := tcvectordb.Indexes{}
	index.FilterIndex = append(index.FilterIndex, tcvectordb.FilterIndex{FieldName: "teststr", FieldType: tcvectordb.String, IndexType: tcvectordb.FILTER})

	db.WithTimeout(time.Second * 30)

	enableWordsEmbedding := true
	_, err := db.CreateCollection(ctx, collection, &tcvectordb.CreateAICollectionOption{
		Description: "desc",
		Indexes:     index,
		AiConfig: &tcvectordb.AiConfig{
			ExpectedFileNum: 100,
			AverageFileSize: 102400,
			Language:        tcvectordb.LanguageChinese,
			DocumentPreprocess: &ai_collection.DocumentPreprocess{
				AppendTitleToChunk:    "1",
				AppendKeywordsToChunk: "0",
			},
			EnableWordsEmbedding: &enableWordsEmbedding,
		},
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
	return nil
}

func (d *AIDemo) UploadFile(ctx context.Context, database, collection, filePath string) (*tcvectordb.UploadAIDocumentResult, error) {
	log.Println("---------------------------- UploadFile ---------------------------")
	coll := d.client.AIDatabase(database).Collection(collection)
	res, err := coll.Upload(ctx, filePath, &tcvectordb.UploadAIDocumentOption{
		MetaData: map[string]tcvectordb.Field{
			"teststr": {Val: "v1"},
			"filekey": {Val: 1024},
			"author":  {Val: "sam"},
		},
	})
	if err != nil {
		return nil, err
	}
	log.Printf("UploadFileResult: %+v", res)
	return res, nil
}

func (d *AIDemo) GetFile(ctx context.Context, database, collection, fileName string) error {
	coll := d.client.AIDatabase(database).Collection(collection)
	log.Println("---------------------------- GetFile by Name ----------------------------")
	result, err := coll.Query(ctx, &tcvectordb.QueryAIDocumentOption{
		FileName: fileName,
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

func (d *AIDemo) QueryAndSearch(ctx context.Context, database, collection string) error {
	db := d.client.AIDatabase(database)
	coll := db.Collection(collection)

	log.Println("---------------------------- Search ----------------------------")
	// 查找与给定查询向量相似的向量。支持输入文本信息检索与输入文本相似的内容，同时，支持搭配标量字段的 Filter 表达式一并检索。
	enableRerank := true
	res, err := coll.Search(ctx, "“什么是向量数据库", &tcvectordb.SearchAIDocumentOption{
		ChunkExpand: []int{1, 0},
		Filter:      tcvectordb.NewFilter(`teststr="v1"`),
		Limit:       2,
		RerankOption: &tcvectordb.RerankOption{
			Enable:                &enableRerank,
			ExpectRecallMultiples: 2.5,
		},
	})
	if err != nil {
		return err
	}
	for _, doc := range res.Documents {
		log.Printf("SearchDocument: %+v", doc)
	}

	log.Println("---------------------------- Update ----------------------------")
	updateRes, err := coll.Update(ctx, &tcvectordb.UpdateAIDocumentOption{
		QueryFilter: tcvectordb.NewFilter(`teststr="v1"`),
		UpdateFields: map[string]interface{}{
			"teststr": "v2",
		},
	})
	if err != nil {
		return err
	}
	log.Printf("updateResult: %+v", updateRes)

	log.Println("---------------------------- Query ----------------------------")
	queryRes, err := coll.Query(ctx, &tcvectordb.QueryAIDocumentOption{
		Filter: tcvectordb.NewFilter(`teststr="v2"`),
	})
	if err != nil {
		return err
	}
	for _, doc := range queryRes.Documents {
		log.Printf("QueryDocument: %+v", doc)
	}
	return nil
}

func (d *AIDemo) Alias(ctx context.Context, database, collection, alias string) error {
	db := d.client.AIDatabase(database)
	log.Println("---------------------------- SetAlias ----------------------------")
	setRes, err := db.SetAlias(ctx, collection, alias)
	if err != nil {
		return err
	}
	log.Printf("SetAlias result: %+v", setRes)

	log.Println("----------------------- DescribeCollection -----------------------")
	collRes, err := db.DescribeCollection(ctx, collection)
	if err != nil {
		return err
	}
	log.Printf("Collection: %+v", collRes)

	log.Println("--------------------------- DeleteAlias ---------------------------")
	delRes, err := db.DeleteAlias(ctx, alias)
	if err != nil {
		return err
	}
	log.Printf("SetAlias result: %+v", delRes)
	return nil
}
