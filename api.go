package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
}

type APIFn func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Error string
}

func writeJsonResponse(
	responseWriter http.ResponseWriter,
	status int,
	value any,
) error {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(status)
	return json.NewEncoder(responseWriter).Encode(value)
}

func makeHTTPHandleFn(apiFn APIFn) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		err := apiFn(responseWriter, httpRequest)
		if err != nil {
			writeJsonResponse(
				responseWriter,
				http.StatusInternalServerError,
				APIError{Error: err.Error()},
			)
		}
	}
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (apiServer *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandleFn(apiServer.handleAccount))

	log.Println("API server is running on", apiServer.listenAddr)
	http.ListenAndServe(apiServer.listenAddr, router)
}

func (apiServer *APIServer) handleAccount(
	responseWriter http.ResponseWriter,
	httpRequest *http.Request,
) error {
	switch httpRequest.Method {
	case http.MethodGet:
		return apiServer.handleGetAccount(responseWriter, httpRequest)
	case http.MethodPost:
		return apiServer.handleCreateAccount(responseWriter, httpRequest)
	case http.MethodDelete:
		return apiServer.handleDeleteAccount(responseWriter, httpRequest)
	default:
		return fmt.Errorf("method not allowed %s", httpRequest.Method)
	}
}

func (apiServer *APIServer) handleGetAccount(
	responseWriter http.ResponseWriter,
	httpRequest *http.Request,
) error {
	account := NewAccount(
		"Hector",
		"Avalos",
		"hg.avalosc97@gmail.com",
	)
	return writeJsonResponse(responseWriter, http.StatusOK, account)
}

func (apiServer *APIServer) handleCreateAccount(
	responseWriter http.ResponseWriter,
	httpRequest *http.Request,
) error {
	return nil
}

func (apiServer *APIServer) handleDeleteAccount(
	responseWriter http.ResponseWriter,
	httpRequest *http.Request,
) error {
	return nil
}

func (apiServer *APIServer) handleTransfer(
	responseWriter http.ResponseWriter,
	httpRequest *http.Request,
) error {
	return nil
}
