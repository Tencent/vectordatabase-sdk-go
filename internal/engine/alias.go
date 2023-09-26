package engine

import (
	"context"
	"errors"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/internal/engine/api/alias"
)

var _ entity.AliasInterface = &implementerAlias{}

type implementerAlias struct {
	entity.SdkClient
	databaseName string
}

func (i *implementerAlias) SetAlias(ctx context.Context, collectionName, aliasName string, option *entity.SetAliasOption) (*entity.AliasResult, error) {
	req := new(alias.SetReq)
	req.Database = i.databaseName
	req.Collection = collectionName
	req.Alias = aliasName
	res := new(alias.SetRes)
	result := new(entity.AliasResult)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return result, err
	}
	result.AffectedCount = int(res.AffectedCount)
	return result, nil
}

func (i *implementerAlias) DeleteAlias(ctx context.Context, aliasName string, option *entity.DeleteAliasOption) (*entity.AliasResult, error) {
	req := new(alias.DeleteReq)
	req.Database = i.databaseName
	req.Alias = aliasName
	res := new(alias.DeleteRes)
	result := new(entity.AliasResult)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return result, err
	}
	result.AffectedCount = int(res.AffectedCount)
	return result, nil
}

func (i *implementerAlias) DescribeAlias(ctx context.Context, aliasName string, option *entity.DescribeAliasOption) (*entity.AliasResult, error) {
	req := new(alias.DescribeReq)
	req.Database = i.databaseName
	req.Alias = aliasName
	res := new(alias.DescribeRes)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return nil, err
	}
	if len(res.Aliases) == 0 {
		return nil, errors.New("alias not found")
	}
	alias := &entity.AliasResult{
		Collection: res.Aliases[0].Collection,
	}
	return alias, nil
}

func (i *implementerAlias) ListAlias(ctx context.Context, option *entity.ListAliasOption) ([]*entity.AliasResult, error) {
	req := new(alias.ListReq)
	req.Database = i.databaseName
	res := new(alias.ListRes)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return nil, err
	}
	if len(res.Aliases) == 0 {
		return nil, errors.New("alias not found")
	}
	var aliases []*entity.AliasResult
	return aliases, nil
}
