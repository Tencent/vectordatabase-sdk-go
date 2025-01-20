package test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/user"
)

func Test_DropUser(t *testing.T) {
	err := cli.DropUser(ctx, tcvectordb.DropUserParams{
		User: "test_user_1",
	})
	printErr(err)

	err = cli.DropUser(ctx, tcvectordb.DropUserParams{
		User: "test_user_2",
	})
	printErr(err)

}

func Test_CreateUser(t *testing.T) {
	err := cli.CreateUser(ctx, tcvectordb.CreateUserParams{
		User:     "test_user_1",
		Password: "123hello!",
	})
	printErr(err)
	err = cli.CreateUser(ctx, tcvectordb.CreateUserParams{
		User:     "test_user_2",
		Password: "123hello!",
	})
	printErr(err)
}

func Test_CreateUser_InvalidPass(t *testing.T) {
	err := cli.CreateUser(ctx, tcvectordb.CreateUserParams{
		User:     "test_user_1",
		Password: "123!",
	})
	printErr(err)
}

func Test_RevokeUser(t *testing.T) {
	db, err := cli.CreateDatabaseIfNotExists(ctx, database)
	printErr(err)
	log.Printf("create database success, %s", db.DatabaseName)

	err = cli.RevokeFromUser(ctx, tcvectordb.RevokeFromUserParams{
		User: "test_user_1",

		Privileges: []*user.Privilege{
			{
				Resource: database + ".*",
				Actions:  []string{"read"},
			},
		},
	})
	printErr(err)
}

func Test_GrantUSer(t *testing.T) {
	db, err := cli.CreateDatabaseIfNotExists(ctx, database)
	printErr(err)
	log.Printf("create database success, %s", db.DatabaseName)

	err = cli.GrantToUser(ctx, tcvectordb.GrantToUserParams{
		User: "test_user_1",

		Privileges: []*user.Privilege{
			{
				Resource: database + ".*",
				Actions:  []string{"read"},
			},
		},
	})
	printErr(err)
}

func Test_DescribeUser(t *testing.T) {
	res, err := cli.DescribeUser(ctx, tcvectordb.DescribeUserParams{
		User: "test_user_1",
	})
	printErr(err)
	fmt.Println(ToJson(res))
}

func Test_ListUser(t *testing.T) {
	res, err := cli.ListUser(ctx)
	printErr(err)
	fmt.Println(ToJson(res))
}

func Test_ChangePassword(t *testing.T) {
	err := cli.ChangePassword(ctx, tcvectordb.ChangePasswordParams{
		User:     "test_user_1",
		Password: "123..A.!",
	})
	printErr(err)
}

func Test_User(t *testing.T) {
	_, err := cli.CreateDatabaseIfNotExists(ctx, database+"1")
	printErr(err)

	res, err := cli.ListDatabase(ctx)
	printErr(err)
	println(ToJson(res))

	newCli, err := tcvectordb.NewRpcClient("http://xx",
		"test_user_1",
		"123..A.!", &tcvectordb.ClientOption{Timeout: 10 * time.Second,
			ReadConsistency: tcvectordb.StrongConsistency})
	printErr(err)
	newCli.Debug(true)

	collRes, err := newCli.Database(database).ListCollection(ctx)
	printErr(err)
	println(ToJson(collRes))

	collRes, err = newCli.Database(database + "1").ListCollection(ctx)
	printErr(err)
	println(ToJson(collRes))

}
