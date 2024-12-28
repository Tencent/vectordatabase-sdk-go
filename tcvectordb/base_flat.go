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

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/user"
	api_user "github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/user"
)

type FlatInterface interface {
	// [Upsert] upserts documents into a collection.
	Upsert(ctx context.Context, databaseName, collectionName string, documents interface{}, params ...*UpsertDocumentParams) (result *UpsertDocumentResult, err error)

	// [Query] queries documents that satisfies the condition from the collection.
	Query(ctx context.Context, databaseName, collectionName string, documentIds []string, params ...*QueryDocumentParams) (result *QueryDocumentResult, err error)

	// [Search] returns the most similar topK vectors by the given vectors.
	// Search is a Batch API.
	Search(ctx context.Context, databaseName, collectionName string, vectors [][]float32, params ...*SearchDocumentParams) (result *SearchDocumentResult, err error)

	// [HybridSearch] retrieves both dense and sparse vectors to return the most similar topK vectors.
	HybridSearch(ctx context.Context, databaseName, collectionName string, params HybridSearchDocumentParams) (result *SearchDocumentResult, err error)

	// [SearchById] returns the most similar topK vectors by the given documentIds.
	SearchById(ctx context.Context, databaseName, collectionName string, documentIds []string, params ...*SearchDocumentParams) (result *SearchDocumentResult, err error)

	// [SearchByText] returns the most similar topK vectors by the given text map.
	// The texts will be firstly embedded into vectors using the embedding model of the collection on the server.
	SearchByText(ctx context.Context, databaseName, collectionName string, text map[string][]string, params ...*SearchDocumentParams) (result *SearchDocumentResult, err error)

	// [Delete] deletes documents by conditions.
	Delete(ctx context.Context, databaseName, collectionName string, param DeleteDocumentParams) (result *DeleteDocumentResult, err error)

	// [Update] updates documents by conditions.
	Update(ctx context.Context, databaseName, collectionName string, param UpdateDocumentParams) (result *UpdateDocumentResult, err error)

	// [Count] counts the number of documents in a collection that satisfy the specified filter conditions.
	Count(ctx context.Context, databaseName, collectionName string,
		params ...CountDocumentParams) (*CountDocumentResult, error)

	// [CreateUser] creates the user with the password.
	CreateUser(ctx context.Context, param CreateUserParams) error

	// [GrantToUser] grants the privileges to the specific user.
	GrantToUser(ctx context.Context, param GrantToUserParams) error

	// [RevokeFromUser] revokes the privileges from the specific user.
	RevokeFromUser(ctx context.Context, param RevokeFromUserParams) error

	// [DescribeUser] describes the specific user's detail, including createTime and privileges.
	DescribeUser(ctx context.Context, param DescribeUserParams) (result *DescribeUserResult, err error)

	// [ListUser] retrieves the details of all users for this instance, including their creation times and privileges.
	ListUser(ctx context.Context) (result *ListUserResult, err error)

	// [DropUser] drops the specific user.
	DropUser(ctx context.Context, param DropUserParams) error

	// [ChangePassword] changes the password for the specific user.
	ChangePassword(ctx context.Context, param ChangePasswordParams) error
}

// [CreateUserParams] holds the parameters for creating the user.
//
// Fields:
//   - User: (Required) The username to create.
//   - Password: (Required) The password of this user.
type CreateUserParams struct {
	User     string
	Password string
}

// [CreateUser] creates the user with the password.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - param: A [CreateUserParams] object that includes the other parameters for creating the user operation.
//     See [CreateUserParams] for more information.
//
// Returns an error if the operation fails.
func (i *implementerFlatDocument) CreateUser(ctx context.Context, param CreateUserParams) error {
	req := new(api_user.CreateReq)
	req.User = param.User
	req.Password = param.Password
	res := new(api_user.CreateReq)
	err := i.Request(ctx, req, res)
	if err != nil {
		return err
	}

	return nil
}

// [GrantToUserParams] holds the parameters for granting the privileges to the specific user.
//
// Fields:
//   - User: (Required) The username to create.
//   - Privileges: (Required) The list of [Privilege] to be granted to the specific user.
type GrantToUserParams struct {
	User       string
	Privileges []*api_user.Privilege
}

// [GrantToUser] grants the privileges to the specific user.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - param: A [GrantToUserParams] object that includes the other parameters for granting the privileges operation.
//     See [GrantToUserParams] for more information.
//
// Returns an error if the operation fails.
func (i *implementerFlatDocument) GrantToUser(ctx context.Context, param GrantToUserParams) error {
	req := new(api_user.GrantReq)
	req.User = param.User
	req.Privileges = param.Privileges
	res := new(api_user.GrantRes)
	err := i.Request(ctx, req, res)
	if err != nil {
		return err
	}

	return nil
}

