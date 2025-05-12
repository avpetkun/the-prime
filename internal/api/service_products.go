package api

import (
	"context"
	"fmt"
	"time"

	"github.com/avpetkun/the-prime/internal/common"
	"github.com/avpetkun/the-prime/pkg/timeu"
)

func (s *Service) getProduct(productID int64) (*common.Product, error) {
	for _, p := range s.allProducts {
		if p.ID == productID {
			return p, nil
		}
	}
	return nil, fmt.Errorf("product %d not found", productID)
}

func (s *Service) startLoadProductsLoop(ctx context.Context) error {
	products, err := s.db.GetAllProducts(ctx)
	if err != nil {
		return err
	}
	s.allProducts = products

	go func() {
		for !timeu.SleepContext(ctx, time.Second*10) {
			products, err = s.db.GetAllProducts(ctx)
			if err != nil {
				s.log.Error().Err(err).Msg("failed to load products")
			} else {
				s.allProducts = products
			}
		}
	}()
	return nil
}
