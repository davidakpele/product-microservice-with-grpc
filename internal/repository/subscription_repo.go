package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"product-microservice/internal/domain"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Save(ctx context.Context, plan *domain.SubscriptionPlan) (*domain.SubscriptionPlan, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.SubscriptionPlan, error)
	FindByProductID(ctx context.Context, productID uuid.UUID) ([]*domain.SubscriptionPlan, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, subscription *domain.SubscriptionPlan) error
	ListAll(ctx context.Context) ([]*domain.SubscriptionPlan, error)
}

// subscriptionRepository implements SubscriptionRepository interface
type subscriptionRepository struct {
	db *gorm.DB
}

// NewSubscriptionRepository creates a new subscription repository
func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

// Save inserts a new subscription plan into the database
func (r *subscriptionRepository) Save(ctx context.Context, plan *domain.SubscriptionPlan) (*domain.SubscriptionPlan, error) {
	// Use GORM's Create method to insert the new record
	if err := r.db.WithContext(ctx).Create(plan).Error; err != nil {
		return nil, err
	}
	return plan, nil
}

// FindByID retrieves a subscription plan by its ID
func (r *subscriptionRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.SubscriptionPlan, error) {
	plan := &domain.SubscriptionPlan{}
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(plan).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("subscription plan not found")
		}
		return nil, err
	}
	return plan, nil
}

// FindByProductID retrieves all subscription plans for a specific product
func (r *subscriptionRepository) FindByProductID(ctx context.Context, productID uuid.UUID) ([]*domain.SubscriptionPlan, error) {
	var plans []*domain.SubscriptionPlan
	if err := r.db.WithContext(ctx).Where("product_id = ?", productID).Find(&plans).Error; err != nil {
		return nil, err
	}
	return plans, nil
}

// Delete removes a subscription plan from the database by its ID
func (r *subscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.SubscriptionPlan{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("subscription plan not found")
		}
		return err
	}
	return nil
}

// Update updates an existing subscription plan in the database
func (r *subscriptionRepository) Update(ctx context.Context, subscription *domain.SubscriptionPlan) error {
	// Update the subscription plan in the database
	if err := r.db.Save(subscription).Error; err != nil {
		return err
	}
	return nil
}

// ListAll fetches all subscription plans from the database
func (r *subscriptionRepository) ListAll(ctx context.Context) ([]*domain.SubscriptionPlan, error) {
	var plans []*domain.SubscriptionPlan

	// Query the database for all subscription plans without conditions
	if err := r.db.Find(&plans).Error; err != nil {
		return nil, err
	}

	return plans, nil
}
