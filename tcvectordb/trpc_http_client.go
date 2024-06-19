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
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"trpc.group/trpc-go/trpc-go/client"
	"trpc.group/trpc-go/trpc-go/codec"
	thttp "trpc.group/trpc-go/trpc-go/http"

	"github.com/pkg/errors"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api"
)

type TrpcHttpClient struct {
	DatabaseInterface
	cli      thttp.Client
	url      string
	username string
	key      string
	option   ClientOption
	debug    bool
}

func NewTrpcHttpClient(url, username, key string, option *ClientOption) (*TrpcHttpClient, error) {
	if option == nil {
		option = &defaultOption
	}
	return newTrpcHttpClient(url, username, key, optionMerge(*option))
}

// newClient new http client with url, username and api key
func newTrpcHttpClient(url, username, key string, option ClientOption) (*TrpcHttpClient, error) {
	if !strings.HasPrefix(url, "http://") {
		return nil, errors.Errorf("invalid url param with: %s", url)
	}
	if username == "" || key == "" {
		return nil, errors.New("username or key is empty")
	}

	cli := new(TrpcHttpClient)
	cli.url = url
	cli.username = username
	cli.key = key
	cli.debug = false

	cli.option = optionMerge(option)

	cli.cli = thttp.NewClientProxy(url[7:],
		client.WithProtocol("http"),
		client.WithSerializationType(codec.SerializationTypeJSON),
		client.WithTimeout(cli.option.Timeout),
	)

	databaseImpl := new(implementerDatabase)
	databaseImpl.SdkClient = cli

	cli.DatabaseInterface = databaseImpl
	return cli, nil
}

func (c *TrpcHttpClient) Request(ctx context.Context, req, res interface{}) error {
	var (
		method = api.Method(req)
		path   = api.Path(req)
	)

	reqBody := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(reqBody)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(req)
	if err != nil {
		return fmt.Errorf("%w, %#v", err, req)
	}

	if c.debug {
		log.Printf("[DEBUG] REQUEST, Method: %s, Path: %s, Body: %s", method, path, strings.TrimSpace(reqBody.String()))
	}

	auth := fmt.Sprintf("Bearer account=%s&api_key=%s", c.username, c.key)

	reqHeader := &thttp.ClientReqHeader{}
	// 必须设置正确的 Method
	reqHeader.Method = strings.ToUpper(method)
	// 为 HTTP Head 添加 request 字段
	reqHeader.AddHeader("Authorization", auth)
	reqHeader.AddHeader("Content-Type", "application/json")
	reqHeader.AddHeader("Sdk-Version", SDKVersion)

	if reqHeader.Method == "POST" {
		err = c.cli.Post(ctx, path, req, res, client.WithReqHead(reqHeader))
	} else if reqHeader.Method == "GET" {
		err = c.cli.Get(ctx, path, res, client.WithReqHead(reqHeader))
	}
	if err != nil {
		return err
	}

	return nil
}

func (c *TrpcHttpClient) WithTimeout(d time.Duration) {
	c.option.Timeout = d
}

// Debug set debug mode to show the request and response info
func (c *TrpcHttpClient) Debug(v bool) {
	c.debug = v
}

func (c *TrpcHttpClient) Close() {
}

func (c *TrpcHttpClient) Options() ClientOption {
	return c.option
}
