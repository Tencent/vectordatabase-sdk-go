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

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb/api/ai_document_set"
	"github.com/pkg/errors"
	"github.com/tencentyun/cos-go-sdk-v5"
)

var _ AIDocumentSetsInterface = &implementerAIDocumentSets{}

type AIDocumentSetsInterface interface {
	SdkClient
	Query(ctx context.Context, params QueryAIDocumentSetParams) (*QueryAIDocumentSetResult, error)
	GetDocumentSetByName(ctx context.Context, documentSetName string) (*GetAIDocumentSetResult, error)
	GetDocumentSetById(ctx context.Context, documentSetId string) (*GetAIDocumentSetResult, error)
	Search(ctx context.Context, param SearchAIDocumentSetsParams) (*SearchAIDocumentSetResult, error)
	DeleteByIds(ctx context.Context, documentSetIds ...string) (result *DeleteAIDocumentSetResult, err error)
	DeleteByNames(ctx context.Context, documentSetNames ...string) (result *DeleteAIDocumentSetResult, err error)
	Delete(ctx context.Context, param DeleteAIDocumentSetParams) (*DeleteAIDocumentSetResult, error)
	Update(ctx context.Context, updateFields map[string]interface{}, param UpdateAIDocumentSetParams) (*UpdateAIDocumentSetResult, error)
	LoadAndSplitText(ctx context.Context, param LoadAndSplitTextParams) (result *LoadAndSplitTextResult, err error)
	GetCosTmpSecret(ctx context.Context, param GetCosTmpSecretParams) (*GetCosTmpSecretResult, error)
}

type AIDocumentSet struct {
	AIDocumentSetInterface
	ai_document_set.QueryDocumentSet
	DatabaseName       string
	CollectionViewName string
}

type implementerAIDocumentSets struct {
	SdkClient
	database       *AIDatabase
	collectionView *AICollectionView
}

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

// Query query the ai_document_set by ai_document_set ids.
// The parameters retrieveVector set true, will return the vector field, but will reduce the api speed.
func (i *implementerAIDocumentSets) Query(ctx context.Context, param QueryAIDocumentSetParams) (*QueryAIDocumentSetResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(ai_document_set.QueryReq)
	res := new(ai_document_set.QueryRes)

	req.Database = i.database.DatabaseName
	req.CollectionView = i.collectionView.CollectionViewName
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

func (i *implementerAIDocumentSets) GetDocumentSetByName(ctx context.Context, documentSetName string) (*GetAIDocumentSetResult, error) {
	return i.get(ctx, GetAIDocumentSetParams{DocumentSetName: documentSetName})
}

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
	req.CollectionView = i.collectionView.CollectionViewName
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
	Documents []ai_document_set.SearchDocument `json:"documents"`
}

// Search search ai_document_set topK by vector. The optional parameters filter will add the filter condition to search.
func (i *implementerAIDocumentSets) Search(ctx context.Context, param SearchAIDocumentSetsParams) (*SearchAIDocumentSetResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(ai_document_set.SearchReq)
	res := new(ai_document_set.SearchRes)

	req.Database = i.database.DatabaseName
	req.CollectionView = i.collectionView.CollectionViewName
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
	result.Documents = res.Documents
	return result, nil
}

type DeleteAIDocumentSetParams struct {
	DocumentSetNames []string `json:"documentSetNames"`
	DocumentSetIds   []string `json:"documentSetIds"`
	Filter           *Filter  `json:"filter"`
}

type DeleteAIDocumentSetResult struct {
	AffectedCount uint64 `json:"affectedCount"`
}

func (i *implementerAIDocumentSets) DeleteByIds(ctx context.Context, documentSetIds ...string) (result *DeleteAIDocumentSetResult, err error) {
	return i.Delete(ctx, DeleteAIDocumentSetParams{DocumentSetIds: documentSetIds})
}

