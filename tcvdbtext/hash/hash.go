package hash

import (
	"github.com/spaolacci/murmur3"
	tcvdbtext "github.com/tencent/vectordatabase-sdk-go/tcvdbtext"
)

type HashInterface interface {
	Hash(text string) int64
	GetHashFuctionName() string
}

type Mmh3 struct{}

func NewMmh3Hash() HashInterface {
	return &Mmh3{}
}

// Use mmh3 to hash text to 32-bit unsigned integer
func (m *Mmh3) Hash(text string) int64 {
	hasher := murmur3.New32()
	hasher.Write([]byte(text))
	hashValue := hasher.Sum32()
	return int64(hashValue)
}

func (m *Mmh3) GetHashFuctionName() string {
	return tcvdbtext.Mmh3HashName
}
