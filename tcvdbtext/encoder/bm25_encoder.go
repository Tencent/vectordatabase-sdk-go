package encoder

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	tcvdbtext "github.com/tencent/vectordatabase-sdk-go/tcvdbtext"
	"github.com/tencent/vectordatabase-sdk-go/tcvdbtext/tokenizer"
)

const (
	BM25Params_ZH_Path = "/../data/bm25_zh_default.json"
	BM25Params_EN_Path = "/../data/bm25_en_default.json"

	BM25_ZH_CONTENT = "zh"
	BM25_EN_CONTENT = "en"
)

type BM25Encoder struct {
	B  float64 `json:"b,omitempty"`
	K1 float64 `json:"k1,omitempty"`
	BM25LearnedParams

	Tokenizer tokenizer.Tokenizer
}

type BM25EncoderParams struct {
	B            *float64
	K1           *float64
	Tokenizer    tokenizer.Tokenizer
	Bm25Language string
}

type BM25LearnedParams struct {
	TokenFreq        map[string]float64 `json:"token_freq,omitempty"`
	DocCount         int64              `json:"doc_count,omitempty"`
	AverageDocLength float64            `json:"average_doc_length,omitempty"`
}

type BM25Params struct {
	B  *float64 `json:"b,omitempty"`
	K1 *float64 `json:"k1,omitempty"`
	tokenizer.TokenizerParams
	BM25LearnedParams
}

func NewBM25Encoder(params *BM25EncoderParams) (SparseEncoder, error) {
	bm25 := new(BM25Encoder)

	bm25.B = tcvdbtext.DefaultBM25EncoderB
	bm25.K1 = tcvdbtext.DefaultBM25EncoderK1

	if params != nil {
		if params.B != nil {
			bm25.B = *params.B
		}
		if params.K1 != nil {
			bm25.K1 = *params.K1
		}
		bm25.Tokenizer = params.Tokenizer
	}

	var err error
	if bm25.Tokenizer == nil {
		bm25.Tokenizer, err = tokenizer.NewJiebaTokenizer(nil)
		if err != nil {
			return nil, err
		}
	}

	if params != nil && params.Bm25Language != "" {
		err := bm25.SetDefaultParams(params.Bm25Language)
		if err != nil {
			return nil, err
		}
	}

	return bm25, nil
}

func (bm25 *BM25Encoder) GetTokenizer() tokenizer.Tokenizer {
	return bm25.Tokenizer
}

func (bm25 *BM25Encoder) SetDefaultParams(bm25Language string) error {
	_, filePath, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filePath)

	bm25ParamsPath := ""
	if bm25Language == BM25_ZH_CONTENT {
		bm25ParamsPath = dir + BM25Params_ZH_Path
	} else if bm25Language == BM25_EN_CONTENT {
		bm25ParamsPath = dir + BM25Params_EN_Path
	} else {
		return fmt.Errorf("input name be 'zh' or 'en'")
	}

	err := bm25.SetParams(bm25ParamsPath)
	if err != nil {
		return fmt.Errorf("use default settings file for language %v to set params failed, err: %v",
			bm25Language, err.Error())
	}

	return nil
}

func (bm25 *BM25Encoder) SetParams(paramsFileLoadPath string) error {
	if !tcvdbtext.FileExists(paramsFileLoadPath) {
		return fmt.Errorf("the filepath %v doesn't exist", paramsFileLoadPath)
	}
	data, err := os.ReadFile(paramsFileLoadPath)
	if err != nil {
		return fmt.Errorf("cannot read file: %v", err)
	}

	bm25ParamsByFile := new(BM25Params)
	err = json.Unmarshal(data, bm25ParamsByFile)
	if err != nil {
		return fmt.Errorf("cannot parse file %v to JSON, err: %v", paramsFileLoadPath, err.Error())
	}

	bm25.B = *bm25ParamsByFile.B
	bm25.K1 = *bm25ParamsByFile.K1
	bm25.BM25LearnedParams = bm25ParamsByFile.BM25LearnedParams

	bm25.Tokenizer.UpdateParameters(tokenizer.TokenizerParams{
		ForSearch: bm25ParamsByFile.ForSearch,
		CutAll:    bm25ParamsByFile.CutAll,
		Hmm:       bm25ParamsByFile.Hmm,

		UserDictFilePath: bm25ParamsByFile.UserDictFilePath,
		StopWords:        bm25ParamsByFile.StopWords,

		HashFunction: bm25ParamsByFile.HashFunction,
	})
	return nil
}

func (bm25 *BM25Encoder) DownloadParams(paramsFileDownloadPath string) error {
	tokenizerParams := bm25.Tokenizer.GetParameters()
	bm25Params := new(BM25Params)
	bm25Params.TokenizerParams = tokenizerParams
	bm25Params.BM25LearnedParams = bm25.BM25LearnedParams
	bm25Params.B = &bm25.B
	bm25Params.K1 = &bm25.K1

	jsonData, err := json.MarshalIndent(bm25Params, "", "  ")
	if err != nil {
		return fmt.Errorf("download bm25 params failed because marshal params failed. err: %v", err.Error())
	}

	err = os.WriteFile(paramsFileDownloadPath, jsonData, os.ModePerm)
	if err != nil {
		return fmt.Errorf("download bm25 params failed because write file failed. err: %v", err.Error())
	}

	return nil
}

