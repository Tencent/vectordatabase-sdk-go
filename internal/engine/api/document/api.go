package document

import (
	"encoding/json"

	"git.woa.com/cloud_nosql/vectordb/vectordb-sdk-go/internal/proto"

	"github.com/gogf/gf/v2/frame/g"
)

type UpsertReq struct {
	g.Meta `path:"/document/upsert" tags:"Document" method:"Post" summary:"插入一条文档数据"`
	proto.UpsertRequest
	Documents []*Document `json:"documents,omitempty"`
}

type UpsertRes struct {
	proto.UpsertResponse
}

type Document struct {
	proto.Document
	Fields map[string]interface{} `json:"-"`
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
		res[len(res)-1] = ','
		res = append(res, field[1:]...)
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
	err = json.Unmarshal(data, &temp.Fields)
	if err != nil {
		return err
	}
	delete(temp.Fields, "vector")

	d.Document = temp.Document
	d.Fields = temp.Fields
	return nil
}

type SearchReq struct {
	g.Meta `path:"/document/search" tags:"Document" method:"Post" summary:"向量查询接口，支持向量检索以及向量+标量混合检索"`
	proto.SearchRequest
	Search *SearchCond `json:"search,omitempty"`
}

type SearchRes struct {
	proto.SearchResponse
	Documents [][]*Document `json:"documents,omitempty"`
}

type SearchCond struct {
	proto.SearchCond
	Vectors [][]float32 `json:"vectors,omitempty"`
	Filter  string      `json:"filter,omitempty"`
}

type QueryReq struct {
	g.Meta `path:"/document/query" tags:"Document" method:"Post" summary:"标量查询接口，当前仅支持主键id查询"`
	proto.QueryRequest
}

type QueryRes struct {
	proto.QueryResponse
	Documents []*Document `json:"documents,omitempty"`
}

type DeleteReq struct {
	g.Meta `path:"/document/delete" tags:"Document" method:"Post" summary:"删除指定id的文档,flat 索引不支持删除"`
	proto.DeleteRequest
}

type DeleteRes struct {
	proto.DeleteResponse
}
