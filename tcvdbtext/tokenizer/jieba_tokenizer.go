//go:build !windows
// +build !windows

package tokenizer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	tcvdbtext "github.com/tencent/vectordatabase-sdk-go/tcvdbtext"
	"github.com/tencent/vectordatabase-sdk-go/tcvdbtext/hash"
	"github.com/yanyiwu/gojieba"
)

const (
	STOP_WORD_PATH = "/../data/stopwords.txt"
)

type JiebaTokenizer struct {
	forSearch bool
	cutAll    bool
	useHmm    bool
	lowerCase bool

	UserDictFilePath  string
	StopWordsFilePath string
	StopWordsEnable   bool

	Jieba        *gojieba.Jieba
	hashFunc     hash.HashInterface
	stopWordsMap map[string]bool
}

func NewJiebaTokenizer(params *TokenizerParams) (Tokenizer, error) {
	defaultForSearch := false
	defaultCutAll := false
	defaultUseHmm := true
	defaultLowerCase := false
	defaultStopWordsEnable := true

	jbt := new(JiebaTokenizer)
	jbt.forSearch = defaultForSearch
	jbt.cutAll = defaultCutAll
	jbt.useHmm = defaultUseHmm
	jbt.lowerCase = defaultLowerCase
	jbt.StopWordsEnable = defaultStopWordsEnable
	jbt.stopWordsMap = make(map[string]bool, 0)

	if params == nil {
		err := jbt.refreshStopWordsMap()
		if err != nil {
			return nil, fmt.Errorf("open file %v for stopwords failed. err: %v", jbt.StopWordsFilePath, err.Error())
		}

		err = jbt.refreshJieba()
		if err != nil {
			return nil, fmt.Errorf("new Tokenizer failed, because refreshing jieba failed. err: %v", err.Error())
		}
		jbt.hashFunc = hash.NewMmh3Hash()
		return jbt, nil
	}

	if params.ForSearch != nil {
		jbt.forSearch = *params.ForSearch
	}

	if params.CutAll != nil {
		jbt.cutAll = *params.ForSearch
	}

	if params.Hmm != nil {
		jbt.useHmm = *params.Hmm
	}

	if params.LowerCase != nil {
		jbt.lowerCase = *params.LowerCase
	}

	boolV, ok := params.StopWords.(bool)
	if ok {
		jbt.StopWordsEnable = boolV
	} else {
		stringV, ok := params.StopWords.(string)
		if ok {
			if stringV != "" && !tcvdbtext.FileExists(stringV) {
				return nil, fmt.Errorf("the StopWordsFilePath in params is invalid, "+
					"because the filepath %v doesn't exist", stringV)
			}
			jbt.StopWordsFilePath = stringV
		}
	}

	if params.UserDictFilePath != "" && !tcvdbtext.FileExists(params.UserDictFilePath) {
		return nil, fmt.Errorf("the UserDictFilePath in params is invalid, "+
			"because the filepath %v doesn't exist", params.UserDictFilePath)
	}

	jbt.UserDictFilePath = params.UserDictFilePath

	jbt.stopWordsMap = make(map[string]bool, 0)
	err := jbt.refreshStopWordsMap()
	if err != nil {
		return nil, fmt.Errorf("open file %v for stopwords failed. err: %v", jbt.StopWordsFilePath, err.Error())
	}

	err = jbt.refreshJieba()
	if err != nil {
		return nil, fmt.Errorf("new Tokenizer failed, because refreshing jieba failed. err: %v", err.Error())
	}

	if params.HashFunction == "" || params.HashFunction == tcvdbtext.Mmh3HashName {
		jbt.hashFunc = hash.NewMmh3Hash()
	} else {
		return nil, fmt.Errorf("not support the hash %v method", params.HashFunction)
	}

	return jbt, nil
}

