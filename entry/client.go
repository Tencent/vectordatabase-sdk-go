package entry

import (
	"context"
	"time"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/model"
)

// SdkClient the http client interface
type SdkClient interface {
	Request(ctx context.Context, req, res interface{}) error
	Options() model.ClientOption
	WithTimeout(d time.Duration)
	Debug(v bool)
	Close()
}

type VectorDBClient interface {
	DatabaseInterface
}
