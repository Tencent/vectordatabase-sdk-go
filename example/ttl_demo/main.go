package main

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
)

type Demo struct {
	client *tcvectordb.Client
}

var (
	vectors = generateRandomVecs(768, 5)
)

func NewDemo(url, username, key string) (*Demo, error) {
	// cli, err := tcvectordb.NewRpcClient(url, username, key, &tcvectordb.ClientOption{
	// 	ReadConsistency: tcvectordb.EventualConsistency})
	cli, err := tcvectordb.NewClient(url, username, key, &tcvectordb.ClientOption{
		ReadConsistency: tcvectordb.EventualConsistency})
	if err != nil {
		return nil, err
	}
	// disable/enable http request log print
	// cli.Debug(false)
	return &Demo{client: cli}, nil
}

func (d *Demo) DropDatabase(ctx context.Context, database, collection string) error {
	log.Println("--------------------------- DropDatabase ---------------------------")
	// 删除db，db下的所有collection都将被删除
	dbDropResult, err := d.client.DropDatabase(ctx, database)
	if err != nil {
		return err
	}
	log.Printf("drop database result: %+v", dbDropResult)
	return nil
}

func (d *Demo) CreateDBAndCollection(ctx context.Context, database, collection string) error {
	// 创建DB
	log.Println("-------------------------- CreateDatabaseIfNotExists --------------------------")
	db, err := d.client.CreateDatabaseIfNotExists(ctx, database)
	if err != nil {
		return err
	}

	log.Println("------------------------- CreateCollection -------------------------")
	// 创建 Collection
	index := tcvectordb.Indexes{}
	index.VectorIndex = append(index.VectorIndex, tcvectordb.VectorIndex{
		FilterIndex: tcvectordb.FilterIndex{
			FieldName: "vector",
			FieldType: tcvectordb.Vector,
			IndexType: tcvectordb.HNSW,
		},
		Dimension:  768,
		MetricType: tcvectordb.IP,
		Params: &tcvectordb.HNSWParam{
			M:              16,
			EfConstruction: 200,
		},
	})

	index.FilterIndex = append(index.FilterIndex,
		tcvectordb.FilterIndex{FieldName: "id", FieldType: tcvectordb.String, IndexType: tcvectordb.PRIMARY})
	index.FilterIndex = append(index.FilterIndex,
		tcvectordb.FilterIndex{FieldName: "expire_at", FieldType: tcvectordb.Uint64, IndexType: tcvectordb.FILTER})

	// 配置TTL，Enable为是否开启TTL，TimeField为配置时间字段unix时间戳：doc将在expire_at设置的时间戳到达后失效
	// doc失效后，会在下个检测周期被清除，检测周期为1小时
	param := &tcvectordb.CreateCollectionParams{
		TtlConfig: &tcvectordb.TtlConfig{
			Enable:    true,
			TimeField: "expire_at",
		},
	}

	db.WithTimeout(time.Second * 30)
	_, err = db.CreateCollectionIfNotExists(ctx, collection, 3, 1, "test collection", index, param)
	if err != nil {
		return err
	}
	return nil
}

func (d *Demo) UpsertData(ctx context.Context, database, collection string) error {
	// 获取 Collection 对象
	coll := d.client.Database(database).Collection(collection)

	log.Println("------------------------------ Upsert ------------------------------")

	documentList := make([]tcvectordb.Document, 0)
	for i := 0; i < 5; i++ {
		id := "000" + strconv.Itoa(i)
		documentList = append(documentList, tcvectordb.Document{
			Id:     id,
			Vector: vectors[i],
			Fields: map[string]tcvectordb.Field{
				"expire_at": {Val: time.Now().Add(10 * time.Second).Unix()},
			},
		})
	}

	result, err := coll.Upsert(ctx, documentList)
	if err != nil {
		return err
	}
	log.Printf("UpsertResult: %+v, time: %v", result, time.Now())
	return nil
}

func (d *Demo) QueryData(ctx context.Context, database, collection string) error {
	// 获取 Collection 对象
	coll := d.client.Database(database).Collection(collection)

	log.Println("------------------------------ Query ------------------------------")

	documentIds := []string{"0000", "0001", "0002", "0003", "0004"}
	outputField := []string{"id", "expire_at"}
	params := tcvectordb.QueryDocumentParams{
		RetrieveVector: false,
		OutputFields:   outputField,
		Limit:          5,
		Offset:         0,
	}

	result, err := coll.Query(ctx, documentIds, &params)
	if err != nil {
		return err
	}
	log.Printf("QueryResult: total: %v, affect: %v", result.Total, result.AffectedCount)
	for _, doc := range result.Documents {
		log.Printf("QueryDocument: %+v", doc)
	}

	time.Sleep(time.Hour)
	log.Println("------------------------------ Query after 1 hour, and expired data will be deleted------------------------------")
	result, err = coll.Query(ctx, documentIds, &params)
	if err != nil {
		return err
	}
	log.Printf("Query after 1 hour: total: %v, affect: %v", result.Total, result.AffectedCount)
	return nil
}

func printErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func generateRandomVecs(dim int, vecNum int) [][]float32 {
	var randGen = rand.New(rand.NewSource(time.Now().UnixNano()))
	arr := make([][]float32, vecNum)
	for i := range arr {
		vector := make([]float32, dim)
		for j := 0; j < dim; j++ {
			vector[j] = randGen.Float32()
		}
		arr[i] = vector
	}
	return arr
}

func main() {
	database := "go-sdk-demo-db"
	collectionName := "go-sdk-demo-col-ttl"

	ctx := context.Background()
	testVdb, err := NewDemo("vdb http url or ip and port", "root", "key get from web console")
	printErr(err)
	err = testVdb.CreateDBAndCollection(ctx, database, collectionName)
	printErr(err)
	err = testVdb.UpsertData(ctx, database, collectionName)
	printErr(err)
	err = testVdb.QueryData(ctx, database, collectionName)
	printErr(err)
	err = testVdb.DropDatabase(ctx, database, collectionName)
	printErr(err)
}
