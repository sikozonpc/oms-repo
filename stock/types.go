package main

import (
	"context"

	pb "github.com/sikozonpc/commons/api"
)

type StockService interface {
	CheckIfItemAreInStock(context.Context, []*pb.ItemsWithQuantity) (bool, []*pb.Item, error)
	GetItems(ctx context.Context, ids []string) ([]*pb.Item, error)
}

type StockStore interface {
	GetItem(ctx context.Context, id string) (*pb.Item, error)
	GetItems(ctx context.Context, ids []string) ([]*pb.Item, error)
}
