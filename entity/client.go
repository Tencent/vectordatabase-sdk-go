package entity

import (
	"context"
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

// SdkClient the http client interface
type SdkClient interface {
	Request(ctx context.Context, req, res interface{}) error
	Options() ClientOption
	WithTimeout(d time.Duration)
	Debug(v bool)
	Close()
}
