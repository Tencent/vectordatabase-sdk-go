# Tencent VectorDB Sparse Encoder SDK

Go SDK for [Tencent VectorDB Sparse Encoder](https://cloud.tencent.com/document/product/1709/111372).

## Getting started

### Prerequisites
1. Go 1.17 or higher

### Install TencentCloud VectorDB Go SDK

1. Use `go get` to install the latest version of the TencentCloud VectorDB Sparse Encoder SDK dependencies: 
```sh
go get -u github.com/tencent/vectordatabase-sdk-go/tcvdbtext
```

2. Try [sparse_vector_demo](examples/sparse_vector_demo/main.go) in an online environment with internet access.

3. Try [sparse_vector_offline_demo](examples/sparse_vector_offline_demo/main.go) in an offline environment without internet access. 
Before running the code, please download files which you need.

    - [Chinese Words Frequency File](https://vectordb-public-1310738255.cos.ap-guangzhou.myqcloud.com/sparsevector/bm25_zh_default.json)
    - [English Words Frequency File](https://vectordb-public-1310738255.cos.ap-guangzhou.myqcloud.com/sparsevector/bm25_en_default.json)
    - [Default Stopwords File](https://vectordb-public-1310738255.cos.ap-guangzhou.myqcloud.com/sparsevector/default_stopwords.txt)
