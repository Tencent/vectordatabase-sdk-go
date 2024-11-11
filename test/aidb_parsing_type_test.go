package test

import (
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/collection_view"
)

func TestAICreateCollectionViewWithDefaultParsingType(t *testing.T) {
	parsingTypeBaseCase(nil, nil, "demo_vision_model_parsing.pdf", []string{})
}

func TestAICreateCollectionViewWithVisionModelParsing(t *testing.T) {
	parsingTypeBaseCase(&api.ParsingProcess{ParsingType: string(tcvectordb.VisionModelParsing)}, nil,
		"demo_vision_model_parsing.pdf", []string{})
}

func TestAICreateCollectionViewWithAlgorithmParsing(t *testing.T) {
	parsingTypeBaseCase(&api.ParsingProcess{ParsingType: string(tcvectordb.AlgorithmParsing)}, nil,
		"demo_vision_model_parsing.pdf", []string{})
}

func TestAIUploadDocWithAlgorithmParsing(t *testing.T) {
	parsingTypeBaseCase(&api.ParsingProcess{ParsingType: string(tcvectordb.VisionModelParsing)},
		&api.ParsingProcess{ParsingType: string(tcvectordb.AlgorithmParsing)},
		"demo_vision_model_parsing.pdf", []string{})
}

func TestAIUploadDocWithVisionModelParsing(t *testing.T) {
	parsingTypeBaseCase(&api.ParsingProcess{ParsingType: string(tcvectordb.AlgorithmParsing)},
		&api.ParsingProcess{ParsingType: string(tcvectordb.VisionModelParsing)},
		"demo_vision_model_parsing.pdf", []string{})
}

func TestAIUploadDocWithMdVisionModelParsing(t *testing.T) {
	parsingTypeBaseCase(&api.ParsingProcess{ParsingType: string(tcvectordb.AlgorithmParsing)},
		&api.ParsingProcess{ParsingType: string(tcvectordb.VisionModelParsing)},
		"tcvdb.md", []string{})
}

func TestAIUploadDocWithMdOutputFields_case1(t *testing.T) {
	parsingTypeBaseCase(&api.ParsingProcess{ParsingType: string(tcvectordb.AlgorithmParsing)},
		&api.ParsingProcess{ParsingType: string(tcvectordb.VisionModelParsing)},
		"tcvdb.md", []string{"parsingProcess"})
}

func TestAIUploadDocWithMdOutputFields_case2(t *testing.T) {
	parsingTypeBaseCase(&api.ParsingProcess{ParsingType: string(tcvectordb.AlgorithmParsing)},
		&api.ParsingProcess{ParsingType: string(tcvectordb.VisionModelParsing)},
		"tcvdb.md", []string{"documentSetId"})
}

func TestAIUploadDocWithMdOutputFields_case3(t *testing.T) {
	parsingTypeBaseCase(&api.ParsingProcess{ParsingType: string(tcvectordb.AlgorithmParsing)},
		&api.ParsingProcess{ParsingType: string(tcvectordb.VisionModelParsing)},
		"tcvdb.md", []string{"parsingType"})
}

func TestAIUploadDocWithMdOutputFields_case4(t *testing.T) {
	parsingTypeBaseCase(&api.ParsingProcess{ParsingType: string(tcvectordb.AlgorithmParsing)},
		&api.ParsingProcess{ParsingType: string(tcvectordb.VisionModelParsing)},
		"tcvdb.md", []string{"parsingType", "splitterPreprocess"})
}

func TestAIUploadDocWithMdWrongType(t *testing.T) {
	parsingTypeBaseCase(&api.ParsingProcess{ParsingType: string(tcvectordb.AlgorithmParsing)},
		&api.ParsingProcess{ParsingType: "hello"},
		"tcvdb.md", []string{})
}

func TestGetFile(t *testing.T) {
	documentName := "demo.pdf"
	coll := cli.AIDatabase("db_xx").CollectionView("coll_xx")
	log.Println("---------------------------- query documentSet ---------------------------")
	param := tcvectordb.QueryAIDocumentSetParams{
		DocumentSetName: []string{documentName},
		//Filter:          tcvectordb.NewFilter(`documentSetName="tcvdb.md"`),
		Limit:  3,
		Offset: 0,
		// 使用OutputFields一定会输出documentSetId、documentSetName便于后续操作
	}
	queryResult, err := coll.Query(ctx, param)
	printErr(err)
	log.Printf("total doc: %d", queryResult.Count)
	for _, doc := range queryResult.Documents {
		println(ToJson(doc))
	}

	log.Println("---------------------------- Get documentSet ---------------------------")

	getResult, err := coll.GetDocumentSetByName(ctx, documentName)
	printErr(err)
	println(ToJson(getResult))
}

