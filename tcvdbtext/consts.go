package tcvdbtext

import "os"

const (
	Mmh3HashName string = "mmh3_hash"
)

const (
	DefaultBM25EncoderB  = 0.75
	DefaultBM25EncoderK1 = 1.2
)

const (
	DefaultStorageDir  string = "/tmp/tencent/vectordatabase/data/"
	CosSparsevectorDir string = "https://vectordb-public-1310738255.cos.ap-guangzhou.myqcloud.com/sparsevector/"

	DefaultStopWordsFileName string = "default_stopwords.txt"
)

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
