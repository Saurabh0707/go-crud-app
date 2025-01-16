package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"go-crud/db"
	"go-crud/models"
	"go-crud/utils"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
)

// ItemOperations is a struct that holds the database interface
type ItemOperations struct {
	DB db.DBInterface
}

// CreateItem handles creating a new item
func (ops *ItemOperations) CreateItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !utils.ValidateItemName(item.Name) {
		http.Error(w, "Invalid item name. Only alphabetic characters are allowed.", http.StatusBadRequest)
		return
	}

	err := ops.DB.QueryRow(context.Background(),
		"INSERT INTO items (name, price) VALUES ($1, $2) RETURNING id",
		item.Name, item.Price).Scan(&item.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(item)
}

// GetItems handles fetching all items
func (ops *ItemOperations) GetItems(w http.ResponseWriter, r *http.Request) {
	rows, err := ops.DB.Query(context.Background(), "SELECT id, name, price FROM items")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		err := rows.Scan(&item.ID, &item.Name, &item.Price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	json.NewEncoder(w).Encode(items)
}

// GetItem handles fetching a single item by ID
func (ops *ItemOperations) GetItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	var item models.Item
	err = ops.DB.QueryRow(context.Background(),
		"SELECT id, name, price FROM items WHERE id=$1", id).Scan(&item.ID, &item.Name, &item.Price)

	if err == pgx.ErrNoRows {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(item)
}

// UpdateItem handles updating an existing item
func (ops *ItemOperations) UpdateItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !utils.ValidateItemName(item.Name) {
		http.Error(w, "Invalid item name. Only alphabetic characters are allowed.", http.StatusBadRequest)
		return
	}

	_, err = ops.DB.Exec(context.Background(),
		"UPDATE items SET name=$1, price=$2 WHERE id=$3",
		item.Name, item.Price, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item.ID = id
	json.NewEncoder(w).Encode(item)
}

// DeleteItem handles deleting an item
func (ops *ItemOperations) DeleteItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	_, err = ops.DB.Exec(context.Background(), "DELETE FROM items WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
