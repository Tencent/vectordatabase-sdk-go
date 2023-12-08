# TencentCloud VectorDB Go SDK

## Getting started

### Prerequisites
Go 1.15 or higher

### Install TencentCloud VectorDB Go SDK

1. Use `go get` to install the latest version of the TencentCloud VectorDB Go SDK and dependencies: 
```sh
go get -u github.com/tencent/vectordatabase-sdk-go/tcvectordb
```

2. Create New VectorDB Client To Start:
```go
import "github.com/tencent/vectordatabase-sdk-go/tcvectordb"

cli, err := tcvectordb.NewClient("http://127.0.0.1", "root", "key get from web console", &client.ClientOption{
		MaxIdldConnPerHost: 50,
		IdleConnTimeout:    time.Second * 30,
	})
if err != nil {
    // handle err
}
defer cli.Close()

db, err := cli.CreateDatabase(context.Background(), "DATABASE NAME")
```

### Examples

See [examples](example_test.go) about how to use this package to communicate with TencentCloud VectorDB