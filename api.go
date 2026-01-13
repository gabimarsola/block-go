package main

import (
	chain "block-go/structure/blockchain"
	"encoding/json"
	"net/http"
)

func registerHandlers(mux *http.ServeMux, bc *chain.Blockchain) {
	/*
		mux.HandleFunc("/tx", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			var tx Transaction
			if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
				http.Error(w, "invalid body", http.StatusBadRequest)
				return
			}
			bc.AddTransaction(tx)
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(map[string]string{"status": "tx accepted"})
		}) */

	mux.HandleFunc("/chain", func(w http.ResponseWriter, r *http.Request) {
		chain := bc.LastHash
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(chain)
	})
}
