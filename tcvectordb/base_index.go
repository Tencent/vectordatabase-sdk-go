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
	RebuildIndex(ctx context.Context, params ...*RebuildIndexParams) (result *RebuildIndexResult, err error)
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

func (i *implementerIndex) RebuildIndex(ctx context.Context, params ...*RebuildIndexParams) (*RebuildIndexResult, error) {
	return i.flat.RebuildIndex(ctx, i.database.DatabaseName, i.collection.CollectionName, params...)
}

func (i *implementerIndex) AddIndex(ctx context.Context, params ...*AddIndexParams) error {
	return i.flat.AddIndex(ctx, i.database.DatabaseName, i.collection.CollectionName, params...)
}
