package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", s.handlerCountAndRedirect)
	mux.HandleFunc("/info", s.handlerInfo)

	// Wrap the mux with CORS middleware
	return s.corsMiddleware(mux)
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "gavmofjäll.se,gavmofjall_se.fly.dev")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false")

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handlerCountAndRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("https://www.svtplay.se/julkalendern-snodrommar?cnt=%d&ref=ywlxuignu@mozmail.com", s.countInMem.Add(1)), http.StatusTemporaryRedirect)
	go func() {
		err := s.repo.IncrementRedirectCount()
		if err != nil {
			fmt.Printf("Error: failed to increment counter, %e", err)
		} else {
			fmt.Printf("counter incremented to %d", s.countInMem.Load())
		}
	}()
}

func (s *Server) handlerInfo(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		RedirectsCount uint32 `json:"redirects_count"`
		DbStatus       string `json:"db_status"`
	}

	var dbStatus string
	if err := s.repo.HealthCheck(); err != nil {
		dbStatus = "error"
		fmt.Printf("ERROR: Database health check failed. Err: %v", err)
	} else {
		dbStatus = "ok"
	}

	response := Response{RedirectsCount: s.countInMem.Load(), DbStatus: dbStatus}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
