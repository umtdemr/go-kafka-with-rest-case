package server

import "net/http"

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

func Run() {
	serv := NewServer(":8000")
	serv.Run()
}
