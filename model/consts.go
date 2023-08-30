package model

type IndexType string

const (
	// vector index type
	FLAT IndexType = "FLAT"
	HNSW IndexType = "HNSW"

	// scalar index type
	PRIMARY IndexType = "primaryKey"
	FILTER  IndexType = "filter"
)

type MetricType string

const (
	L2     MetricType = "L2"
	IP     MetricType = "IP"
	CONINE MetricType = "COSINE"
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
	// BGE_LARGE_ZH 1024
	BGE_LARGE_ZH EmbeddingModel = "bge-large-zh"
	// MULTILINGUAL 768
	MULTILINGUAL_E5_BASE = "multilingual-e5-base"
	// E5_LARGE_V2 1024
	E5_LARGE_V2 = "e5-large-v2"
	// TEXT2VEC_LARGE_CHINESE 1024
	TEXT2VEC_LARGE_CHINESE = "text2vec-large-chinese"
)
