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

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity/api/ai_alias"
)

var _ entity.AIAliasInterface = &implementerAIAlias{}

type implementerAIAlias struct {
	entity.SdkClient
	databaseName string
}

func (i *implementerAIAlias) SetAlias(ctx context.Context, collectionName, aliasName string, option *entity.SetAIAliasOption) (*entity.SetAIAliasResult, error) {
	req := new(ai_alias.SetReq)
	req.Database = i.databaseName
	req.Collection = collectionName
	req.Alias = aliasName
	res := new(ai_alias.SetRes)

	result := new(entity.SetAIAliasResult)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return result, err
	}
	result.AffectedCount = res.AffectedCount
	return result, nil
}

func (i *implementerAIAlias) DeleteAlias(ctx context.Context, aliasName string, option *entity.DeleteAIAliasOption) (*entity.DeleteAIAliasResult, error) {
	req := new(ai_alias.DeleteReq)
	req.Database = i.databaseName
	req.Alias = aliasName
	res := new(ai_alias.DeleteRes)

	result := new(entity.DeleteAIAliasResult)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return result, err
	}
	result.AffectedCount = res.AffectedCount
	return result, nil
}
