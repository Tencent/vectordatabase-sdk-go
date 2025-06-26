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

package index

import "github.com/tencent/vectordatabase-sdk-go/tcvectordb/api"

type RebuildReq struct {
	api.Meta          `path:"/index/rebuild" tags:"Index" method:"Post" summary:"重建整个collection的所有索引"`
	Database          string `json:"database,omitempty"`
	Collection        string `json:"collection,omitempty"`
	DropBeforeRebuild bool   `json:"dropBeforeRebuild,omitempty"`
	Throttle          int32  `json:"throttle"`
	DisableTrain      bool   `json:"disable_train,omitempty"`
	ForceRebuild      bool   `json:"force_rebuild,omitempty"`
	FieldName         string `json:"fieldName,omitempty"`
}

type RebuildRes struct {
	api.CommonRes
	TaskIds []string `json:"task_ids,omitempty"`
}

type AddReq struct {
	api.Meta         `path:"/index/add" tags:"Index" method:"Post" summary:"新增collection的索引"`
	Database         string             `json:"database,omitempty"`
	Collection       string             `json:"collection,omitempty"`
	Indexes          []*api.IndexColumn `json:"indexes,omitempty"`
	BuildExistedData *bool              `json:"buildExistedData,omitempty"`
}

type AddRes struct {
	api.CommonRes
}

type DropReq struct {
	api.Meta   `path:"/index/drop" tags:"Index" method:"Post" summary:"删除collection的索引"`
	Database   string   `json:"database,omitempty"`
	Collection string   `json:"collection,omitempty"`
	FieldNames []string `json:"fieldNames,omitempty"`
}

type DropRes struct {
	api.CommonRes
}

type ModifyVectorIndexReq struct {
	api.Meta      `path:"/index/modifyVectorIndex" tags:"Index" method:"Post" summary:"调整collection的向量索引参数"`
	Database      string             `json:"database,omitempty"`
	Collection    string             `json:"collection,omitempty"`
	VectorIndexes []*api.IndexColumn `json:"vectorIndexes,omitempty"`
	RebuildRules  *RebuildRules      `json:"rebuildRules,omitempty"`
}

type RebuildRules struct {
	DropBeforeRebuild *bool  `json:"dropBeforeRebuild,omitempty"`
	Throttle          *int32 `json:"throttle,omitempty"`
}

type ModifyVectorIndexRes struct {
	api.CommonRes
}