func (jbt *JiebaTokenizer) UpdateParameters(params TokenizerParams) error {
	if params.ForSearch != nil {
		jbt.forSearch = *params.ForSearch
	}

	if params.CutAll != nil {
		jbt.cutAll = *params.CutAll
	}

	if params.Hmm != nil {
		jbt.useHmm = *params.Hmm
	}

	if params.LowerCase != nil {
		jbt.lowerCase = *params.LowerCase
	}

	if params.HashFunction != "" {
		if params.HashFunction == tcvdbtext.Mmh3HashName {
			jbt.hashFunc = hash.NewMmh3Hash()
		} else {
			return fmt.Errorf("not support the hash %v method", params.HashFunction)
		}
	}

	needRefreshJieba := false
	needRefreshStopWordsMap := false
	if params.UserDictFilePath != "" && jbt.UserDictFilePath != params.UserDictFilePath {
		if !tcvdbtext.FileExists(params.UserDictFilePath) {
			return fmt.Errorf("the UserDictFilePath in params is invalid, "+
				"because the filepath %v doesn't exist", params.UserDictFilePath)
		}

		jbt.UserDictFilePath = params.UserDictFilePath
		needRefreshJieba = true
	}

	boolV, ok := params.StopWords.(bool)
	if ok {
		jbt.StopWordsEnable = boolV
		needRefreshStopWordsMap = true
	} else {
		stringV, ok := params.StopWords.(string)
		if ok {
			if stringV != "" && !tcvdbtext.FileExists(stringV) {
				return fmt.Errorf("the StopWordsFilePath in params is invalid, "+
					"because the filepath %v doesn't exist", stringV)
			}
			jbt.StopWordsFilePath = stringV
			needRefreshStopWordsMap = true
		}
	}

	if needRefreshStopWordsMap {
		jbt.refreshStopWordsMap()
		needRefreshJieba = true
	}

	if needRefreshJieba {
		err := jbt.refreshJieba()
		if err != nil {
			return fmt.Errorf("update parameters failed, because refreshing jieba failed. err: %v", err.Error())
		}
	}

	return nil
}

func (jbt *JiebaTokenizer) GetParameters() TokenizerParams {
	forSearch := jbt.forSearch
	CutAll := jbt.cutAll
	UseHmm := jbt.useHmm
	lowerCase := jbt.lowerCase
	var stopWords interface{}

	if jbt.StopWordsFilePath != "" {
		stopWords = jbt.StopWordsFilePath
	} else {
		stopWords = jbt.StopWordsEnable
	}

	return TokenizerParams{
		ForSearch:        &forSearch,
		CutAll:           &CutAll,
		Hmm:              &UseHmm,
		LowerCase:        &lowerCase,
		UserDictFilePath: jbt.UserDictFilePath,
		StopWords:        stopWords,
		HashFunction:     jbt.hashFunc.GetHashFuctionName(),
	}
}

func (jbt *JiebaTokenizer) refreshJieba() error {
	_, filePath, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filePath)

	stopWordPath := ""
	if jbt.StopWordsEnable {
		if jbt.StopWordsFilePath != "" {
			stopWordPath = jbt.StopWordsFilePath
		} else {
			stopWordPath = dir + STOP_WORD_PATH
		}
	}

	jbt.Jieba = gojieba.NewJieba("", "", jbt.UserDictFilePath, "", stopWordPath)
	return nil
}

func (jbt *JiebaTokenizer) refreshStopWordsMap() error {
	if !jbt.StopWordsEnable {
		jbt.stopWordsMap = make(map[string]bool, 0)
		return nil
	}

	_, filePath, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filePath)

	stopWordFilePath := dir + STOP_WORD_PATH
	if jbt.StopWordsFilePath != "" {
		stopWordFilePath = jbt.StopWordsFilePath
	}

	file, err := os.Open(stopWordFilePath)
	if err != nil {
		return fmt.Errorf("open file %v failed. err: %v", stopWordFilePath, err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimRightFunc(line, func(r rune) bool {
			return r == ' ' || r == '\t' || r == '\n' || r == '\r'
		})
		jbt.stopWordsMap[trimmedLine] = true
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read file %v failed. err: %v", stopWordFilePath, err.Error())
	}

	return nil
}

func (jbt *JiebaTokenizer) Tokenize(sentence string) []string {
	if len(sentence) == 0 {
		return []string{}
	}
	if jbt.lowerCase {
		sentence = strings.ToLower(sentence)
	}

	var segs []string
	var words []string
	if jbt.forSearch {
		segs = jbt.Jieba.CutForSearch(sentence, jbt.useHmm)
	} else if jbt.cutAll {
		segs = jbt.Jieba.CutAll(sentence)
	} else {
		segs = jbt.Jieba.Cut(sentence, jbt.useHmm)
	}
	for _, word := range segs {
		if len(word) == 0 || word == " " || jbt.IsStopWord(word) {
			continue
		}
		//print(word + " ")
		words = append(words, word)
	}
	//println()
	return words

}

func (jbt *JiebaTokenizer) Encode(sentence string) []int64 {
	var tokens []int64
	words := jbt.Tokenize(sentence)
	for _, word := range words {
		tokens = append(tokens, jbt.hashFunc.Hash(word))
	}
	return tokens
}

func (jbt *JiebaTokenizer) IsStopWord(word string) bool {
	if _, ok := jbt.stopWordsMap[word]; ok {
		return true
	}
	return false
}

func (jbt *JiebaTokenizer) SetDict(dictFile string) error {
	jbt.UserDictFilePath = dictFile
	err := jbt.refreshJieba()
	if err != nil {
		return fmt.Errorf("set dictionary failed, because refreshing jieba failed. err: %v", err.Error())
	}
	return nil
}
