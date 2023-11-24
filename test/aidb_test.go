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
	"log"
	"testing"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb"
	collection_view "git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb/api/collection_view"
)

var (
	aiDatabase         = "go-sdk-test-ai-db"
	CollectionViewName = "go-sdk-test-ai-coll"
)

func TestAICreateDatabase(t *testing.T) {
	db, err := cli.CreateAIDatabase(ctx, aiDatabase)
	printErr(err)
	t.Logf("create database success, %s", db.DatabaseName)
}

func TestAIDropDatabase(t *testing.T) {
	result, err := cli.DropAIDatabase(ctx, aiDatabase)
	printErr(err)

	t.Logf("drop database result: %+v", result)
}

func TestAICreateCollectionView(t *testing.T) {
	db := cli.AIDatabase(aiDatabase)

	index := tcvectordb.Indexes{
		FilterIndex: []tcvectordb.FilterIndex{
			{
				FieldName: "author",
				FieldType: tcvectordb.String,
				IndexType: tcvectordb.FILTER,
			},
		},
	}

	enableWordsEmbedding := true
	appendTitleToChunk := true
	appendKeywordsToChunk := false

	coll, err := db.CreateCollectionView(ctx, CollectionViewName, &tcvectordb.CreateCollectionViewParams{
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
	col, err := db.DescribeCollectionView(ctx, CollectionViewName)
	printErr(err)
	t.Logf("%+v", col)
}

func TestAIAlias(t *testing.T) {
	db := cli.AIDatabase(aiDatabase)
	_, err := db.SetAlias(ctx, CollectionViewName, collectionAlias)
	printErr(err)

	// 查看 CollectionView 信息
	colRes, err := db.DescribeCollectionView(ctx, CollectionViewName)
	printErr(err)
	t.Logf("%+v", colRes)

	// 删除 CollectionView 的 alias
	db.DeleteAlias(ctx, collectionAlias)

	// 查看 CollectionView 信息
	colRes, err = db.DescribeCollectionView(ctx, CollectionViewName)
	printErr(err)
	t.Logf("%+v", colRes)
}

func TestDropCollectionView(t *testing.T) {
	res, err := cli.AIDatabase(aiDatabase).DropCollectionView(ctx, CollectionViewName)
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
	res, err := cli.AIDatabase(aiDatabase).CollectionView(CollectionViewName).GetCosTmpSecret(ctx, tcvectordb.GetCosTmpSecretParams{
		"tcvdb.md",
	})
	printErr(err)
	t.Logf("%+v", res)
}

func TestLoadAndSplitText(t *testing.T) {
	defer cli.Close()
	col := cli.AIDatabase(aiDatabase).CollectionView(CollectionViewName)

	metaData := map[string]tcvectordb.Field{
		"author":  {Val: "sam"},
		"fileKey": {Val: 1024}}
	result, err := col.LoadAndSplitText(ctx, tcvectordb.LoadAndSplitTextParams{
		LocalFilePath: "../example/tcvdb.md",
		MetaData:      metaData,
	})
	printErr(err)
	t.Logf("%+v", result)
}

func TestAIGetDocumentSet(t *testing.T) {
	col := cli.AIDatabase(aiDatabase).CollectionView(CollectionViewName)
	t.Logf("==============================GetDocumentSetByName==============================")
	res, err := col.GetDocumentSetByName(ctx, "tcvdb.md")
	printErr(err)
	t.Logf("document: %+v", res)

	t.Logf("==============================GetDocumentSetById==============================")
	res, err = col.GetDocumentSetById(ctx, res.DocumentSets.DocumentSetId)
	printErr(err)
	t.Logf("document: %+v", res)
}

func TestAIQuery(t *testing.T) {
	col := cli.AIDatabase(aiDatabase).CollectionView(CollectionViewName)
	param := tcvectordb.QueryAIDocumentSetParams{
		DocumentSetName: []string{"tcvdb.md"},
		// Filter: tcvectordb.NewFilter(`documentSetName="tcvdb.md"`),
		Limit:  3,
		Offset: 0,
	}
	result, err := col.Query(ctx, param)
	printErr(err)
	t.Logf("total doc: %d", result.Count)
	for _, doc := range result.Documents {
		t.Logf("document: %+v", doc)
	}
}

func TestAISearch(t *testing.T) {
	col := cli.AIDatabase(aiDatabase).CollectionView(CollectionViewName)

	// enableRerank := true
	searchRes, err := col.Search(ctx, tcvectordb.SearchAIDocumentSetParams{
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
	col := cli.AIDatabase(aiDatabase).CollectionView(CollectionViewName)
	result, err := col.Update(ctx, map[string]interface{}{
		"author": "jack",
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
	col := cli.AIDatabase(aiDatabase).CollectionView(CollectionViewName)
	result, err := col.Delete(ctx, tcvectordb.DeleteAIDocumentSetParams{
		DocumentSetNames: []string{documentSetName},
		// DocumentSetIds: []string{documentSetId},
	})
	printErr(err)
	t.Logf("%v", result)
}

func TestAITruncate(t *testing.T) {
	db := cli.AIDatabase(aiDatabase)
	result, err := db.TruncateCollectionView(ctx, CollectionViewName)
	printErr(err)
	t.Logf("result: %+v", result)
}
