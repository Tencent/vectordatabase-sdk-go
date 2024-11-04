# TencentCloud VectorDB Go SDK

## Getting started

### Prerequisites
1. Go 1.15 or higher
2. Only support to use bm25 tcvdbtext package when the system is macos or linux, and you should install gcc firstly. Because gojieba in tcvdbtext package use cgo feature in go.  


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
