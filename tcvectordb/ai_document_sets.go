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
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"log"

	"github.com/pkg/errors"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/ai_document_set"
	"github.com/tencentyun/cos-go-sdk-v5"
)

var _ AIDocumentSetsInterface = &implementerAIDocumentSets{}

type AIDocumentSetsInterface interface {
	SdkClient

	// [Query] queries documentSets that satisfies the condition from the collectionView.
	Query(ctx context.Context, params QueryAIDocumentSetParams) (*QueryAIDocumentSetResult, error)

	// [GetDocumentSetByName] gets a documentSet by the name of documentSet.
	GetDocumentSetByName(ctx context.Context, documentSetName string) (*GetAIDocumentSetResult, error)

	// [GetDocumentSetById] gets a documentSet by the documentSet id.
	GetDocumentSetById(ctx context.Context, documentSetId string) (*GetAIDocumentSetResult, error)

	// [GetChunks] gets chunks of a documentSet.
	GetChunks(ctx context.Context, param GetAIDocumentSetChunksParams) (*GetAIDocumentSetChunksResult, error)

	// [Search] returns the most similar topK chunks by the parameters of [SearchAIDocumentSetsParams].
	Search(ctx context.Context, param SearchAIDocumentSetsParams) (*SearchAIDocumentSetResult, error)

	// [DeleteByIds] deletes some documentSets by the list of documentSet ids.
	DeleteByIds(ctx context.Context, documentSetIds ...string) (result *DeleteAIDocumentSetResult, err error)

	// [DeleteByNames] deletes some documentSets by the list of documentSet names.
	DeleteByNames(ctx context.Context, documentSetNames ...string) (result *DeleteAIDocumentSetResult, err error)

	// [Delete] deletes some documentSets by the parameters of [DeleteAIDocumentSetParams].
	Delete(ctx context.Context, param DeleteAIDocumentSetParams) (*DeleteAIDocumentSetResult, error)

	// [Update] updates some documentSets in the collectionView.
	Update(ctx context.Context, updateFields map[string]interface{}, param UpdateAIDocumentSetParams) (*UpdateAIDocumentSetResult, error)

	// [LoadAndSplitText] uploads local file, which will be parsed and saved remotely.
	LoadAndSplitText(ctx context.Context, param LoadAndSplitTextParams) (result *LoadAndSplitTextResult, err error)

	// [GetCosTmpSecret] gets the temp secret to upload specific file.
	GetCosTmpSecret(ctx context.Context, param GetCosTmpSecretParams) (*GetCosTmpSecretResult, error)
}

type AIDocumentSet struct {
	AIDocumentSetInterface `json:"-"`
	DatabaseName           string                           `json:"databaseName"`
	CollectionViewName     string                           `json:"collectionViewName"`
	DocumentSetId          string                           `json:"documentSetId"`
	DocumentSetName        string                           `json:"documentSetName"`
	Text                   string                           `json:"text"`       // assign when use get api
	TextPrefix             string                           `json:"textPrefix"` // assign when use query api
	DocumentSetInfo        *ai_document_set.DocumentSetInfo `json:"documentSetInfo"`
	ScalarFields           map[string]Field
	SplitterPreprocess     *ai_document_set.DocumentSplitterPreprocess `json:"splitterPreprocess,omitempty"`
	ParsingProcess         *api.ParsingProcess                         `json:"parsingProcess,omitempty"`
}

type implementerAIDocumentSets struct {
	SdkClient
	database       *AIDatabase
	collectionView *AICollectionView
}

// [QueryAIDocumentSetParams] holds the parameters for querying documentSets to a collectionView.
//
// Fields:
//   - DocumentSetId: (Optional) The list of documentSet's ids to query. The maximum number of elements in the array is 20.
//   - DocumentSetName: (Optional) The list of documentSet's names to query. The maximum number of elements in the array is 20.
//   - Filter:  (Optional) Filter documentSets by [Filter] conditions before returning the result.
//   - Limit: (Required) Limit the number of documentSets returned.
//   - Offset: (Optional) Skip a specified number of documentSets in the query result set.
//   - OutputFields: (Optional) Return columns specified by the list of column names.
type QueryAIDocumentSetParams struct {
	DocumentSetId   []string `json:"documentSetId"`
	DocumentSetName []string `json:"documentSetName"`
	Filter          *Filter  `json:"filter"`
	Limit           int64    `json:"limit"`
	Offset          int64    `json:"offset"`
	OutputFields    []string `json:"outputFields,omitempty"`
}

