package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/ai_document_set"
	collection_view "github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/collection_view"
)

type AIDemo struct {
	client *tcvectordb.Client
}

func NewAIDemo(url, username, key string) (*AIDemo, error) {
	// cli, err := tcvectordb.NewRpcClient(url, username, key, &tcvectordb.ClientOption{
	// 	ReadConsistency: tcvectordb.EventualConsistency})

	// ReadConsistency can be specified when the client is created,
	// and ReadConsistency will be used in subsequent calls to the sdk interface
	cli, err := tcvectordb.NewClient(url, username, key, &tcvectordb.ClientOption{
		ReadConsistency: tcvectordb.EventualConsistency})
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

func (d *AIDemo) DeleteAndDrop(ctx context.Context, database, collectionView, documentSetName string) error {
	// 删除Document
	log.Println("-------------------------- Delete Document --------------------------")
	cdocDelResult, err := d.client.AIDatabase(database).CollectionView(collectionView).Delete(ctx, tcvectordb.DeleteAIDocumentSetParams{
		DocumentSetNames: []string{documentSetName},
	})
	if err != nil {
		return err
	}
	log.Printf("delete document result: %+v", cdocDelResult)

	// 删除collectionView，删除collectionView的同时，其中的数据也将被全部删除
	log.Println("-------------------------- DropCollectionView --------------------------")
	colDropResult, err := d.client.AIDatabase(database).DropCollectionView(ctx, collectionView)
	if err != nil {
		return err
	}
	log.Printf("drop collection result: %+v", colDropResult)

	log.Println("--------------------------- DropDatabase ---------------------------")
	// 删除db，db下的所有collectionView都将被删除
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

func (d *AIDemo) CreateCollectionView(ctx context.Context, database, collectionView string) error {
	db := d.client.AIDatabase(database)

	log.Println("------------------------- CreateCollectionView -------------------------")
	index := tcvectordb.Indexes{}
	index.FilterIndex = append(index.FilterIndex, tcvectordb.FilterIndex{FieldName: "test_str", FieldType: tcvectordb.String, IndexType: tcvectordb.FILTER})

	db.WithTimeout(time.Second * 30)

	enableWordsEmbedding := false
	appendTitleToChunk := false
	appendKeywordsToChunk := true

	shardNum := uint32(3)
	replicaNum := uint32(1)

	_, err := db.CreateCollectionView(ctx, collectionView, tcvectordb.CreateCollectionViewParams{
		Description: "desc",
		Indexes:     index,
		Embedding: &collection_view.DocumentEmbedding{
			Language:             string(tcvectordb.LanguageChinese),
			EnableWordsEmbedding: &enableWordsEmbedding,
		},
		SplitterPreprocess: &collection_view.SplitterPreprocess{
			AppendTitleToChunk:    &appendTitleToChunk,
			AppendKeywordsToChunk: &appendKeywordsToChunk,
		},
		// parsing files with vision model for all in this collectionView
		// vision model parsing only for pdf filetype, and algorithm parsing for other supported filetypes
		ParsingProcess: &api.ParsingProcess{
			ParsingType: string(tcvectordb.VisionModelParsing),
		},
		ShardNum:   &shardNum,
		ReplicaNum: &replicaNum,
	})

	if err != nil {
		return err
	}

	log.Println("-------------------------- ListCollectionViews --------------------------")
	// 列出所有 CollectionView
	collListRes, err := db.ListCollectionViews(ctx)
	if err != nil {
		return err
	}
	for _, col := range collListRes.CollectionViews {
		log.Printf("ListCollectionViews: %+v", col)
	}
	return nil
}

func (d *AIDemo) LoadAndSplitText(ctx context.Context, database, collection, filePath string) (*tcvectordb.LoadAndSplitTextResult, error) {
	log.Println("---------------------------- UploadFile ---------------------------")
	coll := d.client.AIDatabase(database).CollectionView(collection)

	appendTitleToChunk := true
	appendKeywordsToChunk := false
	chunkSplitter := "\n\n"

	res, err := coll.LoadAndSplitText(ctx, tcvectordb.LoadAndSplitTextParams{
		LocalFilePath: filePath,
		MetaData: map[string]interface{}{
			"test_str": "v1",
			"fileKey":  1024,
			"author":   "sam",
		},
		SplitterPreprocess: ai_document_set.DocumentSplitterPreprocess{
			ChunkSplitter:         &chunkSplitter,
			AppendTitleToChunk:    &appendTitleToChunk,
			AppendKeywordsToChunk: &appendKeywordsToChunk,
		},
		// parsing this file with vision model
		// vision model parsing only for pdf filetype, and algorithm parsing for other supported filetypes
		ParsingProcess: &api.ParsingProcess{
			ParsingType: string(tcvectordb.VisionModelParsing),
		},
	})
	if err != nil {
		return nil, err
	}
	log.Printf("UploadFileResult: %+v", res)
	return res, nil
}

func (d *AIDemo) GetFile(ctx context.Context, database, collection, fileName string) error {
	coll := d.client.AIDatabase(database).CollectionView(collection)
	log.Println("---------------------------- GetFile by Name ----------------------------")
	for {
		result, err := coll.Query(ctx, tcvectordb.QueryAIDocumentSetParams{
			DocumentSetName: []string{fileName},
		})
		if err != nil {
			return err
		}
		if len(result.Documents) == 0 {
			return fmt.Errorf("file %v not found", fileName)
		}

		if result.Documents[0].DocumentSetInfo != nil {
			if *result.Documents[0].DocumentSetInfo.IndexedStatus == "Ready" {
				log.Printf("QueryDocument: %+v", result.Documents[0])
				return nil
			} else {
				log.Printf("file %v is not Ready, status: %v", fileName, *result.Documents[0].DocumentSetInfo.IndexedStatus)
			}
		} else {
			return fmt.Errorf("file %v documentSetInfo is nil", fileName)
		}

		time.Sleep(time.Second * 10)
	}

	return nil
}

func (d *AIDemo) GetChunks(ctx context.Context, database, collection, fileName string) error {
	coll := d.client.AIDatabase(database).CollectionView(collection)
	log.Println("---------------------------- GetChunks by Name ----------------------------")
	result, err := coll.GetChunks(ctx, tcvectordb.GetAIDocumentSetChunksParams{
		DocumentSetName: fileName,
	})
	if err != nil {
		return err
	}
	log.Printf("GetChunks, count: %v", result.Count)
	for _, chunk := range result.Chunks {
		log.Printf("chunk: %+v", chunk)
	}
	return nil
}

func (d *AIDemo) QueryAndSearch(ctx context.Context, database, collectionView string) error {
	db := d.client.AIDatabase(database)
	coll := db.CollectionView(collectionView)

	log.Println("---------------------------- Search ----------------------------")
	//enableRerank := true
	res, err := coll.Search(ctx, tcvectordb.SearchAIDocumentSetsParams{
		Content:     "平安保险的偿付能力是什么水平？",
		ExpandChunk: []int{1, 0},
		Filter:      tcvectordb.NewFilter(`test_str="v1"`),
		Limit:       2,
		// RerankOption: &ai_document_set.RerankOption{
		// 	Enable:                &enableRerank,
		// 	ExpectRecallMultiples: 2.5,
		// },
	})
	if err != nil {
		return err
	}
	for _, doc := range res.Documents {
		log.Printf("SearchDocument: %+v", doc)
	}

	log.Println("---------------------------- Update ----------------------------")
	updateRes, err := coll.Update(ctx, map[string]interface{}{
		"test_str": "v2",
	}, tcvectordb.UpdateAIDocumentSetParams{
		Filter: tcvectordb.NewFilter(`test_str="v1"`),
	})
	if err != nil {
		return err
	}
	log.Printf("updateResult: %+v", updateRes)

	log.Println("---------------------------- Query ----------------------------")
	queryRes, err := coll.Query(ctx, tcvectordb.QueryAIDocumentSetParams{
		Filter: tcvectordb.NewFilter(`test_str="v2"`),
		Limit:  1,
	})
	if err != nil {
		return err
	}
	for _, doc := range queryRes.Documents {
		log.Printf("QueryDocument: %+v", doc)
	}
	return nil
}

func (d *AIDemo) Alias(ctx context.Context, database, collectionView, alias string) error {
	db := d.client.AIDatabase(database)
	log.Println("---------------------------- SetAlias ----------------------------")
	setRes, err := db.SetAlias(ctx, collectionView, alias)
	if err != nil {
		return err
	}
	log.Printf("SetAlias result: %+v", setRes)

	log.Println("----------------------- DescribeCollectionView -----------------------")
	collRes, err := db.DescribeCollectionView(ctx, collectionView)
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

func printErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	database := "go-sdk-demo-ai-db"
	collectionView := "go-sdk-demo-ai-col"
	collectionViewAlias := "go-sdk-demo-ai-alias"

	ctx := context.Background()
	testVdb, err := NewAIDemo("vdb http url or ip and port", "vdb username", "key get from web console")
	printErr(err)
	defer testVdb.client.Close()

	err = testVdb.Clear(ctx, database)
	printErr(err)
	err = testVdb.CreateAIDatabase(ctx, database)
	printErr(err)
	err = testVdb.CreateCollectionView(ctx, database, collectionView)
	printErr(err)
	// 当前支持的文件格式markdown(.md或.markdown)、pdf(.pdf)、ppt(.pptx)、word(.docx)
	loadFileRes, err := testVdb.LoadAndSplitText(ctx, database, collectionView, "../demo_files/demo_vision_model_parsing.pdf")
	printErr(err)
	err = testVdb.GetFile(ctx, database, collectionView, loadFileRes.DocumentSetName)
	printErr(err)
	err = testVdb.GetChunks(ctx, database, collectionView, loadFileRes.DocumentSetName)
	printErr(err)
	err = testVdb.QueryAndSearch(ctx, database, collectionView)
	printErr(err)
	err = testVdb.Alias(ctx, database, collectionView, collectionViewAlias)
	printErr(err)
	err = testVdb.DeleteAndDrop(ctx, database, collectionView, loadFileRes.DocumentSetName)
	printErr(err)
}
