package main

import (
	"context"
	"log"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/index"
)

type Demo struct {
	client *tcvectordb.RpcClient
}

func NewDemo(url, username, key string) (*Demo, error) {
	cli, err := tcvectordb.NewRpcClient(url, username, key, &tcvectordb.ClientOption{
		ReadConsistency: tcvectordb.EventualConsistency})
	if err != nil {
		return nil, err
	}

	// disable/enable rpc request log print
	// cli.Debug(false)
	return &Demo{client: cli}, nil
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

func (d *Demo) CreateDB(ctx context.Context, database string) error {
	log.Println("-------------------------- CreateDatabaseIfNotExists --------------------------")
	_, err := d.client.CreateDatabaseIfNotExists(ctx, database)
	if err != nil {
		return err
	}
	return nil
}

func (d *Demo) CreateCollection(ctx context.Context, database, collection string, fieldType tcvectordb.FieldType) error {
	db := d.client.Database(database)
	log.Println("------------------------- CreateCollection -------------------------")
	index := tcvectordb.Indexes{}
	index.VectorIndex = append(index.VectorIndex, tcvectordb.VectorIndex{
		FilterIndex: tcvectordb.FilterIndex{
			FieldName: "vector",
			FieldType: fieldType,
			IndexType: tcvectordb.HNSW,
		},
		Dimension:  3,
		MetricType: tcvectordb.COSINE,
		Params: &tcvectordb.HNSWParam{
			M:              16,
			EfConstruction: 200,
		},
	})
	index.FilterIndex = append(index.FilterIndex, tcvectordb.FilterIndex{FieldName: "id", FieldType: tcvectordb.String, IndexType: tcvectordb.PRIMARY})

	db.WithTimeout(time.Second * 30)
	_, err := db.CreateCollection(ctx, collection, 3, 1, "test collection", index)
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

func (d *Demo) ModifyVectorIndex(ctx context.Context, database, collection string) error {
	log.Println("------------------------------ ModifyVectorIndex ------------------------------")

	setThrottle := int32(1)
	param := tcvectordb.ModifyVectorIndexParam{
		RebuildRules: &index.RebuildRules{
			Throttle: &setThrottle,
		},
	}
	param.VectorIndexes = make([]tcvectordb.ModifyVectorIndex, 0)
	param.VectorIndexes = append(param.VectorIndexes, tcvectordb.ModifyVectorIndex{
		FieldName:  "vector",
		FieldType:  "float16_vector",
		MetricType: tcvectordb.COSINE,
		Params: &tcvectordb.HNSWParam{
			M:              32,
			EfConstruction: 500,
		}})

	return d.client.ModifyVectorIndex(ctx, database, collection, param)
}

func printErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	database := "go-demo-db-quant"
	collNameWithVec := "go-demo-col-quant"
	collNameWithFloat16Vec := "go-demo-col-quant-float16vec"
	collNameWithBFloat16Vec := "go-demo-col-quant-bfloat16vec"

	ctx := context.Background()
	testVdb, err := NewDemo("vdb http url or ip and port", "vdb username", "key get from web console")
	printErr(err)
	defer testVdb.client.Close()

	err = testVdb.CreateDB(ctx, database)
	printErr(err)
	err = testVdb.CreateCollection(ctx, database, collNameWithVec, tcvectordb.Vector)
	printErr(err)
	err = testVdb.CreateCollection(ctx, database, collNameWithFloat16Vec, tcvectordb.Float16Vector)
	printErr(err)
	err = testVdb.CreateCollection(ctx, database, collNameWithBFloat16Vec, tcvectordb.BFloat16Vector)
	printErr(err)
	err = testVdb.ModifyVectorIndex(ctx, database, collNameWithVec)
	printErr(err)
	err = testVdb.DropDB(ctx, database)
	printErr(err)
}