type QueryAIDocumentSetResult struct {
	Count     uint64          `json:"count"`
	Documents []AIDocumentSet `json:"documents"`
}

// [Query] queries documentSets that satisfies the condition from the collectionView.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - param: A pointer to a [QueryAIDocumentSetParams] object that includes the other parameters for querying documentSets' operation.
//     See [QueryAIDocumentSetParams] for more information.
//
// Notes: The name of the database and the name of collectionView are from the fields of [implementerAIDocumentSets].
//
// Returns a pointer to a [QueryAIDocumentSetResult] object or an error.
func (i *implementerAIDocumentSets) Query(ctx context.Context, param QueryAIDocumentSetParams) (*QueryAIDocumentSetResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(ai_document_set.QueryReq)
	res := new(ai_document_set.QueryRes)

	req.Database = i.database.DatabaseName
	req.CollectionView = i.collectionView.connCollectionViewName
	req.Query = &ai_document_set.QueryCond{
		DocumentSetId:   param.DocumentSetId,
		DocumentSetName: param.DocumentSetName,
		Filter:          param.Filter.Cond(),
		Limit:           param.Limit,
		Offset:          param.Offset,
		OutputFields:    param.OutputFields,
	}

	result := new(QueryAIDocumentSetResult)
	err := i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}
	result.Count = res.Count
	for _, doc := range res.DocumentSets {
		result.Documents = append(result.Documents, *i.toDocumentSet(doc))
	}
	return result, nil
}

type GetAIDocumentSetParams struct {
	DocumentSetId   string `json:"documentSetId"`
	DocumentSetName string `json:"documentSetName"`
}

type GetAIDocumentSetResult struct {
	AIDocumentSet `json:"documentSets"`
	Count         uint64
}

// [GetDocumentSetByName] gets a documentSet by the name of documentSet.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - documentSetName: The name of the documentSet.
//
// Notes: The name of the database and the name of collectionView are from the fields of [implementerAIDocumentSets].
//
// Returns a pointer to a [GetAIDocumentSetResult] object or an error.
func (i *implementerAIDocumentSets) GetDocumentSetByName(ctx context.Context, documentSetName string) (*GetAIDocumentSetResult, error) {
	return i.get(ctx, GetAIDocumentSetParams{DocumentSetName: documentSetName})
}

// [GetDocumentSetById] gets a documentSet by the id of documentSet.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - documentSetId: The id of the documentSet.
//
// Notes: The name of the database and the name of collectionView are from the fields of [implementerAIDocumentSets].
//
// Returns a pointer to a [GetAIDocumentSetResult] object or an error.
func (i *implementerAIDocumentSets) GetDocumentSetById(ctx context.Context, documentSetId string) (*GetAIDocumentSetResult, error) {
	return i.get(ctx, GetAIDocumentSetParams{DocumentSetId: documentSetId})
}

func (i *implementerAIDocumentSets) get(ctx context.Context, param GetAIDocumentSetParams) (*GetAIDocumentSetResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(ai_document_set.GetReq)
	res := new(ai_document_set.GetRes)

	req.Database = i.database.DatabaseName
	req.CollectionView = i.collectionView.connCollectionViewName
	req.DocumentSetId = param.DocumentSetId
	req.DocumentSetName = param.DocumentSetName

	result := new(GetAIDocumentSetResult)
	err := i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}
	result.Count = res.Count
	result.AIDocumentSet = *i.toDocumentSet(res.DocumentSets)
	return result, nil
}

type GetAIDocumentSetChunksParams struct {
	DocumentSetId   string `json:"documentSetId"`
	DocumentSetName string `json:"documentSetName"`
	Limit           *int64 `json:"limit"`
	Offset          int64  `json:"offset"`
}

