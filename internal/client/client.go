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

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/internal/engine/api"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/model"

	"github.com/pkg/errors"
)

type Client struct {
	// cli      *gclient.Client
	cli      *http.Client
	url      string
	username string
	key      string
	timeout  time.Duration
	debug    bool
}

var defaultOption = model.ClientOption{
	Timeout:            time.Second * 5,
	MaxIdldConnPerHost: 10,
	IdleConnTimeout:    time.Minute,
}

// NewClient new http client with url, username and api key
func NewClient(url, username, key string, options *model.ClientOption) (model.SdkClient, error) {
	if !strings.HasPrefix(url, "http") {
		return nil, errors.Errorf("invailid url param with: %s", url)
	}
	if username == "" || key == "" {
		return nil, errors.New("username or key is empty")
	}
	if options == nil {
		options = &defaultOption
	}
	if options.Timeout == 0 {
		options.Timeout = defaultOption.Timeout
	}
	if options.IdleConnTimeout == 0 {
		options.IdleConnTimeout = defaultOption.IdleConnTimeout
	}
	if options.MaxIdldConnPerHost == 0 {
		options.MaxIdldConnPerHost = defaultOption.MaxIdldConnPerHost
	}

	cli := new(Client)
	cli.url = url
	cli.username = username
	cli.key = key
	cli.debug = false

	cli.cli = new(http.Client)
	cli.cli.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConnsPerHost: options.MaxIdldConnPerHost,
		IdleConnTimeout:     options.IdleConnTimeout,
	}
	cli.cli.Timeout = options.Timeout

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
	response, err := c.cli.Do(request)
	// response, err := gclient.New().ContentJson().DoRequest(ctx, method, c.url+path, reqBody.String())
	if err != nil {
		return err
	}
	return c.handleResponse(ctx, response, res)
}

// WithTimeout set client timeout
func (c *Client) WithTimeout(d time.Duration) {
	c.timeout = d
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
	var commenRes model.CommmonResponse

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
