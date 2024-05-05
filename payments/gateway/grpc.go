package gateway

import (
	"context"
	"log"

	pb "github.com/sikozonpc/commons/api"
	"github.com/sikozonpc/commons/discovery"
)

type gateway struct {
	registry discovery.Registry
}

func NewGateway(registry discovery.Registry) *gateway {
	return &gateway{registry}
}

func (g *gateway) UpdateOrderAfterPaymentLink(ctx context.Context, orderID, paymentLink string) error {
	conn, err := discovery.ServiceConnection(context.Background(), "orders", g.registry)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	ordersClient := pb.NewOrderServiceClient(conn)

	_, err = ordersClient.UpdateOrder(ctx, &pb.Order{
		ID:          orderID,
		Status:      "waiting_payment",
		PaymentLink: paymentLink,
	})
	return err
}
