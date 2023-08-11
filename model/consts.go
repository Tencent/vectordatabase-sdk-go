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
