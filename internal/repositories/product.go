package repositories

import (
	"basic-rest-api-go/internal/models"
	"database/sql"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAll() ([]models.Product, error) {
	query := "SELECT id, name, price, stock FROM products"
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]models.Product, 0)
	for rows.Next() {
		var c models.Product
		err := rows.Scan(&c.ID, &c.Name, &c.Price, &c.Stock)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (repo *ProductRepository) Create(Product *models.Product) error {
	query := "INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRow(query, Product.Name, Product.Price, Product.Stock, Product.CategoryID).Scan(&Product.ID)
	return err
}

func (repo *ProductRepository) GetByID(id int) (*models.Product, error) {
	query := "SELECT p.id, p.name, p.price, p.stock, c.id, c.name, c.description FROM products p LEFT JOIN categories c ON p.category_id = c.id WHERE p.id = $1"
	row := repo.db.QueryRow(query, id)

	var p models.Product
	var c models.Category
	err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &c.ID, &c.Name, &c.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if c.ID != 0 {
		p.CategoryID = c.ID
		p.Category = &c
	}

	return &p, nil
}

func (repo *ProductRepository) Update(Product *models.Product) error {
	query := "UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5"
	_, err := repo.db.Exec(query, Product.Name, Product.Price, Product.Stock, Product.CategoryID, Product.ID)
	return err
}

func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	_, err := repo.db.Exec(query, id)
	return err
}