// [RevokeFromUserParams] holds the parameters for revoking the privileges from the specific user.
//
// Fields:
//   - User: (Required) The username to create.
//   - Privileges: (Required) The list of [Privilege] to be revoked from the specific user.
type RevokeFromUserParams struct {
	User       string
	Privileges []*api_user.Privilege
}

// [RevokeFromUser] revokes the privileges from the specific user.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - param: A [RevokeFromUserParams] object that includes the other parameters for revoking the privileges operation.
//     See [RevokeFromUserParams] for more information.
//
// Returns an error if the operation fails.
func (i *implementerFlatDocument) RevokeFromUser(ctx context.Context, param RevokeFromUserParams) error {
	req := new(api_user.RevokeReq)
	req.User = param.User
	req.Privileges = param.Privileges
	res := new(api_user.RevokeRes)
	err := i.Request(ctx, req, res)
	if err != nil {
		return err
	}

	return nil
}

// [DescribeUserParams] holds the parameters for describing the user's detail.
//
// Fields:
//   - User: (Required) The username to describe.
type DescribeUserParams struct {
	User string
}

// [DescribeUserResult] holds the results for describing the user's detail.
//
// Fields:
//   - User: The username.
//   - CreateTime: The creation time of the user.
//   - Privileges: The list of [Privilege] which the user has.
type DescribeUserResult struct {
	User       string               `json:"user,omitempty"`
	CreateTime string               `json:"createTime,omitempty"`
	Privileges []api_user.Privilege `json:"privileges,omitempty"`
}

// [DescribeUser] describes the specific user's detail, including createTime and privileges.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - param: A [DescribeUserParams] object that includes the other parameters for describing the user's detail.
//     See [DescribeUserParams] for more information.
//
// Returns a pointer to a [DescribeUserResult] object or an error.
func (i *implementerFlatDocument) DescribeUser(ctx context.Context, param DescribeUserParams) (result *DescribeUserResult, err error) {
	req := new(api_user.DescribeReq)
	req.User = param.User
	res := new(api_user.DescribeRes)
	err = i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}

	result = new(DescribeUserResult)
	result.User = res.User
	result.CreateTime = res.CreateTime
	result.Privileges = res.Privileges

	return result, nil
}

// [ListUserResult] holds the results for listing the details of all users for this instance.
//
// Fields:
//   - Users: The list of [UserPrivileges] for this instance.
type ListUserResult struct {
	Users []user.UserPrivileges
}

// [ListUser] retrieves the details of all users for this instance, including their creation times and privileges.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//
// Returns a pointer to a [ListUserResult] object or an error.
func (i *implementerFlatDocument) ListUser(ctx context.Context) (result *ListUserResult, err error) {
	req := new(api_user.ListReq)
	res := new(api_user.ListRes)
	err = i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}

	result = new(ListUserResult)
	for _, userPrivileges := range res.Users {
		result.Users = append(result.Users, *userPrivileges)

	}
	return result, nil
}

// [DropUserParams] holds the parameters for dropping the user.
//
// Fields:
//   - User: (Required) The username to drop.
type DropUserParams struct {
	User string
}

// [DropUser] drops the specific user.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - param: A [DropUserParams] object that includes the other parameters for dropping the user.
//     See [DropUserParams] for more information.
//
// Returns an error if the operation fails.
func (i *implementerFlatDocument) DropUser(ctx context.Context, param DropUserParams) error {
	req := new(api_user.DropReq)
	req.User = param.User
	res := new(api_user.DropRes)
	err := i.Request(ctx, req, res)
	if err != nil {
		return err
	}

	return nil
}

// [ChangePasswordParams] holds the parameters for changing the password.
//
// Fields:
//   - User: (Required) The username to change password.
//   - Password: (Required) The password to be changed for the user.
type ChangePasswordParams struct {
	User     string
	Password string
}

// [ChangePassword] changes the password for the specific user.
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - param: A [ChangePasswordParams] object that includes the other parameters for changing the password.
//     See [ChangePasswordParams] for more information.
//
// Returns an error if the operation fails.
func (i *implementerFlatDocument) ChangePassword(ctx context.Context, param ChangePasswordParams) error {
	req := new(api_user.ChangePasswordReq)
	req.User = param.User
	req.Password = param.Password
	res := new(api_user.ChangePasswordRes)
	err := i.Request(ctx, req, res)
	if err != nil {
		return err
	}

	return nil
}
