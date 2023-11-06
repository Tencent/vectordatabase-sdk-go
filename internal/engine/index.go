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
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity/api/index"
)

var _ entity.IndexInterface = &implementerIndex{}

type implementerIndex struct {
	entity.SdkClient
	database entity.Database
}

func (i *implementerIndex) IndexRebuild(ctx context.Context, collectionName string, option *entity.IndexRebuildOption) (*entity.IndexReBuildResult, error) {
	if i.database.IsAIDatabase() {
		return nil, entity.AIDbTypeError
	}
	req := new(index.RebuildReq)
	req.Database = i.database.DatabaseName
	req.Collection = collectionName

	req.DropBeforeRebuild = option.DropBeforeRebuild
	req.Throttle = int32(option.Throttle)

	res := new(index.RebuildRes)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return nil, err
	}
	result := new(entity.IndexReBuildResult)
	result.TaskIds = res.TaskIds
	return result, nil
}
