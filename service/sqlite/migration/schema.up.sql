-- blocks table
CREATE TABLE IF NOT EXISTS blocks (
	block_number INT PRIMARY KEY,
	block_hash CHAR(32) NOT NULL,
	mined_timestamp TIMESTAMP NOT NULL,
	tx_count INT NOT NULL,
	created_at TIMESTAMP DEFAULT current_timestamp
);

-- transactions table
CREATE TABLE IF NOT EXISTS transactions (
	tx_hash CHAR(32) NOT NULL PRIMARY KEY,
	block_number INT NOT NULL,
	tx_from CHAR(32) NOT NULL,
	tx_to CHAR(32) NOT NULL,
	amount NUMERIC NOT NULL,
	nonce INT NOT NULL,
	mined_timestamp TIMESTAMP NOT NULL,
	tx_order INT NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp,
	FOREIGN KEY (block_number) REFERENCES blocks (block_number) ON DELETE CASCADE
);
