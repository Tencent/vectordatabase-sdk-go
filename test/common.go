package test

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
)

var (
	cli *tcvectordb.Client
	//cli                    *tcvectordb.RpcClient
	ctx                    = context.Background()
	database               = "go-sdk-test-db"
	collectionName         = "go-sdk-test-coll"
	collectionAlias        = "go-sdk-test-alias"
	embeddingCollection    = "go-sdk-test-emcoll"
	embedCollWithSparseVec = "go-sdk-test-emcoll-sparse-vec"
)

func init() {
	// 初始化客户端
	var err error
	cli, err = tcvectordb.NewClient("", "root",
		"", &tcvectordb.ClientOption{Timeout: 10 * time.Second,
			ReadConsistency: tcvectordb.StrongConsistency})

	if err != nil {
		log.Println("please input vdb address and authKey, then you can run testcases in test dir")
		panic(err)
	}
	cli.Debug(true)
}

func printErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ToJson(any interface{}) string {
	bytes, err := json.Marshal(any)
	if err != nil {
		return ""
	}
	return string(bytes)
}
