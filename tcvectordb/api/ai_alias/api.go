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

package ai_alias

import "github.com/tencent/vectordatabase-sdk-go/tcvectordb/api"

type SetReq struct {
	api.Meta       `path:"/ai/alias/set" tags:"Alias" method:"Post" summary:"指定集合别名，新增/修改"`
	Database       string `json:"database,omitempty"`
	CollectionView string `json:"collectionView,omitempty"`
	Alias          string `json:"alias,omitempty"`
}

type SetRes struct {
	api.CommonRes
	AffectedCount int `json:"affectedCount,omitempty"`
}

type DeleteReq struct {
	api.Meta `path:"/ai/alias/delete" tags:"Alias" method:"Post" summary:"删除集合别名"`
	Database string `json:"database,omitempty"`
	Alias    string `json:"alias,omitempty"`
}

type DeleteRes struct {
	api.CommonRes
	AffectedCount int `json:"affectedCount,omitempty"`
}
