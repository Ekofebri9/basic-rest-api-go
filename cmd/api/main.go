package main

import (
	"basic-rest-api-go/configs"
	"basic-rest-api-go/databases"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Pruduct struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

var products = []Pruduct{
	{ID: 1, Name: "Indomie Goreng", Price: 3500, Stock: 10},
	{ID: 2, Name: "Indomie Goreng Jumbo", Price: 4500, Stock: 20},
	{ID: 3, Name: "Indomie Soto", Price: 3500, Stock: 30},
}

func main() {

	// initialized config
	config := configs.Init()

	// initialized database
	fmt.Println("Connecting to database...", config.DBConn)
	db, err := databases.Init(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintln(w, "OK")
		// w.Write([]byte("OK"))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/api/products/", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			for _, product := range products {
				if product.ID == id {
					w.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(w).Encode(product)
					w.WriteHeader(http.StatusOK)
					return
				}
			}
			http.Error(w, "Product not found", http.StatusNotFound)
		case http.MethodPut:
			var updatedProduct Pruduct
			err := json.NewDecoder(r.Body).Decode(&updatedProduct)
			if err != nil {
				http.Error(w, "Invalid request payload", http.StatusBadRequest)
				return
			}

			for i, product := range products {
				if product.ID == id {
					updatedProduct.ID = id
					products[i] = updatedProduct
					w.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(w).Encode(updatedProduct)
					w.WriteHeader(http.StatusOK)
					return
				}
			}
		case http.MethodDelete:
			for i, product := range products {
				if product.ID == id {
					products = append(products[:i], products[i+1:]...)
					w.WriteHeader(http.StatusNoContent)
					return
				}
			}
			http.Error(w, "Product not found", http.StatusNotFound)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

		for _, product := range products {
			if product.ID == id {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(product)
				w.WriteHeader(http.StatusOK)
				return
			}
		}
		http.Error(w, "Product not found", http.StatusNotFound)
	})

	http.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(products)
			w.WriteHeader(http.StatusOK)
		case http.MethodPost:
			var newProduct Pruduct
			err := json.NewDecoder(r.Body).Decode(&newProduct)
			if err != nil {
				http.Error(w, "Invalid request payload", http.StatusBadRequest)
				return
			}
			newProduct.ID = len(products) + 1
			products = append(products, newProduct)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(newProduct)
			w.WriteHeader(http.StatusCreated)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	addr := fmt.Sprintf(":%s", config.Port)
	fmt.Println("Starting server on :", config.Port)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
