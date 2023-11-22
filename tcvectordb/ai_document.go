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
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/tcvectordb/api/ai_document"
	"github.com/tencentyun/cos-go-sdk-v5"
)

var _ AIDocumentInterface = &implementerAIDocument{}

type implementerAIDocument struct {
	SdkClient
	database   AIDatabase
	collection AICollection
}

// Query query the ai_document by ai_document ids.
// The parameters retrieveVector set true, will return the vector field, but will reduce the api speed.
func (i *implementerAIDocument) Query(ctx context.Context, options ...*QueryAIDocumentOption) (*QueryAIDocumentsResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(ai_document.QueryReq)
	req.Database = i.database.DatabaseName
	req.CollectionView = i.collection.CollectionName
	if len(options) != 0 && options[0] != nil {
		option := options[0]
		filter := option.Filter
		if filter == nil {
			filter = NewFilter("")
		}
		if option.FileName != "" {
			filter.And(fmt.Sprintf(`_file_name="%s"`, option.FileName))
		}

		req.Query = &ai_document.QueryCond{
			Filter: filter.Cond(),
			Limit:  option.Limit,
			Offset: option.Offset,
		}
		if req.Query.Limit == 0 {
			req.Query.Limit = 1
		}
	}

	res := new(ai_document.QueryRes)
	result := new(QueryAIDocumentsResult)
	err := i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}
	result.AffectedCount = len(res.DocumentSets)
	result.Total = int(res.Count)
	result.Documents = res.DocumentSets
	return result, nil
}

func (i *implementerAIDocument) Get(ctx context.Context, options ...*GetAIDocumentOption) (*GetAIDocumentResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(ai_document.GetReq)
	res := new(ai_document.GetRes)

	req.Database = i.database.DatabaseName
	req.CollectionView = i.collection.CollectionName
	if len(options) != 0 && options[0] != nil {
		option := options[0]
		req.DocumentSetId = option.DocumentSetId
		req.DocumentSetName = option.DocumentSetName
	}

	result := new(GetAIDocumentResult)
	err := i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}
	result.Count = res.Count
	result.DocumentSets = res.DocumentSets
	return result, nil
}

// Search search ai_document topK by vector. The optional parameters filter will add the filter condition to search.
func (i *implementerAIDocument) Search(ctx context.Context, content string, options ...*SearchAIDocumentOption) (*SearchAIDocumentResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(ai_document.SearchReq)
	req.Database = i.database.DatabaseName
	req.CollectionView = i.collection.CollectionName
	req.ReadConsistency = string(i.SdkClient.Options().ReadConsistency)
	req.Search = new(ai_document.SearchCond)
	req.Search.Content = content

	if len(options) != 0 && options[0] != nil {
		option := options[0]
		filter := option.Filter
		if filter == nil {
			filter = NewFilter("")
		}

		if option.FileName != "" {
			filter.And(fmt.Sprintf(`_file_name="%s"`, option.FileName))
		}
		req.Search.Filter = filter.Cond()
		req.Search.Options = ai_document.SearchOption{
			ResultType:  option.ResultType,
			ChunkExpand: option.ChunkExpand,
			// MergeChunk:  option.MergeChunk,
			// Weights: ai_document.SearchOptionWeight{
			// 	ChunkSimilarity: option.Weights.ChunkSimilarity,
			// 	WordSimilarity:  option.Weights.WordSimilarity,
			// 	WordBm25:        option.Weights.WordBm25,
			// },
		}
		if option.RerankOption != nil {
			req.Search.Options.RerankOption = &ai_document.RerankOption{
				Enable:                option.RerankOption.Enable,
				ExpectRecallMultiples: option.RerankOption.ExpectRecallMultiples,
			}
		}
		req.Search.OutputFields = option.OutputFields
		req.Search.Limit = option.Limit
	}

	res := new(ai_document.SearchRes)
	err := i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}
	result := new(SearchAIDocumentResult)
	result.Documents = res.Documents
	return result, nil
}

// Delete delete ai_document by ai_document ids
func (i *implementerAIDocument) Delete(ctx context.Context, options ...*DeleteAIDocumentOption) (result *DeleteAIDocumentResult, err error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(ai_document.DeleteReq)
	req.Database = i.database.DatabaseName
	req.CollectionView = i.collection.CollectionName
	if len(options) != 0 && options[0] != nil {
		option := options[0]
		filter := option.Filter
		if filter == nil {
			filter = NewFilter("")
		}
		if option.FileName != "" {
			filter.And(fmt.Sprintf(`_file_name="%s"`, option.FileName))
		}
		req.Query = &ai_document.DeleteQueryCond{
			DocumentSetId: option.DocumentIds,
			// todo: documentName
			Filter: filter.Cond(),
		}
	}

	res := new(ai_document.DeleteRes)
	result = new(DeleteAIDocumentResult)
	err = i.Request(ctx, req, res)
	if err != nil {
		return
	}
	result.AffectedCount = res.AffectedCount
	return
}

