package api

import "net/http"

func (server *Server) signup(w http.ResponseWriter, r *http.Request) {
	// TODO: real implementation
	w.Write([]byte("signup works!"))
}
