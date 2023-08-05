package main

import (
	"fmt"
	"l0/config"
	"l0/models"
	"l0/nats"
	"os"

	"github.com/nats-io/stan.go"
)

func fatalError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func main() {
	// configs
	cfg, err := config.GetConfig()
	if err != nil {
		fatalError(err)
	}

	// cache and db
	orderModel, err := models.MakeCachedOrderModel(cfg.DB)
	if err != nil {
		fatalError(err)
	}
	defer orderModel.Close()

	// nats
	sc, err := stan.Connect(cfg.NATS.ClusterID, cfg.NATS.ClientID)
	if err != nil {
		fatalError(err)
	}
	sub, err := sc.Subscribe("orders", nats.MakeOrderHandler(orderModel))
	if err != nil {
		fatalError(err)
	}
	defer sc.Close()
	defer sub.Unsubscribe()

	// http

	// system interrupts
	for {

	}
}