func (i *implementerAIDocument) Update(ctx context.Context, options ...*UpdateAIDocumentOption) (*UpdateAIDocumentResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	req := new(ai_document.UpdateReq)
	req.Database = i.database.DatabaseName
	req.CollectionView = i.collection.CollectionName

	if len(options) != 0 && options[0] != nil {
		option := options[0]
		filter := option.QueryFilter
		if filter == nil {
			filter = NewFilter("")
		}
		if option.FileName != "" {
			filter = filter.And(fmt.Sprintf(`_file_name="%s"`, option.FileName))
		}
		req.Query = ai_document.UpdateQueryCond{
			DocumentSetId:   option.QueryIds,
			DocumentSetName: nil,
			Filter:          filter.Cond(),
		}
		req.Update = option.UpdateFields
	}

	res := new(ai_document.UpdateRes)
	result := new(UpdateAIDocumentResult)
	err := i.Request(ctx, req, res)
	if err != nil {
		return result, err
	}
	result.AffectedCount = res.AffectedCount
	return result, nil
}

func (i *implementerAIDocument) GetCosTmpSecret(ctx context.Context, localFilePath string, options ...*GetCosTmpSecretOption) (*GetCosTmpSecretResult, error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}
	fileType := getFileTypeFromFileName(localFilePath)
	if len(options) != 0 && options[0] != nil {
		option := options[0]
		if option.FileType != "" {
			fileType = option.FileType
		}
	}

	if fileType != MarkdownFileType {
		return nil, fmt.Errorf("only support markdown fileType when uploading")
	}

	req := new(ai_document.UploadUrlReq)
	req.Database = i.database.DatabaseName
	req.CollectionView = i.collection.CollectionName
	req.DocumentSetName = filepath.Base(localFilePath)

	res := new(ai_document.UploadUrlRes)
	err := i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}

	if res.UploadCondition == nil || res.Credentials == nil {
		return nil, fmt.Errorf("get file upload url failed")
	}
	result := new(GetCosTmpSecretResult)
	result.FileName = req.DocumentSetName
	result.FileId = res.DocumentSetId
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

func (i *implementerAIDocument) Upload(ctx context.Context, localFilePath string, options ...*UploadAIDocumentOption) (result *UploadAIDocumentResult, err error) {
	if !i.database.IsAIDatabase() {
		return nil, BaseDbTypeError
	}

	// localFilePath string, fileName string, fileType FileType,
	//	metaData map[string]string
	var option *UploadAIDocumentOption
	var metaData map[string]Field
	if len(options) != 0 && options[0] != nil {
		option = options[0]
		metaData = option.MetaData
	}

	res, err := i.GetCosTmpSecret(ctx, localFilePath, &GetCosTmpSecretOption{
		FileType: option.FileType,
	})
	if err != nil {
		return nil, err
	}

	fileSizeIsOk, err := checkFileSize(localFilePath, res.MaxSupportContentLength)
	if err != nil {
		return nil, err
	}
	if !fileSizeIsOk {
		return nil, fmt.Errorf("%v fileSize is invalid, support max content length is %v bytes",
			localFilePath, res.MaxSupportContentLength)
	}

	u, _ := url.Parse(res.CosEndpoint)
	b := &cos.BaseURL{BucketURL: u}

	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     res.TmpSecretID,  // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/ai_document/product/598/37140
			SecretKey:    res.TmpSecretKey, // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/ai_document/product/598/37140
			SessionToken: res.SessionToken,
		},
	})

	header := make(http.Header)
	if metaData == nil {
		metaData = make(map[string]Field)
	}
	// todo: fileid
	metaData["-id"] = Field{Val: res.FileId}
	metaData["-document-set-name"] = Field{Val: res.FileName}

	for k, v := range metaData {
		if v.Type() == "" {
			continue
		}
		header.Add("x-cos-meta-"+string(v.Type())+"-"+k, url.QueryEscape(v.String()))
	}
	fmt.Println(header)

	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			XCosMetaXXX: &header,
		},
	}

	_, err = c.Object.PutFromFile(ctx, res.UploadPath, localFilePath, opt)
	if err != nil {
		return nil, err
	}
	result = new(UploadAIDocumentResult)
	result.FileName = res.FileName
	result.FileId = res.FileId
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
