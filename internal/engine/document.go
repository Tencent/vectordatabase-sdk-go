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
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/tencentyun/cos-go-sdk-v5"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity/api/document"
)

var _ entity.DocumentInterface = &implementerDocument{}

type implementerDocument struct {
	entity.SdkClient
	databaseName   string
	collectionName string
}

// Upsert upsert documents into collection. Support for repeated insertion
func (i *implementerDocument) Upsert(ctx context.Context, documents []entity.Document, option *entity.UpsertDocumentOption) (result *entity.DocumentResult, err error) {
	req := new(document.UpsertReq)
	req.Database = i.databaseName
	req.Collection = i.collectionName
	for _, doc := range documents {
		d := &document.Document{}
		d.Id = doc.Id
		d.Vector = doc.Vector
		d.Fields = make(map[string]interface{})
		for k, v := range doc.Fields {
			d.Fields[k] = v.Val
		}
		req.Documents = append(req.Documents, d)
	}

	if option != nil && option.BuildIndex != nil {
		req.BuildIndex = option.BuildIndex
	}

	res := new(document.UpsertRes)
	result = new(entity.DocumentResult)
	err = i.Request(ctx, req, res)
	if err != nil {
		return
	}
	result.AffectedCount = int(res.AffectedCount)
	return
}

// Query query the document by document ids.
// The parameters retrieveVector set true, will return the vector field, but will reduce the api speed.
func (i *implementerDocument) Query(ctx context.Context, documentIds []string, option *entity.QueryDocumentOption) ([]entity.Document, *entity.DocumentResult, error) {
	req := new(document.QueryReq)
	req.Database = i.databaseName
	req.Collection = i.collectionName
	req.Query = &document.QueryCond{
		DocumentIds: documentIds,
	}
	req.ReadConsistency = string(i.SdkClient.Options().ReadConsistency)
	if option != nil {
		req.Query.Filter = option.Filter.Cond()
		req.Query.RetrieveVector = option.RetrieveVector
		req.Query.OutputFields = option.OutputFields
		req.Query.Offset = option.Offset
		req.Query.Limit = option.Limit
	}

	res := new(document.QueryRes)
	result := new(entity.DocumentResult)
	err := i.Request(ctx, req, res)
	if err != nil {
		return nil, result, err
	}
	var documents []entity.Document
	for _, doc := range res.Documents {
		var d entity.Document
		d.Id = doc.Id
		d.Vector = doc.Vector
		d.Fields = make(map[string]entity.Field)

		for n, v := range doc.Fields {
			d.Fields[n] = entity.Field{Val: v}
		}
		documents = append(documents, d)
	}
	result.AffectedCount = len(documents)
	result.Total = int(res.Count)
	return documents, result, nil
}

// Search search document topK by vector. The optional parameters filter will add the filter condition to search.
// The optional parameters hnswParam only be set with the HNSW vector index type.
func (i *implementerDocument) Search(ctx context.Context, vectors [][]float32, option *entity.SearchDocumentOption) ([][]entity.Document, error) {
	return i.search(ctx, nil, vectors, nil, option)
}

// Search search document topK by document ids. The optional parameters filter will add the filter condition to search.
// The optional parameters hnswParam only be set with the HNSW vector index type.
func (i *implementerDocument) SearchById(ctx context.Context, documentIds []string, option *entity.SearchDocumentOption) ([][]entity.Document, error) {
	return i.search(ctx, documentIds, nil, nil, option)
}

func (i *implementerDocument) SearchByText(ctx context.Context, text map[string][]string, option *entity.SearchDocumentOption) ([][]entity.Document, error) {
	return i.search(ctx, nil, nil, text, option)
}

func (i *implementerDocument) search(ctx context.Context, documentIds []string, vectors [][]float32, text map[string][]string, option *entity.SearchDocumentOption) ([][]entity.Document, error) {
	req := new(document.SearchReq)
	req.Database = i.databaseName
	req.Collection = i.collectionName
	req.ReadConsistency = string(i.SdkClient.Options().ReadConsistency)
	req.Search = new(document.SearchCond)
	req.Search.DocumentIds = documentIds
	req.Search.Vectors = vectors
	for _, v := range text {
		req.Search.EmbeddingItems = v
	}

	if option != nil {
		req.Search.Filter = option.Filter.Cond()
		req.Search.RetrieveVector = option.RetrieveVector
		req.Search.OutputFields = option.OutputFields
		req.Search.Limit = option.Limit

		if option.Params != nil {
			req.Search.Params = new(document.SearchParams)
			req.Search.Params.Nprobe = option.Params.Nprobe
			req.Search.Params.Ef = option.Params.Ef
			req.Search.Params.Radius = option.Params.Radius
		}
	}

	res := new(document.SearchRes)
	err := i.Request(ctx, req, res)
	if err != nil {
		return nil, err
	}
	var documents [][]entity.Document
	for _, result := range res.Documents {
		var vecDoc []entity.Document
		for _, doc := range result {
			d := entity.Document{
				Id:     doc.Id,
				Vector: doc.Vector,
				Score:  doc.Score,
				Fields: make(map[string]entity.Field),
			}
			for n, v := range doc.Fields {
				d.Fields[n] = entity.Field{Val: v}
			}
			vecDoc = append(vecDoc, d)
		}
		documents = append(documents, vecDoc)
	}

	return documents, nil
}

