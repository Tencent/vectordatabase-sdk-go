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

// [SetAlias] sets an alias for collection.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - collectionName: The name of the collection.
//   - aliasName: The alias name to set for the collection.
//
// Returns a pointer to a [SetAliasResult] object or an error.
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

// [DeleteAlias] deletes the alias in the database.
//
// Parameters:
//   - ctx: A context.Context object controls the request's lifetime, allowing for the request
//     to be canceled or to timeout according to the context's deadline.
//   - aliasName: The alias name to delete.
//
// Returns a pointer to a [DeleteAliasResult] object or an error.
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
