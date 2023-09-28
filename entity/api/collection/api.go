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

import "git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity/api"

// CreateReq create collection request
type CreateReq struct {
	api.Meta    `path:"/collection/create" tags:"Collection" method:"Post" summary:"创建collection"`
	Database    string         `protobuf:"bytes,1,opt,name=database,proto3" json:"database,omitempty"`
	Collection  string         `protobuf:"bytes,2,opt,name=collection,proto3" json:"collection,omitempty"`
	ReplicaNum  uint32         `protobuf:"varint,3,opt,name=replicaNum,proto3" json:"replicaNum,omitempty"`
	ShardNum    uint32         `protobuf:"varint,4,opt,name=shardNum,proto3" json:"shardNum,omitempty"`
	Size        uint64         `protobuf:"varint,5,opt,name=size,proto3" json:"size,omitempty"`
	CreateTime  string         `protobuf:"bytes,6,opt,name=createTime,proto3" json:"createTime,omitempty"`
	Description string         `protobuf:"bytes,7,opt,name=description,proto3" json:"description,omitempty"`
	Indexes     []*IndexColumn `protobuf:"bytes,8,rep,name=indexes,proto3" json:"indexes,omitempty"`
	IndexStatus *IndexStatus   `protobuf:"bytes,9,opt,name=indexStatus,proto3" json:"indexStatus,omitempty"`
	AliasList   []string       `protobuf:"bytes,10,rep,name=alias_list,json=aliasList,proto3" json:"alias_list,omitempty"`
	Embedding   Embedding      `json:"embedding"`
}

type Embedding struct {
	Field       string `json:"field,omitempty"`
	VectorField string `json:"vectorField,omitempty"`
	Model       string `json:"model,omitempty"`
}

type IndexColumn struct {
	FieldName    string       `protobuf:"bytes,1,opt,name=fieldName,proto3" json:"fieldName,omitempty"`
	FieldType    string       `protobuf:"bytes,2,opt,name=fieldType,proto3" json:"fieldType,omitempty"`
	IndexType    string       `protobuf:"bytes,3,opt,name=indexType,proto3" json:"indexType,omitempty"`
	Dimension    uint32       `protobuf:"varint,4,opt,name=dimension,proto3" json:"dimension,omitempty"`
	MetricType   string       `protobuf:"bytes,5,opt,name=metricType,proto3" json:"metricType,omitempty"`
	IndexedCount uint64       `protobuf:"varint,6,opt,name=indexedCount,proto3" json:"indexedCount,omitempty"`
	Params       *IndexParams `protobuf:"bytes,8,opt,name=params,proto3" json:"params,omitempty"`
}

type IndexParams struct {
	M              uint32 `protobuf:"varint,1,opt,name=M,proto3" json:"M,omitempty"`
	EfConstruction uint32 `protobuf:"varint,2,opt,name=efConstruction,proto3" json:"efConstruction,omitempty"`
	Nprobe         uint32 `protobuf:"varint,3,opt,name=nprobe,proto3" json:"nprobe,omitempty"`
	Nlist          uint32 `protobuf:"varint,4,opt,name=nlist,proto3" json:"nlist,omitempty"`
}

type IndexStatus struct {
	Status    string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	Progress  string `protobuf:"bytes,2,opt,name=progress,proto3" json:"progress,omitempty"`
	StartTime string `protobuf:"bytes,3,opt,name=startTime,proto3" json:"startTime,omitempty"`
}

// CreateRes create collection response
type CreateRes struct {
	Code          int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg           string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Redirect      string `protobuf:"bytes,3,opt,name=redirect,proto3" json:"redirect,omitempty"`
	AffectedCount int32  `protobuf:"varint,4,opt,name=affectedCount,proto3" json:"affectedCount,omitempty"`
}

// DescribeReq get collection detail request
type DescribeReq struct {
	api.Meta   `path:"/collection/describe" tags:"Collection" method:"Post" summary:"返回collection信息"`
	Database   string `protobuf:"bytes,1,opt,name=database,proto3" json:"database,omitempty"`
	Collection string `protobuf:"bytes,2,opt,name=collection,proto3" json:"collection,omitempty"`
	Transfer   bool   `protobuf:"varint,3,opt,name=transfer,proto3" json:"transfer,omitempty"`
}

// DescribeRes get collection detail response
type DescribeRes struct {
	Code       int32                   `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg        string                  `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Redirect   string                  `protobuf:"bytes,3,opt,name=redirect,proto3" json:"redirect,omitempty"`
	Collection *DescribeCollectionItem `json:"collection"`
}

