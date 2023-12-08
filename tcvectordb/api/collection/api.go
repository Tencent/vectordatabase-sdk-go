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

package collection

import "github.com/tencent/vectordatabase-sdk-go/tcvectordb/api"

// CreateReq create collection request
type CreateReq struct {
	api.Meta    `path:"/collection/create" tags:"Collection" method:"Post" summary:"创建collection"`
	Database    string             `json:"database,omitempty"`
	Collection  string             `json:"collection,omitempty"`
	ReplicaNum  uint32             `json:"replicaNum,omitempty"`
	ShardNum    uint32             `json:"shardNum,omitempty"`
	Size        uint64             `json:"size,omitempty"`
	CreateTime  string             `json:"createTime,omitempty"`
	Description string             `json:"description,omitempty"`
	Indexes     []*api.IndexColumn `json:"indexes,omitempty"`
	IndexStatus *IndexStatus       `json:"indexStatus,omitempty"`
	AliasList   []string           `json:"alias_list,omitempty"`
	Embedding   Embedding          `json:"embedding"`
}

type Embedding struct {
	Field       string `json:"field,omitempty"`
	VectorField string `json:"vectorField,omitempty"`
	Model       string `json:"model,omitempty"`
}

type IndexStatus struct {
	Status    string `json:"status,omitempty"`
	Progress  string `json:"progress,omitempty"`
	StartTime string `json:"startTime,omitempty"`
}

// CreateRes create collection response
type CreateRes struct {
	api.CommonRes
	AffectedCount int `json:"affectedCount,omitempty"`
}

// DescribeReq get collection detail request
type DescribeReq struct {
	api.Meta   `path:"/collection/describe" tags:"Collection" method:"Post" summary:"返回collection信息"`
	Database   string `json:"database,omitempty"`
	Collection string `json:"collection,omitempty"`
}

// DescribeRes get collection detail response
type DescribeRes struct {
	api.CommonRes
	Collection *DescribeCollectionItem `json:"collection"`
}

// DropReq delete collection request
type DropReq struct {
	api.Meta     `path:"/collection/drop" tags:"Collection" method:"Post" summary:"删除collection，并删除collection中的所有文档，如果collectio不经存在返回失败"`
	Database     string `json:"database,omitempty"`
	Collection   string `json:"collection,omitempty"`
	Force        bool   `json:"force,omitempty"`
	WithoutAlias bool   `json:"without_alias,omitempty"`
}

// DropReq delete collection response
type DropRes struct {
	api.CommonRes
	AffectedCount int `json:"affectedCount,omitempty"`
}

type ListReq struct {
	api.Meta `path:"/collection/list" tags:"Collection" method:"Post" summary:"列出指定database中的所有collection"`
	Database string `json:"database,omitempty"`
}

type ListRes struct {
	api.CommonRes
	Collections []*DescribeCollectionItem `json:"collections,omitempty"`
}

type DescribeCollectionItem struct {
	Database      string             `json:"database,omitempty"`
	Collection    string             `json:"collection,omitempty"`
	ReplicaNum    uint32             `json:"replicaNum,omitempty"`
	ShardNum      uint32             `json:"shardNum,omitempty"`
	Size          uint64             `json:"size,omitempty"`
	CreateTime    string             `json:"createTime,omitempty"`
	Description   string             `json:"description,omitempty"`
	Indexes       []*api.IndexColumn `json:"indexes,omitempty"`
	IndexStatus   *IndexStatus       `json:"indexStatus,omitempty"`
	Alias         []string           `json:"alias"`
	DocumentCount int64              `json:"documentCount,omitempty"`
	Embedding     *EmbeddingRes      `json:"embedding,omitempty"`
}

type TruncateReq struct {
	api.Meta          `path:"/collection/truncate" tags:"Collection" method:"Post" summary:"清空 collection 中的所有数据和索引"`
	Database          string `json:"database,omitempty"`
	Collection        string `json:"collection,omitempty"`
	OnlyFlushAnnIndex bool   `json:"only_flush_ann_index,omitempty"`
}

type TruncateRes struct {
	api.CommonRes
	AffectedCount int `json:"affectedCount,omitempty"`
}

type ModifyRes struct {
	api.CommonRes
	TaskIds []string `json:"task_ids,omitempty"`
}

type EmbeddingRes struct {
	Embedding
	Status string
}
