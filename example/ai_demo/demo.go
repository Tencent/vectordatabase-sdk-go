package main

import (
	"context"
	"log"
	"time"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/example"
)

func printErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	database := "go-sdk-demo-ai-db"
	collectionName := "go-sdk-demo-ai-col"
	collectionAlias := "go-sdk-demo-ai-alias"

	ctx := context.Background()
	testVdb, err := example.NewAIDemo("vdb http url or ip and post", "vdb username", "key get from web console")
	printErr(err)
	err = testVdb.Clear(ctx, database)
	printErr(err)
	err = testVdb.CreateAIDatabase(ctx, database)
	printErr(err)
	err = testVdb.CreateCollectionView(ctx, database, collectionName)
	printErr(err)
	fileInfo, err := testVdb.UploadFile(ctx, database, collectionName, "../tcvdb.md")
	printErr(err)
	time.Sleep(time.Second * 30) // 等待后台解析文件完成
	err = testVdb.GetFile(ctx, database, collectionName, fileInfo.FileName)
	printErr(err)
	err = testVdb.QueryAndSearch(ctx, database, collectionName)
	printErr(err)
	err = testVdb.Alias(ctx, database, collectionName, collectionAlias)
	printErr(err)
	err = testVdb.DeleteAndDrop(ctx, database, collectionName, fileInfo.FileName)
	printErr(err)
}
