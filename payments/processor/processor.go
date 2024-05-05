package processor

import pb "github.com/sikozonpc/commons/api"

type PaymentProcessor interface {
	CreatePaymentLink(*pb.Order) (string, error)
}
