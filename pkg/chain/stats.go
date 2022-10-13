package chain

// Stats
type Stats struct {
	Txs         []string `json:"txs"` // array of transactions hash
	TotalAmount float64  `json:"total_amount"`
}
