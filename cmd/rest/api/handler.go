package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// notFound
func (a *API) notFound(w http.ResponseWriter, r *http.Request) {
	a.writer(w, http.StatusNotFound, "not found")
}

// notAllowed
func (a *API) notAllowed(w http.ResponseWriter, r *http.Request) {
	a.writer(w, http.StatusMethodNotAllowed, "method not allowed")
}

// // \\ \\

// handleHealthChech
func (a *API) handleHealthChech(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"version": Version,
		"store":   "up",
	}

	if err := a.store.Ping(); err != nil {
		health["database"] = "down"
		a.writer(w, http.StatusInternalServerError, health)
		return
	}

	a.writer(w, http.StatusOK, health)
}

// handleIndexerCommand
func (a *API) handleIndexerCommand(w http.ResponseWriter, r *http.Request) {
	// Indexer service endpoint // ?auth_token&scan=100:200
	query := r.URL.Query()
	scanRange, found := query["scan"]
	if !found || len(scanRange) < 1 {
		a.writer(w, http.StatusBadRequest, "scan parameter is required")
		return
	}

	endpoint := fmt.Sprintf("%s%s?auth_token=%s&scan=%s", a.conf.Indexer.Host, a.conf.Indexer.Addr, a.conf.Indexer.Token, scanRange[0])

	resp, err := http.Get(endpoint)
	if err != nil {
		a.writer(w, http.StatusBadGateway, err.Error())
		return
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			a.writer(w, http.StatusBadRequest, err.Error())
			return
		}
		var body map[string]any
		if err := json.Unmarshal(b, &body); err != nil {
			a.writer(w, http.StatusBadRequest, err.Error())
			return
		}

		a.writer(w, http.StatusBadRequest, body)
		return
	}

	// success
	a.writer(w, http.StatusOK, "indexer command executed")
}

// handleGetLatestBlock - returns the latest block and all associated transactions
func (a *API) handleGetLatestBlock(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	block, err := a.store.GetLatestBlock(ctx)
	if err != nil {
		a.log.Printf("err, get_latest_block, %s", err.Error())
		a.writer(w, http.StatusInternalServerError, err.Error())
		return
	}

	a.writer(w, http.StatusOK, block)
}

// handleGetBlock
func (a *API) handleGetBlock(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	num, err := strconv.ParseInt(id, 10, 0)
	if err != nil {
		a.writer(w, http.StatusBadRequest, err.Error())
		return
	}

	// is negative
	if num < 0 {
		num = num * -1
	}

	block, err := a.store.GetBlock(ctx, num)
	if err != nil {
		a.writer(w, http.StatusInternalServerError, err.Error())
		return
	}

	a.writer(w, http.StatusOK, block)
}

// handleGetStats
func (a *API) handleGetStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var latest int64 = 0
	block, err := a.store.GetLatestBlock(ctx)
	if err == nil {
		latest = block.Number
	}

	stats, err := a.store.GetStats(ctx, 0, latest)
	if err != nil {
		a.writer(w, http.StatusInternalServerError, err.Error())
		return
	}

	a.writer(w, http.StatusOK, stats)
}

// handleGetRangeStats
func (a *API) handleGetRangeStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	interval := chi.URLParam(r, "range")

	start, end, err := parseRange(interval)
	if err != nil {
		a.writer(w, http.StatusBadRequest, err.Error())
		return
	}

	stats, err := a.store.GetStats(ctx, start, end)
	if err != nil {
		a.writer(w, http.StatusInternalServerError, err.Error())
		return
	}

	a.writer(w, http.StatusOK, stats)
}

// handleGetLatestTx
func (a *API) handleGetLatestTx(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tx, err := a.store.GetLatestTx(ctx)
	if err != nil {
		a.writer(w, http.StatusInternalServerError, err.Error())
		return
	}

	a.writer(w, http.StatusOK, tx)
}

// handleGetTx
func (a *API) handleGetTx(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	hash := chi.URLParam(r, "hash")

	// TODO: validate tx_hash

	tx, err := a.store.GetTx(ctx, hash)
	if err != nil {
		a.writer(w, http.StatusInternalServerError, err.Error())
		return
	}

	a.writer(w, http.StatusOK, tx)
}