func parsingTypeBaseCase(collParsingProcess *api.ParsingProcess, docParPro *api.ParsingProcess,
	documentName string, outputField []string) {
	aiDatabase := "db_" + strconv.FormatInt(time.Now().UnixMicro(), 10)
	_, err := cli.DropAIDatabase(ctx, aiDatabase)
	printErr(err)

	db, err := cli.CreateAIDatabase(ctx, aiDatabase)
	printErr(err)
	log.Printf("create database success, %s", db.DatabaseName)

	log.Println("---------------------------- Create CollectionView ---------------------------")
	index := tcvectordb.Indexes{
		FilterIndex: []tcvectordb.FilterIndex{
			{
				FieldName: "author_name",
				FieldType: tcvectordb.String,
				IndexType: tcvectordb.FILTER,
			},
		},
	}

	enableWordsEmbedding := true
	appendTitleToChunk := true
	appendKeywordsToChunk := false

	collectionName := "coll_" + strconv.FormatInt(time.Now().UnixMicro(), 10)

	coll, err := db.CreateCollectionView(ctx, collectionName, tcvectordb.CreateCollectionViewParams{
		Description: "test ai collectionView",
		Indexes:     index,
		Embedding: &collection_view.DocumentEmbedding{
			Language:             string(tcvectordb.LanguageChinese),
			EnableWordsEmbedding: &enableWordsEmbedding,
		},
		SplitterPreprocess: &collection_view.SplitterPreprocess{
			AppendTitleToChunk:    &appendTitleToChunk,
			AppendKeywordsToChunk: &appendKeywordsToChunk,
		},
		ParsingProcess:  collParsingProcess,
		ExpectedFileNum: 204800,
		AverageFileSize: 10240,
	})
	printErr(err)
	log.Printf("CreateCollectionView success: %v: %v", coll.DatabaseName, coll.CollectionViewName)

	log.Println("---------------------------- List CollectionView ---------------------------")

	colls, err := db.ListCollectionViews(ctx)
	printErr(err)
	for _, col := range colls.CollectionViews {
		log.Printf("%+v", ToJson(col))
	}
	log.Println("---------------------------- Desc CollectionView ---------------------------")
	col, err := db.DescribeCollectionView(ctx, collectionName)
	printErr(err)
	log.Printf("%+v", ToJson(col))

	log.Println("---------------------------- LoadAndSplitText documentSet ---------------------------")
	metaData := map[string]interface{}{
		// 元数据只支持string、uint64类型的值
		"author_name": "sam",
		"fileKey":     1024}

	result, err := col.LoadAndSplitText(ctx, tcvectordb.LoadAndSplitTextParams{
		// DocumentSetName: "tcvdb.md",
		// Reader:          fd,
		LocalFilePath:  "../example/demo_files/" + documentName,
		MetaData:       metaData,
		ParsingProcess: docParPro,
	})
	printErr(err)
	log.Printf("%+v", result)

	time.Sleep(10 * time.Second)

	log.Println("---------------------------- query documentSet ---------------------------")
	param := tcvectordb.QueryAIDocumentSetParams{
		DocumentSetName: []string{documentName},
		//Filter:          tcvectordb.NewFilter(`documentSetName="tcvdb.md"`),
		Limit:  3,
		Offset: 0,
		// 使用OutputFields一定会输出documentSetId、documentSetName便于后续操作
		OutputFields: outputField,
	}
	queryResult, err := coll.Query(ctx, param)
	printErr(err)
	log.Printf("total doc: %d", queryResult.Count)
	for _, doc := range queryResult.Documents {
		println(ToJson(doc))
	}

	log.Println("---------------------------- Get documentSet ---------------------------")

	getResult, err := coll.GetDocumentSetByName(ctx, documentName)
	printErr(err)
	println(ToJson(getResult))

}
