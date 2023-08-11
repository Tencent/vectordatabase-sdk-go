package tcvectordb

import (
	"vectordb-sdk-go/internal/client"
	"vectordb-sdk-go/internal/engine"
)

func NewClient(url, username, key string, option *client.ClientOption) (engine.VectorDBClient, error) {
	sdkCli, err := client.NewClient(url, username, key, option)
	if err != nil {
		return nil, err
	}
	return engine.VectorDB(sdkCli), nil
}
