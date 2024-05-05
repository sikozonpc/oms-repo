package main

import (
	"context"
	"fmt"

	pb "github.com/sikozonpc/commons/api"
	"go.opentelemetry.io/otel/trace"
)

type TelemetryMiddleware struct {
	next StockService
}

func NewTelemetryMiddleware(next StockService) StockService {
	return &TelemetryMiddleware{next}
}

func (s *TelemetryMiddleware) GetItems(ctx context.Context, ids []string) ([]*pb.Item, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("GetItems: %v", ids))

	return s.next.GetItems(ctx, ids)
}

func (s *TelemetryMiddleware) CheckIfItemAreInStock(ctx context.Context, p []*pb.ItemsWithQuantity) (bool, []*pb.Item, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("CheckIfItemAreInStock: %v", p))

	return s.next.CheckIfItemAreInStock(ctx, p)
}
