package tcvectordb

import (
	"context"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/olama"
)

var _ AliasInterface = &rpcImplementerAlias{}

type rpcImplementerAlias struct {
	SdkClient
	rpcClient olama.SearchEngineClient
	database  *Database
}

func (r *rpcImplementerAlias) SetAlias(ctx context.Context, collectionName, aliasName string) (*SetAliasResult, error) {
	if r.database.IsAIDatabase() {
		return nil, AIDbTypeError
	}
	req := &olama.AddAliasRequest{
		Database:   r.database.DatabaseName,
		Collection: collectionName,
		Alias:      aliasName,
	}
	res, err := r.rpcClient.SetAlias(ctx, req)
	if err != nil {
		return nil, err
	}
	return &SetAliasResult{AffectedCount: int(res.AffectedCount)}, nil
}

func (r *rpcImplementerAlias) DeleteAlias(ctx context.Context, aliasName string) (*DeleteAliasResult, error) {
	if r.database.IsAIDatabase() {
		return nil, AIDbTypeError
	}
	req := &olama.RemoveAliasRequest{
		Database: r.database.DatabaseName,
		Alias:    aliasName,
	}
	res, err := r.rpcClient.DeleteAlias(ctx, req)
	if err != nil {
		return nil, err
	}
	return &DeleteAliasResult{AffectedCount: int(res.AffectedCount)}, nil
}
