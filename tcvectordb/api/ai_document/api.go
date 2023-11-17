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

package ai_document

import "git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb/api"

// QueryReq query document request
type QueryReq struct {
	api.Meta   `path:"/ai/document/query" tags:"Document" method:"Post"`
	Database   string     `json:"database"`
	Collection string     `json:"collection"`
	Query      *QueryCond `json:"query"`
}

type QueryCond struct {
	DocumentIds  []string `json:"documentIds,omitempty"`
	Filter       string   `json:"filter,omitempty"`
	Limit        int64    `json:"limit,omitempty"`
	Offset       int64    `json:"offset,omitempty"`
	OutputFields []string `json:"outputFields,omitempty"`
}

// QueryRes query document response
type QueryRes struct {
	api.CommonRes
	Count     uint64          `json:"count"`
	Documents []QueryDocument `json:"documents"`
}

// SearchReq search documents request
type SearchReq struct {
	api.Meta        `path:"/ai/document/search" tags:"Document" method:"Post"`
	Database        string      `protobuf:"bytes,1,opt,name=database,proto3" json:"database"`
	Collection      string      `protobuf:"bytes,2,opt,name=collection,proto3" json:"collection"` // 索引名称
	ReadConsistency string      `json:"readConsistency"`
	Search          *SearchCond `json:"search"`
}

// SearchRes search documents response
type SearchRes struct {
	api.CommonRes
	Documents []SearchDocument `json:"documents"`
}

// SearchCond search filter condition
type SearchCond struct {
	Content      string       `json:"content"`
	Filter       string       `json:"filter"`
	Options      SearchOption `json:"options"`
	OutputFields []string     `json:"outputfields"`    // 输出字段
	Limit        int64        `json:"limit,omitempty"` // 结果数量
}

type SearchOption struct {
	ResultType  string `json:"resultType"`  // chunks|paragraphs|file
	ChunkExpand []int  `json:"chunkExpand"` // 搜索结果中，向前、向后补齐几个chunk的上下文
	// MergeChunk  bool   `json:"mergeChunk"`  // Merge结果中相邻的Chunk
	// Weights     SearchOptionWeight `json:"weights"`     // 多路召回
}

type SearchOptionWeight struct {
	ChunkSimilarity float64 `json:"chunkSimilarity"`
	WordSimilarity  float64 `json:"wordSimilarity"`
	WordBm25        float64 `json:"wordBm25"`
}

type SearchParams struct {
	Nprobe uint32  `protobuf:"varint,1,opt,name=nprobe,proto3" json:"nprobe"`  // 搜索时查找的聚类数量，使用索引默认值即可
	Ef     uint32  `protobuf:"varint,2,opt,name=ef,proto3" json:"ef"`          // HNSW
	Radius float32 `protobuf:"fixed32,3,opt,name=radius,proto3" json:"radius"` // 距离阈值,范围搜索时有效
}

// DeleteReq delete document request
type DeleteReq struct {
	api.Meta   `path:"/ai/document/delete" tags:"Document" method:"Post"`
	Database   string           `json:"database"`
	Collection string           `json:"collection"`
	Query      *DeleteQueryCond `json:"query"`
}

type DeleteQueryCond struct {
	DocumentIds []string `json:"documentIds"`
	Filter      string   `json:"filter"`
}

// DeleteRes delete document request
type DeleteRes struct {
	api.CommonRes
	AffectedCount uint64 `json:"affectedCount"`
}

type UpdateReq struct {
	api.Meta   `path:"/ai/document/update" tags:"Document" method:"Post" summary:"基于[主键查询]和[ Filter 过滤]的部分字段更新或者新增非索引字段"`
	Database   string                 `json:"database"`
	Collection string                 `json:"collection"`
	Query      UpdateQueryCond        `json:"query"`
	Update     map[string]interface{} `json:"update"`
}

type UpdateQueryCond struct {
	DocumentIds []string `json:"documentIds"`
	Filter      string   `json:"filter"`
}

type UpdateRes struct {
	api.CommonRes
	AffectedCount uint64 `json:"affectedCount"`
}

type UploadUrlReq struct {
	api.Meta   `path:"/ai/document/uploadurl" tags:"Document" method:"Post" summary:"获取cos上传签名"`
	Database   string `json:"database"`
	Collection string `json:"collection"`
	FileName   string `json:"fileName"`
	FileType   string `json:"fileType"`
}

type UploadUrlRes struct {
	api.CommonRes
	CosEndpoint     string           `json:"cosEndpoint"`
	CosRegion       string           `json:"cosRegion,omitempty"`
	CosBucket       string           `json:"cosBucket,omitempty"`
	UploadPath      string           `json:"uploadPath"`
	Credentials     *Credentials     `json:"credentials"`
	UploadCondition *UploadCondition `json:"uploadCondition"`
	FileId          string           `json:"fileId"`
}

type UploadCondition struct {
	MaxSupportContentLength int64 `json:"maxSupportContentLength"`
}

type Credentials struct {
	TmpSecretID  string `json:"TmpSecretId"`
	TmpSecretKey string `json:"TmpSecretKey"`
	SessionToken string `json:"Token"`
}
