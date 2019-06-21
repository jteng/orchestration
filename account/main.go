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

var accountDb = map[string]orchestration.Account{
	"1":   {AccountID: "1", AccountType: "Savings", OpeningBalance: 123.4, CurrentBalance: 555.98},
	"11":  {AccountID: "11", AccountType: "checkings", OpeningBalance: 12.4, CurrentBalance: 55.98},
	"111": {AccountID: "111", AccountType: "checkings", OpeningBalance: 1.4, CurrentBalance: 5.98},
}

func main() {
	r := mux.NewRouter()

	r.Handle("/deposits/accounts/{accountId}", NewAccountHandler(findAccount))
	srv := &http.Server{
		IdleTimeout:  10 * time.Second,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
		Addr:         fmt.Sprintf(":%s", os.Getenv("SERV_PORT")),
		Handler:      r,
	}

	log.Fatal(srv.ListenAndServe())
}

type AccountHandler struct {
	Database func(string) (orchestration.Account, bool)
}

func (h *AccountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountId := vars["accountId"]
	acct, ok := h.Database(accountId)
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

func NewAccountHandler(fn func(string) (orchestration.Account, bool)) *AccountHandler {
	return &AccountHandler{
		Database: fn,
	}
}

func findAccount(accountId string) (orchestration.Account, bool) {
	acct, ok := accountDb[accountId]
	return acct, ok
}
