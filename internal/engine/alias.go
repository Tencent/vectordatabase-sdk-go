package engine

import (
	"context"
	"errors"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entry"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/internal/engine/api/alias"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/model"
)

var _ entry.AliasInterface = &implementerAlias{}

type implementerAlias struct {
	entry.SdkClient
	databaseName string
}

func (i *implementerAlias) SetAlias(ctx context.Context, collectionName, aliasName string, option *entry.SetAliasOption) (*entry.AliasResult, error) {
	req := new(alias.SetReq)
	req.Database = i.databaseName
	req.Collection = collectionName
	req.Alias = aliasName
	res := new(alias.SetRes)
	result := new(entry.AliasResult)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return result, err
	}
	result.AffectedCount = int(res.AffectedCount)
	return result, nil
}

func (i *implementerAlias) DeleteAlias(ctx context.Context, aliasName string, option *entry.DeleteAliasOption) (*entry.AliasResult, error) {
	req := new(alias.DeleteReq)
	req.Database = i.databaseName
	req.Alias = aliasName
	res := new(alias.DeleteRes)
	result := new(entry.AliasResult)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return result, err
	}
	result.AffectedCount = int(res.AffectedCount)
	return result, nil
}

func (i *implementerAlias) DescribeAlias(ctx context.Context, aliasName string, option *entry.DescribeAliasOption) (*model.Alias, error) {
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
	alias := &model.Alias{
		Collection: res.Aliases[0].Collection,
	}
	return alias, nil
}

func (i *implementerAlias) ListAlias(ctx context.Context, option *entry.ListAliasOption) ([]*model.Alias, error) {
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
	var aliases []*model.Alias
	return aliases, nil
}
