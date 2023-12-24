package main

import (
	"log"
	"net/http"
	"wb-l0/common"

	"github.com/nats-io/nats.go"

	_ "github.com/lib/pq"
)

func main() {
	var cc common.CachedConnection
	defer cc.Close()

	cc.Connect()

	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}

	defer nc.Close()

	sub, err := nc.Subscribe("insert", cc.NATSInserter)

	if err != nil {
		log.Fatalf("Error subscribing to NATS subject: %v", err)
	}

	log.Printf("Server started!")

	defer sub.Unsubscribe()

	http.HandleFunc("/orders", cc.HTTPGetter)
	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}
