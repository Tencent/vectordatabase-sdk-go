package utils

import (
	"encoding/json"
	"testing"

	"log"
)

func Test_ConvertBinaryArray2Uint8Array(t *testing.T) {
	binaryA := []byte{1, 1, 1, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0}
	res, err := BinaryToUint8(binaryA)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	for _, v := range res {
		print(v, " ")
	}
	println()
}

func ToJson(any interface{}) string {
	bytes, err := json.Marshal(any)
	if err != nil {
		return ""
	}
	return string(bytes)
}
