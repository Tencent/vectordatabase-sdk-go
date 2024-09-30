package encoder

import "github.com/tencent/vectordatabase-sdk-go/tcvdbtext/tokenizer"

type SparseEncoder interface {
	encodeSingleDocument(text string) []SparseVecItem
	EncodeTexts(texts []string) ([][]SparseVecItem, error)
	EncodeText(text string) ([]SparseVecItem, error)

	encodeSingleQuery(text string) []SparseVecItem
	EncodeQueries(texts []string) ([][]SparseVecItem, error)
	EncodeQuery(text string) ([]SparseVecItem, error)

	FitCorpus(corpus []string) error
	DownloadParams(paramsFileDownloadPath string) error

	SetParams(paramsFileLoadPath string) error
	SetDefaultParams(bm25Language string) error

	SetDict(CustomDictLoadPath string) error

	GetTokenizer() tokenizer.Tokenizer
}

type SparseVecItem struct {
	TermId int64
	Score  float32
}
