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
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api"
)

// SdkClient provides the operations of a client.
type SdkClient interface {
	Request(ctx context.Context, req, res interface{}) error
	Options() ClientOption
	WithTimeout(d time.Duration)
	Debug(v bool)
	Close()
}

// [ClientOption] holds the parameters for creating an [Client] to a vectordb instance.
//
// Fields:
//   - Timeout: (Optional) Timeout specifies a time limit for requests made by this Client (defaults to 5s).
//   - MaxIdldConnPerHost:  (Optional) MaxIdleConnsPerHost controls the maximum idle (keep-alive) connections
//     to keep per-host if non-zero (defaults to 2).
//   - IdleConnTimeout: (Optional) IdleConnTimeout is the maximum amount of time an idle (keep-alive) connection
//     will remain idle before closing itself (defaults to 60s). Zero means no limit.
//   - ReadConsistency: (Optional) ReadConsistency represents the consistency level for reads.
//     The default value is "eventualConsistency", but it can be set to "strongConsistency" as well.
//   - Transport: (Optional) Transport specifies the mechanism by which individual HTTP requests are made (defaults to http.Transport).
type ClientOption struct {
	Timeout            time.Duration
	MaxIdldConnPerHost int
	IdleConnTimeout    time.Duration
	ReadConsistency    ReadConsistency
	Transport          http.RoundTripper
}
type Client struct {
	DatabaseInterface
	FlatInterface
	// deprecated:
	FlatIndexInterface

	cli      *http.Client
	url      string
	username string
	key      string
	option   ClientOption
	debug    bool
}

type CommmonResponse struct {
	Code int32  `json:"code,omitempty"` // Code: 0 means success, others mean failure.
	Msg  string `json:"msg,omitempty"`  // Msg: response msg
}

// ClientOption defaultOptions is the default options of a client connected to remote vectordb instance.
var defaultOption = ClientOption{
	Timeout:            time.Second * 5,
	MaxIdldConnPerHost: 2,
	IdleConnTimeout:    time.Minute,
	ReadConsistency:    api.EventualConsistency,
}

// [NewClient] creates and initializes a new instance of [Client] with the given url, username, key and option.
//
// Parameters:
//   - url: The address of vectordb, supporting http only.
//   - username: The username of vectordb, supporting root only currently.
//   - key: The account api key of vectordb, which you can get from console.
//   - option: A [ClientOption] object that includes the configuration for the vectordb client. See
//     [ClientOption] for more information.
//
// Notes:
//   - It is important to handle the error returned by this function to ensure that the
//     vectordb client has been created successfully before attempting to make API calls.
//
// Returns a pointer to an initialized [Client] instance or an error.

func NewClient(url, username, key string, option *ClientOption) (*Client, error) {
	if option == nil {
		option = &defaultOption
	}
	return newClient(url, username, key, optionMerge(*option))
}

func newClient(url, username, key string, option ClientOption) (*Client, error) {
	if !strings.HasPrefix(url, "http") {
		return nil, errors.Errorf("invalid url param with: %s", url)
	}
	if username == "" || key == "" {
		return nil, errors.New("username or key is empty")
	}

	cli := new(Client)
	cli.url = url
	cli.username = username
	cli.key = key
	cli.debug = false

	cli.option = optionMerge(option)

	cli.cli = new(http.Client)
	if option.Transport != nil {
		cli.cli.Transport = option.Transport
	} else {
		cli.cli.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			MaxIdleConnsPerHost: cli.option.MaxIdldConnPerHost,
			IdleConnTimeout:     cli.option.IdleConnTimeout,
		}
	}
	cli.cli.Timeout = cli.option.Timeout

	databaseImpl := new(implementerDatabase)
	databaseImpl.SdkClient = cli
	flatImpl := new(implementerFlatDocument)
	flatImpl.SdkClient = cli
	flatIndexImpl := new(implementerFlatIndex)
	flatIndexImpl.SdkClient = cli

	cli.DatabaseInterface = databaseImpl
	cli.FlatInterface = flatImpl
	cli.FlatIndexInterface = flatIndexImpl
	return cli, nil
}

// Request does request for client.
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

// WithTimeout sets client timeout.
func (c *Client) WithTimeout(d time.Duration) {
	c.option.Timeout = d
	c.cli.Timeout = d
}

// Debug sets debug mode to show the request and response info.
func (c *Client) Debug(v bool) {
	c.debug = v
}

func (c *Client) handleResponse(ctx context.Context, res *http.Response, out interface{}) error {
	responseBytes, err := io.ReadAll(res.Body)
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
		return errors.Errorf("code: %d, message: %s", commenRes.Code, commenRes.Msg)
	}

	if err := json.Unmarshal(responseBytes, &out); err != nil {
		return errors.Wrapf(err, `json.Unmarshal failed with content:%s`, responseBytes)
	}
	return nil
}

// Close closes idle connnections, releasing any open resources.
func (c *Client) Close() {
	c.cli.CloseIdleConnections()
}

// Options returns the option for the client.
func (c *Client) Options() ClientOption {
	return c.option
}

func optionMerge(option ClientOption) ClientOption {
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
