package tokenizer

type Tokenizer interface {
	Tokenize(sentence string) []string
	Encode(sentence string) []int64
	IsStopWord(word string) bool
	UpdateParameters(params TokenizerParams) error
	GetParameters() TokenizerParams
	SetDict(dictFile string) error
}

type TokenizerParams struct {
	HashFunction string `json:"hash_function,omitempty"`
	// bool type: stopWords enbale
	// string type: stopWords filePath
	StopWords interface{} `json:"stop_words,omitempty"`

	UserDictFilePath string `json:"dict_file,omitempty"`
	CutAll           *bool  `json:"cut_all,omitempty"`
	ForSearch        *bool  `json:"for_search,omitempty"`
	Hmm              *bool  `json:"HMM,omitempty"`
	UsePaddle        *bool  `json:"use_paddle,omitempty"`
}
