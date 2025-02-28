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

package document

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/ai_document_set"
)

// UpsertReq upsert document request
type UpsertReq struct {
	api.Meta   `path:"/document/upsert" tags:"Document" method:"Post" summary:"插入一条文档数据"`
	Database   string      `json:"database,omitempty"`
	Collection string      `json:"collection,omitempty"`
	BuildIndex *bool       `json:"buildIndex,omitempty"` // 是否立即构建索引
	Documents  []*Document `json:"documents,omitempty"`
}

// UpsertRes upsert document response
type UpsertRes struct {
	api.CommonRes
	AffectedCount int    `json:"affectedCount,omitempty"`
	Warning       string `json:"warning,omitempty"`
}

// Document document struct for document api
type Document struct {
	Id           string                 `json:"id,omitempty"`
	Vector       []float32              `json:"vector,omitempty"`
	SparseVector [][]interface{}        `json:"sparse_vector,omitempty"`
	Score        float32                `json:"score,omitempty"`
	DocInfo      []byte                 `json:"doc_info,omitempty"`
	Fields       map[string]interface{} `json:"-"`
}

func (d Document) MarshalJSON() ([]byte, error) {
	type Alias Document
	res, err := json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&d),
	})
	if err != nil {
		return nil, err
	}
	if len(d.Fields) != 0 {
		field, err := json.Marshal(d.Fields)
		if err != nil {
			return nil, err
		}
		if len(field) == 0 {
			return res, nil
		}
		// res = {}
		if len(res) == 2 {
			res = append(res[:1], field[1:]...)
		} else {
			res[len(res)-1] = ','
			res = append(res, field[1:]...)
		}
	}
	return res, nil
}

func (d *Document) UnmarshalJSON(data []byte) error {
	type Alias Document
	var temp Alias
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}
	ds := json.NewDecoder(bytes.NewReader(data))
	ds.UseNumber()
	err = ds.Decode(&temp.Fields)
	if err != nil {
		return err
	}
	reflectType := reflect.TypeOf(*d)
	for i := 0; i < reflectType.NumField(); i++ {
		field := reflectType.Field(i)
		tags := strings.Split(field.Tag.Get("json"), ",")
		if tags[0] == "-" {
			continue
		}
		delete(temp.Fields, tags[0])
	}

	*d = Document(temp)
	return nil
}

// SearchReq search documents request
type SearchReq struct {
	api.Meta        `path:"/document/search" tags:"Document" method:"Post" summary:"向量查询接口，支持向量检索以及向量+标量混合检索"`
	Database        string      `json:"database,omitempty"`
	Collection      string      `json:"collection,omitempty"`      // 索引名称
	ReadConsistency string      `json:"readConsistency,omitempty"` // 读取一致性
	Search          *SearchCond `json:"search,omitempty"`
}

// SearchRes search documents response
type SearchRes struct {
	api.CommonRes
	Warning   string        `json:"warning,omitempty"`
	Documents [][]*Document `json:"documents,omitempty"`
}

type HybridSearchReq struct {
	api.Meta        `path:"/document/hybridSearch" tags:"Document" method:"Post" summary:"向量查询接口，支持混合检索"`
	Database        string            `json:"database,omitempty"`
	Collection      string            `json:"collection,omitempty"`      // 索引名称
	ReadConsistency string            `json:"readConsistency,omitempty"` // 读取一致性
	Search          *HybridSearchCond `json:"search,omitempty"`
}

type HybridSearchCond struct {
	RetrieveVector bool     `json:"retrieveVector,omitempty"` // 是否返回原始向量，注意设置为true时会降低性能
	Limit          *int     `json:"limit,omitempty"`          // 结果数量
	OutputFields   []string `json:"outputFields,omitempty"`   // 输出字段

	Filter string `json:"filter,omitempty"`

	AnnParams []*AnnParam    `json:"ann,omitempty"`
	Rerank    *RerankOption  `json:"rerank,omitempty"`
	Match     []*MatchOption `json:"match,omitempty"`
}

type RerankOption struct {
	Method    string    `json:"method,omitempty"`
	FieldList []string  `json:"fieldList,omitempty"`
	Weight    []float32 `json:"weight,omitempty"`
	RrfK      int32     `protobuf:"varint,3,opt,name=rrf_k,json=rrfK,proto3" json:"rrf_k,omitempty"` // for RRF: K参数
}
type MatchOption struct {
	FieldName       string            `protobuf:"bytes,1,opt,name=fieldName,proto3" json:"fieldName,omitempty"`
	Data            [][][]interface{} `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Limit           int               `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
	TerminateAfter  uint32            `json:"terminateAfter,omitempty"`
	CutoffFrequency float64           `json:"cutoffFrequency,omitempty"`
}

type AnnParam struct {
	FieldName   string        `json:"fieldName,omitempty"`
	DocumentIds []string      `json:"documentIds,omitempty"`
	Data        []interface{} `json:"data,omitempty"`
	Params      *SearchParams `json:"params,omitempty"`
	Limit       *int          `json:"limit,omitempty"`
}

// SearchCond search filter condition
type SearchCond struct {
	DocumentIds    []string      `json:"documentIds,omitempty"` // 使用向量id检索
	Params         *SearchParams `json:"params,omitempty"`
	RetrieveVector bool          `json:"retrieveVector,omitempty"` // 是否返回原始向量，注意设置为true时会降低性能
	Limit          int64         `json:"limit,omitempty"`          // 结果数量
	OutputFields   []string      `json:"outputFields,omitempty"`   // 输出字段
	Retrieves      []string      `json:"retrieves,omitempty"`      // 使用字符串检索
	Vectors        [][]float32   `json:"vectors,omitempty"`
	Filter         string        `json:"filter,omitempty"`
	EmbeddingItems []string      `json:"embeddingItems,omitempty"`
	Radius         *float32      `json:"radius,omitempty"` // 距离阈值,范围搜索时有效
}

