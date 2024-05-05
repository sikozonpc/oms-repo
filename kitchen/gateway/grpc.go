package gateway

import (
	"context"
	"log"

	pb "github.com/sikozonpc/commons/api"
	"github.com/sikozonpc/commons/discovery"
)

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

func (g *Gateway) UpdateOrder(ctx context.Context, o *pb.Order) error {
	conn, err := discovery.ServiceConnection(context.Background(), "orders", g.registry)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	ordersClient := pb.NewOrderServiceClient(conn)

	_, err = ordersClient.UpdateOrder(ctx, o)
	return err
}
