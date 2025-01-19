package grpc

import (
	"context"
	"fmt"
	"product-microservice/internal/domain"
	"product-microservice/internal/service"
	pb "product-microservice/proto/product"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProductHandler struct {
	ProductService service.ProductService
	pb.UnimplementedProductServiceServer
}

// NewProductHandler creates a new ProductHandler instance
func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		ProductService: productService,
	}
}

func (h *ProductHandler) CreateProduct(ctx context.Context, req *pb.Product) (*pb.Product, error) {
	createdAt := time.Now()
	updatedAt := time.Now()

	// Create a product object in the domain layer
	domainProduct := &domain.Product{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Price:       float64(req.Price),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	// Assign the correct product type based on the request
	if req.ProductType != nil {
		switch pt := req.ProductType.(type) {
		case *pb.Product_DigitalProduct:
			domainProduct.DigitalProduct = &domain.DigitalProduct{
				FileSize:    pt.DigitalProduct.FileSize,
				DownloadLink: pt.DigitalProduct.DownloadLink,
			}
		case *pb.Product_PhysicalProduct:
			domainProduct.PhysicalProduct = &domain.PhysicalProduct{
				Weight:     pt.PhysicalProduct.Weight,
				Dimensions: pt.PhysicalProduct.Dimensions,
			}
		case *pb.Product_SubscriptionProduct:
			domainProduct.SubscriptionProduct = &domain.SubscriptionProduct{
				SubscriptionPeriod: pt.SubscriptionProduct.SubscriptionPeriod,
				RenewalPrice:      pt.SubscriptionProduct.RenewalPrice,
			}
		default:
			return nil, fmt.Errorf("unsupported product type")
		}
	}

	// Persist the domain product to the database
	newProduct, err := h.ProductService.CreateProduct(domainProduct)
	if err != nil {
		return nil, err
	}

	// Return the created product as a proto response
	return &pb.Product{
		Id:          newProduct.ID.String(),
		Name:        newProduct.Name,
		Description: newProduct.Description,
		Price:       float32(newProduct.Price),
		CreatedAt:   timestamppb.New(createdAt),
		UpdatedAt:   timestamppb.New(updatedAt),
		ProductType: req.ProductType, 
	}, nil
}

// gRPC handler for fetching product by ID
func (h *ProductHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
    // Convert the product ID from string to uuid.UUID
    productID, err := uuid.Parse(req.GetId())
    if err != nil {
        return nil, fmt.Errorf("invalid product ID format: %v", err)
    }

    // Call the service method to get the product
    product, err := h.ProductService.GetProductByID(productID)
    if err != nil {
        return nil, fmt.Errorf("failed to get product: %v", err)
    }

    // Convert the product to the gRPC response format
    productResp := &pb.ProductResponse{
        Id:          product.ID.String(),
        Name:        product.Name,
        Description: product.Description,
        Price:       float32(product.Price),
        CreatedAt:   timestamppb.New(product.CreatedAt),
        UpdatedAt:   timestamppb.New(product.UpdatedAt),
    }

    return productResp, nil
}

func (h *ProductHandler) UpdateProduct(ctx context.Context, req *pb.Product) (*pb.Product, error) {
    // Convert the product ID to uuid.UUID
    productID, err := uuid.Parse(req.GetId())
    if err != nil {
        return nil, fmt.Errorf("invalid product ID format: %v", err)
    }

    // Convert the gRPC product to the domain product
    domainProduct := &domain.Product{
        ID:          productID,
        Name:        req.GetName(),
        Description: req.GetDescription(),
        Price:       float64(req.GetPrice()),
    }

    // Call the service method to update the product
    updatedProduct, err := h.ProductService.UpdateProduct(productID, domainProduct)
    if err != nil {
        return nil, fmt.Errorf("failed to update product: %v", err)
    }

    // Convert the updated product to gRPC response format
    return &pb.Product{
        Id:          updatedProduct.ID.String(),
        Name:        updatedProduct.Name,
        Description: updatedProduct.Description,
        Price:       float32(updatedProduct.Price),
    }, nil
}

func (h *ProductHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*emptypb.Empty, error) {
    // Convert the product ID from string to uuid.UUID
    productID, err := uuid.Parse(req.GetId())
    if err != nil {
        return nil, fmt.Errorf("invalid product ID format: %v", err)
    }

    // Call the service method to delete the product
    err = h.ProductService.DeleteProduct(productID)
    if err != nil {
        return nil, fmt.Errorf("failed to delete product: %v", err)
    }

    // Return an empty response
    return &emptypb.Empty{}, nil
}

// ListProducts handles the ListProducts gRPC method
func (h *ProductHandler) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	// Call the service to get the list of products
	response, err := h.ProductService.ListProducts(ctx, req)
	if err != nil {
		return nil, err
	}

	// Return the list of products
	return response, nil
}