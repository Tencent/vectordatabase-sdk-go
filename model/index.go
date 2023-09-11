package model

import "time"

type HNSWParam struct {
	M              uint32
	EfConstruction uint32
}

type IVFFLATParams struct {
	NList uint32
}

type IVFSQ8Params struct {
	NList uint32
}

type IVFPQParams struct {
	M     uint32
	NList uint32
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
	Dimension     uint32
	MetricType    MetricType
	IndexedCount  uint64
	HNSWParam     HNSWParam
	IVFFLATParams IVFFLATParams
	IVFSQ8Params  IVFSQ8Params
	IVFPQParams   IVFPQParams
}

type Indexes struct {
	VectorIndex []VectorIndex
	FilterIndex []FilterIndex
}

type IndexStatus struct {
	Status    string
	StartTime time.Time
}
