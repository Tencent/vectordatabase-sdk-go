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

package main

import (
	"context"
	"testing"
	"time"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity"
)

var (
	aiDatabase       = "ai-db-test-lqs"
	aiCollectionName = "user_collection-lqs"
)

func TestCreateAIDatabase(t *testing.T) {
	db, err := cli.CreateAIDatabase(context.Background(), aiDatabase, nil)
	printErr(err)
	t.Logf("create database success, %s", db.DatabaseName)
}

func TestListDatabase(t *testing.T) {
	dbList, err := cli.ListDatabase(context.Background(), nil)
	printErr(err)
	t.Logf("base database ======================")
	for _, db := range dbList.Databases {
		t.Logf("database: %s, createTime: %s, dbType: %s", db.DatabaseName, db.Info.CreateTime, db.Info.DbType)
	}

	t.Logf("AI database ======================")
	for _, db := range dbList.AIDatabases {
		t.Logf("database: %s, createTime: %s, dbType: %s", db.DatabaseName, db.Info.CreateTime, db.Info.DbType)
	}
}

func TestDropAIDatabase(t *testing.T) {
	result, err := cli.DropAIDatabase(context.Background(), aiDatabase, nil)
	printErr(err)

	t.Logf("drop database result: %+v", result)
}

func TestCreateCollectionInAIDB(t *testing.T) {
	db := cli.AIDatabase(aiDatabase)

	// 创建 Collection

	// 第一步，设计索引（不是设计 Collection 的结构）
	// 1. 【必须的索引】：主键id字段目前是固定且必须的，参考下面的例子；
	// 2. 【其他索引】：检索时需作为条件查询的字段，比如要按书籍的作者进行过滤，这个时候 author 字段就需要建立索引，
	//     否则无法在查询的时候对 author 字段进行过滤，不需要过滤的字段无需加索引，会浪费内存；

	index := []entity.FilterIndex{
		{
			FieldName: "bookName",
			FieldType: entity.String,
			IndexType: entity.FILTER,
		},
		{
			FieldName: "wordNum",
			FieldType: entity.Uint64,
			IndexType: entity.FILTER,
		},
	}

	// 第二步：创建 Collection
	// 创建collection耗时较长，需要调整客户端的timeout
	db.WithTimeout(time.Second * 30)
	db.Debug(true)
	_, err := db.CreateCollection(context.Background(), aiCollectionName, &entity.CreateAICollectionOption{
		Description: "test ai collection",
		Indexes:     index,
		AiConfig: &entity.AiConfig{
			MaxFiles:        1000,
			AverageFileSize: 1 << 20,
			Language:        entity.LanguageChinese,
		},
	})
	printErr(err)

	// 列出所有 Collection
	colList, err := db.ListCollection(context.Background(), nil)
	printErr(err)
	for _, col := range colList.Collections {
		t.Logf("%+v", col)
	}
}

func TestListAICollection(t *testing.T) {
	db := cli.AIDatabase(aiDatabase)
	t.Logf("ListCollection ================")
	db.Debug(true)
	coll, err := db.ListCollection(context.Background(), nil)
	printErr(err)
	for _, col := range coll.Collections {
		t.Logf("%+v", col)
	}
	t.Logf("DescribeCollection ================")
	col, err := db.DescribeCollection(context.Background(), aiCollectionName, nil)
	printErr(err)
	t.Logf("%+v", col)
}

func TestAIAlias(t *testing.T) {
	db := cli.AIDatabase(aiDatabase)
	db.Debug(true)
	_, err := db.SetAlias(context.Background(), aiCollectionName, collectionAlias, nil)
	printErr(err)

	// 查看 Collection 信息
	colRes, err := db.DescribeCollection(context.Background(), aiCollectionName, nil)
	printErr(err)
	t.Logf("%+v", colRes)

	// 删除 Collection 的 alias
	db.DeleteAlias(context.Background(), collectionAlias, nil)

	// 查看 Collection 信息
	colRes, err = db.DescribeCollection(context.Background(), aiCollectionName, nil)
	printErr(err)
	t.Logf("%+v", colRes)
}

func TestDropAICollection(t *testing.T) {
	res, err := cli.AIDatabase(aiDatabase).DropCollection(context.Background(), aiCollectionName, nil)
	printErr(err)
	t.Logf("%v", res)
	coll, err := cli.AIDatabase(aiDatabase).ListCollection(context.Background(), nil)
	printErr(err)
	for _, col := range coll.Collections {
		t.Logf("%+v", col)
	}
}

