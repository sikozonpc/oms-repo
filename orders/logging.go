package main

import (
	"context"
	"time"

	pb "github.com/sikozonpc/commons/api"
	"go.uber.org/zap"
)

type LoggingMiddleware struct {
	next OrdersService
}

func NewLoggingMiddleware(next OrdersService) OrdersService {
	return &LoggingMiddleware{next}
}

func (s *LoggingMiddleware) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.Order, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("GetOrder", zap.Duration("took", time.Since(start)))
	}()

	return s.next.GetOrder(ctx, p)
}

func (s *LoggingMiddleware) UpdateOrder(ctx context.Context, o *pb.Order) (*pb.Order, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("UpdateOrder", zap.Duration("took", time.Since(start)))
	}()

	return s.next.UpdateOrder(ctx, o)
}

func (s *LoggingMiddleware) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest, items []*pb.Item) (*pb.Order, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("CreateOrder", zap.Duration("took", time.Since(start)))
	}()

	return s.next.CreateOrder(ctx, p, items)
}

func (s *LoggingMiddleware) ValidateOrder(ctx context.Context, p *pb.CreateOrderRequest) ([]*pb.Item, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("ValidateOrder", zap.Duration("took", time.Since(start)))
	}()

	return s.next.ValidateOrder(ctx, p)
}
