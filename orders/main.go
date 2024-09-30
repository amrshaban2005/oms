package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/amrshaban2005/common"
	"github.com/amrshaban2005/common/discovery"
	"github.com/amrshaban2005/common/discovery/consul"
	"github.com/amrshaban2005/oms-orders/gateway"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"google.golang.org/grpc"
)

var (
	serviceName = "orders"
	grpcAddr    = common.EnvString("GRPC_ADDR", "localhost:2000")
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
	mongoUser   = common.EnvString("MONGO_DB_USER", "root")
	mongoPass   = common.EnvString("MONGO_DB_PASS", "example")
	mongoAddr   = common.EnvString("MONGO_DB_HOST", "localhost:27017")
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
				log.Fatal("failed to health check")
			}
			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	// mongo db conn  mongodb://localhost:27017
	uri := fmt.Sprintf("mongodb://%s", mongoAddr)
	mongoClient, err := connectToMongoDB(uri)
	if err != nil {
		log.Fatal("failed to connect to mongo db")
	}

	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {

		log.Fatal("failed to listen")
	}
	defer l.Close()

	gateway := gateway.NewGateway(registry)

	store := NewStore(mongoClient)
	svc := NewService(store, gateway)
	NewGRPCHandler(grpcServer, svc)

	log.Println("GRPC server started at ", grpcAddr)

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}

}

func connectToMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	return client, err
}