func TestGetCosSecret(t *testing.T) {
	res, err := cli.AIDatabase(aiDatabase).Collection(aiCollectionName).GetCosTmpSecret(context.Background(), "./职业规划.md", nil)
	printErr(err)
	t.Logf("%+v", res)
}

func TestUploadFile(t *testing.T) {
	defer cli.Close()
	col := cli.AIDatabase(aiDatabase).Collection(aiCollectionName)
	cli.Debug(true)

	metaData := map[string]entity.Field{
		"fileName": {Val: "职业规划.md"},
		"author":   {Val: "sam"},
		"fileKey":  {Val: 1024}}
	result, err := col.Upload(context.Background(), "./职业规划.md", &entity.UploadAIDocumentOption{
		FileType: "", MetaData: metaData})
	printErr(err)
	t.Logf("%+v", result)
}

func TestAIQuery(t *testing.T) {
	// 查询
	// 1. query 用于查询数据
	// 2. 可以通过传入主键 id 列表或 filter 实现过滤数据的目的
	// 3. 如果没有主键 id 列表和 filter 则必须传入 limit 和 offset，类似 scan 的数据扫描功能
	// 4. 如果仅需要部分 field 的数据，可以指定 output_fields 用于指定返回数据包含哪些 field，不指定默认全部返回

	col := cli.AIDatabase(aiDatabase).Collection(aiCollectionName)
	option := &entity.QueryAIDocumentOption{
		Filter:       nil,
		OutputFields: []string{},
		Limit:        3,
		Offset:       0,
	}
	col.Debug(true)
	result, err := col.Query(context.Background(), option)
	printErr(err)
	t.Logf("total doc: %d", result.Total)
	for _, doc := range result.Documents {
		t.Logf("document: %+v", doc)
	}
}

func TestAISearch(t *testing.T) {
	col := cli.AIDatabase(aiDatabase).Collection(aiCollectionName)

	// searchById
	// 1. searchById 提供按 id 搜索的能力
	// 1. search 提供按照 vector 搜索的能力
	// 2. 支持通过 filter 过滤数据
	// 3. 如果仅需要部分 field 的数据，可以指定 output_fields 用于指定返回数据包含哪些 field，不指定默认全部返回
	// 4. limit 用于限制每个单元搜索条件的条数，如 vector 传入三组向量，limit 为 3，则 limit 限制的是每组向量返回 top 3 的相似度向量

	// 根据主键 id 查找 Top K 个相似性结果，向量数据库会根据ID 查找对应的向量，再根据向量进行TOP K 相似性检索

	filter := entity.NewFilter(`bookName="三国演义"`)
	searchRes, err := col.Search(context.Background(), "文本语句", &entity.SearchAIDocumentOption{
		Filter: filter, // 过滤获取到结果
		Limit:  3,      // 指定 Top K 的 K 值
	})
	printErr(err)
	for _, doc := range searchRes.Documents {
		t.Logf("document: %+v", doc)
	}
}

func TestAIUpdateAndDelete(t *testing.T) {
	col := cli.AIDatabase(aiDatabase).Collection(aiCollectionName)
	fileId := ""
	// update
	// 1. update 提供基于 [主键查询] 和 [Filter 过滤] 的部分字段更新或者非索引字段新增

	// filter 限制仅会更新 id = "0003"
	result, err := col.Update(context.Background(), &entity.UpdateAIDocumentOption{
		QueryIds: []string{fileId},
		UpdateFields: map[string]interface{}{
			"page": 24,
		},
	})
	printErr(err)
	t.Logf("affect count: %d", result.AffectedCount)
	docs, err := col.Query(context.Background(), &entity.QueryAIDocumentOption{
		DocumentIds: []string{fileId},
	})
	printErr(err)
	for _, doc := range docs.Documents {
		t.Logf("query document is: %+v", doc)
	}

	// delete
	// 1. delete 提供基于 [主键查询] 和 [Filter 过滤] 的数据删除能力
	// 2. 删除功能会受限于 collection 的索引类型，部分索引类型不支持删除操作

	// filter 限制只会删除 id="0001" 成功
	col.Delete(context.Background(), &entity.DeleteAIDocumentOption{
		DocumentIds: []string{fileId},
	})
}
