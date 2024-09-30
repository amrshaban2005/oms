package main

import (
	"context"

	pb "github.com/amrshaban2005/common/api"
)

type StockService interface {
	checkIfItemsAreInStock(context.Context, []*pb.ItemsWithQuantity) (bool, []*pb.Item, error)
	GetItems(context.Context, []string) ([]*pb.Item, error)
}

type StockStore interface {
	GetItem(ctx context.Context, id string) (*pb.Item, error)
	GetItems(ctx context.Context, ids []string) ([]*pb.Item, error)
}
