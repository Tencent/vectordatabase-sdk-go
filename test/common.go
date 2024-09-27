package test

import (
	"context"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
)

var (
	cli                    *tcvectordb.Client
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
	cli, err = tcvectordb.NewClient("", "root", "", &tcvectordb.ClientOption{Timeout: 10 * time.Second})
	cli.Debug(true)
	if err != nil {
		panic(err)
	}
}
