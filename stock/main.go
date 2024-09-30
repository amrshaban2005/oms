package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/amrshaban2005/common"
	"github.com/amrshaban2005/common/discovery"
	"github.com/amrshaban2005/common/discovery/consul"
	"google.golang.org/grpc"
)

var (
	serviceName = "stock"
	grpcAddr    = common.EnvString("GRPC_ADDR", "localhost:2002")
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
)

func main() {

	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("Failed to health check")
			}
			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	grpcServer := grpc.NewServer()
	l, err := net.Listen("tcp", grpcAddr)

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer l.Close()

	store := NewStore()
	svc := NewSerive(store)

	NewGrpcHandler(grpcServer, svc)

	log.Println("GRPC server started at ", grpcAddr)

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}

}
