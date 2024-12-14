package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/utils"
)

type Demo struct {
	//client *tcvectordb.Client
	client *tcvectordb.RpcClient
}

func NewDemo(url, username, key string) (*Demo, error) {
	// cli, err := tcvectordb.NewClient(url, username, key, &tcvectordb.ClientOption{
	// 	ReadConsistency: tcvectordb.EventualConsistency})
	cli, err := tcvectordb.NewRpcClient(url, username, key, &tcvectordb.ClientOption{
		ReadConsistency: tcvectordb.StrongConsistency})
	if err != nil {
		return nil, err
	}
	// disable/enable http request log print
	// cli.Debug(false)
	return &Demo{client: cli}, nil
}

func (d *Demo) Clear(ctx context.Context, database string) error {
	log.Println("--------------------------- DropDatabase ---------------------------")
	result, err := d.client.DropDatabase(ctx, database)
	if err != nil {
		return err
	}
	log.Printf("drop database result: %+v", result)
	return nil
}

func (d *Demo) DeleteAndDrop(ctx context.Context, database, collection string) error {
	// 删除collection，删除collection的同时，其中的数据也将被全部删除
	log.Println("-------------------------- DropCollection --------------------------")
	colDropResult, err := d.client.Database(database).DropCollection(ctx, collection)
	if err != nil {
		return err
	}
	log.Printf("drop collection result: %+v", colDropResult)

	log.Println("--------------------------- DropDatabase ---------------------------")
	// 删除db，db下的所有collection都将被删除
	dbDropResult, err := d.client.DropDatabase(ctx, database)
	if err != nil {
		return err
	}
	log.Printf("drop database result: %+v", dbDropResult)
	return nil
}

func (d *Demo) CreateDBAndCollection(ctx context.Context, database, collection, alias string) error {
	// 创建DB--'book'
	log.Println("-------------------------- CreateDatabase --------------------------")
	db, err := d.client.CreateDatabase(ctx, database)
	if err != nil {
		return err
	}

	log.Println("--------------------------- ListDatabase ---------------------------")
	dbList, err := d.client.ListDatabase(ctx)
	if err != nil {
		return err
	}
	for _, db := range dbList.Databases {
		log.Printf("database: %s", db.DatabaseName)
	}

	log.Println("------------------------- CreateCollection -------------------------")
	index := tcvectordb.Indexes{}
	index.VectorIndex = append(index.VectorIndex, tcvectordb.VectorIndex{
		FilterIndex: tcvectordb.FilterIndex{
			FieldName: "vector",
			FieldType: tcvectordb.BinaryVector,
			IndexType: tcvectordb.BIN_FLAT,
		},
		Dimension:  16,
		MetricType: tcvectordb.HAMMING,
	})
	index.FilterIndex = append(index.FilterIndex, tcvectordb.FilterIndex{FieldName: "id", FieldType: tcvectordb.String, IndexType: tcvectordb.PRIMARY})

	maxStrLen := uint32(32)
	param := &tcvectordb.CreateCollectionParams{
		FilterIndexConfig: &tcvectordb.FilterIndexConfig{
			FilterAll:          true,
			FieldsWithoutIndex: []string{"bookName"},
			MaxStrLen:          &maxStrLen,
		},
	}

	db.WithTimeout(time.Second * 30)
	_, err = db.CreateCollection(ctx, collection, 3, 1, "test collection", index, param)
	if err != nil {
		return err
	}

	log.Println("-------------------------- ListCollection --------------------------")
	// 列出所有 Collection
	collListRes, err := db.ListCollection(ctx)
	if err != nil {
		return err
	}
	for _, col := range collListRes.Collections {
		log.Printf("ListCollection: %+v", col)
	}

	log.Println("----------------------------- SetAlias -----------------------------")
	// 设置 Collection 的 alias
	_, err = db.SetAlias(ctx, collection, alias)
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

	log.Println("---------------------------- DeleteAlias ---------------------------")
	// 删除 Collection 的 alias
	delAliasRes, err := db.DeleteAlias(ctx, alias)
	if err != nil {
		return err
	}
	log.Printf("DeleteAliasResult: %v", delAliasRes)
	return nil
}