type GetAIDocumentSetChunksResult struct {
	DocumentSetId   string                  `json:"documentSetId"`
	DocumentSetName string                  `json:"documentSetName"`
	Count           uint64                  `json:"count"`
	Chunks          []ai_document_set.Chunk `json:"chunks"`
}

// [GetChunks] gets chunks of a documentSet.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - param: A pointer to a [GetAIDocumentSetChunksParams] object that includes the other parameters for querying documentSets' operation.
//     See [GetAIDocumentSetChunksParams] for more information.
//
// Notes: The name of the database and the name of collectionView are from the fields of [implementerAIDocumentSets].
//
// Returns a pointer to a [GetAIDocumentSetChunksResult] object or an error.
func (i *implementerAIDocumentSets) GetChunks(ctx context.Context, param GetAIDocumentSetChunksParams) (*GetAIDocumentSetChunksResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(ai_document_set.GetChunksReq)
	res := new(ai_document_set.GetChunksRes)

	req.Database = i.database.DatabaseName
	req.CollectionView = i.collectionView.connCollectionViewName
	req.DocumentSetId = param.DocumentSetId
	req.DocumentSetName = param.DocumentSetName
	req.Limit = param.Limit
	req.Offset = param.Offset

	result := new(GetAIDocumentSetChunksResult)
	err := i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}
	result.Count = res.Count
	result.DocumentSetId = res.DocumentSetId
	result.DocumentSetName = res.DocumentSetName
	result.Count = res.Count
	result.Chunks = append(result.Chunks, res.Chunks...)
	return result, nil
}

// [SearchAIDocumentSetsParams] holds the parameters for searching documentSets in a collectionView.
//
// Fields:
//   - Content: (Optional) The content to apply similarity search.
//   - DocumentSetName: (Optional) The list of documentSet's names to search. The maximum number of elements in the array is 10.
//   - RerankOption: (Optional) A pointer to a [RerankOption] object that includes the other parameters for reranking.
//     See [RerankOption] for more information.
//   - Filter: (Optional) Filter documentSets by [Filter] conditions when searching the results.
//   - Limit: (Required) The value of K for returning the top K most similar items.
type SearchAIDocumentSetsParams struct {
	Content         string                        `json:"content"`
	DocumentSetName []string                      `json:"documentSetName"`
	ExpandChunk     []int                         `json:"expandChunk"`  // 搜索结果中，向前、向后补齐几个chunk的上下文
	RerankOption    *ai_document_set.RerankOption `json:"rerankOption"` // 多路召回
	// MergeChunk  bool
	// Weights      SearchAIOptionWeight
	Filter *Filter `json:"filter"`
	Limit  int64   `json:"limit"`
}

type SearchAIDocumentSetResult struct {
	Documents []AISearchDocumentSet `json:"documents"`
}

// [Search] returns the most similar topK chunks by the parameters of [SearchAIDocumentSetsParams].
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - param: A pointer to a [SearchAIDocumentSetsParams] object that includes the other parameters for searching documentSets' operation.
//     See [SearchAIDocumentSetsParams] for more information.
//
// Notes: The name of the database and the name of collectionView are from the fields of [implementerAIDocumentSets].
//
// Returns a pointer to a [SearchAIDocumentSetResult] object or an error.
func (i *implementerAIDocumentSets) Search(ctx context.Context, param SearchAIDocumentSetsParams) (*SearchAIDocumentSetResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(ai_document_set.SearchReq)
	res := new(ai_document_set.SearchRes)

	req.Database = i.database.DatabaseName
	req.CollectionView = i.collectionView.connCollectionViewName
	req.ReadConsistency = string(i.SdkClient.Options().ReadConsistency)
	req.Search = new(ai_document_set.SearchCond)

	req.Search.Content = param.Content
	req.Search.DocumentSetName = param.DocumentSetName

	req.Search.Options = ai_document_set.SearchOption{
		ChunkExpand: param.ExpandChunk,
		// MergeChunk:  param.MergeChunk,
		// Weights: ai_document_set.SearchOptionWeight{
		// 	ChunkSimilarity: param.Weights.ChunkSimilarity,
		// 	WordSimilarity:  param.Weights.WordSimilarity,
		// 	WordBm25:        param.Weights.WordBm25,
		// },
	}
	if param.RerankOption != nil {
		req.Search.Options.RerankOption = &ai_document_set.RerankOption{
			Enable:                param.RerankOption.Enable,
			ExpectRecallMultiples: param.RerankOption.ExpectRecallMultiples,
		}
	}
	req.Search.Filter = param.Filter.Cond()
	req.Search.Limit = param.Limit

	err := i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}
	result := new(SearchAIDocumentSetResult)
	for _, doc := range res.Documents {
		result.Documents = append(result.Documents, *i.toSearchDocumentSet(doc))
	}
	return result, nil
}

