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
)

var _ IndexInterface = &implementerIndex{}

type IndexInterface interface {
	SdkClient

	// [RebuildIndex] rebuilds all indexes under the specified collection.
	RebuildIndex(ctx context.Context, params ...*RebuildIndexParams) (result *RebuildIndexResult, err error)

	// [AddIndex] adds scalar field index to an existing collection.
	AddIndex(ctx context.Context, params ...*AddIndexParams) (err error)
}

type implementerIndex struct {
	SdkClient
	flat       FlatIndexInterface
	database   *Database
	collection *Collection
}

type RebuildIndexResult struct {
	TaskIds []string
}

// [RebuildIndex] rebuilds all indexes under the specified collection.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - params: A pointer to a [RebuildIndexParams] object that includes the other parameters for the rebuilding indexes operation.
//     See [RebuildIndexParams] for more information.
//
// Notes: The name of the database and the name of collection are from the fields of [implementerIndex].
//
// Returns a pointer to a [RebuildIndexResult] object or an error.
func (i *implementerIndex) RebuildIndex(ctx context.Context, params ...*RebuildIndexParams) (*RebuildIndexResult, error) {
	return i.flat.RebuildIndex(ctx, i.database.DatabaseName, i.collection.CollectionName, params...)
}

// [AddIndex] adds scalar field index to an existing collection.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - params: A pointer to a [AddIndexParams] object that includes the other parameters for the adding scalar field index operation.
//     See [AddIndexParams] for more information.
//
// Notes: The name of the database and the name of collection are from the fields of [implementerIndex].
//
// Returns an error if the addition fails.
func (i *implementerIndex) AddIndex(ctx context.Context, params ...*AddIndexParams) error {
	return i.flat.AddIndex(ctx, i.database.DatabaseName, i.collection.CollectionName, params...)
}
