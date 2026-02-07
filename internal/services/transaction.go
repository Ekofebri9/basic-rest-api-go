package services

import (
	"basic-rest-api-go/internal/dto"
	"basic-rest-api-go/internal/models"
	"basic-rest-api-go/internal/repositories"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []dto.CheckoutItem) (*models.Transaction, error) {
	return s.repo.CreateTransaction(items)
}
