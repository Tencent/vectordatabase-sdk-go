package main

import (
	"context"
	"log"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/ai_service"
)

type Demo struct {
	client *tcvectordb.Client
}

func NewDemo(url, username, key string) (*Demo, error) {
	cli, err := tcvectordb.NewClient(url, username, key, &tcvectordb.ClientOption{
		ReadConsistency: tcvectordb.EventualConsistency})
	if err != nil {
		return nil, err
	}
	// disable/enable http request log print
	// cli.Debug(false)
	return &Demo{client: cli}, nil
}

func printErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (d *Demo) EmbeddingService(ctx context.Context, model string, data []string) error {
	retrieveDenseVector := true
	retrieveSparseVector := true
	embeddingParams := tcvectordb.EmbeddingParams{
		Model:    model,
		DataType: "text",
		ModelParams: &ai_service.ModelParams{
			RetrieveDenseVector:  &retrieveDenseVector,
			RetrieveSparseVector: &retrieveSparseVector,
		},
		Data: data,
	}

	embeddingResult, err := d.client.Embedding(ctx, embeddingParams)
	if err != nil {
		return err
	}
	log.Println(embeddingResult)
	return nil
}

func main() {
	ctx := context.Background()

	testVdb, err := NewDemo("vdb http url or ip and port", "vdb username", "key get from web console")
	printErr(err)
	defer testVdb.client.Close()

	err = testVdb.EmbeddingService(ctx, "bge-m3", []string{"什么是腾讯云向量数据库"})
	printErr(err)
}
