// Copyright (C) 2023 Tencent Cloud.
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the vectordb-sdk-java), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is furnished
// to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
// PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
// SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package entity

import (
	"time"
)

// Collection wrap the collection parameters and document interface to operating the document api
type Collection struct {
	DocumentInterface `json:"-"`
	IndexInterface    `json:"-"`
	DatabaseName      string
	CollectionName    string
	DocumentCount     int64
	Alias             []string
	ShardNum          uint32
	ReplicasNum       uint32
	Indexes           Indexes
	IndexStatus       IndexStatus
	Embedding         Embedding
	Description       string
	Size              uint64
	CreateTime        time.Time
}

func (c *Collection) Debug(v bool) {
	c.DocumentInterface.Debug(v)
}

func (c *Collection) WithTimeout(t time.Duration) {
	c.DocumentInterface.WithTimeout(t)
}

type CreateCollectionOption struct {
	Embedding *Embedding
}

type CreateCollectionResult struct{}

type DescribeCollectionOption struct {
}

type DescribeCollectionResult struct {
	Collection
}

type ListCollectionOption struct{}

type ListCollectionResult struct {
	Collections []*Collection
}

type DropCollectionOption struct{}

type DropCollectionResult struct {
	AffectedCount int
}

type TruncateCollectionOption struct{}

type TruncateCollectionResult struct {
	AffectedCount int
}
