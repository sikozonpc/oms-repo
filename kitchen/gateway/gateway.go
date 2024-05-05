package gateway

import (
	"context"

	pb "github.com/sikozonpc/commons/api"
)

type KitchenGateway interface {
	UpdateOrder(context.Context, *pb.Order) error
}