// [DeleteAIDocumentSetParams] holds the parameters for deleting documentSets in a collectionView.
//
// Fields:
//   - DocumentSetNames: The list of documentSet's names to delete. The maximum number of elements in the array is 20.
//   - DocumentSetIds: The list of documentSet's ids to delete. The maximum number of elements in the array is 20.
//   - Filter: (Optional) Filter documentSets by [Filter] conditions to delete.
type DeleteAIDocumentSetParams struct {
	DocumentSetNames []string `json:"documentSetNames"`
	DocumentSetIds   []string `json:"documentSetIds"`
	Filter           *Filter  `json:"filter"`
}

type DeleteAIDocumentSetResult struct {
	AffectedCount uint64 `json:"affectedCount"`
}

// [DeleteByIds] deletes some documentSets by the list of documentSet ids.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - documentSetIds: The list of documentSet's ids to delete. The maximum number of elements in the array is 20.
//
// Notes: The name of the database and the name of collectionView are from the fields of [implementerAIDocumentSets].
//
// Returns a pointer to a [DeleteAIDocumentSetResult] object or an error.
func (i *implementerAIDocumentSets) DeleteByIds(ctx context.Context, documentSetIds ...string) (result *DeleteAIDocumentSetResult, err error) {
	return i.Delete(ctx, DeleteAIDocumentSetParams{DocumentSetIds: documentSetIds})
}

// [DeleteByNames] deletes some documentSets by the list of documentSet names.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - documentSetNames: The list of documentSet's names to delete. The maximum number of elements in the array is 20.
//
// Notes: The name of the database and the name of collectionView are from the fields of [implementerAIDocumentSets].
//
// Returns a pointer to a [DeleteAIDocumentSetResult] object or an error.
func (i *implementerAIDocumentSets) DeleteByNames(ctx context.Context, documentSetNames ...string) (result *DeleteAIDocumentSetResult, err error) {
	return i.Delete(ctx, DeleteAIDocumentSetParams{DocumentSetNames: documentSetNames})
}

// [Delete] deletes some documentSets by the parameters of [DeleteAIDocumentSetParams].
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - param: A [DeleteAIDocumentSetParams] object that includes the other parameters for deleting documentSets' operation.
//     See [DeleteAIDocumentSetParams] for more information.
//
// Notes: The name of the database and the name of collectionView are from the fields of [implementerAIDocumentSets].
//
// Returns a pointer to a [DeleteAIDocumentSetResult] object or an error.
func (i *implementerAIDocumentSets) Delete(ctx context.Context, param DeleteAIDocumentSetParams) (result *DeleteAIDocumentSetResult, err error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(ai_document_set.DeleteReq)
	res := new(ai_document_set.DeleteRes)

	req.Database = i.database.DatabaseName
	req.CollectionView = i.collectionView.connCollectionViewName
	req.Query = &ai_document_set.DeleteQueryCond{
		DocumentSetId:   param.DocumentSetIds,
		DocumentSetName: param.DocumentSetNames,
		Filter:          param.Filter.Cond(),
	}

	result = new(DeleteAIDocumentSetResult)
	err = i.Request(ctx, req, res)
	if err != nil {
		return
	}
	result.AffectedCount = res.AffectedCount
	return
}

