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

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/alias"
)

var _ AliasInterface = &implementerAlias{}

type AliasInterface interface {
	SdkClient

	// [SetAlias] sets an alias for collection.
	SetAlias(ctx context.Context, collectionName, aliasName string) (result *SetAliasResult, err error)

	// [DeleteAlias] deletes the alias in the database.
	DeleteAlias(ctx context.Context, aliasName string) (result *DeleteAliasResult, err error)
}

type implementerAlias struct {
	SdkClient
	database *Database
}

type SetAliasResult struct {
	AffectedCount int
}

// [SetAlias] sets an alias for collection.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - collectionName: The name of the collection.
//   - aliasName: The alias name to set for the collection.
//
// Notes: The name of the database is from the field of [implementerAlias].
//
// Returns a pointer to a [SetAliasResult] object or an error.
func (i *implementerAlias) SetAlias(ctx context.Context, collectionName, aliasName string) (*SetAliasResult, error) {
	if i.database.IsAIDatabase() {
		return nil, AIDbTypeError
	}
	req := new(alias.SetReq)
	res := new(alias.SetRes)

	req.Database = i.database.DatabaseName
	req.Collection = collectionName
	req.Alias = aliasName

	result := new(SetAliasResult)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return result, err
	}
	result.AffectedCount = res.AffectedCount
	return result, nil
}

type DeleteAliasResult struct {
	AffectedCount int
}

// [DeleteAlias] deletes the alias in the database.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - aliasName: The alias name to delete.
//
// Notes: The name of the database is from the field of [implementerAlias].
//
// Returns a pointer to a [DeleteAliasResult] object or an error.
func (i *implementerAlias) DeleteAlias(ctx context.Context, aliasName string) (*DeleteAliasResult, error) {
	if i.database.IsAIDatabase() {
		return nil, AIDbTypeError
	}
	req := new(alias.DeleteReq)
	res := new(alias.DeleteRes)

	req.Database = i.database.DatabaseName
	req.Alias = aliasName

	result := new(DeleteAliasResult)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return result, err
	}
	result.AffectedCount = res.AffectedCount
	return result, nil
}
