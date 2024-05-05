package inmem

import pb "github.com/sikozonpc/commons/api"

type Inmem struct {}

func NewInmem() *Inmem {
	return &Inmem{}
}

func (i *Inmem) CreatePaymentLink(*pb.Order) (string, error) {
	return "dummy-link", nil
}