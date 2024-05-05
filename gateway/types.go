package main

import pb "github.com/sikozonpc/commons/api"

type CreateOrderRequest struct {
	Order         *pb.Order `"json": order`
	RedirectToURL string    `"json": redirectToURL`
}
