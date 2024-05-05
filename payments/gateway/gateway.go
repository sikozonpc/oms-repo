package gateway

import "context"

type OrdersGateway interface {
	UpdateOrderAfterPaymentLink(ctx context.Context, orderID, paymentLink string) error
}