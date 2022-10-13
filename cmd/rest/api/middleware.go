package api

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// logger
func (a *API) logger(h http.Handler) http.Handler {
	return middleware.Logger(h)
}

// recovery
func (a *API) recovery(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				a.log.Println("HTTP Recovery", err)
				a.writer(w, http.StatusInternalServerError, err)
			}
		}()

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
