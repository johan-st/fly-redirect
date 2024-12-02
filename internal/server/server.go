package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"gavmofjall_se/internal/repo"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port       int
	repo       *repo.TursoRepo
	countInMem atomic.Uint32
}

func NewServer() *http.Server {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Error getting port from environment variable. Err: %v", err)
	}
	repo, err := repo.NewTursoRepo()
	if err != nil {
		log.Fatalf("Error creating Turso repository. Err: %v", err)
	}
	NewServer := &Server{
		port:       port,
		repo:       repo,
		countInMem: atomic.Uint32{},
	}

	count, err := NewServer.repo.GetRedirectCount()
	if err != nil {
		log.Fatalf("Error getting redirect count from database. Err: %v", err)
	}
	NewServer.countInMem.Store(uint32(count))

	addr := fmt.Sprintf("0.0.0.0:%d", NewServer.port)
	fmt.Println("Server listening on", addr)
	fmt.Printf("Redirect count loaded from database: %d\n", NewServer.countInMem.Load())

	// Declare Server config
	server := &http.Server{
		Addr:         addr,
		Handler:      NewServer.RegisterRoutes(),
		ReadTimeout:  250 * time.Millisecond,
		WriteTimeout: 1 * time.Second,
	}

	return server
}
