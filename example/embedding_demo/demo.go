package main

import (
	"context"
	"log"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/example"
)

func printErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	database := "go-sdk-demo-db"
	collectionName := "go-sdk-demo-em-col"
	collectionAlias := "go-sdk-demo-em-alias"

	ctx := context.Background()
	testVdb, err := example.NewEmbeddingDemo("vdb http url or ip and post", "vdb username", "key get from web console")
	printErr(err)
	err = testVdb.Clear(ctx, database)
	printErr(err)
	err = testVdb.CreateDBAndCollection(ctx, database, collectionName, collectionAlias)
	printErr(err)
	err = testVdb.UpsertData(ctx, database, collectionName)
	printErr(err)
	err = testVdb.QueryData(ctx, database, collectionName)
	printErr(err)
	err = testVdb.UpdateAndDelete(ctx, database, collectionName)
	printErr(err)
	err = testVdb.DeleteAndDrop(ctx, database, collectionName)
	printErr(err)
}