type SearchParams struct {
	Nprobe uint32  `json:"nprobe,omitempty"` // 搜索时查找的聚类数量，使用索引默认值即可
	Ef     uint32  `json:"ef,omitempty"`     // HNSW
	Radius float32 `json:"radius,omitempty"` // 距离阈值,范围搜索时有效
}

// QueryReq query document request
type QueryReq struct {
	api.Meta        `path:"/document/query" tags:"Document" method:"Post" summary:"标量查询接口，当前仅支持主键id查询"`
	Database        string     `json:"database,omitempty"`
	Collection      string     `json:"collection,omitempty"`
	Query           *QueryCond `json:"query,omitempty"`
	ReadConsistency string     `json:"readConsistency,omitempty"`
}

type QueryCond struct {
	DocumentIds    []string   `json:"documentIds,omitempty"`
	IndexIds       []uint64   `json:"indexIds,omitempty"`
	RetrieveVector bool       `json:"retrieveVector,omitempty"`
	Filter         string     `json:"filter,omitempty"`
	Limit          int64      `json:"limit,omitempty"`
	Offset         int64      `json:"offset,omitempty"`
	OutputFields   []string   `json:"outputFields,omitempty"`
	Sort           []SortRule `json:"sort,omitempty"`
}

// [SortRule] holds the fields for a single sort rule.
//
// Fields:
//   - FieldName: (Required) The field name for the sort rule, and you can set the field name for the uint64 filter index.
//   - Direction: (Optional) The sort direction, where you can set "desc" or "asc" (default to "asc").
type SortRule struct {
	FieldName string `json:"fieldName,omitempty"`
	Direction string `json:"direction,omitempty"`
}

// QueryRes query document response
type QueryRes struct {
	api.CommonRes
	Count     uint64      `json:"count,omitempty"`
	Documents []*Document `json:"documents,omitempty"`
}

// CountReq query document request
type CountReq struct {
	api.Meta        `path:"/document/count" tags:"Document" method:"Post" summary:"基于Filter统计文档数量"`
	Database        string          `json:"database,omitempty"`
	Collection      string          `json:"collection,omitempty"`
	Query           *CountQueryCond `json:"query,omitempty"`
	ReadConsistency string          `json:"readConsistency,omitempty"`
}

type CountQueryCond struct {
	Filter string `json:"filter,omitempty"`
}

// QueryRes query document response
type CountRes struct {
	api.CommonRes
	Count uint64 `json:"count,omitempty"`
}

// DeleteReq delete document request
type DeleteReq struct {
	api.Meta   `path:"/document/delete" tags:"Document" method:"Post" summary:"删除指定id的文档,flat 索引不支持删除"`
	Database   string     `json:"database,omitempty"`
	Collection string     `json:"collection,omitempty"`
	Query      *QueryCond `json:"query,omitempty"`
}

// DeleteRes delete document request
type DeleteRes struct {
	api.CommonRes
	AffectedCount int `json:"affectedCount,omitempty"`
}

type UpdateReq struct {
	api.Meta   `path:"/document/update" tags:"Document" method:"Post" summary:"基于[主键查询]和[ Filter 过滤]的部分字段更新或者新增非索引字段"`
	Database   string     `json:"database,omitempty"`
	Collection string     `json:"collection,omitempty"`
	Query      *QueryCond `json:"query,omitempty"`
	Update     Document   `json:"update,omitempty"`
}

type UpdateRes struct {
	api.CommonRes
	AffectedCount int    `json:"affectedCount,omitempty"`
	Warning       string `json:"warning,omitempty"`
}

type UploadUrlReq struct {
	api.Meta           `path:"/ai/document/uploadUrl" tags:"Document" method:"Post" summary:"collection表上传文件"`
	Database           string                                      `json:"database"`
	Collection         string                                      `json:"collection"`
	FileName           string                                      `json:"fileName"`
	SplitterPreprocess *ai_document_set.DocumentSplitterPreprocess `json:"splitterPreprocess,omitempty"`
	EmbeddingModel     string                                      `json:"embeddingModel,omitempty"`
	ParsingProcess     *api.ParsingProcess                         `json:"parsingProcess,omitempty"`
	FieldMappings      map[string]string                           `json:"fieldMappings,omitempty"`
}

type UploadUrlRes struct {
	api.CommonRes
	Warning         string                           `json:"warning,omitempty"`
	CosEndpoint     string                           `json:"cosEndpoint,omitempty"`
	CosRegion       string                           `json:"cosRegion,omitempty"`
	CosBucket       string                           `json:"cosBucket,omitempty"`
	UploadPath      string                           `json:"uploadPath,omitempty"`
	Credentials     *ai_document_set.Credentials     `json:"credentials,omitempty"`
	UploadCondition *ai_document_set.UploadCondition `json:"uploadCondition,omitempty"`
}

type GetImageUrlReq struct {
	api.Meta    `path:"/ai/document/getImageUrl" tags:"Document" method:"Post" summary:"获取图片访问地址"`
	Database    string   `json:"database"`
	Collection  string   `json:"collection"`
	FileName    string   `json:"fileName"`
	DocumentIds []string `json:"documentIds"`
}

type GetImageUrlRes struct {
	api.CommonRes
	Images [][]ImageInfo `json:"images"`
}

type ImageInfo struct {
	DocumentId string `json:"documentId"`
	ImageName  string `json:"imageName"`
	ImageUrl   string `json:"imageUrl"`
}
