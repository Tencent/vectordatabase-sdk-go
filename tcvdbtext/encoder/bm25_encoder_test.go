package encoder

import (
	"fmt"
	"log"
	"testing"

	"github.com/tencent/vectordatabase-sdk-go/tcvdbtext/tokenizer"
)

func Test_BM25Encoder_DownloadParams(t *testing.T) {
	bm25Encoder, _ := NewBM25Encoder(nil)
	bm25Encoder.DownloadParams("./bm25_params.json")
}

func Test_BM25Encoder_SetDefaultParams(t *testing.T) {
	bm25Encoder, _ := NewBM25Encoder(nil)
	err := bm25Encoder.SetDefaultParams("zh")
	if err != nil {
		println(err.Error())
	}
	bm25Encoder.DownloadParams("./bm25_params.json")
}

func Test_BM25Encoder_SetParams(t *testing.T) {
	bm25Encoder, _ := NewBM25Encoder(nil)
	err := bm25Encoder.SetDefaultParams("zh")
	if err != nil {
		println(err.Error())
	}
	bm25Encoder.DownloadParams("./bm25_params.json")
}

func Test_BM25Encoder_baseUsage(t *testing.T) {

	bm25, err := NewBM25Encoder(nil)
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = bm25.SetDefaultParams("zh")
	if err != nil {
		log.Fatalf(err.Error())
	}

	println("------------case 1: EncodeTexts------------")
	textVectors, err := bm25.EncodeTexts([]string{
		"腾讯云向量数据库（Tencent Cloud VectorDB）是一款全托管的自研企业级分布式数据库服务，专用于存储、索引、检索、管理由深度神经网络或其他机器学习模型生成的大量多维嵌入向量。",
		"作为专门为处理输入向量查询而设计的数据库，它支持多种索引类型和相似度计算方法，单索引支持10亿级向量规模，百万级 QPS 及毫秒级查询延迟。",
		"不仅能为大模型提供外部知识库，提高大模型回答的准确性，还可广泛应用于推荐系统、NLP 服务、计算机视觉、智能客服等 AI 领域。",
		"腾讯云向量数据库（Tencent Cloud VectorDB）作为一种专门存储和检索向量数据的服务提供给用户， 在高性能、高可用、大规模、低成本、简单易用、稳定可靠等方面体现出显著优势。",
		"腾讯云向量数据库可以和大语言模型 LLM 配合使用。企业的私域数据在经过文本分割、向量化后，可以存储在腾讯云向量数据库中，构建起企业专属的外部知识库，从而在后续的检索任务中，为大模型提供提示信息，辅助大模型生成更加准确的答案。",
		"腾讯云数据库托管机房分布在全球多个位置，这些位置节点称为地域（Region），每个地域又由多个可用区（Zone）构成。每个地域（Region）都是一个独立的地理区域。每个地域内都有多个相互隔离的位置，称为可用区（Zone）。每个可用区都是独立的，但同一地域下的可用区通过低时延的内网链路相连。腾讯云支持用户在不同位置分配云资源，建议用户在设计系统时考虑将资源放置在不同可用区以屏蔽单点故障导致的服务不可用状态。",
	})
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("encode multiple texts: %v\n", textVectors)

	println("------------case 2: EncodeQueries------------")
	queryVectors, err := bm25.EncodeQueries([]string{
		"什么是腾讯云向量数据库？",
		"腾讯云向量数据库有什么优势？",
		"腾讯云向量数据库能做些什么？",
	})
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("encode multiple queries: %v\n", queryVectors)
}

