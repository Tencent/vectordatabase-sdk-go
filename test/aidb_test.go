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
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb/api/ai_collection"
)

var (
	aiDatabase       = "lqs-db"
	aiCollectionName = "lqs-coll"
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

func TestAICreateCollection(t *testing.T) {
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

	enableWordsEmbedding := false
	coll, err := db.CreateCollection(ctx, aiCollectionName, &tcvectordb.CreateAICollectionOption{
		Description: "test ai collection",
		Indexes:     index,
		AiConfig: &tcvectordb.AiConfig{
			ExpectedFileNum:      1000,
			AverageFileSize:      1 << 20,
			Language:             tcvectordb.LanguageChinese,
			DocumentPreprocess:   &ai_collection.SplitterPreprocess{},
			EnableWordsEmbedding: &enableWordsEmbedding,
		},
	})
	printErr(err)
	log.Printf("CreateCollection success: %v: %v", coll.DatabaseName, coll.CollectionName)
}

func TestAIListCollection(t *testing.T) {
	db := cli.AIDatabase(aiDatabase)
	t.Logf("ListCollection ================")
	coll, err := db.ListCollection(ctx, nil)
	printErr(err)
	for _, col := range coll.Collections {
		t.Logf("%+v", col)
	}
	t.Logf("DescribeCollection ================")
	col, err := db.DescribeCollection(ctx, aiCollectionName, nil)
	printErr(err)
	t.Logf("%+v", col)
}

func TestAIAlias(t *testing.T) {
	db := cli.AIDatabase(aiDatabase)
	_, err := db.SetAlias(ctx, aiCollectionName, collectionAlias, nil)
	printErr(err)

	// 查看 Collection 信息
	colRes, err := db.DescribeCollection(ctx, aiCollectionName, nil)
	printErr(err)
	t.Logf("%+v", colRes)

	// 删除 Collection 的 alias
	db.DeleteAlias(ctx, collectionAlias, nil)

	// 查看 Collection 信息
	colRes, err = db.DescribeCollection(ctx, aiCollectionName, nil)
	printErr(err)
	t.Logf("%+v", colRes)
}

func TestDropAICollection(t *testing.T) {
	res, err := cli.AIDatabase(aiDatabase).DropCollection(ctx, aiCollectionName, nil)
	printErr(err)
	t.Logf("%v", res)
	coll, err := cli.AIDatabase(aiDatabase).ListCollection(ctx, nil)
	printErr(err)
	t.Log("list collection:")
	for _, col := range coll.Collections {
		t.Logf("%+v", col)
	}
}

func TestGetCosSecret(t *testing.T) {
	cli.Debug(true)
	res, err := cli.AIDatabase(aiDatabase).Collection(aiCollectionName).GetCosTmpSecret(ctx, "../example/tcvdb.md")
	printErr(err)
	t.Logf("%+v", res)
}

func TestUploadFile(t *testing.T) {
	col := cli.AIDatabase(aiDatabase).Collection(aiCollectionName)

	metaData := map[string]tcvectordb.Field{
		"author":  {Val: "sam"},
		"fileKey": {Val: 1024}}
	result, err := col.Upload(ctx, "../example/tcvdb.md", &tcvectordb.UploadAIDocumentOption{
		FileType: "", MetaData: metaData})
	printErr(err)
	t.Logf("%+v", result)
}

func TestAIGet(t *testing.T) {
	col := cli.AIDatabase(aiDatabase).Collection(aiCollectionName)
	col.Debug(true)
	res, err := col.Get(ctx, &tcvectordb.GetAIDocumentOption{
		DocumentSetName: "tcvdb.md",
	})
	printErr(err)
	t.Logf("document: %+v", res)
}

func TestAIQuery(t *testing.T) {
	col := cli.AIDatabase(aiDatabase).Collection(aiCollectionName)
	option := &tcvectordb.QueryAIDocumentOption{
		// Filter:       tcvectordb.NewFilter(`_file_name="README.md"`),
		OutputFields: []string{},
		Limit:        3,
		Offset:       0,
	}
	result, err := col.Query(ctx, option)
	printErr(err)
	t.Logf("total doc: %d", result.Total)
	for _, doc := range result.Documents {
		t.Logf("document: %+v", doc)
	}
}

func TestAISearch(t *testing.T) {
	col := cli.AIDatabase(aiDatabase).Collection(aiCollectionName)

	// enableRerank := true
	searchRes, err := col.Search(ctx, "什么是向量数据库", &tcvectordb.SearchAIDocumentOption{
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
	col := cli.AIDatabase(aiDatabase).Collection(aiCollectionName)
	result, err := col.Update(ctx, &tcvectordb.UpdateAIDocumentOption{
		FileName: fileName,
		UpdateFields: map[string]interface{}{
			"author": "jack",
		},
	})
	printErr(err)
	t.Logf("affect count: %d", result.AffectedCount)
	docs, err := col.Query(ctx, &tcvectordb.QueryAIDocumentOption{
		FileName: fileName,
		Limit:    1,
	})
	printErr(err)
	for _, doc := range docs.Documents {
		t.Logf("query document is: %+v", doc)
	}
}

func TestAIDelete(t *testing.T) {
	fileName := "tcvdb.md"
	col := cli.AIDatabase(aiDatabase).Collection(aiCollectionName)
	result, err := col.Delete(ctx, &tcvectordb.DeleteAIDocumentOption{
		// DocumentIds: []string{fileId},
		FileName: fileName,
	})
	printErr(err)
	t.Logf("%v", result)
}

func TestAITruncate(t *testing.T) {
	db := cli.AIDatabase(aiDatabase)
	result, err := db.TruncateCollection(ctx, aiCollectionName, nil)
	printErr(err)
	t.Logf("result: %+v", result)
}
