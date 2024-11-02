package tokenizer

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/go-ego/gse"

	tcvdbtext "github.com/tencent/vectordatabase-sdk-go/tcvdbtext"
	"github.com/tencent/vectordatabase-sdk-go/tcvdbtext/hash"
)

const (
	STOP_WORD_PATH = "/../data/stopwords.txt"
)

type JiebaTokenizer struct {
	forSearch bool
	cutAll    bool
	useHmm    bool

	UserDictFilePath  string
	StopWordsFilePath string
	StopWordsEnable   bool

	Jieba    *gse.Segmenter
	hashFunc hash.HashInterface
}

func NewJiebaTokenizer(params *TokenizerParams) (Tokenizer, error) {
	defaultForSearch := false
	defaultCutAll := false
	defaultUseHmm := true
	defaultStopWordsEnable := true

	jbt := new(JiebaTokenizer)
	jbt.forSearch = defaultForSearch
	jbt.cutAll = defaultCutAll
	jbt.useHmm = defaultUseHmm
	jbt.StopWordsEnable = defaultStopWordsEnable
	jbt.Jieba = new(gse.Segmenter)

	_, filePath, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filePath)

	defaultStopWordFilePath := dir + STOP_WORD_PATH

	if params == nil {
		log.Printf("[Waring] Jieba will use default file for stopwords, which is %v", defaultStopWordFilePath)
		jbt.StopWordsFilePath = defaultStopWordFilePath
		err := jbt.Jieba.LoadStop(defaultStopWordFilePath)
		if err != nil {
			return nil, fmt.Errorf("jieba loads file %v for stopwords failed. err: %v", jbt.StopWordsFilePath, err.Error())
		}

		jbt.Jieba.LoadDict()
		jbt.hashFunc = hash.NewMmh3Hash()
		return jbt, nil
	}

	if params.ForSearch != nil {
		jbt.forSearch = *params.ForSearch
	}

	if params.CutAll != nil {
		jbt.cutAll = *params.CutAll
	}

	if params.Hmm != nil {
		jbt.useHmm = *params.Hmm
	}

	stopWordsEnable, ok := params.StopWords.(bool)
	if ok {
		jbt.StopWordsEnable = stopWordsEnable
		if stopWordsEnable {
			log.Printf("[Waring] Jieba will use default file for stopwords, which is %v", defaultStopWordFilePath)
			jbt.StopWordsFilePath = defaultStopWordFilePath
			err := jbt.Jieba.LoadStop(defaultStopWordFilePath)
			if err != nil {
				return nil, fmt.Errorf("jieba loads file %v for stopwords failed. err: %v", jbt.StopWordsFilePath, err.Error())
			}
		}
	} else {
		stopWordFilePath, ok := params.StopWords.(string)
		if ok {
			jbt.StopWordsFilePath = stopWordFilePath
			err := jbt.Jieba.LoadStop(jbt.StopWordsFilePath)
			if err != nil {
				return nil, fmt.Errorf("jieba loads file %v for stopwords failed. err: %v", jbt.StopWordsFilePath, err.Error())
			}
		}
	}

	if params.UserDictFilePath != "" && !tcvdbtext.FileExists(params.UserDictFilePath) {
		return nil, fmt.Errorf("the UserDictFilePath in params is invalid, "+
			"because the filepath %v doesn't exist", params.UserDictFilePath)
	}

	jbt.UserDictFilePath = params.UserDictFilePath
	jbt.Jieba.LoadDict(jbt.UserDictFilePath)

	if params.HashFunction == "" || params.HashFunction == tcvdbtext.Mmh3HashName {
		jbt.hashFunc = hash.NewMmh3Hash()
	} else {
		return nil, fmt.Errorf("not support the hash %v method", params.HashFunction)
	}

	return jbt, nil
}

func (jbt *JiebaTokenizer) Tokenize(sentence string) []string {
	if len(sentence) == 0 {
		return []string{}
	}

	var segs []string
	var words []string
	if jbt.forSearch {
		segs = jbt.Jieba.CutSearch(sentence, jbt.useHmm)
	} else if jbt.cutAll {
		segs = jbt.Jieba.CutAll(sentence)
	} else {
		segs = jbt.Jieba.Cut(sentence, jbt.useHmm)
	}
	for _, word := range segs {
		if len(word) == 0 || word == " " || jbt.Jieba.IsStop(word) {
			continue
		}
		words = append(words, word)
	}
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
	if jbt.Jieba == nil {
		return false
	}
	return jbt.Jieba.IsStop(word)
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

	if params.HashFunction != "" {
		if params.HashFunction == tcvdbtext.Mmh3HashName {
			jbt.hashFunc = hash.NewMmh3Hash()
		} else {
			return fmt.Errorf("not support the hash %v method", params.HashFunction)
		}
	}

	if params.UserDictFilePath != "" {
		if !tcvdbtext.FileExists(params.UserDictFilePath) {
			return fmt.Errorf("the UserDictFilePath in params is invalid, "+
				"because the filepath %v doesn't exist", params.UserDictFilePath)
		}

		jbt.UserDictFilePath = params.UserDictFilePath
		newJieba, err := gse.New(jbt.UserDictFilePath)
		if err != nil {
			return fmt.Errorf("jieba loads file %v for userdict failed. err: %v", jbt.UserDictFilePath, err.Error())
		}
		jbt.Jieba = &newJieba

	}

	if params.StopWords != nil {
		stopWordsEnable, ok := params.StopWords.(bool)
		if ok {
			jbt.StopWordsEnable = stopWordsEnable
		} else {
			stopWordFilePath, stringOk := params.StopWords.(string)
			if stringOk {
				jbt.StopWordsFilePath = stopWordFilePath
			}
		}
	}

	if jbt.StopWordsFilePath != "" {
		err := jbt.Jieba.LoadStop(jbt.StopWordsFilePath)
		if err != nil {
			return fmt.Errorf("jieba loads file %v for stopwords failed. err: %v", jbt.StopWordsFilePath, err.Error())
		}
	} else if jbt.StopWordsEnable {
		_, filePath, _, _ := runtime.Caller(0)
		dir := filepath.Dir(filePath)
		defaultStopWordFilePath := dir + STOP_WORD_PATH

		log.Printf("[Waring] Jieba will use default file for stopwords, which is %v", defaultStopWordFilePath)
		jbt.StopWordsFilePath = defaultStopWordFilePath
		err := jbt.Jieba.LoadStop(defaultStopWordFilePath)
		if err != nil {
			return fmt.Errorf("jieba loads file %v for stopwords failed. err: %v", jbt.StopWordsFilePath, err.Error())
		}
	}

	return nil
}
func (jbt *JiebaTokenizer) GetParameters() TokenizerParams {
	forSearch := jbt.forSearch
	CutAll := jbt.cutAll
	UseHmm := jbt.useHmm
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
		UserDictFilePath: jbt.UserDictFilePath,
		StopWords:        stopWords,
		HashFunction:     jbt.hashFunc.GetHashFuctionName(),
	}
}
func (jbt *JiebaTokenizer) SetDict(dictFile string) error {
	err := jbt.Jieba.LoadDict(dictFile)
	if err != nil {
		return fmt.Errorf("set dictionary failed, because refreshing jieba failed. err: %v", err.Error())
	}
	jbt.UserDictFilePath = dictFile
	return nil
}
