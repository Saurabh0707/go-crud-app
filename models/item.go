package models

// Item represents a simple data structure with ID, Name, and Price fields.
type Item struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}