// Delete delete document by document ids
func (i *implementerDocument) Delete(ctx context.Context, option *entity.DeleteDocumentOption) (result *entity.DocumentResult, err error) {
	req := new(document.DeleteReq)
	req.Database = i.databaseName
	req.Collection = i.collectionName
	if option != nil {
		req.Query = &document.QueryCond{
			DocumentIds: option.DocumentIds,
			Filter:      option.Filter.Cond(),
		}
	}

	res := new(document.DeleteRes)
	result = new(entity.DocumentResult)
	err = i.Request(ctx, req, res)
	if err != nil {
		return
	}
	result.AffectedCount = int(res.AffectedCount)
	return
}

func (i *implementerDocument) Update(ctx context.Context, option *entity.UpdateDocumentOption) (*entity.DocumentResult, error) {
	req := new(document.UpdateReq)
	req.Database = i.databaseName
	req.Collection = i.collectionName
	req.Query = new(document.QueryCond)

	if option != nil {
		req.Query.DocumentIds = option.QueryIds
		req.Query.Filter = option.QueryFilter.Cond()
		req.Update.Vector = option.UpdateVector
		req.Update.Fields = make(map[string]interface{})
		if len(option.UpdateFields) != 0 {
			for k, v := range option.UpdateFields {
				req.Update.Fields[k] = v.Val
			}
		}
	}

	res := new(document.UpdateRes)
	result := new(entity.DocumentResult)
	err := i.Request(ctx, req, res)
	if err != nil {
		return result, err
	}
	result.AffectedCount = int(res.AffectedCount)
	return result, nil
}

func (i *implementerDocument) Upload(ctx context.Context, localFilePath string, option *entity.UploadDocumentOption) (err error) {

	// localFilePath string, fileName string, fileType entity.FileType,
	//	metaData map[string]string

	cosUploadFileName := path.Base(localFilePath)

	fileType := getFileTypeFromFileName(localFilePath)
	if option != nil && option.FileType != "" {
		fileType = option.FileType
	}

	if fileType != entity.MarkdownFileType {
		return fmt.Errorf("only support markdown fileType when uploading")
	}

	req := new(document.UploadUrlReq)
	req.Database = i.databaseName
	req.Collection = i.collectionName
	req.FileName = cosUploadFileName
	req.FileType = fileType

	res := new(document.UploadUrlRes)
	err = i.Request(ctx, req, res)
	if err != nil {
		return err
	}
	fileSizeIsOk, err := checkFileSize(localFilePath, res.UploadCondition.MaxSupportContentLength)
	if err != nil {
		return err
	}
	if !fileSizeIsOk {
		return fmt.Errorf("%v fileSize is invalid, support max content length is %v bytes",
			localFilePath, res.UploadCondition.MaxSupportContentLength)
	}

	u, _ := url.Parse(res.CosEndpoint)
	b := &cos.BaseURL{BucketURL: u}

	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     res.Credentials.TmpSecretID,  // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/document/product/598/37140
			SecretKey:    res.Credentials.TmpSecretKey, // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/document/product/598/37140
			SessionToken: res.Credentials.SessionToken,
		},
	})

	header := make(http.Header)
	if option != nil && option.MetaData != nil {
		for k, v := range option.MetaData {
			header.Add("x-cos-meta-"+k, v)
		}
	}

	header.Add("x-cos-meta-fileType", string(fileType))
	header.Add("x-cos-meta-id", string(res.FileId))

	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			XCosMetaXXX: &header,
		},
	}

	_, err = c.Object.PutFromFile(ctx, res.UploadPath, localFilePath, opt)
	if err != nil {
		return err
	}
	return nil
}

func getFileTypeFromFileName(fileName string) entity.FileType {
	extension := filepath.Ext(fileName)
	extension = strings.ToLower(extension)
	// 不带后缀的文件，默认为markdown文件
	if extension == "" {
		return entity.MarkdownFileType
	} else if extension == ".md" || extension == ".markdown" {
		return entity.MarkdownFileType
	} else {
		return entity.UnSupportFileType
	}
}

func isMarkdownFile(localFilePath string) bool {
	extension := filepath.Ext(localFilePath)
	extension = strings.ToLower(extension)
	return extension == ".md" || extension == ".markdown"
}

func checkFileSize(localFilePath string, maxContentLength int64) (bool, error) {
	fileInfo, err := os.Stat(localFilePath)
	if err != nil {
		return false, err
	}

	if fileInfo.Size() <= maxContentLength {
		return true, nil
	}
	return false, nil
}
