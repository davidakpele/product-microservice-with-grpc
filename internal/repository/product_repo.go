package repository

import (
	"fmt"
	"product-microservice/internal/domain"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductRepository interface
type ProductRepository interface {
	Create(product *domain.Product) error
	GetByID(id uuid.UUID) (*domain.Product, error)
	Update(product *domain.Product) error
	Delete(id uuid.UUID) error
	FindById(id string) (*domain.Product, error)
	GetAllProducts(ctx context.Context) ([]domain.Product, error)
	GetDigitalProducts(ctx context.Context) ([]domain.Product, error)
	GetPhysicalProducts(ctx context.Context) ([]domain.Product, error)
	GetSubscriptionProducts(ctx context.Context) ([]domain.Product, error)
}

// ProductRepositoryImpl struct implements ProductRepository interface
type ProductRepositoryImpl struct {
	DB *gorm.DB
}

// NewProductRepository is the constructor for ProductRepositoryImpl
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &ProductRepositoryImpl{
		DB: db,
	}
}

// Create product in the database
func (r *ProductRepositoryImpl) Create(product *domain.Product) error {
	return r.DB.Create(product).Error
}

// GetByID retrieves a product by its ID, including related data for DigitalProduct, PhysicalProduct, and SubscriptionProduct
func (r *ProductRepositoryImpl) GetByID(id uuid.UUID) (*domain.Product, error) {
	var product domain.Product
	// Preload related entities and fetch the product by ID
	err := r.DB.Preload("DigitalProduct").
		Preload("PhysicalProduct").
		Preload("SubscriptionProduct").
		First(&product, "id = ?", id).Error

	// If no product is found, return a specific error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product with ID %s not found", id)
		}
		// Return other errors as they occur
		return nil, fmt.Errorf("failed to fetch product with ID %s: %v", id, err)
	}

	return &product, nil
}

// Update product in the database
func (r *ProductRepositoryImpl) Update(product *domain.Product) error {
	return r.DB.Save(product).Error
}

// Delete product from the database
func (r *ProductRepositoryImpl) Delete(id uuid.UUID) error {
	return r.DB.Delete(&domain.Product{}, "id = ?", id).Error
}

func (r *ProductRepositoryImpl) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product
	if err := r.DB.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepositoryImpl) GetDigitalProducts(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product
	if err := r.DB.Where("digital_product_id IS NOT NULL").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepositoryImpl) GetPhysicalProducts(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product
	if err := r.DB.Where("physical_product_id IS NOT NULL").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepositoryImpl) GetSubscriptionProducts(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product
	if err := r.DB.Where("subscription_product_id IS NOT NULL").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepositoryImpl) FindById(id string) (*domain.Product, error) {
    var product domain.Product
    if err := r.DB.Where("id = ?", id).First(&product).Error; err != nil {
        return nil, err
    }
    return &product, nil
}