func (bm25 *BM25Encoder) encodeSingleDocument(text string) []SparseVecItem {
	hashTokens, counts := bm25.tf(text)
	var Sum int64
	for _, count := range counts {
		Sum += count
	}

	sparseVector := make([]SparseVecItem, 0)
	tfNormed := make([]float64, len(counts))
	for i, count := range counts {
		tfNormed[i] = float64(count) / ((bm25.K1)*(1.0-bm25.B+bm25.B*(float64(Sum)/bm25.AverageDocLength)) + float64(count))
	}

	for i, v := range tfNormed {
		sparseVector = append(sparseVector, SparseVecItem{hashTokens[i], float32(v)})
	}

	return sparseVector
}

func (bm25 *BM25Encoder) EncodeTexts(texts []string) ([][]SparseVecItem, error) {
	if bm25.AverageDocLength == 0 || bm25.DocCount == 0 || len(bm25.TokenFreq) == 0 {
		return nil, fmt.Errorf("BM25 must be fit before encoding documents")
	}
	sparseVectors := make([][]SparseVecItem, 0)
	for _, text := range texts {
		sparseVectors = append(sparseVectors, bm25.encodeSingleDocument(text))
	}
	return sparseVectors, nil
}

func (bm25 *BM25Encoder) EncodeText(text string) ([]SparseVecItem, error) {
	if bm25.AverageDocLength == 0 || bm25.DocCount == 0 || len(bm25.TokenFreq) == 0 {
		return nil, fmt.Errorf("BM25 must be fit before encoding documents")
	}

	return bm25.encodeSingleDocument(text), nil
}

func (bm25 *BM25Encoder) encodeSingleQuery(text string) []SparseVecItem {
	hashTokens, _ := bm25.tf(text)
	df := make([]float64, len(hashTokens))

	for i, hashToken := range hashTokens {
		df[i] = float64(bm25.TokenFreq[strconv.FormatInt(hashToken, 10)])
	}

	idf := make([]float64, len(df))
	for i, d := range df {
		idf[i] = math.Log(float64(bm25.DocCount+1) / (d + 0.5))
	}

	idfSum := 0.0
	for _, val := range idf {
		idfSum += val
	}

	idfNorm := make([]float64, len(idf))
	for i, val := range idf {
		idfNorm[i] = val / idfSum
	}

	sparseVector := make([]SparseVecItem, 0)

	for i, v := range idfNorm {
		sparseVector = append(sparseVector, SparseVecItem{hashTokens[i], float32(v)})
	}

	return sparseVector

}

func (bm25 *BM25Encoder) EncodeQueries(texts []string) ([][]SparseVecItem, error) {
	if bm25.AverageDocLength == 0 || bm25.DocCount == 0 || len(bm25.TokenFreq) == 0 {
		return nil, fmt.Errorf("BM25 must be fit before encoding documents")
	}
	sparseVectors := make([][]SparseVecItem, 0)
	for _, text := range texts {
		sparseVectors = append(sparseVectors, bm25.encodeSingleQuery(text))
	}
	return sparseVectors, nil
}

func (bm25 *BM25Encoder) EncodeQuery(text string) ([]SparseVecItem, error) {
	if bm25.AverageDocLength == 0 || bm25.DocCount == 0 || len(bm25.TokenFreq) == 0 {
		return nil, fmt.Errorf("BM25 must be fit before encoding documents")
	}

	return bm25.encodeSingleQuery(text), nil
}

func (bm25 *BM25Encoder) FitCorpus(corpus []string) error {
	var docNum int64
	var sumDocLen int64
	tokenFreqCounter := make(map[string]float64)

	for _, doc := range corpus {
		indices, tf := bm25.tf(doc)
		if len(indices) == 0 {
			continue
		}
		docNum++
		var sumTf int64
		sumTf = 0
		for _, v := range tf {
			sumTf += v
		}
		sumDocLen += sumTf

		// Convert indices to strings and update token frequency counter
		for _, index := range indices {
			tokenStr := fmt.Sprintf("%d", index)
			tokenFreqCounter[tokenStr]++
		}
	}

	if bm25.TokenFreq == nil || bm25.DocCount == 0 || bm25.AverageDocLength == 0 {
		bm25.TokenFreq = tokenFreqCounter
		bm25.DocCount = docNum
		bm25.AverageDocLength = float64(sumDocLen) / float64(docNum)
	} else {
		bm25.AverageDocLength = (bm25.AverageDocLength*float64(bm25.DocCount) + float64(sumDocLen)) / float64(bm25.DocCount+docNum)
		bm25.DocCount += docNum
		for k, v := range tokenFreqCounter {
			bm25.TokenFreq[k] += v
		}
	}

	return nil
}

func (bm25 *BM25Encoder) SetDict(CustomDictLoadPath string) error {
	return bm25.Tokenizer.SetDict(CustomDictLoadPath)
}

func (bm25 *BM25Encoder) tf(text string) ([]int64, []int64) {
	tokenizer := bm25.Tokenizer
	tokens := tokenizer.Encode(text)

	counter := make(map[int64]int64, 0)
	for _, token := range tokens {
		if _, ok := counter[token]; !ok {
			counter[token] = 1
		} else {
			counter[token]++
		}
	}

	deduplicatedTokens := make([]int64, 0)
	fres := make([]int64, 0)
	for token, fre := range counter {
		deduplicatedTokens = append(deduplicatedTokens, token)
		fres = append(fres, fre)
	}

	return deduplicatedTokens, fres

}
