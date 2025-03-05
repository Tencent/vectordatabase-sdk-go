package main

import (
	"fmt"
	"log"

	"github.com/tencent/vectordatabase-sdk-go/tcvdbtext/encoder"
)

func main() {
	bm25, err := encoder.NewBM25EncoderByFiles(&encoder.BM25EncoderFileParams{
		StopWordsFile: "./stopwords.txt",
		//WordsFreqFile: "./bm25_zh_default.json",
		UserDictFile: "",
	})
	if err != nil {
		log.Fatalf(err.Error())
	}

	text := "什么是腾讯云向量数据库。"

	// 如需了解分词的情况，可参考下一行代码获取
	tokens := bm25.GetTokenizer().Tokenize(text)
	fmt.Println("tokens: ", tokens)

	// [EncodeText] can be used after set WordsFreqFile in [NewBM25EncoderByFiles]
	// sparse_vectors, err := bm25.EncodeText(text)
	// if err != nil {
	// 	log.Fatalf(err.Error())
	// }
	// fmt.Println("sparse vectors: ", sparse_vectors)
}