func (d *Demo) UpsertData(ctx context.Context, database, collection string) error {
	// 获取 Collection 对象
	coll := d.client.Database(database).Collection(collection)

	log.Println("------------------------------ Upsert ------------------------------")

	binaryVectors := [][]byte{
		{1, 1, 1, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 1, 0, 0, 1, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0, 0, 1, 1, 0, 1, 0, 0, 1, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 1, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 1, 0, 0, 1},
	}

	insertVectors := make([][]float32, 0)
	for i, binaryVector := range binaryVectors {
		res, err := utils.BinaryToUint8(binaryVector)
		if err != nil {
			return errors.New(fmt.Sprintf("convert i: %v, err: %v", i, err))
		}
		insertVectors = append(insertVectors, res)
	}
	documentList := []tcvectordb.Document{
		{
			Id:     "0001",
			Vector: insertVectors[0],
			Fields: map[string]tcvectordb.Field{
				"bookName": {Val: "西游记"},
				"author":   {Val: "吴承恩"},
				"page":     {Val: 21},
				"segment":  {Val: "富贵功名，前缘分定，为人切莫欺心。"},
			},
		},
		{
			Id:     "0002",
			Vector: insertVectors[1],
			Fields: map[string]tcvectordb.Field{
				"bookName": {Val: "西游记"},
				"author":   {Val: "吴承恩"},
				"page":     {Val: 22},
				"segment":  {Val: "正大光明，忠良善果弥深。些些狂妄天加谴，眼前不遇待时临。"},
			},
		},
		{
			Id:     "0003",
			Vector: insertVectors[2],
			Fields: map[string]tcvectordb.Field{
				"bookName": {Val: "三国演义"},
				"author":   {Val: "罗贯中"},
				"page":     {Val: 23},
				"segment":  {Val: "细作探知这个消息，飞报吕布。"},
			},
		},
		{
			Id:     "0004",
			Vector: insertVectors[3],
			Fields: map[string]tcvectordb.Field{
				"bookName": {Val: "三国演义"},
				"author":   {Val: "罗贯中"},
				"page":     {Val: 24},
				"segment":  {Val: "布大惊，与陈宫商议。宫曰：“闻刘玄德新领徐州，可往投之。”布从其言，竟投徐州来。有人报知玄德。"},
			},
		},
		{
			Id:     "0005",
			Vector: insertVectors[4],
			Fields: map[string]tcvectordb.Field{
				"bookName": {Val: "三国演义"},
				"author":   {Val: "罗贯中"},
				"page":     {Val: 25},
				"segment":  {Val: "玄德曰：“布乃当今英勇之士，可出迎之。”糜竺曰：“吕布乃虎狼之徒，不可收留；收则伤人矣。"},
			},
		},
	}
	result, err := coll.Upsert(ctx, documentList)
	if err != nil {
		return err
	}
	log.Printf("UpsertResult: %+v", result)
	return nil
}

func (d *Demo) QueryAndSearchData(ctx context.Context, database, collection string) error {
	// 获取 Collection 对象
	coll := d.client.Database(database).Collection(collection)

	log.Println("------------------------------ Query ------------------------------")
	documentIds := []string{"0001", "0002", "0003", "0004", "0005"}
	filter := tcvectordb.NewFilter(`bookName="三国演义"`)
	outputField := []string{"id", "bookName"}

	result, err := coll.Query(ctx, documentIds, &tcvectordb.QueryDocumentParams{
		Filter:         filter,
		RetrieveVector: false,
		OutputFields:   outputField,
		Limit:          2,
		Offset:         1,
	})
	if err != nil {
		return err
	}
	log.Printf("QueryResult: total: %v, affect: %v", result.Total, result.AffectedCount)
	for _, doc := range result.Documents {
		log.Printf("QueryDocument: %+v", doc)
	}

	log.Println("------------------------------ Query ------------------------------")

	filter = tcvectordb.NewFilter(`author="吴承恩"`)
	_, err = coll.Query(ctx, documentIds, &tcvectordb.QueryDocumentParams{
		Filter:         filter,
		RetrieveVector: false,
		OutputFields:   outputField,
		Limit:          2,
		Offset:         1,
	})
	if err == nil {
		log.Printf("query succ. Because the author filter" +
			" will be ignored, when filterAll is true ")
	}

	log.Println("------------------------------ Search ------------------------------")

	searchVector, err := utils.BinaryToUint8([]byte{0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 1, 1, 0, 1})
	if err != nil {
		return fmt.Errorf("convert failed, err: %v", err.Error())
	}

	searchResult, err := coll.Search(ctx,
		[][]float32{searchVector}, //指定检索向量，最多指定20个
		&tcvectordb.SearchDocumentParams{
			RetrieveVector: true,
			Limit:          2,
		})
	if err != nil {
		return err
	}
	// 输出相似性检索结果，检索结果为二维数组，每一位为一组返回结果，分别对应search时指定的多个向量
	for i, item := range searchResult.Documents {
		log.Printf("SearchDocumentResult, index: %d ==================", i)
		for _, doc := range item {
			log.Printf("SearchDocument: %+v", doc)
		}
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
	err = testVdb.Clear(ctx, database)
	printErr(err)
	err = testVdb.CreateDBAndCollection(ctx, database, collectionName, collectionAlias)
	printErr(err)
	err = testVdb.UpsertData(ctx, database, collectionName)
	printErr(err)
	err = testVdb.QueryAndSearchData(ctx, database, collectionName)
	printErr(err)
	err = testVdb.DeleteAndDrop(ctx, database, collectionName)
	printErr(err)
}
