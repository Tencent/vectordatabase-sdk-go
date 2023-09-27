// Copyright (C) 2023 Tencent Cloud.
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the vectordb-sdk-java), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is furnished
// to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
// PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
// SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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
