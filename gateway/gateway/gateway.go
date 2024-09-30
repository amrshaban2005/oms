package gateway

import (
	"context"

	pb "github.com/amrshaban2005/common/api"
)

type OrdersGateway interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)	
}
