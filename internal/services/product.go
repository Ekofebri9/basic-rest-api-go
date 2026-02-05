package services

import (
	"basic-rest-api-go/internal/dto"
	"basic-rest-api-go/internal/models"
	"basic-rest-api-go/internal/repositories"
	"errors"
)

type ProductService struct {
	repo         *repositories.ProductRepository
	categoryRepo *repositories.CategoryRepository
}

var ErrInvalidCategoryID = errors.New("invalid category ID")

func NewProductService(repo *repositories.ProductRepository, categoryRepo *repositories.CategoryRepository) *ProductService {
	return &ProductService{repo: repo, categoryRepo: categoryRepo}
}

func (s *ProductService) GetAll(search *dto.SearchProductDTO) ([]models.Product, error) {
	return s.repo.GetAll(search)
}

func (s *ProductService) GetByID(id int) (*models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) Create(Product *models.Product) error {
	isValid, err := s.checkCategoryID(Product.CategoryID)
	if err != nil {
		return err
	}
	if !isValid {
		return ErrInvalidCategoryID
	}
	return s.repo.Create(Product)
}

func (s *ProductService) Update(Product *models.Product) error {
	isValid, err := s.checkCategoryID(Product.CategoryID)
	if err != nil {
		return err
	}
	if !isValid {
		return ErrInvalidCategoryID
	}
	return s.repo.Update(Product)
}

func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *ProductService) checkCategoryID(categoryID int) (bool, error) {
	category, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		return false, err
	}
	return category != nil, nil
}
