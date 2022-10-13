package api

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/twiny/blockscan/pkg/chain"
	"github.com/twiny/blockscan/pkg/utils"

	"github.com/ethereum/go-ethereum/core/types"
)

// indexer
func (idx *Indexer) indexer() {
	idx.wg.Add(idx.conf.Indexer.Workers)

	for i := 0; i < idx.conf.Indexer.Workers; i++ {
		go func() {
			defer idx.wg.Done()

			for {
				select {
				case <-idx.ctx.Done():
					return
				case i := <-idx.jobs:
					// scan
					if err := idx.scan(i); err != nil {
						idx.log.Println(err)
						continue
					}

					// check if newer than current head block
					if !idx.subscribed && i >= idx.head {
						// subcribe
						idx.subscribe()

						idx.head = i

						idx.subscribed = true
					}

					// log
					idx.log.Printf("scanned block id %d", i)
				}
			}
		}()
	}
}

// Subscribe
func (idx *Indexer) subscribe() {
	sub, err := idx.client.SubscribeNewHead(context.Background(), idx.events)
	if err != nil {
		idx.log.Println(err)

		idx.subscribed = false
		return
	}

	idx.wg.Add(1)
	go func() {
		defer idx.wg.Done()

		//  log
		idx.log.Println("subscribed to new block")
		for {
			select {
			case err := <-sub.Err():
				idx.log.Println(err)
				idx.subscribed = false
				return
			case <-idx.ctx.Done():
				return
			case header := <-idx.events:
				head := header.Number.Int64()
				idx.jobs <- head

				// update head
				idx.head = head
			}
		}
	}()
}

// scan
func (idx *Indexer) scan(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), idx.conf.Indexer.Timeout)
	defer cancel()

	// already scanned
	if idx.store.HasScanned(ctx, id) {
		return fmt.Errorf("block %d already scanned", id)
	}

	// rate limit
	idx.limiter.Take()

	return idx.storeBlock(ctx, id)
}

// storeBlock
func (idx *Indexer) storeBlock(ctx context.Context, id int64) error {
	// get block
	block, err := idx.client.BlockByNumber(ctx, big.NewInt(id))
	if err != nil {
		return err
	}

	hash := block.Hash()

	txCount, err := idx.client.TransactionCount(ctx, hash)
	if err != nil {
		return err
	}

	b := &chain.Block{
		Number:    id,
		Hash:      hash.Hex(),
		Timestamp: time.Unix(int64(block.Time()), 0),
		TxCount:   txCount,
	}

	if err := idx.store.SaveBlock(ctx, b); err != nil {
		return err
	}

	//
	// get chain id
	chainid, err := idx.client.ChainID(ctx)
	if err != nil {
		return err
	}

	for order, tx := range block.Transactions() {
		msg, err := tx.AsMessage(types.LatestSignerForChainID(chainid), nil)
		if err != nil {
			return err
		}

		t := &chain.Tx{
			BlockNumber: block.Number().Int64(),
			Hash:        tx.Hash().Hex(),
			From:        msg.From().Hex(),
			To:          utils.AddrToHex(msg.To()), // case *common.Address == nil
			Amount:      tx.Value().Int64(),
			Nonce:       tx.Nonce(),
			Timestamp:   time.Unix(int64(block.Time()), 0), // Tx timestamp is same as blocl timestamp
			Order:       order,
		}

		if err := idx.store.SaveTx(ctx, t); err != nil {
			idx.log.Println("idx_store_save_tx", err)
			continue
		}
	}

	return nil
}