// [UpdateAIDocumentSetParams] holds the parameters for updating documentSets in a collectionView.
//
// Fields:
//   - DocumentSetId: The list of documentSet's ids to update. The maximum number of elements in the array is 20.
//   - DocumentSetName: The list of documentSet's names to update. The maximum number of elements in the array is 20.
//   - Filter: (Optional) Filter documentSets by [Filter] conditions to update.
type UpdateAIDocumentSetParams struct {
	DocumentSetId   []string `json:"documentSetId"`
	DocumentSetName []string `json:"documentSetName"`
	Filter          *Filter  `json:"filter"`
}

type UpdateAIDocumentSetResult struct {
	AffectedCount uint64 `json:"affectedCount"`
}

// [Update] updates some documentSets in the collectionView.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - updateFields: The fields with which you want to update the documentSet.
//   - param: A [UpdateAIDocumentSetParams] object that includes the other parameters for updating documentSets' operation.
//     See [UpdateAIDocumentSetParams] for more information.
//
// Notes: The name of the database and the name of collectionView are from the fields of [implementerAIDocumentSets].
//
// Returns a pointer to a [UpdateAIDocumentSetResult] object or an error.
func (i *implementerAIDocumentSets) Update(ctx context.Context, updateFields map[string]interface{},
	param UpdateAIDocumentSetParams) (*UpdateAIDocumentSetResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(ai_document_set.UpdateReq)
	res := new(ai_document_set.UpdateRes)

	req.Database = i.database.DatabaseName
	req.CollectionView = i.collectionView.connCollectionViewName

	req.Query = ai_document_set.UpdateQueryCond{
		DocumentSetId:   param.DocumentSetId,
		DocumentSetName: param.DocumentSetName,
		Filter:          param.Filter.Cond(),
	}
	req.Update = updateFields

	result := new(UpdateAIDocumentSetResult)
	err := i.Request(ctx, req, res)
	if err != nil {
		return result, err
	}
	result.AffectedCount = res.AffectedCount
	return result, nil
}

// [GetCosTmpSecretParams] holds the parameters for getting the temp secret to upload specific file.
//
// Fields:
//   - DocumentSetName: The name of the documentSet(file).
//   - ParsingProcess: A pointer to a [ParsingProcess] object, which includes the parameters
//     for parsing files. See [ParsingProcess] for more information.
//   - ByteLength: (Optional) The size of the file to be uploaded in bytes.
type GetCosTmpSecretParams struct {
	DocumentSetName string              `json:"documentSetName"`
	ParsingProcess  *api.ParsingProcess `json:"parsingProcess,omitempty"`
	ByteLength      *uint64             `json:"byteLength,omitempty"`
}

type GetCosTmpSecretResult struct {
	DocumentSetId           string `json:"documentSetId"`
	DocumentSetName         string `json:"documentSetName"`
	CosEndpoint             string `json:"cosEndpoint"`
	CosRegion               string `json:"cosRegion"`
	CosBucket               string `json:"cosBucket"`
	UploadPath              string `json:"uploadPath"`
	TmpSecretID             string `json:"tmpSecretId"`
	TmpSecretKey            string `json:"tmpSecretKey"`
	SessionToken            string `json:"token"`
	Expiration              string `json:"Expiration,omitempty"`
	ExpiredTime             int    `json:"ExpiredTime,omitempty"`
	MaxSupportContentLength int64  `json:"maxSupportContentLength"`
}

// [GetCosTmpSecret] gets the temp secret to upload specific file.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - param: A [GetCosTmpSecretParams] object that includes the other parameters for getting the
//     temp secret to upload specific file. See [GetCosTmpSecretParams] for more information.
//
// Notes: The name of the database and the name of collectionView are from the fields of [implementerAIDocumentSets].
//
// Returns a pointer to a [GetCosTmpSecretResult] object or an error.
func (i *implementerAIDocumentSets) GetCosTmpSecret(ctx context.Context, param GetCosTmpSecretParams) (*GetCosTmpSecretResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}

	req := new(ai_document_set.UploadUrlReq)
	res := new(ai_document_set.UploadUrlRes)

	req.Database = i.database.DatabaseName
	req.CollectionView = i.collectionView.connCollectionViewName
	req.DocumentSetName = param.DocumentSetName
	req.ParsingProcess = param.ParsingProcess
	req.ByteLength = param.ByteLength

	err := i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}

	if res.UploadCondition == nil || res.Credentials == nil {
		return nil, fmt.Errorf("get file upload url failed")
	}
	result := new(GetCosTmpSecretResult)
	result.DocumentSetName = req.DocumentSetName
	result.DocumentSetId = res.DocumentSetId
	result.CosEndpoint = res.CosEndpoint
	result.CosBucket = res.CosBucket
	result.CosRegion = res.CosRegion
	result.UploadPath = res.UploadPath
	result.TmpSecretID = res.Credentials.TmpSecretID
	result.TmpSecretKey = res.Credentials.TmpSecretKey
	result.SessionToken = res.Credentials.SessionToken
	result.Expiration = res.Credentials.Expiration
	result.ExpiredTime = res.Credentials.ExpiredTime
	result.MaxSupportContentLength = res.UploadCondition.MaxSupportContentLength

	return result, nil
}