// DropReq delete collection request
type DropReq struct {
	api.Meta     `path:"/collection/drop" tags:"Collection" method:"Post" summary:"删除collection，并删除collection中的所有文档，如果collectio不经存在返回失败"`
	Database     string `protobuf:"bytes,1,opt,name=database,proto3" json:"database,omitempty"`
	Collection   string `protobuf:"bytes,2,opt,name=collection,proto3" json:"collection,omitempty"`
	Force        bool   `protobuf:"varint,3,opt,name=force,proto3" json:"force,omitempty"`
	WithoutAlias bool   `protobuf:"varint,4,opt,name=without_alias,json=withoutAlias,proto3" json:"without_alias,omitempty"`
}

// DropReq delete collection response
type DropRes struct {
	Code          int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg           string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Redirect      string `protobuf:"bytes,3,opt,name=redirect,proto3" json:"redirect,omitempty"`
	AffectedCount int32  `protobuf:"varint,4,opt,name=affectedCount,proto3" json:"affectedCount,omitempty"`
}

type ListReq struct {
	api.Meta `path:"/collection/list" tags:"Collection" method:"Post" summary:"列出指定database中的所有collection"`
	Database string `protobuf:"bytes,1,opt,name=database,proto3" json:"database,omitempty"`
	Transfer bool   `protobuf:"varint,2,opt,name=transfer,proto3" json:"transfer,omitempty"`
}

type ListRes struct {
	Code        int32                     `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg         string                    `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Redirect    string                    `protobuf:"bytes,3,opt,name=redirect,proto3" json:"redirect,omitempty"`
	Collections []*DescribeCollectionItem `json:"collections,omitempty"`
}

type DescribeCollectionItem struct {
	Database      string         `protobuf:"bytes,1,opt,name=database,proto3" json:"database,omitempty"`
	Collection    string         `protobuf:"bytes,2,opt,name=collection,proto3" json:"collection,omitempty"`
	ReplicaNum    uint32         `protobuf:"varint,3,opt,name=replicaNum,proto3" json:"replicaNum,omitempty"`
	ShardNum      uint32         `protobuf:"varint,4,opt,name=shardNum,proto3" json:"shardNum,omitempty"`
	Size          uint64         `protobuf:"varint,5,opt,name=size,proto3" json:"size,omitempty"`
	CreateTime    string         `protobuf:"bytes,6,opt,name=createTime,proto3" json:"createTime,omitempty"`
	Description   string         `protobuf:"bytes,7,opt,name=description,proto3" json:"description,omitempty"`
	Indexes       []*IndexColumn `protobuf:"bytes,8,rep,name=indexes,proto3" json:"indexes,omitempty"`
	IndexStatus   *IndexStatus   `protobuf:"bytes,9,opt,name=indexStatus,proto3" json:"indexStatus,omitempty"`
	Alias         []string       `json:"alias"`
	DocumentCount int64          `json:"documentCount,omitempty"`
	Embedding     EmbeddingRes   `json:"embedding"`
}

type TruncateReq struct {
	api.Meta          `path:"/collection/truncate" tags:"Collection" method:"Post" summary:"清空 collection 中的所有数据和索引"`
	Database          string `protobuf:"bytes,1,opt,name=database,proto3" json:"database,omitempty"`
	Collection        string `protobuf:"bytes,2,opt,name=collection,proto3" json:"collection,omitempty"`
	OnlyFlushAnnIndex bool   `protobuf:"varint,3,opt,name=only_flush_ann_index,json=onlyFlushAnnIndex,proto3" json:"only_flush_ann_index,omitempty"`
}

type TruncateRes struct {
	Code          int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg           string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Redirect      string `protobuf:"bytes,3,opt,name=redirect,proto3" json:"redirect,omitempty"`
	AffectedCount int32  `protobuf:"varint,4,opt,name=affectedCount,proto3" json:"affectedCount,omitempty"`
}

type ModifyRes struct {
	Code     int32    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg      string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Redirect string   `protobuf:"bytes,3,opt,name=redirect,proto3" json:"redirect,omitempty"`
	TaskIds  []string `protobuf:"bytes,4,rep,name=task_ids,json=taskIds,proto3" json:"task_ids,omitempty"`
}

type EmbeddingRes struct {
	Embedding
	Status string
}
