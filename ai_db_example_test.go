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
	"encoding/json"
	"testing"
	"time"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity/api/document"
)

var (
	AIDatabase       = "ai-db-test"
	AICollectionName = "user_collection"
)

func TestCreateAIDatabase(t *testing.T) {
	db, err := cli.CreateAiDatabase(context.Background(), AIDatabase, nil)
	printErr(err)
	t.Logf("create database success, %s", db.DatabaseName)
}

func TestListDatabase(t *testing.T) {
	dbList, err := cli.ListDatabase(context.Background(), nil)
	printErr(err)

	for _, db := range dbList {
		t.Logf("database: %s, createTime: %s, dbType: %s", db.DatabaseName, db.Info.CreateTime, db.Info.DbType)
	}
}

func TestDropAIDatabase(t *testing.T) {
	result, err := cli.DropAiDatabase(context.Background(), AIDatabase, nil)
	printErr(err)

	t.Logf("drop collection result: %+v", result)

}

func TestCreateCollectionInAIDB(t *testing.T) {
	db := cli.Database(database)

	// 创建 Collection

	// 第一步，设计索引（不是设计 Collection 的结构）
	// 1. 【必须的索引】：主键id字段目前是固定且必须的，参考下面的例子；
	// 2. 【其他索引】：检索时需作为条件查询的字段，比如要按书籍的作者进行过滤，这个时候 author 字段就需要建立索引，
	//     否则无法在查询的时候对 author 字段进行过滤，不需要过滤的字段无需加索引，会浪费内存；

	index := entity.Indexes{
		FilterIndex: []entity.FilterIndex{
			{
				FieldName: "id",
				FieldType: entity.String,
				IndexType: entity.PRIMARY,
			},
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
		},
	}

	// 第二步：创建 Collection
	// 创建collection耗时较长，需要调整客户端的timeout
	db.WithTimeout(time.Second * 30)
	_, err := db.CreateCollection(context.Background(), AICollectionName, 1, 1, "test collection", index, nil)
	printErr(err)

	// 列出所有 Collection
	colList, err := db.ListCollection(context.Background(), nil)
	printErr(err)
	for _, col := range colList {
		t.Logf("%+v", col)
	}

	// 设置 Collection 的 alias
	_, err = db.SetAlias(context.Background(), collectionName, collectionAlias, nil)
	printErr(err)

	// 查看 Collection 信息
	colRes, err := db.DescribeCollection(context.Background(), collectionName, nil)
	printErr(err)
	t.Logf("%+v", colRes)

	// 删除 Collection 的 alias
	db.DeleteAlias(context.Background(), collectionAlias, nil)
}

func TestUploadFile(t *testing.T) {
	defer cli.Close()
	col := cli.Database(AIDatabase).Collection(AICollectionName)
	cli.Debug(true)

	metaData := map[string]entity.Field{
		"fileName": {Val: "召回率从0.48提升至0.92_back116.md"},
		"author":   {Val: "sam"},
		"fileKey":  {Val: 1024}}
	err := col.Upload(context.Background(), "./召回率从0.48提升至0.92_back116.md", &entity.UploadDocumentOption{
		FileType: "", MetaData: metaData})
	printErr(err)
}

func TestAIQuery(t *testing.T) {
	// 查询
	// 1. query 用于查询数据
	// 2. 可以通过传入主键 id 列表或 filter 实现过滤数据的目的
	// 3. 如果没有主键 id 列表和 filter 则必须传入 limit 和 offset，类似 scan 的数据扫描功能
	// 4. 如果仅需要部分 field 的数据，可以指定 output_fields 用于指定返回数据包含哪些 field，不指定默认全部返回

	col := cli.Database(AIDatabase).Collection(AICollectionName)
	option := &entity.QueryDocumentOption{
		Filter:         nil,
		OutputFields:   []string{},
		RetrieveVector: false,
		Limit:          3,
		Offset:         0,
	}
	col.Debug(true)
	docs, result, err := col.Query(context.Background(), []string{"0001", "0002", "0003", "0004", "0005"}, option)
	printErr(err)
	t.Logf("total doc: %d", result.Total)
	for _, doc := range docs {
		t.Logf("id: %s, vector: %v, field: %+v", doc.Id, doc.Vector, doc.Fields)
	}
}

func Test_json(t *testing.T) {
	str := "{\"id\":\"0001\",\"_indexed_status\":0,\"_cos_address\":\"\",\"_indexed\":0,\"_file_name\":\"罗贯中\"}"
	var doc document.Document
	json.Unmarshal([]byte(str), &doc)

	t.Logf("doc id %v doc field %v", doc.Id, doc.Fields)

}
