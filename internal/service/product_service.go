package service

import (
	"context"
	"errors"
	"fmt"
	"product-microservice/internal/domain"
	"product-microservice/internal/repository"
	pb "product-microservice/proto/product"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProductService interface {
	CreateProduct(product *domain.Product) (*domain.Product, error)
	GetProductByID(id uuid.UUID) (*domain.Product, error)
	UpdateProduct(id uuid.UUID, updatedProduct *domain.Product) (*domain.Product, error)
	DeleteProduct(id uuid.UUID) error
	ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error)
	FindProductById(ctx context.Context, id string) (*domain.Product, error)
}

type productService struct {
	ProductRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) *productService {
	return &productService{
		ProductRepo: productRepo,
	}
}

func (s *productService) CreateProduct(product *domain.Product) (*domain.Product, error) {
	err := s.ProductRepo.Create(product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productService) GetProductByID(id uuid.UUID) (*domain.Product, error) {
	// Calling the repository method to fetch the product by ID
	product, err := s.ProductRepo.GetByID(id)
	if err != nil {
		// Return an error if the product is not found or another error occurs
		return nil, fmt.Errorf("error fetching product with ID %s: %v", id, err)
	}
	return product, nil
}

func (s *productService) UpdateProduct(id uuid.UUID, updatedProduct *domain.Product) (*domain.Product, error) {
	// Get the current product details by ID
	product, err := s.ProductRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %v", err)
	}

	// Update product fields
	product.Name = updatedProduct.Name
	product.Description = updatedProduct.Description
	product.Price = updatedProduct.Price

	// Update product in the database
	err = s.ProductRepo.Update(product)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %v", err)
	}

	return product, nil
}


func (s *productService) DeleteProduct(id uuid.UUID) error {
	return s.ProductRepo.Delete(id)
}

func (s *productService) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	var products []domain.Product
	var err error

	// Determine which product type to fetch
	switch req.Type {
	case "digital":
		products, err = s.ProductRepo.GetDigitalProducts(ctx)
	case "physical":
		products, err = s.ProductRepo.GetPhysicalProducts(ctx)
	case "subscription":
		products, err = s.ProductRepo.GetSubscriptionProducts(ctx)
	default:
		products, err = s.ProductRepo.GetAllProducts(ctx)
	}

	if err != nil {
		return nil, err
	}

	// Convert products to the protobuf response format
	var pbProducts []*pb.Product
	for _, product := range products {
		pbProduct := &pb.Product{
			Id:          product.ID.String(),
			Name:        product.Name,
			Description: product.Description,
			Price:       float32(product.Price),
			CreatedAt:   timestamppb.New(product.CreatedAt),
			UpdatedAt:   timestamppb.New(product.UpdatedAt),
		}

		// Populate the product type-specific fields
		if product.DigitalProduct != nil {
			pbProduct.ProductType = &pb.Product_DigitalProduct{
				DigitalProduct: &pb.DigitalProduct{
					FileSize:     product.DigitalProduct.FileSize,
					DownloadLink: product.DigitalProduct.DownloadLink,
				},
			}
		}
		pbProducts = append(pbProducts, pbProduct)
	}

	// Return the response with products
	return &pb.ListProductsResponse{
		Products: pbProducts,
	}, nil
}

func (s *productService) FindProductById(ctx context.Context, id string) (*domain.Product, error) {
	if id == "" {
		return nil, errors.New("product name cannot be empty")
	}

	return s.ProductRepo.FindById(id)
}