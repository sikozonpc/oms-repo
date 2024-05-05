package main

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	pb "github.com/sikozonpc/commons/api"
	"google.golang.org/grpc"
)

type StockGrpcHandler struct {
	pb.UnimplementedStockServiceServer

	service StockService
	channel *amqp.Channel
}

func NewGRPCHandler(
	server *grpc.Server,
	channel *amqp.Channel,
	stockService StockService,
) {
	handler := &StockGrpcHandler{
		service: stockService,
		channel: channel,
	}

	pb.RegisterStockServiceServer(server, handler)
}

func (s *StockGrpcHandler) CheckIfItemIsInStock(ctx context.Context, p *pb.CheckIfItemIsInStockRequest) (*pb.CheckIfItemIsInStockResponse, error) {
	inStock, items, err := s.service.CheckIfItemAreInStock(ctx, p.Items)
	if err != nil {
		return nil, err
	}

	return &pb.CheckIfItemIsInStockResponse{
		InStock: inStock,
		Items:   items,
	}, nil
}

func (s *StockGrpcHandler) GetItems(ctx context.Context, payload *pb.GetItemsRequest) (*pb.GetItemsResponse, error) {
	items, err := s.service.GetItems(ctx, payload.ItemIDs)
	if err != nil {
		return nil, err
	}

	return &pb.GetItemsResponse{
		Items: items,
	}, nil
}
