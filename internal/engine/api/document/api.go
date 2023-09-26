package document

import (
	"encoding/json"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/internal/engine/api"
)

// UpsertReq upsert document request
type UpsertReq struct {
	api.Meta   `path:"/document/upsert" tags:"Document" method:"Post" summary:"插入一条文档数据"`
	Database   string      `protobuf:"bytes,1,opt,name=database,proto3" json:"database,omitempty"`
	Collection string      `protobuf:"bytes,2,opt,name=collection,proto3" json:"collection,omitempty"`
	BuildIndex bool        `protobuf:"varint,3,opt,name=buildIndex,proto3" json:"buildIndex,omitempty"` // 是否立即构建索引
	Documents  []*Document `json:"documents,omitempty"`
}

// UpsertRes upsert document response
type UpsertRes struct {
	Code          int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg           string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Redirect      string `protobuf:"bytes,3,opt,name=redirect,proto3" json:"redirect,omitempty"`
	AffectedCount int32  `protobuf:"varint,4,opt,name=affectedCount,proto3" json:"affectedCount,omitempty"`
	Warning       string `protobuf:"bytes,5,opt,name=warning,proto3" json:"warning,omitempty"`
}

// Document document struct for document api
type Document struct {
	Id           string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Vector       []float32              `protobuf:"fixed32,2,rep,packed,name=vector,proto3" json:"vector,omitempty"`
	Score        float32                `protobuf:"fixed32,3,opt,name=score,proto3" json:"score,omitempty"`
	IndexId      uint64                 `protobuf:"varint,5,opt,name=index_id,json=indexId,proto3" json:"index_id,omitempty"`
	FromPeer     string                 `protobuf:"bytes,6,opt,name=from_peer,json=fromPeer,proto3" json:"from_peer,omitempty"`
	ShardIdx     int32                  `protobuf:"varint,7,opt,name=shard_idx,json=shardIdx,proto3" json:"shard_idx,omitempty"`
	VectorOffset uint64                 `protobuf:"varint,8,opt,name=vector_offset,json=vectorOffset,proto3" json:"vector_offset,omitempty"`
	DocInfo      []byte                 `protobuf:"bytes,9,opt,name=doc_info,json=docInfo,proto3" json:"doc_info,omitempty"`
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
	err = json.Unmarshal(data, &temp.Fields)
	if err != nil {
		return err
	}
	delete(temp.Fields, "vector")

	*d = Document(temp)
	return nil
}

// SearchReq search documents request
type SearchReq struct {
	api.Meta        `path:"/document/search" tags:"Document" method:"Post" summary:"向量查询接口，支持向量检索以及向量+标量混合检索"`
	Database        string      `protobuf:"bytes,1,opt,name=database,proto3" json:"database,omitempty"`
	Collection      string      `protobuf:"bytes,2,opt,name=collection,proto3" json:"collection,omitempty"`           // 索引名称
	ReadConsistency string      `protobuf:"bytes,4,opt,name=readConsistency,proto3" json:"readConsistency,omitempty"` // 读取一致性
	Search          *SearchCond `json:"search,omitempty"`
}

// SearchRes search documents response
type SearchRes struct {
	Code      int32         `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg       string        `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Redirect  string        `protobuf:"bytes,3,opt,name=redirect,proto3" json:"redirect,omitempty"`
	Warning   string        `protobuf:"bytes,5,opt,name=warning,proto3" json:"warning,omitempty"`
	Documents [][]*Document `json:"documents,omitempty"`
}

// SearchCond search filter condition
type SearchCond struct {
	DocumentIds    []string      `protobuf:"bytes,2,rep,name=documentIds,proto3" json:"documentIds,omitempty"` // 使用向量id检索
	Params         *SearchParams `protobuf:"bytes,3,opt,name=params,proto3" json:"params,omitempty"`
	RetrieveVector bool          `protobuf:"varint,5,opt,name=retrieveVector,proto3" json:"retrieveVector,omitempty"` // 是否返回原始向量，注意设置为true时会降低性能
	Limit          int64         `protobuf:"varint,6,opt,name=limit,proto3" json:"limit,omitempty"`                   // 结果数量
	OutputFields   []string      `protobuf:"bytes,7,rep,name=outputfields,proto3" json:"outputfields,omitempty"`      // 输出字段
	Retrieves      []string      `protobuf:"bytes,8,rep,name=retrieves,proto3" json:"retrieves,omitempty"`            // 使用字符串检索
	Vectors        [][]float32   `json:"vectors,omitempty"`
	Filter         string        `json:"filter,omitempty"`
	EmbeddingItems []string      `json:"embeddingItems,omitempty"`
}

