package main

import (
	"context"

	pb "github.com/sikozonpc/commons/api"
	"github.com/sikozonpc/omsv2-payments/gateway"
	"github.com/sikozonpc/omsv2-payments/processor"
)

type service struct {
	processor processor.PaymentProcessor
	gateway   gateway.OrdersGateway
}

func NewService(processor processor.PaymentProcessor, gateway gateway.OrdersGateway) *service {
	return &service{processor, gateway}
}

func (s *service) CreatePayment(ctx context.Context, o *pb.Order) (string, error) {
	link, err := s.processor.CreatePaymentLink(o)
	if err != nil {
		return "", err
	}

	// update order with the link
	err = s.gateway.UpdateOrderAfterPaymentLink(ctx, o.ID, link)
	if err != nil {
		return "", err
	}

	return link, nil
}
