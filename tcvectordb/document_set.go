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

var _ AIDocumentSetInterface = &implementerAIDocumentSet{}

type AIDocumentSetInterface interface {
	SdkClient
	Query(ctx context.Context, params ...*QueryAIDocumentSetParams) (*QueryAIDocumentSetResult, error)
	GetDocumentSetByName(ctx context.Context, documentSetName string) (*GetAIDocumentSetResult, error)
	GetDocumentSetById(ctx context.Context, documentSetId string) (*GetAIDocumentSetResult, error)
	Search(ctx context.Context, param SearchAIDocumentSetParams) (*SearchAIDocumentSetResult, error)
	DeleteByIds(ctx context.Context, documentSetIds ...string) (result *DeleteAIDocumentSetResult, err error)
	DeleteByNames(ctx context.Context, documentSetNames ...string) (result *DeleteAIDocumentSetResult, err error)
	Delete(ctx context.Context, param DeleteAIDocumentSetParams) (*DeleteAIDocumentSetResult, error)
	Update(ctx context.Context, updateFields map[string]interface{}, param UpdateAIDocumentSetParams) (*UpdateAIDocumentSetResult, error)
	LoadAndSplitText(ctx context.Context, param LoadAndSplitTextParams) (result *LoadAndSplitTextResult, err error)
	GetCosTmpSecret(ctx context.Context, param GetCosTmpSecretParams) (*GetCosTmpSecretResult, error)
}

type implementerAIDocumentSet struct {
	SdkClient
	database       AIDatabase
	collectionView CollectionView
}

