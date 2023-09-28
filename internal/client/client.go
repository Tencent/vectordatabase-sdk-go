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

package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity/api"

	"github.com/pkg/errors"
)

type Client struct {
	// cli      *gclient.Client
	cli      *http.Client
	url      string
	username string
	key      string
	option   entity.ClientOption
	debug    bool
}

type CommmonResponse struct {
	// Code: 0 means success, other means failure.
	Code int32 `json:"code,omitempty"`
	// Msg: response msg
	Msg string `json:"msg,omitempty"`
}

var defaultOption = entity.ClientOption{
	Timeout:            time.Second * 5,
	MaxIdldConnPerHost: 2,
	IdleConnTimeout:    time.Minute,
	ReadConsistency:    entity.EventualConsistency,
}

// NewClient new http client with url, username and api key
func NewClient(url, username, key string, option *entity.ClientOption) (*Client, error) {
	if !strings.HasPrefix(url, "http") {
		return nil, errors.Errorf("invailid url param with: %s", url)
	}
	if username == "" || key == "" {
		return nil, errors.New("username or key is empty")
	}

	cli := new(Client)
	cli.url = url
	cli.username = username
	cli.key = key
	cli.debug = false

	cli.option = *optionMerge(option)

	cli.cli = new(http.Client)
	cli.cli.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConnsPerHost: cli.option.MaxIdldConnPerHost,
		IdleConnTimeout:     cli.option.IdleConnTimeout,
	}
	cli.cli.Timeout = cli.option.Timeout

	return cli, nil
}

// Request do request for client
func (c *Client) Request(ctx context.Context, req, res interface{}) error {
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

	request, err := http.NewRequest(strings.ToUpper(method), c.url+path, reqBody)
	if err != nil {
		return err
	}

	if c.debug {
		log.Printf("[DEBUG] REQUEST, Method: %s, Path: %s, Body: %s", method, path, strings.TrimSpace(reqBody.String()))
	}

	auth := fmt.Sprintf("Bearer account=%s&api_key=%s", c.username, c.key)
	request.Header.Add("Authorization", auth)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Sdk-Version", SDKVersion)
	response, err := c.cli.Do(request)
	if err != nil {
		return err
	}
	return c.handleResponse(ctx, response, res)
}

// WithTimeout set client timeout
func (c *Client) WithTimeout(d time.Duration) {
	c.option.Timeout = d
	c.cli.Timeout = d
}

// Debug set debug mode to show the request and response info
func (c *Client) Debug(v bool) {
	c.debug = v
}

func (c *Client) handleResponse(ctx context.Context, res *http.Response, out interface{}) error {
	responseBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if c.debug {
		log.Printf("[DEBUG] RESPONSE: %d %s", res.StatusCode, string(responseBytes))
	}
	if res.StatusCode/100 != 2 {
		return errors.Errorf("response code is %d, %s", res.StatusCode, string(responseBytes))
	}

	if !json.Valid(responseBytes) {
		return errors.Errorf(`invalid response content: %s`, responseBytes)
	}
	var commenRes CommmonResponse

	if err := json.Unmarshal(responseBytes, &commenRes); err != nil {
		return errors.Wrapf(err, `json.Unmarshal failed with content:%s`, responseBytes)
	}
	if commenRes.Code != 0 {
		return errors.Errorf("server internal error, code: %d, message: %s", commenRes.Code, commenRes.Msg)
	}

	if err := json.Unmarshal(responseBytes, &out); err != nil {
		return errors.Wrapf(err, `json.Unmarshal failed with content:%s`, responseBytes)
	}
	return nil
}

// Close wrap http.Client.CloseIdleConnections
func (c *Client) Close() {
	c.cli.CloseIdleConnections()
}

func (c *Client) Options() entity.ClientOption {
	return c.option
}

func optionMerge(option *entity.ClientOption) *entity.ClientOption {
	if option == nil {
		option = &defaultOption
	}
	if option.Timeout == 0 {
		option.Timeout = defaultOption.Timeout
	}
	if option.IdleConnTimeout == 0 {
		option.IdleConnTimeout = defaultOption.IdleConnTimeout
	}
	if option.MaxIdldConnPerHost == 0 {
		option.MaxIdldConnPerHost = defaultOption.MaxIdldConnPerHost
	}
	if option.ReadConsistency == "" {
		option.ReadConsistency = defaultOption.ReadConsistency
	}
	return option
}
