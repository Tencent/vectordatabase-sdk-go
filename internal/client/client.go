package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/model"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/util/gmeta"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gtag"
)

type Client struct {
	cli      *gclient.Client
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

func NewClient(url, username, key string, options *model.ClientOption) (model.SdkClient, error) {
	if !strings.HasPrefix(url, "http") {
		return nil, gerror.Newf("invailid url param with: %s", url)
	}
	if username == "" || key == "" {
		return nil, gerror.New("username or key is empty")
	}
	if options == nil {
		options = &defaultOption
	}

	cli := new(Client)
	cli.url = url
	cli.username = username
	cli.key = key
	cli.debug = false
	auth := fmt.Sprintf("Bearer account=%s&api_key=%s", username, key)

	cli.cli = gclient.New()
	cli.cli.SetHeader("Authorization", auth)
	cli.cli.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConnsPerHost: options.MaxIdldConnPerHost,
		IdleConnTimeout:     options.IdleConnTimeout,
	}
	cli.cli.Timeout(options.Timeout)

	return cli, nil
}

func (c *Client) Request(ctx context.Context, req, res interface{}) error {
	var (
		method = gmeta.Get(req, gtag.Method).String()
		path   = gmeta.Get(req, gtag.Path).String()
	)
	reqBody := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(reqBody)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(req)
	if err != nil {
		return err
	}

	if c.debug {
		glog.Debugf(ctx, "REQUEST, Method: %s, Path: %s, Body: %s", method, path, strings.TrimSpace(reqBody.String()))
	}

	response, err := c.cli.ContentJson().DoRequest(ctx, method, c.url+path, reqBody.String())
	if err != nil {
		return err
	}
	return c.handleResponse(ctx, response, res)
}

func (c *Client) WithTimeout(d time.Duration) {
	c.timeout = d
	c.cli.Timeout(d)
}

func (c *Client) Debug(v bool) {
	c.debug = v
}

func (c *Client) handleResponse(ctx context.Context, res *gclient.Response, out interface{}) error {
	var (
		responseBytes = res.ReadAll()
	)
	if c.debug {
		glog.Debugf(ctx, "RESPONSE: %s", string(responseBytes))
	}

	if !json.Valid(responseBytes) {
		return gerror.Newf(`invalid response content: %s`, responseBytes)
	}
	var commenRes model.CommmonResponse

	if err := json.Unmarshal(responseBytes, &commenRes); err != nil {
		return gerror.Wrapf(err, `json.Unmarshal failed with content:%s`, responseBytes)
	}
	if commenRes.Code != 0 {
		return gerror.Newf("server internal error, code: %d, message: %s", commenRes.Code, commenRes.Msg)
	}

	if err := json.Unmarshal(responseBytes, &out); err != nil {
		return gerror.Wrapf(err, `json.Unmarshal failed with content:%s`, responseBytes)
	}
	return nil
}

func (c *Client) Close() {
	c.cli.CloseIdleConnections()
}