// Query query the ai_document_set by ai_document_set ids.
// The parameters retrieveVector set true, will return the vector field, but will reduce the api speed.
func (i *implementerAIDocumentSet) Query(ctx context.Context, params ...*QueryAIDocumentSetParams) (*QueryAIDocumentSetResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(ai_document_set.QueryReq)
	res := new(ai_document_set.QueryRes)

	req.Database = i.database.DatabaseName
	req.CollectionView = i.collectionView.CollectionName
	if len(params) != 0 && params[0] != nil {
		param := params[0]

		req.Query = &ai_document_set.QueryCond{
			Filter: param.Filter.Cond(),
			Limit:  param.Limit,
			Offset: param.Offset,
		}
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

func (i *implementerAIDocumentSet) GetDocumentSetByName(ctx context.Context, documentSetName string) (*GetAIDocumentSetResult, error) {
	return i.get(ctx, GetAIDocumentSetParams{DocumentSetName: documentSetName})
}

func (i *implementerAIDocumentSet) GetDocumentSetById(ctx context.Context, documentSetId string) (*GetAIDocumentSetResult, error) {
	return i.get(ctx, GetAIDocumentSetParams{DocumentSetId: documentSetId})
}

func (i *implementerAIDocumentSet) get(ctx context.Context, param GetAIDocumentSetParams) (*GetAIDocumentSetResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(ai_document_set.GetReq)
	res := new(ai_document_set.GetRes)

	req.Database = i.database.DatabaseName
	req.CollectionView = i.collectionView.CollectionName
	req.DocumentSetId = param.DocumentSetId
	req.DocumentSetName = param.DocumentSetName

	result := new(GetAIDocumentSetResult)
	err := i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}
	result.Count = res.Count
	result.DocumentSets = *i.toDocumentSet(res.DocumentSets)
	return result, nil
}

// Search search ai_document_set topK by vector. The optional parameters filter will add the filter condition to search.
func (i *implementerAIDocumentSet) Search(ctx context.Context, param SearchAIDocumentSetParams) (*SearchAIDocumentSetResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(ai_document_set.SearchReq)
	res := new(ai_document_set.SearchRes)

	req.Database = i.database.DatabaseName
	req.CollectionView = i.collectionView.CollectionName
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

func (i *implementerAIDocumentSet) DeleteByIds(ctx context.Context, documentSetIds ...string) (result *DeleteAIDocumentSetResult, err error) {
	return i.Delete(ctx, DeleteAIDocumentSetParams{DocumentSetIds: documentSetIds})
}

func (i *implementerAIDocumentSet) DeleteByNames(ctx context.Context, documentSetNames ...string) (result *DeleteAIDocumentSetResult, err error) {
	return i.Delete(ctx, DeleteAIDocumentSetParams{DocumentSetNames: documentSetNames})
}

// Delete delete documentSet by documentSetId or documentSetName ids
func (i *implementerAIDocumentSet) Delete(ctx context.Context, param DeleteAIDocumentSetParams) (result *DeleteAIDocumentSetResult, err error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(ai_document_set.DeleteReq)
	res := new(ai_document_set.DeleteRes)

	req.Database = i.database.DatabaseName
	req.CollectionView = i.collectionView.CollectionName
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

func (i *implementerAIDocumentSet) Update(ctx context.Context, updateFields map[string]interface{}, param UpdateAIDocumentSetParams) (*UpdateAIDocumentSetResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(ai_document_set.UpdateReq)
	res := new(ai_document_set.UpdateRes)

	req.Database = i.database.DatabaseName
	req.CollectionView = i.collectionView.CollectionName

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

func (i *implementerAIDocumentSet) GetCosTmpSecret(ctx context.Context, param GetCosTmpSecretParams) (*GetCosTmpSecretResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}

	req := new(ai_document_set.UploadUrlReq)
	res := new(ai_document_set.UploadUrlRes)

	req.Database = i.database.DatabaseName
	req.CollectionView = i.collectionView.CollectionName
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

func (i *implementerAIDocumentSet) LoadAndSplitText(ctx context.Context, param LoadAndSplitTextParams) (result *LoadAndSplitTextResult, err error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}

	res, err := i.GetCosTmpSecret(ctx, GetCosTmpSecretParams{
		DocumentSetName: param.DocumentSetName,
	})
	if err != nil {
		return nil, err
	}

	size, err := i.loadAndSplitTextCheckParams(&param)
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
	if metaData == nil {
		metaData = make(map[string]Field)
	}
	metaData["-id"] = Field{Val: res.DocumentSetId}

	for k, v := range metaData {
		if v.Type() == "" {
			continue
		}
		header.Add("x-cos-meta-"+string(v.Type())+"-"+k, url.QueryEscape(v.String()))
	}

	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			XCosMetaXXX: &header,
		},
	}

	_, err = c.Object.PutFromFile(ctx, res.UploadPath, param.LocalFilePath, opt)
	if err != nil {
		return nil, err
	}
	result = new(LoadAndSplitTextResult)
	result.DocumentSetId = res.DocumentSetId
	result.DocumentSetName = res.DocumentSetName
	result.CosEndpoint = res.CosEndpoint
	result.CosRegion = res.CosRegion
	result.CosBucket = res.CosBucket
	result.UploadPath = res.UploadPath
	result.TmpSecretID = res.TmpSecretID
	result.TmpSecretKey = res.TmpSecretKey
	result.SessionToken = res.SessionToken
	result.MaxSupportContentLength = res.MaxSupportContentLength
	return result, nil
}

func (i *implementerAIDocumentSet) loadAndSplitTextCheckParams(param *LoadAndSplitTextParams) (size int64, err error) {
	if param.DocumentSetName == "" {
		if param.LocalFilePath == "" {
			return 0, errors.New("need param: DocumentSetName")
		}
		param.DocumentSetName = filepath.Base(param.LocalFilePath)
	}
	if param.LocalFilePath != "" {
		fileInfo, err := os.Stat(param.LocalFilePath)
		if err != nil {
			return 0, err
		}
		size = fileInfo.Size()

		fd, err := os.Open(param.LocalFilePath)
		if err != nil {
			return 0, err
		}
		param.Reader = fd
	} else {
		bytesBuf := bytes.NewBuffer(nil)
		written, err := io.Copy(bytesBuf, param.Reader)
		if err != nil {
			return 0, err
		}
		size = written
		param.Reader = io.NopCloser(bytesBuf)
	}
	if size == 0 {
		return 0, errors.New("file size cannot be 0")
	}

	return size, nil
}

func (i *implementerAIDocumentSet) toDocumentSet(item ai_document_set.QueryDocumentSet) *AIDocumentSet {

	documentSet := new(AIDocumentSet)

	documentSet.QueryDocumentSet = item
	return documentSet
}
