package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"

	"github.com/amrshaban2005/common"
	"github.com/amrshaban2005/common/discovery"
	"github.com/amrshaban2005/common/discovery/consul"
	"github.com/amrshaban2005/oms-gateway/gateway"
)

var (
	serviceName = "gateway"
	httpAddr    = common.EnvString("HTTP_ADDR", "localhost:8080")
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
)

func main() {
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, httpAddr); err != nil {
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

	ordersGateway := gateway.NewGRPCGateway(registry)

	handler := NewHandler(ordersGateway)

	r := mux.NewRouter()

	handler.reigsterRoutes(r)

	log.Printf("Starting http server at %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, r); err != nil {
		log.Fatal("Failed to start http server")
	}

}
