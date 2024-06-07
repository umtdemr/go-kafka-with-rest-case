package server

import (
	"math/rand"
	"net/http"
	"time"
)

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
}

func Run() {
	serv := NewServer(":8000")
	serv.router.HandleFunc("GET /api/get", handler)
	serv.router.HandleFunc("POST /api/post", handler)
	serv.router.HandleFunc("PUT /api/put", handler)
	serv.router.HandleFunc("DELETE /api/delete", handler)
	serv.Run()
}
