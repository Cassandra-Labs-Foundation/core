package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// CreateAccountRequest represents the expected payload for account creation.
type CreateAccountRequest struct {
	AccountID      string `json:"account_id"`
	InitialBalance int64  `json:"initial_balance"`
}

// TransferRequest represents the expected payload for a funds transfer.
type TransferRequest struct {
	FromAccountID string `json:"from_account_id"`
	ToAccountID   string `json:"to_account_id"`
	Amount        int64  `json:"amount"`
}

// createAccountHandler handles account creation requests.
func createAccountHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	log.Printf("Mock TigerBeetle Server: Creating account: %+v\n", req)
	
	response := map[string]string{
		"message":   "Account created successfully",
		"accountID": req.AccountID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// transferHandler handles funds transfer requests.
func transferHandler(w http.ResponseWriter, r *http.Request) {
	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	log.Printf("Mock TigerBeetle Server: Transferring funds: %+v\n", req)
	
	response := map[string]string{
		"message": "Transfer successful",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/account", createAccountHandler)
	http.HandleFunc("/transfer", transferHandler)
	addr := ":9000"
	fmt.Printf("Starting TigerBeetle mock server on %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}