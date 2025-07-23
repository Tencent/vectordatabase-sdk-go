package test

import (
	"testing"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/ai_service"
)

func Test_embedding(t *testing.T) {
	retrieveDenseVector := true
	retrieveSparseVector := false

	data := []string{}
	for i := 0; i < 100; i++ {
		data = append(data, "hello world")
	}
	result, err := cli.Embedding(ctx, tcvectordb.EmbeddingParams{
		Model:    "bge-base-zh",
		DataType: "text",
		ModelParams: &ai_service.ModelParams{
			RetrieveDenseVector:  &retrieveDenseVector,
			RetrieveSparseVector: &retrieveSparseVector,
		},
		Data: data,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result.TokenUsed)
}
