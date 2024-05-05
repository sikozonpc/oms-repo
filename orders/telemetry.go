package main

import (
	"context"
	"fmt"

	pb "github.com/sikozonpc/commons/api"
	"go.opentelemetry.io/otel/trace"
)

type TelemetryMiddleware struct {
	next OrdersService
}

func NewTelemetryMiddleware(next OrdersService) OrdersService {
	return &TelemetryMiddleware{next}
}

func (s *TelemetryMiddleware) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.Order, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("GetOrder: %v", p))

	return s.next.GetOrder(ctx, p)
}

func (s *TelemetryMiddleware) UpdateOrder(ctx context.Context, o *pb.Order) (*pb.Order, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("UpdateOrder: %v", o))

	return s.next.UpdateOrder(ctx, o)
}

func (s *TelemetryMiddleware) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest, items []*pb.Item) (*pb.Order, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("CreateOrder: %v, items: %v", p, items))

	return s.next.CreateOrder(ctx, p, items)
}

func (s *TelemetryMiddleware) ValidateOrder(ctx context.Context, p *pb.CreateOrderRequest) ([]*pb.Item, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("ValidateOrder: %v", p))

	return s.next.ValidateOrder(ctx, p)
}
