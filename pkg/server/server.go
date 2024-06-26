package server

import (
	"context"
	"github.com/umtdemr/go-kafka-with-rest-case/pkg/logger"
	"github.com/umtdemr/go-kafka-with-rest-case/pkg/store"
	"log"
	"math/rand/v2"
	"net/http"
	"time"
)

var srv *http.Server

type Server struct {
	ListenAddr string
	router     *http.ServeMux
}

func NewServer(listenAddr string) *Server {
	return &Server{listenAddr, http.NewServeMux()}
}

func (s *Server) Run() {
	http.ListenAndServe(s.ListenAddr, s.router)
}

func handler(w http.ResponseWriter, r *http.Request) {
	randomWaitTime := time.Duration(rand.Float64()*3*1000) * time.Millisecond
	time.Sleep(randomWaitTime)
	fileLogger := logger.GetLogger()
	timestamp := time.Now()
	fileLogger.Printf("%s,%v,%v\n", r.Method, randomWaitTime.Milliseconds(), timestamp.Unix())
}

func Run(store *store.Store) {
	serv := NewServer(":8080")
	serv.router.HandleFunc("GET /api/get", handler)
	serv.router.HandleFunc("POST /api/post", handler)
	serv.router.HandleFunc("PUT /api/put", handler)
	serv.router.HandleFunc("DELETE /api/delete", handler)

	fileServer := http.FileServer(http.Dir("./website"))
	serv.router.Handle("/", fileServer)

	srv = &http.Server{
		Addr:    serv.ListenAddr,
		Handler: serv.router,
	}
	StartWebSocketServer(store, serv)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()
}

func Shutdown(ctx context.Context, stop chan struct{}) {
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server Shutdown: %v", err)
	}
	close(stop)
	log.Println("HTTP server shutdown completed.")
}
