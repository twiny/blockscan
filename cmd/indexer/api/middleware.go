package api

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// logger
func (idx *Indexer) logger(h http.Handler) http.Handler {
	return middleware.Logger(h)
}

// recovery
func (idx *Indexer) recovery(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				idx.log.Println("HTTP Recovery", err)
				idx.writer(w, http.StatusInternalServerError, err)
			}
		}()

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