// [LoadAndSplitTextParams] holds the parameters for loading local file.
//
// Fields:
//   - DocumentSetName: The name of the documentSet.
//   - LocalFilePath: The file path of the locally uploaded file.
//   - MetaData: The configuration for the file metadata.
//   - SplitterPreprocess: A pointer to a [DocumentSplitterPreprocess] object, which includes the parameters
//     for splitting document chunks. See [DocumentSplitterPreprocess] for more information.
//   - ParsingProcess: A pointer to a [ParsingProcess] object, which includes the parameters
//     for parsing files. See [ParsingProcess] for more information.
type LoadAndSplitTextParams struct {
	DocumentSetName    string
	Reader             io.Reader
	LocalFilePath      string
	MetaData           map[string]interface{}
	SplitterPreprocess ai_document_set.DocumentSplitterPreprocess
	ParsingProcess     *api.ParsingProcess
}

type cosMetaConfig struct {
	AppendTitleToChunk    *bool               `json:"appendTitleToChunk,omitempty"`
	AppendKeywordsToChunk *bool               `json:"appendKeywordsToChunk,omitempty"`
	ChunkSplitter         *string             `json:"chunkSplitter,omitempty"`
	ParsingProcess        *api.ParsingProcess `json:"parsingProcess,omitempty"`
}

type LoadAndSplitTextResult struct {
	GetCosTmpSecretResult
}

// [LoadAndSplitText] uploads local file, which will be parsed and saved remotely.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - param: A [LoadAndSplitTextParams] object that includes the other parameters for uploading local file.
//     See [LoadAndSplitTextParams] for more information.
//
// Notes: The name of the database and the name of collectionView are from the fields of [implementerAIDocumentSets].
//
// Returns a pointer to a [LoadAndSplitTextResult] object or an error.
func (i *implementerAIDocumentSets) LoadAndSplitText(ctx context.Context, param LoadAndSplitTextParams) (result *LoadAndSplitTextResult, err error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}

	size, reader, err := i.loadAndSplitTextCheckParams(&param)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	byteLength := uint64(size)
	res, err := i.GetCosTmpSecret(ctx, GetCosTmpSecretParams{
		DocumentSetName: param.DocumentSetName,
		ParsingProcess:  param.ParsingProcess,
		ByteLength:      &byteLength,
	})
	if err != nil {
		return nil, err
	}

	if size > res.MaxSupportContentLength {
		return nil, fmt.Errorf("fileSize is invalid, support max content length is %v bytes", res.MaxSupportContentLength)
	}

	u, _ := url.Parse(res.CosEndpoint)
	b := &cos.BaseURL{BucketURL: u}

	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     res.TmpSecretID,  // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/ai_document_set/product/598/37140
			SecretKey:    res.TmpSecretKey, // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/ai_document_set/product/598/37140
			SessionToken: res.SessionToken,
		},
	})

	header := make(http.Header)

	marshalData, err := json.Marshal(param.MetaData)
	if err != nil {
		return nil, fmt.Errorf("put param MetaData into cos header failed, err: %v", err.Error())
	}

	cosMetaConfig := cosMetaConfig{
		AppendTitleToChunk:    param.SplitterPreprocess.AppendTitleToChunk,
		AppendKeywordsToChunk: param.SplitterPreprocess.AppendKeywordsToChunk,
		ChunkSplitter:         param.SplitterPreprocess.ChunkSplitter,
		ParsingProcess:        param.ParsingProcess,
	}

	configMarshalData, err := json.Marshal(cosMetaConfig)
	if err != nil {
		return nil, fmt.Errorf("put param SplitterPreprocess into cos header failed, err: %v", err.Error())
	}

	header.Add("x-cos-meta-data", url.QueryEscape(base64.StdEncoding.EncodeToString(marshalData)))
	header.Add("x-cos-meta-id", res.DocumentSetId)
	header.Add("x-cos-meta-config", url.QueryEscape(base64.StdEncoding.EncodeToString(configMarshalData)))

	headerData, err := json.Marshal(header)
	if err != nil {
		return nil, fmt.Errorf("marshal cos header failed, err: %v", err.Error())
	}
	if len(headerData) > 2048 {
		return nil, fmt.Errorf("cos header for param MetaData is too large, it can not be more than 2k")
	}

	if param.LocalFilePath != "" {
		// upload file by reading local file path, which supports multi parts uploading
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
		// upload file by io.reader
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

	result = new(LoadAndSplitTextResult)
	result.GetCosTmpSecretResult = *res
	return result, nil
}

