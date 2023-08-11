package main

import (
	"context"
	"time"
	"vectordb-sdk-go/internal/client"
	"vectordb-sdk-go/tcvectordb"
)

func main() {
	cli, err := tcvectordb.NewClient("http://127.0.0.1", "root", "key", &client.ClientOption{
		MaxIdldConnPerHost: 50,
		IdleConnTimeout:    time.Second * 30,
	})
	if err != nil {
		panic(err)
	}
	cli.CreateDatabase(context.TODO(), "dbtest")
}
