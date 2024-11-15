# TencentCloud VectorDB Go SDK

## Getting started

### Prerequisites
1. Go 1.17 or higher

### Install TencentCloud VectorDB Go SDK

1. Use `go get` to install the latest version of the TencentCloud VectorDB Go SDK and dependencies: 
```sh
go get -u github.com/tencent/vectordatabase-sdk-go/tcvectordb
```

2. Create New VectorDB Client To Start:
```go
import "github.com/tencent/vectordatabase-sdk-go/tcvectordb"

cli, err := tcvectordb.NewRpcClient("vdb http url or ip and post", "root", "key get from web console", &tcvectordb.ClientOption{
		ReadConsistency: tcvectordb.EventualConsistency,
		MaxIdldConnPerHost: 10,
		IdleConnTimeout:    time.Second * 10,
	})
if err != nil {
    // handle err
}
defer cli.Close()

db, err := cli.CreateDatabase(context.Background(), "DATABASE NAME")
```

### Examples

See [example](example) about how to use this package to communicate with TencentCloud VectorDB
