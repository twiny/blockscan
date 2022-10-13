package chain

import "time"

// Block
type Block struct {
	Number    int64     `json:"number"`
	Hash      string    `json:"hash"`
	Timestamp time.Time `json:"timestamp"` // timestamp when the block was mined
	TxCount   uint      `json:"tx_count"`
	Txs       []string  `json:"txs"` // array of transactions hash
}
