package model

import (
	"encoding/json"
	"time"
)

type IndexParams interface {
	MarshalJson() ([]byte, error)
	Name() string
}

var _ IndexParams = &HNSWParam{}
var _ IndexParams = &IVFFLATParams{}
var _ IndexParams = &IVFSQ8Params{}
var _ IndexParams = &IVFPQParams{}

type HNSWParam struct {
	M              uint32
	EfConstruction uint32
}

func (p *HNSWParam) MarshalJson() ([]byte, error) {
	return json.Marshal(p)
}

func (p *HNSWParam) Name() string {
	return string(HNSW)
}

type IVFFLATParams struct {
	NList uint32
}

func (p *IVFFLATParams) MarshalJson() ([]byte, error) {
	return json.Marshal(p)
}

func (p *IVFFLATParams) Name() string {
	return string(IVF_FLAT)
}

type IVFSQ8Params struct {
	NList uint32
}

func (p *IVFSQ8Params) MarshalJson() ([]byte, error) {
	return json.Marshal(p)
}

func (p *IVFSQ8Params) Name() string {
	return string(IVF_SQ8)
}

type IVFPQParams struct {
	M     uint32
	NList uint32
}

func (p *IVFPQParams) MarshalJson() ([]byte, error) {
	return json.Marshal(p)
}

func (p *IVFPQParams) Name() string {
	return string(IVF_PQ)
}

type FilterIndex struct {
	FieldName string
	FieldType FieldType
	IndexType IndexType
}

func (i *FilterIndex) IsPrimaryKey() bool {
	return i.IndexType == PRIMARY
}

func (i *FilterIndex) IsVectorField() bool {
	return i.FieldType == Vector
}

type VectorIndex struct {
	FilterIndex
	Dimension    uint32
	MetricType   MetricType
	IndexedCount uint64
	Params       IndexParams
}

type Indexes struct {
	VectorIndex []VectorIndex
	FilterIndex []FilterIndex
}

type IndexStatus struct {
	Status    string
	StartTime time.Time
}
