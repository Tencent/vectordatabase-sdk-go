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
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/ai_document_set"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/document"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/user"
	api_user "github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/user"
	"github.com/tencentyun/cos-go-sdk-v5"
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

	UploadFile(ctx context.Context, databaseName, collectionName string, param UploadFileParams) (result *UploadFileResult, err error)

	GetImageUrl(ctx context.Context, databaseName, collectionName string,
		param GetImageUrlParams) (result *GetImageUrlResult, err error)
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

type UploadFileParams struct {
	FileName           string
	LocalFilePath      string
	Reader             io.Reader
	SplitterPreprocess ai_document_set.DocumentSplitterPreprocess
	EmbeddingModel     string
	ParsingProcess     *api.ParsingProcess
	FieldMappings      map[string]string
	MetaData           map[string]interface{}
}

type UploadFileResult struct {
	FileName        string
	CosEndpoint     string
	CosRegion       string
	CosBucket       string
	UploadPath      string
	Credentials     *ai_document_set.Credentials
	UploadCondition *ai_document_set.UploadCondition
}

func (i *implementerFlatDocument) UploadFile(ctx context.Context, databaseName, collectionName string, param UploadFileParams) (result *UploadFileResult, err error) {
	return uploadFile(ctx, i, databaseName, collectionName, param)
}

func checkUploadFileParam(ctx context.Context, param *UploadFileParams) (size int64, reader io.ReadCloser, err error) {
	if param.FileName == "" {
		if param.LocalFilePath == "" {
			return 0, nil, errors.New("need param: FileName or LocalFilePath")
		}
		param.FileName = filepath.Base(param.LocalFilePath)
	}
	fileType := strings.ToLower(filepath.Ext(param.FileName))
	isMarkdown := false
	if fileType == "" || fileType == string(MarkdownFileType) || fileType == string(MdFileType) {
		isMarkdown = true
	}
	if !isMarkdown && param.SplitterPreprocess.ChunkSplitter != nil && *param.SplitterPreprocess.ChunkSplitter != "" {
		log.Printf("[Warning] %s", "param SplitterPreprocess.ChunkSplitter will be ommitted, "+
			"because only markdown filetype supports defining ChunkSplitter")
	}
	if param.LocalFilePath != "" {
		fd, err := os.Open(param.LocalFilePath)
		if err != nil {
			return 0, nil, err
		}
		reader = fd
		fstat, err := fd.Stat()
		if err != nil {
			return 0, nil, err
		}
		size = fstat.Size()
	} else {
		bytesBuf := bytes.NewBuffer(nil)
		written, err := io.Copy(bytesBuf, param.Reader)
		if err != nil {
			return 0, nil, err
		}

		size = written
		reader = io.NopCloser(bytesBuf)
	}

	if size == 0 {
		return 0, nil, errors.New("file size cannot be 0")
	}

	return size, reader, nil
}
func uploadFile(ctx context.Context, cli SdkClient, databaseName, collectionName string,
	param UploadFileParams) (result *UploadFileResult, err error) {
	size, reader, err := checkUploadFileParam(ctx, &param)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	byteLength := uint64(size)

	req := new(document.UploadUrlReq)
	req.Database = databaseName
	req.Collection = collectionName
	req.FileName = param.FileName
	req.EmbeddingModel = param.EmbeddingModel
	req.SplitterPreprocess = &param.SplitterPreprocess
	req.ByteLength = &byteLength

	if param.ParsingProcess != nil {
		req.ParsingProcess = new(api.ParsingProcess)
		req.ParsingProcess.ParsingType = param.ParsingProcess.ParsingType
	}
	req.FieldMappings = make(map[string]string)
	for field, mapping := range param.FieldMappings {
		req.FieldMappings[field] = mapping
	}
	res := new(document.UploadUrlRes)
	err = cli.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}
	if res.Warning != "" {
		log.Printf("[Warning] %s", res.Warning)
	}
	if res.UploadCondition != nil && size > res.UploadCondition.MaxSupportContentLength {
		return nil, fmt.Errorf("fileSize is invalid, support max content length is %v bytes", res.UploadCondition.MaxSupportContentLength)
	}
	if res.Credentials == nil {
		return nil, fmt.Errorf("get credentials for uploading file failed")
	}

	result = new(UploadFileResult)
	result.FileName = param.FileName
	result.CosEndpoint = res.CosEndpoint
	result.CosRegion = res.CosRegion
	result.CosBucket = res.CosBucket
	result.UploadPath = res.UploadPath
	result.Credentials = res.Credentials
	result.UploadCondition = res.UploadCondition

	u, _ := url.Parse(res.CosEndpoint)
	b := &cos.BaseURL{BucketURL: u}

	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     res.Credentials.TmpSecretID,  // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/ai_document_set/product/598/37140
			SecretKey:    res.Credentials.TmpSecretKey, // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/ai_document_set/product/598/37140
			SessionToken: res.Credentials.SessionToken,
		},
	})

	header := make(http.Header)

	marshalData, err := json.Marshal(param.MetaData)
	if err != nil {
		return nil, fmt.Errorf("put param MetaData into cos header failed, err: %v", err.Error())
	}

	header.Add("x-cos-meta-data", url.QueryEscape(base64.StdEncoding.EncodeToString(marshalData)))

	headerData, err := json.Marshal(header)
	if err != nil {
		return nil, fmt.Errorf("marshal cos header failed, err: %v", err.Error())
	}
	if len(headerData) > 2048 {
		return nil, fmt.Errorf("cos header for param MetaData is too large, it can not be more than 2k")
	}

	if param.LocalFilePath != "" {
		opt := &cos.MultiUploadOptions{
			OptIni: &cos.InitiateMultipartUploadOptions{
				nil,
				&cos.ObjectPutHeaderOptions{
					XCosMetaXXX: &header,
					//Listener:    &cos.DefaultProgressListener{},
				},
			},
			// Whether to enable resume from breakpoint, default is false
			CheckPoint: true,
			PartSize:   5,
		}

		_, _, err = c.Object.Upload(ctx, res.UploadPath, param.LocalFilePath, opt)
		if err != nil {
			return nil, err
		}
	} else {

		opt := &cos.ObjectPutOptions{
			ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
				ContentLength: size,
				XCosMetaXXX:   &header,
			},
		}
		_, err = c.Object.Put(ctx, res.UploadPath, reader, opt)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

type GetImageUrlParams struct {
	FileName    string
	DocumentIds []string
}

type GetImageUrlResult struct {
	Images [][]document.ImageInfo
}

func (i *implementerFlatDocument) GetImageUrl(ctx context.Context, databaseName, collectionName string,
	param GetImageUrlParams) (result *GetImageUrlResult, err error) {
	return getImageUrl(ctx, i.SdkClient, databaseName, collectionName, param)
}

func getImageUrl(ctx context.Context, cli SdkClient, databaseName, collectionName string,
	param GetImageUrlParams) (result *GetImageUrlResult, err error) {
	req := new(document.GetImageUrlReq)
	req.Database = databaseName
	req.Collection = collectionName
	req.FileName = param.FileName
	req.DocumentIds = param.DocumentIds

	res := new(document.GetImageUrlRes)
	err = cli.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}

	result = new(GetImageUrlResult)
	result.Images = res.Images
	return result, nil
}
