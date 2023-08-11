package model

type HNSWParam struct {
	M              uint32
	EfConstruction uint32
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
	HNSWParam    HNSWParam
}

type Indexes struct {
	VectorIndex []VectorIndex
	FilterIndex []FilterIndex
}
