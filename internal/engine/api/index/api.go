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

package index

import (
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/internal/engine/api"
)

type RebuildReq struct {
	api.Meta          `path:"/index/rebuild" tags:"Index" method:"Post" summary:"重建整个collection的所有索引"`
	Database          string `protobuf:"bytes,1,opt,name=database,proto3" json:"database,omitempty"`
	Collection        string `protobuf:"bytes,2,opt,name=collection,proto3" json:"collection,omitempty"`
	DropBeforeRebuild bool   `protobuf:"varint,3,opt,name=dropBeforeRebuild,proto3" json:"dropBeforeRebuild,omitempty"`
	Throttle          int32  `protobuf:"varint,4,opt,name=throttle,proto3" json:"throttle,omitempty"`
	DisableTrain      bool   `protobuf:"varint,5,opt,name=disable_train,json=disableTrain,proto3" json:"disable_train,omitempty"`
	ForceRebuild      bool   `protobuf:"varint,6,opt,name=force_rebuild,json=forceRebuild,proto3" json:"force_rebuild,omitempty"`
}

type RebuildRes struct {
	Code     int32    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg      string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Redirect string   `protobuf:"bytes,3,opt,name=redirect,proto3" json:"redirect,omitempty"`
	TaskIds  []string `protobuf:"bytes,4,rep,name=task_ids,json=taskIds,proto3" json:"task_ids,omitempty"`
}
