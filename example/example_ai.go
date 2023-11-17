package example

import (
	"context"
	"log"
	"time"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity/api/ai_collection"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb"
)

type AIDemo struct {
	client *entity.VectorDBClient
}

func NewAIDemo(url, username, key string) (*AIDemo, error) {
	cli, err := tcvectordb.NewClient(url, username, key, &entity.ClientOption{ReadConsistency: entity.EventualConsistency})
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
	cdocDelResult, err := d.client.AIDatabase(database).Collection(collection).Delete(ctx, &entity.DeleteAIDocumentOption{
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
	index := entity.Indexes{}
	index.FilterIndex = append(index.FilterIndex, entity.FilterIndex{FieldName: "teststr", FieldType: entity.String, IndexType: entity.FILTER})

	db.WithTimeout(time.Second * 30)
	_, err := db.CreateCollection(ctx, collection, &entity.CreateAICollectionOption{
		Description: "desc",
		Indexes:     index,
		AiConfig: &entity.AiConfig{
			ExpectedFileNum: 100,
			AverageFileSize: 102400,
			Language:        entity.LanguageChinese,
			DocumentPreprocess: &ai_collection.DocumentPreprocess{
				AppendTitleToChunk:    "1",
				AppendKeywordsToChunk: "0",
			},
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

func (d *AIDemo) UploadFile(ctx context.Context, database, collection, filePath string) (*entity.UploadAIDocumentResult, error) {
	log.Println("---------------------------- UploadFile ---------------------------")
	coll := d.client.AIDatabase(database).Collection(collection)
	res, err := coll.Upload(ctx, filePath, &entity.UploadAIDocumentOption{
		MetaData: map[string]entity.Field{
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
	result, err := coll.Query(ctx, &entity.QueryAIDocumentOption{
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
	res, err := coll.Search(ctx, "“什么是向量数据库", &entity.SearchAIDocumentOption{
		ChunkExpand: []int{1, 0},
		Filter:      entity.NewFilter(`teststr="v1"`),
		Limit:       2,
	})
	if err != nil {
		return err
	}
	for _, doc := range res.Documents {
		log.Printf("SearchDocument: %+v", doc)
	}

	log.Println("---------------------------- Update ----------------------------")
	updateRes, err := coll.Update(ctx, &entity.UpdateAIDocumentOption{
		QueryFilter: entity.NewFilter(`teststr="v1"`),
		UpdateFields: map[string]interface{}{
			"teststr": "v2",
		},
	})
	if err != nil {
		return err
	}
	log.Printf("updateResult: %+v", updateRes)

	log.Println("---------------------------- Query ----------------------------")
	queryRes, err := coll.Query(ctx, &entity.QueryAIDocumentOption{
		Filter: entity.NewFilter(`teststr="v2"`),
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
