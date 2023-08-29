package engine

import (
	"context"
	"errors"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/internal/engine/api/alias"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/model"
)

type implementerAlias struct {
	model.SdkClient
	databaseName string
}

func (i *implementerAlias) AliasSet(ctx context.Context, collectionName, aliasName string) error {
	req := new(alias.SetReq)
	req.Database = i.databaseName
	req.Collection = collectionName
	req.Alias = aliasName
	res := new(alias.SetRes)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return err
	}
	return nil
}

func (i *implementerAlias) AliasDrop(ctx context.Context, aliasName string) error {
	req := new(alias.DropReq)
	req.Database = i.databaseName
	req.Alias = aliasName
	res := new(alias.DropRes)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return err
	}
	return nil
}

func (i *implementerAlias) AliasDescribe(ctx context.Context, aliasName string) (*model.Alias, error) {
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

func (i *implementerAlias) AliasList(ctx context.Context) ([]*model.Alias, error) {
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
	// todo
	return aliases, nil
}
