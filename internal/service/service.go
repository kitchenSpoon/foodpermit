package service

import (
	"fmt"
	"foodpermit/internal/foodpermit"
	"github.com/gorilla/mux"
	"net"
	"net/http"
)

type Service struct {
	*mux.Router
}

func NewService() (Service, error) {
	s, err := foodpermit.NewService()
	if err != nil {
		return Service{}, err
	}
	r := mux.NewRouter()
	r.HandleFunc("/geosearch", s.Geosearch)
	r.HandleFunc("/autocomplete", s.GetSuggestion)
	r.HandleFunc("/", s.Root)
	return Service{r}, nil
}

func StartService() {
	s, err := NewService()
	if err != nil {
		fmt.Errorf("failed to create service: %w", err)
		return
	}

	server := http.Server{Handler: s}
	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Errorf("failed to create listener: %w", err)
		return
	}

	fmt.Println("Starting server")
	err = server.Serve(l)
	if err != nil {
		fmt.Errorf("unexpected error in server.Serve, %w", err)
	}
}
