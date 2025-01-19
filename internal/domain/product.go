package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID                  uuid.UUID `gorm:"primaryKey"`
	Name                string
	Description         string
	Price               float64
	CreatedAt           time.Time
	UpdatedAt           time.Time
	// Associations with specific product types
	DigitalProductID    *uuid.UUID `gorm:"index"`
	PhysicalProductID   *uuid.UUID `gorm:"index"`
	SubscriptionProductID *uuid.UUID `gorm:"index"`

	// Actual product associations (Pointers to structs)
	DigitalProduct      *DigitalProduct
	PhysicalProduct     *PhysicalProduct
	SubscriptionProduct *SubscriptionProduct
}

type DigitalProduct struct {
	ID           uuid.UUID `gorm:"primaryKey"`
	FileSize     int32
	DownloadLink string
}

type PhysicalProduct struct {
	ID         uuid.UUID `gorm:"primaryKey"`
	Weight     float32
	Dimensions string
}

type SubscriptionProduct struct {
	ID                uuid.UUID `gorm:"primaryKey"`
	SubscriptionPeriod string
	RenewalPrice      float32
}

// Hook to automatically set UUID before creating records
func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	if p.DigitalProduct != nil {
		p.DigitalProduct.ID = uuid.New()
	}
	if p.PhysicalProduct != nil {
		p.PhysicalProduct.ID = uuid.New()
	}
	if p.SubscriptionProduct != nil {
		p.SubscriptionProduct.ID = uuid.New()
	}
	return nil
}

func (dp *DigitalProduct) BeforeCreate(tx *gorm.DB) (err error) {
	if dp.ID == uuid.Nil {
		dp.ID = uuid.New()
	}
	return nil
}

func (pp *PhysicalProduct) BeforeCreate(tx *gorm.DB) (err error) {
	if pp.ID == uuid.Nil {
		pp.ID = uuid.New()
	}
	return nil
}

func (sp *SubscriptionProduct) BeforeCreate(tx *gorm.DB) (err error) {
	if sp.ID == uuid.Nil {
		sp.ID = uuid.New()
	}
	return nil
}
