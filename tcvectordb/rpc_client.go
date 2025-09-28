package tcvectordb

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/olama"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type RpcClient struct {
	DatabaseInterface
	FlatInterface
	FlatIndexInterface

	httpImplementer SdkClient
	rpcClient       olama.SearchEngineClient
	cc              *grpc.ClientConn
	url             string
	username        string
	key             string
	option          ClientOption
	debug           bool
}

func NewRpcClient(url, username, key string, option *ClientOption) (*RpcClient, error) {
	if option == nil {
		option = &defaultOption
	}

	var httpTarget string
	var rpcTarget string
	if strings.HasPrefix(url, "http://") {
		httpTarget = url
		rpcTarget = strings.TrimPrefix(url, "http://")
		portIndex := strings.Index(rpcTarget, ":")
		if portIndex == -1 {
			rpcTarget += ":80"
		}
	} else if strings.HasPrefix(url, "https://") {
		return nil, errors.Errorf("invalid url param with %v for not supporting https://", url)
	} else {
		httpTarget = "http://" + url
		rpcTarget = url
	}

	cli := new(RpcClient)
	cli.url = url
	cli.username = username
	cli.key = key
	cli.debug = false
	cli.option = optionMerge(*option)

	cc, err := grpc.Dial(rpcTarget,
		grpc.WithUnaryInterceptor(newInterceptor(cli)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(100*1024*1024)),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(100*1024*1024)),
		grpc.WithInitialWindowSize(100*1024*1024),
		grpc.WithInitialConnWindowSize(100*1024*1024),
		grpc.WithBlock(),
	)
	cli.cc = cc
	if err != nil {
		return nil, err
	}
	cli.rpcClient = olama.NewSearchEngineClient(cc)

	httpc, err := NewClient(httpTarget, username, key, option)
	if err != nil {
		cc.Close()
		return nil, err
	}
	cli.httpImplementer = httpc

	databaseImpl := &rpcImplementerDatabase{
		cli,
		httpc.DatabaseInterface,
		cli.rpcClient,
	}
	flatImpl := &rpcImplementerFlatDocument{
		SdkClient: cli,
		rpcClient: cli.rpcClient,
	}
	flatIndexImpl := &rpcImplementerFlatIndex{
		SdkClient: cli,
		rpcClient: cli.rpcClient,
	}
	cli.DatabaseInterface = databaseImpl
	cli.FlatInterface = flatImpl
	cli.FlatIndexInterface = flatIndexImpl

	return cli, nil
}

func (r *RpcClient) Request(ctx context.Context, req, res interface{}) error {
	return r.httpImplementer.Request(ctx, req, res)
}

func (r *RpcClient) Options() ClientOption {
	return r.option
}

func (r *RpcClient) WithTimeout(d time.Duration) {
	r.httpImplementer.WithTimeout(d)
	r.option.Timeout = d
}

func (r *RpcClient) Debug(v bool) {
	r.httpImplementer.Debug(v)
	r.debug = v
}

func (r *RpcClient) Close() {
	r.httpImplementer.Close()
	r.cc.Close()
}

func (r *RpcClient) GetState() string {
	if r.cc == nil {
		return ""
	}
	return r.cc.GetState().String()
}

func (r *RpcClient) attachCtx(ctx context.Context) context.Context {
	auth := fmt.Sprintf("Bearer account=%s&api_key=%s", r.username, r.key)
	md := metadata.Pairs("authorization", auth)
	attached, _ := context.WithTimeout(ctx, r.option.Timeout)
	attached = metadata.NewOutgoingContext(attached, md)
	return attached
}

func newInterceptor(client *RpcClient) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = client.attachCtx(ctx)
		if client.debug {
			log.Printf("[DEBUG] REQUEST, Method: %s, Content: %v", method, req)
		}
		err := invoker(ctx, method, req, reply, cc, opts...)
		if codeGetter, ok := reply.(interface {
			GetCode() int32
			GetMsg() string
		}); ok {
			if codeGetter.GetCode() != 0 {
				err = errors.Errorf("code: %d, message: %s", codeGetter.GetCode(), codeGetter.GetMsg())
			}
		}
		if client.debug {
			if err != nil {
				log.Printf("[DEBUG] RESPONSE ERROR: %s", err.Error())
			} else {
				log.Printf("[DEBUG] RESPONSE: %v", reply)
			}
		}
		return err
	}
}
