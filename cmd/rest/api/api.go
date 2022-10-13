package api

import (
	"context"
	_ "embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/twiny/blockscan/pkg/config"
	"github.com/twiny/blockscan/service/sqlite"

	"github.com/go-chi/chi/v5"
)

//go:embed version
var Version string

// API
type API struct {
	conf *config.Config
	//
	mux *chi.Mux
	srv *http.Server
	//
	store StoreReader
	//
	log *log.Logger
}

// NewAPI
func NewAPI(path string) (*API, error) {
	conf, err := config.ParseConfig(path)
	if err != nil {
		return nil, err
	}

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

	return &API{
		conf: conf,
		//
		mux: mux,
		srv: &http.Server{
			Addr:         conf.Rest.Addr,
			Handler:      mux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 20 * time.Second,
			IdleTimeout:  10 * time.Second,
		},
		//
		store: store,
		//
		log: log.Default(),
	}, nil
}

// Start
func (a *API) Start() error {
	// add routes
	a.routes()

	log.Println("starting http server on", a.srv.Addr)

	if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Shutdown
func (a *API) Shutdown() {
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

	if err := a.srv.Shutdown(context.TODO()); err != nil {
		a.log.Println(err)
	}

	log.Println("goodbye.")
	os.Exit(0)
}
