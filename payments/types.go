package main

import (
	"context"

	pb "github.com/sikozonpc/commons/api"
)

type PaymentsService interface {
	CreatePayment(context.Context, *pb.Order) (string, error)
}
