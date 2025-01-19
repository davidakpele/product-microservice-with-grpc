package test

import (
	"context"
	"fmt"
	"product-microservice/internal/domain"
	"product-microservice/internal/repository"
	"product-microservice/internal/service"
	"product-microservice/internal/transport/grpc"
	pb "product-microservice/proto/product"
	"testing"
	"time"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	
)
type MockProductRepository struct {
    mock.Mock
}

// Mock GetAllProducts method
func (m *MockProductRepository) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
    args := m.Called(ctx)
    return args.Get(0).([]domain.Product), args.Error(1)
}

// Mock GetDigitalProducts method
func (m *MockProductRepository) GetDigitalProducts(ctx context.Context) ([]domain.Product, error) {
    args := m.Called(ctx)
    return args.Get(0).([]domain.Product), args.Error(1)
}

// Mock Create method
func (m *MockProductRepository) Create(product *domain.Product) error {
    args := m.Called(product)
    return args.Error(0)
}

// Mock GetPhysicalProducts method
func (m *MockProductRepository) GetPhysicalProducts(ctx context.Context) ([]domain.Product, error) {
    args := m.Called(ctx)
    return args.Get(0).([]domain.Product), args.Error(1)
}

// Mock GetSubscriptionProducts method
func (m *MockProductRepository) GetSubscriptionProducts(ctx context.Context) ([]domain.Product, error) {
    args := m.Called(ctx)
    return args.Get(0).([]domain.Product), args.Error(1)
}

// Mock GetByID method
func (m *MockProductRepository) GetByID(id uuid.UUID) (*domain.Product, error) {
    args := m.Called(id)
    if args.Get(0) != nil {
        return args.Get(0).(*domain.Product), args.Error(1)
    }
    return nil, args.Error(1)
}

// Mock Delete method
func (m *MockProductRepository) Delete(id uuid.UUID) error {
    args := m.Called(id)
    return args.Error(0)
}

// Mock Update method
func (m *MockProductRepository) Update(product *domain.Product) error {
    args := m.Called(product)
    return args.Error(0)
}

// Mock FindById method
func (m *MockProductRepository) FindById(id string) (*domain.Product, error) {
    args := m.Called(id)
    if args.Get(0) != nil {
        return args.Get(0).(*domain.Product), args.Error(1)
    }
    return nil, args.Error(1)
}

