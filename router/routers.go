package router

import (
	"go-crud/handlers"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

func NewRouter(db *pgx.Conn) *mux.Router {
	r := mux.NewRouter()
	itemOps := &handlers.ItemOperations{DB: db}

	r.HandleFunc("/api/items", itemOps.CreateItem).Methods("POST")
	r.HandleFunc("/api/items", itemOps.GetItems).Methods("GET")
	r.HandleFunc("/api/items/{id}", itemOps.GetItem).Methods("GET")
	r.HandleFunc("/api/items/{id}", itemOps.UpdateItem).Methods("PUT")
	r.HandleFunc("/api/items/{id}", itemOps.DeleteItem).Methods("DELETE")

	return r
}
