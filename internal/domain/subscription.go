package domain

import (
	"github.com/google/uuid"
)

type SubscriptionPlan struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	ProductID uuid.UUID `gorm:"product_id"`
	PlanName  string    `json:"plan_name"`
	Duration  int       `json:"duration"`
	Price     float64   `json:"price"`
}
