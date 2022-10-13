package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// routes register routes
func (idx *Indexer) routes() {
	// middlewares
	idx.mux.Use(
		idx.recovery,
		idx.logger,
	)

	// Not found & Not Allowed
	idx.mux.NotFound(idx.notFound)
	idx.mux.MethodNotAllowed(idx.notAllowed)

	// endpoints
	idx.mux.Get("/health", idx.handleHealthChech)
	//
	idx.mux.Get("/", idx.handleIndex) // ?auth_token&scan=100:200
}

// Response
type Response struct {
	Status  int         `json:"status"`
	Payload interface{} `json:"payload"`
}

// writer
func (idx *Indexer) writer(w http.ResponseWriter, status int, payload interface{}) {
	response := Response{
		Status:  status,
		Payload: payload,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		idx.log.Println(err)
	}
}

// notFound
func (idx *Indexer) notFound(w http.ResponseWriter, r *http.Request) {
	idx.writer(w, http.StatusNotFound, "not found")
}

// notAllowed
func (idx *Indexer) notAllowed(w http.ResponseWriter, r *http.Request) {
	idx.writer(w, http.StatusMethodNotAllowed, "method not allowed")
}

// // \\ \\
// handleHealthChech
func (idx *Indexer) handleHealthChech(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"version": Version,
		"store":   "up",
	}

	if err := idx.store.Ping(); err != nil {
		health["database"] = "down"
		idx.writer(w, http.StatusInternalServerError, health)
		return
	}

	idx.writer(w, http.StatusOK, health)
}

// handleScan - ?auth_token&scan=100:200
// 100 200 => range
// 100 - the latest => range once reached subscribe to new
// empty => 0:latest
func (idx *Indexer) handleIndex(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	authToken, found := query["auth_token"]
	if !found || len(authToken) < 1 {
		idx.writer(w, http.StatusUnauthorized, fmt.Errorf("auth_token is required"))
		return
	}

	if authToken[0] != idx.conf.Indexer.Token {
		idx.writer(w, http.StatusUnauthorized, fmt.Errorf("auth_token is required"))
	}

	scanRange, found := query["scan"]
	if !found || len(scanRange) < 1 {
		idx.writer(w, http.StatusBadRequest, fmt.Errorf("scan range is required"))
		return
	}

	start, end, err := idx.parseScanQuery(scanRange[0])
	if err != nil {
		idx.writer(w, http.StatusBadRequest, err.Error())
		return
	}

	// add block ids to jobs queue
	go func(s, e int64) {
		for i := s; i <= e; i++ {
			idx.jobs <- i
		}
	}(start, end)

	idx.writer(w, http.StatusOK, "command executed")
}