type SearchParams struct {
	Nprobe uint32  `protobuf:"varint,1,opt,name=nprobe,proto3" json:"nprobe,omitempty"`  // 搜索时查找的聚类数量，使用索引默认值即可
	Ef     uint32  `protobuf:"varint,2,opt,name=ef,proto3" json:"ef,omitempty"`          // HNSW
	Radius float32 `protobuf:"fixed32,3,opt,name=radius,proto3" json:"radius,omitempty"` // 距离阈值,范围搜索时有效
}

// QueryReq query document request
type QueryReq struct {
	api.Meta        `path:"/document/query" tags:"Document" method:"Post" summary:"标量查询接口，当前仅支持主键id查询"`
	Database        string     `protobuf:"bytes,1,opt,name=database,proto3" json:"database,omitempty"`
	Collection      string     `protobuf:"bytes,2,opt,name=collection,proto3" json:"collection,omitempty"`
	Query           *QueryCond `protobuf:"bytes,3,opt,name=query,proto3" json:"query,omitempty"`
	ReadConsistency string     `protobuf:"bytes,4,opt,name=readConsistency,proto3" json:"readConsistency,omitempty"`
}

type QueryCond struct {
	DocumentIds    []string `protobuf:"bytes,1,rep,name=documentIds,proto3" json:"documentIds,omitempty"`
	IndexIds       []uint64 `protobuf:"varint,2,rep,packed,name=indexIds,proto3" json:"indexIds,omitempty"`
	RetrieveVector bool     `protobuf:"varint,3,opt,name=retrieveVector,proto3" json:"retrieveVector,omitempty"`
	Filter         string   `protobuf:"bytes,4,opt,name=filter,proto3" json:"filter,omitempty"`
	Limit          int64    `protobuf:"varint,5,opt,name=limit,proto3" json:"limit,omitempty"`
	Offset         int64    `protobuf:"varint,6,opt,name=offset,proto3" json:"offset,omitempty"`
	OutputFields   []string `protobuf:"bytes,7,rep,name=outputFields,proto3" json:"outputFields,omitempty"`
}

// QueryRes query document response
type QueryRes struct {
	Code      int32       `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg       string      `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Redirect  string      `protobuf:"bytes,3,opt,name=redirect,proto3" json:"redirect,omitempty"`
	Count     uint64      `protobuf:"varint,5,opt,name=count,proto3" json:"count,omitempty"`
	Documents []*Document `json:"documents,omitempty"`
}

// DeleteReq delete document request
type DeleteReq struct {
	api.Meta   `path:"/document/delete" tags:"Document" method:"Post" summary:"删除指定id的文档,flat 索引不支持删除"`
	Database   string     `protobuf:"bytes,1,opt,name=database,proto3" json:"database,omitempty"`
	Collection string     `protobuf:"bytes,2,opt,name=collection,proto3" json:"collection,omitempty"`
	Query      *QueryCond `protobuf:"bytes,3,opt,name=query,proto3" json:"query,omitempty"`
}

// DeleteRes delete document request
type DeleteRes struct {
	Code          int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg           string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Redirect      string `protobuf:"bytes,3,opt,name=redirect,proto3" json:"redirect,omitempty"`
	AffectedCount uint64 `protobuf:"varint,4,opt,name=affectedCount,proto3" json:"affectedCount,omitempty"`
}

type UpdateReq struct {
	api.Meta   `path:"/document/update" tags:"Document" method:"Post" summary:"基于[主键查询]和[ Filter 过滤]的部分字段更新或者新增非索引字段"`
	Database   string     `protobuf:"bytes,1,opt,name=database,proto3" json:"database,omitempty"`
	Collection string     `protobuf:"bytes,2,opt,name=collection,proto3" json:"collection,omitempty"`
	Query      *QueryCond `protobuf:"bytes,3,opt,name=query,proto3" json:"query,omitempty"`
	Update     Document   `json:"update,omitempty"`
}

type UpdateRes struct {
	Code          int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg           string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Redirect      string `protobuf:"bytes,3,opt,name=redirect,proto3" json:"redirect,omitempty"`
	AffectedCount uint64 `protobuf:"varint,4,opt,name=affectedCount,proto3" json:"affectedCount,omitempty"`
	Warning       string `protobuf:"bytes,5,opt,name=warning,proto3" json:"warning,omitempty"`
}
