package tokenizer

import (
	"encoding/json"
	"log"
	"testing"
)

func Test_JiebaTokenizer(t *testing.T) {
	println("------------case 1------------")
	jbt, err := NewJiebaTokenizer(nil)
	if err != nil {
		log.Fatalln(err.Error())
	}

	res := jbt.IsStopWord("的")
	println(res)

	tokenizeRes := jbt.Tokenize("腾讯云的vdb是一款向量数据库")
	println(ToJson(tokenizeRes))
	encodeRes := jbt.Encode("腾讯云的vdb是一款向量数据库")
	println(ToJson(encodeRes))
	params := jbt.GetParameters()
	println(ToJson(params))

	println("------------case 2------------")
	jbtParams := TokenizerParams{
		UserDictFilePath: "../data/userdict_example.txt",
		StopWords:        true,
	}
	jbt.UpdateParameters(jbtParams)
	tokenizeRes = jbt.Tokenize("腾讯云的vdb是一款向量数据库")
	println(ToJson(tokenizeRes))
	encodeRes = jbt.Encode("腾讯云的vdb是一款向量数据库")
	println(ToJson(encodeRes))

	params = jbt.GetParameters()
	println(ToJson(params))

}

func ToJson(any interface{}) string {
	bytes, err := json.Marshal(any)
	if err != nil {
		return ""
	}
	return string(bytes)
}
