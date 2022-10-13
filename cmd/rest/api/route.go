package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// routes register routes
func (a *API) routes() {
	// middlewares
	a.mux.Use(
		a.recovery,
		a.logger,
	)

	// Not found & Not Allowed
	a.mux.NotFound(a.notFound)
	a.mux.MethodNotAllowed(a.notAllowed)

	// endpoints
	a.mux.Get("/health", a.handleHealthChech)

	//
	a.mux.Route("/v1", func(r chi.Router) {
		r.Get("/index", a.handleIndexerCommand)
		//
		r.Get("/block", a.handleGetLatestBlock)
		r.Get("/block/{id}", a.handleGetBlock)

		//
		r.Get("/stats", a.handleGetStats)
		r.Get("/stats/{range}", a.handleGetRangeStats)

		//
		r.Get("/tx", a.handleGetLatestTx)
		r.Get("/tx/{hash}", a.handleGetTx)
	})
}

// Response
type Response struct {
	Status  int         `json:"status"`
	Payload interface{} `json:"payload"`
}

// writer
func (a *API) writer(w http.ResponseWriter, status int, payload interface{}) {
	response := Response{
		Status:  status,
		Payload: payload,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		a.log.Println(err)
	}
}
