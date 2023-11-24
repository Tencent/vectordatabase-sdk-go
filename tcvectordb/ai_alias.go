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

package tcvectordb

import (
	"context"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb/api/ai_alias"
)

var _ AIAliasInterface = &implementerAIAlias{}

type AIAliasInterface interface {
	SdkClient
	SetAlias(ctx context.Context, collectionView, aliasName string) (result *SetAIAliasResult, err error)
	DeleteAlias(ctx context.Context, aliasName string) (result *DeleteAIAliasResult, err error)
}

type implementerAIAlias struct {
	SdkClient
	database *AIDatabase
}

type SetAIAliasResult struct {
	AffectedCount int
}

func (i *implementerAIAlias) SetAlias(ctx context.Context, collectionView, aliasName string) (*SetAIAliasResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(ai_alias.SetReq)
	req.Database = i.database.DatabaseName
	req.CollectionView = collectionView
	req.Alias = aliasName
	res := new(ai_alias.SetRes)

	result := new(SetAIAliasResult)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return result, err
	}
	result.AffectedCount = res.AffectedCount
	return result, nil
}

type DeleteAIAliasResult struct {
	AffectedCount int
}

func (i *implementerAIAlias) DeleteAlias(ctx context.Context, aliasName string) (*DeleteAIAliasResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(ai_alias.DeleteReq)
	req.Database = i.database.DatabaseName
	req.Alias = aliasName
	res := new(ai_alias.DeleteRes)

	result := new(DeleteAIAliasResult)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return result, err
	}
	result.AffectedCount = res.AffectedCount
	return result, nil
}
