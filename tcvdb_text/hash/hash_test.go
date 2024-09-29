package hash

import "testing"

func Test_mmh3(t *testing.T) {
	mmh3 := NewMmh3Hash()
	res := mmh3.Hash("腾讯云")
	println(res)
}
