package controller

import (
	"encoding/json"
	"log"
	"natan/fingo/dbsqlite"
	"natan/fingo/model"
	"natan/fingo/service"
	"net/http"
)

// GetTransactionByIDHandler handles GET /transactions/{id} and returns the transaction with the given ID.
func GetTransactionByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := dbsqlite.NewDBContext()
	defer cancel()

	idStr := r.PathValue("id")
	id, ok := GetID(idStr, w, r)
	if !ok {
		return
	}

	transaction, err := service.GetTransactionByID(ctx, id)
	if err != nil {
		log.Println(err)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "transaction not found"})
		return
	}

	writeJSON(w, http.StatusOK, *transaction)
}

// GetAllTransactionsHandler handles GET /transactions and returns all transactions.
func GetAllTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := dbsqlite.NewDBContext()
	defer cancel()

	transactionsList, err := service.GetAllTransactions(ctx)
	if err != nil {
		log.Println(err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "problem fetching transactions"})
		return
	}

	writeJSON(w, http.StatusOK, transactionsList)
}

// CreateTransactionHandler handles POST /transactions and creates a new transaction from the request body.
func CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := dbsqlite.NewDBContext()
	defer cancel()

	var transaction model.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		log.Printf("could not decode request body: %v", err)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	transactionRec, err := service.CreateTransaction(ctx, transaction)
	if err != nil {
		log.Println(err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "problem when creating transaction"})
		return
	}

	writeJSON(w, http.StatusCreated, *transactionRec)
}

// UpdateTransactionByIDHandler handles PATCH /transactions/{id} and applies a partial update to the given transaction.
func UpdateTransactionByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := dbsqlite.NewDBContext()
	defer cancel()

	var transactionUpdate *model.TransactionUpdate
	idStr := r.PathValue("id")
	id, ok := GetID(idStr, w, r)
	if !ok {
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&transactionUpdate); err != nil {
		log.Printf("could not decode request body: %v", err)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	transaction, err := service.UpdateTransactionByID(ctx, id, transactionUpdate)
	if err != nil {
		log.Println(err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "problem when updating transaction"})
		return
	}

	writeJSON(w, http.StatusOK, *transaction)
}

// DeleteTransactionByIDHandler handles DELETE /transactions/{id} and removes the transaction with the given ID.
func DeleteTransactionByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := dbsqlite.NewDBContext()
	defer cancel()

	idStr := r.PathValue("id")
	id, ok := GetID(idStr, w, r)
	if !ok {
		return
	}

	rows, err := service.DeleteTransactionByID(ctx, id)
	if err != nil {
		log.Println(err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "problem when deleting transaction"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]int64{"rows_affected": rows})
}
