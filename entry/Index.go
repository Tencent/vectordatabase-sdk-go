package entry

import (
	"context"
)

type IndexInterface interface {
	SdkClient
	IndexRebuild(ctx context.Context, collectionName string, option *IndexRebuildOption) (result *IndexReBuildResult, err error)
}

type IndexReBuildResult struct {
	TaskIds []string
}

type IndexRebuildOption struct {
	DropBeforeRebuild bool
	Throttle          int
}
