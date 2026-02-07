package main

import (
	"basic-rest-api-go/configs"
	"basic-rest-api-go/databases"
	"basic-rest-api-go/internal/handlers"
	"basic-rest-api-go/internal/repositories"
	"basic-rest-api-go/internal/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {

	// initialized config
	config := configs.Init()

	// initialized database
	db, err := databases.Init(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	productRepo := repositories.NewProductRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	productsService := services.NewProductService(productRepo, categoryRepo)
	transactionService := services.NewTransactionService(transactionRepo)

	productHandler := handlers.NewProductHandler(productsService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

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

	http.HandleFunc("/api/products/", productHandler.HandleProductByID)
	http.HandleFunc("/api/products", productHandler.HandleProducts)

	http.HandleFunc("/api/checkout", transactionHandler.HandleCheckout) // POST

	addr := fmt.Sprintf(":%s", config.Port)
	fmt.Println("Starting server on :", config.Port)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
