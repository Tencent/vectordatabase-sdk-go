package main

import (
	"context"
	"log"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
)

func main() {
	database := "go-sdk-demo-db"
	collectionName := "go-sdk-demo-col"
	url := "vdb http url or ip and port"
	key := "key get from web console"

	ctx := context.Background()
	cliPool, err := tcvectordb.NewRpcClientPool(url, "root", key, &tcvectordb.ClientOption{
		ReadConsistency: tcvectordb.EventualConsistency,
		RpcPoolSize:     5})
	if err != nil {
		log.Fatalf(err.Error())
	}

	//cliPool.Debug(true)

	log.Println("-------------------------- CreateDatabaseIfNotExists --------------------------")
	_, err = cliPool.CreateDatabaseIfNotExists(ctx, database)
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Println("--------------------------- ListDatabase ---------------------------")
	dbList, err := cliPool.ListDatabase(ctx)
	if err != nil {
		log.Fatalf(err.Error())
	}
	for _, db := range dbList.Databases {
		log.Printf("database: %s", db.DatabaseName)
	}

	log.Println("-------------------------- CreateCollectionIfNotExists --------------------------")
	index := tcvectordb.Indexes{}
	index.VectorIndex = append(index.VectorIndex, tcvectordb.VectorIndex{
		FilterIndex: tcvectordb.FilterIndex{
			FieldName: "vector",
			FieldType: tcvectordb.Vector,
			IndexType: tcvectordb.HNSW,
		},
		Dimension:  3,
		MetricType: tcvectordb.IP,
		Params: &tcvectordb.HNSWParam{
			M:              16,
			EfConstruction: 200,
		},
	})
	index.FilterIndex = append(index.FilterIndex, tcvectordb.FilterIndex{FieldName: "id", FieldType: tcvectordb.String, IndexType: tcvectordb.PRIMARY})
	_, err = cliPool.CreateCollectionIfNotExists(ctx, database, collectionName, 3, 1, "test collection", index)
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Println("------------------------ DescribeCollection ------------------------")
	// 查看 Collection 信息
	colRes, err := cliPool.DescribeCollection(ctx, database, collectionName)
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("DescribeCollection: %+v", colRes)

	log.Println("-------------------------- Upsert --------------------------")

	documentList := []tcvectordb.Document{
		{
			Id:     "0001",
			Vector: []float32{0.2123, 0.21, 0.213},
			Fields: map[string]tcvectordb.Field{
				"bookName": {Val: "西游记"},
				"author":   {Val: "吴承恩"},
				"page":     {Val: 21},
				"segment":  {Val: "富贵功名，前缘分定，为人切莫欺心。"},
			},
		},
	}
	for i := 0; i < 5; i++ {
		res, err := cliPool.Upsert(ctx, database, collectionName, documentList)
		if err != nil {
			log.Fatalf(err.Error())
		}
		log.Printf("UpsertResult: %+v", res)
	}

	cliPool.DropDatabase(ctx, database)
	cliPool.Close()
}
