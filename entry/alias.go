package entry

import (
	"context"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/model"
)

type AliasInterface interface {
	SdkClient
	SetAlias(ctx context.Context, collectionName, aliasName string, option *SetAliasOption) (*AliasResult, error)
	DeleteAlias(ctx context.Context, aliasName string, option *DeleteAliasOption) (*AliasResult, error)
	DescribeAlias(ctx context.Context, aliasName string, option *DescribeAliasOption) (*model.Alias, error)
	ListAlias(ctx context.Context, option *ListAliasOption) ([]*model.Alias, error)
}

type AliasResult struct {
	AffectedCount int
}

type SetAliasOption struct{}

type DeleteAliasOption struct{}

type DescribeAliasOption struct{}

type ListAliasOption struct{}
