# Blockchain Explorer

mini blockchain explorer - composed of two services `Rest` & `Indexer`


### `Indexer`
Is the backend service for the chain explorer, it gets the details of a block along with the details of transactions in it then stores them in a DB.

### `Rest`

An HTTP server that read from a database and exposes public endpoints to view a range of block/transactions as well as statistics about the chain.

```
`GET /health`           - health check endpoint
`GET /v1/index`         - instruct Indexer to perform a scan

`GET /v1/block`         - get latest block in db
`GET /v1/block/{id}`    - get a specific block 

`GET /v1/stats`         - get stats total amount of transactions and all transaction hashes in DB.
`GET /v1/stats/{range}` - get stats for a range of blocks `start:end`

`GET /v1/tx`            - get latest transaction id db.
`GET /v1/tx/{hash}`     - get transaction by hash
```

#### `Indexer Store`

Persistence layer `StoreWriter` an interface expose below API. 

```go
// StoreWriter
type StoreWriter interface {
    Ping() error
    HasScanned(ctx context.Context, id int64) bool
    SaveBlock(ctx context.Context, block *chain.Block) error
    SaveTx(ctx context.Context, tx *chain.Tx) error
}
```

#### `Rest Store`

`StoreReader` reads from a DB and an interface exposing these APIs. 
```go
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
```

## Configuration
```yaml
# rest configuration
rest:
    address: ":8080"

# indexer configuration
indexer:
    address: ":8081"
    host: "http://localhost"
    token: "secret"
    endpoint: "wss://mainnet.infura.io/ws/v3/{api_key}"
    limiter:
        rate: 3
        duration: "1s"
    workers: 5
    timeout: "30s"
    
# store configuration
store:
    path: "./tmp/"
```

## Run

rename file `config/example.config.yaml`  to `config/config.yaml` and uodate config as per requirement.

Then, run `make inderxer` in one terminal to start the indexer service and run `make rest` in another terminal to start the rest service.

once both services are up
```
[rest service] 2022/10/02 19:50:05 starting http server on :8080
[indexer service] 2022/10/02 19:50:06 starting http server on :8081
```

Run in a 3rd terminal window, to instruct the indexer to start scanning the blockchain, from block 15661751 to the latest one.

```
curl -XGET -H "Content-type: application/json" 'http://localhost:8080/v1/index?scan=15661751'
```

A Success Response:

```
{"status":200,"payload":"indexer command executed"}
```

View Postman collection `postman/blockchain_explorer.postman_collection.json` for all `rest` service endpoints/APIs.


## TODO

- [ ] add more tests.
- [ ] add store mocks.
- [ ] -