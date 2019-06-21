package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"orchestration"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.Handle("/customers/{customerId}/accounts/{accountId}", &CustomerAccountHandler{})
	srv := &http.Server{
		IdleTimeout:  10 * time.Second,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
		Addr:         fmt.Sprintf(":%s", os.Getenv("SERV_PORT")),
		Handler:      r,
	}

	log.Fatal(srv.ListenAndServe())
}

type CustomerAccountHandler struct {
}

func (h *CustomerAccountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerId := vars["customerId"]
	accountId := vars["accountId"]
	htclient := &http.Client{
		Timeout: 1 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, getAccountUrl("localhost", "8000", accountId), nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := htclient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	contents, err := ioutil.ReadAll(res.Body)
	var account orchestration.Account
	err = json.Unmarshal(contents, &account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req, err = http.NewRequest(http.MethodGet, getCustomerUrl("localhost", "9000", customerId), nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err = htclient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	contents, err = ioutil.ReadAll(res.Body)
	var customer orchestration.Customer
	err = json.Unmarshal(contents, &customer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	custAcct := orchestration.CustomerAccount{
		Customer: customer,
		Account:  account,
	}
	data, err := json.Marshal(custAcct)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func getAccountUrl(host, port, accountId string) string {
	return fmt.Sprintf("http://%s:%s/deposits/accounts/%s", host, port, accountId)
}

func getCustomerUrl(host, port, customerId string) string {
	return fmt.Sprintf("http://%s:%s/customers/%s", host, port, customerId)
}
