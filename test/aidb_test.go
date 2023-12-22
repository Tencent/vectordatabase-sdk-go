// Copyright (C) 2023 Tencent Cloud.
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the vectordb-sdk-java), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is furnished
// to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
// PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
// SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/ai_document_set"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/collection_view"
)

var (
	aiDatabase         = "go-sdk-test-ai-db"
	collectionViewName = "go-sdk-test-ai-coll"
)

func TestAIDropDatabase(t *testing.T) {
	result, err := cli.DropAIDatabase(ctx, aiDatabase)
	printErr(err)

	t.Logf("drop database result: %+v", result)
}

func TestAICreateDatabase(t *testing.T) {
	db, err := cli.CreateAIDatabase(ctx, aiDatabase)
	printErr(err)
	t.Logf("create database success, %s", db.DatabaseName)
}

func TestAICreateCollectionView(t *testing.T) {
	db := cli.AIDatabase(aiDatabase)

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

	coll, err := db.CreateCollectionView(ctx, collectionViewName, tcvectordb.CreateCollectionViewParams{
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
	})
	printErr(err)
	log.Printf("CreateCollectionView success: %v: %v", coll.DatabaseName, coll.CollectionViewName)
}

func TestAIListCollectionViews(t *testing.T) {
	db := cli.AIDatabase(aiDatabase)
	t.Logf("ListCollectionViews ================")
	coll, err := db.ListCollectionViews(ctx)
	printErr(err)
	for _, col := range coll.CollectionViews {
		t.Logf("%+v", col)
	}
	t.Logf("DescribeCollectionView ================")
	col, err := db.DescribeCollectionView(ctx, collectionViewName)
	printErr(err)
	t.Logf("%+v", col)
}

func TestAIAlias(t *testing.T) {
	db := cli.AIDatabase(aiDatabase)
	_, err := db.SetAlias(ctx, collectionViewName, collectionAlias)
	printErr(err)

	// 查看 CollectionView 信息
	colRes, err := db.DescribeCollectionView(ctx, collectionViewName)
	printErr(err)
	t.Logf("%+v", colRes)

	// 删除 CollectionView 的 alias
	db.DeleteAlias(ctx, collectionAlias)

	// 查看 CollectionView 信息
	colRes, err = db.DescribeCollectionView(ctx, collectionViewName)
	printErr(err)
	t.Logf("%+v", colRes)
}

func TestDropCollectionView(t *testing.T) {
	res, err := cli.AIDatabase(aiDatabase).DropCollectionView(ctx, collectionViewName)
	printErr(err)
	t.Logf("%v", res)
	coll, err := cli.AIDatabase(aiDatabase).ListCollectionViews(ctx)
	printErr(err)
	t.Log("list collectionViews:")
	for _, col := range coll.CollectionViews {
		t.Logf("%+v", col)
	}
}

func TestGetCosSecret(t *testing.T) {
	res, err := cli.AIDatabase(aiDatabase).CollectionView(collectionViewName).GetCosTmpSecret(ctx, tcvectordb.GetCosTmpSecretParams{
		"tcvdb.md",
	})
	printErr(err)
	t.Logf("%+v", res)
}

func TestLoadAndSplitText(t *testing.T) {
	defer cli.Close()
	col := cli.AIDatabase(aiDatabase).CollectionView(collectionViewName)

	metaData := map[string]interface{}{
		// 元数据只支持string、uint64类型的值
		"author_name": "sam",
		"fileKey":     1024}

	// fd, err := os.Open("../example/tcvdb.md")
	// if err != nil {
	// 	t.Log(err)
	// 	return
	// }
	// defer fd.Close()
	appendTitleToChunk := false
	appendKeywordsToChunk := true
	chunkSplitter := "\n\n"

	result, err := col.LoadAndSplitText(ctx, tcvectordb.LoadAndSplitTextParams{
		// DocumentSetName: "tcvdb.md",
		// Reader:          fd,
		LocalFilePath: "../example/tcvdb.md",
		MetaData:      metaData,
		SplitterPreprocess: ai_document_set.DocumentSplitterPreprocess{
			ChunkSplitter:         &chunkSplitter,
			AppendTitleToChunk:    &appendTitleToChunk,
			AppendKeywordsToChunk: &appendKeywordsToChunk,
		},
	})
	printErr(err)
	t.Logf("%+v", result)
}

func TestAIGetDocumentSet(t *testing.T) {
	col := cli.AIDatabase(aiDatabase).CollectionView(collectionViewName)
	t.Logf("==============================GetDocumentSetByName==============================")
	res, err := col.GetDocumentSetByName(ctx, "tcvdb.md")
	printErr(err)
	t.Logf("document: %+v", ToJson(res))

	t.Logf("==============================GetDocumentSetById==============================")
	res, err = col.GetDocumentSetById(ctx, res.DocumentSetId)
	printErr(err)
	t.Logf("document: %+v", ToJson(res))
}

