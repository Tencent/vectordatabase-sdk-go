package main

import (
	"context"
	"log"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/ai_document_set"
)

type Demo struct {
	client *tcvectordb.RpcClient
}

func NewDemo(url, username, key string) (*Demo, error) {
	cli, err := tcvectordb.NewRpcClient(url, username, key, &tcvectordb.ClientOption{
		ReadConsistency: tcvectordb.StrongConsistency})
	if err != nil {
		return nil, err
	}
	// disable/enable http request log print
	// cli.Debug(false)
	return &Demo{client: cli}, nil
}

func (d *Demo) DeleteAndDrop(ctx context.Context, database, collection string) error {
	log.Println("--------------------------- DropDatabase ---------------------------")
	dbDropResult, err := d.client.DropDatabase(ctx, database)
	if err != nil {
		return err
	}
	log.Printf("drop database result: %+v", dbDropResult)
	return nil
}

func (d *Demo) CreateDBAndCollection(ctx context.Context, database, collection string) error {
	log.Println("-------------------------- CreateDatabaseIfNotExists --------------------------")
	db, err := d.client.CreateDatabaseIfNotExists(ctx, database)
	if err != nil {
		return err
	}

	log.Println("------------------------- CreateCollectionIfNotExists -------------------------")
	index := tcvectordb.Indexes{}
	index.VectorIndex = append(index.VectorIndex, tcvectordb.VectorIndex{
		FilterIndex: tcvectordb.FilterIndex{
			FieldName: "vector",
			FieldType: tcvectordb.Vector,
			IndexType: tcvectordb.HNSW,
		},
		Dimension:  768,
		MetricType: tcvectordb.IP,
		Params: &tcvectordb.HNSWParam{
			M:              16,
			EfConstruction: 200,
		},
	})
	index.FilterIndex = append(index.FilterIndex, tcvectordb.FilterIndex{FieldName: "id", FieldType: tcvectordb.String, IndexType: tcvectordb.PRIMARY})
	index.FilterIndex = append(index.FilterIndex, tcvectordb.FilterIndex{FieldName: "file_name", FieldType: tcvectordb.String, IndexType: tcvectordb.FILTER})

	db.WithTimeout(time.Second * 30)
	_, err = db.CreateCollectionIfNotExists(ctx, collection, 3, 1, "test collection", index)
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

func (d *Demo) UploadFile(ctx context.Context, database, collection, localFilePath string) error {
	appendKeywordsToChunk := false
	appendTitleToChunk := false

	// filename := filepath.Base(localFilePath)
	// fd, err := os.Open(localFilePath)
	// if err != nil {
	// 	return err
	// }

	param := tcvectordb.UploadFileParams{
		LocalFilePath: localFilePath,
		//FileName: filename,
		//Reader:   fd,
		SplitterPreprocess: ai_document_set.DocumentSplitterPreprocess{
			AppendKeywordsToChunk: &appendKeywordsToChunk,
			AppendTitleToChunk:    &appendTitleToChunk,
		},
		ParsingProcess: &api.ParsingProcess{
			ParsingType: "AlgorithmParsing",
		},
		EmbeddingModel: "bge-base-zh",
		FieldMappings: map[string]string{
			"filename":  "file_name",
			"text":      "text",
			"imageList": "image_list",
		},
		MetaData: map[string]interface{}{
			"author": "sam",
		},
	}
	result, err := d.client.UploadFile(ctx, database, collection, param)
	if err != nil {
		return err
	}
	log.Printf("UploadFile result: %+v", result)
	return nil
}

func (d *Demo) QueryData(ctx context.Context, database, collection, filename string) error {
	time.Sleep(10 * time.Second)
	log.Println("------------------------------ Query after waiting 10s to parse file ------------------------------")

	result, err := d.client.Query(ctx, database, collection, []string{}, &tcvectordb.QueryDocumentParams{
		Filter:         tcvectordb.NewFilter(`file_name="` + filename + `"`),
		RetrieveVector: false,
		Limit:          1000,
		OutputFields:   []string{"id", "file_name", "name"},
	})
	if err != nil {
		return err
	}
	ids := []string{}
	for _, doc := range result.Documents {
		log.Printf("QueryDocument: %+v", doc)
		ids = append(ids, doc.Id)
	}

	res, err := d.client.GetImageUrl(ctx, database, collection, tcvectordb.GetImageUrlParams{
		FileName:    filename,
		DocumentIds: ids,
	})
	if err != nil {
		return err
	}
	for _, docImages := range res.Images {
		for _, docImage := range docImages {
			log.Printf("docImage: %+v", docImage)
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
	collectionName := "go-sdk-demo-col-test4"

	_, filePath, _, _ := runtime.Caller(0)
	localFilePath := path.Join(path.Dir(filePath), "../demo_files/demo_pdf_image2text_search.pdf")
	filename := filepath.Base(localFilePath)

	ctx := context.Background()
	testVdb, err := NewDemo("vdb http url or ip and port", "vdb username", "key get from web console")
	err = testVdb.CreateDBAndCollection(ctx, database, collectionName)
	printErr(err)
	err = testVdb.UploadFile(ctx, database, collectionName, localFilePath)
	printErr(err)
	err = testVdb.QueryData(ctx, database, collectionName, filename)
	printErr(err)
	err = testVdb.DeleteAndDrop(ctx, database, collectionName)
	printErr(err)
}
