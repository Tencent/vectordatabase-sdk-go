package model

import (
	"time"
)

type ClientOption struct {
	// Timeout: default 5s
	Timeout time.Duration
	// MaxIdldConnPerHost: default 10
	MaxIdldConnPerHost int
	// IdleConnTimeout: default 1m
	IdleConnTimeout time.Duration
}

type CommmonResponse struct {
	// Code: 0 means success, other means failure.
	Code int32 `json:"code,omitempty"`
	// Msg: response msg
	Msg string `json:"msg,omitempty"`
}

type SearchParams struct {
	Nprobe uint32  `json:"nprobe,omitempty"` // 搜索时查找的聚类数量，使用索引默认值即可
	Ef     uint32  `json:"ef,omitempty"`     // HNSW
	Radius float32 `json:"radius,omitempty"` // 距离阈值,范围搜索时有效
}
