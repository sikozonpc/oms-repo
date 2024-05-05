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

func NewGateway(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

func (g *Gateway) CheckIfItemIsInStock(ctx context.Context, customerID string, items []*pb.ItemsWithQuantity) (bool, []*pb.Item, error) {
	conn, err := discovery.ServiceConnection(context.Background(), "stock", g.registry)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	c := pb.NewStockServiceClient(conn)

	res, err := c.CheckIfItemIsInStock(ctx, &pb.CheckIfItemIsInStockRequest{
		Items: items,
	})

	return res.InStock, res.Items, err
}