func TestAIGetDocumentSetChunks(t *testing.T) {
	col := cli.AIDatabase(aiDatabase).CollectionView(collectionViewName)
	t.Logf("==============================GetChunks==============================")
	result, err := col.GetChunks(ctx, tcvectordb.GetAIDocumentSetChunksParams{
		DocumentSetName: "tcvdb.md",
	})
	printErr(err)
	log.Printf("GetChunks, count: %v", result.Count)
	for _, chunk := range result.Chunks {
		log.Printf("chunk: %+v", chunk)
	}
}

func ToJson(any interface{}) string {
	bytes, err := json.Marshal(any)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func TestAIQuery(t *testing.T) {
	col := cli.AIDatabase(aiDatabase).CollectionView(collectionViewName)
	param := tcvectordb.QueryAIDocumentSetParams{
		DocumentSetName: []string{"tcvdb.md"},
		Filter:          tcvectordb.NewFilter(`documentSetName="tcvdb.md"`),
		Limit:           3,
		Offset:          0,
		// 使用OutputFields一定会输出documentSetId、documentSetName便于后续操作
		// OutputFields: []string{"indexedStatus", "textPrefix"},
	}
	result, err := col.Query(ctx, param)
	printErr(err)
	t.Logf("total doc: %d", result.Count)
	for _, doc := range result.Documents {
		b, err := json.Marshal(doc)
		if err != nil {
			return
		}
		fmt.Println(fmt.Sprintf("res %v", string(b)))
	}
}

func TestAISearch(t *testing.T) {
	col := cli.AIDatabase(aiDatabase).CollectionView(collectionViewName)

	// enableRerank := true
	searchRes, err := col.Search(ctx, tcvectordb.SearchAIDocumentSetsParams{
		Content: "什么是向量数据库",
		// FileName: "README.md",
		Filter: nil, // 过滤获取到结果
		// Limit:  3,   // 指定 Top K 的 K 值
		// RerankOption: &ai_document.RerankOption{
		// 	Enable:                &enableRerank,
		// 	ExpectRecallMultiples: 2.5,
		// },
	})
	printErr(err)
	for _, doc := range searchRes.Documents {
		t.Logf("document: %+v", doc)
	}
}

func TestAIUpdate(t *testing.T) {
	fileName := "tcvdb.md"
	col := cli.AIDatabase(aiDatabase).CollectionView(collectionViewName)
	result, err := col.Update(ctx, map[string]interface{}{
		"author_name": "jack",
	}, tcvectordb.UpdateAIDocumentSetParams{
		DocumentSetName: []string{fileName},
	})
	printErr(err)
	t.Logf("affect count: %d", result.AffectedCount)
	docs, err := col.Query(ctx, tcvectordb.QueryAIDocumentSetParams{
		Limit: 1,
	})
	printErr(err)
	for _, doc := range docs.Documents {
		t.Logf("query document is: %+v", doc)
	}
}

func TestAIDelete(t *testing.T) {
	documentSetName := "tcvdb.md"
	// documentSetId := "1177451546364084224"
	col := cli.AIDatabase(aiDatabase).CollectionView(collectionViewName)
	result, err := col.Delete(ctx, tcvectordb.DeleteAIDocumentSetParams{
		DocumentSetNames: []string{documentSetName},
		// DocumentSetIds: []string{documentSetId},
	})
	printErr(err)
	t.Logf("%v", result)
}

func TestDocumentSetSearch(t *testing.T) {
	ds, err := cli.AIDatabase(aiDatabase).CollectionView(collectionViewName).GetDocumentSetByName(ctx, "tcvdb.md")
	printErr(err)
	searchRes, err := ds.Search(ctx, tcvectordb.SearchAIDocumentSetParams{
		Content: "什么是向量数据库",
	})
	printErr(err)
	for _, doc := range searchRes.Documents {
		t.Logf("document: %+v", doc)
	}
}

func TestDocumentSetDelete(t *testing.T) {
	ds, err := cli.AIDatabase(aiDatabase).CollectionView(collectionViewName).GetDocumentSetByName(ctx, "tcvdb.md")
	printErr(err)
	res, err := ds.Delete(ctx)
	printErr(err)
	t.Logf("delete documentset result: %+v", res)
}

func TestAITruncate(t *testing.T) {
	db := cli.AIDatabase(aiDatabase)
	result, err := db.TruncateCollectionView(ctx, collectionViewName)
	printErr(err)
	t.Logf("result: %+v", result)
}
