package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(port string) *http.Server {
	router := mux.NewRouter()
	return &http.Server{
		Addr:    port,
		Handler: router,
	}
}
