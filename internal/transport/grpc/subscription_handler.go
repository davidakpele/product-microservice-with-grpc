package grpc

import (
	"context"
	"log"
	"math"
	"product-microservice/internal/service"
	pb "product-microservice/proto/subscription"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// SubscriptionHandler implements the gRPC service methods
// SubscriptionHandler implements the gRPC service methods
type SubscriptionHandler struct {
	subscriptionService service.SubscriptionService
	productService      service.ProductService
	pb.UnimplementedSubscriptionServiceServer
}

// NewSubscriptionHandler creates a new SubscriptionHandler
func NewSubscriptionHandler(subscriptionService service.SubscriptionService, productService service.ProductService) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
		productService:      productService,
	}
}

// CreateSubscription handles the gRPC request to create a subscription plan
func (h *SubscriptionHandler) CreateSubscription(ctx context.Context, req *pb.CreateSubscriptionPlanRequest) (*pb.CreateSubscriptionPlanResponse, error) {
	// Fetch product by name
	existingProduct, err := h.productService.FindProductById(ctx, req.GetProductId())
	if err != nil {
		log.Printf("Failed to fetch product by name: %v", err)
		return nil, status.Errorf(codes.NotFound, "Product not found: %v", err)
	}

	if existingProduct == nil {
		log.Printf("No product found with the given name: %s", req.GetProductId())
		return nil, status.Errorf(codes.NotFound, "Product not found")
	}

	// Convert DurationDays (int32) to int
	duration := int(req.GetDurationDays())

	// Create a new subscription plan via service layer
	plan, err := h.subscriptionService.CreateSubscriptionPlan(ctx, existingProduct.ID, req.GetPlanName(), duration, float64(req.GetPrice()))
	if err != nil {
		log.Printf("Failed to create subscription plan: %v", err)
		return nil, err
	}

	// Return the created plan as part of the response
	return &pb.CreateSubscriptionPlanResponse{
		SubscriptionPlan: &pb.SubscriptionPlan{
			Id:           plan.ID.String(),
			ProductId:    plan.ProductID.String(),
			PlanName:     plan.PlanName,
			Price:        float32(plan.Price),
			DurationDays: int32(plan.Duration),
		},
	}, nil
}

// GetSubscriptionPlan handles the gRPC request to fetch a subscription plan by its ID
func (h *SubscriptionHandler) GetSubscriptionPlan(ctx context.Context, req *pb.GetSubscriptionPlanRequest) (*pb.SubscriptionPlan, error) {
	// Parse Subscription ID
	subscriptionID, err := uuid.Parse(req.GetId())
	if err != nil {
		log.Printf("Invalid Subscription UUID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Subscription ID")
	}

	// Fetch subscription plan from the repository
	subscriptionPlan, err := h.subscriptionService.GetSubscriptionPlanByID(ctx, subscriptionID)
	if err != nil {
		log.Printf("Failed to fetch subscription plan: %v", err)
		return nil, status.Errorf(codes.NotFound, "Subscription plan not found")
	}

	// Return the fetched subscription plan in the response
	return &pb.SubscriptionPlan{
		Id:           subscriptionPlan.ID.String(),
		ProductId:    subscriptionPlan.ProductID.String(),
		PlanName:     subscriptionPlan.PlanName,
		Price:        float32(subscriptionPlan.Price),
		DurationDays: int32(subscriptionPlan.Duration),
	}, nil
}

// ListSubscriptionPlans handles the gRPC request to list all subscription plans
func (h *SubscriptionHandler) ListSubscriptionPlans(ctx context.Context, req *pb.ListSubscriptionPlansRequest) (*pb.ListSubscriptionPlansResponse, error) {
	// Fetch all subscription plans
	subscriptionPlans, err := h.subscriptionService.ListSubscriptionPlans(ctx)
	if err != nil {
		log.Printf("Failed to list subscription plans: %v", err)
		return nil, status.Errorf(codes.Internal, "Error fetching subscription plans")
	}

	// Map the subscription plans to the protobuf response format
	var pbSubscriptionPlans []*pb.SubscriptionPlan
	for _, plan := range subscriptionPlans {
		pbSubscriptionPlans = append(pbSubscriptionPlans, &pb.SubscriptionPlan{
			Id:           plan.ID.String(),
			ProductId:    plan.ProductID.String(),
			PlanName:     plan.PlanName,
			Price:        float32(plan.Price),
			DurationDays: int32(plan.Duration),
		})
	}

	// Return the response with all subscription plans
	return &pb.ListSubscriptionPlansResponse{
		SubscriptionPlans: pbSubscriptionPlans,
	}, nil
}


func (h *SubscriptionHandler) UpdateSubscriptionPlan(ctx context.Context, req *pb.UpdateSubscriptionPlanRequest) (*pb.SubscriptionPlan, error) {
    // Convert the string ID from the request to a UUID
    id, err := uuid.Parse(req.GetId())
    if err != nil {
        log.Printf("Invalid UUID: %v", err)
        return nil, err
    }

    // Get the price and durationDays directly from the request (no conversion)
    price := req.GetPrice() // price as float32
    durationDays := req.GetDurationDays() // durationDays as int32

    // Round the price to 2 decimal places
    roundedPrice := math.Round(float64(price)*100) / 100.0

    // Update the subscription plan via service layer
    updatedPlan, err := h.subscriptionService.UpdateSubscriptionPlan(ctx, id, req.GetPlanName(), roundedPrice, int(durationDays))
    if err != nil {
        log.Printf("Failed to update subscription plan: %v", err)
        return nil, err
    }

    // Return the updated plan as a response (no conversion needed for price or durationDays)
    return &pb.SubscriptionPlan{
        Id:           updatedPlan.ID.String(),
        ProductId:    updatedPlan.ProductID.String(),
        PlanName:     updatedPlan.PlanName,
        Price:        float32(updatedPlan.Price), // Return as float32
        DurationDays: int32(updatedPlan.Duration), // Return as int32
    }, nil
}



// DeleteSubscription handles the gRPC request to delete a subscription plan
func (h *SubscriptionHandler) DeleteSubscription(ctx context.Context, req *pb.DeleteSubscriptionPlanRequest) (*emptypb.Empty, error) {
	// Step 1: Convert the string ID from the request to a UUID
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		log.Printf("Invalid UUID format: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid UUID: %v", err)
	}

	// Step 2: Call the service to delete the subscription
	err = h.subscriptionService.DeleteSubscriptionPlan(ctx, id)
	if err != nil {
		log.Printf("Failed to delete subscription plan: %v", err)
		return nil, status.Errorf(codes.Internal, "Failed to delete subscription plan: %v", err)
	}

	// Step 3: Return an empty response after successful deletion
	return &emptypb.Empty{}, nil
}


// RegisterHandlers registers the SubscriptionHandler with the gRPC server
func RegisterHandler(server *grpc.Server, subscriptionService service.SubscriptionService, productService service.ProductService) {
	handler := NewSubscriptionHandler(subscriptionService, productService)
	pb.RegisterSubscriptionServiceServer(server, handler)
}
