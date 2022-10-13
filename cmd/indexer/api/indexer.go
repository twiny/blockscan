package api

import (
	"context"
	_ "embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/twiny/blockscan/pkg/config"
	"github.com/twiny/blockscan/service/sqlite"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-chi/chi/v5"
	"github.com/twiny/ratelimit"
)

//go:embed version
var Version string

// Indexer
type Indexer struct {
	wg *sync.WaitGroup
	//
	conf *config.Config
	//
	mux *chi.Mux
	srv *http.Server
	//
	limiter *ratelimit.Limiter
	client  *ethclient.Client
	//
	jobs chan int64 // queue of block ids to scan
	//
	// once used to only subscribe to
	// `client.SubscribeNewHead` once.
	subscribed bool
	head       int64 // latest block
	events     chan *types.Header
	//
	store StoreWriter
	//
	log *log.Logger
	//
	ctx  context.Context
	done context.CancelFunc
}

// NewIndexer
func NewIndexer(path string) (*Indexer, error) {
	conf, err := config.ParseConfig(path)
	if err != nil {
		return nil, err
	}

	//
	mux := chi.NewRouter()

	// db
	store, err := sqlite.NewSQLiteDB(conf.Store.Path)
	if err != nil {
		return nil, err
	}

	if err := store.Ping(); err != nil {
		return nil, err
	}

	//
	if err := store.Migrate("up"); err != nil {
		return nil, err
	}

	client, err := ethclient.Dial(conf.Indexer.Endpoint)
	if err != nil {
		return nil, err
	}

	// get current latest block
	latest, err := client.BlockByNumber(context.Background(), nil) // TODO: timeout context
	if err != nil {
		return nil, err
	}

	// indexer ctx
	ctx, done := context.WithCancel(context.Background())

	// indexer
	idx := &Indexer{
		wg: &sync.WaitGroup{},
		//
		conf: conf,
		//
		mux: mux,
		srv: &http.Server{
			Addr:         conf.Indexer.Addr,
			Handler:      mux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 20 * time.Second,
			IdleTimeout:  10 * time.Second,
		},
		//
		limiter: ratelimit.NewLimiter(conf.Indexer.Limiter.Rate, conf.Indexer.Limiter.Duration),
		client:  client,
		//
		jobs: make(chan int64, 16), // TODO: chan size in conf
		//
		subscribed: false,
		head:       latest.Number().Int64(),
		events:     make(chan *types.Header, 16), // TODO: chan size in conf
		//
		store: store,
		//
		log: log.Default(),
		//
		ctx:  ctx,
		done: done,
	}

	// start indexer
	idx.indexer()

	return idx, nil
}

// Start
func (idx *Indexer) Start() error {
	// add routes
	idx.routes()

	log.Println("starting http server on", idx.srv.Addr)

	if err := idx.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Shutdown
func (idx *Indexer) Shutdown() {
	// signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// 2nd ctrl+c kills program
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-sigs
		log.Println("killing program ...")
		os.Exit(0)
	}()

	// hold
	<-sigs
	log.Println("shutting down ...")

	if err := idx.srv.Shutdown(context.TODO()); err != nil {
		idx.log.Println(err)
	}

	idx.wg.Wait()

	close(idx.jobs)
	close(idx.events)

	log.Println("goodbye.")
	os.Exit(0)
}
