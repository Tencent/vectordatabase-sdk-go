package main

import (
	"context"
	"log"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/user"
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
	// disable/enable http request log print
	// cli.Debug(false)
	return &Demo{client: cli}, nil
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

	db.WithTimeout(time.Second * 30)
	_, err = db.CreateCollectionIfNotExists(ctx, collection, 3, 1, "test collection", index)
	if err != nil {
		return err
	}
	return nil
}

func (d *Demo) UserOperations(ctx context.Context, database, collection, username string) error {
	log.Println("--------------------------- DropUser ---------------------------")
	err := d.client.DropUser(ctx, tcvectordb.DropUserParams{User: username})
	if err != nil {
		log.Printf("drop user failed, err: %+v", err.Error())
	}

	log.Println("--------------------------- CreateUser ---------------------------")
	err = d.client.CreateUser(ctx, tcvectordb.CreateUserParams{User: username, Password: "0dd8e8b3d674"})
	if err != nil {
		return err
	}

	log.Println("--------------------------- DescribeUser ---------------------------")
	result, err := d.client.DescribeUser(ctx, tcvectordb.DescribeUserParams{User: username})
	if err != nil {
		return err
	}
	log.Printf("DescribeUser Result: %+v", result)

	log.Println("--------------------------- GrantToUser ---------------------------")
	err = d.client.GrantToUser(ctx, tcvectordb.GrantToUserParams{
		User: username,
		Privileges: []*user.Privilege{
			{
				Resource: database + ".*",
				Actions:  []string{"read"},
			},
			{
				Resource: database + ".*",
				Actions:  []string{"readWrite"},
			},
		}})
	if err != nil {
		return err
	}

	log.Println("--------------------------- DescribeUser ---------------------------")
	result, err = d.client.DescribeUser(ctx, tcvectordb.DescribeUserParams{User: username})
	if err != nil {
		return err
	}
	log.Printf("DescribeUser Result: %+v", result)

	log.Println("--------------------------- RevokeFromUser ---------------------------")
	err = d.client.RevokeFromUser(ctx, tcvectordb.RevokeFromUserParams{User: username,
		Privileges: []*user.Privilege{
			{
				Resource: database + ".*",
				Actions:  []string{"read"},
			},
		}})
	if err != nil {
		return err
	}

	log.Println("--------------------------- DescribeUser ---------------------------")
	result, err = d.client.DescribeUser(ctx, tcvectordb.DescribeUserParams{User: username})
	if err != nil {
		return err
	}
	log.Printf("DescribeUser Result: %+v", result)

	log.Println("--------------------------- ListUser ---------------------------")
	listRes, err := d.client.ListUser(ctx)
	if err != nil {
		return err
	}
	log.Printf("ListUser Result: %+v", listRes)

	log.Println("--------------------------- DropUser ---------------------------")
	err = d.client.DropUser(ctx, tcvectordb.DropUserParams{User: username})
	if err != nil {
		return err
	}

	return nil
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
func printErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	database := "go-sdk-demo-db"
	collectionName := "go-sdk-demo-col"
	userName := "zhangsan"

	ctx := context.Background()
	testVdb, err := NewDemo("vdb http url or ip and port", "root", "key get from web console")
	printErr(err)
	err = testVdb.CreateDBAndCollection(ctx, database, collectionName)
	printErr(err)
	err = testVdb.UserOperations(ctx, database, collectionName, userName)
	printErr(err)
	err = testVdb.DropDatabase(ctx, database, collectionName)
	printErr(err)
}
