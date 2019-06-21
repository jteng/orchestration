package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"orchestration"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.Handle("/customers/{customerId}", NewAccountHandler(findAccount))
	srv := &http.Server{
		IdleTimeout:  10 * time.Second,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
		Addr:         fmt.Sprintf(":%s", os.Getenv("SERV_PORT")),
		Handler:      r,
	}

	log.Fatal(srv.ListenAndServe())
}

type CustomerHandler struct {
	Database func(string) (orchestration.Customer, bool)
}

func (h *CustomerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerId := vars["customerId"]
	acct, ok := h.Database(customerId)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	data, err := json.Marshal(acct)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func NewAccountHandler(fn func(string) (orchestration.Customer, bool)) *CustomerHandler {
	return &CustomerHandler{
		Database: fn,
	}
}

func findAccount(custId string) (orchestration.Customer, bool) {
	cust, ok := customerDb[custId]
	return cust, ok
}

var customerDb = map[string]orchestration.Customer{
	"c1":  {ID: "c1", FirstName: "Anne", LastName: "Smith"},
	"c11": {ID: "c11", FirstName: "Bonnie", LastName: "Smith"},
	"c12": {ID: "c12", FirstName: "Cindy", LastName: "Smith"},
	"c13": {ID: "c13", FirstName: "Dana", LastName: "Smith"},
}