func (i *implementerAIDocumentSets) loadAndSplitTextCheckParams(param *LoadAndSplitTextParams) (size int64, reader io.ReadCloser, err error) {
	if param.DocumentSetName == "" {
		if param.LocalFilePath == "" {
			return 0, nil, errors.New("need param: DocumentSetName")
		}
		param.DocumentSetName = filepath.Base(param.LocalFilePath)
	}
	fileType := strings.ToLower(filepath.Ext(param.DocumentSetName))
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

func (i *implementerAIDocumentSets) toDocumentSet(item ai_document_set.QueryDocumentSet) *AIDocumentSet {
	documentSet := new(AIDocumentSet)
	documentSet.DatabaseName = i.database.DatabaseName
	documentSet.CollectionViewName = i.collectionView.CollectionViewName
	documentSet.DocumentSetId = item.DocumentSetId
	documentSet.DocumentSetName = item.DocumentSetName
	if item.Text != nil {
		documentSet.Text = *item.Text
	}
	if item.TextPrefix != nil {
		documentSet.TextPrefix = *item.TextPrefix
	}
	documentSet.DocumentSetInfo = item.DocumentSetInfo
	documentSet.ScalarFields = make(map[string]Field)
	for k, v := range item.ScalarFields {
		documentSet.ScalarFields[k] = Field{Val: v}
	}

	docSetImpl := new(implementerAIDocumentSet)
	docSetImpl.SdkClient = i.SdkClient
	docSetImpl.database = i.database
	docSetImpl.collectionView = i.collectionView
	docSetImpl.documentSet = documentSet

	documentSet.AIDocumentSetInterface = docSetImpl
	documentSet.SplitterPreprocess = item.SplitterPreprocess
	documentSet.ParsingProcess = item.ParsingProcess
	return documentSet
}

type AISearchDocumentSet struct {
	DatabaseName       string
	CollectionViewName string
	DocumentSetId      string                     `json:"documentSetId"`
	DocumentSetName    string                     `json:"documentSetName"`
	Score              float64                    `json:"score"`
	SearchData         ai_document_set.SearchData `json:"data"`
	ScalarFields       map[string]Field
}

func (i *implementerAIDocumentSets) toSearchDocumentSet(item ai_document_set.SearchDocument) *AISearchDocumentSet {
	documentSet := new(AISearchDocumentSet)
	documentSet.DatabaseName = i.database.DatabaseName
	documentSet.CollectionViewName = i.collectionView.CollectionViewName
	documentSet.DocumentSetId = item.DocumentSet.DocumentSetId
	documentSet.DocumentSetName = item.DocumentSet.DocumentSetName
	documentSet.Score = item.Score
	documentSet.SearchData = item.Data
	documentSet.ScalarFields = make(map[string]Field)
	for k, v := range item.DocumentSet.ScalarFields {
		documentSet.ScalarFields[k] = Field{Val: v}
	}
	return documentSet
}
