package main

import (
	"context"
	"log"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
)

type Demo struct {
	client *tcvectordb.RpcClient
}

func NewDemo(url, username, key string) (*Demo, error) {
	cli, err := tcvectordb.NewRpcClient(url, username, key, &tcvectordb.ClientOption{
		ReadConsistency: tcvectordb.StrongConsistency})
	if err != nil {
		return nil, err
	}
	// disable/enable http request log print
	// cli.Debug(false)
	return &Demo{client: cli}, nil
}

func (d *Demo) DeleteAndDrop(ctx context.Context, database, collection string) error {
	log.Println("--------------------------- DropDatabase ---------------------------")
	dbDropResult, err := d.client.DropDatabase(ctx, database)
	if err != nil {
		return err
	}
	log.Printf("drop database result: %+v", dbDropResult)
	return nil
}

func (d *Demo) CreateDBAndCollection(ctx context.Context, database, collection, alias string) error {
	log.Println("-------------------------- CreateDatabaseIfNotExists --------------------------")
	db, err := d.client.CreateDatabaseIfNotExists(ctx, database)
	if err != nil {
		return err
	}

	log.Println("------------------------- CreateCollectionIfNotExists -------------------------")
	index := tcvectordb.Indexes{}
	index.VectorIndex = append(index.VectorIndex, tcvectordb.VectorIndex{
		FilterIndex: tcvectordb.FilterIndex{
			FieldName: "vector",
			FieldType: tcvectordb.Vector,
			IndexType: tcvectordb.HNSW,
		},
		Dimension:  1,
		MetricType: tcvectordb.IP,
		Params: &tcvectordb.HNSWParam{
			M:              16,
			EfConstruction: 200,
		},
	})
	index.FilterIndex = append(index.FilterIndex, tcvectordb.FilterIndex{FieldName: "id", FieldType: tcvectordb.String, IndexType: tcvectordb.PRIMARY, AutoId: "uuid"})
	index.FilterIndex = append(index.FilterIndex, tcvectordb.FilterIndex{FieldName: "field_json", FieldType: tcvectordb.Json, IndexType: tcvectordb.FILTER})

	db.WithTimeout(time.Second * 30)
	_, err = db.CreateCollectionIfNotExists(ctx, collection, 3, 1, "test collection", index)
	if err != nil {
		return err
	}

	log.Println("------------------------ DescribeCollection ------------------------")
	// 查看 Collection 信息
	colRes, err := db.DescribeCollection(ctx, collection)
	if err != nil {
		return err
	}
	log.Printf("DescribeCollection: %+v", colRes)
	return nil
}

func (d *Demo) UpsertData(ctx context.Context, database, collection string) error {
	log.Println("------------------------------ Upsert ------------------------------")
	documentList := []tcvectordb.Document{
		{

			Vector: []float32{0},
			Fields: map[string]tcvectordb.Field{
				"field_json": {Val: map[string]interface{}{
					"username": "Alice",
					"age":      28,
				}},
			},
		},
		{
			Vector: []float32{0.1},
			Fields: map[string]tcvectordb.Field{
				"field_json": {Val: map[string]interface{}{
					"username": "Bob",
					"age":      25,
				}},
			},
		},
		{
			Vector: []float32{0.2},
			Fields: map[string]tcvectordb.Field{
				"field_json": {Val: map[string]interface{}{
					"username": "Charlie",
					"age":      35,
				}},
			},
		},
		{
			Vector: []float32{0.3},
			Fields: map[string]tcvectordb.Field{
				"field_json": {Val: map[string]interface{}{
					"username": "David",
					"age":      28,
				}},
			},
		},
	}
	result, err := d.client.Upsert(ctx, database, collection, documentList)
	if err != nil {
		return err
	}
	log.Printf("UpsertResult: %+v", result)
	return nil
}

func (d *Demo) QueryData(ctx context.Context, database, collection string) error {
	log.Println("------------------------------ Query ------------------------------")
	filter := tcvectordb.NewFilter(`field_json.username="David"`)

	result, err := d.client.Query(ctx, database, collection, []string{}, &tcvectordb.QueryDocumentParams{
		Filter:         filter,
		RetrieveVector: true,
		Limit:          10,
	})
	if err != nil {
		return err
	}
	log.Printf("QueryResult: total: %v, affect: %v", result.Total, result.AffectedCount)
	for _, doc := range result.Documents {
		log.Printf("QueryDocument: %+v", doc)
	}
	return nil
}

func printErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	database := "go-sdk-demo-db"
	collectionName := "go-sdk-demo-col"
	collectionAlias := "go-sdk-demo-alias"

	ctx := context.Background()
	testVdb, err := NewDemo("vdb http url or ip and port", "vdb username", "key get from web console")
	printErr(err)
	defer testVdb.client.Close()

	err = testVdb.CreateDBAndCollection(ctx, database, collectionName, collectionAlias)
	printErr(err)
	err = testVdb.UpsertData(ctx, database, collectionName)
	printErr(err)
	err = testVdb.QueryData(ctx, database, collectionName)
	printErr(err)
	err = testVdb.DeleteAndDrop(ctx, database, collectionName)
	printErr(err)
}
