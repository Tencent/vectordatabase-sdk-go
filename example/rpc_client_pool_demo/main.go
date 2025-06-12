package main

import (
	"context"
	"log"
	"sync"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
)

func main() {
	database := "go-sdk-demo-db"
	collectionName := "go-sdk-demo-col"
	url := "vdb http url or ip and port"
	key := "key get from web console"
	username := "vdb username"

	ctx := context.Background()
	cliPool, err := tcvectordb.NewRpcClientPool(url, username, key, &tcvectordb.ClientOption{
		ReadConsistency: tcvectordb.EventualConsistency,
		RpcPoolSize:     9})
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer cliPool.Close()

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
	index.FilterIndex = append(index.FilterIndex, tcvectordb.FilterIndex{FieldName: "id", FieldType: tcvectordb.String, IndexType: tcvectordb.PRIMARY, AutoId: "uuid"})
	_, err = cliPool.CreateCollectionIfNotExists(ctx, database, collectionName, 3, 2, "test collection", index)
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

	log.Println("-------------------------- Concurrent Upsert --------------------------")

	// 并发参数
	const concurrency = 100
	const totalCalls = 1000000
	const callsPerGoroutine = totalCalls / concurrency

	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < callsPerGoroutine; j++ {
				documentList := []tcvectordb.Document{
					{
						Vector: []float32{0.2123, 0.21, 0.213},
						Fields: map[string]tcvectordb.Field{
							"bookName": {Val: "西游记"},
							"author":   {Val: "吴承恩"},
							"page":     {Val: 21},
							"segment":  {Val: "富贵功名，前缘分定，为人切莫欺心。"},
						},
					},
				}

				res, err := cliPool.Upsert(ctx, database, collectionName, documentList)
				if err != nil {
					log.Printf("Goroutine %d, call %d failed: %v", goroutineID, j+1, err)
				} else {
					log.Printf("Goroutine %d, call %d success: %+v", goroutineID, j+1, res)
				}
			}
		}(i)
	}

	// 等待所有goroutine完成
	wg.Wait()
	log.Println("All upsert operations completed")

	cliPool.DropDatabase(ctx, database)

}