func (i *implementerAIDocumentSets) DeleteByNames(ctx context.Context, documentSetNames ...string) (result *DeleteAIDocumentSetResult, err error) {
	return i.Delete(ctx, DeleteAIDocumentSetParams{DocumentSetNames: documentSetNames})
}

// Delete delete documentSet by documentSetId or documentSetName ids
func (i *implementerAIDocumentSets) Delete(ctx context.Context, param DeleteAIDocumentSetParams) (result *DeleteAIDocumentSetResult, err error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(ai_document_set.DeleteReq)
	res := new(ai_document_set.DeleteRes)

	req.Database = i.database.DatabaseName
	req.CollectionView = i.collectionView.CollectionViewName
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

type UpdateAIDocumentSetParams struct {
	DocumentSetId   []string `json:"documentSetId"`
	DocumentSetName []string `json:"documentSetName"`
	Filter          *Filter  `json:"filter"`
}

type UpdateAIDocumentSetResult struct {
	AffectedCount uint64 `json:"affectedCount"`
}

func (i *implementerAIDocumentSets) Update(ctx context.Context, updateFields map[string]interface{}, param UpdateAIDocumentSetParams) (*UpdateAIDocumentSetResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(ai_document_set.UpdateReq)
	res := new(ai_document_set.UpdateRes)

	req.Database = i.database.DatabaseName
	req.CollectionView = i.collectionView.CollectionViewName

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

type GetCosTmpSecretParams struct {
	DocumentSetName string `json:"documentSetName"`
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
	MaxSupportContentLength int64  `json:"maxSupportContentLength"`
}

func (i *implementerAIDocumentSets) GetCosTmpSecret(ctx context.Context, param GetCosTmpSecretParams) (*GetCosTmpSecretResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}

	req := new(ai_document_set.UploadUrlReq)
	res := new(ai_document_set.UploadUrlRes)

	req.Database = i.database.DatabaseName
	req.CollectionView = i.collectionView.CollectionViewName
	req.DocumentSetName = param.DocumentSetName

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
	result.MaxSupportContentLength = res.UploadCondition.MaxSupportContentLength

	return result, nil
}

type LoadAndSplitTextParams struct {
	DocumentSetName string
	Reader          io.Reader
	LocalFilePath   string
	MetaData        map[string]interface{}
}

type LoadAndSplitTextResult struct {
	GetCosTmpSecretResult
}

func (i *implementerAIDocumentSets) LoadAndSplitText(ctx context.Context, param LoadAndSplitTextParams) (result *LoadAndSplitTextResult, err error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}

	size, reader, err := i.loadAndSplitTextCheckParams(&param)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	res, err := i.GetCosTmpSecret(ctx, GetCosTmpSecretParams{
		DocumentSetName: param.DocumentSetName,
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
	metaData := param.MetaData

	marshalData, err := json.Marshal(metaData)
	if err != nil {
		return nil, fmt.Errorf("put param MetaData into cos header failed, err: %v", err.Error())
	}

	header.Add("x-cos-meta-data", url.QueryEscape(base64.StdEncoding.EncodeToString(marshalData)))
	header.Add("x-cos-meta-id", res.DocumentSetId)

	headerData, err := json.Marshal(header)
	if err != nil {
		return nil, fmt.Errorf("marshal cos header failed, err: %v", err.Error())
	}
	if len(headerData) > 2048 {
		return nil, fmt.Errorf("cos header for param MetaData is too large, it can not be more than 2k")
	}

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
	documentSet.QueryDocumentSet = item
	documentSet.DatabaseName = i.database.DatabaseName
	documentSet.CollectionViewName = i.collectionView.CollectionViewName

	docSetImpl := new(implementerAIDocumentSet)
	docSetImpl.SdkClient = i.SdkClient
	docSetImpl.database = i.database
	docSetImpl.collectionView = i.collectionView
	docSetImpl.documentSet = documentSet

	documentSet.AIDocumentSetInterface = docSetImpl
	return documentSet
}
