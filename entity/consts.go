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

package entity

type IndexType string

const (
	// vector index type
	FLAT     IndexType = "FLAT"
	HNSW     IndexType = "HNSW"
	IVF_FLAT IndexType = "IVF_FLAT"
	IVF_PQ   IndexType = "IVF_PQ"
	IVF_SQ4  IndexType = "IVF_SQ4"
	IVF_SQ8  IndexType = "IVF_SQ8"
	IVF_SQ16 IndexType = "IVF_SQ16"

	// scalar index type
	PRIMARY IndexType = "primaryKey"
	FILTER  IndexType = "filter"
)

type MetricType string

const (
	L2     MetricType = "L2"
	IP     MetricType = "IP"
	COSINE MetricType = "COSINE"
)

type FieldType string

const (
	Uint64 FieldType = "uint64"
	String FieldType = "string"
	Vector FieldType = "vector"
)

type EmbeddingModel string

const (
	// M3E_BASE 768
	M3E_BASE EmbeddingModel = "m3e-base"
	// BGE_BASE_ZH 768
	BGE_BASE_ZH EmbeddingModel = "bge-base-zh"
	// MULTILINGUAL 768
	MULTILINGUAL_E5_BASE EmbeddingModel = "multilingual-e5-base"
	// E5_LARGE_V2 1024
	E5_LARGE_V2 EmbeddingModel = "e5-large-v2"
	// TEXT2VEC_LARGE_CHINESE 1024
	TEXT2VEC_LARGE_CHINESE EmbeddingModel = "text2vec-large-chinese"
)

type ReadConsistency string

const (
	// EventualConsistency default value, 选择就近节点
	EventualConsistency ReadConsistency = "eventualConsistency"
	StrongConsistency   ReadConsistency = "strongConsistency"
)
