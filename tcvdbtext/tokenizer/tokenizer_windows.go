//go:build windows
// +build windows

package tokenizer

import "errors"

func NewJiebaTokenizer(params *TokenizerParams) (Tokenizer, error) {
	return nil, errors.New("unsupported windows jieba tokenizer, Please use Linux instead or contact us.")
}
