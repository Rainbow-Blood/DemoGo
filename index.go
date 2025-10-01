package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Demo data structure
type Item struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

// In-memory store
var items = []Item{
	{ID: 1, Name: "Apple", Price: 100},
	{ID: 2, Name: "Banana", Price: 50},
}

// Helper to write JSON response
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// API 1: Get all items
func getItems(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, items)
}

// API 2: Get item by ID
func getItemByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	for _, item := range items {
		if item.ID == id {
			writeJSON(w, http.StatusOK, item)
			return
		}
	}

	writeJSON(w, http.StatusNotFound, map[string]string{"error": "item not found"})
}

// API 3: Create new item
func createItem(w http.ResponseWriter, r *http.Request) {
	var newItem Item
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	newItem.ID = len(items) + 1
	items = append(items, newItem)
	writeJSON(w, http.StatusCreated, newItem)
}

// API 4: Delete item by ID
func deleteItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	for i, item := range items {
		if item.ID == id {
			items = append(items[:i], items[i+1:]...)
			writeJSON(w, http.StatusOK, map[string]string{"message": "item deleted"})
			return
		}
	}

	writeJSON(w, http.StatusNotFound, map[string]string{"error": "item not found"})
}

func main() {
	http.HandleFunc("/items", getItems)       // GET /items
	http.HandleFunc("/item", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getItemByID(w, r)   // GET /item?id=1
		case http.MethodPost:
			createItem(w, r)    // POST /item
		case http.MethodDelete:
			deleteItem(w, r)    // DELETE /item?id=1
		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		}
	})

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
