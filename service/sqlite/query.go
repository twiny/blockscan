package sqlite

// // Rest \\ \\
const selectLatestBlock = `
SELECT
	b1.block_number,
	b1.block_hash,
	b1.mined_timestamp,
	b1.tx_count
FROM
	blocks b1
ORDER BY
	b1.mined_timestamp DESC
LIMIT 1
`

const selectBlock = `
SELECT
	b1.block_number,
	b1.block_hash,
	b1.mined_timestamp,
	b1.tx_count
FROM
	blocks b1
WHERE
	b1.block_number = ?
`

const selectAllTxsHashsByBlockID = `
SELECT
	t1.tx_hash
FROM
	transactions t1
WHERE
	t1.block_number = ?
`

const selectLatestTx = `
SELECT
	t1.tx_hash,
	t1.block_number,
	t1.tx_from,
	t1.tx_to,
	t1.amount,
	t1.nonce,
	t1.mined_timestamp,
	t1.tx_order
FROM
	transactions t1
WHERE
	t1.block_number = (
		SELECT
			b1.block_number
		FROM
			blocks b1
		ORDER BY
			b1.mined_timestamp DESC
		LIMIT 1)
ORDER BY
	t1.tx_order DESC
LIMIT 1
`

const selectTx = `
SELECT
	t1.tx_hash,
	t1.block_number,
	t1.tx_from,
	t1.tx_to,
	t1.amount,
	t1.nonce,
	t1.mined_timestamp,
	t1.tx_order
FROM
	transactions t1
WHERE
	t1.tx_hash = ?
`

const selectSumOfAllTx = `
SELECT
	TOTAL (t1.amount)
FROM
	transactions t1
WHERE (t1.block_number BETWEEN ? AND ?);
`

const selectAllTxHash = `
SELECT
	t1.tx_hash
FROM
	transactions t1
	WHERE (t1.block_number BETWEEN ? AND ?)
`

// // Indexer \\ \\

const hasScanned = `
SELECT EXISTS(SELECT 1 FROM blocks b1 WHERE b1.block_number = ?) AS found;
`

const insertBlock = `
INSERT INTO "blocks"
	(block_number, block_hash, mined_timestamp, tx_count)
VALUES 
	(?,?,?,?);
`

const insertTx = `
INSERT INTO "transactions"
	(tx_hash, block_number, tx_from, tx_to, amount, nonce, mined_timestamp, tx_order)
VALUES 
	(?,?,?,?,?,?,?,?);
`
