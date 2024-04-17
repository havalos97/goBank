package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

type APIFn func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Error string
}

func WriteJSONResponse(
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
			if strings.Contains(err.Error(), "not found") {
				WriteJSONResponse(
					responseWriter,
					http.StatusNotFound,
					APIError{Error: err.Error()},
				)
				return
			}
			WriteJSONResponse(
				responseWriter,
				http.StatusInternalServerError,
				APIError{Error: err.Error()},
			)
		}
	}
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (apiServer *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandleFn(apiServer.handleAccounts))
	router.HandleFunc("/account/{uuid}", makeHTTPHandleFn(apiServer.handleGetAccountById))

	log.Println("API server is running on", apiServer.listenAddr)
	http.ListenAndServe(apiServer.listenAddr, router)
}

func (apiServer *APIServer) handleAccounts(
	responseWriter http.ResponseWriter,
	httpRequest *http.Request,
) error {
	switch httpRequest.Method {
	case http.MethodGet:
		return apiServer.findAllAccounts(responseWriter, httpRequest)
	case http.MethodPost:
		return apiServer.createAccount(responseWriter, httpRequest)
	case http.MethodDelete:
		return apiServer.handleDeleteAccount(responseWriter, httpRequest)
	default:
		return fmt.Errorf("method not allowed %s", httpRequest.Method)
	}
}

func (apiServer *APIServer) findAllAccounts(
	responseWriter http.ResponseWriter,
	_ *http.Request,
) error {
	accountList, err := apiServer.store.FindAllAccounts()

	if err != nil {
		return err
	}
	return WriteJSONResponse(responseWriter, http.StatusOK, accountList)
}

func (apiServer *APIServer) handleGetAccountById(
	responseWriter http.ResponseWriter,
	httpRequest *http.Request,
) error {
	if httpRequest.Method == http.MethodGet {
		accountUuid := mux.Vars(httpRequest)["uuid"]
		account, err := apiServer.store.GetAccountByUUID(accountUuid)
		if err != nil {
			return err
		}
		return WriteJSONResponse(responseWriter, http.StatusOK, account)
	}
	return fmt.Errorf("method not allowed %s", httpRequest.Method)
}

func (apiServer *APIServer) createAccount(
	responseWriter http.ResponseWriter,
	httpRequest *http.Request,
) error {
	createAccReq := new(CreateAccountRequest)
	if err := json.NewDecoder(httpRequest.Body).Decode(createAccReq); err != nil {
		return err
	}

	newlyCreatedAcc, err := apiServer.store.CreateAccount(
		NewAccount(
			createAccReq.FirstName,
			createAccReq.LastName,
			createAccReq.Email,
		),
	)
	if err != nil {
		return err
	}
	return WriteJSONResponse(responseWriter, http.StatusCreated, newlyCreatedAcc)
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
