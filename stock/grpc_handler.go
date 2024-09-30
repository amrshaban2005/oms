package main

import (
	"context"

	pb "github.com/amrshaban2005/common/api"
	"google.golang.org/grpc"
)

type StockGrpcHandler struct {
	pb.UnimplementedStockServiceServer

	service StockService
}

func NewGrpcHandler(server *grpc.Server, stockService StockService) {
	handler := &StockGrpcHandler{
		service: stockService,
	}
	pb.RegisterStockServiceServer(server, handler)
}

func (s *StockGrpcHandler) CheckIfItemIsInStock(ctx context.Context, p *pb.CheckIfItemIsInStockRequest) (*pb.CheckIfItemIsInStockResponse, error) {
	inStock, items, err := s.service.checkIfItemsAreInStock(ctx, p.Items)
	if err != nil {
		return nil, err
	}

	return &pb.CheckIfItemIsInStockResponse{InStock: inStock, Items: items}, nil
}

func (s *StockGrpcHandler) GetItems(ctx context.Context, payload *pb.GetItemsRequest) (*pb.GetItemsResponse, error) {
	items, err := s.service.GetItems(ctx, payload.ItemIDs)
	if err != nil {
		return nil, err
	}

	return &pb.GetItemsResponse{
		Items: items,
	}, nil
}