// setupTestDatabase sets up a PostgreSQL database connection for testing using the existing db and config setup
func setupTestDatabase(t *testing.T) *gorm.DB {
	// Database connection details (from your config)
	dbHost := "localhost"
	dbUser := "postgres"
	dbPassword := "powergrid@2?.net"
	dbName := "product_microservice"
	dbPort := "5432"
	dbSslMode := "disable"
	dbTimezone := "UTC"

	// Construct the connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSslMode, dbTimezone)

	// Open the database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}

	// Auto-migrate models
	err = db.AutoMigrate(&domain.Product{}, &domain.DigitalProduct{}, &domain.PhysicalProduct{}, &domain.SubscriptionProduct{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

// setupServiceAndHandler sets up the service and handler for testing
func setupServiceAndHandler(_ *testing.T, db *gorm.DB) (service.ProductService, *grpc.ProductHandler) {
	repo := repository.NewProductRepository(db)
	service := service.NewProductService(repo)
	handler := grpc.NewProductHandler(service)
	return service, handler
}

// TestCreateProductIntegration tests the integration with PostgreSQL
func TestCreateProductIntegration(t *testing.T) {
	// Setup the database and services
	db := setupTestDatabase(t)
	_, handler := setupServiceAndHandler(t, db)

	// Define multiple products with additional fields for related products
	products := []struct {
		name                string
		price               float32
		productType         string
		digitalProduct      *pb.DigitalProduct
		physicalProduct     *pb.PhysicalProduct
		subscriptionProduct *pb.SubscriptionProduct
	}{
		{"Product A", 19.99, "digital", &pb.DigitalProduct{FileSize: 100, DownloadLink: "http://example.com/a"}, nil, nil},
		{"Product B", 29.99, "physical", nil, &pb.PhysicalProduct{Weight: 2.5, Dimensions: "10x10x5"}, nil},
		{"Product C", 39.99, "subscription", nil, nil, &pb.SubscriptionProduct{SubscriptionPeriod: "1 year", RenewalPrice: 10.0}},
		{"Product D", 49.99, "digital", &pb.DigitalProduct{FileSize: 150, DownloadLink: "http://example.com/d"}, nil, nil},
		{"Product E", 59.99, "physical", nil, &pb.PhysicalProduct{Weight: 5.0, Dimensions: "20x20x10"}, nil},
	}

	// Loop over the products and test the CreateProduct function
	for _, p := range products {
		t.Run(fmt.Sprintf("Test create product: %s", p.name), func(t *testing.T) {
			// Create request dynamically without using 'productType' variable
			req, err := createProductRequest(p)
			if err != nil {
				t.Fatalf("Failed to create product request: %v", err)
			}

			// Call handler
			resp, err := handler.CreateProduct(context.Background(), req)

			// Assert no error occurred
			assert.NoError(t, err)

			// Log success message
			t.Logf("Successfully created product: %s with price: %.2f", resp.GetName(), resp.GetPrice())

			// Ensure the related product data is populated based on type
			switch p.productType {
			case "digital":
				assert.NotNil(t, resp.GetDigitalProduct())
			case "physical":
				assert.NotNil(t, resp.GetPhysicalProduct())
			case "subscription":
				assert.NotNil(t, resp.GetSubscriptionProduct())
			}
		})
	}
}

// Helper function to dynamically create the Product request
func createProductRequest(p struct {
	name                string
	price               float32
	productType         string
	digitalProduct      *pb.DigitalProduct
	physicalProduct     *pb.PhysicalProduct
	subscriptionProduct *pb.SubscriptionProduct
}) (*pb.Product, error) {
	// Return the Product request directly based on the type
	switch p.productType {
	case "digital":
		return &pb.Product{
			Name:        p.name,
			Price:       p.price,
			ProductType: &pb.Product_DigitalProduct{DigitalProduct: p.digitalProduct},
		}, nil
	case "physical":
		return &pb.Product{
			Name:        p.name,
			Price:       p.price,
			ProductType: &pb.Product_PhysicalProduct{PhysicalProduct: p.physicalProduct},
		}, nil
	case "subscription":
		return &pb.Product{
			Name:        p.name,
			Price:       p.price,
			ProductType: &pb.Product_SubscriptionProduct{SubscriptionProduct: p.subscriptionProduct},
		}, nil
	default:
		return nil, fmt.Errorf("invalid product type: %s", p.productType)
	}
}

func TestGetProductIntegration(t *testing.T) {
    // Setup the database and services
    db := setupTestDatabase(t)
    _, handler := setupServiceAndHandler(t, db)

    /** 
	* The productID you want to fetch 
	* NOTE:  Remember to change ID "50ea81af-8b78-4280-be6e-c2eda87b8059" to any product ID in the database.
	*/	
	
    productID := "50ea81af-8b78-4280-be6e-c2eda87b8059"

    t.Run("Test get product by ID", func(t *testing.T) {
        // Call the handler to get the product by ID
        getResp, err := handler.GetProduct(context.Background(), &pb.GetProductRequest{Id: productID})
        if err != nil {
            t.Fatalf("Failed to fetch product: %v", err)
        }

        // Log the fetched product data
        t.Logf("Fetched Product: ID=%s, Name=%s, Price=%.2f", getResp.GetId(), getResp.GetName(), getResp.GetPrice())

        // Add assertions to check that the product was correctly fetched
        assert.Equal(t, productID, getResp.GetId())
        assert.NotEmpty(t, getResp.GetName())
        assert.NotZero(t, getResp.GetPrice())
    })
}

func TestUpdateProductIntegration(t *testing.T) {
    // Setup the database and services
    db := setupTestDatabase(t)
    _, handler := setupServiceAndHandler(t, db)

    // The productID you want to update
    productID := "50ea81af-8b78-4280-be6e-c2eda87b8059"

    // New product details to update
    updatedProduct := &pb.Product{
        Id:          productID,
        Name:        "Updated Product Name",
        Description: "Updated Product Description",
        Price:       99.99, // New price
    }

    t.Run("Test update product by ID", func(t *testing.T) {
        // Call the handler to update the product by ID
        updateResp, err := handler.UpdateProduct(context.Background(), updatedProduct)
        if err != nil {
            t.Fatalf("Failed to update product: %v", err)
        }

        // Log the updated product data
        t.Logf("Updated Product: ID=%s, Name=%s, Price=%.2f", updateResp.GetId(), updateResp.GetName(), updateResp.GetPrice())

        // Add assertions to check that the product was correctly updated
        assert.Equal(t, productID, updateResp.GetId())
        assert.Equal(t, "Updated Product Name", updateResp.GetName())
        assert.Equal(t, "Updated Product Description", updateResp.GetDescription())
        assert.Equal(t, float32(99.99), updateResp.GetPrice()) // Ensure both are float32

        // Fetch the product again to verify changes were applied
        getResp, err := handler.GetProduct(context.Background(), &pb.GetProductRequest{Id: productID})
        if err != nil {
            t.Fatalf("Failed to fetch updated product: %v", err)
        }

        // Assert that the updated values are present
        assert.Equal(t, updatedProduct.GetName(), getResp.GetName())
        assert.Equal(t, updatedProduct.GetDescription(), getResp.GetDescription())
        assert.Equal(t, float32(updatedProduct.GetPrice()), getResp.GetPrice()) // Ensure both are float32
    })
}

func TestDeleteProductIntegration(t *testing.T) {
    // Setup the database and services
    db := setupTestDatabase(t)
    _, handler := setupServiceAndHandler(t, db)

    // The productID you want to delete
    productID := "729998b3-c7cf-4c62-97b6-7f39129f2664"

    t.Run("Test delete product by ID", func(t *testing.T) {
        // Call the handler to delete the product by ID
        _, err := handler.DeleteProduct(context.Background(), &pb.DeleteProductRequest{Id: productID})
        if err != nil {
            t.Fatalf("Failed to delete product: %v", err)
        }

        // Try to fetch the product again to ensure it was deleted
        _, err = handler.GetProduct(context.Background(), &pb.GetProductRequest{Id: productID})
        assert.Error(t, err, "Expected error when fetching deleted product")
    })
}

func TestListProductsByType(t *testing.T) {
    // Create a mock repository
    mockRepo := new(MockProductRepository)

    // Prepare test data for "digital" products
    products := []domain.Product{
        {
            ID:          uuid.New(),
            Name:        "Digital Product 1",
            Description: "Digital Description 1",
            Price:       100.0,
            CreatedAt:   time.Now(),
            UpdatedAt:   time.Now(),
            DigitalProduct: &domain.DigitalProduct{
                FileSize:     100,
                DownloadLink: "http://example.com/download1",
            },
        },
    }

    // Log the prepared test data
    t.Logf("Prepared products: %+v", products)

    // Define the expected response
    expectedResponse := &pb.ListProductsResponse{
        Products: []*pb.Product{
            {
                Id:          products[0].ID.String(),
                Name:        "Digital Product 1",
                Description: "Digital Description 1",
                Price:       100.0,
                CreatedAt:   timestamppb.New(products[0].CreatedAt),
                UpdatedAt:   timestamppb.New(products[0].UpdatedAt),
                ProductType: &pb.Product_DigitalProduct{
                    DigitalProduct: &pb.DigitalProduct{
                        FileSize:     100,
                        DownloadLink: "http://example.com/download1",
                    },
                },
            },
        },
    }

    // Log the expected response
    t.Logf("Expected response: %+v", expectedResponse)

    // Set up mock expectations for "digital" products
    mockRepo.On("GetDigitalProducts", mock.Anything).Return(products, nil)
    
    // Set up mock expectation for Create (even though it's not needed for this test)
    mockRepo.On("Create", mock.Anything).Return(nil) 

    // Create the service instance using the constructor
    productService := service.NewProductService(mockRepo)

    // Create the request for listing "digital" products
    req := &pb.ListProductsRequest{
        Type: "digital",
    }

    // Log the request details
    t.Logf("Request: %+v", req)

    // Call the service method
    resp, err := productService.ListProducts(context.Background(), req)

    // Log the response
    t.Logf("Response: %+v", resp)

    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, expectedResponse.GetProducts()[0].GetId(), resp.GetProducts()[0].GetId())
    assert.Equal(t, expectedResponse.GetProducts()[0].GetName(), resp.GetProducts()[0].GetName())
    assert.Equal(t, expectedResponse.GetProducts()[0].GetDescription(), resp.GetProducts()[0].GetDescription())
    assert.Equal(t, expectedResponse.GetProducts()[0].GetPrice(), resp.GetProducts()[0].GetPrice())
    assert.Equal(t, expectedResponse.GetProducts()[0].GetCreatedAt().AsTime(), resp.GetProducts()[0].GetCreatedAt().AsTime())
    assert.Equal(t, expectedResponse.GetProducts()[0].GetUpdatedAt().AsTime(), resp.GetProducts()[0].GetUpdatedAt().AsTime())
}


