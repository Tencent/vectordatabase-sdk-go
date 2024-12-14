package utils

import (
	"errors"
	"strconv"
)

func BinaryToUint8(binaryArray []byte) ([]float32, error) {
	binaryArrayLen := len(binaryArray)
	if binaryArrayLen%8 != 0 {
		return nil, errors.New("the length of the binaryArray must be a multiple of 8")
	}
	uint8Array := make([]float32, 0)
	for i := 0; i < binaryArrayLen/8; i++ {
		arr := binaryArray[i*8 : (i+1)*8]
		binaryString := ""
		for _, bit := range arr {
			binaryString += strconv.Itoa(int(bit))
		}
		decimalValue, err := strconv.ParseInt(binaryString, 2, 64)
		if err != nil {
			return nil, err
		}
		uint8Value := uint8(decimalValue) & 0xFF
		uint8Array = append(uint8Array, float32(uint8Value))
	}
	return uint8Array, nil
}
