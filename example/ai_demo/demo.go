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
	collectionView := "go-sdk-demo-ai-col"
	collectionViewAlias := "go-sdk-demo-ai-alias"

	ctx := context.Background()
	testVdb, err := example.NewAIDemo("vdb http url or ip and post", "vdb username", "key get from web console")
	printErr(err)
	err = testVdb.Clear(ctx, database)
	printErr(err)
	err = testVdb.CreateAIDatabase(ctx, database)
	printErr(err)
	err = testVdb.CreateCollectionView(ctx, database, collectionView)
	printErr(err)
	loadFileRes, err := testVdb.LoadAndSplitText(ctx, database, collectionView, "../tcvdb.md")
	printErr(err)
	time.Sleep(time.Second * 30) // 等待后台解析文件完成
	err = testVdb.GetFile(ctx, database, collectionView, loadFileRes.DocumentSetName)
	printErr(err)
	err = testVdb.QueryAndSearch(ctx, database, collectionView)
	printErr(err)
	err = testVdb.Alias(ctx, database, collectionView, collectionViewAlias)
	printErr(err)
	err = testVdb.DeleteAndDrop(ctx, database, collectionView, loadFileRes.DocumentSetName)
	printErr(err)
}
