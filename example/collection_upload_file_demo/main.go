package main

import (
	"context"
	"log"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/ai_document_set"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/document"
)

type Demo struct {
	client *tcvectordb.Client
}

func NewDemo(url, username, key string) (*Demo, error) {
	cli, err := tcvectordb.NewClient(url, username, key, &tcvectordb.ClientOption{
		ReadConsistency: tcvectordb.StrongConsistency})
	if err != nil {
		return nil, err
	}
	// disable/enable request log print
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
	index.FilterIndex = append(index.FilterIndex, tcvectordb.FilterIndex{FieldName: "chunk_num", FieldType: tcvectordb.Uint64, IndexType: tcvectordb.FILTER})
	index.FilterIndex = append(index.FilterIndex, tcvectordb.FilterIndex{FieldName: "section_num", FieldType: tcvectordb.Uint64, IndexType: tcvectordb.FILTER})

	ebd := &tcvectordb.Embedding{VectorField: "vector", Field: "text", ModelName: "bge-base-zh"}

	db.WithTimeout(time.Second * 30)
	_, err = db.CreateCollectionIfNotExists(ctx, collection, 3, 1, "test collection", index, &tcvectordb.CreateCollectionParams{
		Embedding: ebd,
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *Demo) UploadFile(ctx context.Context, database, collection, localFilePath string) error {
	appendKeywordsToChunk := true
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
			"filename":   "file_name",
			"text":       "text",
			"imageList":  "image_list",
			"chunkNum":   "chunk_num",
			"sectionNum": "section_num",
		},
		MetaData: map[string]interface{}{
			"testStr":    "v1",
			"testInt":    1024,
			"testDouble": 0.1,
			"testArray":  []string{"one", "two"},
			"testJson": map[string]interface{}{
				"a": 1,
				"b": "str",
			},
		},
	}

	_, err := d.client.UploadFile(ctx, database, collection, param)
	if err != nil {
		log.Printf("UploadFile err: %+v", err.Error())
		return err
	}

	return nil
}

func (d *Demo) QueryData(ctx context.Context, database, collection, filename string) error {
	time.Sleep(15 * time.Second)
	log.Println("------------------------------ Query file details after waiting 15s to parse file ------------------------------")
	//limit := int64(2)
	fileDetialRes, err := d.client.QueryFileDetails(ctx, database, collection, &tcvectordb.QueryFileDetailsParams{
		FileNames: []string{filename},
		// Filter:       tcvectordb.NewFilter(`_indexed_status = \"Ready\"`),
		// Limit:        &limit,
		// Offset:       0,
		// OutputFields: []string{"id", "_indexed_status", "_user_metadata"},
	})
	if err != nil {
		return err
	}
	for _, doc := range fileDetialRes.Documents {
		log.Printf("File detail: %+v", doc)
	}

	log.Println("------------------------------ Query after waiting 15s to parse file ------------------------------")

	result, err := d.client.Query(ctx, database, collection, []string{}, &tcvectordb.QueryDocumentParams{
		Filter: tcvectordb.NewFilter(`file_name="` + filename + `"`),
		Limit:  2,
		Sort: []document.SortRule{
			{
				FieldName: "chunk_num",
				Direction: "asc",
			},
		},
	})
	if err != nil {
		return err
	}
	ids := []string{}
	for _, doc := range result.Documents {
		log.Printf("File chunk: %+v", doc)
		ids = append(ids, doc.Id)
	}

	if len(ids) == 0 {
		return nil
	}

	log.Println("------------------------------ Get file chunks' imageUrls ------------------------------")
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

	log.Println("------------------------------ Query file neighbor chunks ------------------------------")
	if len(result.Documents) != 2 {
		return nil
	}
	chunkDoc := result.Documents[1]
	if _, ok := chunkDoc.Fields["chunk_num"]; !ok || chunkDoc.Fields["chunk_num"].Type() != tcvectordb.Uint64 {
		return nil
	}

	chunkNum := chunkDoc.Fields["chunk_num"].Uint64()
	leftChunkNum := uint64(0)
	if chunkNum >= 2 {
		leftChunkNum = chunkNum - 2
	}
	rightChunkNum := chunkNum + 2

	result, err = d.client.Query(ctx, database, collection, []string{}, &tcvectordb.QueryDocumentParams{
		Filter: tcvectordb.NewFilter(`file_name="` + filename + `"`).And(`chunk_num>=` +
			strconv.Itoa(int(leftChunkNum)) + ` and chunk_num<=` + strconv.Itoa(int(rightChunkNum))),
		Limit: 10,
		Sort: []document.SortRule{
			{
				FieldName: "chunk_num",
				Direction: "asc",
			},
		},
	})
	if err != nil {
		return err
	}

	for _, doc := range result.Documents {
		log.Printf("File expand chunk: %+v", doc)
	}

	log.Println("------------------------------ Query file neighbor chunks  with same section ------------------------------")
	if _, ok := chunkDoc.Fields["section_num"]; !ok || chunkDoc.Fields["section_num"].Type() != tcvectordb.Uint64 {
		return nil
	}
	sectionNum := chunkDoc.Fields["section_num"].Uint64()
	result, err = d.client.Query(ctx, database, collection, []string{}, &tcvectordb.QueryDocumentParams{
		Filter: tcvectordb.NewFilter(`file_name="` + filename + `"`).And(`chunk_num>=` +
			strconv.Itoa(int(leftChunkNum)) + ` and chunk_num<=` + strconv.Itoa(int(rightChunkNum))).And(
			`section_num=` + strconv.Itoa(int(sectionNum))),
		Limit: 10,
		Sort: []document.SortRule{
			{
				FieldName: "chunk_num",
				Direction: "asc",
			},
		},
	})
	if err != nil {
		return err
	}

	for _, doc := range result.Documents {
		log.Printf("File expand chunk with same section: %+v", doc)
	}

	log.Println("------------------------------ search file chunks by text  ------------------------------")
	filter := tcvectordb.NewFilter(`file_name="` + filename + `"`)
	searchResult, err := d.client.SearchByText(ctx, database, collection, map[string][]string{"text": {"商标声明"}},
		&tcvectordb.SearchDocumentParams{
			Filter: filter,
			Params: &tcvectordb.SearchDocParams{Ef: 200}, // 若使用HNSW索引，则需要指定参数ef，ef越大，召回率越高，但也会影响检索速度
			Limit:  2,                                    // 指定 Top K 的 K 值
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

func printErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	database := "test-db"
	collectionName := "test-coll"

	_, filePath, _, _ := runtime.Caller(0)
	localFilePath := path.Join(path.Dir(filePath), "../demo_files/tcvdb.md")
	filename := filepath.Base(localFilePath)

	ctx := context.Background()
	testVdb, err := NewDemo("vdb http url or ip and port", "vdb username", "key get from web console")
	printErr(err)
	err = testVdb.DeleteAndDrop(ctx, database, collectionName)
	printErr(err)
	err = testVdb.CreateDBAndCollection(ctx, database, collectionName)
	printErr(err)
	err = testVdb.UploadFile(ctx, database, collectionName, localFilePath)
	printErr(err)
	err = testVdb.QueryData(ctx, database, collectionName, filename)
	printErr(err)

}
