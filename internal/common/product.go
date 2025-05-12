package common

import (
	"fmt"
	"strings"
	"time"
)

type ProductType string

const (
	ProductTgStars   ProductType = "tg_stars"
	ProductTgPremium ProductType = "tg_premium"
)

type Product struct {
	ID        int64       `json:"id"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt *time.Time  `json:"updatedAt,omitempty"`
	DeletedAt *time.Time  `json:"deletedAt,omitempty"`
	Type      ProductType `json:"type"`
	Name      string      `json:"name"`
	Price     int64       `json:"price"`
	Amount    int         `json:"amount"`
	Badge     string      `json:"badge"`
}

func (p *Product) NeedUsername() bool {
	switch p.Type {
	case ProductTgStars, ProductTgPremium:
		return true
	default:
		return false
	}
}

func (p *Product) Valid() error {
	switch p.Type {
	case ProductTgStars:
		if p.Amount > 0 && p.Amount < 50 {
			return fmt.Errorf("minimum tg stars count 50")
		}
	case ProductTgPremium:
		switch p.Amount {
		case 0, 3, 6, 12:
		default:
			return fmt.Errorf("invalid tg premium months: %d", p.Amount)
		}
	default:
		return fmt.Errorf("invalid product type: %s", p.Type)
	}

	if p.Amount < 0 {
		return fmt.Errorf("invalid product count")
	}
	if p.Price < 0 {
		return fmt.Errorf("invalid product price")
	}

	p.Name = strings.TrimSpace(p.Name)
	if p.Name == "" {
		return fmt.Errorf("product name is empty")
	}

	return nil
}
