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

package tcvectordb

import (
	"encoding/json"
)

type Indexes struct {
	VectorIndex       []VectorIndex
	FilterIndex       []FilterIndex
	SparseVectorIndex []SparseVectorIndex
}

type SparseVectorIndex struct {
	FieldName       string
	FieldType       FieldType
	IndexType       IndexType
	MetricType      MetricType
	DiskSwapEnabled *bool
}

type FilterIndex struct {
	FieldName string
	FieldType FieldType
	ElemType  FieldType
	IndexType IndexType
	AutoId    string
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

type IndexParams interface {
	MarshalJson() ([]byte, error)
	Name() string
}

var _ IndexParams = &HNSWParam{}
var _ IndexParams = &IVFFLATParams{}
var _ IndexParams = &IVFSQParams{}
var _ IndexParams = &IVFPQParams{}
var _ IndexParams = &IVFRabitQParams{}

type IVFRabitQParams struct {
	NList uint32
	Bits  *uint32
}

func (p *IVFRabitQParams) MarshalJson() ([]byte, error) {
	return json.Marshal(p)
}

func (p *IVFRabitQParams) Name() string {
	return string(IVF_RABITQ)
}

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

type IVFSQParams struct {
	NList uint32
}

func (p *IVFSQParams) MarshalJson() ([]byte, error) {
	return json.Marshal(p)
}

func (p *IVFSQParams) Name() string {
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
