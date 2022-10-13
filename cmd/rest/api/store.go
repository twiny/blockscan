package api

import (
	"context"

	"github.com/twiny/blockscan/pkg/chain"
)

// StoreReader
type StoreReader interface {
	Ping() error
	//
	GetLatestBlock(ctx context.Context) (*chain.Block, error)
	GetBlock(ctx context.Context, n int64) (*chain.Block, error)
	//
	GetLatestTx(ctx context.Context) (*chain.Tx, error)
	GetTx(ctx context.Context, hash string) (*chain.Tx, error)
	//
	GetStats(ctx context.Context, i, j int64) (*chain.Stats, error)
}
