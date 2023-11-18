package client

import (
	"encoding/json"
	"fmt"
	"github.com/dwnGnL/ddos-pow/config"
	"github.com/dwnGnL/ddos-pow/internal/application"
	"github.com/dwnGnL/ddos-pow/lib/pow"
	"net/http"
)

type Handler struct {
	conf *config.Config
}

func newHandler(cfg *config.Config) *Handler {
	return &Handler{conf: cfg}
}

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

const (
	SUCCESS = "success"
	ERROR   = "error"
)

func (h Handler) RequestChallenge(w http.ResponseWriter, r *http.Request) {
	app, err := application.GetAppFromRequest(r)
	if err != nil {
		respond(w, http.StatusBadRequest, nil, err)
		return
	}

	clientService := app.GetClient()
	hashCashData, err := clientService.RequestChallenge()

	respond(w, http.StatusOK, hashCashData, nil)
}

func (h Handler) RequestResource(w http.ResponseWriter, r *http.Request) {
	app, err := application.GetAppFromRequest(r)
	if err != nil {
		respond(w, http.StatusBadRequest, nil, err)
		return
	}

	hashCashData := pow.HashcashData{}

	err = json.NewDecoder(r.Body).Decode(&hashCashData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("body hashcashData", hashCashData)

	clientService := app.GetClient()
	quote, err := clientService.RequestResource(hashCashData)

	if err != nil {
		respond(w, http.StatusBadRequest, nil, err)
		return
	}

	respond(w, http.StatusOK, quote, nil)
}

func respond(w http.ResponseWriter, status int, data interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")

	var response Response

	w.WriteHeader(status)

	switch {
	case status >= 200 && status <= 299:
		response.Status = SUCCESS
		response.Data = data
	default:
		response.Status = ERROR
		response.Message = err.Error()
	}

	resp, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
