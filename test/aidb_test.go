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

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity/api/ai_collection"
)

var (
	aiDatabase       = "go-sdk-test-ai-db"
	aiCollectionName = "go-sdk-test-ai-coll"
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

	index := entity.Indexes{
		FilterIndex: []entity.FilterIndex{
			{
				FieldName: "author",
				FieldType: entity.String,
				IndexType: entity.FILTER,
			},
		},
	}

	coll, err := db.CreateCollection(ctx, aiCollectionName, &entity.CreateAICollectionOption{
		Description: "test ai collection",
		Indexes:     index,
		AiConfig: &entity.AiConfig{
			ExpectedFileNum: 1000,
			AverageFileSize: 1 << 20,
			Language:        entity.LanguageChinese,
			DocumentPreprocess: &ai_collection.DocumentPreprocess{
				AppendKeywordsToChunk: "1",
			},
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
	res, err := cli.AIDatabase(aiDatabase).Collection(aiCollectionName).GetCosTmpSecret(ctx, "./README.md", nil)
	printErr(err)
	t.Logf("%+v", res)
}

func TestUploadFile(t *testing.T) {
	defer cli.Close()
	col := cli.AIDatabase(aiDatabase).Collection(aiCollectionName)

	metaData := map[string]entity.Field{
		"author":  {Val: "sam"},
		"fileKey": {Val: 1024}}
	result, err := col.Upload(ctx, "../example/tcvdb.md", &entity.UploadAIDocumentOption{
		FileType: "", MetaData: metaData})
	printErr(err)
	t.Logf("%+v", result)
}

func TestAIQuery(t *testing.T) {
	col := cli.AIDatabase(aiDatabase).Collection(aiCollectionName)
	option := &entity.QueryAIDocumentOption{
		// Filter:       entity.NewFilter(`_file_name="README.md"`),
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

	searchRes, err := col.Search(ctx, "什么是向量数据库", &entity.SearchAIDocumentOption{
		// FileName: "README.md",
		Filter: nil, // 过滤获取到结果
		// Limit:  3,   // 指定 Top K 的 K 值
	})
	printErr(err)
	for _, doc := range searchRes.Documents {
		t.Logf("document: %+v", doc)
	}
}

func TestAIUpdate(t *testing.T) {
	fileName := "tcvdb.md"
	col := cli.AIDatabase(aiDatabase).Collection(aiCollectionName)
	result, err := col.Update(ctx, &entity.UpdateAIDocumentOption{
		FileName: fileName,
		UpdateFields: map[string]interface{}{
			"author": "jack",
		},
	})
	printErr(err)
	t.Logf("affect count: %d", result.AffectedCount)
	docs, err := col.Query(ctx, &entity.QueryAIDocumentOption{
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
	result, err := col.Delete(ctx, &entity.DeleteAIDocumentOption{
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