func Test_BM25Encoder_fitAndLoad(t *testing.T) {
	bm25, err := NewBM25Encoder(nil)
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = bm25.SetDefaultParams("zh")
	if err != nil {
		log.Fatalf(err.Error())
	}

	bm25.SetDict("../data/userdict_example.txt")
	fmt.Println("fit with your own corpus")

	bm25.FitCorpus([]string{
		"腾讯云向量数据库（tencent cloud vectordb）是一款全托管的自研企业级分布式数据库服务，专用于存储、索引、检索、管理由深度神经网络或其他机器学习模型生成的大量多维嵌入向量。",
		"作为专门为处理输入向量查询而设计的数据库，它支持多种索引类型和相似度计算方法，单索引支持10亿级向量规模，高达百万级 qps 及毫秒级查询延迟。",
		"不仅能为大模型提供外部知识库，提高大模型回答的准确性，还可广泛应用于推荐系统、nlp 服务、计算机视觉、智能客服等 AI 领域。",
		"腾讯云向量数据库（tencent cloud vectordb）作为一种专门存储和检索向量数据的服务提供给用户， 在高性能、高可用、大规模、低成本、简单易用、稳定可靠等方面体现出显著优势。 ",
		"腾讯云向量数据库可以和大语言模型 llm 配合使用。企业的私域数据在经过文本分割、向量化后，可以存储在腾讯云向量数据库中，构建起企业专属的外部知识库，从而在后续的检索任务中，为大模型提供提示信息，辅助大模型生成更加准确的答案。",
		"腾讯云数据库托管机房分布在全球多个位置，这些位置节点称为地域（region），每个地域又由多个可用区（zone）构成。每个地域（region）都是一个独立的地理区域。每个地域内都有多个相互隔离的位置，称为可用区（zone）。每个可用区都是独立的，但同一地域下的可用区通过低时延的内网链路相连。腾讯云支持用户在不同位置分配云资源，建议用户在设计系统时考虑将资源放置在不同可用区以屏蔽单点故障导致的服务不可用状态。",
	})
	fmt.Println("download bm25 params")
	bm25.DownloadParams("./bm25_params.json")
	fmt.Println("load bm25 params")
	bm25.SetParams("./bm25_params.json")

	query_vectors, err := bm25.EncodeQueries([]string{"什么是腾讯云向量数据库？", "腾讯云向量数据库有什么优势？", "腾讯云向量数据库能做些什么？"})
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("encode with your own fit params: %v\n", query_vectors)
}

func Test_BM25Encoder_PreciseMode(t *testing.T) {
	bm25, err := NewBM25Encoder(&BM25EncoderParams{Bm25Language: "zh"})
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println(bm25.GetTokenizer().Tokenize("什么是腾讯云向量数据库"))
}

func Test_BM25Encoder_CutAllMode(t *testing.T) {
	bm25, err := NewBM25Encoder(&BM25EncoderParams{Bm25Language: "zh"})
	if err != nil {
		log.Fatalf(err.Error())
	}

	cutAll := true
	jbtParams := tokenizer.TokenizerParams{
		CutAll: &cutAll,
	}
	err = bm25.GetTokenizer().UpdateParameters(jbtParams)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println(bm25.GetTokenizer().Tokenize("什么是腾讯云向量数据库"))
}

func Test_BM25Encoder_NoStopwords(t *testing.T) {
	bm25, err := NewBM25Encoder(&BM25EncoderParams{Bm25Language: "zh"})
	if err != nil {
		log.Fatalf(err.Error())
	}

	jbtParams := tokenizer.TokenizerParams{
		StopWords: false,
	}
	err = bm25.GetTokenizer().UpdateParameters(jbtParams)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println(bm25.GetTokenizer().Tokenize("什么是腾讯云向量数据库。"))
}

func Test_BM25Encoder_WithDefaultStopwords(t *testing.T) {
	bm25, err := NewBM25Encoder(&BM25EncoderParams{Bm25Language: "zh"})
	if err != nil {
		log.Fatalf(err.Error())
	}
	// 默认开启停用词
	jbtParams := tokenizer.TokenizerParams{
		StopWords: true,
	}
	err = bm25.GetTokenizer().UpdateParameters(jbtParams)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println(bm25.GetTokenizer().Tokenize("什么是腾讯云向量数据库。"))
}

func Test_BM25Encoder_WithUserDefineStopwords(t *testing.T) {
	bm25, err := NewBM25Encoder(&BM25EncoderParams{Bm25Language: "zh"})
	if err != nil {
		log.Fatalf(err.Error())
	}
	// 默认开启停用词
	jbtParams := tokenizer.TokenizerParams{
		StopWords: "../data/user_define_stopwords.txt",
	}
	err = bm25.GetTokenizer().UpdateParameters(jbtParams)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println(bm25.GetTokenizer().Tokenize("什么是腾讯云向量数据库。"))
}

func Test_BM25Encoder_WithUserDefineDict(t *testing.T) {
	bm25, err := NewBM25Encoder(&BM25EncoderParams{Bm25Language: "zh"})
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(bm25.GetTokenizer().Tokenize("什么是腾讯云向量数据库。"))

	err = bm25.SetDict("../data/userdict_example.txt")
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println(bm25.GetTokenizer().Tokenize("什么是腾讯云向量数据库。"))
}
