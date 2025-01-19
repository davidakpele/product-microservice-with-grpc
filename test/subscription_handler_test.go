package test

import (
	"context"
	"fmt"
	"product-microservice/internal/domain"
	"product-microservice/internal/repository"
	"product-microservice/internal/service"
	"product-microservice/internal/transport/grpc"
	pb "product-microservice/proto/subscription"

	// "strconv"
	"testing"

	// "github.com/google/uuid"
	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SubscriptionService interface {
    CreateSubscription(ctx context.Context, userID, planID string) (*domain.SubscriptionPlan, error)
    CreateSubscriptionPlan(ctx context.Context, name string, price float64) (*domain.SubscriptionPlan, error)
}

// MockSubscriptionService mocks the SubscriptionService interface
type MockSubscriptionService struct {
    mock.Mock
}

func (m *MockSubscriptionService) CreateSubscription(ctx context.Context, userID string, planID string) (*domain.SubscriptionPlan, error) {
    args := m.Called(ctx, userID, planID)
    return args.Get(0).(*domain.SubscriptionPlan), args.Error(1)
}

func (m *MockSubscriptionService) CreateSubscriptionPlan(ctx context.Context, name string, price float64) (*domain.SubscriptionPlan, error) {
    args := m.Called(ctx, name, price)
    return args.Get(0).(*domain.SubscriptionPlan), args.Error(1)
}

func (m *MockSubscriptionService) GetSubscriptionPlanByID(ctx context.Context, id string) (*domain.SubscriptionPlan, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*domain.SubscriptionPlan), args.Error(1)
}

func (m *MockSubscriptionService) ListSubscriptionPlans(ctx context.Context) ([]*domain.SubscriptionPlan, error) {
    args := m.Called(ctx)
    return args.Get(0).([]*domain.SubscriptionPlan), args.Error(1)
}

func (m *MockSubscriptionService) UpdateSubscriptionPlan(ctx context.Context, id string, name string, price float64) (*domain.SubscriptionPlan, error) {
    args := m.Called(ctx, id, name, price)
    return args.Get(0).(*domain.SubscriptionPlan), args.Error(1)
}

func (m *MockSubscriptionService) DeleteSubscriptionPlan(ctx context.Context, id string) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}



