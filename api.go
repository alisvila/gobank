package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	ListenAddress string
	Store         Storage
}

func writeJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

type apifunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Error string
}

func makeHttpHandler(f apifunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			writeJson(w, 400, err.Error())
		}
	}
}

func newAPIServer(listen string, store Storage) *APIServer {
	return &APIServer{
		ListenAddress: listen,
		Store:         store,
	}
}

func (s *APIServer) run() {
	r := mux.NewRouter()

	r.HandleFunc("/account", makeHttpHandler(s.handleAccount))
	r.HandleFunc("/account/{id}", makeHttpHandler(s.handleAccount))
	http.ListenAndServe(s.ListenAddress, r)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	return fmt.Errorf("method not supported %s", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	fmt.Println(id)
	name := makeAccount("random", "person")
	return writeJson(w, http.StatusOK, name)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {

	return nil
}

// func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
// 	return nil
// }
