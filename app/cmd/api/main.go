package main

import (
	"context"
	"fmt"

	"github.com/isac-caja/circuit-breaker-cases/app/client"
	"github.com/isac-caja/circuit-breaker-cases/app/config"
	"github.com/isac-caja/circuit-breaker-cases/app/http"
	"github.com/isac-caja/circuit-breaker-cases/app/service"
)

func main() {
	c, err := config.New()
	if err != nil {
		panic(fmt.Sprintf("Error to load configuration: %s", err.Error()))
	}

	maker := client.NewMaker(c.RequesterConfig.ClientConfigList)
	requester := service.NewServiceRequester(&c.RequesterConfig, maker.Make())
	service := service.New(c, requester)
	handler := http.NewHandler(service)
	router := http.NewRouter(handler)
	server := http.NewServer(router)

	defer server.Shutdown(context.Background())
	_ = server.ListenAndServe()
	fmt.Printf("Running server on port %s\n", server.Addr)
}
