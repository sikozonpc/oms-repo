package main

import (
	"context"

	pb "github.com/sikozonpc/commons/api"
)

type Service struct {
	store StockStore
}

func NewService(store StockStore) *Service {
	return &Service{store}
}

func (s *Service) CheckIfItemAreInStock(ctx context.Context, p []*pb.ItemsWithQuantity) (bool, []*pb.Item, error) {
	itemIDs := make([]string, 0)
	for _, item := range p {
		itemIDs = append(itemIDs, item.ID)
	}

	itemsInStock, err := s.store.GetItems(ctx, itemIDs)
	if err != nil {
		return false, nil, err
	}

	// Check if all items are in stock
	for _, stockItem := range itemsInStock {
		for _, reqItem := range p {
			if stockItem.ID == reqItem.ID && stockItem.Quantity < reqItem.Quantity {
				return false, itemsInStock, nil
			}
		}
	}

	// create items with prices from stock
	items := make([]*pb.Item, 0)
	for _, stockItem := range itemsInStock {
		for _, reqItem := range p {
			if stockItem.ID == reqItem.ID {
				items = append(items, &pb.Item{
					ID:       stockItem.ID,
					Name:     stockItem.Name,
					PriceID:  stockItem.PriceID,
					Quantity: reqItem.Quantity,
				})
			}
		}
	}

	return true, items, nil
}

func (s *Service) GetItems(ctx context.Context, ids []string) ([]*pb.Item, error) {
	return s.store.GetItems(ctx, ids)
}