// SubscriptionTestDatabaseSetUp sets up a PostgreSQL database connection for testing using the existing db and config setup
func SubscriptionTestDatabaseSetUp(t *testing.T) *gorm.DB {
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
	err = db.AutoMigrate(&domain.SubscriptionPlan{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

// TestCreateSubscription is the test for creating a subscription using a mock service
func TestCreateSubscriptionIntegration(t *testing.T) {
	// Set up the test database connection
	db := SubscriptionTestDatabaseSetUp(t)
	productRepo := repository.NewProductRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	productService := service.NewProductService(productRepo)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)
	handler := grpc.NewSubscriptionHandler(subscriptionService, productService)

	// Define multiple subscription plans
	subscriptionPlans := []struct {
		ProductId string
		name        string
		price       float64
		duration    int32
	}{
		// NOTE: the product id's here are from product table. Product need to exist before you can product create subscriptionPlan.
		{"e96a46ad-8cb8-4484-b8f4-d72face35f86", "Basic Plan", 9.99, 20},
		{"b64c7763-eaf7-4581-bef8-d4a13a317ad3", "Standard Plan", 19.99, 40},
	}

	for _, plan := range subscriptionPlans {
		req := &pb.CreateSubscriptionPlanRequest{
			ProductId:  plan.ProductId, 
			PlanName:     plan.name,
			Price:        float32(plan.price),
			DurationDays: plan.duration,
		}

		// Call handler
		resp, err := handler.CreateSubscription(context.Background(), req)

		// Assert no error occurred
		if err != nil {
			t.Fatalf("Error creating subscription plan: %v", err)
		}

		// Access the subscription plan from the response
		subscriptionPlan := resp

		// Print success message
		t.Logf("Successfully created subscription plan: %s with price: %.2f", subscriptionPlan.SubscriptionPlan.GetPlanName(), subscriptionPlan.SubscriptionPlan.GetPrice())
	}
}

// Test fetching a single subscription.
func TestGetSubscriptionIntegration(t *testing.T) {
	// Set up the test database connection
	db := SubscriptionTestDatabaseSetUp(t)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	service := service.NewSubscriptionService(subscriptionRepo)
	handler := grpc.NewSubscriptionHandler(service, nil) 

	// Assume a subscription plan already exists in the database
	existingSubscriptionID := "32e4182d-a8d6-4c10-9449-5df902cf3b53" 

	// Create a gRPC request to fetch the subscription plan
	req := &pb.GetSubscriptionPlanRequest{
		Id: existingSubscriptionID,
	}

	// Call the GetSubscriptionPlan handler
	resp, err := handler.GetSubscriptionPlan(context.Background(), req)

	// Assert no error occurred
	if err != nil {
		t.Fatalf("Failed to fetch subscription plan: %v", err)
	}

	// Access the subscription plan from the response
	subscriptionPlan := resp

	// Assert subscription plan is not nil
	if subscriptionPlan == nil {
		t.Fatalf("No subscription plan found with ID: %s", existingSubscriptionID)
	}

	// Log the subscription plan details
	t.Logf("Successfully fetched subscription plan: ID: %s, Name: %s, Price: %.2f, Duration: %d days",
		subscriptionPlan.GetId(),
		subscriptionPlan.GetPlanName(),
		subscriptionPlan.GetPrice(),
		subscriptionPlan.GetDurationDays(),
	)
}

// Test to list all subscriptions without creating new data
func TestListSubscriptionsIntegration(t *testing.T) {
	// Set up the test database connection
	db := SubscriptionTestDatabaseSetUp(t)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	service := service.NewSubscriptionService(subscriptionRepo)
	handler := grpc.NewSubscriptionHandler(service, nil) 

	// Create a gRPC request to list all subscription plans 
	req := &pb.ListSubscriptionPlansRequest{}

	// Call the ListSubscriptionPlans handler
	resp, err := handler.ListSubscriptionPlans(context.Background(), req)

	// Assert no error occurred
	if err != nil {
		t.Fatalf("Failed to list subscription plans: %v", err)
	}

	// Assert that there are subscription plans in the response
	if len(resp.GetSubscriptionPlans()) == 0 {
		t.Fatalf("No subscription plans found")
	}

	// Iterate over all subscription plans and log their details
	for _, subscriptionPlan := range resp.GetSubscriptionPlans() {
		t.Logf("Subscription Plan: ID: %s, Product ID: %s, Name: %s, Price: %.2f, Duration: %d days",
			subscriptionPlan.GetId(),
			subscriptionPlan.GetProductId(),
			subscriptionPlan.GetPlanName(),
			subscriptionPlan.GetPrice(),
			subscriptionPlan.GetDurationDays(),
		)
	}
}
func TestUpdateSubscriptionIntegration(t *testing.T) {
    // Set up test database and repository
    db := SubscriptionTestDatabaseSetUp(t)
    repo := repository.NewSubscriptionRepository(db)
    service := service.NewSubscriptionService(repo)
    handler := grpc.NewSubscriptionHandler(service, nil)

    // Step 1: Fetch the subscription with the provided ID from the database
    var subscription domain.SubscriptionPlan
    subscriptionID := "b57103bb-5170-4dff-a843-45e831535cc7"

    // Fetch subscription with the specific ID
    if err := db.First(&subscription, "id = ?", subscriptionID).Error; err != nil {
        t.Fatalf("Failed to fetch subscription: %v", err)
    }

    // Step 2: Prepare the update request with new values
    updatedPlanName := "Updated Plan Name"
    updatedPrice := 99.99
    updatedDurationDays := 30

    req := &pb.UpdateSubscriptionPlanRequest{
        Id:          subscription.ID.String(),
        PlanName:    updatedPlanName,
        Price:       float32(updatedPrice), // Ensuring price is passed as float32
        DurationDays: int32(updatedDurationDays),
    }

    // Step 3: Call the handler to update the subscription
    updatedSubscriptionResp, err := handler.UpdateSubscriptionPlan(context.Background(), req)
    if err != nil {
        t.Fatalf("Failed to update subscription: %v", err)
    }

    // Step 4: Verify the subscription is updated in the database
    var updatedSubscription domain.SubscriptionPlan
    if err := db.First(&updatedSubscription, "id = ?", subscriptionID).Error; err != nil {
        t.Fatalf("Failed to fetch updated subscription: %v", err)
    }

   // Assert that the updated values are reflected in the database
	assert.Equal(t, updatedPlanName, updatedSubscription.PlanName)
	assert.Equal(t, float32(updatedPrice), float32(updatedSubscription.Price))  // Cast to float32
	assert.Equal(t, int32(updatedDurationDays), int32(updatedSubscription.Duration))

	// Optionally, assert the response values as well (check that the updated values match)
	assert.Equal(t, updatedPlanName, updatedSubscriptionResp.GetPlanName())
	assert.Equal(t, float32(updatedPrice), updatedSubscriptionResp.GetPrice())  // Cast to float32
	assert.Equal(t, int32(updatedDurationDays), int32(updatedSubscriptionResp.GetDurationDays()))

}

// TestDeleteSubscriptionIntegration tests deleting a subscription
func TestDeleteSubscriptionIntegration(t *testing.T) {
    // Set up test database and repo
    db := SubscriptionTestDatabaseSetUp(t)
    repo := repository.NewSubscriptionRepository(db)
    service := service.NewSubscriptionService(repo)
    handler := grpc.NewSubscriptionHandler(service, nil)

    // Step 1: Fetch the subscription with the provided ID from the database
    var subscription domain.SubscriptionPlan
    subscriptionID := "32e4182d-a8d6-4c10-9449-5df902cf3b53" 

    // Fetch subscription with the specific ID
    if err := db.First(&subscription, "id = ?", subscriptionID).Error; err != nil {
        t.Fatalf("Failed to fetch subscription: %v", err)
    }

    // Step 2: Create a delete request for the subscription ID
    req := &pb.DeleteSubscriptionPlanRequest{Id: subscription.ID.String()} 

    // Step 3: Call the handler to delete the subscription
    _, err := handler.DeleteSubscription(context.Background(), req)
    if err != nil {
        t.Fatalf("Failed to delete subscription: %v", err)
    }

    // Step 4: Verify the subscription is deleted from the database
    var count int64
    if err := db.Model(&domain.SubscriptionPlan{}).Where("id = ?", subscription.ID).Count(&count).Error; err != nil {
        t.Fatalf("Failed to count subscriptions: %v", err)
    }

    // Assert no subscription with this ID exists
    assert.Equal(t, int64(0), count)
}



