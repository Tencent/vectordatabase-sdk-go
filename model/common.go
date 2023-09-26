package model

import (
	"time"
)

type ClientOption struct {
	// Timeout: default 5s
	Timeout time.Duration
	// MaxIdldConnPerHost: default 2
	MaxIdldConnPerHost int
	// IdleConnTimeout: default 0 means no limit
	IdleConnTimeout time.Duration
	// ReadConsistency: default: EventualConsistency
	ReadConsistency ReadConsistency
}

type CommmonResponse struct {
	// Code: 0 means success, other means failure.
	Code int32 `json:"code,omitempty"`
	// Msg: response msg
	Msg string `json:"msg,omitempty"`
}
