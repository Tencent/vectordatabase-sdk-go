package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/tencent/vectordatabase-sdk-go/tcvdbtext/encoder"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
)

type Demo struct {
	client *tcvectordb.RpcClient
	bm25   encoder.SparseEncoder
}

func NewDemo(url, username, key string) (*Demo, error) {
	client, err := tcvectordb.NewRpcClient(url, username, key, nil)
	if err != nil {
		return nil, err
	}
	bm25, err := encoder.NewBM25Encoder(&encoder.BM25EncoderParams{Bm25Language: "zh"})
	if err != nil {
		return nil, err
	}
	return &Demo{client: client, bm25: bm25}, nil
}

func (d *Demo) DropDB(ctx context.Context, database string) error {
	log.Println("--------------------------- DropDatabase ---------------------------")
	result, err := d.client.DropDatabase(ctx, database)
	if err != nil {
		return err
	}
	log.Printf("drop database result: %+v", result)
	return nil
}

func (d *Demo) CreateDBAndCollection(ctx context.Context, database string, collection string) error {
	log.Println("--------------------------- CreateDatabase ---------------------------")
	db, err := d.client.CreateDatabaseIfNotExists(ctx, database)
	if err != nil {
		return err
	}
	log.Println("--------------------------- CreateCollection ---------------------------")
	diskSwapEnabled := false
	indexes := tcvectordb.Indexes{
		SparseVectorIndex: []tcvectordb.SparseVectorIndex{
			{
				FieldName:       "sparse_vector",
				FieldType:       tcvectordb.SparseVector,
				IndexType:       tcvectordb.SPARSE_INVERTED,
				MetricType:      tcvectordb.IP,
				DiskSwapEnabled: &diskSwapEnabled,
			},
		},
		FilterIndex: []tcvectordb.FilterIndex{
			{
				FieldName: "id",
				FieldType: tcvectordb.String,
				IndexType: tcvectordb.PRIMARY,
			},
		},
	}
	_, err = db.CreateCollectionIfNotExists(ctx, collection, 3, 2, "test collection", indexes)
	if err != nil {
		return err
	}
	return nil
}

func (d *Demo) DescribeCollection(ctx context.Context, database, collection string) error {
	log.Println("------------------------ DescribeCollection ------------------------")
	colRes, err := d.client.Database(database).DescribeCollection(ctx, collection)
	if err != nil {
		return err
	}
	log.Printf("DescribeCollection: %+v", ToJson(colRes))
	return nil
}

func (d *Demo) UpsertData(ctx context.Context, database, collection string) error {
	log.Println("------------------------------ Upsert ------------------------------")
	documentList := []tcvectordb.Document{
		{
			Id: "0001",
			SparseVector: []encoder.SparseVecItem{
				{TermId: 1172076521, Score: 0.71296215},
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

func (d *Demo) ModifySparseVectorIndex(ctx context.Context, database, collection string) error {
	log.Println("------------------------------ ModifySparseVectorIndex ------------------------------")
	diskSwapEnabled := true
	param := tcvectordb.ModifyVectorIndexParam{
		VectorIndexes: []tcvectordb.ModifyVectorIndex{
			{FieldName: "sparse_vector", FieldType: "sparseVector", DiskSwapEnabled: &diskSwapEnabled},
		},
	}
	return d.client.ModifyVectorIndex(ctx, database, collection, param)
}

func main() {
	database := "go-sdk-demo-db"
	collection := "go-sdk-demo-col-sparsevec"

	ctx := context.Background()
	testVdb, err := NewDemo("vdb http url or ip and port", "vdb username", "key get from web console")
	printErr(err)
	defer testVdb.client.Close()

	err = testVdb.DropDB(ctx, database)
	printErr(err)
	err = testVdb.CreateDBAndCollection(ctx, database, collection)
	printErr(err)

	err = testVdb.UpsertData(ctx, database, collection)
	printErr(err)
	err = testVdb.DescribeCollection(ctx, database, collection)
	printErr(err)
	err = testVdb.ModifySparseVectorIndex(ctx, database, collection)
	printErr(err)
	err = testVdb.DescribeCollection(ctx, database, collection)
	printErr(err)
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
