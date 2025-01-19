package service

import (
	"context"
	"errors"
	"product-microservice/internal/domain"
	"product-microservice/internal/repository"
	"github.com/google/uuid"
)

// SubscriptionService defines the interface for subscription-related business logic
type SubscriptionService interface {
	CreateSubscriptionPlan(ctx context.Context, productID uuid.UUID, planName string, duration int, price float64) (*domain.SubscriptionPlan, error)
	GetSubscriptionPlanByID(ctx context.Context, id uuid.UUID) (*domain.SubscriptionPlan, error)
	ListSubscriptionPlans(ctx context.Context) ([]*domain.SubscriptionPlan, error)
	DeleteSubscriptionPlan(ctx context.Context, id uuid.UUID) error
	UpdateSubscriptionPlan(ctx context.Context, id uuid.UUID, planName string, price float64, durationDays int) (*domain.SubscriptionPlan, error)
}

// subscriptionService is the implementation of SubscriptionService
type subscriptionService struct {
	repo repository.SubscriptionRepository
}

// NewSubscriptionService creates a new SubscriptionService
func NewSubscriptionService(repo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{repo: repo}
}

// CreateSubscriptionPlan creates a new subscription plan
func (s *subscriptionService) CreateSubscriptionPlan(ctx context.Context, productID uuid.UUID, planName string, duration int, price float64) (*domain.SubscriptionPlan, error) {
	if planName == "" {
		return nil, errors.New("subscription plan name cannot be empty")
	}

	if duration <= 0 {
		return nil, errors.New("subscription plan duration must be greater than zero")
	}

	if price <= 0 {
		return nil, errors.New("subscription plan price must be greater than zero")
	}

	plan := &domain.SubscriptionPlan{
		ID:        uuid.New(),
		ProductID: productID,
		PlanName:  planName,
		Duration:  duration,
		Price:     price,
	}

	// Save the plan in the repository
	return s.repo.Save(ctx, plan)
}

// GetSubscriptionPlanByID retrieves a subscription plan by its ID
func (s *subscriptionService) GetSubscriptionPlanByID(ctx context.Context, id uuid.UUID) (*domain.SubscriptionPlan, error) {
	return s.repo.FindByID(ctx, id)
}

// ListSubscriptionPlans fetches all subscription plans from the repository
func (s *subscriptionService) ListSubscriptionPlans(ctx context.Context) ([]*domain.SubscriptionPlan, error) {
	// Fetch all subscription plans without any conditions
	plans, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	return plans, nil
}

// DeleteSubscriptionPlan deletes a subscription plan by its ID
func (s *subscriptionService) DeleteSubscriptionPlan(ctx context.Context, id uuid.UUID) error {
	// Delegate to the repository to delete the subscription
	return s.repo.Delete(ctx, id)
}

// UpdateSubscriptionPlan updates a subscription plan by its ID
func (s *subscriptionService) UpdateSubscriptionPlan(ctx context.Context, id uuid.UUID, planName string, price float64, durationDays int) (*domain.SubscriptionPlan, error) {
    // Find the subscription plan by ID
    subscription, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // Update the subscription's fields
    subscription.PlanName = planName
    subscription.Price = price
    subscription.Duration = durationDays

    // Save the updated subscription plan
    if err := s.repo.Update(ctx, subscription); err != nil {
        return nil, err
    }

    return subscription, nil
}

