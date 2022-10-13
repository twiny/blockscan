package sqlite

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"path"

	"github.com/twiny/blockscan/pkg/chain"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migration/schema.up.sql
var schemaUp string

//go:embed migration/schema.down.sql
var schemaDown string

// SQLite
type SQLite struct {
	db *sql.DB
}

// NewSQLiteDB: example: ./tmp/
func NewSQLiteDB(dir string) (*SQLite, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}

	storefile := path.Join(dir, "store.db")

	// create file if not exist
	if _, err := os.Stat(storefile); os.IsNotExist(err) {
		f, err := os.Create(storefile)
		if err != nil {
			return nil, err
		}
		f.Close()
	}

	db, err := sql.Open("sqlite3", storefile+"?cache=shared_sync=1&_cache_size=25000")
	if err != nil {
		return nil, err
	}

	// check connectivity
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &SQLite{
		db: db,
	}, nil
}

// Migrate
func (s *SQLite) Migrate(cmd string) error {
	switch cmd {
	case "up":
		_, err := s.db.ExecContext(context.Background(), schemaUp)
		return err
	case "down":
		_, err := s.db.ExecContext(context.Background(), schemaDown)
		return err
	default:
		return fmt.Errorf("unknown command")
	}
}

// Ping
func (s *SQLite) Ping() error {
	return s.db.Ping()
}

// Rest service \\

// GetLatestBlock
func (s *SQLite) GetLatestBlock(ctx context.Context) (*chain.Block, error) {
	// get latest block
	var b chain.Block
	if err := s.db.QueryRowContext(
		ctx,
		selectLatestBlock,
	).Scan(
		&b.Number,
		&b.Hash,
		&b.Timestamp,
		&b.TxCount,
	); err != nil {
		return nil, err
	}

	// get block transaction
	rows, err := s.db.QueryContext(
		ctx,
		selectAllTxsHashsByBlockID,
		b.Number,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs = []string{}

	for rows.Next() {
		var tx string
		if err := rows.Scan(&tx); err != nil {
			return nil, err
		}

		txs = append(txs, tx)
	}

	b.Txs = txs

	return &b, nil
}

// GetBlock
func (s *SQLite) GetBlock(ctx context.Context, n int64) (*chain.Block, error) {
	// get block
	var b chain.Block
	if err := s.db.QueryRowContext(
		ctx,
		selectBlock,
		n,
	).Scan(
		&b.Number,
		&b.Hash,
		&b.Timestamp,
		&b.TxCount,
	); err != nil {
		return nil, err
	}

	// get block transaction
	rows, err := s.db.QueryContext(
		ctx,
		selectAllTxsHashsByBlockID,
		b.Number,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs = []string{}

	for rows.Next() {
		var tx string
		if err := rows.Scan(&tx); err != nil {
			return nil, err
		}

		txs = append(txs, tx)
	}

	b.Txs = txs

	return &b, nil
}

// GetLatestTx
func (s *SQLite) GetLatestTx(ctx context.Context) (*chain.Tx, error) {
	var t chain.Tx
	if err := s.db.QueryRowContext(
		ctx,
		selectLatestTx,
	).Scan(
		&t.Hash,
		&t.BlockNumber,
		&t.From,
		&t.To,
		&t.Amount,
		&t.Nonce,
		&t.Timestamp,
		&t.Order,
	); err != nil {
		return nil, err
	}
	return &t, nil
}

// GetTx
func (s *SQLite) GetTx(ctx context.Context, hash string) (*chain.Tx, error) {
	var t chain.Tx
	if err := s.db.QueryRowContext(
		ctx,
		selectTx,
		hash,
	).Scan(
		&t.Hash,
		&t.BlockNumber,
		&t.From,
		&t.To,
		&t.Amount,
		&t.Nonce,
		&t.Timestamp,
		&t.Order,
	); err != nil {
		return nil, err
	}
	return &t, nil
}

// GetRangeStats
func (s *SQLite) GetStats(ctx context.Context, i, j int64) (*chain.Stats, error) {
	var status = &chain.Stats{
		Txs:         []string{},
		TotalAmount: 0,
	}

	rows, err := s.db.QueryContext(ctx, selectAllTxHash, i, j)
	if err != nil {
		if err == sql.ErrNoRows {
			return status, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tx string
		if err := rows.Scan(&tx); err != nil {
			return nil, err
		}

		status.Txs = append(status.Txs, tx)
	}

	var total float64
	if err := s.db.QueryRowContext(ctx, selectSumOfAllTx, i, j).Scan(&total); err != nil {
		return nil, err
	}

	status.TotalAmount = total

	return status, nil
}

// Indexer service \\

// HasScanned
func (s *SQLite) HasScanned(ctx context.Context, id int64) bool {
	var found = 0
	if err := s.db.QueryRowContext(
		ctx,
		hasScanned,
		id,
	).Scan(&found); err != nil {
		return false
	}

	return found != 0
}

// SaveBlock
func (s *SQLite) SaveBlock(ctx context.Context, b *chain.Block) error {
	_, err := s.db.ExecContext(
		ctx,
		insertBlock,
		b.Number,
		b.Hash,
		b.Timestamp,
		b.TxCount,
	)
	return err
}

// SaveTx
func (s *SQLite) SaveTx(ctx context.Context, tx *chain.Tx) error {
	_, err := s.db.ExecContext(
		ctx,
		insertTx,
		tx.Hash,
		tx.BlockNumber,
		tx.From,
		tx.To,
		tx.Amount,
		tx.Nonce,
		tx.Timestamp,
		tx.Order,
	)

	return err
}

// // \\ \\
// Close
func (s *SQLite) Close() error {
	return s.db.Close()
}
