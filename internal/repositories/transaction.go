package repositories

import (
	"basic-rest-api-go/internal/dto"
	"basic-rest-api-go/internal/models"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []dto.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	for i := range details {
		details[i].TransactionID = transactionID
		_, err = stmt.Exec(transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

func (repo *TransactionRepository) GetTransactionsByDateRange(startDate, endDate string) ([]*models.Transaction, error) {
	var transactions []*models.Transaction

	rows, err := repo.db.Query("SELECT id, total_amount FROM transactions WHERE created_at BETWEEN $1 AND $2", startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(&transaction.ID, &transaction.TotalAmount)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

func (repo *TransactionRepository) GetTransactionDetails(transactionID []int) ([]models.TransactionDetail, error) {
	var details []models.TransactionDetail

	rows, err := repo.db.Query(
		"SELECT id, transaction_id, product_id, quantity, subtotal FROM transaction_details WHERE transaction_id = ANY($1)",
		pq.Array(transactionID),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var detail models.TransactionDetail
		err := rows.Scan(&detail.ID, &detail.TransactionID, &detail.ProductID, &detail.Quantity, &detail.Subtotal)
		if err != nil {
			return nil, err
		}
		details = append(details, detail)
	}

	return details, nil
}
