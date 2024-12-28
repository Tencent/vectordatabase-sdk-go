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

package user

import "github.com/tencent/vectordatabase-sdk-go/tcvectordb/api"

type CreateReq struct {
	api.Meta `path:"/user/create" tags:"User" method:"Post" summary:"创建用户"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}

type CreateRes struct {
	api.CommonRes
}

type GrantReq struct {
	api.Meta   `path:"/user/grant" tags:"User" method:"Post" summary:"授予用户权限"`
	User       string       `json:"user,omitempty"`
	Privileges []*Privilege `json:"privileges,omitempty"`
}

type Privilege struct {
	Resource string   `json:"resource,omitempty"`
	Actions  []string `json:"actions,omitempty"`
}

type GrantRes struct {
	api.CommonRes
}

type RevokeReq struct {
	api.Meta   `path:"/user/revoke" tags:"User" method:"Post" summary:"撤销用户权限"`
	User       string       `json:"user,omitempty"`
	Privileges []*Privilege `json:"privileges,omitempty"`
}

type RevokeRes struct {
	api.CommonRes
}

type DescribeReq struct {
	api.Meta `path:"/user/describe" tags:"User" method:"Post" summary:"返回用户权限"`
	User     string `json:"user,omitempty"`
}

type DescribeRes struct {
	api.CommonRes
	User       string      `json:"user,omitempty"`
	CreateTime string      `json:"createTime,omitempty"`
	Privileges []Privilege `json:"privileges,omitempty"`
}

type ListReq struct {
	api.Meta `path:"/user/list" tags:"User" method:"Get" summary:"返回实例所有用户权限"`
}

type ListRes struct {
	api.CommonRes
	Users []*UserPrivileges `json:"users,omitempty"`
}

type UserPrivileges struct {
	User       string      `json:"user,omitempty"`
	CreateTime string      `json:"createTime,omitempty"`
	Privileges []Privilege `json:"privileges,omitempty"`
}

type DropReq struct {
	api.Meta `path:"/user/drop" tags:"User" method:"Post" summary:"删除实例用户"`
	User     string `json:"user,omitempty"`
}

type DropRes struct {
	api.CommonRes
}

type ChangePasswordReq struct {
	api.Meta `path:"/user/changePassword" tags:"User" method:"Post" summary:"修改自定义用户的密码"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}

type ChangePasswordRes struct {
	api.CommonRes
}
