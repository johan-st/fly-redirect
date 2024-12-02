package server

import (
	"fmt"
	"net/http"
	"strconv"
)

var count int

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
	count++
	http.Redirect(w, r, "https://www.svtplay.se/julkalendern-snodrommar?cnt="+strconv.Itoa(count)+"&ref=ywlxuignu@mozmail.com", http.StatusTemporaryRedirect)
}

func (s *Server) handlerInfo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Count: %d", count)))
}
